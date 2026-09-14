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
		{"empty flow", "# shared\nrepositories: {}\n", []string{"repositories", "api"}, Repository{"main"}, []string{"# shared", "api", "main"}},
		{"populated flow", "repositories: {api: {base_branch: main}} # shared\n", []string{"repositories", "web"}, Repository{"develop"}, []string{"# shared", "api", "main", "web", "develop"}},
		{"replace flow mapping", "bindings: {api: {path: old}} # local\n", []string{"bindings", "api"}, Binding{"new path"}, []string{"# local", "new path"}},
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
			issues, err := store.knowledgeIssues()
			if err != nil {
				t.Fatal(err)
			}
			if tc.flagged != (len(issues) > 0) {
				t.Fatalf("flagged=%v for %q: %v", tc.flagged, tc.line, issues)
			}
			if tc.flagged && !strings.HasPrefix(issues[0], "context/domains/billing.md:3:") {
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
