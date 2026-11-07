package cli_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kaotypr/context-circuit-source/internal/workspace"
)

func synced(t *testing.T, out string) []workspace.SyncState {
	t.Helper()
	var result struct {
		Repositories []workspace.SyncState `json:"repositories"`
	}
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatal(err, out)
	}
	return result.Repositories
}

func syncOne(t *testing.T, f fixture, id string) workspace.SyncState {
	t.Helper()
	states := synced(t, f.ok("repo", "sync", "--id", id))
	if len(states) != 1 || states[0].ID != id {
		t.Fatalf("expected one state for %s: %+v", id, states)
	}
	return states[0]
}

// A base branch left behind after a merge is the stale start worktree
// preparation reports; sync closes it as a fast-forward, wherever the branch is.
func TestRepoSyncFastForwardsTheBaseBranch(t *testing.T) {
	f := setup(t)
	api := f.repository("api")
	behind(t, f.home, api, "main")
	remote := git(t, api, "rev-parse", "origin/main")

	got := syncOne(t, f, "api")
	if got.Status != "updated" || got.Commits != 1 || got.To != remote || got.Counterpart != "origin/main" {
		t.Fatalf("a checked-out base 1 behind was not fast-forwarded: %+v", got)
	}
	if head := git(t, api, "rev-parse", "main"); head != remote {
		t.Fatalf("main is at %s, not %s", head, remote)
	}
	if again := syncOne(t, f, "api"); again.Status != "current" {
		t.Errorf("a base at its remote is not current: %+v", again)
	}
}

// With another branch checked out, the base moves in place and the checkout's
// own branch and files stay exactly as they were.
func TestRepoSyncMovesABaseNotCheckedOut(t *testing.T) {
	f := setup(t)
	api := f.repository("api")
	behind(t, f.home, api, "main")
	remote := git(t, api, "rev-parse", "origin/main")
	git(t, api, "checkout", "-b", "feature")
	write(t, filepath.Join(api, "wip.txt"), "keep")

	if got := syncOne(t, f, "api"); got.Status != "updated" {
		t.Fatalf("a base checked out nowhere was not moved: %+v", got)
	}
	if git(t, api, "rev-parse", "main") != remote {
		t.Error("main did not move")
	}
	if branch := git(t, api, "symbolic-ref", "--short", "HEAD"); branch != "feature" {
		t.Errorf("sync switched the checkout to %s", branch)
	}
	if read(t, filepath.Join(api, "wip.txt")) != "keep" {
		t.Error("sync touched uncommitted work")
	}
}

// A base checked out in a linked worktree belongs to that worktree, whose files
// would go stale if the ref moved under it.
func TestRepoSyncLeavesABaseCheckedOutElsewhere(t *testing.T) {
	f := setup(t)
	api := f.repository("api")
	behind(t, f.home, api, "main")
	local := git(t, api, "rev-parse", "main")
	git(t, api, "checkout", "-b", "feature")
	linked := filepath.Join(f.home, "linked")
	git(t, api, "worktree", "add", linked, "main")

	got := syncOne(t, f, "api")
	if got.Status != "checked-out-elsewhere" || !strings.Contains(got.Detail, "linked") {
		t.Fatalf("a base checked out in another worktree was not left alone: %+v", got)
	}
	if git(t, api, "rev-parse", "main") != local {
		t.Error("main moved under a linked worktree")
	}
}

// Uncommitted changes in the way stop the fast-forward and are listed; nothing
// is stashed or overwritten to make room.
func TestRepoSyncIsBlockedByChangesInTheWay(t *testing.T) {
	f := setup(t)
	api := f.repository("api")
	behind(t, f.home, api, "main")
	local := git(t, api, "rev-parse", "main")
	// behind's remote commit adds advanced.txt; an untracked copy blocks it.
	write(t, filepath.Join(api, "advanced.txt"), "mine")

	got := syncOne(t, f, "api")
	if got.Status != "blocked" || !strings.Contains(strings.Join(got.Changes, "\n"), "advanced.txt") {
		t.Fatalf("changes in the way did not block: %+v", got)
	}
	if git(t, api, "rev-parse", "main") != local || read(t, filepath.Join(api, "advanced.txt")) != "mine" {
		t.Error("a blocked sync changed the branch or the file")
	}
}

// Local commits are never rewritten: ahead takes nothing, diverged changes nothing.
func TestRepoSyncNeverRewritesLocalCommits(t *testing.T) {
	f := setup(t)
	api := f.repository("api")
	behind(t, f.home, api, "main")
	commitIn(t, api, "local.txt", "test: local work")
	local := git(t, api, "rev-parse", "main")

	if got := syncOne(t, f, "api"); got.Status != "diverged" {
		t.Fatalf("a diverged base was not reported: %+v", got)
	}
	if git(t, api, "rev-parse", "main") != local {
		t.Fatal("a diverged base moved")
	}

	git(t, api, "reset", "--hard", "origin/main")
	commitIn(t, api, "ahead.txt", "test: unpushed work")
	ahead := git(t, api, "rev-parse", "main")
	if got := syncOne(t, f, "api"); got.Status != "ahead" || got.To != ahead {
		t.Fatalf("a base ahead of its remote was not reported: %+v", got)
	}
}

// With no --id every bound repository is synced, in ID order, and one that
// cannot be reached or has no remote is reported without stopping the rest.
func TestRepoSyncCoversEveryBoundRepository(t *testing.T) {
	f := setup(t)
	api := f.repository("api")
	behind(t, f.home, api, "main")
	f.repository("docs")
	web := f.repository("web")
	git(t, web, "remote", "add", "origin", filepath.Join(f.home, "missing.git"))

	states := synced(t, f.ok("repo", "sync"))
	var got []string
	for _, state := range states {
		got = append(got, state.ID+"="+state.Status)
	}
	if strings.Join(got, " ") != "api=updated docs=no-upstream web=unreachable" {
		t.Fatalf("every bound repository was not reported in order: %v", got)
	}

	if out := f.fail("repo", "sync", "--id", "nope"); !strings.Contains(out, "unknown repository: nope") {
		t.Errorf("an unknown ID was not refused: %s", out)
	}
}

// Sync changes the person's checkouts, so it waits for setup like any change.
func TestRepoSyncWaitsForSetup(t *testing.T) {
	f := setup(t)
	f.repository("api")
	if err := os.Remove(filepath.Join(f.root, "member.local.yaml")); err != nil {
		t.Fatal(err)
	}
	if out := f.fail("repo", "sync"); !strings.Contains(out, "setup is incomplete") {
		t.Errorf("sync ran before setup was complete: %s", out)
	}
}

// The workspace repository is stricter than a code checkout: any uncommitted
// path, untracked included, leaves it unfetched, while repo sync still moves the
// code checkouts. Once it moves, the next step and a moved pin are named.
func TestWorkspaceSyncSkipsADirtyWorkspace(t *testing.T) {
	f := setup(t)
	api := f.repository("api")
	behind(t, f.home, api, "main")
	git(t, f.root, "init", "-b", "main")
	f.ok("workspace", "connect", "--base", "main")
	git(t, f.root, "add", "-A")
	git(t, f.root, "commit", "-m", "test: seed workspace")
	bare := filepath.Join(f.home, "workspace.git")
	git(t, f.home, "init", "--bare", "-b", "main", bare)
	git(t, f.root, "remote", "add", "origin", bare)
	git(t, f.root, "push", "-u", "origin", "main")

	teammate := filepath.Join(f.home, "teammate")
	git(t, f.home, "clone", bare, teammate)
	before := read(t, filepath.Join(f.root, ".context-circuit", "CLI_VERSION"))
	write(t, filepath.Join(teammate, ".context-circuit", "CLI_VERSION"), "9.9.9\n")
	git(t, teammate, "commit", "-am", "test: move the pin")
	git(t, teammate, "push", "origin", "main")

	write(t, filepath.Join(f.root, "notes.txt"), "draft")
	var dirty workspace.WorkspaceSyncResult
	if err := json.Unmarshal([]byte(f.ok("workspace", "sync")), &dirty); err != nil {
		t.Fatal(err)
	}
	if dirty.Status != "dirty" || !strings.Contains(strings.Join(dirty.Changes, "\n"), "notes.txt") {
		t.Fatalf("an untracked file did not keep the workspace back: %+v", dirty)
	}
	if out := git(t, f.root, "rev-parse", "--verify", "--quiet", "origin/main"); out != git(t, f.root, "rev-parse", "main") {
		t.Error("a dirty workspace was fetched")
	}
	if got := syncOne(t, f, "api"); got.Status != "updated" {
		t.Errorf("a dirty workspace stopped the code checkouts: %+v", got)
	}

	if err := os.Remove(filepath.Join(f.root, "notes.txt")); err != nil {
		t.Fatal(err)
	}
	var moved workspace.WorkspaceSyncResult
	if err := json.Unmarshal([]byte(f.ok("workspace", "sync")), &moved); err != nil {
		t.Fatal(err)
	}
	if moved.Status != "updated" || moved.CheckRecommended == "" {
		t.Fatalf("a clean workspace behind its remote did not move: %+v", moved)
	}
	if moved.CLIVersionChanged == nil || moved.CLIVersionChanged.From != strings.TrimSpace(before) || moved.CLIVersionChanged.To != "9.9.9" {
		t.Errorf("the moved pin was not named: %+v", moved.CLIVersionChanged)
	}
}
