package cli_test

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

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
	archive := filepath.Join(f.root, "intent", "archive", filepath.Base(archived.Path))
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
	if r := deletedRecord(t, preview, archived.ID); !strings.HasPrefix(r.Path, "intent/archive/") {
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
}
