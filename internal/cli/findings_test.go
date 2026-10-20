package cli_test

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
)

type finding struct {
	Issue   string `json:"issue"`
	Resolve string `json:"resolve"`
}

func checkFindings(t *testing.T, root string) []finding {
	t.Helper()
	_, out, _ := call(root, "check")
	var result struct {
		Issues []finding `json:"issues"`
	}
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("check output is not readable: %v\n%s", err, out)
	}
	return result.Issues
}

// breakThings puts a workspace into several unrelated faulty states at once, so
// the invariant below is asserted across findings from every part of the
// diagnostic rather than one convenient corner of it.
func breakThings(t *testing.T, f fixture) {
	t.Helper()
	git(t, f.root, "init", "-b", "main")

	// A member an earlier version left without a band.
	write(t, filepath.Join(f.root, "members.yaml"), "members:\n  maya:\n    name: Maya\n    band: 0\n  alex:\n    name: Alex\n")

	// A knowledge repository nobody here has obtained.
	shared := read(t, filepath.Join(f.root, "workspace.yaml"))
	write(t, filepath.Join(f.root, "workspace.yaml"), shared+
		"knowledge_repositories:\n  core-service-knowledge:\n    url: https://example.invalid/k.git\n    default_branch: main\n    index: index.md\n")

	// A plan nobody reserved, linked to nothing, naming no repository.
	write(t, filepath.Join(f.root, "plans", "p9999-orphan.md"),
		"---\nid: p9999\ncreated_by: maya\ncreated_at: 2026-09-15T10:53:00Z\n---\n\n# Orphan\n")

	// A note the catalog never lists, beside a catalog entry pointing nowhere.
	write(t, filepath.Join(f.root, "context", "domains", "billing.md"), "# Billing\n\nWhat this settles.\n")
	write(t, filepath.Join(f.root, "context", "INDEX.md"),
		"# Shared project knowledge\n\n## Domains\n\n*What the system does in each area.*\n\n- [Missing](domains/gone.md) {api} — nothing · missing · reviewed 2026-02-04\n")
}

// Every finding carries what discharges it. The issue alone was what this
// reported for its whole life, and it left an agent to derive the remedy in the
// one situation where the workspace is already inconsistent.
func TestEveryFindingNamesWhatDischargesIt(t *testing.T) {
	f := setup(t)
	breakThings(t, f)

	findings := checkFindings(t, f.root)
	if len(findings) < 5 {
		t.Fatalf("expected the broken workspace to produce findings across the diagnostic, got %d: %v", len(findings), findings)
	}
	for _, item := range findings {
		if strings.TrimSpace(item.Issue) == "" {
			t.Errorf("finding with no issue: %+v", item)
		}
		if strings.TrimSpace(item.Resolve) == "" {
			t.Errorf("finding names no resolution: %q", item.Issue)
		}
	}
}

// A resolution is an instruction, and the three kinds are distinguishable: a
// command to run, an edit to make, or a decision that is somebody's to take.
func TestAResolutionIsACommandAnEditOrAPerson(t *testing.T) {
	f := setup(t)
	breakThings(t, f)
	findings := checkFindings(t, f.root)

	want := map[string]string{
		"is not obtained on this machine":       "`knowledge clone --id core-service-knowledge`",
		"member alex holds no band":             "`member band --id alex --band 1`",
		"describes no workspace repository":     "`workspace connect --base BRANCH`",
		"record ID reserved in no ID file":      "needs a person:",
		"note is not listed in the catalog":     "add its entry to context/INDEX.md",
		"catalog entry links to a missing note": "write the note, or correct the link",
	}
	for issue, resolve := range want {
		var seen bool
		for _, item := range findings {
			if strings.Contains(item.Issue, issue) {
				seen = true
				if !strings.Contains(item.Resolve, resolve) {
					t.Errorf("%q resolves to %q, expected it to name %q", issue, item.Resolve, resolve)
				}
			}
		}
		if !seen {
			t.Errorf("the broken workspace produced no finding matching %q; findings were %v", issue, findings)
		}
	}
}

// Two clones assigning bands at once can leave a roster nobody may silently
// repair: the CLI does not pick a winner, so each case is a person's to fix,
// and status still runs so they can see it.
func TestABrokenRosterNeedsAPerson(t *testing.T) {
	for name, roster := range map[string]string{
		"conflict markers": "members:\n  maya:\n    name: Maya\n<<<<<<< HEAD\n    band: 1\n=======\n    band: 2\n>>>>>>> theirs\n",
		"duplicate band":   "members:\n  maya:\n    name: Maya\n    band: 1\n  alex:\n    name: Alex\n    band: 1\n",
		"two band keys":    "members:\n  maya:\n    name: Maya\n    band: 1\n    band: 2\n",
		"listed twice":     "members:\n  maya:\n    name: Maya\n    band: 1\n  maya:\n    name: Maya\n    band: 2\n",
	} {
		f := setup(t)
		write(t, filepath.Join(f.root, "members.yaml"), roster)
		status := orientation(t, f.ok("status"))
		if status.IdentityRequired == "" {
			t.Errorf("%s: status did not ask for setup", name)
		}
		var seen bool
		for _, item := range checkFindings(t, f.root) {
			if strings.HasPrefix(item.Issue, "members.yaml: ") {
				seen = true
				if !strings.HasPrefix(item.Resolve, "needs a person:") {
					t.Errorf("%s: resolves to %q", name, item.Resolve)
				}
				if name == "conflict markers" && !strings.Contains(item.Issue, "conflict markers") {
					t.Errorf("conflict markers were not named: %q", item.Issue)
				}
			}
		}
		if !seen {
			t.Errorf("%s: check reported nothing about the roster", name)
		}
		if out := f.fail("record", "create", "--kind", "intent", "--slug", "x", "--title", "X"); !strings.Contains(out, "setup is incomplete") {
			t.Errorf("%s: a change was not refused: %s", name, out)
		}
	}
}

// A clean workspace reports nothing and exits zero, which is what makes the
// findings above worth reading when they appear.
func TestACleanWorkspaceStillReportsNothing(t *testing.T) {
	f := setup(t)
	f.repository("api")
	code, out, _ := call(f.root, "check")
	if code != 0 {
		t.Fatalf("expected a clean workspace to pass: %s", out)
	}
	if findings := checkFindings(t, f.root); len(findings) != 0 {
		t.Fatalf("expected no findings: %v", findings)
	}
}
