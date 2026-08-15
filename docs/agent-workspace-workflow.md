# Agent Workspace Workflow

Status: agreed workflow contract; pure filesystem implementation

The filesystem record details for this contract live in
runtime-contract.md. This document defines behavior; the runtime contract
defines record shape and ownership mechanics.

## Purpose

Context Circuit is an AI-agent workspace. Its primary job is to give an AI agent a reliable way to enter a project, understand the workspace, coordinate work, continue across sessions, and report what is known, changed, blocked, or still uncertain.

The human-facing interface is an agent session such as:

> Start or resume work in this workspace.

The user does not invoke a bundled CLI. Hosts provide an agent session that
reads these instructions and uses the filesystem records directly.

This document is the normative behavior contract for the workspace.

## Scope and non-goals

Context Circuit owns:

- workspace instructions and Product Knowledge;
- plans, task dependencies, and human approval gates;
- the session and subagent coordination protocol;
- runtime session state, handoffs, leases, prompts, and worktree references;
- safe preparation of isolated repository worktrees;
- evidence-backed agent entry, continuation, and handoff behavior.

The host agent owns model execution and may provide subagent primitives. Product repositories own source code and repository-local conventions.

Context Circuit does not implicitly:

- approve plans or scope changes;
- merge, publish, deploy, or close external work;
- infer completion from tests, Git state, publication, or agent output;
- discard dirty work or repair an uncertain repository automatically;
- make runtime state authoritative over Product Knowledge, plans, human decisions, or repository instructions.

## Workspace model

The workspace contains several distinct kinds of information:

| Area | Responsibility | Authority |
| --- | --- | --- |
| `AGENTS.md` | Agent safety, behavior, and instruction precedence | Normative wrapper instruction |
| `WORKFLOW.md` | Workflow states, gates, and coordination rules | Normative workflow instruction |
| `workspace.yaml` | Repositories, branches, mode, and configuration | Machine-readable configuration |
| `context/` | Product, domain, architecture, conventions, and decisions | Durable project knowledge |
| `sources/` | Raw inputs and authored Idea Brief/PRD artifacts | User/team-owned; passive except request-scoped reads |
| `plans/` | Proposed and approved intended work | Human-reviewed work definition |
| `.runtime/` | Active sessions, subagents, leases, prompts, handoffs, and worktrees | Current execution state only |
| Product repository | Code and repository-local behavior | Repository authority for code |

Runtime state describes what is happening now. It cannot override an approved plan, change a human gate, or redefine the product.

## Authority and evidence

The agent must distinguish instruction, fact, decision, assumption, proposal, and runtime observation.

The effective authority order is:

1. Host and system instructions.
2. Wrapper safety instructions in `AGENTS.md`.
3. Workflow rules in `WORKFLOW.md`.
4. Repository-local instructions for code and repository behavior.
5. Human-approved plans and recorded decisions.
6. Source-cited Product Knowledge and source documents.
7. Runtime state and prior handoffs.
8. Agent assumptions and suggestions.

This is a domain-specific precedence model. Wrapper instructions govern workspace safety and orchestration; repository instructions govern repository code and local conventions. When sources contradict one another, the agent must surface the contradiction and stop before making a consequential assumption.

Every material answer should separate:

- observed state and evidence;
- applicable instructions and decisions;
- assumptions or inferences;
- open questions and blockers;
- proposed next action.

The agent must cite workspace-relative files, plan references, command results, or other available evidence when making claims about the workspace.

## Sessions and subagents

A session is one AI execution context. A session may be the root session started for a human request or a child session spawned by another agent.

```text
Workspace
└── Root session
    ├── Research subagent
    ├── Planning subagent
    ├── Implementation subagent
    └── Verification subagent
```

Sessions form a tree. Every child records its `parent_session_id` and `root_session_id`. A session may delegate further only when its parent scope permits it.

### Root session

The root session:

- owns the human request and overall objective;
- reads the workspace entry context;
- decides whether to discover, plan, execute, verify, resume, or ask for clarification;
- owns coordination of child sessions;
- consolidates child findings and handoffs;
- presents human gates and the final plan-level handoff.

A root session may coordinate multiple plans, but each plan still has an independent execution owner and worktree.

For approved-plan execution, the root claims the lease and exclusive worktree,
then directs a writer child (`write_worktree: true`) and a later independent
verifier child (`write_worktree: false`). Sequential tasks share one writer
child. Independent plans get separate children and worktrees; overlapping
paths are reported before merge or publication. The same writer-child and
verifier-child topology applies when `workspace.yaml` is `mode: solo` and
when it is `mode: team`.

### Subagent session

A subagent session:

- has one explicit parent, objective, scope, role, and expected output;
- reads only the relevant context needed for that scope plus required instructions;
- returns findings, changes, tests, questions, blockers, and recommendations to its parent;
- may be read-only or write-enabled;
- must not approve plans, change canonical statuses, or modify another session's state.

A verifier is read-only for implementation, plan, lease, worktree, and activity
state. It may write only its own session-scoped handoff so its pass/fail result
is durable and resumable.

Typical roles include `researcher`, `planner`, `implementer`, `verifier`, `reviewer`, and `integrator`.

## Work, plans, tasks, and worktrees

Sessions are execution contexts, not work items.

- A plan is the human-reviewed definition of intended work.
- A task is a scoped unit inside a plan.
- A worktree is the writable isolation boundary for code changes.
- A session may own a plan execution, a task, a research question, or a verification scope.

There may be many active plans and sessions in one workspace.

The default ownership rules are:

- one active writing owner per plan execution;
- one writable session per worktree;
- read-only research and verification sessions may inspect a worktree concurrently;
- writing subagents working in parallel must use separate worktrees;
- integration into a shared plan worktree is an explicit parent or human action;
- different plans may run concurrently for the same repository when their worktrees are separate;
- overlapping changes between plans are a risk to report, not a reason to reset or overwrite work.

Plan status remains exactly `draft`, `approved`, and `done`, and `plan.yaml` is
the canonical lifecycle record. Task status is a synchronized projection with
the values `draft`, `ready`, and `done`:

```text
plan draft     → tasks draft
plan approved  → tasks ready
plan done      → tasks done
```

Session status and runtime execution state remain separate from plan and task
status. A task projection is not a second approval gate, execution lease, or
verification result.

## Runtime state

Runtime state is private, local, resumable workspace state. It is preserved until a human chooses cleanup via `cc-cleanup-runtime` and is not treated as Product Knowledge.

The target layout is:

```text
.runtime/
  sessions/
    <session-id>/
      session.yaml
      delegation.yaml
      handoff.md
  plans/
    <plan-id>/
      lease.lock/
        owner.yaml
      lease.yaml
      prompt.md
      completion.yaml
      handoffs/
        <session-id>-<sequence>.md
  worktrees/
    <repository-key>/<plan-id>/
```

There is no single global `current-session.yaml`. Multiple sessions must be discoverable without overwriting one another.

A session record should include, at minimum:

```yaml
session_id: sess-002
parent_session_id: sess-001
root_session_id: sess-001
kind: subagent
role: verifier
objective: Verify checkout validation
plan: plans/app-plans/0010-checkout
task: APP-0042
scope: read-only
write_access: false
worktree: .runtime/worktrees/app/0010-checkout
status: executing
created_at: 2026-08-14T12:00:00Z
updated_at: 2026-08-14T12:20:00Z
next_action: Report verification findings to sess-001
blockers: []
```

A plan lease should identify the owning session, plan, worktree, acquisition time, heartbeat, and release or stale status. Lease acquisition must be atomic enough to prevent two writing sessions from silently claiming the same plan.

Runtime writes must be session-scoped and recoverable. A session must not
rewrite another session's handoff, lease, prompt, or worktree metadata.

See runtime-contract.md for the versioned record fields, safe identifier rules,
atomic lease convention, and handoff format.

## Workspace entry workflow

Every root or child session begins by identifying its session role and loading the minimum authoritative context.

### Root entry

The root session reads, in order:

1. Host and workspace instructions.
2. `AGENTS.md`.
3. `WORKFLOW.md`.
4. `workspace.yaml`.
5. The context index and relevant Product Knowledge.
6. Active session records and handoffs.
7. Active plans and their dependencies.
8. Repository-local instructions and observed Git state for the selected domain.

It then routes the request:

| Situation | Route |
| --- | --- |
| No clear objective | Orient and ask focused questions |
| New PRD or source | Gather evidence and draft Product Knowledge |
| Useful context but no plan | Draft a plan and request approval |
| Approved plan without an owner | Claim the plan and execute |
| Existing owner or child session | Resume or coordinate with that session |
| Verification failure | Continue within approved scope or report a blocker |
| Scope change | Pause and request human decision or reapproval |
| Completed implementation | Verify, hand off, and request human review |

The root agent explains the selected route before taking consequential action.
On a fresh route it creates its own root session record first. On a resume route
it preserves the existing record and reads the latest handoff; it never
reconstructs ownership from conversation history or a global pointer.

Choosing the next action is part of this root entry, not a separate command.
Inspect approved plans with unfinished work, declared dependencies, active
leases and session ownership, source freshness, repository cleanliness,
worktrees, blockers, and pending human gates. Recommend or claim only work
that is dependency-ready, explicitly scoped, and not already owned by another
writing session. When no work is executable, explain whether the session needs
context, a draft plan, human approval, a review, or a decision about a
blocker.

### Child entry

A child session first reads its own session record, parent handoff, delegated scope, required workspace instructions, relevant Product Knowledge, and the exact plan or task it was given.

It must not broaden its scope silently. If the assigned work is contradictory, incomplete, or unsafe, it returns a question or blocker to the parent instead of inventing authority.

## Development loop

The normal loop is:

```text
Orient
  → Gather evidence
  → Draft or revise context
  → Draft or revise plan
  → Human approval (`cc-approve-plan`)
  → Claim plan execution (`cc-run-plan`)
  → Implement in isolated worktree
  → Verify
  → Repair within approved scope when appropriate
  → Human review
  → Record handoff
  → Completion (`cc-finish-plan`)
  → Optional cleanup (`cc-cleanup-runtime`)
```

An agent may iterate between implementation and verification without asking for per-task review when the work remains within the approved plan. It must pause when the plan scope, acceptance criteria, repository boundary, or safety assumptions materially change.

### Plan-driven task lifecycle

The coordinator synchronizes all included task records in one idempotent
operation when a plan is approved or completed. A resumed session performs the
same reconciliation before reporting its next action, repairing stale or
partially synchronized task status without rerunning implementation or
verification checks. The coordinator rewrites only the task `status` field and
preserves other metadata, including an optional provider-owned
`external_status` projection.

An approved plan may execute when task metadata is stale; stale projection is
an observation to repair, not an execution blocker. The coordinator must stop
on an unknown plan status rather than inventing a task status. Plan completion
still requires plan-level implementation evidence, independent verification,
and the human status-change gate.

## Session lifecycle

Session lifecycle is runtime state and is distinct from plan and task status.

Recommended session states are:

- `created`
- `orienting`
- `planning`
- `awaiting-approval`
- `executing`
- `verifying`
- `awaiting-review`
- `blocked`
- `handoff`
- `completed`
- `failed`
- `cancelled`

Only an authorized parent or human may change a child session's delegation or cancel it. A stale session must not be taken over silently; the new session records the takeover or replacement relationship and preserves the old handoff.

## Human gates

Agents may read, inspect, draft, test, create isolated worktrees, implement approved scope, and write runtime handoffs.

Human approval is required for:

- accepting Product Knowledge changes as canonical;
- approving a plan through `cc-approve-plan`;
- materially changing plan scope or acceptance criteria;
- publishing or merging external work;
- deployment or other consequential external actions;
- marking tasks or plans done through `cc-finish-plan`;
- deleting `.runtime/` through `cc-cleanup-runtime`;
- taking over a live or stale session when ownership is ambiguous.

Approval displays the included task list for visibility, but ordinary
whole-plan execution does not require separate task selection or task approval.

Subagents never satisfy a human gate on behalf of the root session or human.

### Completion evidence

Before requesting a plan or task status change, the coordinator records durable
evidence for every task, a passing independent verification handoff, and the
required human gate in the plan runtime directory. It may record that the plan
is ready for human status change (`ready-for-human-status-change`), but it does
not change canonical plan/task status itself. When that evidence is ready, ask
for `cc-finish-plan`. Missing evidence, a failed verifier, or an unresolved
blocker keeps the plan blocked. Finish does not delete `.runtime/` or start
`cc-run-plan`.

## Handoffs and answers

Every subagent handoff must include:

- session and parent identifiers;
- objective and delegated scope;
- evidence inspected;
- decisions or assumptions made;
- files changed, if any;
- tests or verification performed and results;
- open questions and blockers;
- recommended next action.

Every root-session response should identify:

- the session and plan context;
- current state and evidence;
- work performed or deliberately not performed;
- tests and verification;
- blockers and human decisions needed;
- the next safe action.

The agent must say when information is unknown. It must not claim that a plan or task is done merely because an implementation exists, tests pass, or a child session reports completion.

## Recovery and concurrency

When a session stops unexpectedly:

1. Preserve its worktree and runtime files.
2. Detect the last heartbeat and recorded handoff.
3. Do not reset, clean, or overwrite the worktree.
4. Offer resume by the same session or explicit takeover by a new session.
5. Preserve the parent-child relationship and explain any ownership change.

When two sessions request the same plan, only the session holding the valid lease may write. The other session becomes `blocked` or switches to an explicitly read-only role.

When different plans modify overlapping areas, both may continue in isolated worktrees, but the root session must surface the integration risk before merge or publication.

## Acceptance scenarios

The workflow is correct only when it supports these scenarios:

1. A root session enters a new workspace, reads a PRD, drafts Product Knowledge, and proposes a plan without silently approving it.
2. A root session resumes from a prior handoff without relying on conversation history.
3. A root session spawns research and verification subagents with explicit scopes and receives structured handoffs.
4. Two root sessions execute different plans concurrently in separate worktrees.
5. Two sessions attempt the same plan and only one obtains the writing lease.
6. A subagent crashes; a later session resumes or explicitly takes over without losing work.
7. Verification fails; the agent repairs within scope or reports a blocker rather than changing the plan silently.
8. A PRD or Product Knowledge source changes during execution; the agent warns and requests a human-reviewed refresh.
9. A session response clearly separates evidence, assumptions, decisions, blockers, and next actions.

The command-independent acceptance suite must replay these cases from isolated
temporary filesystem and Git fixtures. It must inspect ownership artifacts
directly, prove that a losing lease contender and a verifier cannot mutate
another owner's state, and prove that completion remains blocked until the
human status-change gate is present.

The acceptance suite also verifies the plan-driven task projection: draft plans
leave draft tasks unchanged, approval and completion synchronize all tasks in
bulk, repeated synchronization is idempotent, resume repairs stale metadata
without rerunning checks, and optional `external_status` data is preserved
without requiring a provider.

## Implementation boundary

The primary host workflow exposes one workspace-entry capability for starting or
resuming a root or child session. It must not depend on a command runtime,
generated bundle, package manager, or global pointer.

The sole workspace skill is a thin entry instruction around this contract.
Product repositories remain isolated from wrapper runtime state, and repository
workers modify only their assigned worktree.

## Plan execution capability

`cc-run-plan` is the sole standard execution entry for an approved plan. It is
an agent-session capability backed by the filesystem contract, not a command
runtime. Task status is not a second approval or execution gate. A plan has at
most one active writing owner and an exclusive worktree. Completion evidence
does not change canonical plan or task status.

Before execution, the root coordinator reads the canonical plan and task
contracts, verifies dependencies, repairs stale task projections, checks the
repository and worktree, and claims the plan lease. An unapproved, unknown,
contradictory, or ownership-conflicted plan is blocked.

After preflight, the expected execution records are a writer child packet and
a later independent verifier child packet. Sequential tasks share one writer
child; they are not a reason to skip children. Every writing owner has an
exclusive worktree. A verifier is independent and read-only, and may write
only its own session handoff. Task ordering follows declared dependencies;
task selection is internal coordination rather than a second human approval
gate. Spawn children through the host child-session primitive. If the host
cannot spawn a child, report the missing host primitive to the human. Do not
quietly skip children.

Interruption and recovery preserve the session record, lease, worktree, dirty
state, questions, blockers, and latest handoff. A stale or missing record
requires a visible recovery decision. A takeover, when explicitly authorized,
names the replaced session, reason, and preserved evidence.

Completion requires durable evidence for every task, an independent passing
verification handoff, and the human `status-change` gate through
`cc-finish-plan`. The coordinator may prepare completion evidence, but tests,
Git state, or a verifier never change canonical plan or task status and never
authorize merge, publication, or deployment.
