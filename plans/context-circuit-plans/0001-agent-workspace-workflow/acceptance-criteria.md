# Acceptance criteria

- A fresh workspace-entry session reads AGENTS.md, WORKFLOW.md, workspace.yaml, context/INDEX.md, and the relevant durable context, then reports whether it is orienting, resuming, or blocked.
- A root session can create or resume its own runtime session record and can coordinate two independent plans without either plan overwriting the other’s runtime state.
- A child session can be started only from a delegation packet that records root session, parent session, role, objective, bounded scope, context references, worktree, permissions, acceptance criteria, stop conditions, and handoff schema.
- Two sessions attempting to acquire the same plan lease produce one owner and one explicit contention result; the losing session does not modify plan, task, or worker state.
- Two sessions working on different plans receive distinct plan-scoped worktrees and branches, and no worker writes to the base checkout or another worker’s worktree.
- An interrupted worker can be resumed from its session record and handoff without relying on an implicit global current-session or current-plan pointer.
- A worker that reaches a scope boundary, missing dependency, contradictory source, or failed check records a blocked or failed handoff and stops rather than silently broadening scope.
- An independent verifier can reproduce acceptance checks from the worktree and handoff, record pass/fail evidence, and leave runtime and activity state unchanged.
- Plan completion is refused until task evidence, verification evidence, and required human gates are present; a failed verifier cannot be represented as done.
- The acceptance suite exercises root entry, child delegation, concurrent plans, lease contention, worktree isolation, interruption/recovery, scope safety, contradiction handling, verification failure, and completion gating.
- The replacement workflow is usable through agent instructions and filesystem conventions without requiring the user to invoke node, cc.mjs, or another command as the primary entry point.
- No Node/JavaScript command layer, generated bundle, command-specific host adapter, or package-manager entry point remains in the replacement workspace.
