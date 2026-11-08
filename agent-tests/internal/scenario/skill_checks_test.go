package scenario

import (
	"testing"

	"github.com/kaotypr/context-circuit-source/agent-tests/internal/hosts"
)

// The named stage's points count with their flags in either order, and only
// before followed by after passes.
func TestSkillHooksRun(t *testing.T) {
	before := streamEvent("sh .context-circuit/cli --json extension hooks --when before --stage qa-pass/qa --plan p1000")
	after := streamEvent("sh .context-circuit/cli extension hooks --stage qa-pass/qa --plan p1000 --when after")
	core := streamEvent("sh .context-circuit/cli extension hooks --stage implement --when before --plan p1000")
	for _, tc := range []struct {
		name, stream string
		pass         bool
	}{
		{"both in order", before + after, true},
		{"after only", core + after, false},
		{"before only", before, false},
		{"reversed", after + before, false},
	} {
		if got := checkSkillHooksRun(CheckContext{Stream: tc.stream}).Pass; got != tc.pass {
			t.Errorf("%s: pass %v", tc.name, got)
		}
	}
}

func TestSkillFollowedAndUsed(t *testing.T) {
	listed := CheckContext{Parsed: hosts.Parsed{FinalReply: "Commands run:\n- go vet ./... (passed)\n- go test ./... (passed)"}}
	if r := checkSkillFollowed(listed); !r.Pass {
		t.Fatal(r.Detail)
	}
	if r := checkSkillFollowed(CheckContext{Parsed: hosts.Parsed{FinalReply: "The tests pass."}}); r.Pass {
		t.Fatal("a reply without the commands passed")
	}
	if r := checkSkillUsed(CheckContext{Parsed: hosts.Parsed{Reads: []string{".agents/skills/x-qa-pass-qa/SKILL.md"}}}); !r.Pass {
		t.Fatal(r.Detail)
	}
	if r := checkSkillUsed(CheckContext{Stream: streamEvent("go test ./...")}); r.Pass {
		t.Fatal("a stream without the skill passed")
	}
}
