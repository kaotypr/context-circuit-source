//go:build windows

package workspace

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// scriptFolder writes an extension folder holding one PowerShell script.
func scriptFolder(t *testing.T, body string) string {
	t.Helper()
	folder := t.TempDir()
	if err := os.MkdirAll(filepath.Join(folder, "scripts"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(folder, "scripts", "hook.ps1"), []byte(body), 0644); err != nil {
		t.Fatal(err)
	}
	return folder
}

// On Windows the hook's run_windows command is the one launched; run names a
// program Windows does not have, so launching it would fail the test.
func powershellScript(onFailure string) HookScript {
	return HookScript{
		Run:        []string{"cc-no-such-unix-program", "scripts/hook.sh"},
		RunWindows: []string{"powershell", "-NoProfile", "-NonInteractive", "-ExecutionPolicy", "Bypass", "-File", "scripts/hook.ps1"},
		OnFailure:  onFailure,
	}
}

func testInput() ScriptInput {
	return ScriptInput{Point: "after implement", WorkspaceRoot: `C:\w`, Member: "maya",
		Extension: ScriptExtension{ID: "parity", Settings: map[string]any{"reference": "legacy"}},
		Record:    &ScriptRecord{ID: "p0001", Kind: "plan", Path: `C:\w\plans\p0001-x\plan.md`, Repositories: []string{"api"}}}
}

func TestAScriptReadsItsInputAndRunsInItsDirectory(t *testing.T) {
	dir := t.TempDir()
	folder := scriptFolder(t, `$text = [Console]::In.ReadToEnd()
Set-Content -Path input.json -Value $text
Set-Content -Path cwd.txt -Value (Get-Location).Path
Write-Output 'progress line'
Write-Output '{"status":"ok","summary":"14 behaviors match"}'
`)
	run := runScript(context.Background(), folder, powershellScript("stop"), dir, testInput())
	if run.Status != ScriptOK || run.Summary != "14 behaviors match" {
		t.Fatalf("run: %+v", run)
	}
	data, err := os.ReadFile(filepath.Join(dir, "input.json"))
	if err != nil {
		t.Fatal(err)
	}
	var input ScriptInput
	if err := json.Unmarshal(data, &input); err != nil || input.Point != "after implement" || input.Record.ID != "p0001" {
		t.Fatalf("input: %v %s", err, data)
	}
	cwd, _ := os.ReadFile(filepath.Join(dir, "cwd.txt"))
	// The test's temporary directory may be a short 8.3 name, such as
	// RUNNER~1, while the script sees the long one: compare the directories,
	// not their spelling.
	ran, rerr := os.Stat(strings.TrimSpace(string(cwd)))
	want, werr := os.Stat(dir)
	if rerr != nil || werr != nil || !os.SameFile(ran, want) {
		t.Fatalf("ran in %s, want %s", cwd, dir)
	}
}

func TestAScriptsExitCodeDecidesFailure(t *testing.T) {
	dir := t.TempDir()
	run := runScript(context.Background(), scriptFolder(t, "[Console]::Error.WriteLine('GET /orders/1 differs')\nexit 3\n"), powershellScript("stop"), dir, testInput())
	if run.Status != ScriptFailed || !strings.Contains(run.Summary, "status 3") {
		t.Fatalf("failure: %+v", run)
	}
	run = runScript(context.Background(), scriptFolder(t, `Write-Output '{"status":"skipped","summary":"the plan has no parity list"}'`), powershellScript("stop"), dir, testInput())
	if run.Status != ScriptSkipped || run.Summary != "the plan has no parity list" {
		t.Fatalf("skip: %+v", run)
	}
}

func TestAScriptIsStoppedAtItsTimeout(t *testing.T) {
	script := powershellScript("report")
	script.Timeout = "1s"
	run := runScript(context.Background(), scriptFolder(t, "Start-Sleep -Seconds 30\n"), script, t.TempDir(), testInput())
	if run.Status != ScriptFailed || !strings.Contains(run.Summary, "timed out after 1s") {
		t.Fatalf("timeout: %+v", run)
	}
}
