# Planning and lifecycle

The plan bundle is human intent; the runtime is execution evidence. The
canonical plan schema is `wrapper/contracts/schemas/plan.yaml` and the human
entry is `PLAN.md`.

When an accepted workspace receives a generic build or implementation request
without a named approved plan, the safe next action is `draft-plan`. The agent
may create only the bounded draft plan bundle; it must not create implementation
files or begin execution. Approval and the execution trigger remain separate.

## Review

Read the selected plan, accepted context, provenance, repository evidence,
dependencies, archive sidecar, ownership, and delivery boundary. Return a
Review Card with outcome, summary, exact approval scope, non-effects,
task/dependency table, acceptance → task → verification mapping, evidence,
risks, contradictions, and no more than three focused human decisions. Review
is read-only.

## Approval

`Approve plan <id>` shows a current session-bound Approval Card and changes
nothing. The card states that nothing has changed yet, names the status-only
projection that confirmation would apply, and asks for
`Confirm approval of plan <id>`. It does not imply that confirmation commits
Git or starts `Run approved plan`.

On that exact confirmation only, call `cc_transition_plan_status` once for the
plan and all included task projections (`plan.yaml` draft → approved, tasks
draft → ready). Preserve Markdown bodies and file endings. Do not hand-edit
each task independently. Do not claim a lease, create a worktree, start
children, commit Git, or start execution.

Approval is complete after that status-only transition; execution has not
started. On an instantiated or wrapped workspace, the separate next request is
`Run approved plan <id>`. On `product-source`, if the dirty set is exactly that
approval projection, present the existing `commit-approved-plan` card in the
same session. Do not commit on this turn. Unrelated dirty files remain
`DIRTY_BASE_BLOCKED`.

## Execution

`Run approved plan <id>` is a separate explicit request. Preflight status,
dependencies, archive eligibility, dirty base, wrapper compatibility, ownership,
and child capability. On success, the root claims an exclusive lease, creates a
worktree and receipt, delegates one writer, then one independent verifier.
Sequential tasks share the writer and worktree. Connected plans freeze a DAG
and progress cursor; they do not create a hidden plan or scheduler.

For the maintainer product-source checkout, approval itself changes only the
plan status projection and therefore creates a predictable dirty delta. After
confirmed approval, if that delta is the only dirty change, present the
existing maintainer commit card immediately. If someone later runs before that
commit, execution still reports `MAINTAINER_APPROVAL_COMMIT_REQUIRED`. It never
treats arbitrary dirty source as safe, never treats that class as an execution
exemption, and never commits automatically.

## Completion

When every task has evidence, the verifier passes, and the worktree is clean
and committed, write runtime `completion.yaml` with
`ready-for-human-status-change`. Show a Finish Card. On exact confirmation only,
update approved → done and task projections ready → done, release the owned
lease, and preserve runtime/worktrees/branches for later delivery or cleanup.

Implemented is runtime evidence; Done is human-confirmed canonical status.

## Archive

`archive.yaml` is optional and append-only. A latest `archived` event removes a
plan from ordinary selection without changing status or deleting evidence.
Active owners, dirty/unpushed work, stack membership, and unresolved unfinished
dependencies block archive/restore. Restore is separate from approval and run.
