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
2. Read context/INDEX.md and only the relevant Product Knowledge.
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

## Session hierarchy

A root session owns the human request and coordinates child sessions. A child
session has a parent, root session, role, objective, scope, permissions,
expected output, and stop conditions. A child must not broaden its assignment
or satisfy a human approval gate.

Sessions are not plans. A plan is approved intended work; a task is a unit of
that work; a session is an execution context; a worktree is the writable
isolation boundary.

## Shared context and runtime state

Product Knowledge, decisions, and plans are durable shared context. Runtime
state is local execution state and lives under .runtime/. Runtime state may
describe current progress, leases, handoffs, prompts, and worktrees, but it
cannot override instructions, approved plans, human decisions, or repository
rules.

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
default active branch. A workspace with zero repositories is valid. Delivery
behavior, merge policy, publication, deployment, and external activity tools
are configured separately and are not initialization questions.

## Development loop

The normal loop is:

orient → gather evidence → draft or revise context → draft or revise plan →
human approval → claim plan execution → implement → verify → review or repair
within scope → handoff → resume, continue, or close.

Agents may draft, inspect, test, create isolated worktrees, and implement
approved scope. Humans control Product Knowledge acceptance, plan approval,
material scope changes, publication, merge, deployment, completion, and
ambiguous session takeover.

## Evidence and answers

Answers must distinguish observed evidence, instructions, accepted decisions,
assumptions, proposals, blockers, and next actions. Agents must cite relevant
workspace files, plans, repository state, or verification results. They must
not claim that work is complete merely because tests pass, Git changed, or a
child session reported completion.

## Safety

Preserve dirty or uncertain work. Never reset, stash, clean, merge, deploy,
publish external work, store credentials, or overwrite another session's
runtime state without explicit authorization.
