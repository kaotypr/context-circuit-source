package workspace

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"
)

func TestParseStatusReadsRenamesAsOneEntry(t *testing.T) {
	for _, tc := range []struct {
		name string
		text string
		want []StatusEntry
	}{
		{"empty", "", nil},
		{"untracked and modified", "?? new.md\x00 M old.md\x00", []StatusEntry{{Code: "??", Path: "new.md"}, {Code: " M", Path: "old.md"}}},
		{"deletion", " D gone.md\x00D  staged.md\x00", []StatusEntry{{Code: " D", Path: "gone.md"}, {Code: "D ", Path: "staged.md"}}},
		{"rename", "R  docs.md\x00README.md\x00 M other.md\x00", []StatusEntry{{Code: "R ", Path: "docs.md", From: "README.md"}, {Code: " M", Path: "other.md"}}},
		{"copy", "C  copy.md\x00source.md\x00", []StatusEntry{{Code: "C ", Path: "copy.md", From: "source.md"}}},
		{"rename then modified", "RM b.md\x00a.md\x00", []StatusEntry{{Code: "RM", Path: "b.md", From: "a.md"}}},
		{"spaces and arrows", "?? with space.md\x00R  a -> b.md\x00c d.md\x00", []StatusEntry{{Code: "??", Path: "with space.md"}, {Code: "R ", Path: "a -> b.md", From: "c d.md"}}},
	} {
		if got := ParseStatus(tc.text); !reflect.DeepEqual(got, tc.want) {
			t.Errorf("%s: got %+v, want %+v", tc.name, got, tc.want)
		}
	}
}

func testGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	command := exec.Command("git", append([]string{"-C", dir, "-c", "user.name=Fixture", "-c", "user.email=fixture@example.invalid", "-c", "commit.gpgsign=false"}, args...)...)
	if data, err := command.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v %s", args, err, data)
	}
}

// Approval and completion always rewrite the record they touch, so the only
// place "everything committed lists nothing" can be shown is the check itself.
func TestUncommittedRecordsListsNothingOnACommittedTree(t *testing.T) {
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	t.Setenv("GIT_CONFIG_GLOBAL", filepath.Join(t.TempDir(), "gitconfig"))
	root := t.TempDir()
	for rel, text := range map[string]string{
		"intent/i001-billing.md":         "---\nid: i001\n---\n",
		"plans/p0001-api/plan.md":        "---\nid: p0001\n---\n",
		".context-circuit/ids/maya.yaml": "intents: [i001]\nplans: [p0001]\n",
	} {
		path := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(text), 0644); err != nil {
			t.Fatal(err)
		}
	}
	testGit(t, root, "init", "-b", "main")
	testGit(t, root, "add", ".")
	testGit(t, root, "commit", "-m", "test: seed")
	s := &Store{Root: root}
	files, unchecked := s.UncommittedRecords(context.Background(), "p0001", "i001")
	if unchecked != "" || len(files) != 0 {
		t.Fatalf("a committed tree listed %+v (unchecked %q)", files, unchecked)
	}
	if err := os.WriteFile(filepath.Join(root, "plans", "p0001-api", "plan.md"), []byte("---\nid: p0001\n---\n\nChanged.\n"), 0644); err != nil {
		t.Fatal(err)
	}
	files, _ = s.UncommittedRecords(context.Background(), "p0001")
	if want := []UncommittedFile{{Path: "plans/p0001-api/plan.md", Change: "modified"}}; !reflect.DeepEqual(files, want) {
		t.Fatalf("got %+v, want %+v", files, want)
	}
	// An ID never matches a longer one that starts with it.
	if files, _ = s.UncommittedRecords(context.Background(), "p000"); len(files) != 0 {
		t.Fatalf("a shorter ID matched a longer record: %+v", files)
	}
}

func TestStatusChangeMarks(t *testing.T) {
	for code, want := range map[string]string{
		"??": "new", "A ": "new", "AM": "new",
		" M": "modified", "M ": "modified", "MM": "modified", " T": "modified",
		" D": "deleted", "D ": "deleted", "MD": "deleted",
		"UU": "modified", "AA": "modified", "DD": "modified", "UD": "modified",
	} {
		if got := statusChange(code); got != want {
			t.Errorf("%q: got %s, want %s", code, got, want)
		}
	}
}
