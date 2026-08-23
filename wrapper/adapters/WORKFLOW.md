# Context Circuit workflow

One two-stage router handles human language. Stage A selects one bounded probe;
Stage B emits one normalized action with eligibility, authorization, reason
codes, and any exact human gate. Recommendations never mutate state.

The final action maps to one registered `context_set`. The shared packet loader
validates its exact allowlist, conditional evidence, measured bytes, and receipt
digests. Undeclared, missing, stale, or over-budget evidence stops; adapters
never broaden a packet or copy route policy.

Plans, tasks, sessions, worktrees, and runtime evidence stay distinct. Review
is read-only. Approval, execution, finish, delivery, publication, deployment,
merge, archive, takeover, and cleanup are separate human gates. Writers use
exclusive worktrees; independent verifiers are read-only. A missing child is
`host-blocked`, never self-verification.

Resume checks receipt, primary evidence, wrapper, Git, and ownership. Disabled,
denied, or unavailable providers use the filesystem-only fallback. Codex,
Claude Code, and Cursor Agent CLI share the packet loader and
provider-neutral `host_evidence`; host capability never authorizes a route,
role, lease, gate, or verification result.
