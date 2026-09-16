package cli

import "strings"

// groups are the commands that take a subcommand. Naming them once keeps the
// dispatcher and the help topics from disagreeing about what a group is.
var groups = map[string]bool{
	"member": true, "repo": true, "record": true, "worktree": true,
	"context": true, "template": true, "agent": true,
}

// Topics carry what a signature line cannot: what the command decides, what it
// refuses, and what the caller still owns afterwards. The signatures themselves
// are read out of Help, so a topic cannot drift from the flags it documents.
var Topics = map[string]string{
	"init": `Creates a workspace in place from the embedded seed: instructions, folders,
initial records, and native subagent role definitions for every supported host.
Existing files are preserved. It neither migrates a v1 workspace nor
reinitializes an active v2 one; on a second machine, select an existing member
and connect local checkouts instead.`,

	"status": `Reports the workspace, its members, this machine's bindings, and Git state for
each bound repository. Read-only: it changes nothing and waits on nothing.`,

	"check": `Reports record, binding, dependency, and worktree problems, including a note
that names a record or evidence path it should not. Findings marked readability
report the shape of a note instead: an anchor written into a sentence, a
paragraph carrying a list it never made, a catalog heading that never says what
belongs under it. It is an explicitly invoked diagnostic that reports and exits,
never an implementation or delivery gate.`,

	"member": `The roster in members.yaml is shared; the active member in member.local.yaml is
this machine's alone. A band is a numeric block one member allocates from, which
is what lets clones that cannot see each other avoid colliding. A band is not a
namespace: IDs stay workspace-global and the member is recorded only as
created_by.`,

	"repo": `Binds a logical repository ID to a local checkout. The ID, the repository URL,
and its default branch are shared in workspace.yaml; the checkout path and the
branch this machine starts work from stay in repositories.local.yaml. Worktrees
branch from that local base, so one machine can work off a release branch while
another stays on the default. Connect an existing checkout, clone, or initialize
a new repository; the workspace root itself may be connected as ".". Fetch
updates refs only and implies no rebase or reset.`,

	"record": `Creates and updates intents and plans, allocating IDs that are never reused.
Approve and complete record a decision a person already made; neither command
constitutes that decision, and no command executes a plan. Order derives
dependency waves, starting references, and integration merges without running or
reserving anything.`,

	"context": `Finds catalog entries in context/INDEX.md by substring, or, with --repo, the
entries claiming one repository and the reconciliation they leave outstanding —
what completion names for a plan, for a change made without one. It locates
candidates; judging relevance and writing knowledge stay with the agent.`,

	"worktree": `Prepares, inspects, moves, repairs, and removes isolated working copies.
Preparation may reuse ignored runtime files from the bound checkout through
filesystem cloning, skipping dependencies whose package inputs differ. Git's own
inventory is authoritative; the plan-to-worktree association is a convenience for
resuming. Removal preserves dirty, untracked, and ignored files unless discard is
explicitly authorized.`,

	"template": `Exports blank workspace files to a new directory, for seeding a workspace that
init will complete. It never writes over an existing directory.`,

	"agent": `Resolves per-host role model and effort settings, writes native role definitions,
and composes dispatch specifications. Setup covers every host unless --host
narrows it, because one workspace is opened in several. A dispatch specification
reports the definition it names and whether it is installed, and returns a flag
saying a launch is still required: this executable holds no model credentials,
launches nothing, and a specification is never evidence that an agent ran.`,

	"version": `Prints the running CLI version. A workspace pins its own in
.context-circuit/CLI_VERSION, and the two need not match; the CLI warns when they
differ.`,
}

// CommandHelp answers a question about one command. It returns the command's
// signatures exactly as the full help lists them, followed by that topic's
// detail, and reports false when the name matches neither.
func CommandHelp(topic string) (string, bool) {
	topic = strings.Join(strings.Fields(topic), " ")
	if topic == "" {
		return "", false
	}
	var lines []string
	keep := false
	for _, line := range strings.Split(Help, "\n") {
		switch {
		case line == topic || strings.HasPrefix(line, topic+" "):
			keep = true
		case keep && strings.HasPrefix(line, "    ") && strings.TrimSpace(line) != "":
			// A continuation line carries more flags for the command above it.
		default:
			keep = false
		}
		if keep {
			lines = append(lines, strings.TrimRight(line, " "))
		}
	}
	detail := Topics[strings.Fields(topic)[0]]
	if len(lines) == 0 && detail == "" {
		return "", false
	}
	var text strings.Builder
	if len(lines) > 0 {
		text.WriteString(strings.Join(lines, "\n") + "\n")
	}
	if detail != "" {
		text.WriteString("\n" + detail + "\n")
	}
	return text.String(), true
}
