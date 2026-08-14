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

The session must stop and ask for a human decision when the scope, acceptance
criteria, repository boundary, or safety assumptions materially change. It must
not change plan or task status merely because implementation or tests appear
complete.

The previous run-task command is transitional implementation material. The
filesystem session protocol and the workspace-entry workflow are the target
behavior.
