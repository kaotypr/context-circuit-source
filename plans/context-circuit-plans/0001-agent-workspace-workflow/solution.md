# Solution

- Make AGENTS.md, WORKFLOW.md, workspace.yaml, context/INDEX.md, and role instructions the normative entry contract, with explicit precedence and evidence rules.
- Define a versioned filesystem contract for .runtime/sessions, .runtime/plans, .runtime/worktrees, and durable handoff records, including safe names, required fields, lifecycle statuses, and recovery rules.
- Add host-neutral session-entry and resume instructions that route root sessions, child sessions, workers, and verifiers without introducing a user-facing command layer.
- Add a delegation packet and handoff protocol so root and child sessions can exchange bounded work, context, evidence, stop conditions, and next actions.
- Define plan execution and worktree ownership as coordinator/worker instructions: acquire a plan lease, allocate an exclusive worktree, enforce repository scope, and emit a handoff.
- Build fixture-driven acceptance scenarios that inspect filesystem state and simulate concurrent/interrupted sessions; keep the harness independent of any production command runtime.
- Remove the old Node/JavaScript layer and its command-specific adapters so the filesystem protocol is the only supported workflow.
- Keep provider publication, deployment, merge, credential storage, and activity mutation outside the workspace workflow implementation.
