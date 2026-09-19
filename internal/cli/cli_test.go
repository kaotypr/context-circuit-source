package cli_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	assets "github.com/kaotypr/context-circuit-source"
	"github.com/kaotypr/context-circuit-source/internal/cli"
	"github.com/kaotypr/context-circuit-source/internal/workspace"
)

type fixture struct {
	t    *testing.T
	root string
	home string
}

func setup(t *testing.T) fixture {
	t.Helper()
	home := t.TempDir()
	config := filepath.Join(home, "gitconfig")
	// A fixture identity, and useConfigOnly so Git refuses to invent one from the
	// host instead. Without that guard a commit-creating call that forgot its
	// identity passes on a developer machine Git can guess on and fails only on
	// CI, which is where this was first noticed.
	write(t, config, "[user]\n\tname = Fixture\n\temail = fixture@example.invalid\n\tuseConfigOnly = true\n")
	t.Setenv("GIT_CONFIG_GLOBAL", config)
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	f := fixture{t, filepath.Join(home, "workspace"), home}
	f.ok("init", "--name", "Acme", "--purpose", "Billing software", "--member", "maya", "--member-name", "Maya")
	return f
}
func write(t *testing.T, path, text string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(text), 0644); err != nil {
		t.Fatal(err)
	}
}
func read(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}
func call(root string, args ...string) (int, string, string) {
	var out, err bytes.Buffer
	code := cli.Run(context.Background(), append([]string{"--workspace", root, "--json"}, args...), &out, &err, "2.0.0-test")
	return code, out.String(), err.String()
}
func (f fixture) ok(args ...string) string {
	f.t.Helper()
	code, out, err := call(f.root, args...)
	if code != 0 {
		f.t.Fatalf("%v: exit=%d %s %s", args, code, out, err)
	}
	return out
}

// human runs a command on the default stream, which is what a person sees. The
// rest of these tests ask for --json, so nothing else exercises it.
func (f fixture) human(args ...string) string {
	f.t.Helper()
	var out, errOut bytes.Buffer
	if code := cli.Run(context.Background(), append([]string{"--workspace", f.root}, args...), &out, &errOut, "2.0.0-test"); code != 0 {
		f.t.Fatalf("%v: exit=%d %s %s", args, code, out.String(), errOut.String())
	}
	return out.String()
}

func (f fixture) fail(args ...string) string {
	f.t.Helper()
	code, out, err := call(f.root, args...)
	if code == 0 {
		f.t.Fatalf("expected failure: %v: %s", args, out)
	}
	return out + err
}
func git(t *testing.T, path string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", path, "-c", "commit.gpgsign=false"}, args...)...)
	data, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v %s", args, err, data)
	}
	return strings.TrimSpace(string(data))
}
func (f fixture) repository(id string) string {
	f.t.Helper()
	path := filepath.Join(f.home, id)
	if err := os.Mkdir(path, 0755); err != nil {
		f.t.Fatal(err)
	}
	git(f.t, path, "init", "-b", "main")
	write(f.t, filepath.Join(path, "README.md"), "# Fixture\n")
	git(f.t, path, "add", "README.md")
	git(f.t, path, "-c", "user.name=Fixture", "-c", "user.email=fixture@example.invalid", "commit", "-m", "test: seed repository")
	f.ok("repo", "connect", "--id", id, "--path", path, "--base", "main")
	return path
}
func record(t *testing.T, out string) workspace.Record {
	t.Helper()
	var r workspace.Record
	if err := json.Unmarshal([]byte(out), &r); err != nil {
		t.Fatal(err)
	}
	return r
}
func tree(t *testing.T, out string) workspace.Worktree {
	t.Helper()
	var r workspace.Worktree
	if err := json.Unmarshal([]byte(out), &r); err != nil {
		t.Fatal(err)
	}
	return r
}
func (f fixture) intent(slug string) workspace.Record {
	return record(f.t, f.ok("record", "create", "--kind", "intent", "--slug", slug, "--title", slug))
}
func (f fixture) plan(id, slug string, repos ...string) workspace.Record {
	args := []string{"record", "create", "--kind", "plan", "--intent", id, "--slug", slug, "--title", slug}
	for _, repo := range repos {
		args = append(args, "--repo", repo)
	}
	return record(f.t, f.ok(args...))
}

func TestInitializationAndReleaseSeed(t *testing.T) {
	f := setup(t)
	files, err := assets.Files()
	if err != nil {
		t.Fatal(err)
	}
	for name := range files {
		if _, err := os.Stat(filepath.Join(f.root, filepath.FromSlash(name))); err != nil {
			t.Fatal(err)
		}
	}
	before := read(t, filepath.Join(f.root, "workspace.yaml"))
	f.fail("init", "--name", "Other", "--purpose", "Other", "--member", "other", "--member-name", "Other")
	if after := read(t, filepath.Join(f.root, "workspace.yaml")); after != before {
		t.Fatal("reinitialization changed data")
	}
	blank := filepath.Join(f.home, "blank")
	f.ok("template", "export", "--path", blank)
	if _, err := os.Stat(filepath.Join(blank, ".context-circuit", "local")); !os.IsNotExist(err) {
		t.Fatal("local runtime state leaked into exported seed")
	}
	// A Windows Git checkout may convert the published seed to CRLF.
	for name := range files {
		path := filepath.Join(blank, filepath.FromSlash(name))
		write(t, path, strings.ReplaceAll(read(t, path), "\n", "\r\n"))
	}
	g := fixture{t, blank, f.home}
	g.ok("init", "--name", "Seed", "--purpose", "Seed", "--member", "alex", "--member-name", "Alex")
	protected := filepath.Join(f.home, "protected")
	write(t, filepath.Join(protected, "AGENTS.md"), "Keep my instructions")
	h := fixture{t, protected, f.home}
	h.fail("init", "--name", "Bad", "--purpose", "Bad", "--member", "bad", "--member-name", "Bad")
	if read(t, filepath.Join(protected, "AGENTS.md")) != "Keep my instructions" {
		t.Fatal("existing file overwritten")
	}
	if _, err := os.Stat(filepath.Join(protected, "members.yaml")); !os.IsNotExist(err) {
		t.Fatal("failed preflight partially initialized")
	}
	// A README is not template material: a project adopting a workspace may
	// already have a front page, and initializing around it must leave it whole.
	adopted := filepath.Join(f.home, "adopted")
	write(t, filepath.Join(adopted, "README.md"), "Keep my README")
	a := fixture{t, adopted, f.home}
	a.ok("init", "--name", "Adopted", "--purpose", "Adopted", "--member", "alex", "--member-name", "Alex")
	if read(t, filepath.Join(adopted, "README.md")) != "Keep my README" {
		t.Fatal("initialization replaced a front page it did not write")
	}
	f.ok("check")
}

func TestYAMLEditTrial(t *testing.T) {
	f := setup(t)
	write(t, filepath.Join(f.root, "members.yaml"), "# Team roster\nmembers:\n  # Keep this author\n  maya:\n    name: \"Maya\" # display name\n")
	f.ok("member", "add", "--id", "alex", "--name", "Alex")
	updated := read(t, filepath.Join(f.root, "members.yaml"))
	for _, text := range []string{"# Team roster", "# Keep this author", "# display name", `"Maya"`} {
		if !strings.Contains(updated, text) {
			t.Fatalf("lost %s:\n%s", text, updated)
		}
	}
	f.ok("member", "use", "--id", "alex")
	r := f.intent("first")
	if r.CreatedBy != "alex" || r.ID != "i001" {
		t.Fatalf("unexpected record: %+v", r)
	}
	f.repository("api")
	before := "version: 2\nname: Acme\npurpose: |\n  Billing software\n  and invoicing\nrepositories:\n  api:\n    default_branch: main # preferred target\nrelationships: []\n"
	write(t, filepath.Join(f.root, "workspace.yaml"), before)
	f.ok("repo", "remote", "--id", "api", "--default-branch", "development")
	updated = read(t, filepath.Join(f.root, "workspace.yaml"))
	if !strings.Contains(updated, "and invoicing") || !strings.Contains(updated, "# preferred target") {
		t.Fatalf("lost description/comment: %s", updated)
	}
	write(t, filepath.Join(f.root, "members.yaml"), "members:\n  maya:\n    name: Maya\n  maya:\n    name: Duplicate\n")
	f.fail("member", "list")
	write(t, filepath.Join(f.root, "members.yaml"), "members:\n  maya:\n    name: Maya\n    owner: true\n")
	f.fail("member", "list")
}

// 2.0.0-rc.4 renamed the record gates and converts nothing. The decoder can only
// report an unknown key, which reads as a corrupt file, so the boundary says
// which version wrote the record and what replaced the field.
func TestRecordFromAnEarlierVersionNamesTheBoundary(t *testing.T) {
	f := setup(t)
	f.repository("api")
	i := f.intent("legacy")
	p := f.plan(i.ID, "api", "api")
	legacy := filepath.Join(f.root, p.Path)
	write(t, legacy, strings.Replace(read(t, legacy), "created_at:", "completed: 2026-09-01\ncreated_at:", 1))
	out := f.fail("record", "show", "--id", p.ID)
	for _, want := range []string{"2.0.0-rc.3 or earlier", "completed_at", "start a fresh workspace"} {
		if !strings.Contains(out, want) {
			t.Fatalf("the refusal must name %q: %s", want, out)
		}
	}
}

func TestGlobalRecordsDependenciesAndCompletion(t *testing.T) {
	f := setup(t)
	f.repository("api")
	f.repository("web")
	i := f.intent("billing")
	// A new intent carries the section that holds questions for the person, with
	// its empty state unnumbered so no placeholder claims a pending decision.
	if !strings.Contains(i.Content, "## Open questions\n\nNone known.") {
		t.Fatal("intent is missing the open questions section", i.Content)
	}
	if i.ApprovedAt != "" {
		t.Fatal("a new intent is not approved", i.ApprovedAt)
	}
	// Creation is an event the workspace witnessed, so every record carries it.
	if _, err := time.Parse(workspace.TimeLayout, string(i.CreatedAt)); err != nil {
		t.Fatal("a record records when it was created", i.CreatedAt, err)
	}
	f.ok("record", "approve", "--id", i.ID, "--text", "User approved $1 billing")
	approved := record(t, f.ok("record", "show", "--id", i.ID))
	if _, err := time.Parse(workspace.TimeLayout, string(approved.ApprovedAt)); err != nil {
		t.Fatal("approval records when it happened", approved.ApprovedAt, err)
	}
	if !strings.Contains(approved.Content, "## Approval — "+string(approved.ApprovedAt)) {
		t.Fatal("the frontmatter instant and the section must agree", approved.Content)
	}
	// Renewed approval is a second decision, so both are kept.
	f.ok("record", "approve", "--id", i.ID, "--text", "User approved the wider scope")
	again := record(t, f.ok("record", "show", "--id", i.ID))
	if strings.Count(again.Content, "## Approval — ") != 2 || again.ApprovedAt == "" {
		t.Fatal("renewed approval must not overwrite the first", again.Content)
	}
	p := f.plan(i.ID, "api", "api")
	q := f.plan(i.ID, "web", "web", "api")
	if p.ID != "p0001" || q.ID != "p0002" {
		t.Fatal(p.ID, q.ID)
	}
	if p.ApprovedAt != "" {
		t.Fatal("nothing approves a plan", p.ApprovedAt)
	}
	f.fail("record", "approve", "--id", p.ID, "--text", "Nothing approves a plan")
	parent := record(t, f.ok("record", "show", "--id", i.ID))
	if len(parent.Plans) != 2 || !strings.Contains(parent.Content, "User approved $1 billing") {
		t.Fatal("plan linking changed intent content", parent)
	}
	f.ok("record", "dependencies", "--id", q.ID, "--depends-on", p.ID)
	f.fail("record", "dependencies", "--id", p.ID, "--depends-on", q.ID)
	// A plan record holds the plan. Progress has no append command: the result
	// reaches the person, and completion — which a person requests — carries it.
	if out := f.fail("record", "note", "--id", p.ID, "--text", "Normal checks passed."); !strings.Contains(out, "unknown command") {
		t.Fatalf("unexpected error: %s", out)
	}
	f.ok("record", "complete", "--id", p.ID, "--text", "User requested completion.")
	f.ok("check")
	archive := filepath.Join(f.root, "plans", "archive")
	if err := os.Mkdir(archive, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(filepath.Join(f.root, p.Path), filepath.Join(archive, filepath.Base(p.Path))); err != nil {
		t.Fatal(err)
	}
	r := f.plan(i.ID, "more", "api")
	if r.ID != "p0003" {
		t.Fatal(r.ID)
	}
	f.ok("record", "show", "--id", p.ID)
	// A hand-edited instant is reported rather than read as a gate that never
	// passed, and neither gate may appear on the other kind of record.
	intentPath := filepath.Join(f.root, i.Path)
	sound := read(t, intentPath)
	write(t, intentPath, strings.Replace(sound, "approved_at: ", "approved_at: last ", 1))
	f.fail("check")
	write(t, intentPath, strings.Replace(sound, "approved_at: ", "completed_at: ", 1))
	f.fail("check")
	write(t, intentPath, strings.Replace(sound, "created_at: ", "created_at: today ", 1))
	f.fail("check")
	write(t, intentPath, sound)
	f.ok("check")
	var active []workspace.Record
	if err := json.Unmarshal([]byte(f.ok("record", "list")), &active); err != nil {
		t.Fatal(err)
	}
	for _, item := range active {
		if item.ID == p.ID {
			t.Fatal("active listing included archived record")
		}
	}
	if !strings.Contains(f.ok("record", "list", "--archived"), "p0001") {
		t.Fatal("archive not listed")
	}
	if err := os.Remove(filepath.Join(f.root, r.Path)); err != nil {
		t.Fatal(err)
	}
	x := f.plan(i.ID, "after-delete", "api")
	if x.ID != "p0004" {
		t.Fatal("reused deleted ID", x.ID)
	}
	f.fail("check") // Missing plan still linked from intent is visible.
}

func TestRepositoryBindingsAndCreation(t *testing.T) {
	f := setup(t)
	api := f.repository("api")
	f.fail("repo", "connect", "--id", "duplicate", "--path", api, "--base", "main")
	f.ok("repo", "clone", "--id", "copy", "--path", filepath.Join(f.home, "copy"), "--base", "main", "--url", api)
	f.ok("repo", "fetch", "--id", "copy")
	f.ok("repo", "init", "--id", "docs", "--path", filepath.Join(f.home, "docs"), "--base", "main")
	f.ok("repo", "relate", "--from", "copy", "--to", "api", "--description", "Consumes API")
	f.fail("repo", "relate", "--from", "copy", "--to", "missing", "--description", "Bad")
	git(t, f.root, "init", "-b", "main")
	// The workspace's own checkout is described by its own commands, because
	// every consumer of the repositories map reads an entry as somewhere work
	// happens.
	if out := f.fail("repo", "connect", "--id", "workspace", "--path", ".", "--base", "main"); !strings.Contains(out, "workspace connect") {
		t.Fatal("connecting the workspace as a repository should name workspace connect", out)
	}
	f.ok("workspace", "connect", "--base", "main")
	bindings := read(t, filepath.Join(f.root, "repositories.local.yaml"))
	if !strings.Contains(bindings, "path: .") {
		t.Fatal("workspace binding is not relative", bindings)
	}
	for _, path := range []string{"repositories.local.yaml", "member.local.yaml", ".context-circuit/local/write.lock", ".worktrees/test/file", "repositories/test/file"} {
		git(t, f.root, "check-ignore", path)
	}
	write(t, filepath.Join(api, "dirty.txt"), "preserve")
	f.ok("repo", "inspect", "--id", "api")
	if read(t, filepath.Join(api, "dirty.txt")) != "preserve" {
		t.Fatal("lost dirty work")
	}
}

// A record is quoted into a dispatch brief as authoritative, so a record in one
// language invites a worker to answer it in that language — in comments, in
// commit messages, and worst of all in identifiers. The worker never reads the
// workspace's own instructions, so the boundary travels in the brief.
func TestABriefDrawsTheLanguageBoundary(t *testing.T) {
	f := setup(t)
	api := f.repository("api")
	i := f.intent("tambah-penagihan")
	f.ok("record", "approve", "--id", i.ID, "--text", "Disetujui oleh pengguna")
	p := f.plan(i.ID, "penagihan-api", "api")

	// A member who records no language writes English, which is what every
	// record written before the field assumed.
	brief := f.ok("agent", "dispatch", "--host", "claude-code", "--role", "worker", "--plan", p.ID, "--path", api)
	if !strings.Contains(brief, "stated in English") {
		t.Fatal("an unset language is not reported as English", brief)
	}

	f.ok("member", "language", "--id", "maya", "--language", "Bahasa Indonesia")
	brief = f.ok("agent", "dispatch", "--host", "claude-code", "--role", "worker", "--plan", p.ID, "--path", api)
	for _, want := range []string{
		"stated in Bahasa Indonesia",
		"put into the repository in English",
		"Names are quoted, never translated",
		"Domain vocabulary keeps the form used here",
	} {
		if !strings.Contains(brief, want) {
			t.Fatalf("the brief does not draw the boundary (%q):\n%s", want, brief)
		}
	}

	// The language reported is the one the record was written in, not the one
	// whoever is dispatching happens to write in.
	f.ok("member", "add", "--id", "rina", "--name", "Rina", "--language", "English")
	f.ok("member", "use", "--id", "rina")
	brief = f.ok("agent", "dispatch", "--host", "claude-code", "--role", "worker", "--plan", p.ID, "--path", api)
	if !strings.Contains(brief, "stated in Bahasa Indonesia") {
		t.Fatal("the brief reports the dispatcher's language rather than the author's", brief)
	}

	// It is quoted into a sentence, so it is a language's name and not prose
	// carrying instructions of its own.
	f.fail("member", "language", "--id", "rina", "--language", strings.Repeat("long ", 20))
	f.fail("member", "language", "--id", "nobody", "--language", "English")
	f.fail("member", "add", "--id", "maya", "--name", "Maya", "--language", "English")
}

// The whole of joining a published workspace, driven against a real clone: a
// member arrives, describes this machine, obtains the repositories the record
// already names, and writes the role definitions a clone cannot carry. It ends
// on a clean diagnostic with one shared file touched, which is what the shipped
// instruction promises a joining member.
func TestJoiningAPublishedWorkspace(t *testing.T) {
	f := setup(t)
	api := f.repository("api")
	// A published workspace records where its repositories live; the fixture
	// checkout has no origin to read one from.
	f.ok("repo", "remote", "--id", "api", "--url", api)

	git(t, f.root, "init", "-b", "main")
	f.ok("workspace", "connect", "--base", "main")
	git(t, f.root, "add", "-A")
	git(t, f.root, "-c", "user.name=Fixture", "-c", "user.email=fixture@example.invalid", "commit", "-m", "test: publish the workspace")
	origin := filepath.Join(f.home, "workspace.git")
	git(t, f.home, "clone", "--bare", f.root, origin)

	joined := filepath.Join(f.home, "joined")
	git(t, f.home, "clone", origin, joined)
	ok := func(args ...string) {
		t.Helper()
		if code, out, errOut := call(joined, args...); code != 0 {
			t.Fatalf("%v: exit=%d %s %s", args, code, out, errOut)
		}
	}

	// A clone carries the shared records and none of this machine's state, and
	// the diagnostic is what says so.
	code, out, _ := call(joined, "check")
	for _, want := range []string{"select a workspace member", "connect this machine's checkout for repository api", "workspace connect"} {
		if !strings.Contains(out, want) {
			t.Fatalf("a fresh clone does not report %q:\n%s", want, out)
		}
	}
	if code == 0 {
		t.Fatal("a fresh clone reported nothing to do")
	}

	// Nothing here names a URL or reaches for git: the record carries both.
	ok("member", "add", "--id", "rina", "--name", "Rina", "--band", "2")
	ok("member", "band", "--id", "maya", "--band", "1")
	ok("member", "use", "--id", "rina")
	ok("workspace", "connect", "--base", "main")
	ok("repo", "clone", "--id", "api", "--path", "repositories/api", "--base", "main")
	ok("agent", "setup")

	if code, out, errOut := call(joined, "check"); code != 0 {
		t.Fatalf("joining left the workspace unhealthy: %s %s", out, errOut)
	}
	if read(t, filepath.Join(joined, "repositories/api/README.md")) != "# Fixture\n" {
		t.Fatal("the described repository was not obtained")
	}
	if _, err := os.Stat(filepath.Join(joined, ".claude/agents/cc-worker.md")); err != nil {
		t.Fatal("role definitions a clone cannot carry were not written", err)
	}
	// The roster is the only thing a member arriving has to say to everyone
	// else; a path or a branch of theirs reaching the shared record would be a
	// machine describing the project.
	dirty := git(t, joined, "status", "--porcelain")
	if dirty != "M members.yaml" {
		t.Fatalf("joining changed more than the roster:\n%s", dirty)
	}
}

// Joining a workspace means obtaining checkouts of repositories it already
// describes. The URL is in the shared record, so a clone need not repeat it,
// and nothing shared is rewritten by a member arriving.
func TestCloningARepositoryTheWorkspaceAlreadyDescribes(t *testing.T) {
	f := setup(t)
	api := f.repository("api")
	f.ok("repo", "remote", "--id", "api", "--url", api)
	before := read(t, filepath.Join(f.root, "workspace.yaml"))

	f.ok("repo", "clone", "--id", "api", "--path", "repositories/api", "--base", "main")
	if read(t, filepath.Join(f.root, "repositories/api/README.md")) != "# Fixture\n" {
		t.Fatal("the described repository was not obtained")
	}
	if after := read(t, filepath.Join(f.root, "workspace.yaml")); after != before {
		t.Fatalf("joining rewrote the shared record:\n%s\n%s", before, after)
	}

	// An ID nobody describes still needs a source, and starting one empty is
	// refused where a repository already exists to obtain.
	f.fail("repo", "clone", "--id", "fresh", "--path", filepath.Join(f.home, "fresh"), "--base", "main")
	if out := f.fail("repo", "init", "--id", "api", "--path", filepath.Join(f.home, "empty"), "--base", "main"); !strings.Contains(out, "already registered") {
		t.Fatal("initializing a described repository should be refused", out)
	}
}

// The shared records travel through the workspace's own repository, so its Git
// state is part of what orientation answers. It splits the way a repository
// does, and stays out of the repositories map, whose entries every consumer
// reads as somewhere work happens.
func TestWorkspaceRepositoryIsDescribedSeparately(t *testing.T) {
	f := setup(t)
	f.repository("api")

	// Nothing to describe until the workspace is under version control.
	if out := f.ok("check"); strings.Contains(out, "workspace connect") {
		t.Fatal("asked to describe a workspace that is not a Git checkout", out)
	}
	git(t, f.root, "init", "-b", "main")
	git(t, f.root, "remote", "add", "origin", "git@example.invalid:acme/workspace.git")
	if out := f.fail("check"); !strings.Contains(out, "workspace connect") {
		t.Fatal("an undescribed workspace repository goes unreported", out)
	}

	f.ok("workspace", "connect", "--base", "main")
	shared := read(t, filepath.Join(f.root, "workspace.yaml"))
	if !strings.Contains(shared, "workspace_repository:") || !strings.Contains(shared, "acme/workspace.git") {
		t.Fatalf("the shared record does not describe the workspace repository: %s", shared)
	}
	if strings.Contains(shared, "base_branch") {
		t.Fatalf("the shared record carries one machine's base branch: %s", shared)
	}
	if out := f.ok("check"); strings.Contains(out, "workspace connect") {
		t.Fatal("still asking for a workspace repository that is recorded", out)
	}

	// A base describes this machine alone, exactly as a repository's does.
	f.ok("workspace", "base", "--branch", "integration")
	bindings := read(t, filepath.Join(f.root, "repositories.local.yaml"))
	if !strings.Contains(bindings, "base_branch: integration") {
		t.Fatalf("the workspace base is not local: %s", bindings)
	}
	if strings.Contains(read(t, filepath.Join(f.root, "workspace.yaml")), "integration") {
		t.Fatal("a machine's workspace base reached the shared record")
	}

	f.ok("workspace", "remote", "--default-branch", "trunk")
	if !strings.Contains(read(t, filepath.Join(f.root, "workspace.yaml")), "default_branch: trunk") {
		t.Fatal("the shared default branch did not change")
	}
	f.fail("workspace", "remote")

	var state workspace.Orientation
	if err := json.Unmarshal([]byte(f.ok("status")), &state); err != nil {
		t.Fatal(err)
	}
	if state.Self == nil || state.Self.Branch != "main" || state.Self.BaseBranch != "integration" {
		t.Fatalf("status does not report the workspace's own Git state: %+v", state.Self)
	}
}

// Bands are what let clones that cannot see each other allocate without
// colliding. A solo workspace needs none, so the finding waits for a second
// member to make a collision possible.
func TestUnbandedMembersAreReportedOnlyOnceSharingIsPossible(t *testing.T) {
	f := setup(t)
	if out := f.ok("check"); strings.Contains(out, "shared range") {
		t.Fatal("a solo workspace was asked for allocation bands", out)
	}
	f.ok("member", "add", "--id", "rina", "--name", "Rina")
	out := f.fail("check")
	if !strings.Contains(out, "maya") || !strings.Contains(out, "rina") {
		t.Fatal("unbanded members went unreported", out)
	}
	f.ok("member", "band", "--id", "maya", "--band", "1")
	f.ok("member", "band", "--id", "rina", "--band", "2")
	if out := f.ok("check"); strings.Contains(out, "shared range") {
		t.Fatal("banded members are still reported", out)
	}
}

// A base branch describes one machine, not the repository: the checkout path and
// the branch work starts from travel together in the local binding, while the
// shared record keeps the URL and the default branch everyone agrees on.
func TestBaseBranchBelongsToTheLocalBinding(t *testing.T) {
	f := setup(t)
	api := f.repository("api")
	git(t, api, "checkout", "-b", "release")
	write(t, filepath.Join(api, "release.txt"), "release work")
	git(t, api, "add", "release.txt")
	git(t, api, "-c", "user.name=Fixture", "-c", "user.email=fixture@example.invalid", "commit", "-m", "test: advance the release branch")
	release := git(t, api, "rev-parse", "HEAD")
	git(t, api, "checkout", "main")
	main := git(t, api, "rev-parse", "HEAD")

	f.ok("repo", "base", "--id", "api", "--branch", "release")
	if bindings := read(t, filepath.Join(f.root, "repositories.local.yaml")); !strings.Contains(bindings, "base_branch: release") {
		t.Fatalf("base branch is not recorded locally: %s", bindings)
	}
	shared := read(t, filepath.Join(f.root, "workspace.yaml"))
	if strings.Contains(shared, "base_branch") || !strings.Contains(shared, "default_branch: main") {
		t.Fatalf("shared record carries one machine's base branch: %s", shared)
	}

	i := f.intent("release-work")
	p := f.plan(i.ID, "api", "api")
	w := tree(t, f.ok("worktree", "prepare", "--repo", "api", "--plan", p.ID))
	if w.Head != release {
		t.Fatalf("worktree did not start from the local base branch: %+v", w)
	}
	// An order is derived where it will be run, so it names the same base.
	if start := order(t, f.ok("record", "order", "--mode", "waves")).Layers[0].Plans[0].Start["api"]; start.Base != "release" {
		t.Fatalf("order started from %q rather than the local base branch", start.Base)
	}

	// A binding written before the base moved here still resolves, through the
	// shared default.
	write(t, filepath.Join(f.root, "repositories.local.yaml"), "bindings:\n  api:\n    path: "+api+"\n")
	q := f.plan(i.ID, "default-base", "api")
	w = tree(t, f.ok("worktree", "prepare", "--repo", "api", "--plan", q.ID))
	if w.Head != main {
		t.Fatalf("a binding with no base did not fall back to the default branch: %+v", w)
	}
	// The reported base is the resolved answer, so it names the fallback too;
	// reporting an empty field here sends the reader back to the shared file.
	if w.BaseBranch != "main" {
		t.Fatalf("the fallback base was not reported: %+v", w)
	}
}

// The shared record says where a repository lives. Connecting takes that from
// the checkout Git already has, and a URL carrying a password is refused rather
// than committed.
func TestSharedRepositoryRecordKeepsTheURL(t *testing.T) {
	f := setup(t)
	api := f.repository("api")
	// A recorded URL is compared as a value rather than as a line of the file:
	// a local path is a valid Git URL, and on Windows it carries backslashes
	// that YAML must quote and escape, so the serialized form is not the input.
	recorded := func(id string) string {
		t.Helper()
		cfg, err := (&workspace.Store{Root: f.root}).Config()
		if err != nil {
			t.Fatal(err)
		}
		return cfg.Repositories[id].URL
	}
	f.ok("repo", "clone", "--id", "copy", "--path", filepath.Join(f.home, "copy"), "--base", "main", "--url", api)
	if got := recorded("copy"); got != api {
		t.Fatalf("clone recorded %q rather than where the repository came from, %q", got, api)
	}
	f.ok("repo", "remote", "--id", "api", "--url", "https://example.invalid/api.git")
	if got := recorded("api"); got != "https://example.invalid/api.git" {
		t.Fatalf("repo remote recorded %q", got)
	}
	f.fail("repo", "remote", "--id", "api", "--url", "https://user:secret@example.invalid/api.git")
	f.fail("repo", "remote", "--id", "api")
	f.fail("repo", "remote", "--id", "missing", "--url", "https://example.invalid/missing.git")
	if shared := read(t, filepath.Join(f.root, "workspace.yaml")); strings.Contains(shared, "secret") {
		t.Fatalf("a password reached a shared file: %s", shared)
	}
	// Connecting another machine's checkout describes that machine only.
	f.fail("repo", "connect", "--id", "copy", "--path", filepath.Join(f.home, "copy"), "--base", "main", "--default-branch", "develop")
}

func TestWorktreePrepareResumeMoveRepairRemove(t *testing.T) {
	f := setup(t)
	api := f.repository("api")
	i := f.intent("billing")
	p := f.plan(i.ID, "api", "api")
	write(t, filepath.Join(api, "main-dirty.txt"), "keep main work")
	w := tree(t, f.ok("worktree", "prepare", "--repo", "api", "--plan", p.ID))
	if w.Head != git(t, api, "rev-parse", "HEAD") || w.StartCommit != w.Head || w.Branch != "cc/p0001/api" {
		t.Fatalf("wrong prepared state: %+v", w)
	}
	if _, err := os.Stat(filepath.Join(w.Path, "main-dirty.txt")); !os.IsNotExist(err) {
		t.Fatal("copied uncommitted main work")
	}
	write(t, filepath.Join(w.Path, "change.txt"), "work in progress")
	resumed := tree(t, f.ok("worktree", "prepare", "--repo", "api", "--plan", p.ID))
	if !resumed.Reused || resumed.Path != w.Path {
		t.Fatal("did not reuse prepared worktree")
	}
	f.ok("record", "complete", "--id", p.ID, "--text", "Explicitly marked done")
	if read(t, filepath.Join(w.Path, "change.txt")) != "work in progress" {
		t.Fatal("completion removed work")
	}
	f.fail("worktree", "remove", "--repo", "api", "--path", w.Path)
	f.fail("worktree", "prepare", "--repo", "api", "--plan", p.ID, "--path", filepath.Join(f.home, "collision"))
	moved := filepath.Join(f.home, "moved worktree")
	f.ok("worktree", "move", "--repo", "api", "--path", w.Path, "--to", moved)
	moved, _ = filepath.EvalSymlinks(moved)
	f.ok("worktree", "inspect", "--repo", "api", "--path", moved)
	manual := filepath.Join(f.home, "manual move")
	if err := os.Rename(moved, manual); err != nil {
		t.Fatal(err)
	}
	f.ok("worktree", "repair", "--repo", "api", "--path", manual)
	f.ok("check")
	f.ok("worktree", "remove", "--repo", "api", "--path", manual, "--discard")
	git(t, api, "show-ref", "--verify", "refs/heads/"+w.Branch)
	f.fail("worktree", "remove", "--repo", "api", "--path", api, "--discard")
	if read(t, filepath.Join(api, "main-dirty.txt")) != "keep main work" {
		t.Fatal("main work changed")
	}
}

func TestWorktreeExistingBranchAndIgnoredFiles(t *testing.T) {
	f := setup(t)
	api := f.repository("api")
	git(t, api, "branch", "existing")
	path := filepath.Join(f.home, "existing-worktree")
	f.fail("worktree", "prepare", "--repo", "api", "--branch", "existing", "--path", path)
	w := tree(t, f.ok("worktree", "prepare", "--repo", "api", "--branch", "existing", "--path", path, "--reuse"))
	write(t, filepath.Join(w.Path, ".gitignore"), "local.env\n")
	git(t, w.Path, "add", ".gitignore")
	git(t, w.Path, "-c", "user.name=Fixture", "-c", "user.email=fixture@example.invalid", "commit", "-m", "test: ignore fixture")
	write(t, filepath.Join(w.Path, "local.env"), "test-only")
	f.fail("worktree", "remove", "--repo", "api", "--path", path)
	git(t, api, "worktree", "lock", path)
	f.fail("worktree", "remove", "--repo", "api", "--path", path, "--discard")
	git(t, api, "worktree", "unlock", path)
	if err := os.Remove(filepath.Join(path, "local.env")); err != nil {
		t.Fatal(err)
	}
	f.ok("worktree", "remove", "--repo", "api", "--path", path)
}

func TestWorktreeBaseSelectionAndAdoption(t *testing.T) {
	f := setup(t)
	api := f.repository("api")
	git(t, api, "tag", "main")
	write(t, filepath.Join(api, "next.txt"), "next commit")
	git(t, api, "add", "next.txt")
	git(t, api, "-c", "user.name=Fixture", "-c", "user.email=fixture@example.invalid", "commit", "-m", "test: advance base")
	expected := git(t, api, "rev-parse", "HEAD")
	i := f.intent("choice")
	p := f.plan(i.ID, "api", "api")
	w := tree(t, f.ok("worktree", "prepare", "--repo", "api", "--plan", p.ID))
	if w.Head != expected {
		t.Fatal("selected same-named tag instead of base branch", w)
	}
	manual := filepath.Join(f.home, "native-created")
	git(t, api, "worktree", "add", "-b", "native-branch", manual, "HEAD")
	f.fail("worktree", "prepare", "--repo", "api", "--branch", "native-branch", "--path", manual)
	q := f.plan(i.ID, "adopt", "api")
	adopted := tree(t, f.ok("worktree", "prepare", "--repo", "api", "--plan", q.ID, "--branch", "native-branch", "--path", manual, "--reuse"))
	if adopted.Plan != q.ID || !adopted.Reused {
		t.Fatal("association was not adopted", adopted)
	}
	f.ok("worktree", "prepare", "--repo", "api", "--plan", q.ID, "--branch", "native-branch", "--path", manual)
	git(t, api, "branch", "keep-existing", expected)
	f.ok("repo", "base", "--id", "api", "--branch", "missing-base")
	f.ok("worktree", "prepare", "--repo", "api", "--branch", "keep-existing", "--path", filepath.Join(f.home, "existing"), "--reuse")
	missing := filepath.Join(f.home, "missing-start")
	f.fail("worktree", "prepare", "--repo", "api", "--branch", "invalid-start", "--path", missing, "--start", "absent-ref")
	if _, err := os.Stat(missing); !os.IsNotExist(err) {
		t.Fatal("created directory before resolving start")
	}
	if _, err := workspace.Git(context.Background(), api, "show-ref", "--verify", "refs/heads/invalid-start"); err == nil {
		t.Fatal("created branch before resolving start")
	}
	// Inherited checkout overrides must not redirect a selected repository.
	foreign := f.repository("foreign")
	t.Setenv("GIT_DIR", filepath.Join(foreign, ".git"))
	t.Setenv("GIT_WORK_TREE", foreign)
	// One-shot configuration reaches a child the same way: a host that ran
	// `git -c` exports it to everything underneath. Both spellings are set, so
	// filtering only one of them still fails this.
	t.Setenv("GIT_CONFIG_COUNT", "1")
	t.Setenv("GIT_CONFIG_KEY_0", "core.abbrev")
	t.Setenv("GIT_CONFIG_VALUE_0", "bogus")
	t.Setenv("GIT_CONFIG_PARAMETERS", "'core.abbrev=bogus'")
	var snapshot workspace.Snapshot
	if err := json.Unmarshal([]byte(f.ok("repo", "inspect", "--id", "api")), &snapshot); err != nil {
		t.Fatal(err)
	}
	if snapshot.Head != expected {
		t.Fatal("host environment redirected Git", snapshot)
	}
}

func TestHumanOutputAndDuplicateRecordProtection(t *testing.T) {
	f := setup(t)
	i := f.intent("readable")
	var out, errOut bytes.Buffer
	if code := cli.Run(context.Background(), []string{"--workspace", f.root, "record", "show", "--id", i.ID}, &out, &errOut, "test"); code != 0 {
		t.Fatal(code, errOut.String())
	}
	if out.String() != i.Content {
		t.Fatal("human record display omitted Markdown")
	}
	duplicate := filepath.Join(f.root, "intent", "archive", i.ID+"-duplicate.md")
	write(t, duplicate, i.Content)
	f.fail("record", "show", "--id", i.ID)
	before := read(t, filepath.Join(f.root, ".context-circuit", "ids.yaml"))
	f.fail("record", "create", "--kind", "intent", "--slug", "another", "--title", "Another")
	if read(t, filepath.Join(f.root, ".context-circuit", "ids.yaml")) != before {
		t.Fatal("allocated through an unresolved duplicate")
	}
	if err := os.Remove(duplicate); err != nil {
		t.Fatal(err)
	}
	write(t, filepath.Join(f.root, i.Path), "---\nid: "+i.ID+"\ncreated_by: maya\nowner: another\n---\nbody\n")
	f.fail("record", "show", "--id", i.ID)
}

func TestConcurrentAllocation(t *testing.T) {
	f := setup(t)
	var wg sync.WaitGroup
	errors := make(chan string, 12)
	ids := make(chan string, 12)
	for n := 0; n < 12; n++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			code, out, err := call(f.root, "record", "create", "--kind", "intent", "--slug", fmt.Sprintf("parallel-%d", n), "--title", "Parallel")
			if code != 0 {
				errors <- err
				return
			}
			var r workspace.Record
			if err := json.Unmarshal([]byte(out), &r); err != nil {
				errors <- err.Error()
				return
			}
			ids <- r.ID
		}(n)
	}
	wg.Wait()
	close(errors)
	close(ids)
	for err := range errors {
		t.Error(err)
	}
	seen := map[string]bool{}
	for id := range ids {
		if seen[id] {
			t.Fatal("duplicate ID", id)
		}
		seen[id] = true
	}
	if len(seen) != 12 {
		t.Fatalf("got %d IDs", len(seen))
	}
	f.ok("check")
}

// A mistyped argument is reported, never guessed at and never fatal. An empty
// one names no command and has no first word for the group hint to look up.
func TestMistypedArgumentsAreReported(t *testing.T) {
	f := setup(t)
	if out := f.fail(""); !strings.Contains(out, "unknown command") {
		t.Errorf("an empty command argument: %s", out)
	}
	if out := f.fail("", "trailing"); !strings.Contains(out, "unknown command") {
		t.Errorf("an empty command argument with more behind it: %s", out)
	}
}

// A band is cleared by passing zero, so an omitted --band would otherwise clear
// the member's block without anybody asking for that.
func TestClearingABandTakesTheFlagThatClearsIt(t *testing.T) {
	f := setup(t)
	f.ok("member", "add", "--id", "rina", "--name", "Rina", "--band", "2")
	if out := f.fail("member", "band", "--id", "rina"); !strings.Contains(out, "missing --band") {
		t.Errorf("omitting --band was accepted: %s", out)
	}
	if !strings.Contains(f.ok("member", "list"), `"band": 2`) {
		t.Error("the refused call cleared the band anyway")
	}
	f.ok("member", "band", "--id", "rina", "--band", "0")
	if strings.Contains(f.ok("member", "list"), `"band"`) {
		t.Error("--band 0 no longer clears the block")
	}
}

// A record's body is a document. Rendered as a YAML scalar it arrives as one
// escaped line, which is the shape `record show` already refuses to print.
func TestACreatedRecordPrintsItsBodyToAPerson(t *testing.T) {
	f := setup(t)
	out := f.human("record", "create", "--kind", "intent", "--slug", "billing", "--title", "Billing")
	for _, want := range []string{"id: i001", "path: intent/i001-billing.md", "# Billing", "## Success criteria"} {
		if !strings.Contains(out, want) {
			t.Errorf("the human stream omits %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, `\n`) {
		t.Errorf("the body was rendered as an escaped scalar:\n%s", out)
	}
}

func TestCLIValidationAndSourceIsolation(t *testing.T) {
	f := setup(t)
	for _, args := range [][]string{{"wat"}, {"member", "add", "--id", "alex"}, {"status", "unexpected"}, {"record", "create", "--kind", "intent", "--slug", "../bad", "--title", "Bad"}} {
		f.fail(args...)
	}
	write(t, filepath.Join(f.root, "sources", "private-fixture.md"), "This content must not be searched.")
	write(t, filepath.Join(f.root, "context", "INDEX.md"), "- [Billing](billing.md): invoices and retries\n- [Shipping](shipping.md): tracking\n")
	result := f.ok("context", "find", "--query", "billing")
	if !strings.Contains(result, "Billing") || strings.Contains(result, "Shipping") || strings.Contains(result, "private-fixture") {
		t.Fatal(result)
	}
	outside := filepath.Join(f.home, "outside.yaml")
	write(t, outside, "members: {}\n")
	if err := os.Remove(filepath.Join(f.root, "members.yaml")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(f.root, "members.yaml")); err != nil {
		t.Skip("symlinks unavailable:", err)
	}
	f.fail("member", "add", "--id", "alex", "--name", "Alex")
	if read(t, outside) != "members: {}\n" {
		t.Fatal("followed workspace symlink")
	}
}

// Marking a plan complete is the one step that reconciles durable knowledge, so
// completion hands the agent the catalog entries scoped to that plan's
// repositories. They are candidates to judge, not a list to rewrite.
func TestCompletionReportsScopedKnowledgeCandidates(t *testing.T) {
	f := setup(t)
	f.repository("api")
	f.repository("web")
	intent := f.intent("billing")
	plan := f.plan(intent.ID, "billing-api", "api")
	write(t, filepath.Join(f.root, "context", "INDEX.md"), strings.Join([]string{
		"# Shared project knowledge",
		"",
		"- [Invoice lifecycle](domains/billing/invoice-lifecycle.md) {api} — when an invoice is voided rather than credited · invoices, dunning · reviewed 2026-02-04",
		"- [Checkout funnel](domains/checkout.md) {web} — which steps a buyer may skip · checkout, funnel · reviewed 2026-02-04",
		"- [Session handling](architecture/sessions.md) {api} {web} — where a session is created and expired · sessions, cookies · reviewed 2026-02-04",
		"The braces stop {api} from also matching {api-gateway}.",
		"",
	}, "\n"))
	out := f.ok("--json", "record", "complete", "--id", plan.ID, "--text", "Landed; no durable concept changed.")
	for _, want := range []string{"Invoice lifecycle", "Session handling", plan.ID} {
		if !strings.Contains(out, want) {
			t.Errorf("completion never offered %q: %s", want, out)
		}
	}
	if strings.Contains(out, "Checkout funnel") {
		t.Errorf("offered an entry outside the plan's repositories: %s", out)
	}
	// Index prose may discuss a scope without cataloguing a note.
	if strings.Contains(out, "braces stop") {
		t.Errorf("offered prose as a candidate: %s", out)
	}
	// A repository ID that is a prefix of another must not widen the candidates.
	f.repository("api-gateway")
	other := f.plan(intent.ID, "gateway", "api-gateway")
	out = f.ok("--json", "record", "complete", "--id", other.ID, "--text", "Landed.")
	if strings.Contains(out, "Invoice lifecycle") {
		t.Errorf("{api} matched {api-gateway}: %s", out)
	}
}

// A workspace repository is shared with a team, so its front page says whose
// workspace it is, and its badges report the versions this workspace received
// rather than whatever the product has published since.
func TestInitializationTitlesTheWorkspaceReadme(t *testing.T) {
	f := setup(t)
	readme := read(t, filepath.Join(f.root, "README.md"))
	if !strings.Contains(readme, `<h1 align="center">Context Circuit - Acme</h1>`) {
		t.Errorf("README is not titled with the workspace name:\n%s", readme[:min(len(readme), 600)])
	}
	if strings.Contains(readme, "img.shields.io/github/v/release") {
		t.Error("the README still looks up the newest release instead of the pinned versions")
	}
	version := strings.TrimSpace(read(t, filepath.Join(f.root, ".context-circuit", "VERSION")))
	cliVersion := strings.TrimSpace(read(t, filepath.Join(f.root, ".context-circuit", "CLI_VERSION")))
	for _, pinned := range []string{"label=workspace&message=v" + version, "label=cli&message=v" + cliVersion} {
		if !strings.Contains(readme, pinned) {
			t.Errorf("README does not pin %q", pinned)
		}
	}
	// The name is placed inside markup, so it is escaped rather than trusted.
	g := fixture{t, filepath.Join(f.home, "escaped"), f.home}
	g.ok("init", "--name", `R&D <b>`, "--purpose", "p", "--member", "a", "--member-name", "A")
	if escaped := read(t, filepath.Join(g.root, "README.md")); !strings.Contains(escaped, "R&amp;D &lt;b&gt;") {
		t.Error("the workspace name reached the heading unescaped")
	}
}
