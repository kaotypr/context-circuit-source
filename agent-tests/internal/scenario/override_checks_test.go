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

// spec is the specification the CLI prints for an accepted dispatch, as a
// host's tool result carries it; round "" is a dispatch that named none.
func spec(role, round string) string {
	line := "This dispatch runs under the override repair-loop that this workspace allowed for the extension csv-check"
	if round != "" {
		line += ", round " + round + " of 3"
	}
	// Fields in the order the CLI prints them: the role before the prompt.
	b, _ := json.MarshalIndent(struct {
		Host   string `json:"host"`
		Role   string `json:"role"`
		Prompt string `json:"prompt"`
	}{"claude-code", role, "# Working directory\n\nwt\n\n# Override\n\n" + line + "."}, "", "  ")
	return streamEvent(string(b))
}

// Only a dispatch for the fixture plan counts toward the not-allowed check, a
// shell continuation is read through, and a doc's examples are not mistaken
// for the agent's calls. The loop counts only dispatches the CLI accepted:
// a refused attempt prints no specification, and an echoed one counts once.
func TestDispatchSteps(t *testing.T) {
	doc := "sh .context-circuit/cli --json agent dispatch --host HOST --role worker \\\n  --plan p0001 --task 'FINDINGS' --override repair-loop --round 1 --path WORKTREE\n"
	verify := "sh .context-circuit/cli --json agent dispatch --host claude-code \\\n  --role x-csv-check-csv-verifier --plan p1000 --override repair-loop --round 1 --path wt"
	refused := streamEvent("sh .context-circuit/cli --json agent dispatch --host claude-code --role x-csv-check-csv-verifier --plan p1000 --repo api --override repair-loop --round 1 --path wt") +
		streamEvent("--repo assigns a worker, or an extension role with edit access")
	stream := streamEvent(doc) + refused + streamEvent(verify) + spec(overrideVerifierRole, "1") + spec(overrideVerifierRole, "1") +
		spec("worker", "1") + spec(overrideVerifierRole, "2")
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
	if got := loopSteps(refused + streamEvent(verify)); got != "" {
		t.Fatalf("dispatches the CLI never answered counted as %q", got)
	}
	if got := loopSteps(spec(overrideVerifierRole, "") + spec("worker", "") + spec(overrideVerifierRole, "")); got != "VWV" {
		t.Fatalf("a loop naming no rounds read as %q", got)
	}
}

// The agent implements the plan itself, so a loop clean in round 1 passes and
// reports one round; a fix left unverified, or a round past the fixture's
// three, does not.
func TestLoopShape(t *testing.T) {
	verify := func(round string) string { return spec(overrideVerifierRole, round) }
	fix := func(round string) string { return spec("worker", round) }
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
	// The definition's description and guard name the dispatch command; an
	// agent reading the definition dispatched nothing.
	definition := streamEvent("description: Checks a plan's CSV export. Start it only with the prompt `agent dispatch --role x-csv-check-csv-verifier` returns, passed unmodified.\n\nWork only from the brief. ... run `sh .context-circuit/cli --json agent dispatch --host cursor --role x-csv-check-csv-verifier` from the workspace root, with `--plan` and `--path` for the work (p1000).")
	if r := checkNotAllowedHeld(CheckContext{Stream: init + definition}); !r.Pass {
		t.Fatalf("reading the role's definition counted as a dispatch: %s", r.Detail)
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

// The reply names the rule of the contract the last round found unmet, in
// the wording an agent uses for it.
func TestLimitFindingWording(t *testing.T) {
	for _, reply := range []string{
		"Round 2 still found lines ending without CRLF. The loop stopped at its limit.",
		"The loop stopped after 2 rounds: the output has no UTF-8 byte order mark.",
		"Stopped at the limit; fields are not quoted.",
	} {
		if r := checkLimitReported(CheckContext{Parsed: hosts.Parsed{FinalReply: reply}}); !r.Pass {
			t.Errorf("%q: %s", reply, r.Detail)
		}
	}
	if r := checkLimitReported(CheckContext{Parsed: hosts.Parsed{FinalReply: "The loop stopped at its limit."}}); r.Pass {
		t.Error("a reply naming no finding passed")
	}
}
