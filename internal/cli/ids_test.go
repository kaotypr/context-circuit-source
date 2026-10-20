package cli_test

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
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
// reservations to, so nothing splits until it can.
func TestTheSplitWaitsForSetup(t *testing.T) {
	f := setup(t)
	schema2(t, f, []string{"i050"}, nil)
	upgraded(t, f, "members:\n  maya:\n    name: Maya\n    band: 0\n")
	f.ok("member", "use", "--id", "maya")
	if !exists(filepath.Join(f.root, ".context-circuit", "ids.yaml")) {
		t.Fatal("the split ran before a member was selected")
	}
	f.intent("next")
	if exists(filepath.Join(f.root, ".context-circuit", "ids.yaml")) {
		t.Fatal("the first change after setup did not split")
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
