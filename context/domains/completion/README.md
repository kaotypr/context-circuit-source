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
  - Standard completion is inferred from candidate acceptance plus delivery; Critical completion is explicit after a verified execution.
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
done plan, plus the Product Knowledge reconciliation it starts. Standard infers
completion from acceptance plus delivery; Critical requires an explicit human
completion request. Verification alone never completes a plan. Route completion
and context-update review here. Owned by the `cc-complete` skill.

## Scope

Inside: the `draft → done` transition gated on a verified latest execution,
the durable completion record, and reconciliation that stages
`context/proposals/` entries without applying them.

Outside: verification itself ([verification](../verification/README.md)) and
delivery ([delivery](../delivery/README.md)).

## Behavior

Standard completion is inferred after the latest execution passes independent
verification, the human accepts its candidate, and delivery is recorded. Critical
completion additionally requires the explicit human completion request; verification
alone never marks a plan complete (INV-COMPLETE-01). Completion records the plan
revision, candidate, execution, per-repository commits, verifier result, and
acceptor, then starts Product Knowledge reconciliation
(INV-COMPLETE-02).

Reconciliation may produce context update proposals or a stale/conflict warning,
but Product Knowledge changes only through explicit separate human acceptance; a
plan may be `done` while a proposal is pending (INV-KNOWLEDGE-02). Product
Knowledge is optimized first for agent retrieval and the context index is a
retrieval catalog, not a full copy of page content (INV-KNOWLEDGE-01).

## Workflows

- Mark complete: `docs/getting-started.md`
- Knowledge reconciliation and proposals: `docs/product-knowledge.md`

## Interfaces

- Human requests: "Mark `<id>` complete", "Accept the context update for `<id>`"
- Records: `completion.yaml`, `context-impact.yaml`
- Proposal staging: `context/proposals/`

## Constraints and edge cases

Completion refuses unless the latest execution is `verified`; a `failed` or
`blocked` result is reported plainly. The runtime preserves reconciliation
references but never interprets Product Knowledge.

## Reconciliation, impact status, and staleness

A successful verifier creates pending completion evidence — execution status
`verified`, plan status `draft` — before candidate acceptance and, for Standard,
delivery inference; Critical still needs the explicit human completion request.

Reconciliation reads the final plan and task files and revisions, the Product
Knowledge references and grounding summary, the changed paths and commits per
repository, the worker handoffs and verifier evidence, and the current revisions
of relevant context units. Each impact moves through `not-assessed` ->
`review-needed` -> `accepted` / `deferred` / `conflict`, plus `no-update-needed`;
only `review-needed`, `deferred`, and `conflict` get a stored proposal file. A
deferred proposal stays visible to future plan creation when its context is
relevant.

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
  `wrapper/contracts/schemas/context-index.yaml`,
  `wrapper/contracts/schemas/context-proposal.yaml`
- `wrapper/contracts/invariants.yaml`: INV-COMPLETE-01, INV-COMPLETE-02,
  INV-KNOWLEDGE-01, INV-KNOWLEDGE-02

## Verification

`sh test/acceptance.sh` (completion suite).

## Provenance

Authored from the current wrapper at HEAD `4b8ac0b`. Raw `sources/` was not
scanned.

## Acceptance notes

Accepted 2026-08-24 from proposal `0013-domain-completion`. This refresh is
itself an instance of the reconciliation/acceptance flow this domain describes.
