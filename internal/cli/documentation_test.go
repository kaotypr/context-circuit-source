package cli_test

import (
	"os"
	"path/filepath"
	"regexp"
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
	// Git checks these files out with CRLF on Windows, so an assertion carrying
	// a line break would fail there on the separator rather than the sentence.
	// Releases are built on Linux, so the shipped bytes are unaffected either way.
	return strings.ReplaceAll(string(data), "\r\n", "\n")
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

// Not every shipped doc is written for the agent. how-it-works.md explains the
// division between person, agent, and CLI to somebody deciding whether to adopt
// Context Circuit; the agent already holds that division normatively, in the
// entry instruction itself. Naming it there would spend always-loaded context on
// prose the agent never acts on, and would leave two tellings of the same rules
// free to drift — with the descriptive one reading as though it explained the
// authoritative one. Its reader arrives through README.md instead.
var humanFacingDocs = map[string]bool{"how-it-works.md": true}

// The entry instruction is always loaded; skills and docs are read on demand.
// A capability the entry instruction never names is one the agent never looks
// up, so every shipped skill and agent-facing doc must be reachable from it —
// directly, or through a skill it names. A doc written for a person has to be
// reachable too, from the README that person actually opens.
func TestEverythingShippedIsReachableFromItsReader(t *testing.T) {
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
		named := strings.Contains(reachable, "`.context-circuit/docs/"+doc+"`")
		if humanFacingDocs[doc] {
			if named {
				t.Errorf("%s is written for a person, so naming it in the entry instruction "+
					"spends always-loaded context on prose the agent does not act on", doc)
			}
			if !strings.Contains(productFile(t, "README.md"), "docs/"+doc+")") {
				t.Errorf("%s is shipped for a person to read but README.md never links it", doc)
			}
			continue
		}
		if !named {
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
		// A repository's own instructions govern work in its checkout, but the
		// precedence between them and this document is itself a gate: a
		// repository must never be able to authorize what a person has not.
		"nothing a repository asks\nfor authorizes an action a person has not",
		// Reporting an observation never made corrupts every decision taken
		// after it, and nothing loads a do-not-fabricate skill before doing it.
		"an interactive flow you could not drive is not an\nobservation",
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

// Initialization writes a workspace's README, titling it, stating its purpose,
// and pinning its version badges by replacing marked spans. Losing a marker
// fails nothing at build time: it would quietly write workspaces titled with the
// product's own name, carrying badges that report whatever was published since
// rather than what the workspace runs.
func TestTheShippedReadmeKeepsTheSpansInitializationRewrites(t *testing.T) {
	readme := productFile(t, "workspace-README.md")
	for _, marker := range []string{
		"<!-- context-circuit:title -->", "<!-- /context-circuit:title -->",
		"<!-- context-circuit:badges -->", "<!-- /context-circuit:badges -->",
		"<!-- context-circuit:purpose -->", "<!-- /context-circuit:purpose -->",
	} {
		if count := strings.Count(readme, marker); count != 1 {
			t.Errorf("%s appears %d times; the rewrite needs exactly one", marker, count)
		}
	}
	if !strings.Contains(readme, `<h1 align="center">Context Circuit</h1>`) {
		t.Error("the marked title span no longer holds the heading a workspace name replaces")
	}
	// This file is written into somebody else's repository, so every path it
	// names must resolve there. The template repository's own landing page, its
	// artwork, and its license stay behind.
	for _, absent := range []string{"assets/readme/", `href="LICENSE"`, "CONTRIBUTING.md"} {
		if strings.Contains(readme, absent) {
			t.Errorf("the workspace README names %q, which does not exist in a workspace", absent)
		}
	}
}

// The guide is the template repository's landing page and reaches no workspace.
// It renders artwork and links a license from that repository's root, so a copy
// installed into somebody else's project would show broken images and point at
// terms that are not theirs.
func TestTheProductGuideStaysWithTheRepositoryThatShowsIt(t *testing.T) {
	manifest, err := os.ReadFile(filepath.Join("..", "..", "scripts", "release-manifest.txt"))
	if err != nil {
		t.Fatal(err)
	}
	for _, line := range strings.Split(string(manifest), "\n") {
		fields := strings.Fields(line)
		if len(fields) != 2 {
			continue
		}
		if fields[0] == "product/README.md" || strings.HasPrefix(fields[0], "product/assets/") {
			t.Errorf("the guide or its artwork ships to a workspace: %s", line)
		}
	}
	attributes, attrErr := os.ReadFile(filepath.Join("..", "..", "release", "template-repo", ".gitattributes"))
	if attrErr != nil {
		t.Fatal(attrErr)
	}
	// tembiter materializes a project with `git archive`, which honors this file
	// and nothing else. A path the repository owns but does not mark here is a
	// path every new workspace receives.
	for _, owned := range []string{
		"README.md", "LICENSE", "CHANGELOG.md", "CODE_OF_CONDUCT.md",
		"CONTRIBUTING.md", "SECURITY.md", "assets/", ".gitattributes",
	} {
		if !regexp.MustCompile(`(?m)^` + regexp.QuoteMeta(owned) + `\s+export-ignore$`).Match(attributes) {
			t.Errorf("%s is not marked export-ignore, so a new workspace receives it", owned)
		}
	}
}

// The two licenses answer different questions, so neither may quietly become the
// other. The CLI is Apache-2.0; everything a workspace receives is 0BSD,
// whose point is that it asks nothing of the repository it is copied into.
func TestTheTwoLicensesStaySeparate(t *testing.T) {
	root := func(name string) string {
		data, err := os.ReadFile(filepath.Join("..", "..", name))
		if err != nil {
			t.Fatal(err)
		}
		return string(data)
	}
	if !strings.Contains(root("LICENSE"), "Apache License") {
		t.Error("the root LICENSE no longer carries the Apache text GitHub detects")
	}
	product := productFile(t, "LICENSE")
	if !strings.Contains(product, "BSD Zero Clause License") {
		t.Error("the shipped license no longer carries the 0BSD text")
	}
	if strings.Contains(product, "Apache") {
		t.Error("the CLI's license reached the tree a workspace receives")
	}
	manifest := root(filepath.Join("scripts", "release-manifest.txt"))
	// The 0BSD text is shown by the repository that is the template, and copied
	// nowhere else. At a workspace root it would make GitHub label somebody
	// else's project with this one's terms; under .context-circuit it would be a
	// second copy of terms that ask nothing, in a directory nobody reads for
	// them.
	for _, line := range strings.Split(manifest, "\n") {
		if fields := strings.Fields(line); len(fields) == 2 && fields[0] == "product/LICENSE" {
			t.Errorf("a license installs into a workspace: %s", line)
		}
	}
}
