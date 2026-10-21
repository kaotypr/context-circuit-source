package cli_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/kaotypr/context-circuit-source/internal/workspace"
)

// schema2 rewrites a fixture into what a 2.3.0 CLI left behind: schema 2, and
// every reservation in the one shared ID file.
func schema2(t *testing.T, f fixture, intents, plans []string) {
	t.Helper()
	if err := os.RemoveAll(filepath.Join(f.root, ".context-circuit", "ids")); err != nil {
		t.Fatal(err)
	}
	list := func(ids []string) string {
		if len(ids) == 0 {
			return " []\n"
		}
		return "\n  - " + strings.Join(ids, "\n  - ") + "\n"
	}
	write(t, filepath.Join(f.root, ".context-circuit", "ids.yaml"), "intents:"+list(intents)+"plans:"+list(plans))
	config := filepath.Join(f.root, "workspace.yaml")
	write(t, config, strings.Replace(read(t, config), "version: 3", "version: 2", 1))
}

func TestTheFirstChangeSplitsASchema2Workspace(t *testing.T) {
	f := setup(t)
	first := f.intent("first")
	// A record an earlier version wrote with no author, and a number whose
	// record was never written: both belong to whoever splits.
	write(t, filepath.Join(f.root, "intent", "i002-old.md"), "---\nid: i002\ncreated_at: 2026-09-15T10:53:00Z\n---\n\n# Old\n")
	schema2(t, f, []string{first.ID, "i002", "i050"}, nil)
	legacy := filepath.Join(f.root, ".context-circuit", "ids.yaml")
	config := filepath.Join(f.root, "workspace.yaml")

	// Reading, and a deletion preview, change nothing.
	f.ok("status")
	call(f.root, "check")
	deletion(t, f.ok("record", "delete", "--member", "maya"))
	if !exists(legacy) || !strings.Contains(read(t, config), "version: 2") {
		t.Fatal("a command that changes nothing split the ID file")
	}

	code, out, errOut := call(f.root, "record", "create", "--kind", "intent", "--slug", "next", "--title", "Next")
	if code != 0 {
		t.Fatalf("exit=%d %s %s", code, out, errOut)
	}
	if created := record(t, out); created.ID != "i003" {
		t.Fatalf("the first record after the split should be i003: %+v", created)
	}
	if !strings.Contains(errOut, "notice: moved this workspace's ID reservations") || !strings.Contains(errOut, "from 2 to 3") || !strings.Contains(errOut, "person's decision") {
		t.Fatalf("the split was not reported on the error stream: %s", errOut)
	}
	if exists(legacy) || !strings.Contains(read(t, config), "version: 3") {
		t.Fatal("the split did not remove the shared file and move the schema")
	}
	if l := idFile(t, f, ".context-circuit/ids/maya.yaml"); !slices.Equal(l.Intents, []string{first.ID, "i002", "i050", "i003"}) {
		t.Fatalf("maya's file does not hold every reservation: %+v", l)
	}

	// A teammate's older clone merges the shared file back with a new number.
	write(t, legacy, "intents:\n  - "+first.ID+"\n  - i060\nplans: []\n")
	f.ok("status")
	if !exists(legacy) {
		t.Fatal("reading re-split the shared file")
	}
	code, _, errOut = call(f.root, "record", "create", "--kind", "intent", "--slug", "later", "--title", "Later")
	if code != 0 || !strings.Contains(errOut, "Already reserved in a member file: "+first.ID) || !strings.Contains(errOut, "Commit and push") || strings.Contains(errOut, "from 2 to 3") {
		t.Fatalf("the re-split was not reported: %d %s", code, errOut)
	}
	if exists(legacy) {
		t.Fatal("the re-split left the shared file")
	}
	if l := idFile(t, f, ".context-circuit/ids/maya.yaml"); !slices.Contains(l.Intents, "i060") || len(l.Intents) != 6 {
		t.Fatalf("the re-split lost or duplicated a reservation: %+v", l)
	}
	f.ok("check")
}

// A workspace that cannot say who is working has nobody to give authorless
// reservations to, so nothing splits until it can. The selection that completes
// setup then splits it itself, rather than leaving it to the next command.
func TestTheSplitWaitsForSetup(t *testing.T) {
	f := setup(t)
	schema2(t, f, []string{"i050"}, nil)
	upgraded(t, f, "members:\n  maya:\n    name: Maya\n    band: 0\n")
	legacy := filepath.Join(f.root, ".context-circuit", "ids.yaml")
	f.fail("record", "create", "--kind", "intent", "--slug", "early", "--title", "Early")
	if !exists(legacy) {
		t.Fatal("the split ran before a member was selected")
	}
	code, out, errOut := call(f.root, "member", "use", "--id", "maya")
	if code != 0 {
		t.Fatalf("exit=%d %s %s", code, out, errOut)
	}
	if exists(legacy) || !strings.Contains(errOut, "notice: moved this workspace's ID reservations") {
		t.Fatalf("the selection that completed setup did not split: %s", errOut)
	}
	// member.local.yaml is this machine's alone and is never offered for commit.
	if strings.Contains(errOut, "member.local.yaml") || strings.Contains(errOut, "members.yaml") {
		t.Fatalf("the notice named a file the selection did not change for commit: %s", errOut)
	}
	if l := idFile(t, f, ".context-circuit/ids/maya.yaml"); !slices.Equal(l.Intents, []string{"i050"}) {
		t.Fatalf("maya's file does not hold the reservation: %+v", l)
	}
}

// The band that completes setup makes the split in the same command, and the
// notice names every shared file that command changed: the roster it banded, as
// well as what the split wrote and removed. A second call has nothing to split.
func TestTheBandThatCompletesSetupSplits(t *testing.T) {
	f := setup(t)
	schema2(t, f, []string{"i050"}, nil)
	write(t, filepath.Join(f.root, "members.yaml"), "members:\n  maya:\n    name: Maya\n    band: 0\n  alex:\n    name: Alex\n")
	legacy := filepath.Join(f.root, ".context-circuit", "ids.yaml")

	code, out, errOut := call(f.root, "member", "band", "--id", "alex", "--band", "2")
	if code != 0 {
		t.Fatalf("exit=%d %s %s", code, out, errOut)
	}
	if exists(legacy) || !strings.Contains(read(t, filepath.Join(f.root, "workspace.yaml")), "version: 3") {
		t.Fatal("the band that completed setup did not split")
	}
	want := "Commit and push .context-circuit/ids.yaml, .context-circuit/ids/maya.yaml, members.yaml and workspace.yaml so every clone sees the same reservations"
	if !strings.Contains(errOut, "notice: moved this workspace's ID reservations") || !strings.Contains(errOut, "from 2 to 3") || !strings.Contains(errOut, want) {
		t.Fatalf("the notice does not name every file the command changed:\n%s", errOut)
	}
	if l := idFile(t, f, ".context-circuit/ids/maya.yaml"); !slices.Equal(l.Intents, []string{"i050"}) {
		t.Fatalf("maya's file does not hold the reservation: %+v", l)
	}

	code, out, errOut = call(f.root, "member", "band", "--id", "alex", "--band", "3")
	if code != 0 || strings.Contains(errOut, "notice:") {
		t.Fatalf("a second band change split again: exit=%d %s %s", code, out, errOut)
	}
	f.ok("check")
}

// Adding the last member with a band completes setup the same way.
func TestAddingABandedMemberSplitsWhenItCompletesSetup(t *testing.T) {
	f := setup(t)
	schema2(t, f, []string{"i050"}, nil)
	write(t, filepath.Join(f.root, "members.yaml"), "members:\n  maya:\n    name: Maya\n    band: 0\n  alex:\n    name: Alex\n")
	// Refused: alex still has no band, and a refused call splits nothing.
	f.fail("member", "add", "--id", "sam", "--name", "Sam")
	if !exists(filepath.Join(f.root, ".context-circuit", "ids.yaml")) {
		t.Fatal("a refused add split the ID file")
	}
	// Setup is still incomplete after this add, so it does not split either.
	_, _, errOut := call(f.root, "member", "add", "--id", "sam", "--name", "Sam", "--band", "3")
	if strings.Contains(errOut, "notice:") || !exists(filepath.Join(f.root, ".context-circuit", "ids.yaml")) {
		t.Fatalf("an add that left setup incomplete split: %s", errOut)
	}
	code, out, errOut := call(f.root, "member", "band", "--id", "alex", "--band", "2")
	if code != 0 || !strings.Contains(errOut, "members.yaml") || exists(filepath.Join(f.root, ".context-circuit", "ids.yaml")) {
		t.Fatalf("the band completing setup did not split: exit=%d %s %s", code, out, errOut)
	}
}

// A roster complete except for the member being added splits on that add.
func TestAnAddThatCompletesSetupSplits(t *testing.T) {
	f := setup(t)
	schema2(t, f, []string{"i050"}, nil)
	upgraded(t, f, "members:\n  maya:\n    name: Maya\n    band: 0\n")
	write(t, filepath.Join(f.root, "member.local.yaml"), "member: alex\n")
	code, out, errOut := call(f.root, "member", "add", "--id", "alex", "--name", "Alex", "--band", "2")
	if code != 0 {
		t.Fatalf("exit=%d %s %s", code, out, errOut)
	}
	if exists(filepath.Join(f.root, ".context-circuit", "ids.yaml")) || !strings.Contains(errOut, "members.yaml") || !strings.Contains(errOut, ".context-circuit/ids/alex.yaml") {
		t.Fatalf("the add that completed setup did not split and name its files: %s", errOut)
	}
	// The authorless reservation goes to the member selected here.
	if l := idFile(t, f, ".context-circuit/ids/alex.yaml"); !slices.Equal(l.Intents, []string{"i050"}) {
		t.Fatalf("alex's file does not hold the reservation: %+v", l)
	}
}

func TestCheckReadsTheMemberFiles(t *testing.T) {
	f := setup(t)
	mine := f.intent("mine")
	ids := filepath.Join(f.root, ".context-circuit", "ids")
	// A number far outside maya's band is valid in her file.
	write(t, filepath.Join(ids, "maya.yaml"), "intents:\n  - "+mine.ID+"\nplans:\n  - p5000\n")
	write(t, filepath.Join(ids, "alex.yaml"), "intents:\n  - "+mine.ID+"\nplans: []\n")
	write(t, filepath.Join(ids, "x_y.yaml"), "intents: []\nplans: []\n")
	write(t, filepath.Join(ids, "zed.yaml"), "intents:\n  - x1\nplans: []\n")
	write(t, filepath.Join(ids, "broken.yaml"), "intents: [\n")
	findings := checkFindings(t, f.root)
	for _, want := range []string{
		mine.ID + " is reserved in more than one ID file: .context-circuit/ids/alex.yaml, .context-circuit/ids/maya.yaml",
		".context-circuit/ids/x_y.yaml is not named after a member ID",
		`.context-circuit/ids/zed.yaml: invalid reserved ID "x1"`,
		".context-circuit/ids/broken.yaml",
	} {
		var named bool
		for _, item := range findings {
			if strings.Contains(item.Issue, want) && strings.HasPrefix(item.Resolve, "needs a person:") {
				named = true
			}
		}
		if !named {
			t.Errorf("check did not report %q for a person: %v", want, findings)
		}
	}
	for _, item := range findings {
		if strings.Contains(item.Issue, "p5000") || strings.Contains(item.Issue, "reserved in no ID file") {
			t.Errorf("check reported a valid reservation: %v", item)
		}
	}
}

// twoMembers gives maya band 1 and adds alex with band 2.
func twoMembers(f fixture) {
	f.ok("member", "band", "--id", "maya", "--band", "1")
	f.ok("member", "add", "--id", "alex", "--name", "Alex", "--band", "2")
}

// published commits the workspace to a bare origin and returns it.
func published(t *testing.T, f fixture) string {
	t.Helper()
	git(t, f.root, "init", "-b", "main")
	git(t, f.root, "add", "-A")
	git(t, f.root, "commit", "-m", "test: publish the workspace")
	origin := filepath.Join(f.home, "workspace.git")
	git(t, f.home, "clone", "--bare", f.root, origin)
	return origin
}

// member clones the workspace for one member and returns its fixture.
func member(t *testing.T, f fixture, origin, id string) fixture {
	t.Helper()
	g := fixture{t, filepath.Join(f.home, id), f.home}
	git(t, f.home, "clone", origin, g.root)
	g.ok("member", "use", "--id", id)
	return g
}

func commitAll(t *testing.T, g fixture, message string) {
	t.Helper()
	git(t, g.root, "add", "-A")
	git(t, g.root, "commit", "-m", message)
}

func idFindings(t *testing.T, root string) []finding {
	t.Helper()
	var out []finding
	for _, item := range checkFindings(t, root) {
		if strings.Contains(item.Issue, "ID file") || strings.Contains(item.Issue, "reserved") || strings.Contains(item.Issue, "duplicate record ID") {
			out = append(out, item)
		}
	}
	return out
}

// The point of one file per member: two members creating records on their own
// clones merge without a conflict in any ID file.
func TestTwoMembersMergeWithoutConflict(t *testing.T) {
	f := setup(t)
	f.repository("api")
	twoMembers(f)
	origin := published(t, f)
	maya, alex := member(t, f, origin, "maya"), member(t, f, origin, "alex")

	for _, g := range []fixture{maya, alex} {
		intent := g.intent("work")
		g.plan(intent.ID, "work-api", "api")
		g.intent("more")
		commitAll(t, g, "docs(records): add work")
	}
	git(t, maya.root, "push", "origin", "main")
	git(t, alex.root, "-c", "pull.rebase=false", "pull", "--no-edit", "origin", "main")
	if found := idFindings(t, alex.root); len(found) != 0 {
		t.Fatalf("a clean merge left ID findings: %v", found)
	}
	if l := ledger(t, alex); !slices.Equal(l.Intents, []string{"i200", "i201", "i100", "i101"}) || !slices.Equal(l.Plans, []string{"p2000", "p1000"}) {
		t.Fatalf("the merged ID files lost a reservation: %+v", l)
	}
}

// Two clones that each split before pulling the other each give the same
// authorless number to their own member. The merge still succeeds, nothing is
// handed out twice, and check names the duplicate for a person.
func TestTwoClonesSplittingApartLeaveADuplicateForCheck(t *testing.T) {
	f := setup(t)
	twoMembers(f)
	// Written after the roster changes, which would otherwise split it here.
	schema2(t, f, []string{"i050"}, nil)
	origin := published(t, f)
	maya, alex := member(t, f, origin, "maya"), member(t, f, origin, "alex")

	for _, g := range []fixture{maya, alex} {
		g.intent("work")
		commitAll(t, g, "docs(records): split and add work")
	}
	git(t, maya.root, "push", "origin", "main")
	git(t, alex.root, "-c", "pull.rebase=false", "pull", "--no-edit", "origin", "main")
	found := idFindings(t, alex.root)
	if len(found) != 1 || !strings.Contains(found[0].Issue, "i050 is reserved in more than one ID file") {
		t.Fatalf("check did not report exactly the duplicate reservation: %v", found)
	}
	if next := alex.intent("next"); next.ID != "i201" {
		t.Fatalf("allocation after the merge handed out %s", next.ID)
	}
}

// A call its own command refuses must not split on its way to the refusal.
func TestARefusedCallLeavesASchema2WorkspaceUnsplit(t *testing.T) {
	f := setup(t)
	schema2(t, f, []string{"i050"}, nil)
	legacy := filepath.Join(f.root, ".context-circuit", "ids.yaml")
	config := filepath.Join(f.root, "workspace.yaml")
	for _, tc := range []struct {
		args []string
		want string
	}{
		{[]string{"member", "add", "--id", "alex", "--name", "Alex"}, "the next free band is"},
		{[]string{"member", "roles", "--id", "maya"}, "--clear"},
		{[]string{"member", "roles", "--id", "maya", "--role", "docs", "--clear"}, "not both"},
		{[]string{"member", "roles", "--id", "maya", "--role", ""}, "single-line"},
		{[]string{"record", "delete", "--id", "i999", "--confirm"}, "no record or reserved ID matches i999"},
	} {
		if out := f.fail(tc.args...); !strings.Contains(out, tc.want) {
			t.Fatalf("%v: expected a refusal naming %q: %s", tc.args, tc.want, out)
		}
		if !exists(legacy) || !strings.Contains(read(t, config), "version: 2") || exists(filepath.Join(f.root, ".context-circuit", "ids")) {
			t.Fatalf("%v was refused but split the ID file", tc.args)
		}
	}
}

// A confirmed deletion by ID is a change, so on a schema 2 workspace it splits
// first and then releases the number from the member file the split wrote.
func TestAConfirmedDeletionByIDOnASchema2Workspace(t *testing.T) {
	f := setup(t)
	kept := f.intent("kept")
	gone := f.intent("gone")
	schema2(t, f, []string{kept.ID, gone.ID, "i050"}, nil)

	code, out, errOut := call(f.root, "record", "delete", "--id", gone.ID, "--id", "i050", "--confirm")
	if code != 0 {
		t.Fatalf("exit=%d %s %s", code, out, errOut)
	}
	if d := deletion(t, out); !d.Deleted || !slices.Equal(d.Released, []string{gone.ID, "i050"}) {
		t.Fatalf("the deletion did not release what it named: %+v", d)
	}
	if !strings.Contains(errOut, "notice: moved this workspace's ID reservations") {
		t.Fatalf("the split was not reported: %s", errOut)
	}
	if exists(filepath.Join(f.root, ".context-circuit", "ids.yaml")) || !strings.Contains(read(t, filepath.Join(f.root, "workspace.yaml")), "version: 3") {
		t.Fatal("the confirmed deletion did not split the workspace")
	}
	if l := ledger(t, f); !slices.Equal(l.Intents, []string{kept.ID}) {
		t.Fatalf("the member file holds %v after the deletion", l.Intents)
	}
	if exists(filepath.Join(f.root, gone.Path)) || !exists(filepath.Join(f.root, kept.Path)) {
		t.Fatal("the deletion removed the wrong record")
	}
}

// A 2.3.0 roster could hold a member whose ID is a Windows device name. Status
// and check say so before anything fails, and the split names who is at fault
// instead of only refusing.
func TestADeviceNameMemberIsNamedBeforeAndAtTheSplit(t *testing.T) {
	f := setup(t)
	write(t, filepath.Join(f.root, "members.yaml"), "members:\n  maya:\n    name: Maya\n    band: 0\n  con:\n    name: Con\n    band: 1\n")
	write(t, filepath.Join(f.root, "intent", "i100-theirs.md"), "---\nid: i100\ncreated_by: con\ncreated_at: 2026-09-15T10:53:00Z\n---\n\n# Theirs\n")
	write(t, filepath.Join(f.root, "intent", "i002-stray.md"), "---\nid: i002\ncreated_by: aux\ncreated_at: 2026-09-15T10:53:00Z\n---\n\n# Stray\n")
	schema2(t, f, []string{"i002", "i100"}, nil)

	for _, want := range []string{"member con is a Windows device name", "i002: created_by aux is a Windows device name"} {
		var named bool
		for _, item := range checkFindings(t, f.root) {
			if strings.Contains(item.Issue, want) && strings.HasPrefix(item.Resolve, "needs a person:") {
				named = true
			}
		}
		if !named {
			t.Errorf("check did not report %q for a person", want)
		}
	}
	if issues := orientation(t, f.ok("status")).Issues; !slices.ContainsFunc(issues, func(item workspace.Finding) bool {
		return strings.Contains(item.Issue, "member con is a Windows device name") && strings.Contains(item.Resolve, "members.yaml")
	}) {
		t.Errorf("status did not report the device-name member: %v", issues)
	}

	out := f.fail("record", "create", "--kind", "intent", "--slug", "next", "--title", "Next")
	for _, want := range []string{`i002 name "aux" in created_by`, "Windows device name", "members.yaml", "nothing was changed"} {
		if !strings.Contains(out, want) {
			t.Fatalf("the split refusal does not name %q: %s", want, out)
		}
	}
	if !exists(filepath.Join(f.root, ".context-circuit", "ids.yaml")) || exists(filepath.Join(f.root, ".context-circuit", "ids")) {
		t.Fatal("a refused split wrote member files or removed the shared file")
	}

	// The author on the roster is named the same way once the stray one is fixed.
	write(t, filepath.Join(f.root, "intent", "i002-stray.md"), "---\nid: i002\ncreated_by: maya\ncreated_at: 2026-09-15T10:53:00Z\n---\n\n# Stray\n")
	if out := f.fail("record", "create", "--kind", "intent", "--slug", "next", "--title", "Next"); !strings.Contains(out, `i100 name "con" in created_by`) {
		t.Fatalf("the split refusal does not name the roster author: %s", out)
	}

	// The member selected on this machine is named as that.
	upgraded(t, f, "members:\n  con:\n    name: Con\n    band: 1\n")
	// The selection completes setup and is kept; the split it cannot make is a
	// warning, and the next changing command refuses with the reason.
	if code, out, errOut := call(f.root, "member", "use", "--id", "con"); code != 0 || !strings.Contains(errOut, "warning: member use was done") || !strings.Contains(errOut, `"con"`) {
		t.Fatalf("the selection was refused or did not warn: exit=%d %s %s", code, out, errOut)
	}
	if out := f.fail("record", "create", "--kind", "intent", "--slug", "next", "--title", "Next"); !strings.Contains(out, `the member selected on this machine, "con"`) {
		t.Fatalf("the split refusal does not name the selected member: %s", out)
	}
}

// The split writes every author's file, not only the file of the member making
// it. Two clones that each make their first change before pulling the other's
// split therefore both add the same member files with different content, and
// Git stops the merge with an add/add conflict. That is expected, and is why
// one person upgrades first and pushes; this pins it and the documented
// resolution: keep the union of both sides, and check names what both took.
func TestTwoClonesSplittingApartConflictOnEveryAuthorsFile(t *testing.T) {
	f := setup(t)
	twoMembers(f)
	mine := f.intent("mine")
	write(t, filepath.Join(f.root, "intent", "i200-theirs.md"), "---\nid: i200\ncreated_by: alex\ncreated_at: 2026-09-15T10:53:00Z\n---\n\n# Theirs\n")
	schema2(t, f, []string{mine.ID, "i200", "i050"}, nil)
	origin := published(t, f)
	maya, alex := member(t, f, origin, "maya"), member(t, f, origin, "alex")
	for _, g := range []fixture{maya, alex} {
		g.intent("work")
		commitAll(t, g, "docs(records): split and add work")
	}

	for _, file := range []string{".context-circuit/ids/maya.yaml", ".context-circuit/ids/alex.yaml"} {
		if read(t, filepath.Join(maya.root, file)) == read(t, filepath.Join(alex.root, file)) {
			t.Fatalf("both clones wrote %s alike, so nothing here conflicts", file)
		}
	}

	git(t, maya.root, "push", "origin", "main")
	merge := exec.Command("git", "-C", alex.root, "-c", "commit.gpgsign=false", "-c", "pull.rebase=false", "pull", "--no-edit", "origin", "main")
	data, err := merge.CombinedOutput()
	if err == nil || !strings.Contains(string(data), "CONFLICT (add/add)") {
		t.Fatalf("expected an add/add conflict in the member files: %v %s", err, data)
	}

	// The documented resolution: each member file keeps both sides' entries.
	write(t, filepath.Join(alex.root, ".context-circuit", "ids", "maya.yaml"), read(t, filepath.Join(maya.root, ".context-circuit", "ids", "maya.yaml")))
	write(t, filepath.Join(alex.root, ".context-circuit", "ids", "alex.yaml"), "intents:\n  - i200\n  - i050\n  - i201\nplans: []\n")
	commitAll(t, alex, "docs(records): merge the ID split")
	found := idFindings(t, alex.root)
	if len(found) != 1 || !strings.Contains(found[0].Issue, "i050 is reserved in more than one ID file") {
		t.Fatalf("check did not report exactly the number both splits took: %v", found)
	}
}
