package scenario

import (
	"encoding/json"
	"testing"

	"github.com/kaotypr/context-circuit-source/agent-tests/internal/hosts"
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
