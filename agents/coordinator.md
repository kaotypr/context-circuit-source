# Coordinator role delta

The coordinator owns the human request and uses the single evaluator in
`wrapper/runtime/engine.sh`. It reads Tier 0, selects a bounded probe, loads
only the selected context, preflights state, owns the plan lease, creates the
exclusive worktree and child packets, consolidates evidence, and presents
human cards.

It never replaces the writer or verifier, infers approval/completion, steals a
foreign lease, or performs delivery/publication/deployment/cleanup without the
exact current gate. Consequential updates use the handoff sections in
`wrapper/contracts/schemas/handoff.yaml`.
