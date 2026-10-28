package hosts

import (
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestClaudeParseRecorded(t *testing.T) {
	h, err := ForName("claude")
	if err != nil {
		t.Fatal(err)
	}
	p := parseFile(t, h, "claude-result.jsonl")
	if !p.Terminal {
		t.Fatal("expected a terminal result")
	}
	if p.IsError {
		t.Fatal("is_error")
	}
	if p.SessionID != "ee304eef-1733-4b5d-a940-7faa6dc4f53b" {
		t.Fatalf("session %s", p.SessionID)
	}
	if p.FinalReply != "OK" {
		t.Fatalf("reply %q", p.FinalReply)
	}
	if p.Turns != 1 {
		t.Fatalf("turns %d", p.Turns)
	}
	if p.Model != "claude-opus-5-5" || len(p.ModelsUsed) != 1 || p.ModelsUsed[0] != "claude-opus-5-5" {
		t.Fatalf("model %q used %v", p.Model, p.ModelsUsed)
	}
	if p.Subagents != 0 || p.ToolCalls != 0 || !p.ReadsAvailable {
		t.Fatalf("subagents %d tools %d reads %v", p.Subagents, p.ToolCalls, p.ReadsAvailable)
	}
	// input 2 + cache_creation 31194 + cache_read 11304
	eqInt(t, "input_total", p.Tokens.InputTotal, 2+31194+11304)
	eqInt(t, "input_uncached", p.Tokens.InputUncached, 2)
	eqInt(t, "cache_write", p.Tokens.CacheWrite, 31194)
	eqInt(t, "cache_read", p.Tokens.CacheRead, 11304)
	eqInt(t, "output", p.Tokens.Output, 4)
	if p.Tokens.Reasoning != nil {
		t.Fatalf("reasoning %d", *p.Tokens.Reasoning)
	}
	if p.NullReasons["reasoning"] != "not reported" {
		t.Fatalf("reasoning reason %q", p.NullReasons["reasoning"])
	}
	if p.UsageSource != "reported" || p.CostBasis != "reported" || p.CostUSD == nil {
		t.Fatalf("usage %s basis %s cost %v", p.UsageSource, p.CostBasis, p.CostUSD)
	}
	near(t, *p.CostUSD, 0.2519008)
}

func TestClaudeParseSyntheticTwoModels(t *testing.T) {
	h, err := ForName("claude")
	if err != nil {
		t.Fatal(err)
	}
	const stream = `
this is not json
{"type":"assistant","message":{"role":"assistant","content":[{"type":"text","text":"hi"},{"type":"tool_use","name":"Read","input":{"file_path":"a.go"}},{"type":"tool_use","name":"Bash","input":{"command":"echo hi\nthere"}}]}}
{"type":"nope"}
{"type":"result","session_id":"sess-synth","is_error":false,"result":"done","num_turns":2,"total_cost_usd":1.5,"usage":{"input_tokens":1,"cache_creation_input_tokens":2,"cache_read_input_tokens":3,"output_tokens":4,"thinking_tokens":9},"modelUsage":{"zeta":{"inputTokens":10,"outputTokens":1,"cacheReadInputTokens":2,"cacheCreationInputTokens":3},"alpha":{"inputTokens":4,"outputTokens":5,"cacheReadInputTokens":6,"cacheCreationInputTokens":7}}}
`
	p := parseText(t, h, stream)
	if !p.Terminal || p.FinalReply != "done" || p.SessionID != "sess-synth" || p.Turns != 2 {
		t.Fatalf("%+v", p)
	}
	if p.Subagents != 1 {
		t.Fatalf("subagents %d", p.Subagents)
	}
	wantModels := []string{"alpha", "zeta"}
	if !reflect.DeepEqual(p.ModelsUsed, wantModels) || p.Model != "alpha" {
		t.Fatalf("model %q used %v", p.Model, p.ModelsUsed)
	}
	eqInt(t, "input_uncached", p.Tokens.InputUncached, 14)
	eqInt(t, "cache_read", p.Tokens.CacheRead, 8)
	eqInt(t, "cache_write", p.Tokens.CacheWrite, 10)
	eqInt(t, "output", p.Tokens.Output, 6)
	eqInt(t, "input_total", p.Tokens.InputTotal, 32)
	eqInt(t, "reasoning", p.Tokens.Reasoning, 9)
	if p.ToolCalls != 2 || !p.ReadsAvailable {
		t.Fatalf("tools %d reads %v", p.ToolCalls, p.Reads)
	}
	wantReads := []string{"Tool: Read a.go", "Tool: Bash echo hi there"}
	if !reflect.DeepEqual(p.Reads, wantReads) {
		t.Fatalf("reads %q", p.Reads)
	}
	if p.CostBasis != "reported" || p.CostUSD == nil {
		t.Fatal("cost")
	}
	near(t, *p.CostUSD, 1.5)
}

func TestClaudeParseSyntheticTruncated(t *testing.T) {
	h, err := ForName("claude")
	if err != nil {
		t.Fatal(err)
	}
	const stream = `
not json
{"type":"assistant","message":{"content":[{"type":"tool_use","name":"Read","input":{"file_path":"a.go"}}]}}
`
	p := parseText(t, h, stream)
	if p.Terminal {
		t.Fatal("terminal on a truncated stream")
	}
	if p.ToolCalls != 1 {
		t.Fatalf("tools %d", p.ToolCalls)
	}
}

func TestClaudeBuildArgv(t *testing.T) {
	h, err := ForName("claude")
	if err != nil {
		t.Fatal(err)
	}
	run := "/work/run"
	in := LaunchInput{
		Prompt:    "Reply OK",
		Model:     "claude-opus-5-5",
		Workspace: "/work/ws",
		RunRoot:   run,
		SessionID: "11111111-1111-1111-1111-111111111111",
		DecoyDir:  "/work/run/decoy",
		Cfg: HostConfig{
			Command:      "claude",
			Model:        "claude-opus-5-5",
			MaxBudgetUSD: 3,
			AuthEnv:      []string{"ANTHROPIC_API_KEY", "CLAUDE_CODE_OAUTH_TOKEN"},
		},
		Caller: map[string]string{
			"PATH":              "/usr/bin:/bin",
			"HOME":              "/home/caller",
			"ANTHROPIC_API_KEY": "claude-auth-value",
			"GITHUB_TOKEN":      "gh-should-not-pass",
		},
	}
	launch, err := h.Build(in)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		"claude",
		"-p", "Reply OK",
		"--model", "claude-opus-5-5",
		"--output-format", "stream-json",
		"--verbose",
		"--dangerously-skip-permissions",
		"--strict-mcp-config",
		"--max-budget-usd", "3",
		"--session-id", "11111111-1111-1111-1111-111111111111",
	}
	if !reflect.DeepEqual(launch.Argv, want) {
		t.Fatalf("argv %q", launch.Argv)
	}
	if launch.Dir != "/work/ws" {
		t.Fatalf("dir %s", launch.Dir)
	}
	home := filepath.Join(run, "claude-home")
	homes := 0
	for _, e := range launch.Env {
		if strings.HasPrefix(e, "HOME=") {
			homes++
			if e != "HOME="+home {
				t.Fatalf("HOME %s", e)
			}
		}
		if strings.Contains(e, "/home/caller") || strings.Contains(e, "GITHUB_TOKEN") || strings.Contains(e, "gh-should-not-pass") {
			t.Fatalf("caller leaked: %s", e)
		}
	}
	if homes != 1 {
		t.Fatalf("HOME count %d env %v", homes, launch.Env)
	}
	env := envMap(t, launch.Env)
	if env["CLAUDE_CONFIG_DIR"] != filepath.Join(run, "claude") {
		t.Fatalf("config dir %s", env["CLAUDE_CONFIG_DIR"])
	}
	if env["ANTHROPIC_API_KEY"] != "claude-auth-value" {
		t.Fatalf("auth %v", env["ANTHROPIC_API_KEY"])
	}
	wantDirs := []string{filepath.Join(run, "claude"), home}
	if !reflect.DeepEqual(launch.MkDirs, wantDirs) {
		t.Fatalf("mkdirs %v", launch.MkDirs)
	}

	in.Effort = "high"
	withEffort, err := h.Build(in)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(withEffort.Argv[len(withEffort.Argv)-2:], []string{"--effort", "high"}) {
		t.Fatalf("effort argv %q", withEffort.Argv)
	}

	in.Effort = ""
	in.Cfg.MaxBudgetUSD = 0
	in.Cfg.Command = ""
	zeroBudget, err := h.Build(in)
	if err != nil {
		t.Fatal(err)
	}
	if zeroBudget.Argv[0] != "claude" || !containsArg(zeroBudget.Argv, "3") {
		t.Fatalf("default budget argv %q", zeroBudget.Argv)
	}
}

func TestClaudeAuthReady(t *testing.T) {
	h, err := ForName("claude")
	if err != nil {
		t.Fatal(err)
	}
	cfg := HostConfig{AuthEnv: []string{"ANTHROPIC_API_KEY", "CLAUDE_CODE_OAUTH_TOKEN"}}
	if ok, _ := h.AuthReady(map[string]string{"ANTHROPIC_API_KEY": "x"}, cfg); !ok {
		t.Fatal("api key")
	}
	if ok, _ := h.AuthReady(map[string]string{"CLAUDE_CODE_OAUTH_TOKEN": "x"}, cfg); !ok {
		t.Fatal("oauth")
	}
	if ok, reason := h.AuthReady(map[string]string{}, cfg); ok {
		t.Fatal(reason)
	}
}

func containsArg(argv []string, want string) bool {
	for _, a := range argv {
		if a == want {
			return true
		}
	}
	return false
}

func envMap(t *testing.T, env []string) map[string]string {
	t.Helper()
	out := make(map[string]string, len(env))
	for _, e := range env {
		k, v, ok := strings.Cut(e, "=")
		if !ok {
			t.Fatalf("bad env %q", e)
		}
		out[k] = v
	}
	return out
}

func TestClaudeModelUsageMissingCountStaysNull(t *testing.T) {
	stream := `{"type":"result","session_id":"s","is_error":false,"result":"OK","usage":{"input_tokens":5},"modelUsage":{"m1":{"inputTokens":5,"outputTokens":2}}}` + "\n"
	p, err := (claudeHost{}).Parse(strings.NewReader(stream))
	if err != nil {
		t.Fatal(err)
	}
	if p.Tokens.CacheRead != nil || p.Tokens.CacheWrite != nil || p.Tokens.InputTotal != nil {
		t.Fatalf("a count the host did not report became a number: %+v", p.Tokens)
	}
	if p.Tokens.Output == nil || *p.Tokens.Output != 2 {
		t.Fatalf("output = %v", p.Tokens.Output)
	}
}
