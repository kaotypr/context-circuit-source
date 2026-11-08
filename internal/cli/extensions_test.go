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
	if len(h.Skipped) != 1 || !strings.Contains(h.Skipped[0].Reason, "runs no extension scripts") {
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
	if !strings.Contains(spec.Prompt, filepath.Join(f.root, "extensions", "parity", "instructions", "after-implement.md")) {
		t.Fatalf("the brief does not name the instruction file:\n%s", spec.Prompt)
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
