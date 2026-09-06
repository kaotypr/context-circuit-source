---
kind: domain
status: accepted
title: Completion
slug: completion
owners: []
sources: []
source_revisions:
  - wrapper: HEAD
    commit: 4b8ac0b
    basis: current-wrapper
generated_at: 2026-08-24T00:00:00Z
review_date: 2026-11-24
freshness: accepted-from-current-wrapper
assumptions:
  - Standard and Critical become done only on an explicit mark-done; Explore is planless; verification, candidate acceptance, and delivery do not mark a plan done.
unknowns: []
contradictions: []
acceptance:
  state: accepted
  accepted_at: 2026-08-24
  accepted_by: maintainer
workflows:
  - docs/getting-started.md
  - docs/product-knowledge.md
---

# Completion

## Summary

Completion is the candidate-bound transition from a verified implementation to a
done plan. Standard and Critical become done only on an explicit mark-done.
Explore is planless. Verification, candidate acceptance, and delivery never
complete a plan. When the plan affected Product Knowledge, mark-done then
updates live context files in place — the same reconcile as gathering context.
Route mark-done here. Owned by the `cc-complete` skill.

## Scope

Inside: the `draft → done` transition gated on a verified latest execution,
the durable completion record, and in-place Product Knowledge reconcile when
the plan affected knowledge.

Outside: verification itself ([verification](../verification/README.md)) and
delivery ([delivery](../delivery/README.md)). Delivery does not start
reconcile and does not mark a plan done.

## Behavior

After the latest execution passes independent verification and the human
accepts its candidate, an explicit mark-done request (`plan-complete`) flips
Standard or Critical `draft → done` when `completion-ready` still holds
(INV-COMPLETE-01). Verification, candidate acceptance, and delivery never mark
a plan complete. Completion records the plan revision, candidate, execution,
per-repository commits, verifier result, and acceptor.

When the plan affected Product Knowledge, that same mark-done starts in-place
reconcile of live `context/` files and keeps `INDEX.md` consistent
(INV-COMPLETE-02, INV-KNOWLEDGE-02). If the plan did not affect Product
Knowledge, context files stay as they are. The runtime never writes or
interprets Product Knowledge; in-place edits are a coordinator act after
mark-done. A later plan may start even if that update has not landed.
Product Knowledge is optimized first for agent retrieval and the context
index is a retrieval catalog, not a full copy of page content
(INV-KNOWLEDGE-01).

## Workflows

- Mark complete: `docs/getting-started.md`
- In-place knowledge updates: `docs/product-knowledge.md`

## Interfaces

- Human requests: "Mark `<id>` complete", "Mark `<id>` done"
- Records: `completion.yaml`, `context-impact.yaml`

## Constraints and edge cases

Completion refuses unless the latest execution is `verified`; a `failed` or
`blocked` result is reported plainly. The runtime preserves optional
reconciliation references but never interprets Product Knowledge.

## Reconciliation, impact status, and staleness

A successful verifier creates pending completion evidence — execution status
`verified`, plan status `draft` — before candidate acceptance and the explicit
mark-done.

Reconcile reads the final plan and task files and revisions, the Product
Knowledge references and grounding summary, the changed paths and commits per
repository, the worker handoffs and verifier evidence, and the current
revisions of relevant context units. The coordinator edits live context files
only when there is a durable knowledge change, or records no-update-needed.
There is no extra knowledge-acceptance gate.

Staleness: a page is stale when its freshness rule expired, a cited repository
revision materially changed, or a relevant decision changed. Stale context may
guide read-only orientation, but a consequential plan surfaces the stale
reference and refreshes it before execution.

## Implementation references

- `.agents/skills/cc-complete/SKILL.md`
- `wrapper/runtime/engine.sh`: `cc_completion_ready`, `cc_latest_execution`,
  `cc_plan_complete`, `cc_context_impact_record`
- `wrapper/contracts/schemas/completion.yaml`,
  `wrapper/contracts/schemas/context-impact.yaml`,
  `wrapper/contracts/schemas/context-index.yaml`
- `wrapper/contracts/invariants.yaml`: INV-COMPLETE-01, INV-COMPLETE-02,
  INV-KNOWLEDGE-01, INV-KNOWLEDGE-02

## Verification

`sh test/acceptance.sh` (completion suite).

## Provenance

Authored from the current wrapper at HEAD `4b8ac0b`. Raw `sources/` was not
scanned.

## Acceptance notes

Accepted 2026-08-24 from proposal `0013-domain-completion`. Refreshed for
in-place knowledge updates: mark-done is the only done trigger; reconcile
writes live context files rather than a sidecar.
