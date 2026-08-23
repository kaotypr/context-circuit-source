# Planning and lifecycle

The plan bundle is human intent; the runtime is execution evidence. The
canonical plan schema is `wrapper/contracts/schemas/plan.yaml` and the human
entry is `PLAN.md`.

When an accepted workspace receives a generic build or implementation request
without a named approved plan, the safe next action is `draft-plan`. The agent
may create only the bounded draft plan bundle; it must not create implementation
files or begin execution. Approval and the execution trigger remain separate.

## Review

`Review plan <id>`, `Walk me through plan <id>`, and equivalent named-plan
phrasing route to `review-plan` through the shared `cc-plan` discovery
adapter. A review request with no usable plan id routes to `clarify-target`
instead of guessing a bundle, the same way an unnamed approval request does.

Read the selected plan, accepted context, provenance, repository evidence,
dependencies, archive sidecar, ownership, and delivery boundary. Return a
Review Card with outcome, summary, exact approval scope, non-effects,
task/dependency table, acceptance → task → verification mapping, evidence,
risks, contradictions, and no more than three focused human decisions. Review
is read-only; see `docs/plan-review.md` for the exact card and the optional
host question-prompt procedure that presents those same decisions through the
current host's native question UI when one is present, falling back to the
card text otherwise. Approval confirmation itself remains owned by `cc-gates`
and `docs/gates.md`.

## Approval

`Approve plan <id>` shows a session-bound card and changes nothing. Card
wording is in `docs/gates.md`. On exact `Confirm approval of plan <id>` only,
call `cc_transition_plan_status` once. That transition is status-only. The
next gate is in `docs/gates.md`.

## Execution

`Run approved plan <id>` is a separate explicit request. Preflight status,
dependencies, archive eligibility, dirty base, wrapper compatibility, ownership,
and child capability. On success, the root claims an exclusive lease, creates a
worktree and receipt, delegates one writer, then one independent verifier.
Sequential tasks share the writer and worktree. Connected plans freeze a DAG
and progress cursor; they do not create a hidden plan or scheduler.

If someone runs before the product-source maintainer commit, execution reports
`MAINTAINER_APPROVAL_COMMIT_REQUIRED`. That class is not an execution
exemption. The commit card is in `docs/gates.md`. It never treats arbitrary
dirty source as safe and never commits automatically.

## Completion

When every task has evidence, the verifier passes, and the worktree is clean
and committed, write runtime `completion.yaml` with
`ready-for-human-status-change`. Show a Finish Card. On exact confirmation only,
update approved → done and task projections ready → done, release the owned
lease, and preserve runtime/worktrees/branches for later delivery or cleanup.

Only an independent `passed` outcome satisfies completion. Exact-match
evidence-layer comparison is owned by `wrapper/contracts/schemas/plan.yaml`.
Writer handoff mappings are claims. `failed`, `blocked`, and `waived` remain
non-passing and are never a finish authorization.

Implemented is runtime evidence; Done is human-confirmed canonical status.

## Archive

`archive.yaml` is optional and append-only. A latest `archived` event removes a
plan from ordinary selection without changing status or deleting evidence.
Active owners, dirty/unpushed work, stack membership, and unresolved unfinished
dependencies block archive/restore. Restore is separate from approval and run.
