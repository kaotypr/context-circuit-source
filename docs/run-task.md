# Plan execution in the agent workspace

Plan execution is one stage of the workspace session workflow, not the primary
user interface.

A root session executes an approved plan by:

1. Confirming the plan is approved and has unfinished work.
2. Acquiring the plan lease.
3. Creating or reusing the plan's exclusive repository worktree.
4. Writing a session-scoped delegation packet or plan prompt.
5. Delegating bounded implementation and verification work.
6. Continuing through dependency-ready tasks within approved scope.
7. Recording tests, findings, blockers, and the next handoff.

A plan has at most one active writing owner. Read-only research and verification
sessions may inspect its worktree, but two writing sessions must never share it.

The runtime contract uses an exclusive plan lease and records the owning
session, root session, worktree, heartbeat, and release state. A same-plan
contender becomes blocked; it must not modify the plan, task, lease, or worker
worktree. Different plans may proceed concurrently.

The session must stop and ask for a human decision when the scope, acceptance
criteria, repository boundary, or safety assumptions materially change. It must
not change plan or task status merely because implementation or tests appear
complete.

The previous run-task command is transitional implementation material. The
filesystem session protocol and the workspace-entry workflow are the target
behavior. During migration, the command accepts an explicit session identity
for internal adapters and acquires the same lease used by the session
workflow. When no identity is supplied, it uses a clearly labeled
legacy-run-task compatibility identity; this path is retained only until the
replacement acceptance suite passes.

See docs/runtime-contract.md for record fields, atomic acquisition, recovery,
and handoff rules.
