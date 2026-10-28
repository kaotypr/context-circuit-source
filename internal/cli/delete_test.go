package cli_test

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/kaotypr/context-circuit-source/internal/cli"
	"github.com/kaotypr/context-circuit-source/internal/workspace"
)

func deletedIDs(d workspace.Deletion) []string {
	ids := []string{}
	for _, r := range d.Records {
		ids = append(ids, r.ID)
	}
	slices.Sort(ids)
	return ids
}

func deletedRecord(t *testing.T, d workspace.Deletion, id string) workspace.DeletedRecord {
	t.Helper()
	for _, r := range d.Records {
		if r.ID == id {
			return r
		}
	}
	t.Fatalf("%s is not in the deletion: %+v", id, d.Records)
	return workspace.DeletedRecord{}
}

// One record goes by its ID, its plan leaves its intent's list, and the number
// it held is the next one handed out.
func TestDeletingOneRecordByID(t *testing.T) {
	f := team(t)
	intent := f.intent("billing")
	plan := f.plan(intent.ID, "billing-api", "api")
	other := f.intent("reporting")
	f.ok("record", "approve", "--id", intent.ID, "--text", "yes")

	preview := deletion(t, f.ok("record", "delete", "--id", other.ID))
	if preview.Deleted || preview.ConfirmationRequired == "" || !slices.Equal(preview.Released, []string{other.ID}) {
		t.Fatalf("a deletion by ID without --confirm must be a preview: %+v", preview)
	}
	if r := deletedRecord(t, preview, other.ID); r.CreatedBy != "maya" || r.Status != "draft" || r.Path != other.Path {
		t.Fatalf("the preview must show the author, status, and path: %+v", r)
	}
	if !exists(filepath.Join(f.root, other.Path)) || !slices.Contains(ledger(t, f).Intents, other.ID) {
		t.Fatal("a preview deleted something")
	}

	done := deletion(t, f.ok("record", "delete", "--id", plan.ID, "--confirm"))
	if !done.Deleted || done.SyncRequired == "" || !slices.Equal(done.Released, []string{plan.ID}) {
		t.Fatalf("a confirmed deletion must say so and release the plan: %+v", done)
	}
	if exists(filepath.Join(f.root, filepath.Dir(plan.Path))) {
		t.Fatal("the plan folder survived its deletion")
	}
	if !slices.Equal(done.Unlinked, []string{plan.ID + " from intent " + intent.ID}) {
		t.Fatalf("the result does not name the unlinked plan: %v", done.Unlinked)
	}
	if listed := record(t, f.ok("record", "show", "--id", intent.ID)); slices.Contains(listed.Plans, plan.ID) || listed.ApprovedAt == "" {
		t.Fatalf("the intent still lists the deleted plan or lost its approval: %+v", listed)
	}
	if again := f.plan(intent.ID, "billing-api-again", "api"); again.ID != plan.ID {
		t.Fatalf("the released number was not handed out again: %s", again.ID)
	}

	// An intent and its plan named together: the link between them is inside
	// the set, so it does not block, and approval does not protect the intent.
	preview = deletion(t, f.ok("record", "delete", "--id", intent.ID, "--id", plan.ID))
	if r := deletedRecord(t, preview, intent.ID); r.Status != "approved" {
		t.Fatalf("the preview must show the intent is approved: %+v", r)
	}
	done = deletion(t, f.ok("record", "delete", "--id", plan.ID, "--id", intent.ID, "--id", intent.ID, "--confirm"))
	if !slices.Equal(deletedIDs(done), []string{intent.ID, plan.ID}) || len(done.Unlinked) != 0 {
		t.Fatalf("the set was not deleted together: %+v", done)
	}
	if exists(filepath.Join(f.root, intent.Path)) {
		t.Fatal("the intent survived its deletion")
	}
	if fresh := f.intent("fresh"); fresh.ID != intent.ID {
		t.Fatalf("the lowest free number was not handed out again: %s", fresh.ID)
	}
	if _, out, _ := call(f.root, "check"); strings.Contains(out, "ledger") || strings.Contains(out, "broken plan link") {
		t.Fatalf("a deletion by ID left the workspace inconsistent:\n%s", out)
	}
}

// By ID, a person reaches what --member cannot: an archived record, another
// member's record, a record with no author, and a reserved number with no file.
func TestDeletingByIDReachesEveryKindOfRecord(t *testing.T) {
	f := team(t)
	archived := f.intent("old-idea")
	archive := filepath.Join(f.root, "intent", "_archive", filepath.Base(archived.Path))
	write(t, archive, read(t, filepath.Join(f.root, archived.Path)))
	os.Remove(filepath.Join(f.root, archived.Path))
	f.ok("member", "use", "--id", "alex")
	theirs := f.intent("reporting")
	f.ok("member", "use", "--id", "maya")
	write(t, filepath.Join(f.root, "intent", "i050-old.md"), "---\nid: i050\ncreated_at: 2026-09-15T10:53:00Z\n---\n\n# Old\n")
	// i250 is reserved with no file, as an interrupted create leaves it, and
	// i100 is held twice, as two clones that split separately can leave it.
	write(t, filepath.Join(f.root, ".context-circuit/ids/alex.yaml"), "intents:\n  - i200\n  - i250\n  - i100\nplans: []\n")

	preview := deletion(t, f.ok("record", "delete", "--id", archived.ID, "--id", theirs.ID, "--id", "i050", "--id", "i250"))
	if want := []string{"i050", "i100", "i200", "i250"}; !slices.Equal(deletedIDs(preview), want) || !slices.Equal(preview.Released, []string{"i100", "i200", "i250"}) {
		t.Fatalf("preview selected %v and releases %v", deletedIDs(preview), preview.Released)
	}
	if r := deletedRecord(t, preview, "i250"); r.Status != "reservation only" || r.Path != "" {
		t.Fatalf("a bare reservation must show as reservation only: %+v", r)
	}
	if r := deletedRecord(t, preview, theirs.ID); r.CreatedBy != "alex" {
		t.Fatalf("the preview must show another member as the author: %+v", r)
	}
	if r := deletedRecord(t, preview, "i050"); r.CreatedBy != "" || r.Status != "draft" {
		t.Fatalf("an authorless record must show no author: %+v", r)
	}
	if r := deletedRecord(t, preview, archived.ID); !strings.HasPrefix(r.Path, "intent/_archive/") {
		t.Fatalf("the archived record's path is not shown: %+v", r)
	}

	deletion(t, f.ok("record", "delete", "--id", archived.ID, "--id", theirs.ID, "--id", "i050", "--id", "i250", "--confirm"))
	for _, gone := range []string{archive, filepath.Join(f.root, theirs.Path), filepath.Join(f.root, "intent", "i050-old.md")} {
		if exists(gone) {
			t.Errorf("%s survived its deletion", gone)
		}
	}
	for _, rel := range []string{".context-circuit/ids/alex.yaml", ".context-circuit/ids/maya.yaml"} {
		if l := idFile(t, f, rel); len(l.Intents)+len(l.Plans) != 0 {
			t.Fatalf("%s kept a released number: %+v", rel, l)
		}
	}
}

func TestDeletingByIDRefusesWhatItCannotFinish(t *testing.T) {
	f := team(t)
	intent := f.intent("billing")
	first := f.plan(intent.ID, "billing-api", "api")
	second := f.plan(intent.ID, "billing-web", "api")
	f.ok("record", "dependencies", "--id", second.ID, "--depends-on", first.ID)
	before := idFiles(t, f)

	for _, c := range []struct {
		args []string
		want string
	}{
		{[]string{"--id", first.ID, "--member", "maya"}, "without --member or --all"},
		{[]string{"--id", first.ID, "--all"}, "without --member or --all"},
		{[]string{"--id", "i999"}, "no record or reserved ID matches i999"},
		{[]string{"--id", "x12", "--confirm"}, `"x12" is not an intent or plan ID`},
		// A dependant outside the selection blocks, and the refusal names what
		// to add.
		{[]string{"--id", first.ID, "--confirm"}, second.ID + " (maya) depends on " + first.ID},
		{[]string{"--id", first.ID, "--confirm"}, "--id " + second.ID + " to this call"},
		// A plan's intent: is not the link that is edited out.
		{[]string{"--id", intent.ID, "--confirm"}, "--id " + first.ID + " --id " + second.ID + " to this call"},
	} {
		out := f.fail(append([]string{"record", "delete"}, c.args...)...)
		if !strings.Contains(out, c.want) {
			t.Errorf("%v: want %q in %s", c.args, c.want, out)
		}
	}
	if idFiles(t, f) != before || !exists(filepath.Join(f.root, first.Path)) || !exists(filepath.Join(f.root, intent.Path)) {
		t.Fatal("a refused deletion changed something")
	}
	if listed := record(t, f.ok("record", "show", "--id", intent.ID)); !slices.Equal(listed.Plans, []string{first.ID, second.ID}) {
		t.Fatalf("a refused deletion unlinked a plan: %v", listed.Plans)
	}
	done := deletion(t, f.ok("record", "delete", "--id", first.ID, "--id", second.ID, "--confirm"))
	if !done.Deleted || len(done.Unlinked) != 2 {
		t.Fatalf("plans named together with their dependency were not deleted: %+v", done)
	}
	if listed := record(t, f.ok("record", "show", "--id", intent.ID)); len(listed.Plans) != 0 {
		t.Fatalf("the intent still lists deleted plans: %v", listed.Plans)
	}
}

// A refused deletion still shows the preview: the records asked for, the ones
// blocking them and how they link, the reason, and the IDs to add. It deletes
// nothing and exits 1, with --confirm or without, in JSON and in YAML.
func TestARefusedDeletionStillShowsItsPreview(t *testing.T) {
	f := team(t)
	intent := f.intent("billing")
	plan := f.plan(intent.ID, "billing-api", "api")
	f.ok("record", "approve", "--id", intent.ID, "--text", "yes")
	before := idFiles(t, f)

	for _, confirm := range [][]string{nil, {"--confirm"}} {
		code, out, errOut := call(f.root, append([]string{"record", "delete", "--id", intent.ID}, confirm...)...)
		if code != 1 {
			t.Fatalf("%v: a refusal must exit 1: exit=%d %s %s", confirm, code, out, errOut)
		}
		d := deletion(t, out)
		if d.Deleted || len(d.Released) != 0 || d.ConfirmationRequired != "" {
			t.Fatalf("%v: a refusal must not read as a deletion or a preview to confirm: %+v", confirm, d)
		}
		if r := deletedRecord(t, d, intent.ID); r.CreatedBy != "maya" || r.Status != "approved" || r.Path != intent.Path {
			t.Fatalf("%v: the requested record is not shown: %+v", confirm, r)
		}
		want := workspace.BlockingRecord{ID: plan.ID, CreatedBy: "maya", Path: filepath.ToSlash(filepath.Dir(plan.Path)) + "/", Status: "open", Link: "belongs to intent " + intent.ID}
		if !slices.Equal(d.BlockedBy, []workspace.BlockingRecord{want}) {
			t.Fatalf("%v: the blocking record is not shown: %+v", confirm, d.BlockedBy)
		}
		if !slices.Equal(d.AddIDs, []string{plan.ID}) || !strings.Contains(d.Refused, "--id "+plan.ID+" to this call") {
			t.Fatalf("%v: the refusal does not name what to add: %+v", confirm, d)
		}
		if !strings.Contains(errOut, "nothing was deleted") {
			t.Fatalf("%v: the reason is not on the error stream: %s", confirm, errOut)
		}
	}

	// A person reading YAML sees the same preview.
	var out, errOut bytes.Buffer
	code := cli.Run(context.Background(), []string{"--workspace", f.root, "record", "delete", "--id", intent.ID}, &out, &errOut, testVersion)
	for _, want := range []string{"blocked_by:", "id: " + plan.ID, "link: belongs to intent " + intent.ID, "add_ids:", "refused:", "status: approved"} {
		if code != 1 || !strings.Contains(out.String(), want) {
			t.Fatalf("the YAML preview lacks %q: exit=%d %s %s", want, code, out.String(), errOut.String())
		}
	}

	if idFiles(t, f) != before || !exists(filepath.Join(f.root, intent.Path)) || !exists(filepath.Join(f.root, plan.Path)) {
		t.Fatal("a refused deletion changed something")
	}
}

// A worktree this machine prepared for the plan refuses, as it does for
// --member. What Git still holds once it is gone is warned about, not refused.
func TestDeletingAPlanByIDWarnsAboutLeftoverBranches(t *testing.T) {
	f := team(t)
	api := filepath.Join(f.home, "api")
	intent := f.intent("billing")
	plan := f.plan(intent.ID, "billing-api", "api")
	w := tree(t, f.ok("worktree", "prepare", "--repo", "api", "--plan", plan.ID))
	if out := f.fail("record", "delete", "--id", plan.ID, "--confirm"); !strings.Contains(out, "prepared for "+plan.ID) || !strings.Contains(out, "worktree remove --repo api --path "+w.Path) {
		t.Fatalf("a worktree prepared for the plan did not refuse: %s", out)
	}
	if !exists(filepath.Join(f.root, plan.Path)) {
		t.Fatal("a refused deletion removed the plan")
	}
	f.ok("worktree", "remove", "--repo", "api", "--path", w.Path)
	git(t, api, "update-ref", "refs/remotes/origin/cc/"+plan.ID+"/api", "HEAD")
	git(t, api, "worktree", "add", "-b", "cc/"+plan.ID+"/spike", filepath.Join(f.home, "spike"), "main")
	if err := os.MkdirAll(filepath.Join(f.root, ".worktrees", plan.ID, "leftover"), 0755); err != nil {
		t.Fatal(err)
	}

	preview := deletion(t, f.ok("record", "delete", "--id", plan.ID))
	done := deletion(t, f.ok("record", "delete", "--id", plan.ID, "--confirm"))
	for _, d := range []workspace.Deletion{preview, done} {
		warnings := strings.Join(d.Warnings, "\n")
		for _, want := range []string{
			"branch cc/" + plan.ID + "/api in api",
			"remote-tracking branch origin/cc/" + plan.ID + "/api in api",
			"is on branch cc/" + plan.ID + "/spike",
			"is left from an earlier " + plan.ID,
		} {
			if !strings.Contains(warnings, want) {
				t.Errorf("deleted=%v: want %q in warnings:\n%s", d.Deleted, want, warnings)
			}
		}
	}
	if !done.Deleted || exists(filepath.Join(f.root, filepath.Dir(plan.Path))) {
		t.Fatalf("a warning stopped the deletion: %+v", done)
	}
	// The deletion never touches a repository.
	git(t, api, "rev-parse", "--verify", "refs/heads/cc/"+plan.ID+"/api")
	if !exists(filepath.Join(f.home, "spike")) {
		t.Fatal("the deletion removed a worktree directory")
	}
	if !strings.Contains(done.SyncRequired, "worktree prune") || !strings.Contains(done.SyncRequired, "`check`") {
		t.Fatalf("the result does not tell other machines what to clean up: %s", done.SyncRequired)
	}
}

// --member warns about leftovers the way --id does.
func TestDeletingByMemberWarnsAboutLeftoverBranches(t *testing.T) {
	f := team(t)
	api := filepath.Join(f.home, "api")
	plan := f.plan(f.intent("billing").ID, "billing-api", "api")
	git(t, api, "branch", "cc/"+plan.ID+"/api", "main")
	d := deletion(t, f.ok("record", "delete", "--member", "maya"))
	if !strings.Contains(strings.Join(d.Warnings, "\n"), "branch cc/"+plan.ID+"/api in api") {
		t.Fatalf("a deletion by member did not warn about the leftover branch: %v", d.Warnings)
	}
}

// A plan number no file holds any more, still listed by its intent, is taken
// out of the list rather than refused.
func TestDeletingAReservationOnlyPlanUnlinksIt(t *testing.T) {
	f := team(t)
	intent := f.intent("billing")
	plan := f.plan(intent.ID, "billing-api", "api")
	if err := os.RemoveAll(filepath.Join(f.root, filepath.Dir(plan.Path))); err != nil {
		t.Fatal(err)
	}
	done := deletion(t, f.ok("record", "delete", "--id", plan.ID, "--confirm"))
	if !slices.Equal(done.Unlinked, []string{plan.ID + " from intent " + intent.ID}) || !slices.Equal(done.Released, []string{plan.ID}) {
		t.Fatalf("the reservation-only plan was not unlinked and released: %+v", done)
	}
	if listed := record(t, f.ok("record", "show", "--id", intent.ID)); len(listed.Plans) != 0 {
		t.Fatalf("the intent still lists the plan: %v", listed.Plans)
	}
}

// A plan given a released number does not quietly resume the worktree an
// earlier plan with that number left behind; the same plan still resumes its own.
func TestPreparingAReusedPlanNumberRefusesTheOldWorktree(t *testing.T) {
	f := team(t)
	plan := f.plan(f.intent("billing").ID, "billing-api", "api")
	w := tree(t, f.ok("worktree", "prepare", "--repo", "api", "--plan", plan.ID))
	write(t, filepath.Join(w.Path, "old.txt"), "old work\n")
	git(t, w.Path, "add", "old.txt")
	git(t, w.Path, "commit", "-m", "old work")
	if again := tree(t, f.ok("worktree", "prepare", "--repo", "api", "--plan", plan.ID)); !again.Reused || again.Path != w.Path {
		t.Fatalf("the same plan did not resume its worktree: %+v", again)
	}

	// Another machine pulled a deletion and the number came back for a new
	// plan: same ID, same path and branch, a later created_at.
	path := filepath.Join(f.root, plan.Path)
	original := read(t, path)
	write(t, path, strings.Replace(original, "created_at: "+string(plan.CreatedAt), "created_at: 2099-01-01T00:00:00Z", 1))
	out := f.fail("worktree", "prepare", "--repo", "api", "--plan", plan.ID)
	for _, want := range []string{"not this plan's to resume", "earlier plan given the number " + plan.ID, "worktree remove --repo api", "--reuse"} {
		if !strings.Contains(out, want) {
			t.Errorf("want %q in the refusal: %s", want, out)
		}
	}

	// An association written before 2.4.0 has no created_at; its commits,
	// made before the new plan existed, say the same.
	assocs := filepath.Join(f.root, ".context-circuit", "local", "worktrees.yaml")
	lines := []string{}
	for line := range strings.SplitSeq(read(t, assocs), "\n") {
		if !strings.Contains(line, "plan_created_at") {
			lines = append(lines, line)
		}
	}
	write(t, assocs, strings.Join(lines, "\n"))
	if out := f.fail("worktree", "prepare", "--repo", "api", "--plan", plan.ID); !strings.Contains(out, "holds commits made before this "+plan.ID) {
		t.Fatalf("a legacy association with older commits was resumed: %s", out)
	}
	// The same plan from before the field existed still resumes, and is
	// recorded from then on.
	write(t, path, original)
	f.ok("worktree", "prepare", "--repo", "api", "--plan", plan.ID)
	if !strings.Contains(read(t, assocs), "plan_created_at: "+string(plan.CreatedAt)) {
		t.Fatalf("the resumed association did not record the plan's created_at:\n%s", read(t, assocs))
	}

	// --reuse takes the old commits deliberately and adopts the worktree.
	write(t, path, strings.Replace(original, "created_at: "+string(plan.CreatedAt), "created_at: 2099-01-01T00:00:00Z", 1))
	if adopted := tree(t, f.ok("worktree", "prepare", "--repo", "api", "--plan", plan.ID, "--reuse")); !adopted.Reused {
		t.Fatalf("--reuse did not select the worktree: %+v", adopted)
	}
	f.ok("worktree", "prepare", "--repo", "api", "--plan", plan.ID)

	// A bare branch left from the earlier plan is named as that.
	f.ok("worktree", "remove", "--repo", "api", "--path", w.Path)
	write(t, path, strings.Replace(original, "created_at: "+string(plan.CreatedAt), "created_at: 2099-01-02T00:00:00Z", 1))
	if out := f.fail("worktree", "prepare", "--repo", "api", "--plan", plan.ID); !strings.Contains(out, "likely left from an earlier plan given the number "+plan.ID) {
		t.Fatalf("a branch from an earlier plan was not named as one: %s", out)
	}
}

// legacyAssociations strips plan_created_at, as an association written before
// 2.4.0 has none.
func legacyAssociations(t *testing.T, f fixture) {
	t.Helper()
	path := filepath.Join(f.root, ".context-circuit", "local", "worktrees.yaml")
	lines := []string{}
	for line := range strings.SplitSeq(read(t, path), "\n") {
		if !strings.Contains(line, "plan_created_at") {
			lines = append(lines, line)
		}
	}
	write(t, path, strings.Join(lines, "\n"))
}

// A worktree from before 2.4.0 whose branch holds no commit has nothing to date
// it by: it resumes only when clean and where a fresh prepare would start it.
// Once recorded, the same plan resumes it with no flag wherever the base is.
func TestPreparingALegacyWorktreeWithNoCommits(t *testing.T) {
	f := team(t)
	api := filepath.Join(f.home, "api")
	plan := f.plan(f.intent("billing").ID, "billing-api", "api")
	w := tree(t, f.ok("worktree", "prepare", "--repo", "api", "--plan", plan.ID))

	legacyAssociations(t, f)
	draft := filepath.Join(w.Path, "draft.txt")
	write(t, draft, "uncommitted\n")
	out := f.fail("worktree", "prepare", "--repo", "api", "--plan", plan.ID)
	for _, want := range []string{"not this plan's to resume", "uncommitted changes", "worktree remove --repo api", "--reuse"} {
		if !strings.Contains(out, want) {
			t.Errorf("want %q in the refusal of a dirty legacy worktree: %s", want, out)
		}
	}
	if err := os.Remove(draft); err != nil {
		t.Fatal(err)
	}

	if again := tree(t, f.ok("worktree", "prepare", "--repo", "api", "--plan", plan.ID)); !again.Reused || again.Path != w.Path {
		t.Fatalf("a clean legacy worktree at the base tip was not resumed: %+v", again)
	}
	if assocs := read(t, filepath.Join(f.root, ".context-circuit", "local", "worktrees.yaml")); !strings.Contains(assocs, "plan_created_at: "+string(plan.CreatedAt)) {
		t.Fatalf("the resumed association did not record the plan's created_at:\n%s", assocs)
	}

	legacyAssociations(t, f)
	git(t, api, "commit", "--allow-empty", "-m", "later base work")
	if out := f.fail("worktree", "prepare", "--repo", "api", "--plan", plan.ID); !strings.Contains(out, "not at the tip of main") {
		t.Fatalf("a legacy worktree behind the base tip was resumed: %s", out)
	}
	f.ok("worktree", "prepare", "--repo", "api", "--plan", plan.ID, "--reuse")
	if again := tree(t, f.ok("worktree", "prepare", "--repo", "api", "--plan", plan.ID)); !again.Reused {
		t.Fatalf("a recorded association for the same plan was not resumed: %+v", again)
	}
}

// check names a worktree whose association was recorded for an earlier plan
// that held the same number, and how to be rid of it.
func TestCheckNamesAWorktreeFromAnEarlierPlan(t *testing.T) {
	f := team(t)
	plan := f.plan(f.intent("billing").ID, "billing-api", "api")
	f.ok("worktree", "prepare", "--repo", "api", "--plan", plan.ID)
	f.ok("check")
	path := filepath.Join(f.root, plan.Path)
	write(t, path, strings.Replace(read(t, path), "created_at: "+string(plan.CreatedAt), "created_at: 2099-01-01T00:00:00Z", 1))
	out := f.fail("check")
	for _, want := range []string{"prepared for an earlier " + plan.ID, "worktree remove --repo api", "cc/" + plan.ID + "/api", "worktree prune --repo api"} {
		if !strings.Contains(out, want) {
			t.Errorf("want %q in the check finding: %s", want, out)
		}
	}
}

// An empty plan folder, as an interrupted deletion can leave, is reported as
// one to remove rather than as a plan.md to restore.
func TestCheckNamesAnEmptyPlanFolder(t *testing.T) {
	f := team(t)
	if err := os.MkdirAll(filepath.Join(f.root, "plans", "p0999-left-over"), 0o755); err != nil {
		t.Fatal(err)
	}
	out := f.fail("check")
	if !strings.Contains(out, "empty plan folder: plans/p0999-left-over") || !strings.Contains(out, "remove the empty folder") || strings.Contains(out, "restore a regular plan.md") {
		t.Fatalf("check did not name the empty plan folder: %s", out)
	}
}

// A clone with only origin's base branch dates a bare plan branch from that
// base, not from the start of history.
func TestABarePlanBranchIsDatedFromOriginsBase(t *testing.T) {
	f := team(t)
	api := filepath.Join(f.home, "api")
	plan := f.plan(f.intent("billing").ID, "billing-api", "api")
	git(t, api, "update-ref", "refs/remotes/origin/trunk", "HEAD")
	git(t, api, "branch", "cc/"+plan.ID+"/api", "HEAD")
	f.ok("repo", "base", "--id", "api", "--branch", "trunk")
	path := filepath.Join(f.root, plan.Path)
	write(t, path, strings.Replace(read(t, path), "created_at: "+string(plan.CreatedAt), "created_at: 2099-01-01T00:00:00Z", 1))
	out := f.fail("worktree", "prepare", "--repo", "api", "--plan", plan.ID)
	if strings.Contains(out, "earlier plan") || !strings.Contains(out, "use --reuse") {
		t.Fatalf("a branch holding nothing past origin's base was dated by its whole history: %s", out)
	}
}
