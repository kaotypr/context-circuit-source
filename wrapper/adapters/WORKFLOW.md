# Context Circuit entry workflow

Human language enters one two-stage router. Stage A selects exactly one bounded
probe from Tier 0. Stage B reads that probe and emits one normalized action
decision with eligibility, authorization, reason codes, and any exact human
gate. A recommendation is never a mutation.

Plans, sessions, tasks, worktrees, and runtime evidence are distinct. Review is
read-only. `Approve plan <id>` presents a card and mutates nothing;
`Confirm approval of plan <id>` is the status-only gate and does not commit
Git or start `Run approved plan`. On product-source, that confirmed turn may
then present the existing maintainer commit card; commit remains a separate
exact confirmation. Instantiated or wrapped workspaces keep
`Run approved plan <id>` as the next explicit request. Completion requires a
separate status-change confirmation. Writers use exclusive worktrees and
verifiers are independent and read-only.

On interruption, preserve state and resume only after receipt, ownership, Git,
and wrapper compatibility checks. Offline filesystem operation is complete;
external delivery, publication, deployment, merge, archive, takeover, and
destructive cleanup remain explicit human gates.

Use the owner files named by `wrapper/contracts/routes.yaml` and
`wrapper/contracts/invariants.yaml`; do not create parallel policy.

Host mapping is evidence around this workflow. A Codex native subagent, Claude
Task/subagent, or Cursor Task/subagent may receive the same bounded writer or
independent verifier packet. The packet includes host evidence and preserves
its role permissions. A missing child primitive emits `host-blocked`; it never
downgrades verification or bypasses a human gate.

Resume is host-neutral: re-read the session receipt, latest handoff, wrapper
version, Git state, and ownership before re-entering the same route. Provider
status `disabled`, `denied`, or `unavailable` uses the filesystem-only fallback.
