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

// upgraded rewrites a fixture into what an earlier version left behind: the
// given roster, written by hand, and no member selected on this machine.
func upgraded(t *testing.T, f fixture, roster string) {
	t.Helper()
	write(t, filepath.Join(f.root, "members.yaml"), roster)
	if err := os.Remove(filepath.Join(f.root, "member.local.yaml")); err != nil {
		t.Fatal(err)
	}
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

// A person working alone is a team of one, set up once inside init: the roster
// lists them, this machine selects them, and they hold band 1 unless they pick
// another.
func TestInitSetsUpItsMember(t *testing.T) {
	f := setup(t)
	g := fixture{t, filepath.Join(f.home, "alone"), f.home}
	if out := g.fail("init", "--name", "Alone", "--purpose", "Side project"); !strings.Contains(out, "missing --member") {
		t.Fatalf("init without a member was not refused: %s", out)
	}
	g.ok("init", "--name", "Alone", "--purpose", "Side project", "--member", "ana", "--member-name", "Ana")
	status := orientation(t, g.ok("status"))
	if status.ActiveMember != "ana" || status.IdentityRequired != "" || status.BandRequired != "" {
		t.Fatalf("init left setup outstanding: %+v", status)
	}
	if band := status.Members.Members["ana"].Band; band == nil || *band != 1 {
		t.Fatalf("init did not give its member band 1: %v", band)
	}
	if status.NextBand != 2 {
		t.Fatalf("the next band offered should be 2, got %d", status.NextBand)
	}
	if created := g.intent("first"); created.ID != "i100" || created.CreatedBy != "ana" {
		t.Fatalf("the first record should be i100 by ana: %+v", created)
	}
	g.ok("check")
}

// A workspace an earlier version left with an empty roster keeps its records,
// tells the person what to set up, and refuses every change until they have.
func TestAnUpgradedEmptyRosterIsSetUpBeforeAnyChange(t *testing.T) {
	f := setup(t)
	api := filepath.Join(f.home, "api")
	if err := os.Mkdir(api, 0755); err != nil {
		t.Fatal(err)
	}
	git(t, api, "init", "-b", "main")
	// A record an earlier version wrote with no author stays valid.
	write(t, filepath.Join(f.root, "intent", "i001-old.md"), "---\nid: i001\ncreated_at: 2026-09-15T10:53:00Z\n---\n\n# Old\n")
	write(t, filepath.Join(f.root, ".context-circuit", "ids.yaml"), "intents:\n  - i001\nplans: []\n")
	upgraded(t, f, "members: {}\n")

	status := orientation(t, f.ok("status"))
	if !strings.Contains(status.IdentityRequired, "member add --id ID --name NAME --band N") || status.NextBand != 1 {
		t.Fatalf("status did not ask for a member and offer band 1: %+v", status)
	}
	for _, args := range [][]string{
		{"record", "create", "--kind", "intent", "--slug", "x", "--title", "X"},
		{"repo", "connect", "--id", "api", "--path", api, "--base", "main"},
	} {
		if out := f.fail(args...); !strings.Contains(out, "setup is incomplete") || !strings.Contains(out, "member add") {
			t.Fatalf("%v was not refused naming setup: %s", args, out)
		}
	}
	if code, out, _ := call(f.root, "check"); code == 0 || !strings.Contains(out, "select a workspace member") || strings.Contains(out, "unknown created_by") {
		t.Fatalf("check did not run and report the missing member alone: %s", out)
	}

	f.ok("member", "add", "--id", "ana", "--name", "Ana", "--band", "1")
	f.ok("member", "use", "--id", "ana")
	if created := f.intent("next"); created.ID != "i100" || created.CreatedBy != "ana" {
		t.Fatalf("a record made after setup should carry its author: %+v", created)
	}
	f.ok("repo", "connect", "--id", "api", "--path", api, "--base", "main")
	if _, err := os.Stat(filepath.Join(f.root, "intent", "i001-old.md")); err != nil {
		t.Fatal("the old record did not survive setup")
	}
}

// Members an earlier version left without a band are each offered one that
// starts empty, and changes wait until every member holds one.
func TestAMemberWithoutABandIsOfferedOne(t *testing.T) {
	f := setup(t)
	write(t, filepath.Join(f.root, "members.yaml"), "members:\n  maya:\n    name: Maya\n    band: 0\n  rina:\n    name: Rina\n  ben:\n    name: Ben\n")
	// i120 already sits in band 1, so band 1 is not offered.
	write(t, filepath.Join(f.root, ".context-circuit", "ids.yaml"), "intents:\n  - i120\nplans: []\n")

	status := orientation(t, f.ok("status"))
	if status.BandRequired == "" || status.BandSuggestions["ben"] != 2 || status.BandSuggestions["rina"] != 3 || len(status.BandSuggestions) != 2 || status.NextBand != 4 {
		t.Fatalf("unexpected band suggestions: %+v %v next %d", status.BandRequired, status.BandSuggestions, status.NextBand)
	}
	if out := f.fail("record", "create", "--kind", "intent", "--slug", "x", "--title", "X"); !strings.Contains(out, "member band --id ID --band N") || !strings.Contains(out, "ben band 2") {
		t.Fatalf("a change with an unbanded member was not refused naming the step: %s", out)
	}
	findings := checkFindings(t, f.root)
	var named bool
	for _, item := range findings {
		if item.Issue == "member rina holds no band" && strings.Contains(item.Resolve, "`member band --id rina --band 3`") {
			named = true
		}
	}
	if !named {
		t.Fatalf("check did not offer rina band 3: %v", findings)
	}
	f.ok("member", "band", "--id", "rina", "--band", "3")
	f.ok("member", "band", "--id", "ben", "--band", "1")
	if status := orientation(t, f.ok("status")); status.BandRequired != "" {
		t.Fatalf("banding every member did not settle setup: %+v", status)
	}
	// Picked by hand, a band holding old records allocates past them.
	f.ok("member", "use", "--id", "ben")
	if created := f.intent("ben"); created.ID != "i100" {
		t.Fatalf("band 1 allocated %s, want i100", created.ID)
	}
}

// Every command that changes something is refused before setup, while reading
// keeps working so the person can see why.
func TestEveryChangeWaitsForSetup(t *testing.T) {
	f := setup(t)
	intent := f.intent("billing")
	upgraded(t, f, "members:\n  maya:\n    name: Maya\n    band: 0\n")
	path := filepath.Join(f.home, "anywhere")
	for _, args := range [][]string{
		{"record", "create", "--kind", "intent", "--slug", "x", "--title", "X"},
		{"record", "approve", "--id", intent.ID, "--text", "yes"},
		{"record", "complete", "--id", "p0001", "--text", "done"},
		{"record", "dependencies", "--id", "p0001"},
		{"record", "delete", "--all"},
		{"repo", "connect", "--id", "api", "--path", path, "--base", "main"},
		{"repo", "clone", "--id", "api", "--path", path, "--base", "main", "--url", "https://example.invalid/api.git"},
		{"repo", "init", "--id", "api", "--path", path, "--base", "main"},
		{"repo", "base", "--id", "api", "--branch", "main"},
		{"repo", "remote", "--id", "api", "--url", "https://example.invalid/api.git"},
		{"repo", "relate", "--from", "api", "--to", "web", "--description", "calls"},
		{"repo", "fetch", "--id", "api"},
		{"workspace", "connect", "--base", "main"},
		{"workspace", "base", "--branch", "main"},
		{"workspace", "remote", "--url", "https://example.invalid/w.git"},
		{"knowledge", "connect", "--id", "k", "--path", path},
		{"knowledge", "clone", "--id", "k", "--url", "https://example.invalid/k.git"},
		{"knowledge", "remote", "--id", "k", "--url", "https://example.invalid/k.git"},
		{"knowledge", "sync", "--id", "k"},
		{"worktree", "prepare", "--repo", "api"},
		{"worktree", "move", "--repo", "api", "--path", path, "--to", path + "2"},
		{"worktree", "repair", "--repo", "api", "--path", path},
		{"worktree", "remove", "--repo", "api", "--path", path},
		{"worktree", "prune", "--repo", "api"},
		{"agent", "setup"},
		{"agent", "configure", "--host", "claude-code", "--role", "worker", "--model", "inherit", "--effort", "inherit"},
		{"member", "language", "--id", "maya", "--language", "English"},
		{"member", "tone", "--id", "maya", "--clear"},
	} {
		if out := f.fail(args...); !strings.Contains(out, "setup is incomplete") {
			t.Errorf("%v was not refused for setup: %s", args, out)
		}
	}
	for _, args := range [][]string{{"status"}, {"record", "list"}, {"member", "list"}, {"record", "show", "--id", intent.ID}} {
		f.ok(args...)
	}
	if _, out, _ := call(f.root, "check"); !strings.Contains(out, "select a workspace member") {
		t.Fatalf("check did not run: %s", out)
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

func TestInitTakesAWholeMember(t *testing.T) {
	f := setup(t)
	for i, half := range [][]string{{"--member", "ana"}, {"--member-name", "Ana"}} {
		g := fixture{t, filepath.Join(f.home, "half"+string(rune('a'+i))), f.home}
		args := append([]string{"init", "--name", "Half", "--purpose", "Half"}, half...)
		if out := g.fail(args...); !strings.Contains(out, "missing --member") {
			t.Fatalf("%v: half a member was not refused by name: %s", half, out)
		}
	}
}

// Adding a member writes nothing without a band the person chose, and the
// refusal offers the next free one.
func TestAddingAMemberTakesABand(t *testing.T) {
	f := setup(t)
	if out := f.fail("member", "add", "--id", "ana", "--name", "Ana"); !strings.Contains(out, "pass --band N; the next free band is 1") {
		t.Fatalf("member add without --band did not offer the next band: %s", out)
	}
	if strings.Contains(read(t, filepath.Join(f.root, "members.yaml")), "ana") {
		t.Fatal("a refused member add wrote the member")
	}
	f.ok("member", "add", "--id", "ana", "--name", "Ana", "--band", "7")
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

// --all deletes every record, authorless ones from an earlier version
// included, and only while the roster lists exactly one member.
func TestDeletingEverythingNeedsARosterOfOne(t *testing.T) {
	f := team(t)
	if out := f.fail("record", "delete", "--all"); !strings.Contains(out, "--member ID") || !strings.Contains(out, "exactly one member") {
		t.Fatalf("--all with two members was not refused toward --member: %s", out)
	}
	if out := f.fail("record", "delete", "--member", "maya", "--all"); !strings.Contains(out, "not both") {
		t.Fatalf("--member with --all was not refused: %s", out)
	}
	f.fail("record", "delete")

	g := setup(t)
	g.repository("api")
	intent := g.intent("first")
	g.plan(intent.ID, "first-api", "api")
	// A record written by an earlier version, with no author.
	write(t, filepath.Join(g.root, "intent", "i050-old.md"), "---\nid: i050\ncreated_at: 2026-09-15T10:53:00Z\n---\n\n# Old\n")
	l := ledger(t, g)
	write(t, filepath.Join(g.root, ".context-circuit", "ids.yaml"), "intents:\n  - "+strings.Join(append(l.Intents, "i050"), "\n  - ")+"\nplans:\n  - "+strings.Join(l.Plans, "\n  - ")+"\n")
	preview := deletion(t, g.ok("record", "delete", "--all"))
	if len(preview.Records) != 3 || preview.Deleted {
		t.Fatalf("a preview should list every record and delete none: %+v", preview)
	}
	deletion(t, g.ok("record", "delete", "--all", "--confirm"))
	if l := ledger(t, g); len(l.Intents)+len(l.Plans) != 0 {
		t.Fatalf("deleting everything did not reset the ledger: %+v", l)
	}
	if paths, _ := filepath.Glob(filepath.Join(g.root, "intent", "i*")); len(paths) != 0 {
		t.Fatalf("records survived deleting everything: %v", paths)
	}
	if next := g.intent("restart"); next.ID != "i001" {
		t.Fatalf("a reset workspace should count from i001 again: %s", next.ID)
	}
	g.ok("record", "delete", "--member", "maya", "--confirm")
	if out := g.fail("record", "delete", "--all"); !strings.Contains(out, "no intents, plans, or reserved IDs") {
		t.Fatalf("an empty deletion did not say there was nothing to delete: %s", out)
	}
}
