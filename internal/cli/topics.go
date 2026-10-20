package cli

import "strings"

// groups are the commands that take a subcommand. Naming them once keeps the
// dispatcher and the help topics from disagreeing about what a group is.
var groups = map[string]bool{
	"member": true, "repo": true, "record": true, "worktree": true,
	"context": true, "template": true, "agent": true, "workspace": true,
	"knowledge": true,
}

// Topics carry what a signature line cannot: what the command decides, what it
// refuses, and what the caller still owns afterwards. The signatures themselves
// are read out of Help, so a topic cannot drift from the flags it documents.
var Topics = map[string]string{
	"init": `Creates a workspace in place from the embedded seed: instructions, folders,
initial records, and native subagent role definitions for every supported host.
Existing files are preserved. It requires the first member, writes them to the
roster with band 1 unless --member-band names another, and selects them on this
machine, so a person working alone is a team of one with nothing left to set up.
It neither migrates a v1 workspace nor
reinitializes an active v2 one; on a second machine, select an existing member
and connect local checkouts instead.`,

	"status": `Reports the workspace, its members, this machine's bindings, and Git state for
each bound repository. Read-only: it changes nothing and waits on nothing.`,

	"check": `Reports record, binding, dependency, and worktree problems, including a note
that names a record or evidence path it should not.

Every finding carries both what is wrong and what discharges it: the command to
run where one exists, the edit to make where an edit is the whole of it, or a
resolution opening with "needs a person" where the next step is somebody's
judgment and no command should be improvised in its place. Findings marked readability
report the shape of a note instead: an anchor written into a sentence, a
paragraph carrying a list it never made, a catalog heading that never says what
belongs under it. It is an explicitly invoked diagnostic that reports and exits,
never an implementation or delivery gate.`,

	"member": `The roster in members.yaml is shared; the active member in member.local.yaml is
this machine's alone. Every member holds a band, a numeric block that member
allocates from, which is what lets clones that cannot see each other avoid
colliding. Band 0 is a band like any other, covering i001-i099 and p0001-p0999;
no two members share one. A band is not a namespace: IDs stay workspace-global
and the member is recorded only as created_by. Add requires --band and, without
it, names the next free band. Until the roster lists this machine's selected
member and every member holds a band, status reports identity_required or
band_required with band_suggestions, and every changing command except these
setup steps is refused. Selection is taken as given: nothing verifies who a
person is.

Roles are free text saying what a member contributes, such as maintainer or
docs, shown wherever the roster is. They are informational and decide nothing:
approval, allocation, deletion, and dispatch never read them, and they are not
the subagent --role that agent configure and agent dispatch take. They are
written only at workspace schema 3, so a member added while an upgrade is still
mid-setup is added without roles, which member roles sets afterwards.`,

	"knowledge": `Mounts a repository of knowledge this workspace reads and never owns: an
organization's knowledge center, shared by many workspaces and changed through
its own repository. The ID, URL, tracked branch, and index file are shared in
workspace.yaml; the checkout path stays in repositories.local.yaml. It is not a
repository where work happens, so it holds no base branch, takes no worktree,
and is never named by a plan; the two ID namespaces are kept apart so a note's
anchor means one thing.

Read-only here is enforced, not requested: obtaining one disables its push URL,
and sync reasserts that. Sync fetches, then fast-forwards only when the checkout
is clean and on the shared branch. It never merges, rebases, resets, or
discards: a checkout holding local work is reported and left exactly as it is,
because moving that work is a decision for the person who made it. Improvements
go upstream through that repository's own review, from a separate checkout.

Retrieval reads it beside this workspace's own catalog, and reports any index it
could not read rather than narrowing the result in silence. Nothing validates
its contents the way check validates context/: the catalog rules exist so a
person here can fix what they break, and neither half holds for a repository
this workspace does not own. Completion never names a borrowed entry as
knowledge to reconcile, because nobody here can move it.`,

	"repo": `Binds a logical repository ID to a local checkout. The ID, the repository URL,
and its default branch are shared in workspace.yaml; the checkout path and the
branch this machine starts work from stay in repositories.local.yaml. Worktrees
branch from that local base, so one machine can work off a release branch while
another stays on the default. Connect an existing checkout, clone, or initialize
a new repository. The workspace's own checkout is refused here and described by
the workspace commands instead, because an entry in this map is somewhere work
happens. Fetch updates refs only and implies no rebase or reset.`,

	"record": `Creates and updates intents and plans, allocating IDs from a permanent ledger.
Approve and complete record a decision a person already made; neither command
constitutes that decision, and no command executes a plan. Delete removes one
member's intents and plans, active and archived, every record while the roster
lists one member, or the records named with repeatable --id, whoever wrote them,
and is the only operation that releases a reserved ID. Without --confirm it
previews what would go, with each record's author and status, and deletes
nothing. It refuses while a record outside the selection, or a worktree on this
machine, still points at one of them; by ID, a deleted plan is taken out of its
own intent's plans list instead. Leftover cc/pNNNN branches and worktree
directories are warned about rather than refused. IDs are released only after
the files are gone, so a deletion that stops part-way names the --id call that
finishes it. Order derives
dependency waves, starting references, and integration merges without running or
reserving anything.`,

	"context": `Finds catalog entries in context/INDEX.md by substring, or, with --repo, the
entries claiming one repository and the reconciliation they leave outstanding —
what completion names for a plan, for a change made without one. It locates
candidates; judging relevance and writing knowledge stay with the agent.`,

	"worktree": `Prepares, inspects, moves, repairs, and removes isolated working copies, and
prunes stale plan-to-worktree associations. Preparation may reuse ignored runtime files from the bound checkout through
filesystem cloning, skipping dependencies whose package inputs differ. Git's own
inventory is authoritative; the plan-to-worktree association is a convenience for
resuming. Removal preserves dirty, untracked, and ignored files unless discard is
explicitly authorized. Prune only drops the associations of worktrees Git no
longer lists, such as ones removed with git itself; it touches no worktree,
branch, or file, and reports what it dropped as dropped. Preparation refuses a
branch already checked out elsewhere, says where and what is safe to do next, and
changes nothing. Without --reuse it also refuses a worktree or branch left from an
earlier plan that held the same number, judged by the plan's created_at, and a
worktree from before 2.4.0 holding no commits unless it is clean and at the base
tip; preparing again for the same plan resumes its worktree.`,

	"template": `Exports blank workspace files to a new directory, for seeding a workspace that
init will complete. It never writes over an existing directory.`,

	"agent": `Resolves per-host role model and effort settings, writes native role definitions,
and composes dispatch specifications. Setup covers every host unless --host
narrows it, because one workspace is opened in several. A dispatch specification
reports the definition it names and whether it is installed, and returns a flag
saying a launch is still required: this CLI holds no model credentials,
launches nothing, and a specification is never evidence that an agent ran.`,

	"workspace": `Describes the Git repository carrying the workspace itself, so status reports
its branch and working state and a second machine is told where to clone it. It
splits like a repository and for the same reason: the URL and default branch are
shared in workspace.yaml under workspace_repository, while the checkout root and
the branch this machine works from stay in repositories.local.yaml. It is kept
out of the repositories map because a plan, a relationship, and a worktree all
read an entry there as somewhere work happens. No command commits, pushes, or
merges the workspace.`,

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
