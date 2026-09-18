package workspace

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestYAMLContainerEdits(t *testing.T) {
	for _, tc := range []struct {
		name, input string
		keys        []string
		value       any
		preserve    []string
	}{
		{"empty flow", "# shared\nrepositories: {}\n", []string{"repositories", "api"}, Repository{DefaultBranch: "main"}, []string{"# shared", "api", "main"}},
		{"populated flow", "repositories: {api: {default_branch: main}} # shared\n", []string{"repositories", "web"}, Repository{DefaultBranch: "develop"}, []string{"# shared", "api", "main", "web", "develop"}},
		{"replace flow mapping", "bindings: {api: {path: old}} # local\n", []string{"bindings", "api"}, Binding{Path: "new path"}, []string{"# local", "new path"}},
		{"sequence", "plans: [] # links\n", []string{"plans"}, []string{"p0001", "p0002"}, []string{"# links", "p0001", "p0002"}},
		{"scalar", "base: main # target\nother: 'quoted'\n", []string{"base"}, "next", []string{"# target", "'quoted'", "next"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Edit([]byte(tc.input), tc.keys, tc.value)
			if err != nil {
				t.Fatal(err)
			}
			var decoded map[string]any
			if err := Decode(got, &decoded); err != nil {
				t.Fatalf("%v: %s", err, got)
			}
			for _, fragment := range tc.preserve {
				if !strings.Contains(string(got), fragment) {
					t.Fatalf("lost %q: %s", fragment, got)
				}
			}
		})
	}
}

// A seed's empty `{}` container must grow into block style rather than extend a
// single flow line as entries accumulate.
func TestEmptyContainerGrowsInBlockStyle(t *testing.T) {
	got, err := Edit([]byte("bindings: {}\n"), []string{"bindings", "tembiter"}, Binding{"/ws/tembiter", "main"})
	if err != nil {
		t.Fatal(err)
	}
	want := "bindings:\n  tembiter:\n    path: /ws/tembiter\n    base_branch: main\n"
	if string(got) != want {
		t.Fatalf("first entry\nwant %q\ngot  %q", want, got)
	}
	got, err = Edit(got, []string{"bindings", "other"}, Binding{"/ws/other", "main"})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(got), "{") {
		t.Fatalf("second entry fell back to flow style: %s", got)
	}
	var decoded Bindings
	if err := Decode(got, &decoded); err != nil {
		t.Fatalf("%v: %s", err, got)
	}
	if decoded.Bindings["tembiter"].Path != "/ws/tembiter" || decoded.Bindings["other"].Path != "/ws/other" {
		t.Fatalf("lost a binding: %#v", decoded.Bindings)
	}

	// The workspace seed keeps its sibling keys and header comment.
	seed := "# shared\nversion: 2\nrepositories: {}\nrelationships: []\n"
	ws, err := Edit([]byte(seed), []string{"repositories", "tembiter"}, Repository{DefaultBranch: "main"})
	if err != nil {
		t.Fatal(err)
	}
	for _, fragment := range []string{"# shared", "version: 2", "relationships: []", "\n  tembiter:\n    default_branch: main\n"} {
		if !strings.Contains(string(ws), fragment) {
			t.Fatalf("lost %q: %s", fragment, ws)
		}
	}

	// A populated flow container still receives a flow value, because mixing
	// styles through MergeFromReader does not serialize.
	mixed, err := Edit([]byte("bindings: {api: {path: /ws/api}} # local\n"), []string{"bindings", "web"}, Binding{"/ws/web", "main"})
	if err != nil {
		t.Fatal(err)
	}
	var reparsed Bindings
	if err := Decode(mixed, &reparsed); err != nil {
		t.Fatalf("%v: %s", err, mixed)
	}
	if len(reparsed.Bindings) != 2 {
		t.Fatalf("lost a binding: %s", mixed)
	}
}

func TestAmbiguousYAML(t *testing.T) {
	for _, input := range []string{"a: 1\na: 2\n", "a: 1\n---\na: 2\n"} {
		var decoded map[string]any
		if err := Decode([]byte(input), &decoded); err == nil {
			t.Fatalf("accepted ambiguous YAML: %s", input)
		}
	}
}

// A subprocess holds the real OS lock until killed, exercising crash recovery
// across processes rather than only goroutines in one executable.
func TestLockProcessHelper(t *testing.T) {
	root := os.Getenv("CC_LOCK_TEST_ROOT")
	if root == "" {
		return
	}
	s, err := Open(root)
	if err != nil {
		t.Fatal(err)
	}
	_, err = s.WithLock(context.Background(), func() (any, error) {
		if err := os.WriteFile(filepath.Join(root, "ready"), []byte("ready"), 0600); err != nil {
			return nil, err
		}
		time.Sleep(30 * time.Second)
		return nil, nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestLockSurvivesProcessInterruption(t *testing.T) {
	root := t.TempDir()
	child := exec.Command(os.Args[0], "-test.run=^TestLockProcessHelper$")
	child.Env = append(os.Environ(), "CC_LOCK_TEST_ROOT="+root)
	if err := child.Start(); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = child.Process.Kill(); _ = child.Wait() }()
	deadline := time.Now().Add(5 * time.Second)
	for {
		if _, err := os.Stat(filepath.Join(root, "ready")); err == nil {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("child never acquired lock")
		}
		time.Sleep(10 * time.Millisecond)
	}
	s, err := Open(root)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	if _, err := s.WithLock(ctx, func() (any, error) { t.Fatal("entered while another process held lock"); return nil, nil }); err == nil {
		t.Fatal("concurrent process lock accepted")
	}
	if err := child.Process.Kill(); err != nil {
		t.Fatal(err)
	}
	_ = child.Wait()
	if _, err := s.WithLock(context.Background(), func() (any, error) { return nil, nil }); err != nil {
		t.Fatalf("lock was not released after process termination: %v", err)
	}
}

func TestKnowledgeIssues(t *testing.T) {
	for _, tc := range []struct {
		name, line string
		flagged    bool
	}{
		{"concept word", "Task progress and temporary results belong in plans.", false},
		{"repository path", "Code: `api@internal/billing/dunning/`.", false},
		{"patterned repository path", "Every feature is `web@src/features/<feature>/` with a `routes.ts`.", false},
		{"repository directory that shares a reserved name", "Loaders live in `api@src/sources/loader.go`.", false},
		{"external URL", "Pricing is published at https://example.com/plans/pricing.", false},
		{"plural business noun", "A customer may hold several subscription plans.", false},
		{"record file", "See plans/p0007-billing.md for the rollout.", true},
		{"bare plan ID", "Reworded during p0007 to match the new flow.", true},
		{"bare intent ID", "Approved as i001.", true},
		{"intent record path", "Superseded by intent/i002-billing.md.", true},
		{"raw evidence path", "Derived from sources/vendor-spec.pdf.", true},
		{"workspace machinery path", "Described in .context-circuit/docs/working.md.", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			if err := os.MkdirAll(filepath.Join(root, "context", "domains"), 0o755); err != nil {
				t.Fatal(err)
			}
			note := filepath.Join(root, "context", "domains", "billing.md")
			if err := os.WriteFile(note, []byte("# Billing\n\n"+tc.line+"\n"), 0o644); err != nil {
				t.Fatal(err)
			}
			store := &Store{Root: root}
			found, err := store.knowledgeIssues()
			if err != nil {
				t.Fatal(err)
			}
			// This boundary is about what a note may name. Where the same line
			// also reads badly, the readability pass reports that separately.
			issues := []Finding{}
			for _, issue := range found {
				if !strings.Contains(issue.Issue, "readability:") {
					issues = append(issues, issue)
				}
			}
			if tc.flagged != (len(issues) > 0) {
				t.Fatalf("flagged=%v for %q: %v", tc.flagged, tc.line, issues)
			}
			if tc.flagged && !strings.HasPrefix(issues[0].Issue, "context/domains/billing.md:3:") {
				t.Fatalf("expected a located finding: %v", issues)
			}
		})
	}
}

func TestKnowledgeIssuesWithoutContextDirectory(t *testing.T) {
	store := &Store{Root: t.TempDir()}
	issues, err := store.knowledgeIssues()
	if err != nil || len(issues) != 0 {
		t.Fatalf("%v %v", err, issues)
	}
}

func TestCatalogConsistency(t *testing.T) {
	for _, tc := range []struct {
		name, index string
		notes       []string
		want        string
	}{
		{"entry and note agree", "- [Billing](domains/billing.md) — voiding rules", []string{"domains/billing.md"}, ""},
		{"entry without a note", "- [Billing](domains/billing.md) — voiding rules", nil, "links to a missing note"},
		{"note without an entry", "Nothing catalogued yet.", []string{"domains/billing.md"}, "not listed in the catalog"},
		{"fenced example catalogs nothing", "```\n- [Example](domains/example.md) — shape only\n```", nil, ""},
		{"prose link is not an entry", "See [Billing](domains/billing.md) for the shape.", nil, ""},
		{"README is navigation", "Nothing catalogued yet.", []string{"domains/README.md"}, ""},
		{"fragment resolves to the note", "- [Billing](domains/billing.md#voiding) — voiding rules", []string{"domains/billing.md"}, ""},
		{"external link is skipped", "- [Vendor](https://example.invalid/spec) — vendor behavior", nil, ""},
		{"entry leaving the catalog", "- [Elsewhere](../README.md) — outside", nil, "links outside the catalog"},
		// An entry naming repositories claims to describe their code, so it
		// says when that claim was last confirmed; one naming none does not.
		{"repository entry without a date", "- [Billing](domains/billing.md) {api} — voiding rules · billing", []string{"domains/billing.md"}, "without a `reviewed"},
		{"repository entry with a date", "- [Billing](domains/billing.md) {api} — voiding rules · billing · reviewed 2026-02-04", []string{"domains/billing.md"}, ""},
		{"several repositories still need one", "- [Billing](domains/billing.md) {api, web} — voiding rules · billing", []string{"domains/billing.md"}, "without a `reviewed"},
		{"an entry naming none needs no date", "- [Glossary](glossary.md) — project vocabulary · glossary", []string{"glossary.md"}, ""},
		{"a fenced example is still exempt", "```\n- [Example](domains/example.md) {api} — shape only\n```", nil, ""},
		// A group says what belongs under it, so the next note has somewhere
		// obvious to go, and one directory answers to one group.
		{"heading without a scope line", "## Domains\n\n- [Billing](domains/billing.md) — voiding rules", []string{"domains/billing.md"}, "no line saying what belongs under it"},
		{"heading with a scope line", "## Domains\n\n*What the system does in each area.*\n\n- [Billing](domains/billing.md) — voiding rules", []string{"domains/billing.md"}, ""},
		{"one concern under two headings", "## Domains\n\n*Scope.*\n\n- [Billing](domains/billing.md) — voiding rules\n\n## Elsewhere\n\n*Scope.*\n\n- [Dunning](domains/dunning.md) — retry schedule", []string{"domains/billing.md", "domains/dunning.md"}, "catalogued under two headings"},
		{"uncatalogued top-level note needs no group", "- [Glossary](glossary.md) — project vocabulary", []string{"glossary.md"}, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			for _, note := range tc.notes {
				full := filepath.Join(root, "context", filepath.FromSlash(note))
				if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(full, []byte("# Note\n"), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			index := filepath.Join(root, "context", "INDEX.md")
			if err := os.MkdirAll(filepath.Dir(index), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(index, []byte("# Catalog\n\n"+tc.index+"\n"), 0o644); err != nil {
				t.Fatal(err)
			}
			issues, err := (&Store{Root: root}).knowledgeIssues()
			if err != nil {
				t.Fatal(err)
			}
			if tc.want == "" {
				if len(issues) != 0 {
					t.Fatalf("expected a consistent catalog: %v", issues)
				}
				return
			}
			if len(issues) != 1 || !strings.Contains(issues[0].Issue, tc.want) {
				t.Fatalf("expected %q: %v", tc.want, issues)
			}
			// Every finding carries what discharges it, or says a person must.
			if issues[0].Resolve == "" {
				t.Fatalf("finding names no resolution: %v", issues[0])
			}
		})
	}
}

// Without a catalog the agent falls back to filenames, so notes are not faults.
func TestCatalogConsistencySkippedWithoutAnIndex(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "context"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "context", "billing.md"), []byte("# Billing\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	issues, err := (&Store{Root: root}).knowledgeIssues()
	if err != nil || len(issues) != 0 {
		t.Fatalf("%v %v", err, issues)
	}
}
