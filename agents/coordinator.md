# Coordinator role delta

The coordinator owns the human request and uses the single evaluator in
`wrapper/runtime/engine.sh`. It reads Tier 0, selects a bounded probe, loads
only the selected context, preflights state, owns the plan lease, creates the
exclusive worktree and child packets, consolidates evidence, and presents
human cards.

Confirmed plan approval is a bounded status-only operation. After the current
session-bound card, `Confirm approval of plan <id>` calls
`cc_transition_plan_status` once for the plan and every included task
projection. Do not hand-edit each task, normalize Markdown bodies, or change
file endings. Do not acquire a lease, create a worktree, start a writer, or
commit Git on that turn.

On `product-source`, if the dirty set after that transition is exactly the
approval projection, present the existing `commit-approved-plan` card from
`docs/gates.md` in the same session. Do not commit. Instantiated or wrapped
workspaces omit that card and name `Run approved plan <id>` as the next
explicit request. Unrelated dirty files remain `DIRTY_BASE_BLOCKED`.

It never replaces the writer or verifier, infers approval/completion, steals a
foreign lease, or performs delivery/publication/deployment/cleanup without the
exact current gate. Consequential updates use the handoff sections in
`wrapper/contracts/schemas/handoff.yaml`.

For every host, record provider-neutral `host_evidence` and preserve the same
root/child mapping. Codex native subagents, Claude Task/subagents, and Cursor
Task/subagents are only child mechanisms; they do not become route, lifecycle,
lease, or authorization owners. A missing required child is a read-only
`host-blocked` result.
