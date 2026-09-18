package cli_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// commitAt dates a commit, which is what lets these tests put code changes on
// either side of the day a note was last confirmed.
func commitAt(t *testing.T, repo, date, name, body, message string) {
	t.Helper()
	write(t, filepath.Join(repo, filepath.FromSlash(name)), body)
	run := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", append([]string{"-C", repo, "-c", "commit.gpgsign=false",
			"-c", "user.name=Fixture", "-c", "user.email=fixture@example.invalid"}, args...)...)
		cmd.Env = append(os.Environ(), "GIT_AUTHOR_DATE="+date+"T12:00:00Z", "GIT_COMMITTER_DATE="+date+"T12:00:00Z")
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v %s", args, err, out)
		}
	}
	run("add", "-A")
	run("commit", "-m", message)
}

// anchoredWorkspace is a workspace with one repository and one note anchored
// into it, catalogued with the date it was last confirmed.
func anchoredWorkspace(t *testing.T, f fixture, reviewed string) string {
	t.Helper()
	repo := filepath.Join(f.home, "api")
	if err := os.Mkdir(repo, 0755); err != nil {
		t.Fatal(err)
	}
	git(t, repo, "init", "-b", "main")
	commitAt(t, repo, "2026-01-05", "internal/billing/dunning/retry.go", "package dunning\n", "feat: seed")
	f.ok("repo", "connect", "--id", "api", "--path", repo, "--base", "main")

	write(t, filepath.Join(f.root, "context", "domains", "invoice-lifecycle.md"),
		"# Invoice lifecycle\n\nWhen an invoice is voided rather than credited.\n\n"+
			"## When is an invoice voided?\n\nBefore the first payment attempt.\n\n"+
			"Owner:\n\n- `api@internal/billing/dunning/` — the retry schedule.\n")
	write(t, filepath.Join(f.root, "context", "INDEX.md"),
		"# Shared project knowledge\n\n## Domains\n\n*What the system does in each area.*\n\n"+
			"- [Invoice lifecycle](domains/invoice-lifecycle.md) {api} — when an invoice is voided · invoices, dunning · reviewed "+reviewed+"\n")
	return repo
}

func reviewFindings(t *testing.T, root string) []finding {
	t.Helper()
	var reviews []finding
	for _, item := range checkFindings(t, root) {
		if strings.Contains(item.Issue, "was last confirmed") || strings.Contains(item.Issue, "anchors to no code") {
			reviews = append(reviews, item)
		}
	}
	return reviews
}

// The point of reading the date against the code rather than the calendar: a
// note whose anchors nobody has touched is not stale, however old its date is.
func TestAnUntouchedAnchorIsNotStaleHoweverOldTheDate(t *testing.T) {
	f := setup(t)
	anchoredWorkspace(t, f, "2026-06-01")
	if reviews := reviewFindings(t, f.root); len(reviews) != 0 {
		t.Fatalf("a note whose code has not moved must stay silent: %v", reviews)
	}
}

func TestCodeMovingUnderAnAnchorIsReportedWithItsEvidence(t *testing.T) {
	f := setup(t)
	repo := anchoredWorkspace(t, f, "2026-06-01")
	commitAt(t, repo, "2026-08-02", "internal/billing/dunning/retry.go", "package dunning\n\n// wind-down\n", "fix: shorten the wind-down")
	commitAt(t, repo, "2026-08-03", "internal/billing/dunning/retry.go", "package dunning\n\n// wind-down\n// again\n", "fix: again")

	reviews := reviewFindings(t, f.root)
	if len(reviews) != 1 {
		t.Fatalf("expected one review finding: %v", reviews)
	}
	for _, want := range []string{"domains/invoice-lifecycle.md", "last confirmed 2026-06-01", "api@internal/billing/dunning/", "2 commits since"} {
		if !strings.Contains(reviews[0].Issue, want) {
			t.Errorf("finding does not carry %q: %s", want, reviews[0].Issue)
		}
	}
	if !strings.HasPrefix(reviews[0].Resolve, "needs a person:") {
		t.Errorf("re-confirming a note is a judgment, not a command: %s", reviews[0].Resolve)
	}

	// A change somewhere else in the repository is not this note's business.
	commitAt(t, repo, "2026-08-04", "internal/shipping/label.go", "package shipping\n", "feat: unrelated")
	if reviews := reviewFindings(t, f.root); !strings.Contains(reviews[0].Issue, "2 commits since") {
		t.Errorf("only commits under the note's own anchors count: %s", reviews[0].Issue)
	}
}

// A note confirmed on the day the code changed was confirmed against that
// change, so the day itself does not count against it.
func TestTheDayANoteWasConfirmedDoesNotCountAgainstIt(t *testing.T) {
	f := setup(t)
	repo := anchoredWorkspace(t, f, "2026-06-01")
	commitAt(t, repo, "2026-06-01", "internal/billing/dunning/retry.go", "package dunning\n\n// same day\n", "fix: same day")
	if reviews := reviewFindings(t, f.root); len(reviews) != 0 {
		t.Fatalf("a same-day commit must not report the note it was confirmed against: %v", reviews)
	}
}

// A note anchoring to no code has no evidence to read, so calendar age is all
// there is — and it is reported only where a workspace asks for it.
func TestAnUnanchoredNoteIsSilentUntilAWorkspaceAsks(t *testing.T) {
	f := setup(t)
	write(t, filepath.Join(f.root, "context", "references", "billing-contract.md"),
		"# Billing contract\n\nWhat the provider promises.\n\n## What fields does it send?\n\nAn amount and a currency.\n")
	write(t, filepath.Join(f.root, "context", "INDEX.md"),
		"# Shared project knowledge\n\n## References\n\n*Facts looked up rather than read.*\n\n"+
			"- [Billing contract](references/billing-contract.md) — what the provider promises · billing, contract · reviewed 2020-01-01\n")

	if reviews := reviewFindings(t, f.root); len(reviews) != 0 {
		t.Fatalf("age alone reports nothing until a workspace sets a threshold: %v", reviews)
	}

	shared := read(t, filepath.Join(f.root, "workspace.yaml"))
	write(t, filepath.Join(f.root, "workspace.yaml"), shared+"knowledge_review_days: 180\n")

	reviews := reviewFindings(t, f.root)
	if len(reviews) != 1 {
		t.Fatalf("expected the threshold to report it: %v", reviews)
	}
	for _, want := range []string{"references/billing-contract.md", "anchors to no code", "last confirmed 2020-01-01"} {
		if !strings.Contains(reviews[0].Issue, want) {
			t.Errorf("finding does not carry %q: %s", want, reviews[0].Issue)
		}
	}
}

// A repository this machine has not obtained already has a finding of its own.
// Reporting every note anchored into it as well would bury that one.
func TestAnUnobtainedRepositoryIsNotReportedTwice(t *testing.T) {
	f := setup(t)
	anchoredWorkspace(t, f, "2026-06-01")
	local := read(t, filepath.Join(f.root, "repositories.local.yaml"))
	write(t, filepath.Join(f.root, "repositories.local.yaml"), strings.SplitN(local, "bindings:", 2)[0]+"bindings: {}\n")

	if reviews := reviewFindings(t, f.root); len(reviews) != 0 {
		t.Fatalf("expected no review findings for an unbound repository: %v", reviews)
	}
	var bound bool
	for _, item := range checkFindings(t, f.root) {
		if strings.Contains(item.Issue, "api") && strings.Contains(item.Resolve, "repo clone") {
			bound = true
		}
	}
	if !bound {
		t.Error("expected the binding itself to be reported instead")
	}
}
