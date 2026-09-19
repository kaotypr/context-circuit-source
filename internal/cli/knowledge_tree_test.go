package cli_test

import (
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// `context/` is this project's own product knowledge, and the only thing that
// validates it is the product's own diagnostic. Running that used to require a
// maintainer to have registered this checkout as a workspace by hand, which put
// the validation behind two machine-local files: it ran on one laptop, never in
// CI, and a fresh clone got silence rather than a warning. The workspace is
// built here instead, so every clone and every pull request gets the same
// checks without the source repository having to pretend to be a workspace.
func TestTheProjectsOwnKnowledgeIsValid(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, ".git")); err != nil {
		t.Skip("not a Git checkout; the reviewed-date pass would have no history to read")
	}

	f := setup(t)
	if err := os.RemoveAll(filepath.Join(f.root, "context")); err != nil {
		t.Fatal(err)
	}
	copyTree(t, filepath.Join(root, "context"), filepath.Join(f.root, "context"))

	// Every note anchors to this repository by ID, so the reviewed-date pass
	// needs a binding to this checkout to resolve them. A branch name that does
	// not exist locally is fine — the pass falls back to HEAD, which is what a
	// detached CI checkout has — but it must be well formed.
	f.ok("repo", "connect", "--id", "context-circuit-source", "--path", root, "--base", currentBranch(root))

	var mechanical, judgment []finding
	for _, issue := range checkFindings(t, f.root) {
		// A finding naming an unobtained checkout means the anchors resolved to
		// nothing, which makes the reviewed-date pass silent rather than clean.
		// Reporting that as success is the one failure this test cannot afford.
		if strings.Contains(issue.Issue, "connect this machine's checkout") {
			t.Fatalf("the binding did not resolve, so the knowledge checks ran against nothing:\n  %s", issue.Issue)
		}
		if strings.HasPrefix(issue.Resolve, "needs a person: ") {
			judgment = append(judgment, issue)
			continue
		}
		mechanical = append(mechanical, issue)
	}

	for _, issue := range mechanical {
		t.Errorf("context/ is not clean:\n  %s\n  %s", issue.Issue, issue.Resolve)
	}

	// A finding the product marks `needs a person:` is one no command and no
	// edit discharges — a note to re-read against code that moved, most often.
	// Blocking a pull request on somebody else's judgment would make this red
	// for every contributor who touched anchored code, so it is reported and
	// left to a maintainer.
	for _, issue := range judgment {
		t.Logf("needs a maintainer's judgment: %s", issue.Issue)
	}
	if depth := historyDepth(root); depth < 2 {
		t.Logf("this checkout has %d commit(s); a shallow clone makes the reviewed-date pass vacuous", depth)
	}
}

// currentBranch reports the checked-out branch, or a well-formed placeholder
// when HEAD is detached, which is what CI produces for a pull request.
func currentBranch(root string) string {
	out, err := exec.Command("git", "-C", root, "symbolic-ref", "--quiet", "--short", "HEAD").Output()
	if branch := strings.TrimSpace(string(out)); err == nil && branch != "" {
		return branch
	}
	return "main"
}

func historyDepth(root string) int {
	out, err := exec.Command("git", "-C", root, "rev-list", "--count", "HEAD").Output()
	if err != nil {
		return 0
	}
	depth := 0
	for _, r := range strings.TrimSpace(string(out)) {
		if r < '0' || r > '9' {
			return 0
		}
		depth = depth*10 + int(r-'0')
	}
	return depth
}

func copyTree(t *testing.T, from, to string) {
	t.Helper()
	if err := filepath.WalkDir(from, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(from, path)
		if err != nil {
			return err
		}
		target := filepath.Join(to, rel)
		if entry.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, 0644)
	}); err != nil {
		t.Fatal(err)
	}
}
