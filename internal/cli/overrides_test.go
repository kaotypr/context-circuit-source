package cli_test

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kaotypr/context-circuit-source/internal/workspace"
)

const overrideManifest = `id: parity
version: 1.0.0
members: default-on
hooks: []
roles:
  - id: parity-verifier
    description: Checks a rebuilt service behaves like the one it replaces.
    access: read-only
    instructions: roles/parity-verifier.md
    returns: [Verdict, Findings]
overrides:
  - default: repair-loop
    instructions: instructions/repair-loop.md
    verifier: parity-verifier
    rounds: 3
  - default: risk-review
    instructions: instructions/risk-review.md
  - default: delegation
    stage: implement
    role: worker
`

var overrideFiles = []string{"roles/parity-verifier.md", "instructions/repair-loop.md", "instructions/risk-review.md"}

// overrideWorkspace is a workspace listing the parity extension with its
// overrides declared and none allowed, and a plan in api.
func overrideWorkspace(t *testing.T) (fixture, string) {
	t.Helper()
	f := setup(t)
	f.repository("api")
	extensionFolder(t, filepath.Join(f.root, "extensions", "parity"), overrideManifest, overrideFiles...)
	write(t, filepath.Join(f.root, "extensions.yaml"), "# Our extensions.\nextensions:\n  - id: parity\n    source: {path: extensions/parity}\n    # Only the api.\n    repositories: [api]\n")
	i := f.intent("orders")
	p := f.plan(i.ID, "orders", "api")
	return f, p.ID
}

func afterImplement(t *testing.T, f fixture, plan string) workspace.HookResult {
	t.Helper()
	return hooks(t, f.ok("extension", "hooks", "--stage", "implement", "--when", "after", "--plan", plan))
}

// Nothing a manifest declares applies until a person allows it, and an allowed
// override reaches the coordinator through the hook point it belongs to.
func TestOverridesApplyOnlyWhenAllowed(t *testing.T) {
	f, plan := overrideWorkspace(t)
	if h := afterImplement(t, f, plan); len(h.Overrides) != 0 {
		t.Fatalf("an override applied without an allowance: %+v", h)
	}
	if _, ok := decodeMap(t, f.ok("status"))["overrides"]; ok {
		t.Fatal("status lists overrides nobody allowed")
	}
	if out := f.fail("agent", "dispatch", "--host", "claude-code", "--role", "x-parity-parity-verifier", "--plan", plan, "--override", "repair-loop", "--path", f.root); !strings.Contains(out, "repair-loop is not allowed in this workspace") {
		t.Fatalf("an unallowed override was dispatched: %s", out)
	}

	allowed := decodeMap(t, f.ok("extension", "override", "--id", "parity", "--allow", "repair-loop", "--rounds", "2"))
	if overrides, _ := allowed["overrides"].(map[string]any); overrides["repair-loop"] == nil {
		t.Fatalf("allow result: %v", allowed)
	}
	file := read(t, filepath.Join(f.root, "extensions.yaml"))
	for _, want := range []string{"# Our extensions.", "# Only the api.", "repair-loop:", "rounds: 2"} {
		if !strings.Contains(file, want) {
			t.Fatalf("extensions.yaml lost or lacks %q:\n%s", want, file)
		}
	}
	f.ok("extension", "override", "--id", "parity", "--allow", "delegation")

	h := afterImplement(t, f, plan)
	if len(h.Overrides) != 1 {
		t.Fatalf("after implement overrides: %+v", h.Overrides)
	}
	loop := h.Overrides[0]
	if loop.Default != "repair-loop" || loop.Extension != "parity" || loop.Rounds != 2 || loop.Verifier != "x-parity-parity-verifier" || loop.Instructions != "extensions/parity/instructions/repair-loop.md" {
		t.Fatalf("repair-loop block: %+v", loop)
	}
	if before := hooks(t, f.ok("extension", "hooks", "--stage", "implement", "--when", "before", "--plan", plan)); len(before.Overrides) != 1 || before.Overrides[0].Default != "delegation" || before.Overrides[0].Role != "worker" {
		t.Fatalf("before implement overrides: %+v", before.Overrides)
	}
	if before := hooks(t, f.ok("extension", "hooks", "--stage", "plan", "--when", "before", "--plan", plan)); len(before.Overrides) != 0 {
		t.Fatalf("delegation of implement applied at plan: %+v", before.Overrides)
	}

	var status struct {
		Overrides []workspace.OverrideStatus `json:"overrides"`
	}
	if err := json.Unmarshal([]byte(f.ok("status")), &status); err != nil {
		t.Fatal(err)
	}
	want := []workspace.OverrideStatus{{Default: "repair-loop", Extension: "parity", Rounds: 2, OnForMember: true}, {Default: "delegation", Extension: "parity", OnForMember: true}}
	if len(status.Overrides) != 2 || status.Overrides[0] != want[0] || status.Overrides[1] != want[1] {
		t.Fatalf("status overrides: %+v", status.Overrides)
	}

	// Off for the member, the override is listed but returned nowhere.
	f.ok("extension", "disable", "--id", "parity")
	if h := afterImplement(t, f, plan); len(h.Overrides) != 0 {
		t.Fatalf("an override applied with its extension off: %+v", h.Overrides)
	}
	if !strings.Contains(f.ok("status"), `"on_for_member": false`) {
		t.Fatal("status did not say the override is off for the member")
	}
	if out := f.fail("agent", "dispatch", "--host", "claude-code", "--role", "worker", "--plan", plan, "--task", "Fix the header.", "--override", "repair-loop", "--path", f.root); !strings.Contains(out, "off for member") {
		t.Fatalf("a dispatch ran under an override of an extension that is off: %s", out)
	}
	f.ok("extension", "enable", "--id", "parity")

	// Outside the extension's repositories, nothing applies.
	f.repository("web")
	other := f.plan(f.intent("landing").ID, "landing", "web")
	if h := hooks(t, f.ok("extension", "hooks", "--stage", "implement", "--when", "after", "--plan", other.ID)); len(h.Overrides) != 0 {
		t.Fatalf("an override applied outside its repositories: %+v", h.Overrides)
	}

	// Revoking the last allowance drops the key and keeps the comments.
	f.ok("extension", "override", "--id", "parity", "--revoke", "delegation")
	f.ok("extension", "override", "--id", "parity", "--revoke", "repair-loop")
	file = read(t, filepath.Join(f.root, "extensions.yaml"))
	if strings.Contains(file, "overrides") || !strings.Contains(file, "# Only the api.") {
		t.Fatalf("revoke left:\n%s", file)
	}
	if h := afterImplement(t, f, plan); len(h.Overrides) != 0 {
		t.Fatalf("a revoked override applied: %+v", h.Overrides)
	}
}

// extension override refuses what a person cannot allow.
func TestOverrideAllowanceRefusals(t *testing.T) {
	f, _ := overrideWorkspace(t)
	for _, tc := range []struct {
		args []string
		want string
	}{
		{[]string{"--allow", "approve"}, `"approve" is not an overridable default`},
		{[]string{"--allow", "verify-during-execution"}, "does not declare the override verify-during-execution"},
		{[]string{"--allow", "repair-loop", "--rounds", "4"}, "never raise them"},
		{[]string{"--allow", "repair-loop", "--rounds", "0"}, ""},
		{[]string{"--allow", "risk-review", "--rounds", "2"}, "--rounds belongs to repair-loop"},
		{[]string{"--revoke", "risk-review"}, "no allowance for risk-review"},
		{[]string{}, "pass --allow NAME"},
		{[]string{"--allow", "risk-review", "--revoke", "risk-review"}, "pass --allow NAME"},
	} {
		args := append([]string{"extension", "override", "--id", "parity"}, tc.args...)
		if tc.want == "" {
			f.ok(args...)
			continue
		}
		if out := f.fail(args...); !strings.Contains(out, tc.want) {
			t.Errorf("%v: want %q, got %s", tc.args, tc.want, out)
		}
	}
	if out := f.fail("extension", "override", "--id", "nobody", "--allow", "repair-loop"); !strings.Contains(out, "lists no extension nobody") {
		t.Errorf("an unlisted extension was given an allowance: %s", out)
	}
}

// An allowed override unlocks exactly the dispatches its default needs, and
// every other use is refused, naming what is missing.
func TestOverrideDispatch(t *testing.T) {
	f, plan := overrideWorkspace(t)
	f.ok("extension", "override", "--id", "parity", "--allow", "repair-loop", "--rounds", "2")
	f.ok("extension", "override", "--id", "parity", "--allow", "risk-review")
	dispatch := func(args ...string) []string {
		return append([]string{"agent", "dispatch", "--host", "claude-code", "--path", f.root}, args...)
	}

	spec := dispatchSpec(t, f.ok(dispatch("--role", "x-parity-parity-verifier", "--plan", plan, "--override", "repair-loop", "--round", "1")...))
	if !spec.ReadOnly {
		t.Fatalf("verifier spec: %+v", spec)
	}
	inOrder(t, spec.Prompt, "# Plan "+plan, "# Override\n\nThis dispatch runs under the override repair-loop that this workspace allowed for the extension parity, round 1 of 2.",
		filepath.Join("extensions", "parity", "instructions", "repair-loop.md"), "# What to return")
	spec = dispatchSpec(t, f.ok(dispatch("--role", "worker", "--plan", plan, "--task", "Findings: the CSV has no header row.", "--override", "repair-loop", "--round", "1")...))
	inOrder(t, spec.Prompt, "# Task\n\nFindings: the CSV has no header row.", "round 1 of 2", "# What to return\n\n"+workspace.Roles["worker"].Instructions)
	if strings.Contains(workspace.Roles["worker"].Instructions, "override") {
		t.Fatal("the worker's own limits changed")
	}
	// risk-review unlocks the shipped reviewer without a person's request.
	f.fail(dispatch("--role", "reviewer", "--plan", plan, "--task", "Review the diff.")...)
	spec = dispatchSpec(t, f.ok(dispatch("--role", "reviewer", "--plan", plan, "--task", "Review the diff.", "--override", "risk-review")...))
	if !strings.Contains(spec.Prompt, "the override risk-review") {
		t.Fatalf("reviewer brief:\n%s", spec.Prompt)
	}

	for _, tc := range []struct {
		args []string
		want string
	}{
		{[]string{"--role", "x-parity-parity-verifier", "--plan", plan, "--override", "approve"}, "--override names an overridable default"},
		{[]string{"--role", "x-parity-parity-verifier", "--plan", plan, "--override", "verify-during-execution"}, "verify-during-execution is not allowed in this workspace"},
		{[]string{"--role", "x-parity-parity-verifier", "--task", "Check.", "--override", "repair-loop"}, "pass --plan"},
		{[]string{"--role", "worker", "--plan", plan, "--override", "repair-loop"}, "pass the findings as --task"},
		{[]string{"--role", "worker", "--plan", plan, "--task", "Findings.", "--override", "repair-loop", "--round", "2"}, "round 2 is the last of 2"},
		{[]string{"--role", "x-parity-parity-verifier", "--plan", plan, "--override", "repair-loop", "--round", "3"}, "--round is between 1 and 2"},
		{[]string{"--role", "reviewer", "--plan", plan, "--task", "Review.", "--override", "repair-loop"}, "unlocks x-parity-parity-verifier, or a worker with the findings as --task, not reviewer"},
		{[]string{"--role", "worker", "--plan", plan, "--task", "Findings.", "--override", "risk-review"}, "only repair-loop dispatches a worker"},
		{[]string{"--role", "explorer", "--task", "Where?", "--round", "1"}, "--round goes with --override"},
	} {
		if out := f.fail(dispatch(tc.args...)...); !strings.Contains(out, tc.want) {
			t.Errorf("%v: want %q, got %s", tc.args, tc.want, out)
		}
	}
	f.ok("extension", "override", "--id", "parity", "--allow", "delegation")
	if out := f.fail(dispatch("--role", "worker", "--plan", plan, "--override", "delegation")...); !strings.Contains(out, "unlocks no dispatch") {
		t.Errorf("delegation unlocked a dispatch: %s", out)
	}
}

// check names each allowance and manifest problem, every one for a person.
func TestOverrideCheckFindings(t *testing.T) {
	manifest := func(overrides string) string {
		return "id: parity\nmembers: default-on\nhooks: []\nroles:\n  - id: verifier\n    description: Checks parity.\n    access: read-only\n    instructions: roles/verifier.md\n  - id: fixer\n    description: Fixes parity.\n    access: edit\n    instructions: roles/verifier.md\noverrides:\n" + overrides
	}
	loop := "  - default: repair-loop\n    verifier: verifier\n    rounds: 3\n"
	for _, tc := range []struct{ name, manifest, allow, want string }{
		{"approve override", manifest("  - default: approve\n"), "", `override \"approve\" is not an overridable default`},
		{"edit verifier", manifest("  - default: repair-loop\n    verifier: fixer\n    rounds: 2\n"), "", "a verifier is read-only"},
		{"unknown verifier", manifest("  - default: verify-during-execution\n    verifier: nobody\n"), "", "verifier nobody is neither the shipped reviewer"},
		{"rounds out of range", manifest("  - default: repair-loop\n    verifier: verifier\n    rounds: 11\n"), "", "rounds is between 1 and 10"},
		{"missing review instructions", manifest("  - default: risk-review\n"), "", "names no instructions"},
		{"delegation to read-only", manifest("  - default: delegation\n    stage: implement\n    role: verifier\n"), "", "needs a role with edit access"},
		{"undeclared allowance", manifest(loop), "verify-during-execution: {}", "allows verify-during-execution, which its manifest does not declare"},
		{"allowance outside the list", manifest(loop), "push: {}", "allows push, which is not an overridable default"},
		{"rounds above the manifest", manifest(loop), "repair-loop: {rounds: 5}", "above the 3 its manifest asks for"},
		{"rounds out of range allowed", manifest(loop), "repair-loop: {rounds: 12}", "outside 1-10"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := setup(t)
			extensionFolder(t, filepath.Join(f.root, "extensions", "parity"), tc.manifest, "roles/verifier.md")
			entry := "extensions:\n  - id: parity\n    source: {path: extensions/parity}\n"
			if tc.allow != "" {
				entry += "    overrides:\n      " + tc.allow + "\n"
			}
			write(t, filepath.Join(f.root, "extensions.yaml"), entry)
			out := f.fail("check")
			if !strings.Contains(out, tc.want) {
				t.Fatalf("check did not report %q: %s", tc.want, out)
			}
			if tc.allow != "" && !strings.Contains(out, "needs a person") {
				t.Fatalf("an allowance finding is not for a person: %s", out)
			}
		})
	}

	t.Run("unreadable extension", func(t *testing.T) {
		f := setup(t)
		write(t, filepath.Join(f.root, "extensions.yaml"), "extensions:\n  - id: parity\n    source: {path: extensions/parity}\n    overrides:\n      repair-loop: {}\n")
		if out := f.fail("check"); !strings.Contains(out, "allows repair-loop, but its manifest cannot be read") {
			t.Fatalf("check: %s", out)
		}
	})

	t.Run("conflicting delegation", func(t *testing.T) {
		f := setup(t)
		extensionFolder(t, filepath.Join(f.root, "extensions", "one"), "id: one\nmembers: default-on\nhooks: []\noverrides:\n  - default: delegation\n    stage: implement\n    role: worker\n")
		extensionFolder(t, filepath.Join(f.root, "extensions", "two"), "id: two\nmembers: default-on\nhooks: []\nroles:\n  - id: fixer\n    description: Fixes.\n    access: edit\n    instructions: roles/fixer.md\noverrides:\n  - default: delegation\n    stage: implement\n    role: fixer\n", "roles/fixer.md")
		write(t, filepath.Join(f.root, "extensions.yaml"), "extensions:\n  - id: one\n    source: {path: extensions/one}\n    overrides:\n      delegation: {}\n  - id: two\n    source: {path: extensions/two}\n    overrides:\n      delegation: {}\n")
		out := f.fail("check")
		if !strings.Contains(out, "extensions one and two both allow delegation of the implement stage, to worker and to x-two-fixer") || !strings.Contains(out, "needs a person: choose which role") {
			t.Fatalf("check: %s", out)
		}
	})
}
