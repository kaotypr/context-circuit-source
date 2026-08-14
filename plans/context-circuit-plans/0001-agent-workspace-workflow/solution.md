# Solution

- Make AGENTS.md, WORKFLOW.md, workspace.yaml, context/INDEX.md, and role instructions the normative entry contract, with explicit precedence and evidence rules.
- Define a versioned filesystem contract for .runtime/sessions, .runtime/plans, .runtime/worktrees, and durable handoff records, including safe names, required fields, lifecycle statuses, and recovery rules.
- Add host-neutral session-entry and resume instructions that route root sessions, child sessions, workers, and verifiers without introducing a user-facing command layer.
- Add a delegation packet and handoff protocol so root and child sessions can exchange bounded work, context, evidence, stop conditions, and next actions.
- Adapt the existing plan execution/worktree behavior behind the workflow contract: acquire a plan lease, allocate an exclusive worktree, enforce repository scope, and emit a handoff while retaining the legacy adapter.
- Build fixture-driven acceptance scenarios that inspect filesystem state and simulate concurrent/interrupted sessions; keep the harness independent of the production session protocol.
- Document the migration gate: the legacy Node/JavaScript layer is retained until all replacement acceptance criteria pass in a clean review, after which a separate cleanup plan may remove it.
- Keep provider publication, deployment, merge, credential storage, and activity mutation outside the workspace workflow implementation.
