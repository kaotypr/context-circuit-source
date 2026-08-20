# Runtime contract

Status: implementation contract for the Agent Workspace Workflow

This document defines the filesystem records used to coordinate root sessions,
child sessions, plans, leases, handoffs, worktrees, and stack runs. It is
execution state, not Product Knowledge. It must never override AGENTS.md,
WORKFLOW.md, workspace.yaml, an approved plan, or a human decision. Runtime
state must not override `plan.yaml`.

## Layout

The runtime root is .runtime/ in the workspace that owns the plan:

~~~
.runtime/
  sessions/
    <session-id>/
      session.yaml
      delegation.yaml
      handoff.yaml
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
  stacks/
    <stack-id>/
      graph.yaml
      progress.yaml
      lease.lock/
        owner.yaml
      lease.yaml
  worktrees/
    <repository-key>/<plan-id>/
~~~

delegation.yaml is required for a child session and absent for a root session.
`handoff.yaml` is the preferred structured handoff owned by that session.
`handoff.md` is optional explanatory Markdown for a structured handoff and is
also retained as a read-only legacy fallback for historical sessions that have
no YAML handoff. The plan-level `handoffs/` directory is an append-only
evidence view; it does not replace the session-owned handoff.

Runtime YAML, including `handoff.yaml`, is execution state outside OKF. It does
not become Product Knowledge, a plan/task lifecycle record, or a second
authority for a lease, stack, or human gate.

There is deliberately no .runtime/current-session.yaml,
.runtime/current-plan.yaml, .runtime/current-stack.yaml, or other global
pointer. Every lookup is scoped by an explicit session, root session, plan,
task, stack, or repository key. A stack run is runtime execution state, not a
plan and not a session. There is no `plans/<repository-key>-stacks/` layout
and no durable `stack.yaml`.

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
session directory, the plan lease it owns, the stack records it owns, and the
assigned worktree.

Write rules:

1. Create parent directories before writing.
2. Write a temporary file in the same directory, flush it when the host
   supports flushing, and atomically rename it into place.
3. Never overwrite another session's handoff.yaml, handoff.md, delegation.yaml,
   lease.lock/owner.yaml, or worktree ownership record.
4. Treat a missing or partially written required record as blocked, not as
   permission to reconstruct facts.
5. Preserve old handoffs and released lease evidence until a human chooses
   runtime cleanup via `cc-cleanup-runtime`.

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
`status: ready-for-human-status-change` before the human status-change gate,
and `status: completed` after `cc-finish-plan` records that the gate was
satisfied. It must identify the verifier and human gate evidence.
`ready-for-human-status-change` remains the pre-finish evidence state.
`cc-finish-plan` is the named status-change skill.

Its minimum evidence fields are:

~~~yaml
schema_version: 1
plan: plans/context-circuit-plans/0001-example
status: ready-for-human-status-change
verifier_session_id: sess-002
verification_handoff: .runtime/sessions/sess-002/handoff.yaml
human_gate: status-change
canonical_status_changed: false
~~~

The completion record is evidence, not canonical intent. The canonical plan/task status remains unchanged until the human status-change gate is explicitly satisfied through `cc-finish-plan`. After confirmed finish, `completion.yaml` records:

~~~yaml
status: completed
human_gate: status-change
canonical_status_changed: true
~~~

A failed verifier, missing task evidence, or missing human gate keeps completion blocked. Finish releases `lease.lock/` when this coordinator owns that lease and does not delete `.runtime/`.

## Stack run records

A connected set of already-approved plans may execute as one stack run under
`.runtime/stacks/<stack-id>/`. Invoking `cc-run-stack` starts or resumes that
run. There is no stack-approval gate and no scheduler, daemon, or queue.

`graph.yaml` is frozen once for the run: members, directed edges, and leaves.
The DAG must be acyclic and IDs must resolve. Multi-parent nodes are joins.
Do not rebuild a different tree on resume.

`progress.yaml` is the resume cursor. Member states are `pending`,
`waiting-parents`, `ready`, `running`, `joining`, `implemented`, `failed`, or
`blocked`. Record `frozen_sha`, `worktree`, `base`, and `join_parents` as
members move.

The stack lease lives under that directory as `lease.lock/owner.yaml` and
`lease.yaml`. Per-plan leases remain `.runtime/plans/<plan-id>/lease.lock/`
and are claimed by the same root session that holds the stack lease. A live
foreign stack or plan lease blocks the run.

```yaml
# graph.yaml
schema_version: 1
stack_id: from-0001
repository: example-repository
members:
  - id: plan-a
    path: plans/example-repository-plans/0001-plan-a
  - id: plan-b
    path: plans/example-repository-plans/0002-plan-b
  - id: plan-c
    path: plans/example-repository-plans/0003-plan-c
edges:
  - from: plan-a
    to: plan-b
  - from: plan-a
    to: plan-c
leaves:
  - plan-b
  - plan-c
```

```yaml
# progress.yaml
schema_version: 1
stack_id: from-0001
status: running
owner_session_id: sess-a
replaced_session_id: null
members:
  plan-a:
    state: implemented
    frozen_sha: abc123
    worktree: .runtime/worktrees/example-repository/plan-a
    base:
      kind: default-branch
      commit: def456
    join_parents: []
  plan-c:
    state: waiting-parents
    frozen_sha: null
    worktree: .runtime/worktrees/example-repository/plan-c
    base:
      kind: join
    join_parents:
      - plan-x
      - plan-y
```

Member runtime records do not change `plan.yaml`. Canonical plan status stays
`draft`, `approved`, and `done`. There is no implemented plan status.

## Implemented versus done

Implemented is runtime evidence, not a canonical status:

- the writer finished;
- the independent verifier passed;
- the worktree HEAD is committed and clean;
- `.runtime/plans/<plan-id>/completion.yaml` is
  `ready-for-human-status-change`;
- `plan.yaml` remains `approved`.

After verifier pass, freeze the parent SHA on `progress.yaml`. Dependents wait
on implemented parents, not `done`. `cc-run-stack` must not write `plan.yaml`
`done` and must not run `cc-finish-plan`. `cc-finish-plan` remains the per-plan
human gate after the stack stops.

Require a local commit on the exclusive member branch before freeze. That
commit is stack-internal. It is not merge to the default branch and not a
delivery authorization.

## Stack worktrees

Each member keeps `.runtime/worktrees/<repository-key>/<plan-id>/`. Writable
worktrees stay exclusive. Overlapping paths across members are a reported
integration risk, not a shared worktree. The new worktree is a Git checkout
of that commit-ish, not a copy of the parent directory.

Worktree rules:

1. No parent: `git worktree add` from the repository default or active branch,
   matching standalone `cc-run-plan`.
2. One parent: add a new branch at the parent's frozen SHA.
3. Several parents: sort parent IDs, add from the first frozen SHA, merge the
   remaining SHAs in that order, then the writer runs. Routine joins do not
   pause for a human and do not wait for `default_branch`.

## Structured handoffs

The session owner writes `.runtime/sessions/<session-id>/handoff.yaml` as the
canonical resumable handoff. It records execution-specific deltas and
references canonical records instead of copying plan requirements or unchanged
session metadata:

```yaml
schema_version: 1
handoff_schema: session-handoff-v1
session_id: sess-002
parent_session_id: sess-001
root_session_id: sess-001
status: completed
plan: plans/app-plans/0010-checkout
task: APP-0042
worktree: .runtime/worktrees/app/example
result:
  outcome: implemented-and-verified
  changed_files:
    - src/checkout.ts
  tests:
    - command: sh test/acceptance.sh
      result: passed
blockers: []
decisions:
  - "Kept runtime YAML outside OKF and canonical plan status."
evidence:
  - kind: task-evidence
    path: .runtime/plans/0010-checkout/task-evidence/APP-0042.yaml
    result: completed
  - kind: verification
    path: .runtime/sessions/sess-003/handoff.yaml
    result: passed
next_action:
  route: review
  description: Request independent review and human status-change decision.
canonical_refs:
  plan: plans/app-plans/0010-checkout/plan.yaml
  tasks:
    - plans/app-plans/0010-checkout/tasks/APP-0042.md
  acceptance:
    - plans/app-plans/0010-checkout/plan.yaml#acceptance_criteria
  ownership:
    session: .runtime/sessions/sess-002/session.yaml
    delegation: .runtime/sessions/sess-002/delegation.yaml
    lease: .runtime/plans/0010-checkout/lease.yaml
    worktree: .runtime/worktrees/app/example
```

The required resumable fields are `schema_version`, `handoff_schema`,
`session_id`, `parent_session_id`, `root_session_id`, `status`, `plan`, `task`,
`worktree`, `result`, `blockers`, `decisions`, `evidence`, `next_action`, and
`canonical_refs`. `status` is one of `completed`, `blocked`, `failed`, or
`awaiting-human-gate`. `result`, `blockers`, `decisions`, `evidence`, and
`next_action` must describe the current execution delta; canonical references
must identify the authoritative plan, task, acceptance, evidence, and
ownership records.

On resume, readers prefer and validate `handoff.yaml`. If it is absent, they
may read a historical Markdown-only `handoff.md` using the legacy handoff
format below. If `handoff.yaml` is present but malformed or incomplete, the
route is blocked; readers must not fall back to Markdown and must not
reconstruct missing structured fields. When both files exist, YAML is
authoritative and Markdown is explanatory only. Historical runtime is not
rewritten in place.

Structured handoffs remain runtime records: they cannot override
`plan.yaml`, task projections, leases, stack `graph.yaml` or `progress.yaml`,
completion evidence, verifier isolation, or any human gate. A handoff can
recommend `cc-finish-plan`, but cannot satisfy its status-change gate.

## Handoffs

The optional Markdown view is explanatory only. A new session may write
`.runtime/sessions/<session-id>/handoff.md` as a concise human summary beside
the structured handoff, but it is not a competing canonical record:

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

The legacy Markdown status must be one of completed, blocked, failed, or
awaiting-human-gate. A handoff never changes canonical plan/task status.
Parents copy or reference structured handoff evidence; they do not rewrite the
child's record.

## Recovery and concurrency

On interruption:

1. Preserve the worktree, lease, session record, and latest handoff.
2. Read `handoff.yaml` first, or the historical Markdown-only `handoff.md`
   when no structured handoff exists. Then read the heartbeat and inspect the
   worktree before deciding whether work can resume.
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
dependencies, preferred route (`delegated`), optional `writer-count: 1` when
sequential tasks share one worktree, and stop conditions.

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

## Worktree removal and runtime cleanup

`cc-cleanup-runtime` is the named skill for deleting workspace `.runtime/`.
Cleanup is workspace-wide: confirmation must say that every session, lease,
handoff, stack run, and worktree under `.runtime/` will be removed. After
cleanup there is no stack to resume; a new run writes a new graph.

Before deletion, inspect every registered runtime worktree for uncommitted
changes, including untracked files, and for commits not present on the
worktree's upstream, or local-only commits when no upstream exists. Classify
each worktree as `clean`, `uncommitted`, `unpushed`, or both. Live
non-terminal sessions are listed because deleting `.runtime/` removes resume
records; they do not invent permission to discard Git work.

If dirty or unpushed work exists, stop with a confirmation list. Uncommitted
work is destructive if cleanup proceeds. Unpushed commits are not deleted by
`git worktree remove` because the branch remains in the parent repository;
still confirm because the checkout is going away. `--force` is allowed only
after the human confirmed discarding the listed dirty work.

Confirmed cleanup removes each registered Git worktree with `git worktree
remove`, then deletes remaining `.runtime/` records. It must not modify the
base checkout, delete Git branches, reset or stash product repositories, or
change plan or task status. Path traversal, symlinks, and identifiers outside
this contract remain rejected. Missing or already-empty `.runtime/` is a
successful no-op.
