---
kind: domain
status: accepted
title: System-design authoring
slug: system-design-authoring
owners: []
sources: []
source_revisions:
  - wrapper: HEAD
    commit: cb84870
    basis: current-wrapper
generated_at: 2026-08-27T00:00:00Z
review_date: 2026-11-27
freshness: accepted-from-current-wrapper
assumptions:
  - A system design is one kind of source, not a lifecycle stage.
unknowns: []
contradictions: []
acceptance:
  state: accepted
  accepted_at: 2026-08-27
  accepted_by: maintainer
workflows: []
---

# System-design authoring

## Summary

Authoring and structuring a **system design** as source material via the shipped
`cc-system-design` skill. A system design is one kind of source, not a lifecycle
stage: it lives under `sources/system-design/` and carries no status, acceptance
gate, or runtime record. Route "design the system for `<X>`" / "structure the
system design" here; the skill is purely skill-invoked (`/cc-system-design`).

## Scope

Inside: where a system design lives, the three-tier layout, the detail-per-file
rubric, scope separation by concern, and embedded diagrams.

Outside: everything with a lifecycle — planning ([plan-review](../plan-review/README.md)),
context acceptance, execution, delivery. The skill only shapes source files.

## Behavior

- **A source, not a stage.** A system design lives at
  `sources/system-design/<product>/<version>/<scope>/`; it is read like any named
  source (passive, INV-SEC-02) and has no status, gate, or record.
- **Three-tier layout.** Every folder has a `README.md` index; each scope's
  normative content is `design.md`; detail splits into files or sub-folders as it
  grows (scale-triggered — don't pre-fragment).
- **Altitude.** `README.md` orients only; `design.md` is the readable overview a
  reviewer can stop at (states what/why and the decisions, defers how); a
  `<concern>.md` carries one concern at depth. Every fact lives once.
- **Scope by concern, never by repository.** The top level keys the
  product/initiative; deep per-repo detail is a scope, not a separate design.
- **Feeds the existing flow.** A system design is read during ordinary context
  gathering and motivates context proposals through the existing path; plans
  ground in the resulting Product Knowledge via existing `product_knowledge`
  references. There is no design-acceptance gate and no new plan field.
- **Boundary.** The skill drafts structured source files only; it never approves,
  accepts, plans, executes, or writes Product Knowledge (INV-SKILL-01), and adds
  no engine function, runtime record, schema, or invariant.

## Interfaces

- Human request: "Design the system for `<X>`" / `/cc-system-design`
- Output: authored files under `sources/system-design/<product>/<version>/<scope>/`

## Constraints and edge cases

A bounded change needs no system design — go straight from Product Knowledge to a
plan. The reference model is this repository's own `sources/system-design/`
(dogfooded).

## Implementation references

- `.agents/skills/cc-system-design/SKILL.md`
- `test/contracts/test-contracts.sh` (skill present + rubric + no-runtime-surface checks)
- `wrapper/contracts/invariants.yaml`: INV-SKILL-01 (skill packaging/resolution)

## Verification

`sh test/acceptance.sh` (contracts suite asserts the skill ships, carries the
rubric, and adds no runtime surface). No live harness case — it is authoring
guidance, not an execution behavior.

## Provenance

Authored from the current wrapper at HEAD `cb84870` (implementation `048eb67`).
Design source `sources/system-design/context-circuit/v0.6/system-design-authoring/`
was named by the accepting request.

## Acceptance notes

Accepted 2026-08-27 from proposal `0020-domain-system-design-authoring`.
