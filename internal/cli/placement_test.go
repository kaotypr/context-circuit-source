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

// move renames a record entry, or a plan's whole folder, to the slash path to
// (the entry path it will have), creating the folders it needs.
func move(t *testing.T, f fixture, r workspace.Record, to string) workspace.Record {
	t.Helper()
	from, dst := filepath.Join(f.root, filepath.FromSlash(r.Path)), filepath.Join(f.root, filepath.FromSlash(to))
	if filepath.Base(from) == "plan.md" {
		from, dst = filepath.Dir(from), filepath.Dir(dst)
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(from, dst); err != nil {
		t.Fatal(err)
	}
	r.Path = to
	return r
}

// listed is the IDs record list returns with the given flags.
func listed(t *testing.T, f fixture, args ...string) []string {
	t.Helper()
	var records []workspace.Record
	if err := json.Unmarshal([]byte(f.ok(append([]string{"record", "list"}, args...)...)), &records); err != nil {
		t.Fatal(err)
	}
	var ids []string
	for _, r := range records {
		ids = append(ids, r.ID)
	}
	return ids
}

func shownAt(t *testing.T, f fixture, id, want string) {
	t.Helper()
	if got := record(t, f.ok("record", "show", "--id", id)).Path; got != want {
		t.Fatalf("%s shown at %s, want %s", id, got, want)
	}
}

// checkOnlyGit runs check in a fixture that is a Git checkout, where the one
// finding expected is that it describes no workspace repository.
func checkOnlyGit(t *testing.T, f fixture) {
	t.Helper()
	_, out, _ := call(f.root, "check")
	var result struct {
		Issues []workspace.Finding `json:"issues"`
	}
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatal(err)
	}
	for _, issue := range result.Issues {
		if !strings.Contains(issue.Issue, "describes no workspace repository") {
			t.Fatalf("check reported %q", issue.Issue)
		}
	}
}

// Records moved into one level of group folder, and into _archive/ above or
// below that group, keep working with every command that reads them, and
// check finds nothing to report.
func TestRecordsInGroupFoldersWork(t *testing.T) {
	f := setup(t)
	f.repository("api")
	git(t, f.root, "init", "-q", "-b", "main")
	i := f.intent("thing")
	p := f.plan(i.ID, "thing", "api")
	q := f.plan(i.ID, "next", "api")
	f.ok("record", "dependencies", "--id", q.ID, "--depends-on", p.ID)
	i = move(t, f, i, "intent/alex/"+filepath.Base(i.Path))
	p = move(t, f, p, "plans/alex/"+p.ID+"-thing/plan.md")
	commitRecords(t, f.root)

	shownAt(t, f, i.ID, i.Path)
	shownAt(t, f, p.ID, p.Path)
	if ids := listed(t, f); !slices.Equal(ids, []string{i.ID, p.ID, q.ID}) {
		t.Fatalf("record list = %v", ids)
	}
	if o := order(t, f.ok("record", "order", "--plan", q.ID, "--plan", p.ID)); !slices.Equal(o.PlanIDs, []string{p.ID, q.ID}) {
		t.Fatalf("record order = %+v", o)
	}
	files, _ := f.gate("record", "approve", "--id", i.ID, "--text", "Go ahead.")
	expectFiles(t, files, map[string]string{i.Path: "modified"})
	commitRecords(t, f.root)
	files, _ = f.gate("record", "complete", "--id", p.ID, "--text", "Done.")
	expectFiles(t, files, map[string]string{p.Path: "modified"})
	checkOnlyGit(t, f)

	// Archived at both levels: below the group and above it.
	commitRecords(t, f.root)
	p = move(t, f, p, "plans/alex/_archive/"+p.ID+"-thing/plan.md")
	i = move(t, f, i, "intent/_archive/alex/"+filepath.Base(i.Path))
	if ids := listed(t, f); !slices.Equal(ids, []string{q.ID}) {
		t.Fatalf("record list showed archived records: %v", ids)
	}
	if ids := listed(t, f, "--archived"); !slices.Equal(ids, []string{i.ID, p.ID, q.ID}) {
		t.Fatalf("record list --archived = %v", ids)
	}
	shownAt(t, f, i.ID, i.Path)
	shownAt(t, f, p.ID, p.Path)
	if o := order(t, f.ok("record", "order", "--plan", q.ID)); !slices.Equal(o.PlanIDs, []string{q.ID}) {
		t.Fatalf("record order = %+v", o)
	}
	files, _ = f.gate("record", "complete", "--id", p.ID, "--text", "Done again.")
	expectFiles(t, files, map[string]string{
		p.Path:                                  "new",
		"plans/alex/" + p.ID + "-thing/plan.md": "deleted",
		i.Path:                                  "new",
		"intent/alex/" + filepath.Base(i.Path):  "deleted",
	})
	checkOnlyGit(t, f)

	preview := deletion(t, f.ok("record", "delete", "--id", q.ID, "--id", p.ID, "--id", i.ID))
	if r := deletedRecord(t, preview, p.ID); r.Path != strings.TrimSuffix(p.Path, "plan.md") {
		t.Fatalf("deletion preview shows %s at %s", p.ID, r.Path)
	}
	deletion(t, f.ok("record", "delete", "--id", q.ID, "--id", p.ID, "--id", i.ID, "--confirm"))
	for _, gone := range []string{i.Path, filepath.Dir(p.Path)} {
		if exists(filepath.Join(f.root, filepath.FromSlash(gone))) {
			t.Fatalf("delete left %s", gone)
		}
	}
}

// What the classifier would pass over silently is reported by check instead.
func TestCheckReportsWhatNoCommandReads(t *testing.T) {
	f := setup(t)
	i := f.intent("thing")
	record := "---\nid: p0500\ncreated_by: maya\ncreated_at: 2026-09-15T10:53:00Z\n---\n\n# Deep\n"
	write(t, filepath.Join(f.root, "plans", "_drafts", "p0500-deep", "plan.md"), record)
	write(t, filepath.Join(f.root, "plans", "alex", "deeper", "p0501-deep", "plan.md"), strings.Replace(record, "p0500", "p0501", 1))
	write(t, filepath.Join(f.root, "plans", "p0600-group", "p0502-deep", "plan.md"), strings.Replace(record, "p0500", "p0502", 1))
	out := f.fail("check")
	for _, want := range []string{
		"reserved folder name: plans/_drafts",
		"record file outside the folders records are read from: plans/alex/deeper/p0501-deep/plan.md",
		"folder named like a record holds records: plans/p0600-group",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("check did not report %q: %s", want, out)
		}
	}
	if strings.Contains(out, "missing or invalid plan entry: plans/p0600-group") {
		t.Errorf("check called a group folder a plan without plan.md: %s", out)
	}
	if ids := listed(t, f, "--archived"); !slices.Equal(ids, []string{i.ID}) {
		t.Fatalf("records outside the read folders were listed: %v", ids)
	}
}

// legacyArchive puts an approved intent and a completed plan under the legacy
// archive/ folders an earlier version used.
func legacyArchive(t *testing.T, f fixture) (workspace.Record, workspace.Record) {
	t.Helper()
	f.repository("api")
	i := f.intent("old")
	f.ok("record", "approve", "--id", i.ID, "--text", "Go ahead.")
	p := f.plan(i.ID, "old", "api")
	f.ok("record", "complete", "--id", p.ID, "--text", "Done.")
	i = move(t, f, i, "intent/archive/"+filepath.Base(i.Path))
	p = move(t, f, p, "plans/archive/"+p.ID+"-old/plan.md")
	return i, p
}

func TestTheFirstChangeRenamesTheLegacyArchive(t *testing.T) {
	f := setup(t)
	git(t, f.root, "init", "-q", "-b", "main")
	current := f.intent("current")
	i, p := legacyArchive(t, f)
	// An archive/ inside a group is an ordinary folder, and stays where it is.
	grouped := filepath.Join(f.root, "plans", "alex", "archive", "notes.md")
	write(t, grouped, "# Notes\n")
	commitRecords(t, f.root)

	// Before the rename, reads already treat the legacy folders as archived,
	// and check names the rename that is pending.
	if ids := listed(t, f); !slices.Equal(ids, []string{current.ID}) {
		t.Fatalf("record list before the rename = %v", ids)
	}
	if ids := listed(t, f, "--archived"); !slices.Equal(ids, []string{i.ID, current.ID, p.ID}) {
		t.Fatalf("record list --archived before the rename = %v", ids)
	}
	f.ok("status")
	out := f.fail("check")
	if !strings.Contains(out, "legacy archive folder awaits its rename to _archive/: intent/archive") || !strings.Contains(out, "plans/archive") {
		t.Fatalf("check did not report the pending rename: %s", out)
	}
	if !exists(filepath.Join(f.root, "intent", "archive")) || exists(filepath.Join(f.root, "intent", "_archive")) {
		t.Fatal("a read renamed the archive")
	}

	code, result, errOut := call(f.root, "record", "approve", "--id", current.ID, "--text", "Go ahead.")
	if code != 0 {
		t.Fatalf("exit=%d %s %s", code, result, errOut)
	}
	if !strings.Contains(errOut, "notice: renamed intent/archive to intent/_archive; renamed plans/archive to plans/_archive") || !strings.Contains(errOut, "person's decision") {
		t.Fatalf("the rename was not reported: %s", errOut)
	}
	var gate gateResult
	if err := json.Unmarshal([]byte(result), &gate); err != nil {
		t.Fatal(err)
	}
	files := map[string]string{}
	for _, file := range gate.CommitRecommended.Files {
		files[file.Path] = file.Change
	}
	expectFiles(t, files, map[string]string{
		current.Path: "modified",
		i.Path:       "deleted",
		p.Path:       "deleted",
		"intent/_archive/" + filepath.Base(i.Path): "new",
		"plans/_archive/" + p.ID + "-old/plan.md":  "new",
	})
	for _, gone := range []string{"intent/archive", "plans/archive"} {
		if exists(filepath.Join(f.root, filepath.FromSlash(gone))) {
			t.Fatalf("%s is still there", gone)
		}
	}
	if !exists(grouped) {
		t.Fatal("an archive/ inside a group was moved")
	}
	shownAt(t, f, i.ID, "intent/_archive/"+filepath.Base(i.Path))
	shownAt(t, f, p.ID, "plans/_archive/"+p.ID+"-old/plan.md")
	if ids := listed(t, f); !slices.Equal(ids, []string{current.ID}) {
		t.Fatalf("record list after the rename = %v", ids)
	}
	l := ledger(t, f)
	if !slices.Contains(l.Intents, i.ID) || !slices.Contains(l.Plans, p.ID) {
		t.Fatalf("the rename released a number: %+v", l)
	}
	if next := f.intent("next"); next.ID == i.ID {
		t.Fatal("an archived number was handed out again")
	}
	checkOnlyGit(t, f)

	// A clone that had not pulled brings archive/ back; the next change moves
	// it again, into the _archive/ that is there now.
	back := f.intent("back")
	back = move(t, f, back, "intent/archive/"+filepath.Base(back.Path))
	if out := f.fail("check"); !strings.Contains(out, "intent/archive") {
		t.Fatalf("check did not report the returned archive: %s", out)
	}
	_, _, errOut = call(f.root, "record", "create", "--kind", "intent", "--slug", "later", "--title", "Later")
	if !strings.Contains(errOut, "moved the contents of intent/archive into intent/_archive") {
		t.Fatalf("the returned archive was not merged: %s", errOut)
	}
	shownAt(t, f, back.ID, "intent/_archive/"+filepath.Base(back.Path))
	if exists(filepath.Join(f.root, "intent", "archive")) {
		t.Fatal("the returned archive is still there")
	}
}

// Both folders present are merged, a group folder both hold included; a name
// both hold with different content stops the rename before anything moves.
func TestTheRenameMergesAndRefusesAClash(t *testing.T) {
	f := setup(t)
	a, b, c := f.intent("a"), f.intent("b"), f.intent("c")
	a = move(t, f, a, "intent/archive/"+filepath.Base(a.Path))
	b = move(t, f, b, "intent/archive/alex/"+filepath.Base(b.Path))
	c = move(t, f, c, "intent/_archive/alex/"+filepath.Base(c.Path))
	clash := filepath.Join(f.root, "intent", "_archive", "README.md")
	write(t, clash, "kept\n")
	write(t, filepath.Join(f.root, "intent", "archive", "README.md"), "different\n")

	out := f.fail("record", "create", "--kind", "intent", "--slug", "d", "--title", "D")
	if !strings.Contains(out, "intent/archive/README.md") || !strings.Contains(out, "nothing was moved") {
		t.Fatalf("the clash was not named: %s", out)
	}
	if !exists(filepath.Join(f.root, filepath.FromSlash(a.Path))) || !exists(filepath.Join(f.root, filepath.FromSlash(b.Path))) {
		t.Fatal("a refused rename moved something")
	}
	if ids := listed(t, f, "--archived"); len(ids) != 3 {
		t.Fatalf("a refused rename created a record: %v", ids)
	}

	write(t, filepath.Join(f.root, "intent", "archive", "README.md"), "kept\n")
	f.ok("record", "create", "--kind", "intent", "--slug", "d", "--title", "D")
	shownAt(t, f, a.ID, "intent/_archive/"+filepath.Base(a.Path))
	shownAt(t, f, b.ID, "intent/_archive/alex/"+filepath.Base(b.Path))
	shownAt(t, f, c.ID, c.Path)
	if exists(filepath.Join(f.root, "intent", "archive")) || read(t, clash) != "kept\n" {
		t.Fatal("the merge left the legacy folder or changed what _archive/ held")
	}
}

// Only the pinned CLI makes the rename, as only it makes the ID split; reads
// from another CLI still run.
func TestAnotherCLIRefusesToRenameTheArchive(t *testing.T) {
	f := setup(t)
	i, _ := legacyArchive(t, f)
	pinned(t, f.root, "2.5.0")
	if code, out, errOut := runAs(f.root, "2.5.1", "record", "list", "--archived"); code != 0 || !strings.Contains(out, i.ID) {
		t.Fatalf("a read was refused: exit=%d %s %s", code, out, errOut)
	}
	code, out, errOut := runAs(f.root, "2.5.1", "record", "create", "--kind", "intent", "--slug", "a", "--title", "A")
	if code != 1 || out != "" || !strings.Contains(errOut, "this change would rename the workspace's archive/ folders to _archive/, which the pinned CLI must make; nothing was changed") {
		t.Fatalf("exit=%d %s %s", code, out, errOut)
	}
	if !exists(filepath.Join(f.root, "intent", "archive")) || exists(filepath.Join(f.root, "intent", "_archive")) {
		t.Fatal("a refused call renamed the archive")
	}
}
