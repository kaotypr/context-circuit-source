package cli_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

type gateResult struct {
	CommitRecommended *struct {
		Files []struct {
			Path   string `json:"path"`
			Change string `json:"change"`
		} `json:"files"`
		Note string `json:"note"`
	} `json:"commit_recommended"`
	CommitUnchecked string `json:"commit_unchecked"`
}

// gate runs approve or complete and returns what it recommends committing as
// path -> change.
func (f fixture) gate(args ...string) (map[string]string, gateResult) {
	f.t.Helper()
	var result gateResult
	if err := json.Unmarshal([]byte(f.ok(args...)), &result); err != nil {
		f.t.Fatal(err)
	}
	files := map[string]string{}
	if result.CommitRecommended != nil {
		for _, file := range result.CommitRecommended.Files {
			files[file.Path] = file.Change
		}
		if !strings.Contains(result.CommitRecommended.Note, "authorization") {
			f.t.Fatalf("the recommendation does not say committing needs authorization: %q", result.CommitRecommended.Note)
		}
	}
	return files, result
}

func commitRecords(t *testing.T, dir string) {
	t.Helper()
	git(t, dir, "add", "-A")
	git(t, dir, "commit", "-q", "-m", "test: commit records")
}

func expectFiles(t *testing.T, got, want map[string]string) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("commit_recommended files: got %v, want %v", got, want)
	}
}

func TestApprovalRecommendsCommittingItsRecords(t *testing.T) {
	f := setup(t)
	git(t, f.root, "init", "-q", "-b", "main")
	commitRecords(t, f.root)

	// An intent never committed, and an unrelated one left uncommitted beside it.
	i := f.intent("billing")
	other := f.intent("other")
	files, _ := f.gate("record", "approve", "--id", i.ID, "--text", "Go ahead.")
	expectFiles(t, files, map[string]string{i.Path: "new", ".context-circuit/ids/maya.yaml": "modified"})
	if _, listed := files[other.Path]; listed {
		t.Fatalf("an unrelated intent was listed: %v", files)
	}

	// A committed intent that approval rewrites.
	commitRecords(t, f.root)
	files, _ = f.gate("record", "approve", "--id", other.ID, "--text", "Go ahead.")
	expectFiles(t, files, map[string]string{other.Path: "modified"})

	// A renamed record is listed once under each path, not as a stray old path.
	commitRecords(t, f.root)
	archived := "intent/archive/" + filepath.Base(i.Path)
	if err := os.MkdirAll(filepath.Join(f.root, "intent", "archive"), 0755); err != nil {
		t.Fatal(err)
	}
	git(t, f.root, "mv", filepath.FromSlash(i.Path), filepath.FromSlash(archived))
	files, _ = f.gate("record", "approve", "--id", i.ID, "--text", "Still go ahead.")
	expectFiles(t, files, map[string]string{archived: "new", i.Path: "deleted"})
}

func TestCompletionRecommendsCommittingThePlanAndItsIntent(t *testing.T) {
	f := setup(t)
	f.repository("api")
	git(t, f.root, "init", "-q", "-b", "main")
	i := f.intent("billing")
	p := f.plan(i.ID, "api", "api")
	standalone := record(t, f.ok("record", "create", "--kind", "plan", "--slug", "alone", "--title", "Alone", "--repo", "api"))
	commitRecords(t, f.root)

	files, _ := f.gate("record", "complete", "--id", p.ID, "--text", "Delivered.")
	expectFiles(t, files, map[string]string{p.Path: "modified"})

	commitRecords(t, f.root)
	write(t, filepath.Join(f.root, filepath.FromSlash(i.Path)), read(t, filepath.Join(f.root, filepath.FromSlash(i.Path)))+"\nEdited.\n")
	files, _ = f.gate("record", "complete", "--id", p.ID, "--text", "Delivered again.")
	expectFiles(t, files, map[string]string{p.Path: "modified", i.Path: "modified"})

	commitRecords(t, f.root)
	files, _ = f.gate("record", "complete", "--id", standalone.ID, "--text", "Done.")
	expectFiles(t, files, map[string]string{standalone.Path: "modified"})
}

// Before the per-member split the one shared file is an ID file too, and the
// split approval itself triggers removes it.
func TestApprovalNamesBothIDLayouts(t *testing.T) {
	f := setup(t)
	i := f.intent("billing")
	schema2(t, f, []string{i.ID}, nil)
	git(t, f.root, "init", "-q", "-b", "main")
	commitRecords(t, f.root)
	files, _ := f.gate("record", "approve", "--id", i.ID, "--text", "Go ahead.")
	expectFiles(t, files, map[string]string{
		i.Path:                           "modified",
		".context-circuit/ids.yaml":      "deleted",
		".context-circuit/ids/maya.yaml": "new",
	})
}

// A workspace kept inside another repository is reported relative to its own
// root, and what lies outside it is not its business.
func TestANestedWorkspaceReportsItsOwnPaths(t *testing.T) {
	f := setup(t)
	git(t, f.home, "init", "-q", "-b", "main")
	i := f.intent("billing")
	write(t, filepath.Join(f.home, filepath.FromSlash(i.Path)), "---\nid: "+i.ID+"\n---\n")
	files, _ := f.gate("record", "approve", "--id", i.ID, "--text", "Go ahead.")
	expectFiles(t, files, map[string]string{i.Path: "new", ".context-circuit/ids/maya.yaml": "new"})
}

func TestTheGatesSucceedWhenGitCannotSay(t *testing.T) {
	f := setup(t)
	f.repository("api")
	i := f.intent("billing")
	p := f.plan(i.ID, "api", "api")

	// The default fixture is not a Git repository.
	_, result := f.gate("record", "approve", "--id", i.ID, "--text", "Go ahead.")
	if result.CommitRecommended != nil || !strings.Contains(result.CommitUnchecked, "not in a Git repository") {
		t.Fatalf("outside a repository: %+v", result)
	}

	// Git missing altogether. Each call looks git up on PATH, so an empty
	// directory removes it; on Windows this replaces Path as well.
	t.Setenv("PATH", t.TempDir())
	for _, args := range [][]string{
		{"record", "approve", "--id", i.ID, "--text", "Again."},
		{"record", "complete", "--id", p.ID, "--text", "Delivered."},
	} {
		code, out, errOut := call(f.root, args...)
		if code != 0 {
			t.Fatalf("%v failed without Git: %s %s", args, out, errOut)
		}
		var missing gateResult
		if err := json.Unmarshal([]byte(out), &missing); err != nil {
			t.Fatal(err)
		}
		if missing.CommitRecommended != nil || !strings.Contains(missing.CommitUnchecked, "Git is not available") {
			t.Fatalf("%v without Git: %s", args, out)
		}
	}
	if _, err := os.Stat(filepath.Join(f.root, ".git")); !os.IsNotExist(err) {
		t.Fatal("the check created a repository")
	}
}

// The rename fix is in the one reader `status` and `repo inspect` share.
func TestInspectListsARenameOnce(t *testing.T) {
	f := setup(t)
	api := f.repository("api")
	git(t, api, "mv", "README.md", "docs.md")
	var snapshot struct {
		Changes []string `json:"changes"`
	}
	if err := json.Unmarshal([]byte(f.ok("repo", "inspect", "--id", "api")), &snapshot); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(snapshot.Changes, []string{"R  docs.md"}) {
		t.Fatalf("changes: %q", snapshot.Changes)
	}
}
