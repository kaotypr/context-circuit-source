package cli_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// tryGit runs a command expected to fail, so the failure can be asserted rather
// than ending the test the way the shared git helper would.
func tryGit(path string, args ...string) (string, error) {
	cmd := exec.Command("git", append([]string{"-C", path}, args...)...)
	data, err := cmd.CombinedOutput()
	return string(data), err
}

// upstream builds a repository standing in for an organization's knowledge
// center: something this workspace reads and another team owns.
func (f fixture) upstream(name string, files map[string]string) string {
	f.t.Helper()
	path := filepath.Join(f.home, name)
	if err := os.Mkdir(path, 0755); err != nil {
		f.t.Fatal(err)
	}
	git(f.t, path, "init", "-b", "main")
	for rel, body := range files {
		write(f.t, filepath.Join(path, filepath.FromSlash(rel)), body)
	}
	git(f.t, path, "add", "-A")
	git(f.t, path, "-c", "user.name=Fixture", "-c", "user.email=fixture@example.invalid", "commit", "-m", "docs: seed knowledge")
	return path
}

func (f fixture) advance(path, rel, body string) {
	f.t.Helper()
	write(f.t, filepath.Join(path, filepath.FromSlash(rel)), body)
	git(f.t, path, "add", "-A")
	git(f.t, path, "-c", "user.name=Fixture", "-c", "user.email=fixture@example.invalid", "commit", "-m", "docs: advance knowledge")
}

func TestBorrowedKnowledgeIsMountedReadOnlyAndSearched(t *testing.T) {
	f := setup(t)
	origin := f.upstream("sheknows", map[string]string{
		"index.md":               "# Knowledge\n\n- [File service](services/file-svc.md) — uploads, retention\n",
		"services/file-svc.md":   "# file-svc\n\nUploads go through the signed URL flow.\n",
		"services/ums-notify.md": "# UMS notifications\n\nInsert writes in-app rows only; it does not send email.\n",
	})

	f.ok("knowledge", "clone", "--id", "sheknows", "--url", origin)
	local := filepath.Join(f.root, "knowledge", "sheknows")
	if _, err := os.Stat(filepath.Join(local, "index.md")); err != nil {
		t.Fatalf("expected the knowledge checkout at the default path: %v", err)
	}

	// The shared record carries the mount; the path stays local to this machine.
	if shared := read(t, filepath.Join(f.root, "workspace.yaml")); !strings.Contains(shared, "knowledge_repositories:") || !strings.Contains(shared, "index: index.md") {
		t.Fatalf("expected a shared knowledge record with its detected index:\n%s", shared)
	}
	if local := read(t, filepath.Join(f.root, "repositories.local.yaml")); !strings.Contains(local, "knowledge:") {
		t.Fatalf("expected a local knowledge binding:\n%s", local)
	}

	// Read-only is enforced by Git rather than asked for in prose.
	if pushURL := git(t, local, "remote", "get-url", "--push", "origin"); pushURL != "READ-ONLY-KNOWLEDGE-REPOSITORY" {
		t.Fatalf("expected the push URL disabled, got %q", pushURL)
	}
	f.advance(local, "services/file-svc.md", "# file-svc\n\nEdited locally, which must not reach anyone.\n")
	if out, err := tryGit(local, "push", "origin", "main"); err == nil {
		t.Fatalf("expected the push to fail against a disabled URL: %s", out)
	}

	// Retrieval reads the borrowed index beside this workspace's own catalog,
	// and says which side a match came from.
	found := f.ok("context", "find", "--query", "file-svc")
	if !strings.Contains(found, `"borrowed": true`) || !strings.Contains(found, "sheknows/index.md") {
		t.Fatalf("expected a borrowed match naming its source:\n%s", found)
	}
}

func TestBorrowedKnowledgeIsNeverNamedAsKnowledgeToReconcile(t *testing.T) {
	f := setup(t)
	origin := f.upstream("sheknows", map[string]string{
		"index.md":   "# Knowledge\n\n- [Billing contract](billing.md) {api} — invoice fields · billing, invoice · reviewed 2026-01-05\n",
		"billing.md": "# Billing contract\n",
	})
	f.ok("knowledge", "clone", "--id", "sheknows", "--url", origin)
	f.repository("api")

	// The borrowed entry names {api} in this workspace's own vocabulary, and is
	// still not something completion may ask anyone here to move.
	out := f.ok("context", "find", "--repo", "api")
	if strings.Contains(out, "Billing contract") {
		t.Fatalf("a borrowed entry must never be offered as knowledge to reconcile:\n%s", out)
	}
}

func TestKnowledgeSyncFastForwardsOnlyWhenItIsSafe(t *testing.T) {
	f := setup(t)
	origin := f.upstream("sheknows", map[string]string{"index.md": "# Knowledge\n"})
	f.ok("knowledge", "clone", "--id", "sheknows", "--url", origin)
	local := filepath.Join(f.root, "knowledge", "sheknows")

	if out := f.ok("knowledge", "sync", "--id", "sheknows"); !strings.Contains(out, `"status": "current"`) {
		t.Fatalf("expected a clean checkout to report current:\n%s", out)
	}

	f.advance(origin, "index.md", "# Knowledge\n\n- [New](new.md) — added upstream\n")
	if out := f.ok("knowledge", "sync", "--id", "sheknows"); !strings.Contains(out, `"status": "updated"`) {
		t.Fatalf("expected a fast-forward:\n%s", out)
	}
	if body := read(t, filepath.Join(local, "index.md")); !strings.Contains(body, "added upstream") {
		t.Fatalf("expected the upstream change to arrive:\n%s", body)
	}

	// Local work is reported and kept, never merged away.
	write(t, filepath.Join(local, "scratch.md"), "# in progress\n")
	out := f.ok("knowledge", "sync", "--id", "sheknows")
	if !strings.Contains(out, `"status": "dirty"`) || !strings.Contains(out, "nothing was changed") {
		t.Fatalf("expected a dirty checkout to be reported and left alone:\n%s", out)
	}
	if _, err := os.Stat(filepath.Join(local, "scratch.md")); err != nil {
		t.Fatalf("a dirty file must survive a sync: %v", err)
	}
	if err := os.Remove(filepath.Join(local, "scratch.md")); err != nil {
		t.Fatal(err)
	}

	// Commits made here, with the upstream also moved, are a divergence this
	// command has no standing to resolve.
	f.advance(local, "local-only.md", "# written in the wrong checkout\n")
	f.advance(origin, "index.md", "# Knowledge\n\n- [Newer](newer.md) — added upstream again\n")
	out = f.ok("knowledge", "sync", "--id", "sheknows")
	if !strings.Contains(out, `"status": "diverged"`) || !strings.Contains(out, "nothing was changed or discarded") {
		t.Fatalf("expected a diverged checkout to be reported and left alone:\n%s", out)
	}
	if _, err := os.Stat(filepath.Join(local, "local-only.md")); err != nil {
		t.Fatalf("a local commit must survive a sync: %v", err)
	}
}

func TestKnowledgeAndRepositoryIDsShareOneNamespace(t *testing.T) {
	f := setup(t)
	origin := f.upstream("sheknows", map[string]string{"index.md": "# Knowledge\n"})
	f.ok("knowledge", "clone", "--id", "sheknows", "--url", origin)

	// An anchor names a repository by ID, so one ID may not mean two things.
	other := f.upstream("other", map[string]string{"README.md": "# Other\n"})
	if out := f.fail("repo", "connect", "--id", "sheknows", "--path", other, "--base", "main"); !strings.Contains(out, "already names a knowledge repository") {
		t.Fatalf("expected the ID clash to be refused:\n%s", out)
	}
	api := f.repository("api")
	if out := f.fail("knowledge", "connect", "--id", "api", "--path", origin); !strings.Contains(out, "already names a repository where work happens") {
		t.Fatalf("expected the reverse clash to be refused:\n%s", out)
	}
	// The same checkout cannot be both, either.
	if out := f.fail("knowledge", "connect", "--id", "api-docs", "--path", api); !strings.Contains(out, "already connected as repository api") {
		t.Fatalf("expected a work checkout to be refused as knowledge:\n%s", out)
	}
}

func TestCheckReportsBorrowedKnowledgeThisMachineLacks(t *testing.T) {
	f := setup(t)
	origin := f.upstream("sheknows", map[string]string{"index.md": "# Knowledge\n"})
	f.ok("knowledge", "clone", "--id", "sheknows", "--url", origin)

	// A clone carries the shared record and none of this machine's state, which
	// is what joining a workspace looks like from the second machine.
	local := read(t, filepath.Join(f.root, "repositories.local.yaml"))
	write(t, filepath.Join(f.root, "repositories.local.yaml"), strings.SplitN(local, "knowledge:", 2)[0])

	code, out, _ := call(f.root, "check")
	if code == 0 {
		t.Fatalf("expected check to report the missing knowledge checkout:\n%s", out)
	}
	if !strings.Contains(out, "is not obtained on this machine") || !strings.Contains(out, "knowledge clone --id sheknows") {
		t.Fatalf("expected the finding to name what to run:\n%s", out)
	}
}
