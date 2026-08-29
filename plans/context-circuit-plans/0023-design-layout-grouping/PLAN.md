# 0023 — Generalize the system-design layout to a grouping dimension

- **Plan ID:** `0023-design-layout-grouping`
- **Status:** done
- **Repository:** `context-circuit-source`
- **Depends on:** `0019-system-design-authoring-skill`
- **Owns:** the refined layout convention (in the `cc-system-design` skill)

## Original request

Retroactive plan for the v0.6.1 `design-layout-grouping` scope, authored as if
from an empty repo. Source design:
`sources/system-design/context-circuit/v0.6.1/design-layout-grouping/design.md`.

## Objective and desired behavior

- Generalize the middle segment of the system-design source path from a strict
  *version* to a *grouping* dimension:
  `sources/system-design/<product-or-project>/<grouping>/<scope>/`.
- `<grouping>` defaults to a **version**; a **named grouping** (`phase-2`, `Q1`,
  a milestone) is allowed as an explicit author choice — one grouping dimension
  per product, not mixed arbitrarily.
- Version groupings use **3-number semver with a `v` prefix going forward**
  (`v0.6.1`). Legacy 2-number folders (`v0.5`, `v0.6`) stay, not renamed.
- When the author does not specify a grouping, the agent derives it from the
  version the design targets — the project's **next** version for forthcoming
  design, the **current** version for as-built — read from the project's own
  source of truth (`workspace.yaml` `template_version`, release manifests), never
  invented; genuine ambiguity is surfaced, not guessed.

## Constraints and non-goals

- `README.md` + `design.md` structure is unchanged at every level; only the
  meaning and formatting of the middle segment are refined.
- Adds **no** lifecycle, status, gate, runtime, or new owner.
- Never one folder per git repository (unchanged rule).

## Product Knowledge grounding

- `invariants` (`wrapper/contracts/invariants.yaml`) — the skill as a
  read-as-procedure packet under INV-SKILL-01.

## Tasks

1. **LAYOUT-001** — refine the layout convention in the `cc-system-design`
   skill.

## Acceptance & verification

- The skill carries the grouping dimension, version default, named groupings,
  3-number-semver-going-forward, the legacy no-rename rule, and the default
  derivation.
- `sh test/release/test-release.sh`.

## Assumptions, open questions, risks

- A tree may hold both 2- and 3-number groupings during transition — acceptable
  and expected; the 3-number rule binds new groupings only.

## Expected commits and delivery notes

A single skill-text refinement; no contract, schema, or runtime change.
