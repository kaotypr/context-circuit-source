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
`draft` to `ready`, and completion changes them to `done`. Task status is not a
second approval gate and does not provide an execution lease or verification
result.

On approval, completion, and session entry or resume, the coordinator
reconciles every included task to the projection expected by the plan. The
operation is bulk, idempotent, preserves task metadata such as
`external_status`, and does not rerun implementation or verification checks.
Runtime session state may say that a plan is being executed, blocked, or
awaiting review, but it must never silently change the canonical plan status.

For large projects, prefer several coherent plans by domain or repository
boundary rather than one unbounded plan. Multiple approved plans may execute
concurrently when their worktrees and ownership are distinct.
