# Planning

Plans define human-reviewed intended work. They are not the workspace session
itself and they do not replace runtime execution state.

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

Plans begin as `draft` and require explicit human approval. `plan.yaml` is the
canonical lifecycle record: approval changes the included task projections from
`draft` to `ready`, and completion changes them to `done`. Task status is not a second approval or execution gate. It does not provide an execution lease or
verification result. There is no user-facing `cc-run-task` workflow. There is
no run-task command.

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
