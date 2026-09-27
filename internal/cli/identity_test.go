package cli_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/kaotypr/context-circuit-source/internal/workspace"
)

// solo initializes a workspace with no member, which is what a solo workspace
// is: an empty roster, nobody to select, and records with no author.
func solo(t *testing.T) fixture {
	t.Helper()
	f := setup(t)
	g := fixture{t, filepath.Join(f.home, "solo"), f.home}
	g.ok("init", "--name", "Solo", "--purpose", "Side project")
	return g
}

func orientation(t *testing.T, out string) workspace.Orientation {
	t.Helper()
	var o workspace.Orientation
	if err := json.Unmarshal([]byte(out), &o); err != nil {
		t.Fatal(err)
	}
	return o
}

func deletion(t *testing.T, out string) workspace.Deletion {
	t.Helper()
	var d workspace.Deletion
	if err := json.Unmarshal([]byte(out), &d); err != nil {
		t.Fatal(err)
	}
	return d
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func TestASoloWorkspaceAsksNobodyWhoTheyAre(t *testing.T) {
	f := solo(t)
	if exists(filepath.Join(f.root, "member.local.yaml")) {
		t.Fatal("a solo workspace selected a member nobody named")
	}
	if roster := read(t, filepath.Join(f.root, "members.yaml")); !strings.Contains(roster, "members: {}") {
		t.Fatalf("a solo workspace recorded a roster:\n%s", roster)
	}
	status := orientation(t, f.ok("status"))
	if status.IdentityRequired != "" || status.ActiveMember != "" {
		t.Fatalf("a solo workspace asked for an identity: %+v", status)
	}
	created := f.intent("first")
	if created.ID != "i001" || created.CreatedBy != "" {
		t.Fatalf("a solo record should be i001 with no author: %+v", created)
	}
	if text := read(t, filepath.Join(f.root, created.Path)); strings.Contains(text, "created_by") {
		t.Fatalf("a solo record carries an author field:\n%s", text)
	}

	// The first member turns it into a team workspace. The record written
	// before that stays valid, and this machine now has to say who it is.
	f.ok("member", "add", "--id", "ana", "--name", "Ana", "--band", "1")
	status = orientation(t, f.ok("status"))
	if status.IdentityRequired == "" {
		t.Fatal("a workspace that now lists members did not ask this machine who is working")
	}
	if out := f.fail("record", "create", "--kind", "intent", "--slug", "second", "--title", "Second"); !strings.Contains(out, "member use --id ID") {
		t.Fatalf("an unidentified create did not name the remedy: %s", out)
	}
	_, out, _ := call(f.root, "check")
	if strings.Contains(out, "unknown created_by") {
		t.Fatalf("a record written before the roster existed is reported as unattributed:\n%s", out)
	}
	f.ok("member", "use", "--id", "ana")
	if status = orientation(t, f.ok("status")); status.IdentityRequired != "" || status.ActiveMember != "ana" {
		t.Fatalf("selecting a member did not settle the identity: %+v", status)
	}
	if next := f.intent("second"); next.ID != "i100" || next.CreatedBy != "ana" {
		t.Fatalf("an identified member should allocate from their band: %+v", next)
	}
}

func TestAnUnidentifiedCloneIsToldToSayWhoItIs(t *testing.T) {
	f := setup(t)
	if status := orientation(t, f.ok("status")); status.IdentityRequired != "" {
		t.Fatalf("the member who initialized the workspace was asked who they are: %+v", status)
	}
	local := filepath.Join(f.root, "member.local.yaml")
	for name, state := range map[string]func(){
		"absent":  func() { os.Remove(local) },
		"unknown": func() { write(t, local, "member: stranger\n") },
	} {
		state()
		status := orientation(t, f.ok("status"))
		if !strings.Contains(status.IdentityRequired, "member use --id ID") || !strings.Contains(status.IdentityRequired, "member add") {
			t.Fatalf("%s: identity_required does not name both remedies: %q", name, status.IdentityRequired)
		}
		f.fail("record", "create", "--kind", "intent", "--slug", "x", "--title", "X")
	}
	f.ok("member", "use", "--id", "maya")
	if status := orientation(t, f.ok("status")); status.IdentityRequired != "" {
		t.Fatalf("selecting an existing member did not settle the identity: %+v", status)
	}
}

func TestInitTakesAMemberWholeOrNotAtAll(t *testing.T) {
	f := setup(t)
	for i, half := range [][]string{{"--member", "ana"}, {"--member-name", "Ana"}} {
		g := fixture{t, filepath.Join(f.home, "half"+string(rune('a'+i))), f.home}
		args := append([]string{"init", "--name", "Half", "--purpose", "Half"}, half...)
		if out := g.fail(args...); !strings.Contains(out, "together") {
			t.Fatalf("%v: half a member was not refused by name: %s", half, out)
		}
	}
}

// team gives maya band 1 and alex band 2 on one machine, switching between
// them with member use, which is how two clones look to the ledger.
func team(t *testing.T) fixture {
	t.Helper()
	f := setup(t)
	f.repository("api")
	f.ok("member", "band", "--id", "maya", "--band", "1")
	f.ok("member", "add", "--id", "alex", "--name", "Alex", "--band", "2")
	return f
}

func ledger(t *testing.T, f fixture) workspace.Ledger {
	t.Helper()
	var l workspace.Ledger
	if err := workspace.Decode([]byte(read(t, filepath.Join(f.root, ".context-circuit/ids.yaml"))), &l); err != nil {
		t.Fatal(err)
	}
	return l
}

func TestDeletingAMembersRecordsReleasesTheirNumbers(t *testing.T) {
	f := team(t)
	kept := f.intent("billing")
	plan := f.plan(kept.ID, "billing-api", "api")
	archived := f.intent("old-idea")
	archive := filepath.Join(f.root, "intent/archive", filepath.Base(archived.Path))
	write(t, archive, read(t, filepath.Join(f.root, archived.Path)))
	os.Remove(filepath.Join(f.root, archived.Path))
	f.ok("member", "use", "--id", "alex")
	other := f.intent("reporting")

	preview := deletion(t, f.ok("record", "delete", "--member", "maya"))
	if preview.Deleted || preview.ConfirmationRequired == "" {
		t.Fatalf("a deletion without --confirm must be a preview naming the confirmation owed: %+v", preview)
	}
	ids := []string{}
	for _, r := range preview.Records {
		ids = append(ids, r.ID)
	}
	slices.Sort(ids)
	if want := []string{"i100", "i101", "p1000"}; !slices.Equal(ids, want) || !slices.Equal(preview.Released, want) {
		t.Fatalf("preview selected %v and releases %v, want %v including the archived intent", ids, preview.Released, want)
	}
	if !exists(filepath.Join(f.root, kept.Path)) || !slices.Contains(ledger(t, f).Intents, "i100") {
		t.Fatal("a preview deleted something")
	}

	done := deletion(t, f.ok("record", "delete", "--member", "maya", "--confirm"))
	if !done.Deleted || done.SyncRequired == "" {
		t.Fatalf("a confirmed deletion must say so and name the sync it leaves outstanding: %+v", done)
	}
	for _, gone := range []string{kept.Path, filepath.Dir(plan.Path), archive} {
		if !filepath.IsAbs(gone) {
			gone = filepath.Join(f.root, gone)
		}
		if exists(gone) {
			t.Errorf("%s survived its deletion", gone)
		}
	}
	if !exists(filepath.Join(f.root, other.Path)) {
		t.Fatal("another member's record was deleted")
	}
	if l := ledger(t, f); !slices.Equal(l.Intents, []string{"i200"}) || len(l.Plans) != 0 {
		t.Fatalf("the ledger should keep only alex's reservation: %+v", l)
	}
	f.ok("member", "use", "--id", "maya")
	if again := f.intent("fresh"); again.ID != "i100" {
		t.Fatalf("a released number was not handed out again: %s", again.ID)
	}
	if _, out, _ := call(f.root, "check"); strings.Contains(out, "ledger") || strings.Contains(out, "broken plan link") {
		t.Fatalf("a completed deletion left the workspace inconsistent:\n%s", out)
	}
	if out := f.fail("record", "delete", "--member", "nobody"); !strings.Contains(out, "no intents or plans were created by member nobody") {
		t.Fatalf("deleting for a member with no records did not say so: %s", out)
	}
}

func TestDeletionRefusesToStrandAnotherMembersWork(t *testing.T) {
	f := team(t)
	intent := f.intent("billing")
	f.ok("member", "use", "--id", "alex")
	theirs := f.plan(intent.ID, "billing-web", "api")
	out := f.fail("record", "delete", "--member", "maya", "--confirm")
	if !strings.Contains(out, theirs.ID+" (alex) belongs to intent "+intent.ID) || !strings.Contains(out, "nothing was deleted") {
		t.Fatalf("a plan left under a deleted intent was not named: %s", out)
	}
	if !exists(filepath.Join(f.root, intent.Path)) || !slices.Contains(ledger(t, f).Intents, intent.ID) {
		t.Fatal("a refused deletion changed something")
	}
	// Alex's own plan is linked from maya's intent, so deleting alex strands
	// that link the other way round.
	if out := f.fail("record", "delete", "--member", "alex"); !strings.Contains(out, intent.ID+" (maya) links plan "+theirs.ID) {
		t.Fatalf("an intent left linking a deleted plan was not named: %s", out)
	}
}

func TestDeletionRefusesAPlanThisMachineIsWorkingOn(t *testing.T) {
	f := team(t)
	intent := f.intent("billing")
	plan := f.plan(intent.ID, "billing-api", "api")
	w := tree(t, f.ok("worktree", "prepare", "--repo", "api", "--plan", plan.ID))
	out := f.fail("record", "delete", "--member", "maya")
	if !strings.Contains(out, "prepared for "+plan.ID) || !strings.Contains(out, "worktree remove --repo api --path "+w.Path) {
		t.Fatalf("a worktree prepared for a deleted plan was not named with its removal: %s", out)
	}
}

func TestDeletionMatchesTheWorkspaceItIsIn(t *testing.T) {
	f := team(t)
	if out := f.fail("record", "delete", "--all"); !strings.Contains(out, "--member ID") {
		t.Fatalf("--all in a team workspace was not refused toward --member: %s", out)
	}
	if out := f.fail("record", "delete", "--member", "maya", "--all"); !strings.Contains(out, "not both") {
		t.Fatalf("--member with --all was not refused: %s", out)
	}
	f.fail("record", "delete")

	g := solo(t)
	g.repository("api")
	intent := g.intent("first")
	g.plan(intent.ID, "first-api", "api")
	g.intent("second")
	if out := g.fail("record", "delete", "--member", "maya"); !strings.Contains(out, "--all") {
		t.Fatalf("--member in a solo workspace was not refused toward --all: %s", out)
	}
	preview := deletion(t, g.ok("record", "delete", "--all"))
	if len(preview.Records) != 3 || preview.Deleted {
		t.Fatalf("a solo preview should list every record and delete none: %+v", preview)
	}
	deletion(t, g.ok("record", "delete", "--all", "--confirm"))
	if l := ledger(t, g); len(l.Intents)+len(l.Plans) != 0 {
		t.Fatalf("a solo deletion did not reset the ledger: %+v", l)
	}
	if paths, _ := filepath.Glob(filepath.Join(g.root, "intent", "i*")); len(paths) != 0 {
		t.Fatalf("records survived a solo deletion: %v", paths)
	}
	if next := g.intent("restart"); next.ID != "i001" {
		t.Fatalf("a reset solo workspace should count from i001 again: %s", next.ID)
	}
	g.ok("record", "delete", "--all", "--confirm")
	if out := g.fail("record", "delete", "--all"); !strings.Contains(out, "no intents, plans, or reserved IDs") {
		t.Fatalf("an empty solo deletion did not say there was nothing to delete: %s", out)
	}
}
