# Planning

Plans define human-reviewed intended work. They are not the workspace session
itself and they do not replace runtime execution state. Every-session Product
Knowledge includes `context/WORKSPACE.md` for workspace identity and
`context/PROJECT.md` for the project being built.

A plan should have:

- one repository domain where code changes are expected;
- a clear objective and source;
- implementation scope and non-goals;
- Product Knowledge references;
- dependencies and connections;
- acceptance criteria;
- test scope and verification commands;
- risks, assumptions, and open questions;
- tasks with explicit dependencies and bounded scopes.

Plans begin as `draft` and require explicit human approval through
`cc-approve-plan`. `plan.yaml` is the canonical lifecycle record: approval
changes the included task projections from `draft` to `ready`, and completion
through `cc-finish-plan` changes them to `done`. Task status is not a second approval or execution gate. It does not provide an execution lease or
verification result. Approved work enters through `cc-run-plan`.

`cc-approve-plan` is the named plan-approval skill. After explicit confirmation
in the current session it changes `plan.yaml` from `draft` to `approved` and
reconciles included tasks to `ready`. Approval does not claim a lease, create a
worktree, or start `cc-run-plan`. A prior `cc-review-plan` run is not required.
Already-approved plans are reported as approved and may enter `cc-run-plan`
without a second status rewrite.

A plan has at most one active writing owner. The runtime contract uses an
exclusive plan lease and an exclusive worktree. Two writing sessions must never
share that worktree. Completion evidence, tests, Git state, or a verifier
handoff do not change canonical plan or task status. Missing evidence or a
failed verifier blocks the completion record; the canonical status remains
unchanged.

On approval, completion, and session entry or resume, the coordinator
reconciles every included task to the projection expected by the plan. The
operation is bulk, idempotent, preserves task metadata such as
`external_status`, and does not rerun implementation or verification checks.
Runtime session state may say that a plan is being executed, blocked, or
awaiting review, but it must never silently change the canonical plan status.

For large projects, prefer several coherent plans by domain or repository
boundary rather than one unbounded plan. Multiple approved plans may execute
concurrently when their worktrees and ownership are distinct.

Approved execution enters through `cc-run-plan`. The root directs a writer
child and an independent verifier child. Sequential tasks share one writer
child and one worktree; independent plans use separate children and
worktrees. Overlapping paths are reported before merge or publication.
