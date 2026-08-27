# Proposal 0020 — add domain: system-design-authoring

- id: 0020-domain-system-design-authoring
- target_context_unit: context/domains/system-design-authoring/README.md
- operation: add
- confidence: high
- status: review-needed
- related_plan: —
- affected_commits: [048eb67]
- evidence_refs:
  - sources/system-design/context-circuit/v0.6/system-design-authoring/
  - .agents/skills/cc-system-design/SKILL.md
  - test/contracts/test-contracts.sh (skill present + rubric + no-runtime-surface checks)

## Statement

v0.6 ships **`cc-system-design`**, a skill for authoring and structuring a system
design as *source material*. A system design is **one kind of source**, not a
lifecycle stage: it lives under `sources/system-design/<product>/<version>/<scope>/`
and carries **no status, no acceptance gate, no runtime record**. The skill's whole
value is getting the structure and altitude right; it warrants a small retrieval
unit so a reader asking "how do I write a system design here?" finds the rubric.

## Proposed unit (summary to author on acceptance)

- **What it is.** A read-as-procedure skill (INV-SKILL-01), purely skill-invoked
  (`/cc-system-design`, no WORKFLOW action). It drafts structured source files and
  nothing else — it never approves, accepts, plans, executes, or writes Product
  Knowledge, and adds no engine function, `.runtime` record, schema, or invariant.
- **Layout.** Three-tier under `sources/`: a `README.md` index, a normative
  `design.md` per scope, detail files/sub-folders split scale-triggered.
- **Altitude rubric.** `README.md` orients; `design.md` is the readable overview a
  reviewer can stop at (states what/why + decisions, defers how); a `<concern>.md`
  carries one concern at depth. Every fact lives once; don't pre-fragment.
- **Scope by concern, never by repository.** The top level keys the
  product/initiative; deep per-repo detail is a scope, not a separate design.
- **Feeds the existing flow.** A system design is read during ordinary context
  gathering and motivates context proposals through the existing path; plans ground
  in the resulting Product Knowledge via the existing `product_knowledge`
  references — no new gate and no new plan field.
- **Dogfood.** This repository's own `sources/system-design/` is the reference model.

## Verification

`sh test/acceptance.sh` (contracts suite asserts the skill ships, carries the
rubric, and adds no runtime surface). No live harness case — it is authoring
guidance, not an execution behavior.
