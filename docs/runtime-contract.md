# Runtime contract

Status: implementation contract for the Agent Workspace Workflow

This document defines the filesystem records used to coordinate root sessions,
child sessions, plans, leases, handoffs, and worktrees. It is execution state,
not Product Knowledge. It must never override AGENTS.md, WORKFLOW.md,
workspace.yaml, an approved plan, or a human decision.

## Layout

The runtime root is .runtime/ in the workspace that owns the plan:

~~~
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
~~~

delegation.yaml is required for a child session and absent for a root session.
handoff.md is the latest handoff owned by that session. The plan-level
handoffs/ directory is an append-only evidence view; it does not replace the
session-owned handoff.

There is deliberately no .runtime/current-session.yaml,
.runtime/current-plan.yaml, or other global pointer. Every lookup is scoped by
an explicit session, root session, plan, task, or repository key.

On fresh root entry, the coordinator creates its own session record before
delegating work or claiming a plan. The root record uses `kind: root`,
`parent_session_id: null`, and a unique session ID. If no session record is
present, the route is `orienting`; it must not be treated as an implicit active
session. Resume reads the explicit record and latest handoff before choosing
the next action.

## Safe identifiers and write rules

Identifiers used in runtime paths must match:

~~~
^[a-z0-9][a-z0-9._-]{0,63}$
~~~

Agents must reject path traversal, absolute paths, symlinks, and identifiers
that do not match this rule. Runtime writes are limited to the session's own
session directory, the plan lease it owns, and the assigned worktree.

Write rules:

1. Create parent directories before writing.
2. Write a temporary file in the same directory, flush it when the host
   supports flushing, and atomically rename it into place.
3. Never overwrite another session's handoff.md, delegation.yaml,
   lease.lock/owner.yaml, or worktree ownership record.
4. Treat a missing or partially written required record as blocked, not as
   permission to reconstruct facts.
5. Preserve old handoffs and released lease evidence until a human chooses
   runtime cleanup.

## Session record

Every session has .runtime/sessions/<session-id>/session.yaml:

~~~
schema_version: 1
session_id: sess-002
parent_session_id: sess-001
root_session_id: sess-001
kind: subagent
role: implementer
objective: Implement the assigned task
plan: plans/context-circuit-plans/0010-example
task: EXAMPLE-0001
scope: repository files listed in the task
non_goals:
  - Do not modify wrapper context
write_access: true
worktree: .runtime/worktrees/app/example
status: executing
created_at: 2026-08-14T12:00:00Z
updated_at: 2026-08-14T12:20:00Z
next_action: Run the task verification commands
blockers: []
~~~

Root sessions use kind: root, parent_session_id: null, and a workspace-relative
objective. Child sessions must have both parent and root IDs. A session's role,
scope, write_access, and worktree are authoritative for that session only.

Allowed session statuses are:

~~~
created, orienting, planning, awaiting-approval, executing, verifying,
awaiting-review, blocked, handoff, completed, failed, cancelled
~~~

Only the root session or an explicitly authorized runtime capability may update
another session's status. A worker may write its own handoff but must not mark
its plan or task done.

## Delegation packet

A child session's .runtime/sessions/<child-id>/delegation.yaml is created by
the parent before the child starts:

~~~
schema_version: 1
session_id: sess-002
parent_session_id: sess-001
root_session_id: sess-001
role: implementer
objective: Implement the assigned task
scope:
  plan: plans/context-circuit-plans/0010-example
  task: EXAMPLE-0001
  paths:
    - src/
non_goals:
  - Do not change plan scope
context_refs:
  - AGENTS.md
  - WORKFLOW.md
  - context/ARCHITECTURE.md
repository: app
worktree: .runtime/worktrees/app/example
permissions:
  write_worktree: true
  write_runtime_session: true
  write_plan: false
  write_activity: false
acceptance_criteria:
  - The assigned behavior is implemented and verified
stop_conditions:
  - A required change falls outside paths
  - A source contradiction affects the result
handoff_schema: session-handoff-v1
~~~

The child must refuse an incomplete packet. A parent may delegate multiple
children concurrently, but each write-enabled child receives a different
worktree and a non-overlapping ownership assignment unless the parent has an
explicit integration step.

For a verifier, `write_worktree`, `write_plan`, and `write_activity` must be
false. `write_runtime_session: true` permits writes only to that verifier's own
session directory and handoff; it does not permit changing worker, lease, plan,
worktree, or activity records.

## Plan lease

The live lease is .runtime/plans/<plan-id>/lease.yaml, guarded by the
exclusive-create directory .runtime/plans/<plan-id>/lease.lock/. The first
session that creates lease.lock owns the plan. It then writes owner.yaml and
lease.yaml atomically. If the lock already exists, the contender must not
modify the plan, task, worktree, or owner's runtime records.

~~~
schema_version: 1
plan: plans/context-circuit-plans/0010-example
session_id: sess-001
root_session_id: sess-001
worktree: .runtime/worktrees/app/example
status: active
acquired_at: 2026-08-14T12:00:00Z
heartbeat_at: 2026-08-14T12:20:00Z
released_at: null
stale_after_seconds: 1800
~~~

The lock directory is the ownership primitive; lease.yaml is the
human-readable state and evidence. A missing owner.yaml, a mismatched session,
an invalid timestamp, or an unexpired heartbeat makes acquisition unsafe and
therefore blocked. Only the owner may release its lock. Stale recovery
requires a human decision or an explicitly authorized takeover that records
the previous owner and reason before creating a new lease.

## Worktree ownership

Each write-enabled plan has one path:

~~~
.runtime/worktrees/<repository-key>/<plan-id>/
~~~

The plan lease and worktree record must name the same session and path. A
worker must verify its current Git root, branch, and path before writing. The
base checkout, another plan worktree, wrapper context, and external activity
state are outside the worker's write scope. A verifier may inspect a worktree
but has write_worktree: false.

## Plan/task lifecycle projection

`plan.yaml` is the canonical lifecycle record. Plan status uses exactly
`draft`, `approved`, and `done`. Included task records use a synchronized
projection with exactly `draft`, `ready`, and `done`:

| Plan status | Expected task status |
| --- | --- |
| `draft` | `draft` |
| `approved` | `ready` |
| `done` | `done` |

The coordinator reconciles every included task in one bulk operation when a
plan is approved or completed and when a session enters or resumes. The
operation is idempotent: it rewrites only a stale task `status` field, keeps
other task metadata intact, and does not rerun implementation or verification
checks. A stale projection is an observation to repair and must not block
approved-plan execution. An unknown plan status is a blocker rather than a
reason to invent a projection.

Task status is not a second human approval, execution lease, or verification
result. Approval may display the included task list, but ordinary whole-plan
execution does not require a separate task-selection ceremony.

Provider-specific activity state is optional and separate from canonical task
status. An explicitly configured adapter may maintain an opaque
`external_status` projection; core synchronization preserves it and never
uses it to decide lifecycle state. The core workflow remains usable without a
provider and does not store credentials or external activity records in
ordinary workspace state.

## Completion gate

The coordinator may prepare `.runtime/plans/<plan-id>/completion.yaml` only
after every task has a durable evidence record, the independent verifier has a
completed passing handoff, and no blocking or failed evidence remains. The
runtime completion record may use `status: blocked` or
`status: ready-for-human-status-change` and must identify the verifier and
human gate evidence.

Its minimum evidence fields are:

~~~yaml
schema_version: 1
plan: plans/context-circuit-plans/0001-example
status: ready-for-human-status-change
verifier_session_id: sess-002
verification_handoff: .runtime/sessions/sess-002/handoff.md
human_gate: status-change
canonical_status_changed: false
~~~

The completion record is evidence, not canonical intent. The canonical plan/task status remains unchanged until the human status-change gate is explicitly satisfied. A failed verifier, missing task evidence, or missing human gate keeps completion blocked.

## Handoffs

The session owner writes .runtime/sessions/<session-id>/handoff.md with:

~~~markdown
# Session handoff

- session_id: sess-002
- parent_session_id: sess-001
- root_session_id: sess-001
- status: completed
- plan: plans/context-circuit-plans/0010-example
- task: EXAMPLE-0001
- worktree: .runtime/worktrees/app/example

## Objective and scope

...

## Evidence inspected

...

## Changed files

...

## Tests and verification

...

## Decisions and assumptions

...

## Questions, blockers, and limitations

...

## Recommended next action

...
~~~

The status must be one of completed, blocked, failed, or awaiting-human-gate.
A handoff never changes canonical plan/task status. Parents copy or reference
handoff evidence; they do not rewrite the child's record.

## Recovery and concurrency

On interruption:

1. Preserve the worktree, lease, session record, and latest handoff.
2. Read the heartbeat and inspect the worktree before deciding whether work
   can resume.
3. Resume with the same session when possible.
4. If ownership is ambiguous, stop and request human takeover.
5. If takeover is authorized, create a new session and handoff that names the
   replaced session, evidence, reason, and remaining scope.

Independent plans may execute concurrently because each has a separate plan
lease and worktree. Same-plan contenders do not race; the loser is blocked or
read-only. Overlapping file changes across different plans remain an
integration risk and must be reported before merge or publication.

## Plan execution records

The plan runner may keep an optional plan prompt or task evidence files under
the owning plan runtime directory. These records are append-oriented evidence;
they do not replace `plan.yaml`, the task contracts, the session record, or the
lease. A prompt identifies the plan, root session, repository, worktree,
dependencies, route (`solo` or `delegated`), and stop conditions.

Task evidence should identify the task, owning session, assigned paths, changed
files, checks run, result, limitations, and handoff path. A task is not
complete because its worker reports success; the coordinator still requires
the plan's acceptance and independent verification evidence.

The coordinator must validate repository scope by comparing the current Git
root and branch with the lease and delegation packet before writing. For
multi-repository plans, each repository has an explicit role, path boundary,
branch expectation, and exclusive worktree. A worker cannot write the base
checkout, another repository, another plan's worktree, or runtime records it
does not own.

Same-plan contention is resolved by exclusive creation of `lease.lock/`. The
losing session becomes blocked or read-only and must not mutate the plan,
tasks, lease, worktree, or another session. Stale recovery requires explicit
human authorization and a new session that records the replaced owner and
reason; a heartbeat timeout alone is not permission to steal ownership.
