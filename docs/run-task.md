# Plan execution in the agent workspace

Plan execution is one stage of the workspace session workflow, not a command
the user invokes.

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

Before requesting completion, the coordinator records task evidence, a passing
independent verification handoff, and the required human status-change gate in
the plan runtime directory. Missing evidence or a failed verifier blocks the
completion record; the canonical plan/task status remains unchanged.

There is no run-task command. The coordinator and repository-worker
instructions define the execution behavior, while the filesystem lease,
worktree, prompt, and handoff records make the state resumable and inspectable.

See docs/runtime-contract.md for record fields, atomic acquisition, recovery,
and handoff rules.
