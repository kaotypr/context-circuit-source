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

// authorless writes an intent the way an earlier version wrote one while the
// workspace had no roster: with no created_by.
func authorless(t *testing.T, f fixture, id string) string {
	t.Helper()
	rel := "intent/" + id + "-old.md"
	write(t, filepath.Join(f.root, filepath.FromSlash(rel)), "---\nid: "+id+"\ncreated_at: 2026-09-15T10:53:00Z\n---\n\n# Old\n")
	return rel
}

// finishedWork gives the active member an approved intent with one completed
// plan, which makes both of them finished.
func finishedWork(f fixture, slug string) (workspace.Record, workspace.Record) {
	i := f.intent(slug)
	f.ok("record", "approve", "--id", i.ID, "--text", "Go ahead.")
	p := f.plan(i.ID, slug, "api")
	f.ok("record", "complete", "--id", p.ID, "--text", "Done.")
	return record(f.t, f.ok("record", "show", "--id", i.ID)), record(f.t, f.ok("record", "show", "--id", p.ID))
}

func TestRecordListShowsYourOwnRecordsFirst(t *testing.T) {
	f := setup(t)
	twoMembers(f)
	mine := f.intent("mine")
	f.ok("member", "use", "--id", "alex")
	theirs := f.intent("theirs")
	f.ok("member", "use", "--id", "maya")
	old := authorless(t, f, "i050")

	if ids := listed(t, f); !slices.Equal(ids, []string{mine.ID}) {
		t.Fatalf("record list = %v, want only maya's", ids)
	}
	if ids := listed(t, f, "--member", "alex"); !slices.Equal(ids, []string{theirs.ID}) {
		t.Fatalf("record list --member alex = %v", ids)
	}
	if ids := listed(t, f, "--all"); !slices.Equal(ids, []string{"i050", mine.ID, theirs.ID}) {
		t.Fatalf("record list --all = %v", ids)
	}
	if out := f.fail("record", "list", "--all", "--member", "alex"); !strings.Contains(out, "not both") {
		t.Fatalf("--all with --member: %s", out)
	}

	// Reading is allowed before a member is selected: everyone's records are
	// listed, and the notice says why.
	if err := os.Remove(filepath.Join(f.root, "member.local.yaml")); err != nil {
		t.Fatal(err)
	}
	code, out, errOut := call(f.root, "record", "list")
	if code != 0 || !strings.Contains(errOut, "no member is selected") {
		t.Fatalf("record list without a member: exit=%d %s", code, errOut)
	}
	for _, id := range []string{mine.ID, theirs.ID, strings.TrimSuffix(filepath.Base(old), "-old.md")} {
		if !strings.Contains(out, `"`+id+`"`) {
			t.Fatalf("record list without a member left out %s: %s", id, out)
		}
	}
}

type myStatus struct {
	MyWork *workspace.MyWork   `json:"my_work"`
	Issues []workspace.Finding `json:"issues"`
}

func myWorkStatus(t *testing.T, f fixture) myStatus {
	t.Helper()
	out := f.ok("status")
	if strings.Contains(out, `"my_work"`) && !strings.HasPrefix(strings.TrimSpace(out), "{\n  \"my_work\"") {
		t.Fatalf("status does not open with my_work: %.80s", out)
	}
	var o myStatus
	if err := json.Unmarshal([]byte(out), &o); err != nil {
		t.Fatal(err)
	}
	return o
}

func ids(items []workspace.WorkItem) []string {
	out := []string{}
	for _, item := range items {
		out = append(out, item.ID)
	}
	return out
}

func TestStatusOpensWithYourOwnWork(t *testing.T) {
	f := setup(t)
	f.repository("api")
	twoMembers(f)
	waiting := f.intent("waiting")
	unplanned := f.intent("unplanned")
	f.ok("record", "approve", "--id", unplanned.ID, "--text", "Go ahead.")
	planned := f.intent("planned")
	f.ok("record", "approve", "--id", planned.ID, "--text", "Go ahead.")
	open := f.plan(planned.ID, "open", "api")
	done, donePlan := finishedWork(f, "done")
	// Alex's records are nobody else's business in maya's status.
	f.ok("member", "use", "--id", "alex")
	f.intent("alex-waiting")
	finishedWork(f, "alex-done")
	f.ok("member", "use", "--id", "maya")

	work := myWorkStatus(t, f).MyWork
	if work == nil || work.Member != "maya" {
		t.Fatalf("my_work = %+v", work)
	}
	if got := ids(work.AwaitingApproval); !slices.Equal(got, []string{waiting.ID}) {
		t.Errorf("awaiting_approval = %v", got)
	}
	if got := ids(work.ApprovedWithoutPlans); !slices.Equal(got, []string{unplanned.ID}) {
		t.Errorf("approved_without_plans = %v", got)
	}
	if got := ids(work.OpenPlans); !slices.Equal(got, []string{open.ID}) {
		t.Errorf("open_plans = %v", got)
	}
	if work.UnarchivedFinished != 2 {
		t.Errorf("unarchived_finished = %d, want 2 (%s and %s)", work.UnarchivedFinished, done.ID, donePlan.ID)
	}

	// A record that cannot be read becomes a finding; status still answers.
	write(t, filepath.Join(f.root, "intent", "i090-broken.md"), "no frontmatter\n")
	o := myWorkStatus(t, f)
	if o.MyWork != nil {
		t.Errorf("my_work reported from unreadable records: %+v", o.MyWork)
	}
	reported := false
	for _, issue := range o.Issues {
		reported = reported || strings.Contains(issue.Issue, "i090-broken.md")
	}
	if !reported {
		t.Errorf("an unreadable record was not reported: %+v", o.Issues)
	}
}

type archiveResult struct {
	Archived        []workspace.ArchivedRecord `json:"archived"`
	AlreadyArchived []string                   `json:"already_archived"`
	Unfinished      []string                   `json:"unfinished"`
	LinksRewritten  []workspace.LinkChange     `json:"links_rewritten"`
	LinksReported   []workspace.LinkChange     `json:"links_reported"`
}

func archived(t *testing.T, out string) archiveResult {
	t.Helper()
	var result archiveResult
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatal(err)
	}
	return result
}

func TestArchiveMovesYourFinishedRecords(t *testing.T) {
	f := setup(t)
	f.repository("api")
	git(t, f.root, "init", "-q", "-b", "main")
	twoMembers(f)
	i, p := finishedWork(f, "done")
	draft := f.intent("draft")
	f.ok("member", "use", "--id", "alex")
	ai, ap := finishedWork(f, "theirs")
	f.ok("member", "use", "--id", "maya")
	// A finished plan in a group folder is archived inside that group.
	i2, p2 := finishedWork(f, "grouped")
	p2 = move(t, f, p2, "plans/team/"+p2.ID+"-grouped/plan.md")
	i2 = move(t, f, i2, "intent/team/"+filepath.Base(i2.Path))
	commitRecords(t, f.root)
	ledger := read(t, filepath.Join(f.root, ".context-circuit", "ids", "maya.yaml"))

	files, _ := f.gate("record", "archive")
	want := map[string]string{
		i.ID:  "intent/_archive/" + filepath.Base(i.Path),
		p.ID:  "plans/_archive/" + p.ID + "-done/plan.md",
		i2.ID: "intent/team/_archive/" + filepath.Base(i2.Path),
		p2.ID: "plans/team/_archive/" + p2.ID + "-grouped/plan.md",
	}
	for id, path := range want {
		shownAt(t, f, id, path)
	}
	shownAt(t, f, draft.ID, draft.Path)
	shownAt(t, f, ai.ID, ai.Path)
	shownAt(t, f, ap.ID, ap.Path)
	if got := read(t, filepath.Join(f.root, ".context-circuit", "ids", "maya.yaml")); got != ledger {
		t.Fatalf("archiving changed the ID file:\n%s\nwant\n%s", got, ledger)
	}
	expectFiles(t, files, map[string]string{
		i.Path:      "deleted",
		want[i.ID]:  "new",
		p.Path:      "deleted",
		want[p.ID]:  "new",
		i2.Path:     "deleted",
		want[i2.ID]: "new",
		p2.Path:     "deleted",
		want[p2.ID]: "new",
	})
	if ids := listed(t, f); !slices.Equal(ids, []string{draft.ID}) {
		t.Fatalf("record list after archiving = %v", ids)
	}
	if work := myWorkStatus(t, f).MyWork; work.UnarchivedFinished != 0 {
		t.Fatalf("unarchived_finished after archiving = %d", work.UnarchivedFinished)
	}
	// Nothing finished is left, so a second run moves nothing.
	if again := archived(t, f.ok("record", "archive")); len(again.Archived) != 0 {
		t.Fatalf("a second archive moved %+v", again.Archived)
	}
	// The numbers stay reserved: the next intent does not reuse them.
	if next := f.intent("next"); next.ID == i.ID || next.ID == i2.ID {
		t.Fatalf("an archived number was handed out again: %s", next.ID)
	}

	// Another member's records, and unfinished ones, only by ID.
	commitRecords(t, f.root)
	byID := archived(t, f.ok("record", "archive", "--id", ap.ID, "--id", draft.ID, "--id", i.ID))
	if len(byID.Archived) != 2 || !slices.Equal(byID.Unfinished, []string{draft.ID}) || !slices.Equal(byID.AlreadyArchived, []string{i.ID}) {
		t.Fatalf("archive by ID = %+v", byID)
	}
	shownAt(t, f, ap.ID, "plans/_archive/"+ap.ID+"-theirs/plan.md")
	shownAt(t, f, draft.ID, "intent/_archive/"+filepath.Base(draft.Path))
	if out := f.fail("record", "archive", "--id", "i099"); !strings.Contains(out, "no record matches i099") {
		t.Fatalf("archiving an unknown ID: %s", out)
	}
	checkOnlyGit(t, f)
}

func TestArchiveKeepsLinksWorking(t *testing.T) {
	f := setup(t)
	f.repository("api")
	i, p := finishedWork(f, "done")
	other := f.intent("other")
	planDir := filepath.Join(f.root, filepath.FromSlash(filepath.Dir(p.Path)))
	write(t, filepath.Join(planDir, "design.md"), "# Design\n\nSee [the plan](plan.md) and [the intent](../../"+i.Path+").\n")
	otherFile := filepath.Join(f.root, filepath.FromSlash(other.Path))
	write(t, otherFile, read(t, otherFile)+
		"\nFollows [the plan](../"+p.Path+" \"plan\") and [design](../"+filepath.ToSlash(filepath.Dir(p.Path))+"/design.md#top).\n"+
		"Code `[not a link](../"+p.Path+")` and a [reference][r].\n\n[r]: ../"+p.Path+"\n")
	planFile := filepath.Join(planDir, "plan.md")
	write(t, planFile, read(t, planFile)+"\nBack to [the other intent](../../"+other.Path+").\n")

	result := archived(t, f.ok("record", "archive"))
	if len(result.Archived) != 2 {
		t.Fatalf("archive = %+v", result)
	}
	newPlanDir := filepath.Join(f.root, "plans", "_archive", filepath.Base(planDir))
	design := read(t, filepath.Join(newPlanDir, "design.md"))
	if !strings.Contains(design, "[the plan](plan.md)") || !strings.Contains(design, "[the intent](../../../intent/_archive/"+filepath.Base(i.Path)+")") {
		t.Errorf("design.md links: %s", design)
	}
	plan := read(t, filepath.Join(newPlanDir, "plan.md"))
	if !strings.Contains(plan, "[the other intent](../../../"+other.Path+")") {
		t.Errorf("plan.md links: %s", plan)
	}
	text := read(t, otherFile)
	archivedPlan := "../plans/_archive/" + filepath.Base(planDir)
	for _, want := range []string{
		"[the plan](" + archivedPlan + "/plan.md \"plan\")",
		"[design](" + archivedPlan + "/design.md#top)",
		"`[not a link](../" + p.Path + ")`",
		"[r]: ../" + p.Path,
	} {
		if !strings.Contains(text, want) {
			t.Errorf("%s does not hold %q:\n%s", other.Path, want, text)
		}
	}
	if len(result.LinksReported) != 1 || result.LinksReported[0].To != archivedPlan+"/plan.md" {
		t.Errorf("links_reported = %+v", result.LinksReported)
	}
	if len(result.LinksRewritten) != 4 {
		t.Errorf("links_rewritten = %+v", result.LinksRewritten)
	}
}

func TestCompletionOffersToArchive(t *testing.T) {
	f := setup(t)
	f.repository("api")
	i := f.intent("thing")
	f.ok("record", "approve", "--id", i.ID, "--text", "Go ahead.")
	p := f.plan(i.ID, "first", "api")
	q := f.plan(i.ID, "second", "api")
	var result struct {
		ArchiveOffer *struct {
			UnarchivedFinished int    `json:"unarchived_finished"`
			Offer              string `json:"offer"`
		} `json:"archive_offer"`
	}
	if err := json.Unmarshal([]byte(f.ok("record", "complete", "--id", p.ID, "--text", "Done.")), &result); err != nil {
		t.Fatal(err)
	}
	if result.ArchiveOffer == nil || result.ArchiveOffer.UnarchivedFinished != 1 || !strings.Contains(result.ArchiveOffer.Offer, "record archive") {
		t.Fatalf("archive_offer after the first plan = %+v", result.ArchiveOffer)
	}
	result.ArchiveOffer = nil
	if err := json.Unmarshal([]byte(f.ok("record", "complete", "--id", q.ID, "--text", "Done.")), &result); err != nil {
		t.Fatal(err)
	}
	// Both plans, and now the intent, are finished.
	if result.ArchiveOffer == nil || result.ArchiveOffer.UnarchivedFinished != 3 {
		t.Fatalf("archive_offer after the last plan = %+v", result.ArchiveOffer)
	}
	// Declining is doing nothing: every record stays where it was.
	for _, r := range []workspace.Record{i, p, q} {
		shownAt(t, f, r.ID, r.Path)
	}
}
