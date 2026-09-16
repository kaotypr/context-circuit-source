package cli_test

import (
	"path/filepath"
	"strings"
	"testing"
)

// commitIn writes a file and commits it, so a fixture can advance a branch
// without repeating the identity flags every call.
func commitIn(t *testing.T, repo, name, message string) {
	t.Helper()
	write(t, filepath.Join(repo, name), "content\n")
	git(t, repo, "add", name)
	git(t, repo, "-c", "user.name=Fixture", "-c", "user.email=fixture@example.invalid", "commit", "-m", message)
}

// behind gives repo an origin, publishes branch, advances it there, and leaves
// the local branch one commit short of the counterpart it already fetched.
func behind(t *testing.T, home, repo, branch string) {
	t.Helper()
	bare := filepath.Join(home, filepath.Base(repo)+".git")
	git(t, home, "init", "--bare", "-b", branch, bare)
	git(t, repo, "remote", "add", "origin", bare)
	git(t, repo, "push", "-u", "origin", branch)
	commitIn(t, repo, "advanced.txt", "test: advance the published branch")
	git(t, repo, "push", "origin", branch)
	git(t, repo, "reset", "--hard", "HEAD~1")
}

// A local base branch nobody synchronized produces a worktree missing work that
// already exists, and preparation otherwise reports only which branch it used.
// The transcript this guards against ended in a stash, a conflict, and a rebase.
func TestPreparingFromAStaleBaseSaysSo(t *testing.T) {
	f := setup(t)
	api := f.repository("api")
	behind(t, f.home, api, "main")

	i := f.intent("billing")
	p := f.plan(i.ID, "ledger", "api")
	w := tree(t, f.ok("worktree", "prepare", "--repo", "api", "--plan", p.ID))
	if w.SyncRequired == "" {
		t.Fatalf("preparing from a base 1 commit behind its remote said nothing: %+v", w)
	}
	for _, want := range []string{"1 commit(s) behind", "origin/main", w.Branch} {
		if !strings.Contains(w.SyncRequired, want) {
			t.Errorf("sync_required omits %q: %s", want, w.SyncRequired)
		}
	}
}

// Reporting a branch that is current would train the caller to ignore the field,
// so a synchronized base and a repository with no remote both stay silent.
func TestPreparingFromACurrentBaseIsSilent(t *testing.T) {
	f := setup(t)
	f.repository("api")
	i := f.intent("billing")

	p := f.plan(i.ID, "no-remote", "api")
	if w := tree(t, f.ok("worktree", "prepare", "--repo", "api", "--plan", p.ID)); w.SyncRequired != "" {
		t.Errorf("a repository with no remote reported staleness: %s", w.SyncRequired)
	}

	web := f.repository("web")
	bare := filepath.Join(f.home, "web.git")
	git(t, f.home, "init", "--bare", "-b", "main", bare)
	git(t, web, "remote", "add", "origin", bare)
	git(t, web, "push", "-u", "origin", "main")
	q := f.plan(i.ID, "current", "web")
	if w := tree(t, f.ok("worktree", "prepare", "--repo", "web", "--plan", q.ID)); w.SyncRequired != "" {
		t.Errorf("a base level with its remote reported staleness: %s", w.SyncRequired)
	}
}

// Bands separate members, not one member's own clones. The ledger that would
// show the collision is the file that was not pulled, so the number is reported
// where it is handed out rather than by `check` after both records exist.
func TestAllocatingInABehindWorkspaceSaysSo(t *testing.T) {
	f := setup(t)
	git(t, f.root, "init", "-b", "main")
	git(t, f.root, "add", "-A")
	git(t, f.root, "-c", "user.name=Fixture", "-c", "user.email=fixture@example.invalid", "commit", "-m", "test: seed workspace")
	behind(t, f.home, f.root, "main")

	i := f.intent("billing")
	if i.SyncRequired == "" {
		t.Fatalf("allocating in a workspace 1 commit behind its remote said nothing: %+v", i)
	}
	for _, want := range []string{"1 commit(s) behind", "origin/main", i.ID} {
		if !strings.Contains(i.SyncRequired, want) {
			t.Errorf("sync_required omits %q: %s", want, i.SyncRequired)
		}
	}
	// The warning belongs to the allocation, never to the record it describes.
	if strings.Contains(read(t, filepath.Join(f.root, i.Path)), "sync_required") {
		t.Errorf("the allocation warning was written into the record file")
	}
}

// A workspace that is not a repository, or is level with its remote, allocates
// without comment: most workspaces are one clone and must pay nothing for this.
func TestAllocatingInACurrentWorkspaceIsSilent(t *testing.T) {
	f := setup(t)
	if i := f.intent("billing"); i.SyncRequired != "" {
		t.Errorf("a workspace outside Git reported staleness: %s", i.SyncRequired)
	}
}
