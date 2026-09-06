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
  - Standard and Critical become done only on an explicit mark-done with no unreadiness look; Explore is planless; verification, candidate acceptance, and delivery do not mark a plan done.
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

Completion is the `draft → done` flip on an explicit mark-done ask. Standard and
Critical become done only on that ask, with no look at work or evidence and no
unreadiness refusal. Explore is planless. Verification, candidate acceptance,
and delivery never complete a plan. When the plan affected Product Knowledge,
mark-done then updates live context files in place — the same reconcile as
gathering context. Route mark-done here. Owned by the `cc-complete` skill.

## Scope

Inside: the `draft → done` transition on ask, the optional completion record
written when execution evidence files exist, and in-place Product Knowledge
reconcile when the plan affected knowledge.

Outside: verification itself ([verification](../verification/README.md)) and
delivery ([delivery](../delivery/README.md)). Delivery does not start
reconcile and does not mark a plan done.

## Behavior

An explicit mark-done request (`plan-complete`) flips Standard or Critical
`draft → done` with no look at work or evidence (INV-COMPLETE-01). Asking for
several named plans does that for each named plan on the same path. A plan
that was never built, failed a check, or has no evidence still becomes done
when the human asks. Verification, candidate acceptance, and delivery never
mark a plan complete. When execution evidence files exist, mark-done may also
write the implementation completion record (plan revision, candidate,
execution, per-repository commits, verifier result, acceptor); missing files
do not block the status change.

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

Asking to mark a plan done is enough. There is no unreadiness refusal for a
never-built, failed, blocked, or unevidenced plan. The runtime preserves
optional reconciliation references but never interprets Product Knowledge.

## Reconciliation, impact status, and staleness

A successful verifier leaves execution status `verified` and plan status
`draft` until the explicit mark-done. That evidence is not a gate on the
status flip.

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
- `wrapper/runtime/engine.sh`: `cc_plan_complete`, `cc_completion_finalize`,
  `cc_latest_execution`, `cc_context_impact_record`. `cc_completion_ready` is
  an eligibility query, not a mark-done gate.
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
in-place knowledge updates, then for ungated mark-done: asking flips status
with no unreadiness look; reconcile still writes live context files rather
than a sidecar.
