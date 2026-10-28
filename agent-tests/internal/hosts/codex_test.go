package hosts

import (
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestCodexParseRecorded(t *testing.T) {
	h, err := ForName("codex")
	if err != nil {
		t.Fatal(err)
	}
	p := parseFile(t, h, "codex-exec.jsonl")
	if !p.Terminal {
		t.Fatal("expected a terminal turn")
	}
	if p.SessionID != "01a1670e-b8f6-7671-a852-b960bc093e4b" {
		t.Fatalf("session %s", p.SessionID)
	}
	if p.FinalReply != "OK" {
		t.Fatalf("reply %q", p.FinalReply)
	}
	if p.Turns != 1 || p.ToolCalls != 0 {
		t.Fatalf("turns %d tools %d", p.Turns, p.ToolCalls)
	}
	if p.Model != "" || p.NullReasons["model"] != "not in the stream; pinned model is used" {
		t.Fatalf("model %q reason %q", p.Model, p.NullReasons["model"])
	}
	eqInt(t, "input_total", p.Tokens.InputTotal, 15974)
	eqInt(t, "cache_read", p.Tokens.CacheRead, 0)
	eqInt(t, "cache_write", p.Tokens.CacheWrite, 0)
	eqInt(t, "input_uncached", p.Tokens.InputUncached, 15974)
	eqInt(t, "output", p.Tokens.Output, 5)
	eqInt(t, "reasoning", p.Tokens.Reasoning, 0)
	if p.UsageSource != "reported" || p.CostUSD != nil || p.CostBasis != "" {
		t.Fatalf("usage %s cost %v basis %q", p.UsageSource, p.CostUSD, p.CostBasis)
	}
	if len(p.Warnings) != 2 || !strings.Contains(p.Warnings[0], "unrecognized configuration") {
		t.Fatalf("warnings %q", p.Warnings)
	}
}

func TestCodexParseSyntheticTwoTurns(t *testing.T) {
	h, err := ForName("codex")
	if err != nil {
		t.Fatal(err)
	}
	const stream = `
Reading additional input from stdin...
{"type":"thread.started","thread_id":"thread-synth"}
{"type":"item.completed","item":{"id":"1","type":"agent_message","text":"first"}}
{"type":"turn.completed","usage":{"input_tokens":100,"cached_input_tokens":40,"output_tokens":10,"reasoning_output_tokens":3}}
{"type":"item.completed","item":{"id":"2","type":"command_execution","command":"ls -l"}}
{"type":"item.completed","item":{"id":"3","type":"agent_message","text":"second"}}
{"type":"turn.completed","usage":{"input_tokens":50,"cached_input_tokens":10,"output_tokens":5,"reasoning_output_tokens":1,"cache_write_input_tokens":7}}
not-json-either
`
	p := parseText(t, h, stream)
	if !p.Terminal || p.Turns != 2 || p.FinalReply != "second" || p.SessionID != "thread-synth" {
		t.Fatalf("terminal %v turns %d reply %q session %s", p.Terminal, p.Turns, p.FinalReply, p.SessionID)
	}
	eqInt(t, "input_total", p.Tokens.InputTotal, 150)
	eqInt(t, "cache_read", p.Tokens.CacheRead, 50)
	eqInt(t, "cache_write", p.Tokens.CacheWrite, 7)
	eqInt(t, "input_uncached", p.Tokens.InputUncached, 93)
	eqInt(t, "output", p.Tokens.Output, 15)
	eqInt(t, "reasoning", p.Tokens.Reasoning, 4)
	if p.ToolCalls != 1 || !reflect.DeepEqual(p.Reads, []string{"Tool: command_execution ls -l"}) {
		t.Fatalf("tools %d reads %q", p.ToolCalls, p.Reads)
	}
	if p.CostUSD != nil {
		t.Fatal("cost should be unset")
	}
}

func TestCodexParseSyntheticTruncated(t *testing.T) {
	h, err := ForName("codex")
	if err != nil {
		t.Fatal(err)
	}
	const stream = `
Reading additional input from stdin...
{"type":"thread.started","thread_id":"thread-trunc"}
`
	p := parseText(t, h, stream)
	if p.Terminal {
		t.Fatal("terminal on a truncated stream")
	}
	if p.SessionID != "thread-trunc" || p.Turns != 0 {
		t.Fatalf("%+v", p)
	}
}

func TestCodexBuildArgvSecretNotInEnv(t *testing.T) {
	h, err := ForName("codex")
	if err != nil {
		t.Fatal(err)
	}
	const secret = "codex-secret-value"
	run := "/work/run"
	in := LaunchInput{
		Prompt:    "Reply OK",
		Model:     "gpt-6-sol",
		Effort:    "medium",
		Workspace: "/work/ws",
		RunRoot:   run,
		DecoyDir:  "/work/run/decoy",
		Cfg: HostConfig{
			Command: "codex",
			Model:   "gpt-6-sol",
			AuthEnv: []string{"OPENAI_API_KEY"},
		},
		Caller: map[string]string{
			"PATH":           "/usr/bin:/bin",
			"HOME":           "/home/caller",
			"OPENAI_API_KEY": secret,
		},
	}
	launch, err := h.Build(in)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		"codex",
		"exec",
		"--json",
		"--ignore-user-config",
		"--skip-git-repo-check",
		"-C", "/work/ws",
		"-m", "gpt-6-sol",
		"--sandbox", "workspace-write",
		"--add-dir", filepath.Join(run, "gocache"),
		"--add-dir", filepath.Join("/work/ws", ".codex"),
		"-c", `model_reasoning_effort="medium"`,
		"Reply OK",
	}
	if !reflect.DeepEqual(launch.Argv, want) {
		t.Fatalf("argv %#v", launch.Argv)
	}
	if launch.Stdin != "" {
		t.Fatalf("stdin %q", launch.Stdin)
	}
	for _, a := range launch.Argv {
		if strings.Contains(a, secret) {
			t.Fatalf("argv leaked secret: %s", a)
		}
	}
	for _, e := range launch.Env {
		if strings.Contains(e, secret) || strings.Contains(e, "OPENAI_API_KEY") || strings.Contains(e, "/home/caller") {
			t.Fatalf("env leaked: %s", e)
		}
	}
	if !reflect.DeepEqual(launch.SetupArgv, [][]string{{"codex", "login", "--with-api-key"}}) {
		t.Fatalf("setup argv %#v", launch.SetupArgv)
	}
	if !reflect.DeepEqual(launch.SetupStdin, []string{secret}) {
		t.Fatalf("setup stdin %#v", launch.SetupStdin)
	}
	env := envMap(t, launch.Env)
	if env["HOME"] != filepath.Join(run, "codex-home") {
		t.Fatalf("HOME %s", env["HOME"])
	}
	if env["CODEX_HOME"] != filepath.Join(run, "codex") {
		t.Fatalf("CODEX_HOME %s", env["CODEX_HOME"])
	}
	wantDirs := []string{filepath.Join(run, "codex"), filepath.Join(run, "codex-home"), filepath.Join("/work/ws", ".codex")}
	if !reflect.DeepEqual(launch.MkDirs, wantDirs) {
		t.Fatalf("mkdirs %v", launch.MkDirs)
	}

	in.Caller = map[string]string{"PATH": "/usr/bin:/bin"}
	in.Effort = ""
	plain, err := h.Build(in)
	if err != nil {
		t.Fatal(err)
	}
	if plain.SetupArgv != nil || plain.SetupStdin != nil {
		t.Fatalf("unexpected setup %#v %#v", plain.SetupArgv, plain.SetupStdin)
	}
	if len(launch.CopyIn) != 0 || len(plain.CopyIn) != 0 {
		t.Fatalf("copy without subscription %#v %#v", launch.CopyIn, plain.CopyIn)
	}
	if containsArg(plain.Argv, "-c") {
		t.Fatalf("effort flag without effort: %q", plain.Argv)
	}
}

func TestCodexAuthReady(t *testing.T) {
	h, err := ForName("codex")
	if err != nil {
		t.Fatal(err)
	}
	cfg := HostConfig{AuthEnv: []string{"OPENAI_API_KEY"}}
	if ok, _ := h.AuthReady(map[string]string{"OPENAI_API_KEY": "x"}, cfg); !ok {
		t.Fatal("expected ready")
	}
	if ok, _ := h.AuthReady(map[string]string{}, cfg); ok {
		t.Fatal("expected missing")
	}
}

// A run allowed to commit and push names the workspace's .git and its remote
// as writable roots after the fixed ones; without them workspace-write keeps
// .git read-only and the remote outside the sandbox.
func TestCodexBuildNamesWritableDirs(t *testing.T) {
	h, err := ForName("codex")
	if err != nil {
		t.Fatal(err)
	}
	in := LaunchInput{
		Prompt:       "Reply OK",
		Workspace:    "/work/ws",
		RunRoot:      "/work/run",
		Cfg:          HostConfig{Command: "codex", Model: "gpt-6-sol"},
		Caller:       map[string]string{"PATH": "/usr/bin:/bin"},
		WritableDirs: []string{"/work/ws/.git", "/work/run/remotes/ws.git"},
	}
	launch, err := h.Build(in)
	if err != nil {
		t.Fatal(err)
	}
	argv := strings.Join(launch.Argv, " ")
	for _, want := range []string{
		"--add-dir " + filepath.Join("/work/ws", ".codex") + " --add-dir /work/ws/.git --add-dir /work/run/remotes/ws.git",
	} {
		if !strings.Contains(argv, want) {
			t.Fatalf("argv %q lacks %q", argv, want)
		}
	}
}
