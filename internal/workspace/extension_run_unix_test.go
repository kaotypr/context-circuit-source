//go:build !windows

package workspace

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// scriptFolder writes an extension folder holding one shell script.
func scriptFolder(t *testing.T, body string) string {
	t.Helper()
	folder := t.TempDir()
	if err := os.MkdirAll(filepath.Join(folder, "scripts"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(folder, "scripts", "hook.sh"), []byte(body), 0755); err != nil {
		t.Fatal(err)
	}
	return folder
}

func shellScript(onFailure string) HookScript {
	return HookScript{Run: []string{"sh", "scripts/hook.sh"}, OnFailure: onFailure}
}

func testInput() ScriptInput {
	return ScriptInput{Point: "after implement", WorkspaceRoot: "/w", Member: "maya",
		Extension: ScriptExtension{ID: "parity", Settings: map[string]any{"reference": "legacy"}},
		Record:    &ScriptRecord{ID: "p0001", Kind: "plan", Path: "/w/plans/p0001-x/plan.md", Repositories: []string{"api"}},
		Worktrees: map[string]string{"api": "/w/.worktrees/p0001/api"}}
}

// The script reads its input on stdin, runs in the directory it is given with
// no shell between, and reports through its final JSON object.
func TestAScriptReadsItsInputAndRunsInItsDirectory(t *testing.T) {
	dir := t.TempDir()
	folder := scriptFolder(t, `cat > input.json
pwd > cwd.txt
printf 'progress line\n'
printf '{"status":"ok","summary":"14 behaviors match","details_path":"parity.md"}\n'
`)
	t.Setenv("GIT_DIR", "/elsewhere")
	t.Setenv("CC_TEST_KEPT", "kept")
	if err := os.WriteFile(filepath.Join(folder, "scripts", "env.sh"), []byte("env > env.txt\n"), 0755); err != nil {
		t.Fatal(err)
	}
	run := runScript(context.Background(), folder, shellScript("stop"), dir, testInput())
	if run.Status != ScriptOK || run.Summary != "14 behaviors match" || run.DetailsPath != "parity.md" || run.OnFailure != "stop" || run.ID != "parity" {
		t.Fatalf("run: %+v", run)
	}
	var input ScriptInput
	data, err := os.ReadFile(filepath.Join(dir, "input.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, &input); err != nil || input.Point != "after implement" || input.Record.ID != "p0001" ||
		input.Extension.Settings["reference"] != "legacy" || input.Worktrees["api"] == "" || input.Confirmed {
		t.Fatalf("input: %v %s", err, data)
	}
	cwd, _ := os.ReadFile(filepath.Join(dir, "cwd.txt"))
	realDir, _ := filepath.EvalSymlinks(dir)
	if got, _ := filepath.EvalSymlinks(strings.TrimSpace(string(cwd))); got != realDir {
		t.Fatalf("ran in %s, want %s", cwd, dir)
	}
	// The member's environment, minus the Git variables the Git runner strips.
	envRun := runScript(context.Background(), folder, HookScript{Run: []string{"sh", "scripts/env.sh"}, OnFailure: "report"}, dir, testInput())
	if envRun.Status != ScriptOK {
		t.Fatalf("env run: %+v", envRun)
	}
	env, _ := os.ReadFile(filepath.Join(dir, "env.txt"))
	if !strings.Contains(string(env), "CC_TEST_KEPT=kept") || strings.Contains(string(env), "GIT_DIR=") {
		t.Fatalf("environment:\n%s", env)
	}
}

func TestAScriptsExitCodeDecidesFailure(t *testing.T) {
	dir := t.TempDir()
	run := runScript(context.Background(), scriptFolder(t, "echo 'GET /orders/1 differs' >&2\nexit 3\n"), shellScript("stop"), dir, testInput())
	if run.Status != ScriptFailed || !strings.Contains(run.Summary, "status 3") || !strings.Contains(run.Summary, "GET /orders/1 differs") {
		t.Fatalf("failure: %+v", run)
	}
	// A skip the script reports itself is carried through.
	run = runScript(context.Background(), scriptFolder(t, `echo '{"status":"skipped","summary":"the plan has no parity list"}'`), shellScript("stop"), dir, testInput())
	if run.Status != ScriptSkipped || run.Summary != "the plan has no parity list" {
		t.Fatalf("skip: %+v", run)
	}
	// A program that cannot start is a failure, not an error.
	run = runScript(context.Background(), scriptFolder(t, ""), HookScript{Run: []string{"cc-no-such-program-here"}, OnFailure: "report"}, dir, testInput())
	if run.Status != ScriptFailed || !strings.Contains(run.Summary, "not started") {
		t.Fatalf("missing program: %+v", run)
	}
}

func TestAScriptIsStoppedAtItsTimeout(t *testing.T) {
	script := shellScript("report")
	script.Timeout = "200ms"
	run := runScript(context.Background(), scriptFolder(t, "sleep 30\n"), script, t.TempDir(), testInput())
	if run.Status != ScriptFailed || !strings.Contains(run.Summary, "timed out after 200ms") {
		t.Fatalf("timeout: %+v", run)
	}
}

func TestAScriptsStdoutIsCut(t *testing.T) {
	run := runScript(context.Background(), scriptFolder(t, "head -c 200000 /dev/zero | tr '\\0' 'x'\necho\necho '{\"status\":\"ok\",\"summary\":\"lost\"}'\n"), shellScript("report"), t.TempDir(), testInput())
	if run.Status != ScriptOK || !strings.Contains(run.Summary, "stdout over 64 KiB was cut") || strings.Contains(run.Summary, "lost") {
		t.Fatalf("cut: %+v", run)
	}
}
