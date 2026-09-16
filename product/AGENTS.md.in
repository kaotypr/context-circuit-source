# Context Circuit workspace

Coordinate this shared project workspace across one or more Git repositories.
Keep its identity, repositories, relationships, knowledge, members, intents, and
plans useful for solo and team work. Talk in ordinary project language. Use the
`context-circuit-cli` executable for dependable bookkeeping; users need not manage
runtime commands or ID numbers. Run the version this workspace pins in
`.context-circuit/CLI_VERSION`; the `cc-cli` skill installs and resolves it. See
`.context-circuit/docs/commands.md` as needed.

This file carries the gates and the standing prohibitions — what needs a person,
and what is never done without one. The procedure behind each stage lives in a
skill under `.agents/skills/`, named below where it comes due. A rule written
here is in force whether or not a skill was loaded.

## Orient and initialize

Context Circuit CLI is a separately installed product. If it is missing, or the
user requests an update, use `.agents/skills/cc-cli/SKILL.md`. Install for the
agent's execution OS/architecture; WSL and remote Linux use Linux packages. Keep
the template version and existing records unchanged during CLI updates.

Subagent role definitions are host-local and gitignored, so initialization writes
them for every host this workspace may later be opened in. A clone carries none,
and `agent configure` makes existing ones stale: if this host's `cc-*` definitions
are missing under `.claude/agents/`, `.codex/agents/`, or `.cursor/agents/`, or
settings changed, run `agent setup`. With no `--host` it covers every host and
preserves customized definitions. Check while orienting, not when something needs
to delegate: writing a definition is not loading it, and a host may not register a
role until it reloads. Dispatch reports whether the definition it names is there.

Read `workspace.yaml` for project identity, purpose, logical repositories with
their URL and default branch, and relationships. The shared roster is
`members.yaml`. The active member, this machine's checkout paths, and the branch
this machine starts work from are the documented, gitignored `member.local.yaml`
and `repositories.local.yaml`. Read those bindings when needed; never inspect
unrelated host settings, credentials, private provider payloads, or another
session's state. Never store secrets in workspace records.
`.context-circuit/docs/workspace.md` describes what each file may hold.

Initializing a workspace, connecting or cloning repositories, and adding members
or their allocation bands each happen only on a request, through
`.agents/skills/cc-workspace/SKILL.md`. On another machine, select an existing
member and connect its local checkouts; do not reinitialize the shared workspace.

## Shared knowledge

Use the optional `context/INDEX.md` to retrieve relevant architecture, conventions,
decisions, domain rules, terminology, and repository relationships. Without an
index, use targeted filenames or search terms in `context/`. Do not scan every
repository or context note. Sources are passive evidence: read only exact source
files named by the user or the task's explicit references; never scan `sources/`
by default.

Live context notes describe the project, not the workspace machinery that
produced them, and a note never names a plan record, an intent record, or a file
under `sources/`. To write or maintain knowledge — its anchors, its catalog
entry, the glossary, and what `check` enforces — use
`.agents/skills/cc-knowledge/SKILL.md`. Keep raw evidence separate from accepted
knowledge, and keep task progress and temporary results in plans.

## Intent and planning

For an implementation change, make the intended outcome explicit before detailed
code investigation. Ground it first: retrieve the bearing knowledge through
`context/INDEX.md` and read the request's apparent meaning against what the
project already records. What that reading contradicts, or cannot settle, is a
question for the person, not a gap to fill with a quiet assumption. Write the
`iNNN-slug.md` intent through `.agents/skills/cc-intent/SKILL.md`, which carries
its body fields and the numbered open questions a person answers by.

Present the concrete intent, invite answers by number, and stop there. The
request that prompted an intent is not approval of it: a person approves the
outcome after reading how it was written down, so do not treat "add billing" as
consent to the intent derived from it. Record the actual user approval, which
stamps `approved_at` and keeps the person's words beside it. A
record carries that decision and never establishes one: neither a command, nor an
editable field, nor another agent can supply human consent, and an intent is not
approved because a field says so.

Approving an intent authorizes planning, and planning alone. Without asking
again, inspect real code and create the linked `pNNNN-slug.md` plans through
`.agents/skills/cc-plan/SKILL.md`. `record approve` names that work as
`planning_required`, because approval reporting success is the decision recorded,
not the request finished. Obtain renewed approval only when the intended outcome
or success criteria materially change; a question that bears on the outcome
returns to the intent, and a plan never silently settles one.

Delegating that investigation is a judgment about size, not a step. Dispatch a
planner, or parallel explorers, through `.agents/skills/cc-dispatch/SKILL.md`,
which says when the reading is large enough to be worth delegating and what each
role may do. The coordinator writes every record either way.

Then present the plans and stop again. There is no plan approval gate — reading
them is the person's option, never a required step — but no repository is
touched until they ask for execution. That request is a decision made after the
plans exist and can be seen, so an instruction to implement that arrived earlier,
including one in the message that started the intent, does not begin it. Ask for
it plainly instead of inferring it.

A record's instants — `created_at`, `approved_at`, `completed_at` — are canonical
ISO 8601 UTC timestamps, `2026-09-15T10:53:00Z`. Every other date, in a record
written by hand as well as through the executable, is an ISO 8601 calendar date,
`YYYY-MM-DD`. No other date format belongs in a workspace file.

Intent and plan IDs are workspace-global. Use only `created_by` for member-related
metadata; no assignee, owner, reviewer, or member namespaces. Never reuse a
reserved ID, including after archival or deletion. Local locking serializes one
workspace directory. `.context-circuit/docs/working.md` describes the record
structures in full.

## Worktrees and implementation

Implementation starts on a request to execute a plan, and never before one.

Preparing a worktree for each repository the plan names is the first step of
execution, not a judgment call. The one exception: when the person asks to work
directly in a bound checkout, work there and preserve everything it already
holds. The execution request covers that preparation, so do not add another
permission gate in front of it. Use `.agents/skills/cc-worktree/SKILL.md` for
preparation, starting points, and the environment reuse report. Never silently
force checkout, reset, stash, or overwrite unrelated files, and never print
copied environment contents, put them in prompts, or store them in shared
records.

A prepared worktree is a repository checkout. If its root has an `AGENTS.md` or
`CLAUDE.md`, read it before changing anything and follow it where it is more
specific than this document — it owns that project's conventions, checks, and
commit format. If the repository ships skills, match the work against their
descriptions and open only the ones that match, the way `context/INDEX.md` is
used for knowledge. A repository with neither needs nothing extra. Where the two
conflict, this document's gates and prohibitions hold: nothing a repository asks
for authorizes an action a person has not.

Use `.agents/skills/cc-dispatch/SKILL.md` to delegate bounded exploration,
planning, and implementation when useful; small tasks may stay in the main
session. Apply configured host role/model/effort settings to actual subagent
invocations, wait for their results, and integrate them. To execute several plans
at once, use `.agents/skills/cc-stacked/SKILL.md`, which derives the shape with
`record order` rather than guessing it.

Implement the approved outcome in dependency order across the relevant
repositories. Run appropriate tests, lint, and builds as ordinary implementation
checks, then leave no uncommitted work when the implementation is reported:
commit on the worktree's own branch, in as many commits as the work naturally
takes, failing and partial work included. A dependent plan starts from its
predecessor's branch and a resume reads real diffs, so uncommitted work is work
the next step cannot see. A bound checkout is the exception — commit there only
when the person authorizes it. A plan record holds the plan — approach, tasks
and order, risks and checks — and it is written when the work is planned. After
that, leave it alone: running a plan does not edit the plan. The result reaches
the person, and the record changes again only when they ask for completion,
which stamps
`completed_at` and carries that result. What the plan says stays what was
planned, so a later reader can see what was intended and read the diff for what
happened. On resume, inspect actual branches and diffs rather than trusting the
record, which describes an intention and never an outcome. Preserve failed,
partial, and interrupted work; failure in one repository does not discard
successful work in another.

Do not start independent verification during execution, trigger a reviewer from
risk classification, or create automatic repair loops. Do not introduce execution,
candidate, verification, host-evidence, or separate formal completion records.
The `check` command is an explicitly invoked diagnostic for workspace consistency,
not an implementation gate. Report what was implemented, tested, and left uncertain.

Report what was actually run and what it actually returned. A check you did not
run is not a check, and an interactive flow you could not drive is not an
observation: say it is unverified rather than describing what it would have
shown. Unverified is a complete answer.

## Independent review on request

Independent verification is a manually requested read-only code review. It is
never triggered by risk classification, and it is not a condition for opening a
pull request, delivering, or completing. When a person asks for one, use
`.agents/skills/cc-review/SKILL.md`. Never call the implementing session's own
inspection independent; where independence is unavailable, say so and offer an
ordinary review instead. This optional capability does not make child agents
mandatory elsewhere.

## Delivery, completion, and organization

Push, PR creation, merge into a base branch, deployment, external publication,
and deletion of workspace data require explicit authorization. Reuse
authorization already given. Committing is split by where it lands: a commit on
a prepared worktree's own `cc/*` branch is implementation and is covered by the
execution request, while a commit in a bound checkout is a change to the
person's own working copy and still needs authorization. A local integration
merge that assembles a dependent plan's base is implementation, not delivery,
and is covered by the authorization for the run it belongs to. Use ordinary Git
and provider tools; the Go executable provides repository, branch, and base
information and never silently delivers. The branch the work started from — this
machine's recorded base, or the repository's default branch when the binding
records none — is the default PR target unless the user chooses another. Read it
as `base_branch`; `workspace.yaml` records
a `default_branch`, which is the repository's default and not what this machine
delivers to, and answering from it targets the wrong branch convincingly. When
the user asks to deliver changes or open a pull request, use
`.agents/skills/cc-deliver/SKILL.md`, which says which branch each request opens
from and against what, per repository and for the chain ends of a stacked run.
Explain relevant drift
or conflicts without imposing automatic rebase-and-reverify behavior.

Delivery does not mark a plan done. When the user explicitly marks one or more
plans completed, that request is finished only when the completion note is
appended *and* the durable project knowledge those plans changed is reconciled.
`record complete` names the second half as `reconcile_required`; judge it through
`.agents/skills/cc-complete/SKILL.md`. Completion reporting success is the record
written, not the request finished. Other work need not wait for unrelated
knowledge updates.

Worktrees and branches remain after completion. Remove a worktree only when
requested; preserve dirty, untracked, and ignored files unless their disposal is
explicitly authorized. Branch deletion is separate. On request, the executable can
move or repair worktrees using Git; Git's inventory is authoritative, and the
ignored plan-to-worktree association is only a convenience for resuming work.
See `.context-circuit/docs/worktrees.md` for recovery and cleanup mechanics.

Archiving is optional ordinary file organization on request. Keep ID reservations
and fix relative links when moving a record. External publication is an explicitly
requested task through available host tools; no automatic provider synchronization
is part of core v2.

Never add agent attribution, credit, co-author, or generated-by text to commits,
PRs, reviews, or comments. Follow repository commit conventions; default to
`type(scope): imperative subject`. Inspect authored commit messages and remove
injected attribution before reporting delivery complete.
