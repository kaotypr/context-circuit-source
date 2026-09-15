package cli_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/kaotypr/context-circuit-source/internal/workspace"
)

func order(t *testing.T, out string) workspace.Order {
	t.Helper()
	var result workspace.Order
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatal(err)
	}
	return result
}

// diamond builds the five-plan graph that a single repository makes hardest:
// two independent roots, two fan-ins, one repository throughout.
func diamond(f fixture) {
	f.repository("api")
	f.ok("record", "create", "--kind", "intent", "--slug", "billing", "--title", "Billing")
	plan := func(slug string, dependencies ...string) {
		args := []string{"record", "create", "--kind", "plan", "--slug", slug, "--title", slug,
			"--intent", "i001", "--repo", "api"}
		for _, dependency := range dependencies {
			args = append(args, "--depends-on", dependency)
		}
		f.ok(args...)
	}
	plan("one")
	plan("two")
	plan("three", "p0001", "p0002")
	plan("four", "p0002")
	plan("five", "p0003", "p0004")
}

func TestOrderRecommendsLinearForSingleRepositoryFanIn(t *testing.T) {
	f := setup(t)
	diamond(f)
	result := order(t, f.ok("record", "order", "--intent", "i001"))

	if result.Selected != "linear" || result.Strategy.Recommended != "linear" {
		t.Fatalf("expected linear, got %s (recommended %s)", result.Selected, result.Strategy.Recommended)
	}
	if result.Strategy.Repositories != 1 || result.Strategy.FanIns != 2 {
		t.Fatalf("expected one repository and two fan-ins, got %+v", result.Strategy)
	}
	if result.Strategy.ConflictSurface != "high" {
		t.Fatalf("single repository with parallel plans is a high conflict surface, got %s", result.Strategy.ConflictSurface)
	}
	if result.Strategy.Waves.Layers != 3 || result.Strategy.Waves.IntegrationMerges != 2 {
		t.Fatalf("expected three waves costing two merges, got %+v", result.Strategy.Waves)
	}
	if result.Strategy.Linear.IntegrationMerges != 0 || len(result.Chain) != 5 {
		t.Fatalf("a chain never needs a merge: %+v %d", result.Strategy.Linear, len(result.Chain))
	}
	if len(result.Layers) != 0 || len(result.Integration) != 0 {
		t.Fatal("linear selection must not also emit waves")
	}
	// Every link stacks on its predecessor, so each one already contains its
	// own dependencies and no integration merge is ever required.
	want := []string{"main", "cc/p0001/api", "cc/p0002/api", "cc/p0003/api", "cc/p0004/api"}
	for i, plan := range result.Chain {
		if got := plan.Start["api"].Base; got != want[i] {
			t.Fatalf("%s starts at %s, want %s", plan.ID, got, want[i])
		}
		if len(plan.Start["api"].Merge) != 0 {
			t.Fatalf("%s asked for a merge in linear mode", plan.ID)
		}
	}
}

func TestOrderWavesReportIntegrationPoints(t *testing.T) {
	f := setup(t)
	diamond(f)
	result := order(t, f.ok("record", "order", "--intent", "i001", "--mode", "waves"))

	if result.Selected != "waves" || result.Strategy.Recommended != "linear" {
		t.Fatalf("an explicit mode selects without changing the recommendation: %+v", result.Strategy)
	}
	if len(result.Layers) != 3 {
		t.Fatalf("expected three waves, got %d", len(result.Layers))
	}
	if ids := waveIDs(result.Layers[0]); strings.Join(ids, ",") != "p0001,p0002" {
		t.Fatalf("wave 1 should hold both roots, got %v", ids)
	}
	if ids := waveIDs(result.Layers[1]); strings.Join(ids, ",") != "p0003,p0004" {
		t.Fatalf("wave 2 should hold both dependents, got %v", ids)
	}
	// Two predecessors in the same repository need a base the CLI will not
	// invent; it names the branches and leaves the merge to the coordinator.
	three := result.Layers[1].Plans[0]
	if three.Start["api"].Base != "cc/p0001/api" || strings.Join(three.Start["api"].Merge, ",") != "cc/p0002/api" {
		t.Fatalf("p0003 needs an integration base, got %+v", three.Start["api"])
	}
	// One predecessor is ordinary ancestry, not an integration point.
	if four := result.Layers[1].Plans[1]; len(four.Start["api"].Merge) != 0 || four.Start["api"].Base != "cc/p0002/api" {
		t.Fatalf("p0004 stacks directly on p0002, got %+v", four.Start["api"])
	}
	if len(result.Integration) != 2 {
		t.Fatalf("expected two integration points, got %+v", result.Integration)
	}
	if result.Integration[0].Plan != "p0003" || result.Integration[1].Plan != "p0005" {
		t.Fatalf("unexpected integration points: %+v", result.Integration)
	}
	// Plans sharing a repository inside one wave cannot corrupt each other, but
	// they will diverge before delivery; the report has to say so.
	if strings.Join(three.Shares, ",") != "p0004" {
		t.Fatalf("p0003 shares api with p0004, got %v", three.Shares)
	}
}

func TestOrderHoldsDependentsUntilCompletion(t *testing.T) {
	f := setup(t)
	diamond(f)
	f.ok("record", "complete", "--id", "p0001", "--text", "User requested completion.")

	result := order(t, f.ok("record", "order", "--intent", "i001", "--mode", "waves"))
	if strings.Join(result.Completed, ",") != "p0001" || result.Remaining != 4 {
		t.Fatalf("expected p0001 complete and four remaining, got %v/%d", result.Completed, result.Remaining)
	}
	// p0002 is unfinished, so everything behind it stays held. That is failure
	// propagation: it falls out of the graph instead of being tracked.
	if ids := waveIDs(result.Layers[0]); strings.Join(ids, ",") != "p0002" {
		t.Fatalf("only p0002 is ready, got %v", ids)
	}
	if ids := waveIDs(result.Layers[1]); strings.Join(ids, ",") != "p0003,p0004" {
		t.Fatalf("dependents stay held behind p0002, got %v", ids)
	}

	f.ok("record", "complete", "--id", "p0002", "--text", "User requested completion.")
	result = order(t, f.ok("record", "order", "--intent", "i001", "--mode", "waves"))
	if ids := waveIDs(result.Layers[0]); strings.Join(ids, ",") != "p0003,p0004" {
		t.Fatalf("completing p0002 releases both dependents, got %v", ids)
	}
	if result.Remaining != 3 {
		t.Fatalf("expected three remaining, got %d", result.Remaining)
	}
}

func TestRecordCompleteWritesMachineReadableDate(t *testing.T) {
	f := setup(t)
	diamond(f)
	f.ok("record", "complete", "--id", "p0001", "--text", "User requested completion.")

	shown := record(t, f.ok("record", "show", "--id", "p0001"))
	if shown.CompletedAt == "" {
		t.Fatal("completion must be readable without parsing prose")
	}
	// The human-readable section is what a person reads; neither replaces the other.
	if !strings.Contains(shown.Content, "## Completion — "+string(shown.CompletedAt)) {
		t.Fatalf("completion note missing from body: %s", shown.Content)
	}
	// Instants are written bare, so every one in a record reads the same way
	// whether it sits in frontmatter or in prose.
	if !strings.Contains(shown.Content, "completed_at: "+string(shown.CompletedAt)+"\n") {
		t.Fatalf("completion instant is not bare: %s", shown.Content)
	}
	if _, err := time.Parse(workspace.TimeLayout, string(shown.CompletedAt)); err != nil {
		t.Fatalf("completion instant is not YYYY-MM-DD HH:mm: %v", err)
	}
}

func TestOrderRecommendsWavesAcrossRepositories(t *testing.T) {
	f := setup(t)
	f.repository("api")
	f.repository("web")
	f.ok("record", "create", "--kind", "intent", "--slug", "billing", "--title", "Billing")
	f.ok("record", "create", "--kind", "plan", "--slug", "api-work", "--title", "API",
		"--intent", "i001", "--repo", "api")
	f.ok("record", "create", "--kind", "plan", "--slug", "web-work", "--title", "Web",
		"--intent", "i001", "--repo", "web")
	f.ok("record", "create", "--kind", "plan", "--slug", "web-follow", "--title", "Follow",
		"--intent", "i001", "--repo", "web", "--depends-on", "p0001")

	result := order(t, f.ok("record", "order", "--intent", "i001"))
	if result.Selected != "waves" {
		t.Fatalf("independent repositories overlap for free, got %s: %s", result.Selected, result.Strategy.Reason)
	}
	if result.Strategy.ConflictSurface != "none" {
		t.Fatalf("disjoint repositories per wave have no conflict surface, got %s", result.Strategy.ConflictSurface)
	}
	// A cross-repository dependency is ordering only: p0003 is in web, its
	// predecessor is in api, so there is no branch to build on.
	follow := result.Layers[1].Plans[0]
	if follow.ID != "p0003" || follow.Start["web"].Base != "main" {
		t.Fatalf("cross-repository dependency must not invent a base: %+v", follow)
	}
}

func TestOrderHoldsPlansBehindUnfinishedOutsideWork(t *testing.T) {
	f := setup(t)
	diamond(f)
	f.ok("record", "create", "--kind", "intent", "--slug", "later", "--title", "Later")
	f.ok("record", "create", "--kind", "plan", "--slug", "outside", "--title", "Outside",
		"--intent", "i002", "--repo", "api", "--depends-on", "p0001")

	result := order(t, f.ok("record", "order", "--intent", "i002"))
	if len(result.Blocked) != 1 || result.Blocked[0].Plan != "p0006" {
		t.Fatalf("a dependency in another intent must hold its dependent, got %+v", result.Blocked)
	}
	if result.Remaining != 0 {
		t.Fatalf("nothing is runnable, got %d", result.Remaining)
	}

	f.ok("record", "complete", "--id", "p0001", "--text", "User requested completion.")
	result = order(t, f.ok("record", "order", "--intent", "i002"))
	if len(result.Blocked) != 0 || result.Remaining != 1 {
		t.Fatalf("completing the outside dependency releases it, got %+v/%d", result.Blocked, result.Remaining)
	}
}

func TestOrderRejectsUnknownMode(t *testing.T) {
	f := setup(t)
	diamond(f)
	if out := f.fail("record", "order", "--intent", "i001", "--mode", "parallel"); !strings.Contains(out, "auto, waves, or linear") {
		t.Fatalf("unexpected error: %s", out)
	}
}

func TestDispatchBriefsWorkerOnPlanAncestry(t *testing.T) {
	f := setup(t)
	diamond(f)
	var spec workspace.Dispatch
	out := f.ok("agent", "dispatch", "--host", "claude-code", "--role", "worker",
		"--plan", "p0003", "--path", f.root, "--task", "Implement p0003 tasks in order.")
	if err := json.Unmarshal([]byte(out), &spec); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"sole owner of this working directory",
		"Plan p0003 of intent i001",
		"Depends on: p0001, p0002",
		"expected in this branch's ancestry",
		"invisible in a separate worktree",
		"The full plan follows and is authoritative",
		// A brief names the file in the host's own path syntax, so the
		// expectation is built the same way rather than assuming a separator.
		filepath.Join("plans", "p0003-three.md"),
	} {
		if !strings.Contains(spec.Prompt, want) {
			t.Fatalf("brief missing %q:\n%s", want, spec.Prompt)
		}
	}
	if !spec.LaunchRequired {
		t.Fatal("a specification is not a launched agent")
	}
	// An agent reads the brief top to bottom with nothing else in context, so
	// where the work happens precedes what the work is.
	if !strings.HasPrefix(spec.Prompt, "# Working directory\n\n"+spec.WorkingDirectory) {
		t.Fatalf("a brief opens with its working directory:\n%s", spec.Prompt)
	}
	if strings.Index(spec.Prompt, "# Plan p0003") > strings.Index(spec.Prompt, "\n# Task\n") {
		t.Fatalf("the plan precedes the task it assigns:\n%s", spec.Prompt)
	}
	// Quoting the record is what stops a worker from searching the repository
	// for a plan file it was never handed.
	for _, heading := range []string{"## Approach", "## Tasks and order", "## Risks and checks"} {
		if !strings.Contains(spec.Prompt, heading) {
			t.Fatalf("brief missing quoted %q:\n%s", heading, spec.Prompt)
		}
	}
	if strings.Contains(spec.Prompt, "## Progress and result") {
		t.Fatalf("a plan holds the plan; progress is not a template section:\n%s", spec.Prompt)
	}
}

// The quoted record is the assignment, so a role that receives one needs no task
// to restate it. The exception is a worker taking one slice of a plan several
// workers share, which is named in the task and nowhere else.
func TestARecordAssignsTheRolesItBriefs(t *testing.T) {
	f := setup(t)
	diamond(f)
	f.ok("record", "approve", "--id", "i001", "--text", "User approved the billing outcome.")

	// A planner plans the whole intent, so there is nothing left to narrow.
	if out := f.fail("agent", "dispatch", "--host", "claude-code", "--role", "planner",
		"--intent", "i001", "--path", f.root, "--task", "Plan it."); !strings.Contains(out, "takes no --task") {
		t.Fatalf("unexpected error: %s", out)
	}

	whole := dispatch(t, f, "agent", "dispatch", "--host", "claude-code", "--role", "worker",
		"--plan", "p0003", "--path", f.root)
	// The quoted plan carries a "## Tasks and order" heading of its own, so the
	// section is matched exactly rather than by substring.
	if strings.Contains(whole.Prompt, "\n# Task\n") {
		t.Fatalf("a worker owning its whole plan is assigned by the plan:\n%s", whole.Prompt)
	}
	if !strings.Contains(whole.Prompt, "The full plan follows and is authoritative") {
		t.Fatalf("dropping the task must not drop the plan:\n%s", whole.Prompt)
	}
	// The prohibition the coordinator kept restating belongs to the role.
	for _, want := range []string{"open a pull request", "merge into the base branch"} {
		if !strings.Contains(whole.Prompt, want) {
			t.Fatalf("worker instructions missing %q:\n%s", want, whole.Prompt)
		}
	}

	slice := dispatch(t, f, "agent", "dispatch", "--host", "claude-code", "--role", "worker",
		"--plan", "p0003", "--path", f.root, "--shared", "--task", "Own the ledger package only.")
	if !strings.Contains(slice.Prompt, "# Task\n\nOwn the ledger package only.") {
		t.Fatalf("a shared worker's slice is named nowhere else:\n%s", slice.Prompt)
	}

	// A role with no record to read still states its assignment.
	if out := f.fail("agent", "dispatch", "--host", "claude-code", "--role", "explorer",
		"--path", f.root); !strings.Contains(out, "nonempty") {
		t.Fatalf("unexpected error: %s", out)
	}
}

// A working directory with no statement of what may be done in it leaves how to
// read it open, and a host that also offers GUI automation and web search will
// sometimes answer that question badly.
func TestAReadOnlyBriefSaysHowToRead(t *testing.T) {
	f := setup(t)
	diamond(f)
	f.ok("record", "approve", "--id", "i001", "--text", "User approved the billing outcome.")
	spec := dispatch(t, f, "agent", "dispatch", "--host", "codex", "--role", "planner",
		"--intent", "i001", "--path", f.root)
	for _, want := range []string{
		"# Ownership",
		"You are reading, not changing",
		"rather than driving another application",
		"consult no external source",
	} {
		if !strings.Contains(spec.Prompt, want) {
			t.Fatalf("read-only brief missing %q:\n%s", want, spec.Prompt)
		}
	}
	// Reading the repository is the work; only checks are off limits.
	if !strings.Contains(spec.Prompt, "Run no tests, linters, or builds") {
		t.Fatalf("a planner is told which running is forbidden:\n%s", spec.Prompt)
	}
	if strings.Contains(spec.Prompt, "You cannot run anything") {
		t.Fatalf("no host is known to enforce that, so it states conduct instead:\n%s", spec.Prompt)
	}
	// An implementer's ownership is unchanged by any of this.
	worker := dispatch(t, f, "agent", "dispatch", "--host", "codex", "--role", "worker",
		"--plan", "p0003", "--path", f.root)
	if !strings.Contains(worker.Prompt, "sole owner of this working directory") {
		t.Fatalf("worker ownership changed:\n%s", worker.Prompt)
	}
}

func dispatch(t *testing.T, f fixture, args ...string) workspace.Dispatch {
	t.Helper()
	var spec workspace.Dispatch
	if err := json.Unmarshal([]byte(f.ok(args...)), &spec); err != nil {
		t.Fatal(err)
	}
	return spec
}

// A planner returns the plan shape, so handing it a numbered plan settles the
// split it was dispatched to propose, and an intent that turns out to hold
// several plans has no single ID to pass.
func TestPlannerDispatchesAgainstAnApprovedIntent(t *testing.T) {
	f := setup(t)
	diamond(f)
	if out := f.fail("agent", "dispatch", "--host", "claude-code", "--role", "planner",
		"--plan", "p0001", "--path", f.root); !strings.Contains(out, "not an already numbered --plan") {
		t.Fatalf("unexpected error: %s", out)
	}
	// Planning an unapproved intent plans an outcome nobody agreed to.
	if out := f.fail("agent", "dispatch", "--host", "claude-code", "--role", "planner",
		"--intent", "i001", "--path", f.root); !strings.Contains(out, "not approved") {
		t.Fatalf("unexpected error: %s", out)
	}
	f.ok("record", "approve", "--id", "i001", "--text", "User approved the billing outcome.")
	var spec workspace.Dispatch
	out := f.ok("agent", "dispatch", "--host", "claude-code", "--role", "planner",
		"--intent", "i001", "--path", f.root)
	if err := json.Unmarshal([]byte(out), &spec); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"# Intent i001",
		filepath.Join("intent", "i001-billing.md"),
		"The full intent follows and is authoritative",
		"none is yours to number",
	} {
		if !strings.Contains(spec.Prompt, want) {
			t.Fatalf("planner brief missing %q:\n%s", want, spec.Prompt)
		}
	}
	if strings.Contains(spec.Prompt, "# Plan p") {
		t.Fatalf("a planner is handed no plan number:\n%s", spec.Prompt)
	}
	// --intent is the planner's pairing; a worker implements a numbered plan.
	if out := f.fail("agent", "dispatch", "--host", "claude-code", "--role", "worker",
		"--intent", "i001", "--path", f.root, "--task", "Implement."); !strings.Contains(out, "other roles take --plan") {
		t.Fatalf("unexpected error: %s", out)
	}
}

func TestDispatchSharedWorktreeChangesOwnership(t *testing.T) {
	f := setup(t)
	diamond(f)
	var spec workspace.Dispatch
	out := f.ok("agent", "dispatch", "--host", "claude-code", "--role", "worker", "--shared",
		"--path", f.root, "--task", "Implement the assigned part of p0001.")
	if err := json.Unmarshal([]byte(out), &spec); err != nil {
		t.Fatal(err)
	}
	// Workers sharing one worktree really can overwrite each other; workers in
	// separate worktrees cannot. The brief must not claim the wrong one.
	if !strings.Contains(spec.Prompt, "Preserve their edits") {
		t.Fatalf("shared dispatch needs shared-ownership language:\n%s", spec.Prompt)
	}
	if strings.Contains(spec.Prompt, "sole owner") {
		t.Fatalf("shared dispatch must not claim sole ownership:\n%s", spec.Prompt)
	}
	if out := f.fail("agent", "dispatch", "--host", "claude-code", "--role", "explorer", "--shared",
		"--path", f.root, "--task", "Investigate."); !strings.Contains(out, "workers sharing one worktree") {
		t.Fatalf("unexpected error: %s", out)
	}
}

func waveIDs(wave workspace.OrderWave) []string {
	var ids []string
	for _, plan := range wave.Plans {
		ids = append(ids, plan.ID)
	}
	return ids
}

// Role definitions are gitignored, so a clone of an initialized workspace has
// none of them however the creating machine got them. A specification then names
// an agent type the host cannot resolve, and saying so turns a dispatch that
// launches nothing into a fact the caller can act on.
func TestDispatchReportsMissingRoleDefinition(t *testing.T) {
	f := setup(t)
	diamond(f)
	for _, host := range workspace.Hosts {
		directory := "." + host
		if host == "claude-code" {
			directory = ".claude"
		}
		if err := os.RemoveAll(filepath.Join(f.root, directory, "agents")); err != nil {
			t.Fatal(err)
		}
	}
	f.ok("record", "approve", "--id", "i001", "--text", "User approved the billing outcome.")
	read := func() workspace.Dispatch {
		t.Helper()
		var spec workspace.Dispatch
		out := f.ok("agent", "dispatch", "--host", "claude-code", "--role", "planner",
			"--intent", "i001", "--path", f.root)
		if err := json.Unmarshal([]byte(out), &spec); err != nil {
			t.Fatal(err)
		}
		return spec
	}

	before := read()
	if before.DefinitionInstalled {
		t.Fatal("a clone carries no gitignored definition, so none can be installed")
	}
	if before.DefinitionPath != ".claude/agents/cc-planner.md" {
		t.Fatalf("definition path: %q", before.DefinitionPath)
	}
	if !strings.Contains(before.SetupRequired, "agent setup --host claude-code") {
		t.Fatalf("missing remedy: %q", before.SetupRequired)
	}
	// Reporting is not refusing: the prompt still carries what a live spawn tool
	// needs when native roles are unavailable.
	if before.Prompt == "" || !before.LaunchRequired {
		t.Fatal("an absent definition must not empty the specification")
	}

	f.ok("agent", "setup", "--host", "claude-code")
	after := read()
	if !after.DefinitionInstalled {
		t.Fatal("setup wrote the definition; dispatch still reports it missing")
	}
	if after.SetupRequired != "" {
		t.Fatalf("remedy offered for an installed definition: %q", after.SetupRequired)
	}
	// Another host's definitions are separate; setup for one says nothing of it.
	var other workspace.Dispatch
	out := f.ok("agent", "dispatch", "--host", "cursor", "--role", "planner",
		"--intent", "i001", "--path", f.root)
	if err := json.Unmarshal([]byte(out), &other); err != nil {
		t.Fatal(err)
	}
	if other.DefinitionInstalled {
		t.Fatal("claude-code setup must not install cursor definitions")
	}
}

// Initialization covers every host, because the same workspace is opened in
// more than one and the host that created it is not the host that plans or
// executes in it.
func TestInitInstallsRolesForEveryHost(t *testing.T) {
	f := setup(t)
	for _, want := range []string{
		".codex/agents/cc-planner.toml",
		".claude/agents/cc-planner.md",
		".cursor/agents/cc-planner.md",
		".claude/agents/cc-explorer.md",
		".claude/agents/cc-worker.md",
		".claude/agents/cc-reviewer.md",
	} {
		if _, err := os.Stat(filepath.Join(f.root, want)); err != nil {
			t.Errorf("init left %s missing: %v", want, err)
		}
	}
	// Definitions written by init record an inventory, so a later version
	// replaces them instead of reading them as a customization.
	if out := f.ok("agent", "setup"); !strings.Contains(out, "cc-planner") {
		t.Fatalf("setup after init: %s", out)
	}
}

// A base branch is what work starts from and delivers back to. The shared
// configuration records only default_branch, which is a plausible and wrong
// answer, so every shape describing a repository or worktree reports the
// recorded base rather than leaving a coordinator to find it.
func TestTheRecordedBaseBranchIsReported(t *testing.T) {
	f := setup(t)
	api := f.repository("api")
	git(t, api, "branch", "release")
	f.ok("repo", "base", "--id", "api", "--branch", "release")

	var snapshot workspace.Snapshot
	if err := json.Unmarshal([]byte(f.ok("repo", "inspect", "--id", "api")), &snapshot); err != nil {
		t.Fatal(err)
	}
	if snapshot.BaseBranch != "release" {
		t.Fatalf("repo inspect reported base %q", snapshot.BaseBranch)
	}
	// The checked-out branch and the base differ routinely; both are reported.
	if snapshot.Branch == snapshot.BaseBranch {
		t.Fatalf("fixture cannot distinguish current from base: %+v", snapshot)
	}

	var orientation workspace.Orientation
	if err := json.Unmarshal([]byte(f.ok("status")), &orientation); err != nil {
		t.Fatal(err)
	}
	if got := orientation.Repositories["api"].BaseBranch; got != "release" {
		t.Fatalf("status reported base %q", got)
	}

	// Delivery happens from a worktree, so it carries the base too.
	i := f.intent("billing")
	p := f.plan(i.ID, "api", "api")
	w := tree(t, f.ok("worktree", "prepare", "--repo", "api", "--plan", p.ID))
	if w.BaseBranch != "release" {
		t.Fatalf("prepared worktree did not name its base: %+v", w)
	}
	var inspected workspace.Snapshot
	if err := json.Unmarshal([]byte(f.ok("worktree", "inspect", "--repo", "api", "--path", w.Path)), &inspected); err != nil {
		t.Fatal(err)
	}
	if inspected.BaseBranch != "release" {
		t.Fatalf("worktree inspect reported base %q", inspected.BaseBranch)
	}
}

// Completion returns candidates to judge, and a caller just told "completed"
// treats the request as finished. The outstanding act is named beside the data.
func TestCompletionNamesTheOutstandingReconciliation(t *testing.T) {
	f := setup(t)
	f.repository("api")
	i := f.intent("billing")
	p := f.plan(i.ID, "api", "api")
	var done map[string]any
	if err := json.Unmarshal([]byte(f.ok("record", "complete", "--id", p.ID, "--text", "Delivered.")), &done); err != nil {
		t.Fatal(err)
	}
	if done["completed"] != p.ID {
		t.Fatalf("completion did not name the plan: %+v", done)
	}
	candidates, _ := done["knowledge_candidates"].([]any)
	required, _ := done["reconcile_required"].(string)
	if len(candidates) > 0 && required == "" {
		t.Fatalf("candidates returned with no outstanding act named: %+v", done)
	}
	if len(candidates) == 0 && required != "" {
		t.Fatalf("nothing to reconcile, yet an act was named: %+v", done)
	}
}

// Approval is the gate that hands work back instead of finishing it, so a bare
// success reads as the request being done. It names what it authorized.
func TestApprovalNamesWhatItAuthorized(t *testing.T) {
	f := setup(t)
	i := f.intent("billing")
	var approved map[string]any
	if err := json.Unmarshal([]byte(f.ok("record", "approve", "--id", i.ID, "--text", "Go ahead.")), &approved); err != nil {
		t.Fatal(err)
	}
	if approved["approved"] != i.ID {
		t.Fatalf("approval did not name the intent: %+v", approved)
	}
	required, _ := approved["planning_required"].(string)
	if !strings.Contains(required, "plan") {
		t.Fatalf("approval named no outstanding act: %+v", approved)
	}
}
