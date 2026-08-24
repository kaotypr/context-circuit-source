---
id: 0012-accept-wrapper-terminology
target_context_unit: context/TERMINOLOGY.md
operation: change
statement: >
  Keep "wrapper" as accepted vocabulary meaning the universal project workspace
  product (a synonym for the workspace product). The design source has been
  updated (08 §1/§5, 09) to accept "wrapper" instead of deprecating it. Add the
  term to context/TERMINOLOGY.md and the shipped docs/terminology.md; context's
  existing "wrapper" usage is retained (no stripping).
evidence_refs:
  - sources/context-circuit-v0.5-design/08-terminology.md  # §1 Wrapper row, §5 rewritten (updated 2026-08-24)
  - sources/context-circuit-v0.5-design/09-source-and-template.md  # wrapper accepted (updated 2026-08-24)
  - docs/terminology.md
affected_repositories:
  - context-circuit-source
affected_commits:
  - 4b8ac0b
confidence: high
status: review-needed
---

# Accept "wrapper" as a glossary term (synonym for the workspace product)

Direction reversed from the earlier draft: instead of stripping "wrapper" from
context to match the design, the DESIGN was updated to accept "wrapper," and
context keeps using it. The design edits are already made in `sources/` (08 §1
adds a Wrapper row; 08 §5 retitled "The term 'wrapper'" and no longer deprecates
it; 09 accepts it).

## Context changes to accept

1. `context/TERMINOLOGY.md` — add a glossary row: **Wrapper** = accepted synonym
   for the universal project workspace (the Context Circuit product); the
   directory `wrapper/` holds its shipped layer.
2. `docs/terminology.md` (shipped glossary, Table 1) — add the same **Wrapper**
   term so the shipped projection matches the design.
3. No stripping: the existing "wrapper" / "wrapper contracts" usage across
   `ARCHITECTURE.md`, `PROJECT.md`, `WORKSPACE.md`, `maintainer.md`, the domain
   pages, and the 2026-08-21 decision stays as-is.

## Nuance to preserve

The term names the product; the directory `wrapper/` names one part (the shipped
layer). Keep directory references as explicit paths so the two are not confused.

## Acceptance action

Add the Wrapper row to `context/TERMINOLOGY.md` and `docs/terminology.md`. The
`sources/` design edits are already applied (direct maintainer change).
