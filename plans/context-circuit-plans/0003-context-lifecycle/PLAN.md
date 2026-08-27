# 0003 — Product Knowledge and context lifecycle

- **Plan ID:** `0003-context-lifecycle`
- **Status:** done
- **Repository:** `context-circuit-source`
- **Depends on:** `0001-contracts-runtime-foundation`
- **Owns (invariants):** INV-KNOWLEDGE-01, INV-KNOWLEDGE-02

## Original request

Retroactive plan for the context/Product-Knowledge increment, built as if from
an empty repo. Source design:
`sources/system-design/context-circuit/v0.5/core/04-context-lifecycle.md`.

## Objective

Hold durable, retrieval-optimized Product Knowledge and gate every change to it
behind an explicit human decision.

## Desired behavior

- The context index is a **retrieval catalog** — stable ids, summaries, domains,
  repositories, decisions, constraints, freshness, provenance — while staying
  human-readable. It is not a full copy of page content.
- A context update is accepted only by an explicit, separate human decision.
  Completion, verification, worker claims, and code changes never silently accept
  one. A plan may be `done` while a proposal is still pending.

## Scope / deliverables

- Schemas: `context-index.yaml`, `context-proposal.yaml`, `context-impact.yaml`.
- `context/INDEX.md` retrieval catalog + the `context/proposals/` area.

## Non-goals

- Planning logic (0004) or completion-time reconciliation (0007) — this plan
  provides the shapes and the acceptance gate they use.

## Product Knowledge grounding

- `architecture` (`context/ARCHITECTURE.md`), `invariants`
  (`wrapper/contracts/invariants.yaml`), `sources` (`context/SOURCES.md`).

## Tasks

1. **CTX-001** — context index/proposal/impact schemas.
2. **CTX-002** — seed `context/INDEX.md` and the proposals area.

## Acceptance & verification

- Index carries metadata not bodies; proposals carry provenance + impact status;
  retrieval is selective; pending proposals never mutate accepted knowledge.
- `sh test/contracts/test-contracts.sh`.

## Assumptions / risks

- Risk: an agent treating a proposal as accepted knowledge — prevented by the
  no-silent-accept invariant and the separate acceptance decision.

## Delivery notes

Grounds planning (0004) and is reconciled by completion (0007).
