---
kind: domain
status: accepted
title: System-design authoring
slug: system-design-authoring
owners: []
sources: []
source_revisions:
  - wrapper: HEAD
    commit: c67dd23
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
stage: it lives under `sources/system-design/` (product-level) or, for one intent,
under `intent/<id>/detail/` (intent detail). Both homes use the same three-tier
layout and split-when-earned rule; neither carries status, an acceptance gate, or
a runtime record. Route "design the system for `<X>`" / "structure the system
design" / "write this change out by topic" here; the skill is purely
skill-invoked (`/cc-system-design`). `cc-intent` may recommend or honor a request
for the intent-detail home; it does not duplicate this rubric.

## Scope

Inside: where a system design lives (both homes), the three-tier layout, the
detail-per-file rubric, scope separation by concern, and embedded diagrams.

Outside: everything with a lifecycle — planning ([plan-review](../plan-review/README.md)),
context acceptance, execution, delivery, and Gate 1 itself
([intent](../intent/README.md)). The skill only shapes source files. Intent detail
is not Product Knowledge and has no approval of its own.

## Behavior

- **A source, not a stage.** A product-level system design lives at
  `sources/system-design/<product-or-project>/<grouping>/<scope>/`; an intent's
  fuller write-up lives at `intent/<id>/detail/`. Both are read as named sources
  (passive, INV-SEC-02) and have no status, gate, or record. Choose the
  product-level home for a picture that will spawn several intents; choose intent
  detail for the shape of one already-drafted (or about-to-be-drafted) intent.
- **Grouping segment.** The middle segment is a *grouping* dimension, default a
  version (3-number semver with a `v` prefix going forward, e.g. `v0.6.1`); a named
  grouping (`phase-2`, `Q1`, a milestone) is an explicit author choice, one
  dimension per product. When unspecified, the agent derives it from the version
  the design targets — the project's **next** version for forthcoming design, the
  **current** version for as-built state — read from the project's own source of
  truth (`workspace.yaml` `template_version`, release manifests), never invented,
  and surfaces a genuinely ambiguous choice rather than guessing. Legacy 2-number
  folders (`v0.5`, `v0.6`) are not renamed; `<product-or-project>` may collapse to
  one segment in a single-project workspace.
- **Three-tier layout.** Every folder has a `README.md` index; each scope's
  normative content is `design.md`; detail splits into files or sub-folders as it
  grows (scale-triggered — don't pre-fragment).
- **Altitude.** `README.md` orients only; `design.md` is the readable overview a
  reviewer can stop at (states what/why and the decisions, defers how); a
  `<concern>.md` carries one concern at depth. Every fact lives once.
- **Scope by concern, never by repository.** The top level keys the
  product/initiative; deep per-repo detail is a scope, not a separate design.
- **Feeds the existing flow.** A product-level system design is read during
  ordinary context gathering and motivates context proposals through the existing
  path; plans ground in the resulting Product Knowledge via existing
  `product_knowledge` references. There is no design-acceptance gate and no new
  plan field. Intent detail is the confirmed topic shape of one intent: after
  Gate 1, planning uses it when present, without replacing the post-approval
  read of the real code and without entering `contract_digest`.
- **Boundary.** The skill drafts structured source files only; it never approves,
  accepts, plans, executes, or writes Product Knowledge (INV-SKILL-01), and adds
  no engine function, runtime record, schema, or invariant.

## Interfaces

- Human request: "Design the system for `<X>`" / "write this change out by topic" / `/cc-system-design`
- Output: authored files under `sources/system-design/<product-or-project>/<grouping>/<scope>/` or `intent/<id>/detail/`

## Constraints and edge cases

A bounded change needs no product-level system design — go straight from Product
Knowledge to a plan. A small single-outcome intent needs no intent detail. The
reference model is this repository's own `sources/system-design/` (dogfooded);
intent detail reuses that inner rubric at `intent/<id>/detail/`.

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
Updated 2026-08-29 from proposal `0027-system-design-grouping` (v0.6.1
design-layout-grouping): the middle path segment is a grouping dimension
(version-default, 3-number semver going forward; named groupings allowed) with an
agent current/next default derivation. Implementation `c67dd23`.
