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

// The refusal rc.9 added is armed by `depends_on`; a plan that builds on another
// without declaring it got a worktree cut from the base, holding none of the
// work it was meant to extend. That is the shape of the missed fix in PR #31.
func TestPreparingBesideUnmergedPlanWorkSaysSo(t *testing.T) {
	f := setup(t)
	f.repository("api")
	i := f.intent("billing")

	first := f.plan(i.ID, "ledger", "api")
	w := tree(t, f.ok("worktree", "prepare", "--repo", "api", "--plan", first.ID))
	commitIn(t, w.Path, "ledger.go", "test: implement the ledger")

	second := f.plan(i.ID, "release", "api")
	if second.UnmergedPlans == "" {
		t.Fatalf("creating a plan beside unmerged work said nothing: %+v", second)
	}
	for _, want := range []string{first.ID, "cc/" + first.ID + "/api", "record dependencies"} {
		if !strings.Contains(second.UnmergedPlans, want) {
			t.Errorf("create-time report omits %q: %s", want, second.UnmergedPlans)
		}
	}

	next := tree(t, f.ok("worktree", "prepare", "--repo", "api", "--plan", second.ID))
	if next.UnmergedPlans == "" {
		t.Fatalf("preparing beside unmerged work said nothing: %+v", next)
	}
	for _, want := range []string{first.ID, "cc/" + first.ID + "/api", "1 commit(s)"} {
		if !strings.Contains(next.UnmergedPlans, want) {
			t.Errorf("prepare-time report omits %q: %s", want, next.UnmergedPlans)
		}
	}
}

// Declaring the dependency is the answer the report asks for, so it must stop
// asking. An explicit --start is the coordinator saying what they meant, and a
// sibling whose work is already in the base is not missing from anything.
func TestUnmergedPlanWorkIsSilentOnceAnswered(t *testing.T) {
	f := setup(t)
	api := f.repository("api")
	i := f.intent("billing")

	first := f.plan(i.ID, "ledger", "api")
	w := tree(t, f.ok("worktree", "prepare", "--repo", "api", "--plan", first.ID))
	commitIn(t, w.Path, "ledger.go", "test: implement the ledger")

	declared := record(t, f.ok("record", "create", "--kind", "plan", "--intent", i.ID,
		"--slug", "declared", "--title", "declared", "--repo", "api", "--depends-on", first.ID))
	if declared.UnmergedPlans != "" {
		t.Errorf("a declared dependency was still reported as unmerged: %s", declared.UnmergedPlans)
	}

	// An explicit start is honored without comment.
	other := f.plan(i.ID, "explicit", "api")
	if t2 := tree(t, f.ok("worktree", "prepare", "--repo", "api", "--plan", other.ID,
		"--start", "cc/"+first.ID+"/api")); t2.UnmergedPlans != "" {
		t.Errorf("an explicit --start was questioned: %s", t2.UnmergedPlans)
	}

	// Once the work is in the base branch, nothing is missing from it.
	git(t, api, "merge", "--no-ff", "--no-edit", "-m", "test: land the ledger", "cc/"+first.ID+"/api")
	landed := f.plan(i.ID, "after-merge", "api")
	if landed.UnmergedPlans != "" {
		t.Errorf("merged work was reported as unmerged: %s", landed.UnmergedPlans)
	}
	if t3 := tree(t, f.ok("worktree", "prepare", "--repo", "api", "--plan", landed.ID)); t3.UnmergedPlans != "" {
		t.Errorf("merged work was reported at preparation: %s", t3.UnmergedPlans)
	}
}

// With several unmerged siblings there is no single branch to start from: the
// correct start is a base plus integration merges, which `record order` derives
// from declared dependencies and cannot derive from none. Naming one branch
// would quietly exclude the rest, so the several case names the declaration.
func TestSeveralUnmergedPlansNameTheDeclarationNotABranch(t *testing.T) {
	f := setup(t)
	f.repository("api")
	i := f.intent("billing")

	for _, slug := range []string{"ledger", "invoices"} {
		p := f.plan(i.ID, slug, "api")
		w := tree(t, f.ok("worktree", "prepare", "--repo", "api", "--plan", p.ID))
		commitIn(t, w.Path, slug+".go", "test: implement "+slug)
	}

	third := f.plan(i.ID, "release", "api")
	w := tree(t, f.ok("worktree", "prepare", "--repo", "api", "--plan", third.ID))
	if w.UnmergedPlans == "" {
		t.Fatalf("two unmerged siblings said nothing: %+v", w)
	}
	for _, want := range []string{"p0001", "p0002", "record dependencies --id " + third.ID, "record order"} {
		if !strings.Contains(w.UnmergedPlans, want) {
			t.Errorf("several-sibling report omits %q: %s", want, w.UnmergedPlans)
		}
	}
	// Naming one branch as the start is the failure this wording exists to
	// avoid: following it would exclude the other sibling entirely.
	if strings.Contains(w.UnmergedPlans, "Start from") {
		t.Errorf("several siblings were answered with one arbitrary start: %s", w.UnmergedPlans)
	}
}

// The two reasons a schema does not match need different answers. A workspace
// written by a newer CLI is this CLI being old, and answering it with migration
// advice sends the reader to rewrite a workspace that is not broken.
func TestSchemaMismatchSaysWhichSideIsBehind(t *testing.T) {
	f := setup(t)
	config := filepath.Join(f.root, "workspace.yaml")
	current := read(t, config)

	write(t, config, strings.Replace(current, "version: 2", "version: 3", 1))
	ahead := f.fail("status")
	if !strings.Contains(ahead, "requires a newer CLI") {
		t.Errorf("a schema this CLI is too old for read as a migration problem: %s", ahead)
	}
	if strings.Contains(ahead, "v1 migration") {
		t.Errorf("a forward schema was answered with migration advice: %s", ahead)
	}

	write(t, config, strings.Replace(current, "version: 2", "version: 1", 1))
	if behind := f.fail("status"); !strings.Contains(behind, "v1 migration is not automatic") {
		t.Errorf("an older workspace lost its migration answer: %s", behind)
	}
}
