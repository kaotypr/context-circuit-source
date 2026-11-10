//go:build !windows

package cli_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kaotypr/context-circuit-source/internal/workspace"
)

// guardScript logs where it ran and what it read, one line per run, and fails
// at the point CC_HOOK_FAIL names.
const guardScript = `input=$(cat)
printf '%s\t%s\n' "$(pwd)" "$input" >> "$CC_HOOK_LOG"
case "$input" in *"\"point\":\"$CC_HOOK_FAIL\""*) echo "it failed at $CC_HOOK_FAIL" >&2; exit 1;; esac
echo 'progress'
echo '{"status":"ok","summary":"checked"}'
`

const guardManifest = `id: guard
members: default-on
hooks:
  - at: before approve
    script: {run: [sh, scripts/hook.sh], on_failure: report}
  - at: after implement
    script: {run: [sh, scripts/hook.sh], on_failure: stop}
  - at: before complete
    script: {run: [sh, scripts/hook.sh], on_failure: stop}
  - at: after complete
    script: {run: [sh, scripts/hook.sh], on_failure: stop}
  - at: before deliver
    script:
      run: [sh, scripts/hook.sh]
      on_failure: report
      consent: standing
      action: Move the task to In Review
  - at: after deliver
    script:
      run: [sh, scripts/hook.sh]
      on_failure: report
      consent: each-time
      action: Post the pull request to the team channel
`

type hookLine struct {
	dir   string
	input workspace.ScriptInput
}

// hookLog reads every run the guard script logged, and empties the log.
func hookLog(t *testing.T, path string) []hookLine {
	t.Helper()
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	var lines []hookLine
	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		dir, input, _ := strings.Cut(line, "\t")
		var parsed workspace.ScriptInput
		if err := json.Unmarshal([]byte(input), &parsed); err != nil {
			t.Fatalf("%v: %s", err, line)
		}
		lines = append(lines, hookLine{dir, parsed})
	}
	return lines
}

func evalPath(path string) string {
	real, err := filepath.EvalSymlinks(path)
	if err != nil {
		return path
	}
	return real
}

func guardFixture(t *testing.T) (fixture, string) {
	f := setup(t)
	f.repository("api")
	folder := filepath.Join(f.root, "extensions", "guard")
	write(t, filepath.Join(folder, "extension.yaml"), guardManifest)
	write(t, filepath.Join(folder, "scripts", "hook.sh"), guardScript)
	f.ok("extension", "add", "--id", "guard", "--path", "extensions/guard")
	log := filepath.Join(f.home, "hook.log")
	t.Setenv("CC_HOOK_LOG", log)
	t.Setenv("CC_HOOK_FAIL", "")
	return f, log
}

func extensionStates(t *testing.T, out string) map[string]workspace.ExtensionState {
	t.Helper()
	var states []workspace.ExtensionState
	if err := json.Unmarshal([]byte(out), &states); err != nil {
		t.Fatalf("%v: %s", err, out)
	}
	byID := map[string]workspace.ExtensionState{}
	for _, state := range states {
		byID[state.ID] = state
	}
	return byID
}

// A script runs only after this member allowed this exact version here, and
// any change to the folder asks again. The allow shows what it allows and runs
// nothing.
func TestAScriptRunsOnlyOnceItsVersionIsAllowed(t *testing.T) {
	f, log := guardFixture(t)
	i := f.intent("billing")
	p := f.plan(i.ID, "ledger", "api")

	h := hooks(t, f.ok("extension", "hooks", "--stage", "implement", "--when", "after", "--plan", p.ID))
	if len(h.Ran) != 0 || len(h.Skipped) != 1 || !strings.Contains(h.Skipped[0].Reason, "not allowed on this machine") ||
		!strings.Contains(h.Skipped[0].Reason, "extension allow --id guard") || h.Stopped {
		t.Fatalf("an unallowed script: %+v", h)
	}
	if got := extensionStates(t, f.ok("extension", "list"))["guard"]; got.Trust != workspace.TrustNotAllowed || got.Allow != "extension allow --id guard" || len(got.Scripts) != 6 {
		t.Fatalf("list before allow: %+v", got)
	}
	if hookLog(t, log) != nil {
		t.Fatal("a script ran before it was allowed")
	}

	var allowed workspace.AllowedExtension
	if err := json.Unmarshal([]byte(f.ok("extension", "allow", "--id", "guard")), &allowed); err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(allowed.Version, "sha256:") || allowed.Member != "maya" || len(allowed.Scripts) != 6 || allowed.Source.Path != "extensions/guard" {
		t.Fatalf("allow: %+v", allowed)
	}
	var standing, eachTime bool
	for _, script := range allowed.Scripts {
		standing = standing || script.Consent == "standing" && script.Action == "Move the task to In Review"
		eachTime = eachTime || script.Consent == "each-time" && script.Action == "Post the pull request to the team channel"
		if script.Timeout != "10m0s" || len(script.Run) != 2 {
			t.Fatalf("script shown: %+v", script)
		}
	}
	if !standing || !eachTime {
		t.Fatalf("the consent actions were not shown: %+v", allowed.Scripts)
	}
	if hookLog(t, log) != nil {
		t.Fatal("allow ran a script")
	}
	if !strings.Contains(read(t, filepath.Join(f.root, ".context-circuit", "local", "extension-trust.yaml")), allowed.Version) {
		t.Fatal("the allow was not recorded in the local trust file")
	}
	if got := extensionStates(t, f.ok("extension", "list"))["guard"]; got.Trust != workspace.TrustAllowed || got.Allow != "" {
		t.Fatalf("list after allow: %+v", got)
	}

	// Allowed, the script runs with its input, in the plan's worktree.
	tree := tree(t, f.ok("worktree", "prepare", "--repo", "api", "--plan", p.ID))
	hookLog(t, log)
	h = hooks(t, f.ok("extension", "hooks", "--stage", "implement", "--when", "after", "--plan", p.ID))
	if len(h.Ran) != 1 || h.Ran[0].Status != workspace.ScriptOK || h.Ran[0].Summary != "checked" || h.Stopped {
		t.Fatalf("an allowed script: %+v", h)
	}
	runs := hookLog(t, log)
	if len(runs) != 1 {
		t.Fatalf("runs: %+v", runs)
	}
	run := runs[0]
	realTree, _ := filepath.EvalSymlinks(tree.Path)
	if dir, _ := filepath.EvalSymlinks(run.dir); dir != realTree {
		t.Fatalf("ran in %s, want the worktree %s", run.dir, tree.Path)
	}
	if run.input.Point != "after implement" || run.input.Member != "maya" || run.input.Extension.ID != "guard" ||
		run.input.Record == nil || run.input.Record.ID != p.ID || run.input.Record.Kind != "plan" ||
		run.input.Worktrees["api"] != tree.Path || run.input.Confirmed || run.input.WorkspaceRoot == "" {
		t.Fatalf("input: %+v", run.input)
	}

	// Any change to the folder asks again.
	write(t, filepath.Join(f.root, "extensions", "guard", "scripts", "hook.sh"), guardScript+"# changed\n")
	h = hooks(t, f.ok("extension", "hooks", "--stage", "implement", "--when", "after", "--plan", p.ID))
	if len(h.Ran) != 0 || len(h.Skipped) != 1 || !strings.Contains(h.Skipped[0].Reason, "changed since allowed") {
		t.Fatalf("a changed folder still ran: %+v", h)
	}
	if got := extensionStates(t, f.ok("extension", "list"))["guard"]; got.Trust != workspace.TrustChanged {
		t.Fatalf("list after a change: %+v", got)
	}
	// Line endings alone are not a change.
	write(t, filepath.Join(f.root, "extensions", "guard", "scripts", "hook.sh"), guardScript)
	write(t, filepath.Join(f.root, "extensions", "guard", "extension.yaml"), strings.ReplaceAll(guardManifest, "\n", "\r\n"))
	if got := extensionStates(t, f.ok("extension", "list"))["guard"]; got.Trust != workspace.TrustAllowed {
		t.Fatalf("CRLF line endings asked again: %+v", got)
	}

	// Another member on this machine allows separately.
	f.ok("member", "add", "--id", "lee", "--name", "Lee", "--band", "2")
	f.ok("member", "use", "--id", "lee")
	if got := extensionStates(t, f.ok("extension", "list"))["guard"]; got.Trust != workspace.TrustNotAllowed {
		t.Fatalf("one member's allow covered another: %+v", got)
	}
}

// A required extension whose script is not allowed stops the stage.
func TestARequiredExtensionsUnallowedScriptStops(t *testing.T) {
	f, _ := guardFixture(t)
	write(t, filepath.Join(f.root, "extensions", "guard", "extension.yaml"), strings.Replace(guardManifest, "members: default-on", "members: required", 1))
	h := hooks(t, f.ok("extension", "hooks", "--stage", "implement", "--when", "after", "--repo", "api"))
	if !h.Stopped || !strings.Contains(h.Resolve, "extension allow --id guard") {
		t.Fatalf("an unallowed required script did not stop: %+v", h)
	}
}

// standing runs, each-time waits for the person's yes, and --confirm runs only
// the script that asked, telling it it was confirmed.
func TestConsentDecidesWhenAScriptRuns(t *testing.T) {
	f, log := guardFixture(t)
	repoPath := filepath.Join(f.home, "api")
	f.ok("extension", "allow", "--id", "guard")

	h := hooks(t, f.ok("extension", "hooks", "--stage", "deliver", "--when", "before", "--repo", "api"))
	if len(h.Ran) != 1 || h.Ran[0].Status != workspace.ScriptOK {
		t.Fatalf("a standing script did not run: %+v", h)
	}
	realRepo, _ := filepath.EvalSymlinks(repoPath)
	if runs := hookLog(t, log); len(runs) != 1 || evalPath(runs[0].input.Worktrees["api"]) != realRepo || evalPath(runs[0].dir) != realRepo {
		t.Fatalf("a direct change's script did not get the bound checkout: %+v", runs)
	}

	h = hooks(t, f.ok("extension", "hooks", "--stage", "deliver", "--when", "after", "--repo", "api"))
	if len(h.Ran) != 0 || len(h.ConfirmationRequired) != 1 || h.Stopped {
		t.Fatalf("an each-time script did not wait: %+v", h)
	}
	asked := h.ConfirmationRequired[0]
	if asked.ID != "guard" || asked.Action != "Post the pull request to the team channel" ||
		asked.Confirm != "extension hooks --stage deliver --when after --repo api --confirm guard" {
		t.Fatalf("confirmation: %+v", asked)
	}
	if hookLog(t, log) != nil {
		t.Fatal("an each-time script ran without a yes")
	}
	h = hooks(t, f.ok("extension", "hooks", "--stage", "deliver", "--when", "after", "--repo", "api", "--confirm", "guard"))
	if len(h.Ran) != 1 || len(h.ConfirmationRequired) != 0 {
		t.Fatalf("a confirmed script did not run: %+v", h)
	}
	if runs := hookLog(t, log); len(runs) != 1 || !runs[0].input.Confirmed {
		t.Fatalf("a confirmed script was not told: %+v", runs)
	}
	// --confirm names an each-time script, and runs nothing else.
	if out := f.fail("extension", "hooks", "--stage", "deliver", "--when", "before", "--repo", "api", "--confirm", "guard"); !strings.Contains(out, "no each-time script") {
		t.Fatalf("a confirm of a standing script: %s", out)
	}
	if hookLog(t, log) != nil {
		t.Fatal("a refused confirm ran a script")
	}
	// A CLI-fired stage's scripts run only when its command fires them.
	i := f.intent("billing")
	h = hooks(t, f.ok("extension", "hooks", "--stage", "approve", "--when", "before", "--intent", i.ID))
	if len(h.Ran) != 0 || len(h.Skipped) != 1 || !strings.Contains(h.Skipped[0].Reason, "runs when `record approve` fires") {
		t.Fatalf("before approve ran a script through hooks: %+v", h)
	}
	if out := f.fail("extension", "hooks", "--stage", "approve", "--when", "before", "--intent", i.ID, "--confirm", "guard"); !strings.Contains(out, "record approve") {
		t.Fatalf("confirm at a CLI-fired stage: %s", out)
	}
	if hookLog(t, log) != nil {
		t.Fatal("hooks ran a CLI-fired stage's script")
	}
}

// A CLI-fired stage runs its before point before it writes and its after
// point after: a before stop leaves the command unrun, an after stop leaves
// what it wrote, and a report failure lets the stage go on.
func TestACLIFiredStageRunsItsScriptsAroundItsWrite(t *testing.T) {
	f, log := guardFixture(t)
	f.ok("extension", "allow", "--id", "guard")
	i := f.intent("billing")
	p := f.plan(i.ID, "ledger", "api")

	// A report failure before approval is shown, and approval goes ahead.
	t.Setenv("CC_HOOK_FAIL", "before approve")
	approved := decodeMap(t, f.ok("record", "approve", "--id", i.ID, "--text", "Go ahead."))
	if approved["approved"] != i.ID {
		t.Fatalf("approve: %v", approved)
	}
	blocks, _ := json.Marshal(approved["extensions"])
	if !strings.Contains(string(blocks), `"status":"failed"`) || !strings.Contains(string(blocks), "it failed at before approve") || strings.Contains(string(blocks), `"stopped":true`) {
		t.Fatalf("a report failure: %s", blocks)
	}
	if runs := hookLog(t, log); len(runs) != 1 || runs[0].input.Record == nil || runs[0].input.Record.ID != i.ID || runs[0].input.Record.Kind != "intent" {
		t.Fatalf("approve runs: %+v", runs)
	}

	// A stop before completion: exit 1, nothing written, the block shown.
	t.Setenv("CC_HOOK_FAIL", "before complete")
	code, out, errText := call(f.root, "record", "complete", "--id", p.ID, "--text", "Done.")
	if code != 1 || !strings.Contains(out, `"stopped": true`) || !strings.Contains(errText, "record complete was not run") {
		t.Fatalf("a before stop: exit=%d %s %s", code, out, errText)
	}
	if strings.Contains(read(t, filepath.Join(f.root, filepath.FromSlash(p.Path))), "completed_at") {
		t.Fatal("a stopped completion was written")
	}
	if runs := hookLog(t, log); len(runs) != 1 || runs[0].input.Point != "before complete" {
		t.Fatalf("only the before point runs: %+v", runs)
	}

	// A stop after completion: reported as stopped, and the completion stays.
	t.Setenv("CC_HOOK_FAIL", "after complete")
	completed := decodeMap(t, f.ok("record", "complete", "--id", p.ID, "--text", "Done."))
	blocks, _ = json.Marshal(completed["extensions"])
	var points []workspace.HookResult
	if err := json.Unmarshal(blocks, &points); err != nil {
		t.Fatal(err)
	}
	if completed["completed"] != p.ID || len(points) != 2 || points[0].Stopped || points[1].Point != "after complete" || !points[1].Stopped ||
		len(points[1].Ran) != 1 || points[1].Ran[0].Status != workspace.ScriptFailed {
		t.Fatalf("an after stop: %v", completed)
	}
	if !strings.Contains(read(t, filepath.Join(f.root, filepath.FromSlash(p.Path))), "completed_at") {
		t.Fatal("an after stop undid the completion")
	}
	if runs := hookLog(t, log); len(runs) != 2 || runs[0].input.Point != "before complete" || runs[1].input.Point != "after complete" {
		t.Fatalf("complete runs: %+v", runs)
	}
}

// standing or each-time consent without the action it covers is a check
// finding naming the manifest.
func TestConsentWithoutAnActionIsAFinding(t *testing.T) {
	f, _ := guardFixture(t)
	write(t, filepath.Join(f.root, "extensions", "guard", "extension.yaml"), strings.Replace(guardManifest, "      action: Move the task to In Review\n", "", 1))
	if out := f.fail("check"); !strings.Contains(out, "consent standing needs the action it covers") || !strings.Contains(out, "extensions/guard/extension.yaml") {
		t.Fatalf("check: %s", out)
	}
}
