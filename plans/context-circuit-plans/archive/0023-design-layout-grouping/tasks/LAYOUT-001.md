---
id: LAYOUT-001
title: Refine the layout convention in the cc-system-design skill
repositories: [context-circuit-source]
paths: [.agents/skills/cc-system-design/SKILL.md]
depends_on: []
acceptance: [LAYOUT-AC-01, LAYOUT-AC-02]
verification: [LAYOUT-VT-01]
---

## Intended behavior

The `cc-system-design` skill teaches the generalized layout: the middle path
segment is a *grouping* dimension, not a strict version. A version is the default
grouping; a named grouping (a slug such as `phase-2`, `Q1`, or a milestone) is an
explicit author choice, with one grouping dimension chosen per product rather
than release and phase folders mixed arbitrarily.

## Concrete change and affected surfaces

Update `.agents/skills/cc-system-design/SKILL.md` so it: (1) names the path
`sources/system-design/<product-or-project>/<grouping>/<scope>/` and defines the
grouping dimension; (2) states that version groupings use 3-number semver with a
`v` prefix going forward (`v0.6.1`, not `v0.6`) while legacy `v0.5`/`v0.6` folders
remain and are not retroactively renamed; (3) specifies the agent's default
derivation — the project's next version for forthcoming design, the current
version for as-built, read from the project's own source of truth and never
invented; and (4) says the agent surfaces genuine ambiguity instead of guessing.
The `README.md` + `design.md` structure at every level is unchanged.

## Inputs and outputs

Inputs: the v0.6 authoring rubric and the observed gaps (non-version groupings;
inconsistent version depth). Outputs: the refined convention in the skill text.
No lifecycle, status, gate, runtime, schema, or new owner is added.

## Context references

- sources/system-design/context-circuit/v0.6.1/design-layout-grouping/design.md

## Risks and open questions

- During transition a tree holds both 2- and 3-number groupings; the skill must
  present this as expected, not an error.

## Stop conditions

- Do not retroactively rename legacy folders; do not add a runtime, schema, or
  first-class `design/` area.
