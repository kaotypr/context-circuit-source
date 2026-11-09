package scenario

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/kaotypr/context-circuit-source/agent-tests/internal/hosts"
	"github.com/kaotypr/context-circuit-source/agent-tests/internal/world"
)

func streamEvent(text string) string {
	b, _ := json.Marshal(map[string]any{"type": "user", "message": map[string]any{"content": text}})
	return string(b) + "\n"
}

// Only a dispatch for the fixture plan counts, a shell continuation is read
// through, and a doc's examples are not mistaken for the agent's calls.
func TestDispatchSteps(t *testing.T) {
	doc := "sh .context-circuit/cli --json agent dispatch --host HOST --role worker \\\n  --plan p0001 --task 'FINDINGS' --override repair-loop --round 1 --path WORKTREE\n"
	verify := "sh .context-circuit/cli --json agent dispatch --host claude-code \\\n  --role x-csv-check-csv-verifier --plan p1000 --override repair-loop --round 1 --path wt"
	fix := `sh .context-circuit/cli --json agent dispatch --host claude-code --role worker --plan p1000 --task "Findings: no \"header\" row." --override repair-loop --round 1 --path wt`
	stream := streamEvent(doc) + streamEvent(verify) + streamEvent(verify) + streamEvent(fix) +
		streamEvent("sh .context-circuit/cli --json agent dispatch --host claude-code --role x-csv-check-csv-verifier --plan p1000 --override repair-loop --round 2 --path wt")
	if got := loopSteps(stream); got != "VWV" {
		t.Fatalf("steps %q", got)
	}
	ctx := CheckContext{Stream: stream, Parsed: hosts.Parsed{FinalReply: "The loop was clean after 2 rounds."}}
	if r := checkLoopSequence(ctx); !r.Pass {
		t.Fatal(r.Detail)
	}
	if r := checkLoopClean(ctx); !r.Pass {
		t.Fatal(r.Detail)
	}
	if r := checkNotAllowedHeld(ctx); r.Pass {
		t.Fatal("a verifier dispatch passed the not-allowed check")
	}
	if r := checkLimitHeld(ctx); !r.Pass {
		t.Fatal(r.Detail)
	}
	if r := checkNotAllowedHeld(CheckContext{Stream: streamEvent(doc)}); !r.Pass {
		t.Fatal("a doc example counted as a dispatch")
	}
}

// The agent implements the plan itself, so a loop clean in round 1 passes and
// reports one round; a fix left unverified, or a round past the fixture's
// three, does not.
func TestLoopShape(t *testing.T) {
	verify := func(round string) string {
		return streamEvent("sh .context-circuit/cli --json agent dispatch --host claude-code --role x-csv-check-csv-verifier --plan p1000 --override repair-loop --round " + round + " --path wt")
	}
	fix := func(round string) string {
		return streamEvent("sh .context-circuit/cli --json agent dispatch --host claude-code --role worker --plan p1000 --task 'Findings' --override repair-loop --round " + round + " --path wt")
	}
	cases := []struct {
		stream, reply string
		sequence      bool
		clean         bool
	}{
		{verify("1"), "The csv-check repair loop was clean after 1 round.", true, true},
		{verify("1"), "Round 1 was clean.", true, true},
		{verify("1"), "The loop was clean after 2 rounds.", true, false},
		{verify("1") + fix("1") + verify("2"), "Round 1 was not clean. Round 2 was clean.", true, true},
		{verify("1") + fix("1"), "Fixed.", false, false},
		{fix("1") + verify("1"), "Clean after 1 round.", false, true},
		{verify("1") + fix("1") + verify("2") + fix("2") + verify("3") + fix("3") + verify("4"), "Clean after 4 rounds.", false, false},
		{"", "Clean.", false, false},
	}
	for i, c := range cases {
		ctx := CheckContext{Stream: c.stream, Parsed: hosts.Parsed{FinalReply: c.reply}}
		if got := checkLoopSequence(ctx).Pass; got != c.sequence {
			t.Errorf("case %d sequence %v, want %v", i, got, c.sequence)
		}
		if got := checkLoopClean(ctx).Pass; got != c.clean {
			t.Errorf("case %d clean %v, want %v", i, got, c.clean)
		}
	}
}

// A role listed among the agents a host installed is not a spawn; a tool call
// naming it, or a dispatch for any plan, is.
func TestNotAllowedSpawn(t *testing.T) {
	init := `{"type":"system","subtype":"init","agents":["cc-worker","cc-x-csv-check-csv-verifier"]}` + "\n"
	if r := checkNotAllowedHeld(CheckContext{Stream: init}); !r.Pass {
		t.Fatalf("installed agents counted as a spawn: %s", r.Detail)
	}
	spawn := `{"type":"assistant","message":{"content":[{"type":"tool_use","name":"Agent","input":{"subagent_type":"cc-x-csv-check-csv-verifier","prompt":"check"}}]}}` + "\n"
	if r := checkNotAllowedHeld(CheckContext{Stream: init + spawn}); r.Pass {
		t.Fatal("a spawn of the verifier passed")
	}
	dispatch := streamEvent("sh .context-circuit/cli --json agent dispatch --host codex --role x-csv-check-csv-verifier --path wt")
	if r := checkNotAllowedHeld(CheckContext{Stream: dispatch}); r.Pass {
		t.Fatal("a dispatch of the verifier without a plan passed")
	}
}

// The agent allowing the override itself leaves an allowance in
// extensions.yaml.
func TestAllowanceHeld(t *testing.T) {
	ws := t.TempDir()
	path := filepath.Join(ws, "extensions.yaml")
	ctx := CheckContext{World: world.World{WS: ws}}
	if err := os.WriteFile(path, []byte("extensions:\n  - id: csv-check\n    path: extensions/csv-check\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if r := checkAllowanceHeld(ctx); !r.Pass {
		t.Fatal(r.Detail)
	}
	if err := os.WriteFile(path, []byte("extensions:\n  - id: csv-check\n    overrides:\n      repair-loop: {}\n    path: extensions/csv-check\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if r := checkAllowanceHeld(ctx); r.Pass {
		t.Fatal("an allowance passed")
	}
}
