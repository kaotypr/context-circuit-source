# Planning and lifecycle

The plan bundle is human intent; the runtime is execution evidence. The
canonical plan schema is `wrapper/contracts/schemas/plan.yaml` and the human
entry is `PLAN.md`.

## Review

Read the selected plan, accepted context, provenance, repository evidence,
dependencies, archive sidecar, ownership, and delivery boundary. Return a
Review Card with outcome, summary, exact approval scope, non-effects,
task/dependency table, acceptance → task → verification mapping, evidence,
risks, contradictions, and no more than three focused human decisions. Review
is read-only.

## Approval

Show an Approval Card tied to the named plan and current session. On exact
confirmation only, update `plan.yaml` draft → approved and reconcile included
task projections draft → ready atomically/idempotently. Do not claim a lease,
create a worktree, start children, or perform Git work.

## Execution

`Run approved plan <id>` is a separate explicit request. Preflight status,
dependencies, archive eligibility, dirty base, wrapper compatibility, ownership,
and child capability. On success, the root claims an exclusive lease, creates a
worktree and receipt, delegates one writer, then one independent verifier.
Sequential tasks share the writer and worktree. Connected plans freeze a DAG
and progress cursor; they do not create a hidden plan or scheduler.

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
