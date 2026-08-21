# Context Circuit entry workflow

Human language enters one two-stage router. Stage A selects exactly one bounded
probe from Tier 0. Stage B reads that probe and emits one normalized action
decision with eligibility, authorization, reason codes, and any exact human
gate. A recommendation is never a mutation.

Plans, sessions, tasks, worktrees, and runtime evidence are distinct. Review is
read-only; approval changes status only; execution requires a separate explicit
request; completion requires a separate status-change confirmation. Writers use
exclusive worktrees and verifiers are independent and read-only.

On interruption, preserve state and resume only after receipt, ownership, Git,
and wrapper compatibility checks. Offline filesystem operation is complete;
external delivery, publication, deployment, merge, archive, takeover, and
destructive cleanup remain explicit human gates.

Use the owner files named by `wrapper/contracts/routes.yaml` and
`wrapper/contracts/invariants.yaml`; do not create parallel policy.
