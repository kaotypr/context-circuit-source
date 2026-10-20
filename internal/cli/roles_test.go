package cli_test

import (
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"
)

func TestAMemberRecordsTheRolesTheyContribute(t *testing.T) {
	f := setup(t)
	roster := filepath.Join(f.root, "members.yaml")

	f.ok("member", "language", "--id", "maya", "--language", "Bahasa Indonesia")
	f.ok("member", "tone", "--id", "maya", "--tone", "semi-formal")
	f.ok("member", "roles", "--id", "maya", "--role", "maintainer", "--role", "docs")
	if body := read(t, roster); !strings.Contains(body, "- maintainer") || !strings.Contains(body, "- docs") {
		t.Fatalf("expected both roles in the roster:\n%s", body)
	}
	if listed := strings.Join(strings.Fields(f.ok("member", "list")), ""); !strings.Contains(listed, `"roles":["maintainer","docs"]`) {
		t.Fatalf("expected member list to show the roles:\n%s", listed)
	}
	if got := orientation(t, f.ok("status")).Members.Members["maya"].Roles; !slices.Equal(got, []string{"maintainer", "docs"}) {
		t.Fatalf("expected status to show the roles: %v", got)
	}
	if human := f.human("member", "list"); !strings.Contains(human, "roles:") {
		t.Fatalf("expected the human roster to show the roles:\n%s", human)
	}

	// Setting roles again replaces them rather than adding to them.
	f.ok("member", "roles", "--id", "maya", "--role", "reviewer")
	if got := orientation(t, f.ok("status")).Members.Members["maya"].Roles; !slices.Equal(got, []string{"reviewer"}) {
		t.Fatalf("expected the roles to be replaced: %v", got)
	}

	f.ok("member", "roles", "--id", "maya", "--clear")
	body := read(t, roster)
	if strings.Contains(body, "roles:") {
		t.Fatalf("expected --clear to remove the roles:\n%s", body)
	}
	for _, kept := range []string{"band: 0", "language: Bahasa Indonesia", "tone: semi-formal"} {
		if !strings.Contains(body, kept) {
			t.Fatalf("clearing the roles must leave %q alone:\n%s", kept, body)
		}
	}

	// A role is free text: nothing is validated, an empty one included.
	f.ok("member", "add", "--id", "alex", "--name", "Alex", "--band", "2", "--role", "Release & QA: weekends", "--role", "")
	if got := orientation(t, f.ok("status")).Members.Members["alex"].Roles; !slices.Equal(got, []string{"Release & QA: weekends", ""}) {
		t.Fatalf("expected member add to carry the roles as given: %v", got)
	}
	// Adding the same member again is a no-op only when nothing would change.
	f.ok("member", "add", "--id", "alex", "--name", "Alex", "--band", "2")
	if out := f.fail("member", "add", "--id", "alex", "--name", "Alex", "--band", "2", "--role", "docs"); !strings.Contains(out, "use member roles") {
		t.Fatalf("re-adding with other roles should point at member roles: %s", out)
	}

	g := fixture{t, filepath.Join(f.home, "fresh"), f.home}
	g.ok("init", "--name", "Fresh", "--purpose", "Side project", "--member", "ana", "--member-name", "Ana", "--member-role", "maintainer")
	if got := orientation(t, g.ok("status")).Members.Members["ana"].Roles; !slices.Equal(got, []string{"maintainer"}) {
		t.Fatalf("expected init to record the first member's roles: %v", got)
	}
}

func TestRolesRefuseAnIncompleteCall(t *testing.T) {
	f := setup(t)
	for _, tc := range []struct {
		name string
		args []string
		want string
	}{
		{"neither value", []string{"member", "roles", "--id", "maya"}, "--clear"},
		{"both values", []string{"member", "roles", "--id", "maya", "--role", "docs", "--clear"}, "not both"},
		{"unknown member", []string{"member", "roles", "--id", "nobody", "--role", "docs"}, "unknown member"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if out := f.fail(tc.args...); !strings.Contains(out, tc.want) {
				t.Fatalf("expected a refusal naming %q:\n%s", tc.want, out)
			}
		})
	}
}

// A roster nobody gave roles to reads exactly as it did before roles existed.
func TestARosterWithoutRolesShowsNoRoles(t *testing.T) {
	f := setup(t)
	for _, out := range []string{read(t, filepath.Join(f.root, "members.yaml")), f.ok("member", "list"), f.ok("status"), f.human("member", "list"), f.human("status")} {
		if strings.Contains(out, "roles") {
			t.Fatalf("a roster without roles printed a roles key:\n%s", out)
		}
	}
}

// A schema 2 roster cannot carry roles, because a 2.3.0 CLI reading it would
// stop at the field instead of asking for an upgrade.
func TestRolesWaitForSchema3(t *testing.T) {
	f := setup(t)
	schema2(t, f, nil, nil)
	roster := filepath.Join(f.root, "members.yaml")
	upgraded(t, f, "members:\n  maya:\n    name: Maya\n    band: 0\n")
	before := read(t, roster)
	out := f.fail("member", "add", "--id", "alex", "--name", "Alex", "--band", "2", "--role", "docs")
	if !strings.Contains(out, "schema 3") || !strings.Contains(out, "member roles") {
		t.Fatalf("the refusal should name schema 3 and the next step: %s", out)
	}
	if read(t, roster) != before {
		t.Fatal("a refused member add changed members.yaml")
	}

	// Once setup is complete, the change that writes roles moves the workspace
	// to schema 3 first.
	f.ok("member", "use", "--id", "maya")
	f.ok("member", "roles", "--id", "maya", "--role", "maintainer")
	if !strings.Contains(read(t, filepath.Join(f.root, "workspace.yaml")), "version: 3") {
		t.Fatal("roles were written without moving the workspace to schema 3")
	}
	if !strings.Contains(read(t, roster), "- maintainer") {
		t.Fatalf("the roles were not written:\n%s", read(t, roster))
	}
	f.ok("check")
}

var instant = regexp.MustCompile(`\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z`)

// Roles decide nothing: the same work with and without them allocates the same
// IDs, approves and deletes the same records, and briefs agents the same way.
func TestRolesDecideNothing(t *testing.T) {
	run := func(roles bool) (ids []string, prompt, preview string) {
		f := setup(t)
		if roles {
			f.ok("member", "roles", "--id", "maya", "--role", "approver-of-everything")
		}
		i := f.intent("billing")
		next := f.intent("invoices")
		f.ok("record", "approve", "--id", i.ID, "--text", "Approved")
		spec := dispatch(t, f, "agent", "dispatch", "--host", "codex", "--role", "planner", "--intent", i.ID, "--path", f.root)
		// Each fixture lives in its own directory, which the brief names.
		home := f.home
		if real, err := filepath.EvalSymlinks(home); err == nil {
			home = real
		}
		prompt = strings.ReplaceAll(strings.ReplaceAll(spec.Prompt, home, "HOME"), f.home, "HOME")
		// The brief quotes the approval instant, which differs whenever the two
		// runs straddle a second, as they often do on a slow runner.
		prompt = instant.ReplaceAllString(prompt, "INSTANT")
		return []string{i.ID, next.ID}, prompt, strings.ReplaceAll(f.ok("record", "delete", "--member", "maya"), f.home, "HOME")
	}
	plainIDs, plainPrompt, plainPreview := run(false)
	roleIDs, rolePrompt, rolePreview := run(true)
	if !slices.Equal(plainIDs, roleIDs) {
		t.Fatalf("roles changed allocation: %v and %v", plainIDs, roleIDs)
	}
	if strings.Contains(rolePrompt, "approver-of-everything") {
		t.Fatal("a dispatch brief quoted a member's role")
	}
	if plainPrompt != rolePrompt {
		t.Fatal("roles changed a dispatch brief")
	}
	if plainPreview != rolePreview {
		t.Fatalf("roles changed what deletion would take:\n%s\n%s", plainPreview, rolePreview)
	}
}

// A hand-written roster may use flow style. Rewriting one member of it must
// still produce a roster that reads back.
func TestRolesCanBeSetInAFlowStyleRoster(t *testing.T) {
	f := setup(t)
	write(t, filepath.Join(f.root, "members.yaml"), "members: {maya: {name: Maya, band: 0}}\n")
	f.ok("member", "roles", "--id", "maya", "--role", "maintainer", "--role", "docs")
	if got := orientation(t, f.ok("status")).Members.Members["maya"].Roles; !slices.Equal(got, []string{"maintainer", "docs"}) {
		t.Fatalf("expected the roles to read back from a flow-style roster: %v", got)
	}
}
