# 0019 — cc-system-design authoring skill

- **Plan ID:** `0019-system-design-authoring-skill`
- **Status:** done
- **Repository:** `context-circuit-source`
- **Depends on:** `0010-release-assembly-and-publication`
- **Owns:** the `cc-system-design` skill (under INV-SKILL-01)

## Original request

Retroactive plan for the system-design authoring skill, built as if from an
empty repo. Source design:
`sources/system-design/context-circuit/v0.6/system-design-authoring/`.

## Objective and desired behavior

- A skill that teaches an agent to author and structure a system design as source
  material: location (`sources/system-design/<product>/<version>/<scope>/`), the
  three-tier layout (README index / normative `design.md` / detail files),
  detail-per-file altitude, scale-triggered splitting, and scope-by-concern
  separation (never one folder per git repo).
- A system design is one kind of source, not a lifecycle stage: the skill adds
  no engine function, `.runtime` record, schema, invariant, WORKFLOW action, or
  first-class `design/` area. It is purely skill-invoked and points authors at
  this repo's own `sources/system-design/` as the dogfooded reference.

## Constraints and non-goals

- Non-goal: any lifecycle, status, gate, or runtime surface.

## Product Knowledge grounding

- `invariants` (`wrapper/contracts/invariants.yaml`), `conventions`
  (`context/CONVENTIONS.md`).

## Tasks

1. **SDA-001** — author the skill + wire the release allowlist and tests.

## Acceptance & verification

- Skill body carries the full rubric; skill is shipped and resolved by path.
- `sh test/release/test-release.sh`, `sh test/contracts/test-contracts.sh`.

## Assumptions, open questions, risks

- Copies the `cc-run-stack` allowlist-wiring pattern exactly; no contract bump.

## Expected commits and delivery notes

A skill plus four allowlist edits; no core contract change.
