package cli_test

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/kaotypr/context-circuit-source/internal/cli"
)

// run calls the CLI with exactly the given arguments, so the position of the
// global options is the test's to choose.
func run(args ...string) (int, string, string) {
	var out, errOut bytes.Buffer
	code := cli.Run(context.Background(), args, &out, &errOut, "2.0.0-test")
	return code, out.String(), errOut.String()
}

// The global options work after the command as well as before it, for reading
// and changing commands alike. Without --workspace taking effect, these would
// run against the test's own directory, which is no workspace.
func TestGlobalOptionsAfterTheCommand(t *testing.T) {
	f := setup(t)
	for _, args := range [][]string{
		{"status", "--workspace", f.root, "--json"},
		{"--json", "status", "--workspace", f.root},
		{"--workspace", f.root, "status", "--json"},
		{"member", "list", "--json", "--workspace", f.root},
		{"status", "--workspace=" + f.root, "--json=true"},
	} {
		code, out, errOut := run(args...)
		var decoded map[string]any
		if code != 0 || json.Unmarshal([]byte(out), &decoded) != nil {
			t.Fatalf("%v: exit=%d %s %s", args, code, out, errOut)
		}
	}
	code, out, errOut := run("record", "create", "--kind", "intent", "--slug", "late", "--title", "Late", "--workspace", f.root, "--json")
	if code != 0 || record(t, out).ID != "i001" {
		t.Fatalf("a changing command with the options after it: exit=%d %s %s", code, out, errOut)
	}
	// Without --json anywhere the result is YAML, wherever --workspace sits.
	if code, out, errOut := run("record", "list", "--workspace", f.root); code != 0 || strings.HasPrefix(strings.TrimSpace(out), "[") || !strings.Contains(out, "id: i001") {
		t.Fatalf("YAML expected: exit=%d %s %s", code, out, errOut)
	}
}

// An unknown flag names itself and then the command's own usage, rather than
// an empty "Usage of status:" block.
func TestAnUnknownFlagPrintsTheCommandsUsage(t *testing.T) {
	f := setup(t)
	for _, c := range []struct {
		args []string
		want string
	}{
		{[]string{"status", "--nope"}, "status                inspect workspace"},
		{[]string{"record", "delete", "--nope"}, "record delete         --member ID | --all | --id ID"},
		{[]string{"worktree", "prepare", "--nope"}, "worktree prepare      --repo ID [--plan ID]"},
	} {
		code, out, errOut := run(append([]string{"--workspace", f.root}, c.args...)...)
		command := strings.Join(c.args[:len(c.args)-1], " ")
		for _, want := range []string{"flag provided but not defined: -nope", c.want, "may also follow the full command", "try: context-circuit-cli help " + command} {
			if code != 2 || !strings.Contains(errOut, want) {
				t.Fatalf("%v: want %q: exit=%d %s %s", c.args, want, code, out, errOut)
			}
		}
		if strings.Contains(errOut, "Usage of") {
			t.Fatalf("%v: the flag package's own usage is still printed: %s", c.args, errOut)
		}
	}
	// Asking for help is an answer, on the result stream.
	if code, out, _ := run("status", "-h"); code != 0 || !strings.Contains(out, "Read-only: it changes nothing") {
		t.Fatalf("status -h: exit=%d %s", code, out)
	}
}

// Every hint names the command the running version is installed as, so a
// variant build never sends its caller to the public command.
func TestHelpHintsNameTheRunningCommand(t *testing.T) {
	for version, command := range map[string]string{
		"2.5.0":        "context-circuit-cli",
		"2.5.0-acme.1": "acme-context-circuit-cli",
		"2.5.0-rc.1":   "context-circuit-cli",
	} {
		for _, c := range []struct {
			args []string
			want string
		}{
			{nil, "Usage: " + command + " [--workspace PATH]"},
			{[]string{"help"}, "Usage: " + command + " [--workspace PATH]"},
			{[]string{"--help"}, "Usage: " + command + " [--workspace PATH]"},
			{[]string{"record"}, "Usage: " + command + " [--workspace PATH]"},
			{[]string{"help", "bogus"}, "try: " + command + " help\n"},
			{[]string{"record", "--nope"}, "try: " + command + " help record\n"},
			{[]string{"record", "bogus"}, "try: " + command + " help record\n"},
			{[]string{"status", "--nope"}, "usage: " + command + " [--workspace PATH] [--json] status [OPTIONS]"},
			{[]string{"status", "--nope"}, "try: " + command + " help status\n"},
			{[]string{"--nope"}, "Usage of " + command + ":"},
		} {
			var out, errOut bytes.Buffer
			cli.Run(context.Background(), c.args, &out, &errOut, version)
			text := out.String() + errOut.String()
			if !strings.Contains(text, c.want) {
				t.Errorf("%s %v: want %q in:\n%s", version, c.args, c.want, text)
			}
			if command == "context-circuit-cli" && strings.Contains(text, "-context-circuit-cli") {
				t.Errorf("%s %v: a labelled command in public output:\n%s", version, c.args, text)
			}
		}
	}
	// The signatures a topic prints are read out of the public text unchanged.
	if text, ok := cli.CommandHelp("status"); !ok || !strings.HasPrefix(text, "status ") {
		t.Fatalf("status topic: %v %s", ok, text)
	}
}
