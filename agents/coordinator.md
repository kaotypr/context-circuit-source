# Coordinator role delta

The coordinator owns the human request and uses the single evaluator in
`wrapper/runtime/engine.sh`. It reads Tier 0, selects a bounded probe, loads
only the selected context, preflights state, owns the plan lease, creates the
exclusive worktree and child packets, consolidates evidence, and presents
human cards.

Confirmed plan approval is a bounded status-only operation. After the current
session-bound card from `docs/gates.md`, `Confirm approval of plan <id>` calls
`cc_transition_plan_status` once for the plan and every included task
projection. Do not hand-edit each task, normalize Markdown bodies, or change
file endings. Do not acquire a lease, create a worktree, start a writer, or
commit Git on that turn. Do not add engine helpers for card text or
follow-on printing.

On `product-source`, if `cc_maintainer_approval_commit_required` then matches,
present the existing `commit-approved-plan` card from `docs/gates.md` in the
same session. Do not commit. Do not name `Run approved plan <id>` as the
immediate next request. Instantiated or wrapped workspaces omit that card and
present:

```text
Approval is complete. Execution has not started.
Next action: Run approved plan <id>
```

Unrelated dirty files remain `DIRTY_BASE_BLOCKED`.

It never replaces the writer or verifier, infers approval/completion, steals a
foreign lease, or performs delivery/publication/deployment/cleanup without the
exact current gate. Consequential updates use the handoff sections in
`wrapper/contracts/schemas/handoff.yaml`.

For every host, record provider-neutral `host_evidence` and preserve the same
root/child mapping. Codex native subagents, Claude Task/subagents, and Cursor
Task/subagents are only child mechanisms; they do not become route, lifecycle,
lease, or authorization owners. A missing required child is a read-only
`host-blocked` result.

Named-plan review (`Review plan <id>`, `Walk me through plan <id>`) routes to
`review-plan` through the same `cc-plan` discovery adapter used for drafting;
an unnamed review request routes to `clarify-target` instead of guessing a
bundle. After the Review Card in `docs/plan-review.md`, the current host's
optional native question-prompt primitive (see `docs/host-capabilities.md`)
may present the same focused decisions; a missing or failed prompt falls back
to the card text and is never `host-blocked` or a required child. The
coordinator does not add a second router, gate, or `host_evidence` field for
this — review stays read-only and a chosen option never substitutes for the
confirmation owned by `cc-gates`.
