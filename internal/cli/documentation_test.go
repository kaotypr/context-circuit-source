package cli_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kaotypr/context-circuit-source/internal/cli"
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

// documentedCommands returns every complete `context-circuit-cli ...` invocation in
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
		if pending == "" && !strings.HasPrefix(line, "context-circuit-cli ") {
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
	sources := []string{"docs/commands.md", "README.md"}
	for _, name := range shippedSkills(t) {
		sources = append(sources, filepath.Join("skills", name, "SKILL.md"))
	}
	for _, source := range sources {
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

// The entry instruction is always loaded; skills and docs are read on demand.
// A capability the entry instruction never names is one the agent never looks
// up, so every shipped skill and doc must be reachable from it — directly, or
// through a skill it names.
func TestEverythingShippedIsReachableFromTheEntryInstruction(t *testing.T) {
	agents := productFile(t, "AGENTS.md.in")
	reachable := agents
	for _, name := range shippedSkills(t) {
		if strings.Contains(agents, "`.agents/skills/"+name+"/SKILL.md`") {
			reachable += productFile(t, filepath.Join("skills", name, "SKILL.md"))
		}
	}
	for _, name := range shippedSkills(t) {
		if !strings.Contains(agents, "`.agents/skills/"+name+"/SKILL.md`") {
			t.Errorf("entry instruction never names the %s skill", name)
		}
	}
	for _, doc := range shippedDocs(t) {
		if !strings.Contains(reachable, "`.context-circuit/docs/"+doc+"`") {
			t.Errorf("%s is shipped but named by nothing an agent loads", doc)
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

// shippedSkills lists the skill directories carrying a canonical SKILL.md.
func shippedSkills(t *testing.T) []string {
	t.Helper()
	entries, err := os.ReadDir(filepath.Join("..", "..", "product", "skills"))
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, entry := range entries {
		if entry.IsDir() {
			names = append(names, entry.Name())
		}
	}
	if len(names) < 2 {
		t.Fatalf("expected shipped skills, found %d", len(names))
	}
	return names
}

// shippedDocs lists the on-demand reference the product installs.
func shippedDocs(t *testing.T) []string {
	t.Helper()
	entries, err := os.ReadDir(filepath.Join("..", "..", "product", "docs"))
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
			names = append(names, entry.Name())
		}
	}
	return names
}

// Placement matters: rules that apply to all implementation must stay in the
// always-loaded instruction, not move into a skill an agent loads only when it
// is running several plans at once.
func TestGeneralRulesStayOutOfTheStackedSkill(t *testing.T) {
	agents := productFile(t, "AGENTS.md.in")
	stacked := productFile(t, filepath.Join("skills", "cc-stacked", "SKILL.md"))
	for _, general := range []string{
		"Do not start independent verification during execution",
		"The `check` command is an explicitly invoked diagnostic",
		"Report what was implemented, tested, and left uncertain.",
	} {
		if !strings.Contains(agents, general) {
			t.Errorf("general implementation rule left the entry instruction: %q", general)
		}
		if strings.Contains(stacked, general) {
			t.Errorf("general implementation rule is trapped in cc-stacked: %q", general)
		}
	}
	// The gates and prohibitions are the entry instruction's job. A skill that
	// is never loaded must not be the only place one of them is written.
	for _, gate := range []string{
		"Implementation starts on a request to execute a plan",
		"require explicit authorization",
		// Committing moved sides; both halves of that split are gates.
		"leave no uncommitted work when the implementation is reported",
		"a commit in a bound checkout",
		"Never add agent attribution",
		"an intent is not\napproved because a field says so",
	} {
		if !strings.Contains(agents, gate) {
			t.Errorf("gate or prohibition missing from the entry instruction: %q", gate)
		}
	}
}

// An agent asking about one command must get that command. Answering a narrow
// question with the whole manual, and exiting zero, reads as success and sends
// the caller round again with a different phrasing.
func TestHelpAnswersTheCommandAsked(t *testing.T) {
	f := setup(t)
	full := f.ok("help")
	for _, topic := range []string{"member", "repo", "record", "worktree", "agent", "init", "version"} {
		out := f.ok("help", topic)
		if out == full {
			t.Errorf("help %s returned the whole manual", topic)
		}
		if !strings.Contains(out, topic) {
			t.Errorf("help %s never names it: %s", topic, out)
		}
		if detail := cli.Topics[topic]; detail != "" && !strings.Contains(out, detail) {
			t.Errorf("help %s omits its detail", topic)
		}
	}
	// A narrower name still resolves to that command's own signature.
	if out := f.ok("help", "record", "order"); !strings.Contains(out, "--mode auto|waves|linear") {
		t.Errorf("help record order: %s", out)
	} else if strings.Contains(out, "record approve") {
		t.Errorf("help record order widened to the whole group: %s", out)
	}
	// A name with no topic fails loudly rather than printing everything.
	if out := f.fail("help", "nonsense"); !strings.Contains(out, "no help topic: nonsense") {
		t.Errorf("unknown topic: %s", out)
	}
}

// A group followed by an option is a caller looking for help. Reporting it as an
// unknown command names the group as the fault when the group was correct.
func TestGroupWithOptionReportsTheOption(t *testing.T) {
	f := setup(t)
	out := f.fail("member", "--help")
	for _, want := range []string{"unknown option for member: --help", "help member"} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q: %s", want, out)
		}
	}
	if out := f.fail("record", "bogus"); !strings.Contains(out, "try: context-circuit-cli help record") {
		t.Errorf("unknown subcommand offers no topic: %s", out)
	}
}

// Every command the help lists must have detail behind it, and every topic must
// describe commands that exist, or one side has drifted from the other.
func TestEveryCommandHasATopic(t *testing.T) {
	seen := map[string]bool{}
	for _, line := range strings.Split(cli.Help, "\n") {
		if line == "" || strings.HasPrefix(line, " ") || !strings.Contains(line, "  ") && line != "version" {
			continue
		}
		name := strings.Fields(line)[0]
		if strings.HasSuffix(name, ":") || strings.Contains(name, "-") {
			continue
		}
		seen[name] = true
		if cli.Topics[name] == "" {
			t.Errorf("help lists %q with no topic detail", name)
		}
	}
	for name := range cli.Topics {
		if !seen[name] {
			t.Errorf("topic %q describes no listed command", name)
		}
		if _, ok := cli.CommandHelp(name); !ok {
			t.Errorf("topic %q resolves to nothing", name)
		}
	}
}
