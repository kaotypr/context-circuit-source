# 0010 — Source/template boundary, seed, release assembly and publication

- **Plan ID:** `0010-release-assembly-and-publication`
- **Status:** done
- **Repository:** `context-circuit-source`
- **Depends on:** `0001`–`0009` (the full product surface it ships)
- **Owns (invariants):** release_boundary (`wrapper/manifest.yaml`); cites INV-RUNTIME-01, INV-SEC-01, INV-COMMIT-01

## Original request

Retroactive plan for the source/template boundary and release increment, built
as if from an empty repo. Source design:
`sources/system-design/context-circuit/v0.5/core/09-source-and-template.md` and
`10-template-publication.md`.

## Objective and desired behavior

- One repository is both the **source** (factory) and a working example of the
  **template** it produces. `wrapper/manifest.yaml` makes the boundary precise:
  shipped (replaceable on upgrade), mutable_seed (copied once), never_ship.
- An upgrade replaces only the machinery and always preserves `workspace.yaml`,
  `context/`, `sources/`, `plans/`, `repositories.local.yaml`, `repositories/`,
  and `.runtime/`.
- The assembler ships only template-owned files (explicit allowlist), never
  secrets. Publication is workflow-only — the runtime never pushes or publishes.
- The template carries its own SemVer line, stamped at publish, gated on the full
  acceptance suite passing first.

## Constraints and non-goals

- Non-goal: the human-simulated harness (0011–0012), which consumes the assembled
  artifact but is source-only tooling.

## Product Knowledge grounding

- `invariants` (`wrapper/contracts/invariants.yaml`), `architecture`
  (`context/ARCHITECTURE.md`).

## Tasks

1. **REL-001** — release boundary, manifest, blank template seed.
2. **REL-002** — assembler + upgrade migrations.
3. **REL-003** — publication workflow behind the acceptance gate.

## Acceptance & verification

- Manifest classifies every path; upgrade preserves user state; assembly rejects
  non-allowlisted files; publish runs only after acceptance passes; no
  attribution/credential in the record.
- `sh test/acceptance.sh`, `sh test/release/test-release.sh`,
  `sh test/release/test-publish.sh`.

## Assumptions, open questions, risks

- Risk: an upgrade clobbering user knowledge — prevented by the upgrade boundary
  and the release test.

## Expected commits and delivery notes

Publication is one commit plus a tag on the template line; the maintainer is the
sole recorded author (INV-COMMIT-01).
