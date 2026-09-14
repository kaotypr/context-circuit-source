package cli_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// An agent learns this CLI from the shipped instruction and docs, so a command
// that drifted from its documentation is a command the agent will get wrong.
// These tests read the real product files and check the surface they teach.

func productFile(t *testing.T, rel string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "..", "product", rel))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

// documentedCommands returns every complete `context-circuit ...` invocation in
// fenced examples, rejoining the backslash continuations the docs wrap with.
func documentedCommands(text string) []string {
	var commands []string
	var pending string
	fenced := false
	for _, line := range strings.Split(text, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "```") {
			fenced = !fenced
			continue
		}
		if !fenced {
			continue
		}
		line = strings.TrimSpace(line)
		if pending == "" && !strings.HasPrefix(line, "context-circuit ") {
			continue
		}
		if strings.HasSuffix(line, "\\") {
			pending += strings.TrimSuffix(line, "\\") + " "
			continue
		}
		commands = append(commands, strings.TrimSpace(pending+line))
		pending = ""
	}
	return commands
}

// A usage template teaches shape, not a runnable line; its metavariables stand
// in for arguments rather than naming a command.
func isTemplate(command string) bool {
	for _, token := range strings.Fields(command) {
		if token == "COMMAND" || token == "help" {
			return true
		}
	}
	return false
}

func TestDocumentedCommandsExist(t *testing.T) {
	f := setup(t)
	f.repository("api")
	var checked int
	for _, source := range []string{"docs/commands.md", "skills/cc-dispatch/SKILL.md", "README.md"} {
		for _, command := range documentedCommands(productFile(t, source)) {
			if isTemplate(command) {
				continue
			}
			args := strings.Fields(command)[1:]
			// Drop a documented --workspace so every example runs against the
			// fixture; the rest of the line is checked exactly as written.
			for i := 0; i+1 < len(args); i++ {
				if args[i] == "--workspace" {
					args = append(args[:i], args[i+2:]...)
					break
				}
			}
			args = unquote(args)
			checked++
			_, out, errOut := call(f.root, args...)
			// Placeholders make many examples fail at execution, which is fine.
			// A surface error means the documented command or flag is not real.
			for _, surface := range []string{
				"unknown command",
				"flag provided but not defined",
				"unexpected positional arguments",
				"missing --",
			} {
				if strings.Contains(errOut, surface) {
					t.Errorf("%s teaches a command the CLI does not accept (%s):\n  %s\n  %s",
						source, surface, command, strings.TrimSpace(out+errOut))
					break
				}
			}
		}
	}
	if checked < 20 {
		t.Fatalf("expected the docs to teach a real command surface, only found %d", checked)
	}
}

// unquote collapses the shell quoting the examples use for readability.
func unquote(args []string) []string {
	var out []string
	var joining bool
	var buffer string
	for _, arg := range args {
		switch {
		case joining:
			buffer += " " + arg
			if strings.HasSuffix(arg, "'") {
				out = append(out, strings.Trim(buffer, "'"))
				joining = false
			}
		case strings.HasPrefix(arg, "'") && !strings.HasSuffix(arg, "'"):
			joining, buffer = true, arg
		default:
			out = append(out, strings.Trim(arg, "'"))
		}
	}
	if joining {
		out = append(out, strings.Trim(buffer, "'"))
	}
	return out
}

// The entry instruction is always loaded; the detailed docs are read on demand.
// A capability the entry instruction never names is one the agent never looks up.
func TestEntryInstructionNamesTheStackedFlow(t *testing.T) {
	agents := productFile(t, "AGENTS.md.in")
	for _, want := range []string{
		"`record order`", // the command, on one line so it can be searched
		"`.context-circuit/docs/commands.md`",
		"`.agents/skills/cc-dispatch/SKILL.md`",
		"## Stacked plans",
	} {
		if !strings.Contains(agents, want) {
			t.Errorf("entry instruction never mentions %s", want)
		}
	}
	// A backtick command split over a line break cannot be found by search and
	// reads as two words.
	for _, command := range []string{"record order", "record complete", "worktree prepare", "agent dispatch"} {
		parts := strings.Fields(command)
		if strings.Contains(agents, "`"+parts[0]+"\n") {
			t.Errorf("%q is split across a line break in the entry instruction", command)
		}
	}
}

// Placement matters: rules that apply to all implementation must not sit under a
// heading an agent skips when it is running a single plan.
func TestStackedSectionDoesNotCaptureGeneralRules(t *testing.T) {
	agents := productFile(t, "AGENTS.md.in")
	start := strings.Index(agents, "## Stacked plans")
	if start < 0 {
		t.Fatal("missing stacked plans section")
	}
	section := agents[start:]
	if next := strings.Index(section[3:], "\n## "); next >= 0 {
		section = section[:next+3]
	}
	for _, general := range []string{
		"Do not start independent verification during execution",
		"The `check` command is an explicitly invoked diagnostic",
		"Report what was implemented, tested, and left uncertain.",
	} {
		if strings.Contains(section, general) {
			t.Errorf("general implementation rule is trapped under Stacked plans: %q", general)
		}
	}
}
