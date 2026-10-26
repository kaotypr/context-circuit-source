package hosts

import (
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestCursorParseRecordedStream(t *testing.T) {
	h, err := ForName("cursor")
	if err != nil {
		t.Fatal(err)
	}
	p := parseFile(t, h, "cursor-stream.jsonl")
	if !p.Terminal || p.IsError {
		t.Fatalf("terminal %v error %v", p.Terminal, p.IsError)
	}
	if p.SessionID != "adc3832e-c388-4dbd-9291-d87c608760c8" {
		t.Fatalf("session %s", p.SessionID)
	}
	if p.Model != "Claude Opus 4.8 300K High" {
		t.Fatalf("model %q", p.Model)
	}
	if !strings.HasPrefix(p.FinalReply, "OK\n") || !strings.Contains(p.FinalReply, "example-marketplace") {
		t.Fatalf("reply %q", p.FinalReply)
	}
	eqInt(t, "input_uncached", p.Tokens.InputUncached, 2)
	eqInt(t, "output", p.Tokens.Output, 73)
	eqInt(t, "cache_read", p.Tokens.CacheRead, 0)
	eqInt(t, "cache_write", p.Tokens.CacheWrite, 21364)
	eqInt(t, "input_total", p.Tokens.InputTotal, 2+21364)
	if p.ReadsAvailable || p.NullReasons["reads"] != "reads list not available" {
		t.Fatalf("reads %v reason %q", p.ReadsAvailable, p.NullReasons["reads"])
	}
	if p.CostUSD != nil || p.UsageSource != "reported" || p.ToolCalls != 0 {
		t.Fatalf("cost %v source %s tools %d", p.CostUSD, p.UsageSource, p.ToolCalls)
	}
}

func TestCursorParseRecordedResult(t *testing.T) {
	h, err := ForName("cursor")
	if err != nil {
		t.Fatal(err)
	}
	p := parseFile(t, h, "cursor-result.json")
	if !p.Terminal {
		t.Fatal("expected a terminal result")
	}
	if p.SessionID != "4bec7943-97d7-43b4-b478-a2defbe814bb" {
		t.Fatalf("session %s", p.SessionID)
	}
	if p.Model != "" {
		t.Fatalf("model %q", p.Model)
	}
	if !strings.HasPrefix(p.FinalReply, "OK\n") {
		t.Fatalf("reply %q", p.FinalReply)
	}
	eqInt(t, "input_uncached", p.Tokens.InputUncached, 2)
	eqInt(t, "output", p.Tokens.Output, 75)
	eqInt(t, "cache_write", p.Tokens.CacheWrite, 21364)
	eqInt(t, "input_total", p.Tokens.InputTotal, 2+21364)
}

func TestCursorParseSyntheticToolCall(t *testing.T) {
	h, err := ForName("cursor")
	if err != nil {
		t.Fatal(err)
	}
	const stream = `
hello
{"type":"system","subtype":"init","session_id":"cur-1","model":"Test Model"}
{"type":"tool_call","name":"read","input":{"file_path":"main.go"}}
{"type":"result","is_error":false,"result":"OK","session_id":"cur-1","usage":{"inputTokens":1,"outputTokens":2,"cacheReadTokens":3,"cacheWriteTokens":4}}
`
	p := parseText(t, h, stream)
	if !p.Terminal || p.FinalReply != "OK" || p.Model != "Test Model" || p.SessionID != "cur-1" {
		t.Fatalf("%+v", p)
	}
	if !p.ReadsAvailable || p.ToolCalls != 1 {
		t.Fatalf("reads %v tools %d", p.ReadsAvailable, p.ToolCalls)
	}
	if _, ok := p.NullReasons["reads"]; ok {
		t.Fatal("reads marked unavailable")
	}
	if !reflect.DeepEqual(p.Reads, []string{"Tool: read main.go"}) {
		t.Fatalf("reads %q", p.Reads)
	}
	eqInt(t, "input_total", p.Tokens.InputTotal, 8)
	eqInt(t, "input_uncached", p.Tokens.InputUncached, 1)
	eqInt(t, "cache_read", p.Tokens.CacheRead, 3)
	eqInt(t, "cache_write", p.Tokens.CacheWrite, 4)
	eqInt(t, "output", p.Tokens.Output, 2)
	if p.CostUSD != nil {
		t.Fatal("cost")
	}
}

func TestCursorParseSyntheticTruncated(t *testing.T) {
	h, err := ForName("cursor")
	if err != nil {
		t.Fatal(err)
	}
	const stream = `
garbage
{"type":"system","subtype":"init","session_id":"cur-2","model":"Test Model"}
`
	p := parseText(t, h, stream)
	if p.Terminal {
		t.Fatal("terminal on a truncated stream")
	}
	if p.ReadsAvailable || p.NullReasons["reads"] != "reads list not available" {
		t.Fatalf("reads %v reason %q", p.ReadsAvailable, p.NullReasons["reads"])
	}
	if p.SessionID != "cur-2" {
		t.Fatalf("session %s", p.SessionID)
	}
}

func TestCursorBuildArgv(t *testing.T) {
	h, err := ForName("cursor")
	if err != nil {
		t.Fatal(err)
	}
	run := "/work/run"
	in := LaunchInput{
		Prompt:    "Reply OK",
		Model:     "auto",
		Effort:    "high",
		Workspace: "/work/ws",
		RunRoot:   run,
		DecoyDir:  "/work/run/decoy",
		Cfg: HostConfig{
			Command: "cursor-agent",
			Model:   "auto",
			AuthEnv: []string{"CURSOR_API_KEY"},
		},
		Caller: map[string]string{
			"PATH":            "/usr/bin:/bin",
			"HOME":            "/home/caller",
			"XDG_CONFIG_HOME": "/home/caller/.config",
			"CURSOR_API_KEY":  "cursor-auth-value",
		},
	}
	launch, err := h.Build(in)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		"cursor-agent",
		"-p", "Reply OK",
		"--model", "auto",
		"--output-format", "stream-json",
		"--force",
		"--trust",
		"--workspace", "/work/ws",
	}
	if !reflect.DeepEqual(launch.Argv, want) {
		t.Fatalf("argv %#v", launch.Argv)
	}
	if containsArg(launch.Argv, "--effort") || containsArg(launch.Argv, "high") {
		t.Fatalf("effort leaked into argv %#v", launch.Argv)
	}
	home := filepath.Join(run, "cursor-home")
	env := envMap(t, launch.Env)
	if env["HOME"] != home {
		t.Fatalf("HOME %s", env["HOME"])
	}
	if env["XDG_CONFIG_HOME"] != filepath.Join(home, ".config") {
		t.Fatalf("xdg %s", env["XDG_CONFIG_HOME"])
	}
	if env["CURSOR_API_KEY"] != "cursor-auth-value" {
		t.Fatalf("auth %s", env["CURSOR_API_KEY"])
	}
	for _, e := range launch.Env {
		if strings.Contains(e, "/home/caller") {
			t.Fatalf("caller home leaked: %s", e)
		}
	}
	wantDirs := []string{home, filepath.Join(home, ".config")}
	if !reflect.DeepEqual(launch.MkDirs, wantDirs) || launch.Dir != "/work/ws" {
		t.Fatalf("dir %s mkdirs %v", launch.Dir, launch.MkDirs)
	}
}

func TestCursorAuthReady(t *testing.T) {
	h, err := ForName("cursor")
	if err != nil {
		t.Fatal(err)
	}
	cfg := HostConfig{AuthEnv: []string{"CURSOR_API_KEY"}}
	if ok, _ := h.AuthReady(map[string]string{"CURSOR_API_KEY": "x"}, cfg); !ok {
		t.Fatal("expected ready")
	}
	if ok, _ := h.AuthReady(map[string]string{}, cfg); ok {
		t.Fatal("expected missing")
	}
}
