# Agent Workspace Workflow

Context Circuit is an AI-agent workspace. Its primary workflow is a session
that enters or resumes the workspace, understands the current state, performs
one safe development loop, and leaves a durable handoff.

The full agreed contract is documented in
docs/agent-workspace-workflow.md. This file is the normative entry summary
that every root and child session must read.
Runtime record shapes and lease ownership are defined in
docs/runtime-contract.md.

## Session entry

The agent must:

1. Read AGENTS.md, WORKFLOW.md, and workspace.yaml.
2. Read context/INDEX.md, context/WORKSPACE.md, context/PROJECT.md, and only the relevant Product Knowledge.
3. Identify whether it is a root session or a child session.
4. Read its runtime session record and parent handoff when applicable.
5. Inspect relevant plans, repository instructions, branch, worktree, and Git
   state.
6. Explain whether it is orienting, gathering context, planning, awaiting
   approval, executing, verifying, blocked, or handing off.

On fresh root entry, create the root session record before claiming work. On
resume, use the explicit session record and latest handoff; never infer a
current session from conversation history or a global pointer.

The user should be able to say “start or resume work in this workspace” without
knowing the internal command implementation.

Choosing the next action is part of this root entry, not a separate command.
Inspect approved plans with unfinished work, declared dependencies, active
leases and session ownership, source freshness, repository cleanliness,
worktrees, blockers, and pending human gates. Recommend or claim only work
that is dependency-ready, explicitly scoped, and not already owned by another
writing session. When no work is executable, explain whether the session needs
context, a draft plan, human approval, a review, or a decision about a
blocker. The recommendation is evidence-backed and read-only until the root
session or human explicitly performs the next consequential action.

## Session hierarchy

A root session owns the human request and coordinates child sessions. A child
session has a parent, root session, role, objective, scope, permissions,
expected output, and stop conditions. A child must not broaden its assignment
or satisfy a human approval gate.

For approved-plan execution, the root claims the lease and exclusive worktree,
then directs a writer child and a later independent verifier child. Sequential
tasks share one writer child. Workspace `mode: solo` or `mode: team` does not
change that topology.

Sessions are not plans. A plan is approved intended work; a task is a unit of
that work; a session is an execution context; a worktree is the writable
isolation boundary.

Plan status is the canonical lifecycle authority. Plans use `draft`, `approved`,
and `done`; included task status is a synchronized projection using `draft`,
`ready`, and `done`. `cc-approve-plan` is the named plan-approval skill:
confirmed approval synchronizes all included tasks to `ready`. `cc-finish-plan`
is the named status-change skill: confirmed completion synchronizes them to
`done`. Ordinary whole-plan execution does not require a separate task approval
or task-selection ceremony. A one-plan request enters through `cc-run-plan`.
Connected approved unimplemented plans, or an interrupted
`.runtime/stacks/<stack-id>/` run, enter through `cc-run-stack`. Refuse to
treat a stack run as one `cc-run-plan`. The stack run freezes a runtime
`graph.yaml` and resumes from `progress.yaml`; it is not a scheduler. A resumed
session repairs stale task projections idempotently without rerunning
implementation or verification checks. Completion evidence does not change
canonical plan or task status.

## Shared context and runtime state

Product Knowledge, decisions, and plans are durable shared context. Runtime
state is local execution state and lives under .runtime/. Runtime state may
describe current progress, leases, handoffs, prompts, stack runs, and
worktrees, but it cannot override instructions, approved plans, human
decisions, or repository rules.

There is no single global current session. Multiple root and child sessions
may coexist. Each session has its own runtime record. A plan has at most one
active writing owner and a writable worktree is exclusive.

The source inbox is passive. Session entry reads the context index and the
smallest relevant Product Knowledge, not every file under `sources/`. A
source-based request identifies and reads only the selected source files,
explains why they are relevant, and records an evidence trail without copying
raw source text into accepted context.

Initialization establishes only core workspace identity: solo or team mode,
known repositories or project items, their roles, and each repository's
default active branch. Record confirmed identity in `workspace.yaml` and
`context/WORKSPACE.md`, not `context/PROJECT.md`. A workspace with zero
repositories is valid. Delivery behavior, merge policy, publication,
deployment, and external activity tools are configured separately and are not
initialization questions.

## Development loop

The normal loop is:

orient → gather evidence → draft or revise context → draft or revise plan →
human approval (`cc-approve-plan`) → claim plan execution (`cc-run-plan`) or
connected approved plans (`cc-run-stack`) →
implement → verify → review or repair within scope → handoff →
completion (`cc-finish-plan`) → optional cleanup (`cc-cleanup-runtime`).

Agents may draft, inspect, test, create isolated worktrees, and implement
approved scope. Humans control Product Knowledge acceptance, plan approval
through `cc-approve-plan`, material scope changes, publication, merge,
deployment, completion through `cc-finish-plan`, runtime cleanup through
`cc-cleanup-runtime`, and ambiguous session takeover.

## Evidence and answers

Answers must distinguish observed evidence, instructions, accepted decisions,
assumptions, proposals, blockers, and next actions. Agents must cite relevant
workspace files, plans, repository state, or verification results. They must
not claim that work is complete merely because tests pass, Git changed, or a
child session reported completion.

## Safety

Preserve dirty or uncertain work. Never reset, stash, clean, merge, deploy,
publish external work, store credentials, or overwrite another session's
runtime state without explicit authorization. Preserve `.runtime/` until a
human chooses cleanup through `cc-cleanup-runtime`.

## Archive eligibility

Archive eligibility is separate from plan status. A missing `archive.yaml`
sidecar means the historical active behavior; a latest `archived` event excludes
the unchanged bundle from ordinary discovery, approval, execution, stack
membership, and finishing while preserving explicit historical reads. Only
`cc-archive-plan`, after its own human `archive` gate, may append archive or
restore evidence. An archived done dependency remains resolved as done; an
archived draft or approved dependency stays visible and blocks. Archive never
changes status, cleans runtime, or cascades to related plans.
