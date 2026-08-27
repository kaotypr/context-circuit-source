# Proposal 0019 — change domain: delivery (drift guard)

- id: 0019-change-delivery
- target_context_unit: context/domains/delivery/README.md
- operation: change
- confidence: high
- status: review-needed
- related_plan: —
- affected_commits: [74eb510]
- evidence_refs:
  - wrapper/contracts/invariants.yaml (INV-DELIVER-01, extended with the drift-guard clause)
  - wrapper/runtime/engine.sh (cc_delivery_drift, cc_delivery_rebase)
  - .agents/skills/cc-deliver/SKILL.md (Drift guard section)
  - test/delivery/test-delivery.sh
  - sources/system-design/context-circuit/v0.6/run-stack/delivery.md

## Statement

Delivery stays separate, explicit, per repository, per plan, and targets each
repository's recorded `anchor_branch` — unchanged in every respect but one: v0.6
adds a **drift guard** to INV-DELIVER-01. When a plan is delivered and its recorded
base has diverged from the current `anchor_branch` tip (because a sibling plan
already merged), the plan is rebased onto the current tip and re-verified before
its pull request opens; a plan is never merged from a base that no longer reflects
the branch it will land on.

## Change (edits to the delivery page on acceptance)

- In **Behavior**: add the drift guard — before opening a pull request, check
  `cc_delivery_drift`; if drift is detected, `cc_delivery_rebase` rebases the
  execution branch onto the current anchor tip, flags re-verification, and a rebase
  conflict is reported as blocked (`DELIVERY_REBASE_CONFLICT`) with work preserved.
- Clarify the runtime still performs no delivery action itself: the only merge the
  runtime authors is the integration **base** on a plan's own branch (run-stack,
  INV-CONCURRENCY-02), never a delivery merge; delivery remains report-only + the
  drift rebase.
- In **Implementation references**: add `cc_delivery_drift`, `cc_delivery_rebase`,
  and the extended INV-DELIVER-01.
