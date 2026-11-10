package workspace

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// A script hook is code from someone else's commit, run with the member's own
// permissions. It is launched as a process with no shell: the first element of
// its command is the program and the rest its arguments, and a relative path
// among them is resolved inside the extension folder, never outside it. It
// reads one JSON object on stdin, says what it did through its exit code and an
// optional final JSON object on stdout, and is stopped at its timeout.

// scriptOutputLimit is how much of a script's stdout, and of its stderr, is kept.
const scriptOutputLimit = 64 << 10

// defaultScriptTimeout applies where a manifest names no timeout.
const defaultScriptTimeout = 10 * time.Minute

// scriptGOOS is the platform a script command is chosen for. A test sets it to
// read the other platform's choice.
var scriptGOOS = runtime.GOOS

// ScriptRecord is the record a hook fires for, as a script receives it.
type ScriptRecord struct {
	ID           string   `json:"id"`
	Kind         string   `json:"kind"`
	Path         string   `json:"path"`
	Repositories []string `json:"repositories"`
}

// ScriptExtension is the extension a script belongs to, as it receives it.
type ScriptExtension struct {
	ID       string         `json:"id"`
	Settings map[string]any `json:"settings"`
}

// ScriptInput is the JSON a script reads on stdin. Record is null where the
// point carries no record; an archive passes every record it moved in Records.
type ScriptInput struct {
	Point         string            `json:"point"`
	WorkspaceRoot string            `json:"workspace_root"`
	Member        string            `json:"member"`
	Extension     ScriptExtension   `json:"extension"`
	Record        *ScriptRecord     `json:"record"`
	Records       []ScriptRecord    `json:"records,omitempty"`
	Worktrees     map[string]string `json:"worktrees"`
	Confirmed     bool              `json:"confirmed"`
}

// HookRun is one script's outcome: ok, failed, or skipped by the script
// itself, with what its on_failure makes of a failure.
type HookRun struct {
	ID          string `json:"id"`
	Status      string `json:"status"`
	OnFailure   string `json:"on_failure"`
	Summary     string `json:"summary,omitempty"`
	DetailsPath string `json:"details_path,omitempty"`
}

// Script statuses.
const (
	ScriptOK      = "ok"
	ScriptFailed  = "failed"
	ScriptSkipped = "skipped"
)

// scriptCommand chooses the command for the platform, or says why a hook has
// none there.
func scriptCommand(script HookScript, goos string) ([]string, string) {
	if goos == "windows" {
		if len(script.RunWindows) == 0 {
			return nil, "no Windows command: the manifest names no run_windows for this hook"
		}
		return script.RunWindows, ""
	}
	return script.Run, ""
}

// resolveScriptArgs makes each relative path in a command absolute inside the
// extension folder, because the script runs in a worktree or the workspace
// root, not in its own folder. A program is a path when it names a directory;
// otherwise it is looked up on PATH. An argument is a path when it names a
// directory and the file is there; any other argument is passed as written.
// A path that resolves outside the folder is refused.
func resolveScriptArgs(folder string, command []string) ([]string, error) {
	resolved := make([]string, len(command))
	for i, arg := range command {
		resolved[i] = arg
		if filepath.IsAbs(arg) || strings.HasPrefix(arg, "-") || strings.Contains(arg, "://") || !strings.ContainsAny(arg, `/\`) {
			continue
		}
		relative := filepath.ToSlash(arg)
		if err := insidePath(relative); err != nil {
			return nil, fmt.Errorf("%s %w", arg, errEscapes)
		}
		if _, err := os.Lstat(filepath.Join(folder, filepath.FromSlash(relative))); err != nil {
			if i == 0 {
				return nil, fmt.Errorf("the program %s is not in the extension folder", arg)
			}
			continue
		}
		path, err := resolveWithin(folder, relative)
		if err != nil {
			return nil, err
		}
		resolved[i] = path
	}
	return resolved, nil
}

// capped keeps the first limit bytes written to it and counts the rest, so a
// script that writes without end never fills memory and never blocks.
type capped struct {
	buf       bytes.Buffer
	limit     int
	truncated bool
}

func (c *capped) Write(p []byte) (int, error) {
	if room := c.limit - c.buf.Len(); room > 0 {
		if len(p) > room {
			c.buf.Write(p[:room])
			c.truncated = true
		} else {
			c.buf.Write(p)
		}
	} else if len(p) > 0 {
		c.truncated = true
	}
	return len(p), nil
}

// scriptOutput is the optional final JSON object a script writes on stdout.
type scriptOutput struct {
	Status      string `json:"status"`
	Summary     string `json:"summary"`
	DetailsPath string `json:"details_path"`
}

// finalObject reads the JSON object stdout ends with, if it ends with one: the
// shortest tail that starts a line and decodes whole.
func finalObject(out []byte) (scriptOutput, bool) {
	text := bytes.TrimSpace(out)
	if len(text) == 0 || text[len(text)-1] != '}' {
		return scriptOutput{}, false
	}
	for i := len(text) - 1; i >= 0; i-- {
		if text[i] != '{' || i > 0 && text[i-1] != '\n' {
			continue
		}
		var object scriptOutput
		if json.Unmarshal(text[i:], &object) == nil {
			return object, true
		}
	}
	return scriptOutput{}, false
}

// tail is the last part of a script's stream, short enough for a summary.
func tail(text string) string {
	text = strings.TrimSpace(strings.ReplaceAll(text, "\r\n", "\n"))
	lines := strings.Split(text, "\n")
	if len(lines) > 10 {
		lines = lines[len(lines)-10:]
	}
	text = strings.Join(lines, "\n")
	if len(text) > 2000 {
		text = "…" + text[len(text)-2000:]
	}
	return text
}

// scriptTimeout is the manifest's timeout, or the default.
func scriptTimeout(script HookScript) time.Duration {
	if duration, err := time.ParseDuration(script.Timeout); err == nil && duration > 0 {
		return duration
	}
	return defaultScriptTimeout
}

// runScript launches one script hook and folds its outcome into a HookRun. It
// never returns an error: what went wrong is the run's failure summary, so a
// script that cannot start is a failure its on_failure decides about.
func runScript(ctx context.Context, folder string, script HookScript, dir string, input ScriptInput) HookRun {
	run := HookRun{ID: input.Extension.ID, OnFailure: script.OnFailure}
	fail := func(summary string) HookRun {
		run.Status, run.Summary = ScriptFailed, summary
		return run
	}
	command, _ := scriptCommand(script, scriptGOOS)
	args, err := resolveScriptArgs(folder, command)
	if err != nil {
		return fail("not started: " + err.Error())
	}
	stdin, err := json.Marshal(input)
	if err != nil {
		return fail("not started: " + err.Error())
	}
	timeout := scriptTimeout(script)
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	process := exec.CommandContext(ctx, args[0], args[1:]...)
	process.Dir = dir
	process.Env = gitFreeEnviron()
	process.Stdin = bytes.NewReader(stdin)
	stdout := &capped{limit: scriptOutputLimit}
	stderr := &capped{limit: scriptOutputLimit}
	process.Stdout, process.Stderr = stdout, stderr
	// A script that leaves a child holding its output open would otherwise
	// keep the CLI waiting after the script itself was stopped.
	process.WaitDelay = 2 * time.Second
	err = process.Run()
	object, parsed := finalObject(stdout.buf.Bytes())
	if parsed {
		run.Summary, run.DetailsPath = object.Summary, object.DetailsPath
	}
	cut := ""
	if stdout.truncated {
		cut = fmt.Sprintf("stdout over %d KiB was cut", scriptOutputLimit>>10)
	}
	note := func(summary string) string {
		if cut == "" {
			return summary
		}
		if summary == "" {
			return cut
		}
		return summary + " (" + cut + ")"
	}
	switch {
	case errors.Is(ctx.Err(), context.DeadlineExceeded):
		return fail(note(fmt.Sprintf("timed out after %s", timeout)))
	case err != nil:
		summary := run.Summary
		if summary == "" {
			summary = tail(stderr.buf.String())
		}
		var exit *exec.ExitError
		if errors.As(err, &exit) {
			reason := fmt.Sprintf("exited with status %d", exit.ExitCode())
			if summary == "" {
				summary = reason
			} else {
				summary = reason + ": " + summary
			}
		} else {
			summary = "not started: " + err.Error()
		}
		return fail(note(summary))
	}
	run.Status = ScriptOK
	switch object.Status {
	case ScriptSkipped:
		run.Status = ScriptSkipped
	case ScriptFailed:
		// Exit code 0 with a failure said in so many words is still a failure.
		run.Status = ScriptFailed
	}
	run.Summary = note(run.Summary)
	return run
}
