# Authoring, dispatch, and compatibility

## Standalone authoring

The person asks for a detailed plan for a specified outcome. The agent first
retrieves relevant accepted knowledge, then reads the actual repository code
and its instructions, identifies implementation choices and unresolved product
decisions, and creates the plan folder. Unlike an intent, this route does not
ask a person to approve an agent-derived outcome before code investigation:
the requested outcome is already their specification. If investigation reveals
that the outcome itself is unsettled, the agent asks the person instead of
burying a guessed decision in an API spec or diagram.

The agent writes the overview first and adds only detail files that make the
proposal easier to inspect. A concrete API request/response contract, field
matrix, migration sequence, or diagram belongs in a concern file when it would
overload `plan.md`. The entry links it with a short reason to read it and marks
it required for all workers through `shared` or for the repositories that need
it through `by_repository`. The agent
presents the complete plan, then incorporates feedback in the folder. An
execution request comes after that review. A material change to the proposed
outcome goes back to the person before work
starts; an implementation detail can be revised in the plan.

Presentation names the plan entry, links each supporting file, and calls out
the choices that deserve the person's attention. The agent does not treat a
folder's existence as evidence that the person saw its contents. Feedback may
change the entry and its details before execution; the agent presents material
revisions again so the subsequent execution request refers to the current
proposal.

Intent-linked authoring still begins with an approved intent and then reads
code. The same folder layout applies, and the plan entry links its detail.
The intent remains the approved statement of outcome; the plan may explain
implementation at greater depth but must not quietly revise that outcome.

## Dispatch and implementation

The coordinator prepares a worktree for each repository named in the plan.
For a multi-repository plan it may dispatch one sole-owner worker per
repository, each with its actual worktree path and a task naming that
repository's slice. A worker dispatch names its `--repo` explicitly, so the
CLI can check that the repository belongs to the plan and select its reading
set. A single-repository plan may infer that sole repository. The CLI quotes
the plan's Markdown body and the metadata relevant to the worker, as the
current brief does; it parses but does not quote the entire `required_files`
frontmatter map. The brief states the local absolute
workspace root once, then lists only files under `required_files.shared` and
`required_files.by_repository.<repo>`, using paths relative to that root. It
instructs the
worker to resolve and read every listed file before editing. A workspace-root
relative path is never presented as if it were relative to the worker's
repository worktree. The coordinator must supply any cross-repository interface
decisions that are not settled by the plan; workers do not guess them from
sibling worktrees.

For example, the API worker's prompt carries the plan body and this selected
reading list, not the entire frontmatter map or the web-only file:

```text
Working directory: /local/workspace/.worktrees/p0001/api
Workspace root: /local/workspace
Plan: p0001
Repository assignment: api

Required files relative to workspace root:
- plans/p0001-billing/data-model.md
- plans/p0001-billing/api-contract.md

Read every required file before editing. Work only in the API worktree.
```

The plan body may still link to files assigned to another repository so the
person can review the whole plan; those links are not an instruction for this
worker to open them. The CLI resolves the plan-folder paths into the
workspace-relative paths above and validates them before returning a brief.

The CLI composes the brief and returns `launch_required`; the host launches the
worker. Large detail files are read from disk by the worker rather than copied
into every prompt. Binary illustrations are referenced by path, never embedded
in a structured CLI response. A missing required file prevents dispatch rather
than letting a worker implement from an incomplete brief. If a worker's host
cannot open a required illustration or other file type, it reports that limit
before editing rather than silently skipping the design. The coordinator waits
for results and integrates them in dependency order, then reports what changed
and which checks ran.
Neither the plan entry nor its details become a progress log during execution.
Completion is still a later human request and knowledge reconciliation still
follows it.

For a single standalone plan, ordering and delivery select it with
`record order --plan pNNNN`. For a run of several plans, the coordinator passes
each selected ID with repeatable `--plan`; the selected set is shown to the
person before the run and reused when deriving its delivery branches. An
unfinished dependency outside that set is reported as blocked, not quietly
added. Intent-linked runs may continue using `--intent`.

## Existing workspaces and versions

The workspace template and CLI are separate products. This design targets
v2.2.0 behavior in both, while the current source still declares v2.1.0 in
each version file. Release work must make the template pin a CLI that reads
folder plans before it creates them. A newer template is applied to an existing
workspace as a delta; it must not recopy the blank seed over that workspace.

The v2.2.0 CLI continues reading, editing, ordering, dispatching, and completing
legacy single-file plans. New plans use folders; existing files need no
migration. An optional explicit conversion could be designed later, but is
outside this proposal because moving a plan can break handwritten relative
links and concurrent references. The permanent ID ledger and branch naming
remain unchanged, so a legacy plan and a folder plan can depend on one another.
Dispatch of a legacy plan retains the existing quoted-body brief and has no
supporting-file reading map.

Implementation must update the shipped entry instruction, `cc-plan`,
`cc-dispatch`, `cc-stacked`, and `cc-deliver` skills and the workspace and
command documentation together. In particular, stacked execution and delivery
must accept explicitly selected standalone plan IDs instead of assuming an
intent groups every run. The CLI's record discovery, strict YAML schema,
ordering, dispatch, and diagnostic must agree with those instructions. Release
assembly must verify the updated shipped text and embedded seed; no design
file under `sources/` ships to a workspace.

Acceptance work should exercise: creating a standalone folder plan; creating
an intent-linked folder plan and backlink; mixed old/new listing and ID lookup;
duplicate ID refusal; completion and dependency ordering across both forms;
explicit plan-ID ordering and blocked outside dependencies; dispatch reading
selection and prompt for multi-file, multi-repository plans; unsafe, missing,
and unassigned supporting-file findings; archived folder lookup; interrupted
creation leaving a reserved ID; and release seed
inventory. Native host loading of the instructions and prompts remains
unverified until exercised in that host.
