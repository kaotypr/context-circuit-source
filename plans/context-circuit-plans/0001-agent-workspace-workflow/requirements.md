# Requirements

- The workspace entry point is an AI-agent session that reads authoritative instructions and durable context before deciding what to do.
- Root sessions and subagent sessions have explicit identities, parent/root relationships, roles, scopes, and lifecycle records.
- A root session can coordinate multiple concurrent plans without a global current-plan or current-session file.
- Every child session receives a structured delegation packet with objective, scope, context references, permissions, acceptance criteria, stop conditions, and handoff schema.
- Runtime state is filesystem-based, inspectable, append-safe where practical, and separated from durable product knowledge and plan intent.
- Plan leases prevent two writing sessions from concurrently owning the same plan while allowing unrelated plans to proceed concurrently.
- Repository workers use an exclusive plan-scoped branch and worktree; the base checkout and unrelated worktrees remain untouched.
- Workers return a durable handoff containing status, evidence, changed paths, tests, unresolved risks, and the next recommended action.
- Verifiers are independent and read-only; they can evaluate a worker handoff and worktree without mutating plan or runtime state.
- Human gates remain explicit for approval, delegation where required, conflict resolution, and completion.
- The workflow must recover from interruption, stale leases, missing handoffs, contradictory context, and failed verification without inventing state.
- The old Node/JavaScript command layer remains intact and usable during migration and is removed only after a separately reviewed cleanup plan follows a passing replacement acceptance suite.
