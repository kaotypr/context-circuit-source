# Context Circuit workflow

One two-stage router selects a bounded probe and normalized action; recommendations
never mutate state. Action selects one registered `context_set`; the shared packet loader enforces exact allowlists, selected evidence, byte budgets, and
receipts. Undeclared, missing, stale, or over-budget evidence stops and adapters
never broaden packets or copy route policy.

Plans, tasks, sessions, worktrees, and runtime evidence stay distinct. Review is
read-only; approval, execution, finish, delivery, publication, deployment,
merge, archive, takeover, and cleanup remain separate gates. Writers use
exclusive worktrees and verifiers are independent/read-only. Resume checks
receipt, primary evidence, wrapper, Git, and ownership; disabled, denied, or
unavailable providers remain filesystem-only. Hosts share the packet loader and
provider-neutral `host_evidence`; capability never authorizes route, role,
lease, gate, or verification. Committed graph is authoritative; children
consume delegation/child-start/receipt/handoff via
`cc_runtime_launch_projection`; `cc_runtime_graph_authoritative` blocks
stale/unmarked/foreign resume; missing child is `host-blocked`; provider failure
remains filesystem-only.
