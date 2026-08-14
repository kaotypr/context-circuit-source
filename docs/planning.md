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

Plans and tasks begin as draft. Human approval is explicit. Runtime session
state may say that a plan is being executed, blocked, or awaiting review, but it
must never silently change the plan or task status.

For large projects, prefer several coherent plans by domain or repository
boundary rather than one unbounded plan. Multiple approved plans may execute
concurrently when their worktrees and ownership are distinct.
