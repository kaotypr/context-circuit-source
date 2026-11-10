package cli_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kaotypr/context-circuit-source/internal/workspace"
)

// extensionFolder writes an extension folder with a manifest and the
// instruction files it names.
func extensionFolder(t *testing.T, dir, manifest string, files ...string) {
	t.Helper()
	write(t, filepath.Join(dir, "extension.yaml"), manifest)
	for _, file := range files {
		write(t, filepath.Join(dir, filepath.FromSlash(file)), "# "+file+"\n\nCheck parity for what this plan changed.\n")
	}
}

const parityManifest = `id: parity
version: 1.0.0
name: Behavior parity
summary: Checks a rebuilt service behaves like the one it replaces.
members: default-on
settings:
  reference: { required: true, summary: The previous app }
hooks:
  - at: after plan
    instructions: instructions/after-plan.md
  - at: after implement
    instructions: instructions/after-implement.md
    script:
      run: [sh, scripts/parity.sh]
      on_failure: stop
  - at: after approve
    instructions: instructions/after-approve.md
`

func hooks(t *testing.T, out string) workspace.HookResult {
	t.Helper()
	var result struct {
		Extensions []workspace.HookResult `json:"extensions"`
	}
	if err := json.Unmarshal([]byte(out), &result); err != nil || len(result.Extensions) != 1 {
		t.Fatalf("%v: %s", err, out)
	}
	return result.Extensions[0]
}

func decodeMap(t *testing.T, out string) map[string]any {
	t.Helper()
	var result map[string]any
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("%v: %s", err, out)
	}
	return result
}

// A workspace without extensions.yaml behaves as it did before extensions: no
// command's result grows a key, and status says nothing about them.
func TestNoExtensionsChangesNothing(t *testing.T) {
	f := setup(t)
	f.repository("api")
	i := f.intent("billing")
	approved := decodeMap(t, f.ok("record", "approve", "--id", i.ID, "--text", "Go ahead."))
	p := f.plan(i.ID, "api", "api")
	prepared := decodeMap(t, f.ok("worktree", "prepare", "--repo", "api", "--plan", p.ID))
	completed := decodeMap(t, f.ok("record", "complete", "--id", p.ID, "--text", "Done."))
	archived := decodeMap(t, f.ok("record", "archive"))
	status := decodeMap(t, f.ok("status"))
	for name, result := range map[string]map[string]any{"approve": approved, "prepare": prepared, "complete": completed, "archive": archived, "status": status} {
		if _, ok := result["extensions"]; ok {
			t.Errorf("%s returns extensions in a workspace without any: %v", name, result)
		}
	}
	out := f.ok("extension", "hooks", "--stage", "implement", "--when", "after")
	if h := hooks(t, out); len(h.Instructions) != 0 || h.Stopped {
		t.Fatalf("hooks without extensions: %s", out)
	}
	if !strings.Contains(out, "lists no extensions") {
		t.Fatalf("hooks did not say the workspace lists none: %s", out)
	}
	if out := f.ok("extension", "list"); strings.TrimSpace(out) != "[]" {
		t.Fatalf("list: %s", out)
	}
}

func TestALocalExtensionReturnsItsInstructions(t *testing.T) {
	f := setup(t)
	f.repository("api")
	f.repository("web")
	extensionFolder(t, filepath.Join(f.root, "extensions", "parity"), parityManifest,
		"instructions/after-plan.md", "instructions/after-implement.md", "instructions/after-approve.md")

	// A mismatched id is refused before anything is written.
	if out := f.fail("extension", "add", "--id", "other", "--path", "extensions/parity"); !strings.Contains(out, "is extension parity") {
		t.Fatalf("mismatched id: %s", out)
	}
	if _, err := os.Stat(filepath.Join(f.root, "extensions.yaml")); !os.IsNotExist(err) {
		t.Fatal("a refused add wrote extensions.yaml")
	}
	if out := f.fail("extension", "add", "--id", "parity", "--path", "extensions/parity", "--repo", "nope"); !strings.Contains(out, "unknown repository") {
		t.Fatalf("unknown repository: %s", out)
	}
	f.ok("extension", "add", "--id", "parity", "--path", "extensions/parity", "--repo", "api")
	if out := f.fail("extension", "add", "--id", "parity", "--path", "extensions/parity"); !strings.Contains(out, "already lists") {
		t.Fatalf("second add: %s", out)
	}
	// Settings are edited by hand; check names the required one until then.
	if out := f.fail("check"); !strings.Contains(out, "missing the required setting reference") {
		t.Fatalf("check did not name the missing setting: %s", out)
	}
	write(t, filepath.Join(f.root, "extensions.yaml"), "# Our extensions.\nextensions:\n  - id: parity\n    source: {path: extensions/parity}\n    repositories: [api]\n    settings:\n      reference: legacy-orders\n")
	if out := f.ok("check"); strings.Contains(out, "extensions.yaml") {
		t.Fatalf("a valid extension is a finding: %s", out)
	}

	i := f.intent("billing")
	inApi := f.plan(i.ID, "orders", "api")
	inWeb := f.plan(i.ID, "pages", "web")

	h := hooks(t, f.ok("extension", "hooks", "--stage", "implement", "--when", "after", "--plan", inApi.ID))
	if len(h.Instructions) != 1 || h.Instructions[0].File != "extensions/parity/instructions/after-implement.md" {
		t.Fatalf("instructions: %+v", h)
	}
	if h.Instructions[0].Settings["reference"] != "legacy-orders" {
		t.Fatalf("settings did not travel: %+v", h.Instructions[0])
	}
	// The script is not allowed on this machine yet (on Windows it names no
	// Windows command), so it is skipped and says why, and nothing ran.
	if len(h.Skipped) != 1 || len(h.Ran) != 0 || !strings.Contains(h.Skipped[0].Reason, "extension allow --id parity") && !strings.Contains(h.Skipped[0].Reason, "no Windows command") {
		t.Fatalf("the script was not reported as not run: %+v", h)
	}
	h = hooks(t, f.ok("extension", "hooks", "--stage", "implement", "--when", "after", "--plan", inWeb.ID))
	if len(h.Instructions) != 0 || len(h.Skipped) != 1 || !strings.Contains(h.Skipped[0].Reason, "outside repositories") {
		t.Fatalf("a web plan got an api extension: %+v", h)
	}
	// An intent carries no repository list, so a repository-limited extension
	// is skipped there and says why.
	h = hooks(t, f.ok("extension", "hooks", "--stage", "intent", "--when", "after", "--intent", i.ID))
	if len(h.Instructions) != 0 {
		t.Fatalf("intent got instructions: %+v", h)
	}

	// The CLI fires approve itself; hooks refuses its after point and answers
	// its before point with instructions only.
	if out := f.fail("extension", "hooks", "--stage", "approve", "--when", "after", "--intent", i.ID); !strings.Contains(out, "record approve") {
		t.Fatalf("after approve was not refused: %s", out)
	}
	f.ok("extension", "hooks", "--stage", "approve", "--when", "before", "--intent", i.ID)

	// Member switches: off for this member skips the extension.
	f.ok("extension", "disable", "--id", "parity")
	h = hooks(t, f.ok("extension", "hooks", "--stage", "implement", "--when", "after", "--plan", inApi.ID))
	if len(h.Instructions) != 0 || len(h.Skipped) != 1 || h.Skipped[0].Reason != "off for member maya" {
		t.Fatalf("off: %+v", h)
	}
	if !strings.Contains(read(t, filepath.Join(f.root, "extensions.yaml")), "maya:") {
		t.Fatal("disable wrote no switch")
	}
	f.ok("extension", "enable", "--id", "parity")
	if out := f.fail("extension", "enable", "--id", "parity", "--member", "ghost"); !strings.Contains(out, "not on the roster") {
		t.Fatalf("unknown member: %s", out)
	}

	// The worktree stage returns its blocks beside the worktree it prepared.
	prepared := f.ok("worktree", "prepare", "--repo", "api", "--plan", inApi.ID)
	if w := tree(t, prepared); w.Branch == "" {
		t.Fatalf("the worktree is still the result: %s", prepared)
	}
	// Implementation is a stage the CLI never sees, so preparing a worktree
	// where an implement hook applies names the two commands to run there.
	next, _ := decodeMap(t, prepared)["next"].([]any)
	if len(next) != 2 || !strings.Contains(next[0].(string), "--stage implement --when before --plan "+inApi.ID+" --repo api") ||
		!strings.Contains(next[1].(string), "--stage implement --when after --plan "+inApi.ID+" --repo api") {
		t.Fatalf("prepare did not name the implement hooks: %s", prepared)
	}
	// Outside the extension's repositories, nothing is named.
	if web := decodeMap(t, f.ok("worktree", "prepare", "--repo", "web", "--plan", inWeb.ID)); web["next"] != nil {
		t.Fatalf("a web worktree was told to run api hooks: %v", web)
	}
	// Off for this member, nothing is named either.
	f.ok("extension", "disable", "--id", "parity")
	if again := decodeMap(t, f.ok("worktree", "prepare", "--repo", "api", "--plan", inApi.ID)); again["next"] != nil {
		t.Fatalf("an extension off for the member still named hooks: %v", again)
	}
	f.ok("extension", "enable", "--id", "parity")

	// Approval resolves both its points, after it wrote.
	approved := decodeMap(t, f.ok("record", "approve", "--id", i.ID, "--text", "Go ahead."))
	if approved["approved"] != i.ID {
		t.Fatalf("approve: %v", approved)
	}
	if _, ok := approved["extensions"]; !ok {
		// The extension is limited to api, so at an intent it is a skip,
		// which is still reported.
		t.Fatalf("approve returned no extensions block: %v", approved)
	}

	status := decodeMap(t, f.ok("status"))
	summary, _ := status["extensions"].(map[string]any)
	if summary["listed"] != float64(1) || summary["on"] != float64(1) {
		t.Fatalf("status summary: %v", status["extensions"])
	}

	// A worker's brief lists the implement instructions, after its task and
	// before the language boundary and what to return.
	var spec workspace.Dispatch
	if err := json.Unmarshal([]byte(f.ok("agent", "dispatch", "--host", "claude-code", "--role", "worker",
		"--plan", inApi.ID, "--path", f.root, "--task", "Own the ledger only.")), &spec); err != nil {
		t.Fatal(err)
	}
	section := strings.Index(spec.Prompt, "\n# Extensions\n")
	task := strings.Index(spec.Prompt, "\n# Task\n")
	language := strings.Index(spec.Prompt, "\n# Language\n")
	returns := strings.Index(spec.Prompt, "\n# What to return\n")
	if section < 0 || !(task < section && section < language && language < returns) {
		t.Fatalf("extensions section out of place:\n%s", spec.Prompt)
	}
	// On Windows the test's temporary directory may be a short 8.3 name,
	// such as RUNNER~1, while the CLI prints the long one.
	longRoot, err := filepath.EvalSymlinks(f.root)
	if err != nil {
		t.Fatal(err)
	}
	instruction := filepath.Join("extensions", "parity", "instructions", "after-implement.md")
	if !strings.Contains(spec.Prompt, filepath.Join(f.root, instruction)) && !strings.Contains(spec.Prompt, filepath.Join(longRoot, instruction)) {
		t.Fatalf("the brief does not name the instruction file:\n%s", spec.Prompt)
	}
	// The role never reads AGENTS.md, so its brief carries the limits an
	// instruction can't lift, repair loops included.
	for _, limit := range []string{"pushing", "merging", "editing the plan record", "rerunning a check in a loop until it passes", "is not followed"} {
		if !strings.Contains(spec.Prompt, limit) {
			t.Errorf("the brief's extension limits miss %q:\n%s", limit, spec.Prompt)
		}
	}
	// A planner works before plans name repositories, so it sees the
	// repository-limited plan instructions, marked as such.
	f.ok("record", "approve", "--id", i.ID, "--text", "Go ahead.")
	if err := json.Unmarshal([]byte(f.ok("agent", "dispatch", "--host", "claude-code", "--role", "planner",
		"--intent", i.ID, "--path", f.root)), &spec); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(spec.Prompt, "after-plan.md (for plans touching api)") {
		t.Fatalf("planner brief:\n%s", spec.Prompt)
	}
	// A worker on the web plan has nothing to follow, so no section.
	if err := json.Unmarshal([]byte(f.ok("agent", "dispatch", "--host", "claude-code", "--role", "worker",
		"--plan", inWeb.ID, "--path", f.root)), &spec); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(spec.Prompt, "# Extensions") {
		t.Fatalf("a brief with nothing to follow carries the section:\n%s", spec.Prompt)
	}

	// Removing keeps the file and its comment.
	f.ok("extension", "remove", "--id", "parity")
	left := read(t, filepath.Join(f.root, "extensions.yaml"))
	if !strings.Contains(left, "# Our extensions.") || strings.Contains(left, "parity") {
		t.Fatalf("remove: %s", left)
	}
}

func TestRequiredExtensionsAndCheckFindings(t *testing.T) {
	f := setup(t)
	f.repository("api")
	manifest := strings.Replace(parityManifest, "members: default-on", "members: required", 1)
	extensionFolder(t, filepath.Join(f.root, "extensions", "parity"), manifest,
		"instructions/after-plan.md", "instructions/after-implement.md", "instructions/after-approve.md")
	f.ok("extension", "add", "--id", "parity", "--path", "extensions/parity")
	if out := f.fail("extension", "disable", "--id", "parity"); !strings.Contains(out, "required") {
		t.Fatalf("a required extension was switched off: %s", out)
	}
	write(t, filepath.Join(f.root, "extensions.yaml"), `extensions:
  - id: parity
    source: { path: extensions/parity }
    repositories: [nowhere]
    settings:
      reference: legacy
      colour: blue
    members:
      maya: off
      ghost: on
  - id: parity
    source: { path: extensions/parity }
    settings: { reference: legacy }
  - id: remote
    source:
      url: https://example.invalid/ext.git
      ref: v1
`)
	out := f.fail("check")
	for _, want := range []string{
		"names repository nowhere",
		"switches member ghost, who is not on the roster",
		"switches off maya, but its manifest makes it required",
		"sets colour, which its manifest does not declare",
		"lists parity more than once",
		"needs a person: decide which entry stays",
		"records no resolved commit",
		"extension update --id remote --ref v1",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("check missing %q:\n%s", want, out)
		}
	}
	// An off switch on a required extension has no effect.
	h := hooks(t, f.ok("extension", "hooks", "--stage", "plan", "--when", "after"))
	if len(h.Instructions) != 1 || h.Instructions[0].ID != "parity" {
		t.Fatalf("required extension skipped: %+v", h)
	}
}

func TestAManifestBeyondThisCLIOrEscapingItsFolder(t *testing.T) {
	f := setup(t)
	extensionFolder(t, filepath.Join(f.root, "extensions", "future"), "id: future\nmembers: opt-in\nrequires_cli: \">=99.0.0\"\nhooks: []\nfuture_key: true\n")
	if out := f.fail("extension", "add", "--id", "future", "--path", "extensions/future"); !strings.Contains(out, "needs a newer CLI") {
		t.Fatalf("newer manifest: %s", out)
	}
	extensionFolder(t, filepath.Join(f.root, "extensions", "leaky"), "id: leaky\nmembers: default-on\nhooks:\n  - at: after plan\n    instructions: ../../workspace.yaml\n")
	if out := f.fail("extension", "add", "--id", "leaky", "--path", "extensions/leaky"); !strings.Contains(out, "leaves its folder") {
		t.Fatalf("escaping path: %s", out)
	}
	write(t, filepath.Join(f.root, "extensions.yaml"), "extensions:\n  - id: future\n    source: {path: extensions/future}\n  - id: leaky\n    source: {path: extensions/leaky}\n")
	out := f.fail("check")
	for _, want := range []string{"needs a newer CLI", "cc-cli", "needs a person: the extension's author"} {
		if !strings.Contains(out, want) {
			t.Errorf("check missing %q:\n%s", want, out)
		}
	}
	// What cannot be read cannot be shown to be optional, unless it says so:
	// the opt-in one is skipped, the default-on one with a bad path is too,
	// because its header says it may be switched off.
	h := hooks(t, f.ok("extension", "hooks", "--stage", "plan", "--when", "after"))
	if h.Stopped || len(h.Skipped) != 2 {
		t.Fatalf("hooks: %+v", h)
	}
}

// remoteExtension makes a Git repository holding an extension in a subfolder,
// tagged v1 with an annotated tag.
func remoteExtension(t *testing.T, f fixture) (string, string) {
	t.Helper()
	path := filepath.Join(f.home, "registry")
	git(t, f.home, "init", "-q", "-b", "main", path)
	extensionFolder(t, filepath.Join(path, "parity"), strings.Replace(parityManifest, "members: default-on", "members: opt-in", 1),
		"instructions/after-plan.md", "instructions/after-implement.md", "instructions/after-approve.md")
	git(t, path, "add", ".")
	git(t, path, "commit", "-q", "-m", "feat: parity")
	git(t, path, "tag", "-a", "v1", "-m", "v1")
	return path, git(t, path, "rev-parse", "HEAD")
}

func TestARemoteExtensionIsPinnedAndCached(t *testing.T) {
	f := setup(t)
	registry, first := remoteExtension(t, f)
	out := f.ok("extension", "add", "--id", "parity", "--url", registry, "--ref", "v1", "--subpath", "parity")
	var entry workspace.ExtensionEntry
	if err := json.Unmarshal([]byte(out), &entry); err != nil {
		t.Fatal(err)
	}
	// The annotated tag is peeled to the commit it names.
	if entry.Source.Commit != first {
		t.Fatalf("recorded %s, want %s", entry.Source.Commit, first)
	}
	cache := filepath.Join(f.root, ".context-circuit", "local", "extensions", "parity", first)
	if _, err := os.Stat(filepath.Join(cache, "parity", "extension.yaml")); err != nil {
		t.Fatalf("not cached: %v", err)
	}
	if push := git(t, cache, "remote", "get-url", "--push", "origin"); !strings.Contains(push, "READ-ONLY") {
		t.Fatalf("the cached source can push: %s", push)
	}
	// Opt-in: off until this member switches it on.
	h := hooks(t, f.ok("extension", "hooks", "--stage", "plan", "--when", "after"))
	if len(h.Instructions) != 0 || !strings.Contains(h.Skipped[0].Reason, "opt-in") {
		t.Fatalf("opt-in: %+v", h)
	}
	f.ok("extension", "enable", "--id", "parity")
	h = hooks(t, f.ok("extension", "hooks", "--stage", "plan", "--when", "after"))
	if len(h.Instructions) != 1 || !strings.HasPrefix(h.Instructions[0].File, ".context-circuit/local/extensions/parity/"+first+"/parity/") {
		t.Fatalf("remote instructions: %+v", h)
	}

	// The tag moves upstream: sync reports it and keeps the recorded commit.
	write(t, filepath.Join(registry, "parity", "README.md"), "# Parity\n")
	git(t, registry, "add", ".")
	git(t, registry, "commit", "-q", "-m", "docs: readme")
	git(t, registry, "tag", "-f", "-a", "v1", "-m", "v1 moved")
	second := git(t, registry, "rev-parse", "HEAD")
	if out := f.ok("extension", "sync"); !strings.Contains(out, "ref-moved") || !strings.Contains(out, first) {
		t.Fatalf("sync did not report the move: %s", out)
	}
	if !strings.Contains(read(t, filepath.Join(f.root, "extensions.yaml")), first) {
		t.Fatal("sync changed the recorded commit")
	}
	f.ok("extension", "update", "--id", "parity", "--ref", "v1")
	if data := read(t, filepath.Join(f.root, "extensions.yaml")); !strings.Contains(data, second) || !strings.Contains(data, "maya:") {
		t.Fatalf("update: %s", data)
	}
	out = f.ok("extension", "sync", "--prune")
	if !strings.Contains(out, first) || !strings.Contains(out, "pruned") {
		t.Fatalf("prune: %s", out)
	}
	if _, err := os.Stat(cache); !os.IsNotExist(err) {
		t.Fatal("the old commit was not pruned")
	}

	// A clone that has not fetched the commit: check cannot read it, says so
	// apart from its findings, and hooks fetch it.
	if err := os.RemoveAll(filepath.Join(f.root, ".context-circuit", "local", "extensions")); err != nil {
		t.Fatal(err)
	}
	checked := decodeMap(t, f.ok("check"))
	if _, ok := checked["unchecked"]; !ok {
		t.Fatalf("check did not name the uncached remote: %v", checked)
	}
	h = hooks(t, f.ok("extension", "hooks", "--stage", "plan", "--when", "after"))
	if len(h.Instructions) != 1 {
		t.Fatalf("hooks did not fetch the commit: %+v", h)
	}
}

// A required remote this machine cannot read stops the stage, naming sync.
func TestAnUnreadableRequiredExtensionStops(t *testing.T) {
	f := setup(t)
	registry, _ := remoteExtension(t, f)
	f.ok("extension", "add", "--id", "parity", "--url", registry, "--ref", "v1", "--subpath", "parity")
	if err := os.RemoveAll(filepath.Join(f.root, ".context-circuit", "local", "extensions")); err != nil {
		t.Fatal(err)
	}
	if err := os.RemoveAll(registry); err != nil {
		t.Fatal(err)
	}
	h := hooks(t, f.ok("extension", "hooks", "--stage", "plan", "--when", "after"))
	if !h.Stopped || !strings.Contains(h.Resolve, "extension sync") {
		t.Fatalf("an unreadable extension did not stop: %+v", h)
	}
}

const roleManifest = `id: parity
version: 1.0.0
members: default-on
hooks:
  - at: after implement
    instructions: instructions/after-implement.md
roles:
  - id: parity-verifier
    description: Checks a rebuilt service behaves like the one it replaces.
    access: read-only
    stages: [implement]
    like: reviewer
    instructions: roles/parity-verifier.md
    returns: [Verdict, Findings, Not checked]
  - id: parity-fixer
    description: Brings a rebuilt service back in line with the one it replaces.
    access: edit
    like: worker
    instructions: roles/parity-fixer.md
    returns: [Fixed]
role_additions:
  - role: worker
    instructions: roles/worker-parity.md
    returns: [Parity]
`

func dispatchSpec(t *testing.T, out string) workspace.Dispatch {
	t.Helper()
	var spec workspace.Dispatch
	if err := json.Unmarshal([]byte(out), &spec); err != nil {
		t.Fatalf("%v: %s", err, out)
	}
	return spec
}

// inOrder fails unless each section appears in the prompt after the one before.
func inOrder(t *testing.T, prompt string, sections ...string) {
	t.Helper()
	last := -1
	for _, section := range sections {
		at := strings.Index(prompt, section)
		if at < 0 || at < last {
			t.Fatalf("%q is missing or out of order:\n%s", section, prompt)
		}
		last = at
	}
}

// An extension role is installed beside the shipped four, briefed like the
// shipped role its access makes it, and removed again with its extension,
// except where a person edited its definition.
func TestExtensionRolesAreInstalledDispatchedAndRemoved(t *testing.T) {
	f := setup(t)
	f.repository("api")
	extensionFolder(t, filepath.Join(f.root, "extensions", "parity"), roleManifest,
		"instructions/after-implement.md", "roles/parity-verifier.md", "roles/parity-fixer.md", "roles/worker-parity.md")
	worker := read(t, filepath.Join(f.root, ".claude/agents/cc-worker.md"))

	added := decodeMap(t, f.ok("extension", "add", "--id", "parity", "--path", "extensions/parity"))
	if note, _ := added["setup_required"].(string); !strings.Contains(note, "agent setup") {
		t.Fatalf("add did not name agent setup: %v", added)
	}
	if note, _ := decodeMap(t, f.ok("status"))["setup_required"].(string); !strings.Contains(note, "agent setup") {
		t.Fatal("status did not report the stale definitions")
	}

	f.ok("agent", "configure", "--host", "claude-code", "--role", "reviewer", "--model", "review-model", "--effort", "high")
	setupOut := decodeMap(t, f.ok("agent", "setup"))
	if written, _ := setupOut["written"].([]any); len(written) < 18 {
		t.Fatalf("setup did not write every definition: %v", setupOut)
	}
	if _, ok := decodeMap(t, f.ok("status"))["setup_required"]; ok {
		t.Fatal("status still reports setup after it ran")
	}
	verifier := filepath.Join(f.root, ".claude/agents/cc-x-parity-parity-verifier.md")
	text := read(t, verifier)
	inOrder(t, text, "description: ", "Start it only with the prompt `agent dispatch --role x-parity-parity-verifier` returns, passed unmodified.", "---\n\nWork only from the brief `agent dispatch` writes for this role.", "agent dispatch --host claude-code --role x-parity-parity-verifier", "Check parity")
	for _, want := range []string{"name: cc-x-parity-parity-verifier", "tools: Read, Glob, Grep", "review-model", "Check parity", "- Verdict\n- Findings\n- Not checked", "added by the workspace extension parity"} {
		if !strings.Contains(text, want) {
			t.Errorf("verifier definition misses %q:\n%s", want, text)
		}
	}
	codex := read(t, filepath.Join(f.root, ".codex/agents/cc-x-parity-parity-verifier.toml"))
	if !strings.Contains(codex, `name = "cc_x_parity_parity_verifier"`) || !strings.Contains(codex, `sandbox_mode = "read-only"`) || !strings.Contains(codex, "agent dispatch --host codex --role x-parity-parity-verifier") || !strings.Contains(codex, "Start it only with the prompt") {
		t.Fatalf("codex definition:\n%s", codex)
	}
	if cursor := read(t, filepath.Join(f.root, ".cursor/agents/cc-x-parity-parity-fixer.md")); !strings.Contains(cursor, "readonly: false") {
		t.Fatalf("an edit role is read-only on cursor:\n%s", cursor)
	}
	if fixer := read(t, filepath.Join(f.root, ".claude/agents/cc-x-parity-parity-fixer.md")); strings.Contains(fixer, "tools:") {
		t.Fatalf("an edit role was given the read-only tools:\n%s", fixer)
	}
	// Additions travel in the brief only: the worker's definition is untouched.
	if read(t, filepath.Join(f.root, ".claude/agents/cc-worker.md")) != worker {
		t.Fatal("an addition changed the worker's definition")
	}

	// Its own tiering replaces what it takes from like.
	f.ok("agent", "configure", "--host", "claude-code", "--role", "x-parity-parity-verifier", "--model", "verifier-model", "--effort", "low")
	if note, _ := decodeMap(t, f.ok("status"))["setup_required"].(string); note == "" {
		t.Fatal("a changed setting left the definition current")
	}
	f.ok("agent", "setup", "--host", "claude-code")
	if text := read(t, verifier); !strings.Contains(text, "verifier-model") || strings.Contains(text, "review-model") {
		t.Fatalf("the role's own setting was not applied:\n%s", text)
	}
	if out := f.fail("agent", "configure", "--host", "claude-code", "--role", "x-parity-nobody", "--model", "m", "--effort", "low"); !strings.Contains(out, "declares agent role x-parity-nobody") {
		t.Fatalf("an unknown extension role was configured: %s", out)
	}

	i := f.intent("orders")
	p := f.plan(i.ID, "orders", "api")

	// A read-only role is briefed like the explorer.
	spec := dispatchSpec(t, f.ok("agent", "dispatch", "--host", "claude-code", "--role", "x-parity-parity-verifier", "--plan", p.ID, "--path", f.root))
	if !spec.ReadOnly || spec.AgentType != "cc-x-parity-parity-verifier" || !spec.DefinitionInstalled || spec.Setting.Model != "verifier-model" ||
		spec.Report != "Report this role's result to the person under its headings, in this order: Verdict, Findings, Not checked." {
		t.Fatalf("verifier spec: %+v", spec)
	}
	inOrder(t, spec.Prompt, "# Working directory", "# Ownership\n\nYou are reading, not changing.", "# Plan "+p.ID,
		"# Extensions", "after implement, parity:", "# Language", "# What to return\n\n# roles/parity-verifier.md",
		"- Verdict\n- Findings\n- Not checked", "no push, pull or merge request", "launches no other agent", "Change nothing in the working directory.")
	if strings.Contains(spec.Prompt, "Repository assignment") {
		t.Fatalf("a read-only role was given a worker's assignment:\n%s", spec.Prompt)
	}
	// Naming the plan's own repository is accepted and changes nothing; any
	// other is refused.
	if same := dispatchSpec(t, f.ok("agent", "dispatch", "--host", "claude-code", "--role", "x-parity-parity-verifier", "--plan", p.ID, "--repo", "api", "--path", f.root)); same.Prompt != spec.Prompt {
		t.Fatalf("--repo changed a read-only role's brief:\n%s", same.Prompt)
	}
	if out := f.fail("agent", "dispatch", "--host", "claude-code", "--role", "x-parity-parity-verifier", "--plan", p.ID, "--repo", "web", "--path", f.root); !strings.Contains(out, "is not one of "+p.ID+"'s repositories") {
		t.Fatalf("a read-only role took a repository outside its plan: %s", out)
	}
	f.fail("agent", "dispatch", "--host", "claude-code", "--role", "x-parity-parity-verifier", "--repo", "api", "--task", "Compare.", "--path", f.root)
	f.fail("agent", "dispatch", "--host", "claude-code", "--role", "x-parity-parity-verifier", "--path", f.root)
	f.ok("agent", "dispatch", "--host", "codex", "--role", "x-parity-parity-verifier", "--task", "Compare the orders endpoints.", "--path", f.root)

	// An edit role is briefed like the worker, against a plan.
	if out := f.fail("agent", "dispatch", "--host", "claude-code", "--role", "x-parity-parity-fixer", "--task", "Fix it.", "--path", f.root); !strings.Contains(out, "pass --plan") {
		t.Fatalf("an edit role ran without a plan: %s", out)
	}
	spec = dispatchSpec(t, f.ok("agent", "dispatch", "--host", "codex", "--role", "x-parity-parity-fixer", "--plan", p.ID, "--repo", "api", "--path", f.root))
	if spec.ReadOnly || spec.AgentType != "cc_x_parity_parity_fixer" {
		t.Fatalf("fixer spec: %+v", spec)
	}
	inOrder(t, spec.Prompt, "# Ownership\n\nYou are the sole owner", "Repository assignment: api", "# What to return", "- Fixed", "Commit only on the branch")

	// A worker's brief carries the addition, its heading after the core ones,
	// and its access is unchanged.
	spec = dispatchSpec(t, f.ok("agent", "dispatch", "--host", "claude-code", "--role", "worker", "--plan", p.ID, "--path", f.root))
	if spec.ReadOnly {
		t.Fatal("an addition changed the worker's access")
	}
	if spec.Report != "" {
		t.Fatalf("a shipped role was given a report line: %q", spec.Report)
	}
	inOrder(t, spec.Prompt, "# Ownership\n\nYou are the sole owner", "# Extensions\n\nThis workspace's extensions add these instructions to the implement stage and the worker role.",
		"after implement, parity:", "worker role, parity: ", filepath.Join("extensions", "parity", "roles", "worker-parity.md"),
		"# What to return\n\n"+workspace.Roles["worker"].Instructions, "- Parity")
	// The other shipped roles get nothing.
	spec = dispatchSpec(t, f.ok("agent", "dispatch", "--host", "claude-code", "--role", "explorer", "--task", "Where are orders?", "--path", f.root))
	if strings.Contains(spec.Prompt, "Parity") || strings.Contains(spec.Prompt, "# Extensions") {
		t.Fatalf("an explorer got the worker's addition:\n%s", spec.Prompt)
	}

	// Off for this member, the roles leave: unedited definitions are removed,
	// an edited one is kept and named.
	fixer := filepath.Join(f.root, ".claude/agents/cc-x-parity-parity-fixer.md")
	write(t, fixer, "my own fixer\n")
	disabled := decodeMap(t, f.ok("extension", "disable", "--id", "parity"))
	if _, ok := disabled["setup_required"]; !ok {
		t.Fatalf("disable did not name setup: %v", disabled)
	}
	if out := f.fail("agent", "dispatch", "--host", "claude-code", "--role", "x-parity-parity-verifier", "--plan", p.ID, "--path", f.root); !strings.Contains(out, "no extension on for this member") {
		t.Fatalf("a role of a disabled extension was dispatched: %s", out)
	}
	var result workspace.AgentSetup
	if err := json.Unmarshal([]byte(f.ok("agent", "setup")), &result); err != nil {
		t.Fatal(err)
	}
	if len(result.Kept) != 1 || result.Kept[0].Path != ".claude/agents/cc-x-parity-parity-fixer.md" || len(result.Removed) != 5 {
		t.Fatalf("setup after disable: %+v", result)
	}
	if read(t, fixer) != "my own fixer\n" {
		t.Fatal("an edited definition was overwritten")
	}
	for _, gone := range []string{verifier, filepath.Join(f.root, ".codex/agents/cc-x-parity-parity-fixer.toml")} {
		if _, err := os.Stat(gone); !os.IsNotExist(err) {
			t.Fatalf("%s was not removed", gone)
		}
	}
	if _, ok := decodeMap(t, f.ok("status"))["setup_required"]; ok {
		t.Fatal("status reports setup after the definitions were removed")
	}
	// Back on, setup will not overwrite the person's file.
	if enabled := decodeMap(t, f.ok("extension", "enable", "--id", "parity")); enabled["setup_required"] == nil {
		t.Fatalf("enable did not name setup: %v", enabled)
	}
	if out := f.fail("agent", "setup", "--host", "claude-code"); !strings.Contains(out, "preserve customized agent definition") {
		t.Fatalf("setup overwrote an edited definition: %s", out)
	}
	if err := os.Remove(fixer); err != nil {
		t.Fatal(err)
	}
	f.ok("agent", "setup")
	removed := decodeMap(t, f.ok("extension", "remove", "--id", "parity"))
	if removed["setup_required"] == nil {
		t.Fatalf("remove did not name setup: %v", removed)
	}
	result = workspace.AgentSetup{}
	if err := json.Unmarshal([]byte(f.ok("agent", "setup")), &result); err != nil {
		t.Fatal(err)
	}
	if len(result.Removed) != 6 || len(result.Kept) != 0 {
		t.Fatalf("setup after remove: %+v", result)
	}
}

// A workspace without extension roles has nothing to set up, and its
// extension commands say nothing about setup.
func TestHookOnlyExtensionsNeedNoSetup(t *testing.T) {
	f := setup(t)
	extensionFolder(t, filepath.Join(f.root, "extensions", "style"), "id: style\nmembers: default-on\nhooks:\n  - at: after implement\n    instructions: a.md\n", "a.md")
	if added := decodeMap(t, f.ok("extension", "add", "--id", "style", "--path", "extensions/style")); added["setup_required"] != nil {
		t.Fatalf("a hook-only extension asked for setup: %v", added)
	}
	if _, ok := decodeMap(t, f.ok("status"))["setup_required"]; ok {
		t.Fatal("status asked for setup without extension roles")
	}
}

// check names each manifest problem with a role, and the clash between two
// extensions that no single manifest shows.
func TestExtensionRoleCheckFindings(t *testing.T) {
	base := "id: parity\nmembers: default-on\nhooks: []\n"
	role := func(fields string) string {
		return base + "roles:\n  - id: verifier\n    description: Checks parity.\n    access: read-only\n    instructions: roles/verifier.md\n" + fields
	}
	for _, tc := range []struct{ name, manifest, want string }{
		{"id not a slug", base + "roles:\n  - id: Parity Verifier\n    description: x\n    access: read-only\n    instructions: roles/verifier.md\n", "role id"},
		{"id twice", role("  - id: verifier\n    description: Again.\n    access: edit\n    instructions: roles/verifier.md\n"), "role verifier appears twice"},
		{"like not shipped", role("    like: boss\n"), "like names a shipped role"},
		{"bad access", base + "roles:\n  - id: verifier\n    description: x\n    access: write\n    instructions: roles/verifier.md\n", "access is read-only or edit"},
		{"unknown stage", role("    stages: [ship]\n"), "unknown stage"},
		{"core heading", base + "role_additions:\n  - role: planner\n    returns: [Approach]\n", "core heading of that role"},
		{"not a shipped role", base + "role_additions:\n  - role: manager\n    returns: [Notes]\n", "not a shipped role"},
		{"escaping file", base + "roles:\n  - id: verifier\n    description: x\n    access: read-only\n    instructions: ../outside.md\n", "leaves its folder"},
		{"missing file", base + "roles:\n  - id: verifier\n    description: x\n    access: read-only\n    instructions: roles/absent.md\n", "cannot be read"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := setup(t)
			extensionFolder(t, filepath.Join(f.root, "extensions", "parity"), tc.manifest, "roles/verifier.md")
			write(t, filepath.Join(f.root, "extensions.yaml"), "extensions:\n  - id: parity\n    source: {path: extensions/parity}\n")
			if out := f.fail("check"); !strings.Contains(out, tc.want) {
				t.Fatalf("check did not report %q: %s", tc.want, out)
			}
		})
	}

	f := setup(t)
	for _, id := range []string{"one", "two"} {
		extensionFolder(t, filepath.Join(f.root, "extensions", id), "id: "+id+"\nmembers: default-on\nhooks: []\nrole_additions:\n  - role: worker\n    returns: [Parity]\n")
	}
	write(t, filepath.Join(f.root, "extensions.yaml"), "extensions:\n  - id: one\n    source: {path: extensions/one}\n  - id: two\n    source: {path: extensions/two}\n")
	if out := f.fail("check"); !strings.Contains(out, `extensions one and two both add the heading \"Parity\" to the worker role`) {
		t.Fatalf("check did not report the shared heading: %s", out)
	}
}
