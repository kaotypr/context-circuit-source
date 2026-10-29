package scenario

import (
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/kaotypr/context-circuit-source/agent-tests/internal/record"
)

func TestCatalogue(t *testing.T) {
	dir := filepath.Join("..", "..", "scenarios")
	loaded, err := LoadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(loaded) != 12 {
		t.Fatalf("scenarios %d", len(loaded))
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	byID := map[string]Scenario{}
	for _, sc := range loaded {
		byID[sc.ID] = sc
	}
	seenName := map[string]bool{}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yaml") {
			continue
		}
		if seenName[entry.Name()] {
			t.Fatalf("duplicate file %s", entry.Name())
		}
		seenName[entry.Name()] = true
		sc, err := Load(filepath.Join(dir, entry.Name()))
		if err != nil {
			t.Fatal(err)
		}
		if entry.Name() != sc.ID+".yaml" {
			t.Fatalf("file %s id %s", entry.Name(), sc.ID)
		}
		// A scenario's version is bumped when its prompt, setup or check
		// semantics change, so earlier runs are not compared with later ones.
		wantVersion := map[string]int{"plan": 2, "intent": 2, "intent-lang": 2, "borrowed": 2, "archive": 2, "complete": 2}
		want := wantVersion[sc.ID]
		if want == 0 {
			want = 1
		}
		if sc.Version != want {
			t.Fatalf("%s version %d, want %d", sc.ID, sc.Version, want)
		}
		if strings.TrimSpace(sc.Title) == "" || strings.TrimSpace(sc.Prompt) == "" {
			t.Fatalf("%s missing title or prompt", sc.ID)
		}
		wantSetup := map[string]bool{
			"join": true, "intent-lang": true, "plan": true,
			"execute": true, "complete": true, "borrowed": true, "archive": true,
			"push-records": true, "push-records-conflict": true,
		}
		if wantSetup[sc.ID] {
			if len(sc.Setup) == 0 || strings.TrimSpace(strings.Join(sc.Setup, "")) == "" {
				t.Fatalf("%s setup is empty", sc.ID)
			}
		} else if len(sc.Setup) != 0 {
			t.Fatalf("%s has setup", sc.ID)
		}
		shared := []string{
			"launcher-only", "nothing-pushed", "no-stray-commits", "no-attribution",
			"check-clean", "sources-passive", "no-outside-writes",
		}
		if len(sc.Checks) < len(shared) {
			t.Fatalf("%s checks %v", sc.ID, sc.Checks)
		}
		for i, name := range shared {
			if sc.Checks[i] != name {
				t.Fatalf("%s check %d is %s", sc.ID, i, sc.Checks[i])
			}
		}
		names := map[string]bool{}
		for _, name := range sc.Checks {
			if names[name] {
				t.Fatalf("%s duplicate check %s", sc.ID, name)
			}
			names[name] = true
			if Registry[name] == nil {
				t.Fatalf("%s check %s is not registered", sc.ID, name)
			}
		}
		for _, art := range sc.DesignArtifacts {
			if art == "" || filepath.IsAbs(art) || strings.Contains(art, "..") {
				t.Fatalf("%s artifact %q", sc.ID, art)
			}
			clean := path.Clean(art)
			if clean == ".." || strings.HasPrefix(clean, "../") {
				t.Fatalf("%s artifact %q escapes", sc.ID, art)
			}
		}
	}
	if byID["join"].Workspace != "ws2" {
		t.Fatalf("join workspace %q", byID["join"].Workspace)
	}
	if len(byID["execute"].AllowedBranches) != 1 || byID["execute"].AllowedBranches[0] != "cc/p1000/api" {
		t.Fatalf("execute branches %v", byID["execute"].AllowedBranches)
	}
	for id, sc := range byID {
		push := id == "push-records" || id == "push-records-conflict"
		if sc.PushWorkspaceBase != push {
			t.Fatalf("%s push_workspace_base %v", id, sc.PushWorkspaceBase)
		}
	}
	if Registry["not-a-check"] != nil {
		t.Fatal("an unknown check name is registered")
	}
	digest, err := record.Digest(filepath.Join("..", "..", "fixtures"))
	if err != nil {
		t.Fatal(err)
	}
	hashes := map[string]string{}
	for _, sc := range loaded {
		sum := sc.Hash(digest)
		if prev, ok := hashes[sum]; ok {
			t.Fatalf("hash collision %s and %s", prev, sc.ID)
		}
		hashes[sum] = sc.ID
	}
	banned := regexp.MustCompile(`(?i)\b(pilot|maintainer)\b|/Users`)
	roots := []string{dir, filepath.Join("..", "..", "fixtures")}
	for _, root := range roots {
		err := filepath.WalkDir(root, func(p string, d os.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return err
			}
			body, err := os.ReadFile(p)
			if err != nil {
				return err
			}
			if banned.Match(body) {
				rel, _ := filepath.Rel(root, p)
				t.Errorf("%s contains a banned word or a local path", rel)
			}
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestSetbodyKeepsFrontmatter(t *testing.T) {
	dir := t.TempDir()
	recordPath := filepath.Join(dir, "record.md")
	bodyPath := filepath.Join(dir, "body.md")
	record := "---\nid: i100\ncreated_by: maya\n---\n\n# Old\n\n## Goal\n\nDescribe the intended outcome.\n"
	body := "\n# Invoice CSV export\n\n## Goal\n\nCustomers can download invoices.\n"
	if err := os.WriteFile(recordPath, []byte(record), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(bodyPath, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	script := filepath.Join("..", "..", "fixtures", "bin", "setbody.sh")
	cmd := execCommand(t, "sh", script, recordPath, bodyPath)
	if cmd != "" {
		t.Fatal(cmd)
	}
	got, err := os.ReadFile(recordPath)
	if err != nil {
		t.Fatal(err)
	}
	text := string(got)
	if !strings.Contains(text, "id: i100") || !strings.Contains(text, "created_by: maya") {
		t.Fatalf("frontmatter lost:\n%s", text)
	}
	if strings.Contains(text, "Describe the intended outcome.") || strings.Contains(text, "# Old") {
		t.Fatalf("old body kept:\n%s", text)
	}
	if !strings.Contains(text, "Customers can download invoices.") {
		t.Fatalf("new body missing:\n%s", text)
	}
}

func TestPlanArtifactGlob(t *testing.T) {
	ok, err := path.Match("plans/p1000-*/plan.md", "plans/p1000-invoice-csv-export/plan.md")
	if err != nil || !ok {
		t.Fatalf("match %v %v", ok, err)
	}
	ok, err = path.Match("cc/p1000/*", "cc/p1000/api")
	if err != nil || !ok {
		t.Fatalf("branch %v %v", ok, err)
	}
}
