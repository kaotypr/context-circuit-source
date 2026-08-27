# 0008 — Archive and restore

- **Plan ID:** `0008-archive-and-restore`
- **Status:** done
- **Repository:** `context-circuit-source`
- **Depends on:** `0004-planning`
- **Owns (invariants):** INV-ARCHIVE-01, INV-ARCHIVE-02

## Original request

Retroactive plan for the archive/restore increment, built as if from an empty
repo. Source design:
`sources/system-design/context-circuit/v0.5/core/05-planning-and-execution.md`
(plan organization).

## Objective and desired behavior

- Archiving moves `plans/<id>/` to `plans/archive/<id>/`, removes the row from
  `plans/INDEX.md`, and preserves plan files, status, and runtime evidence. It
  performs no plan-status or execution validation and implies no completion,
  delivery, cleanup, or deletion.
- The normal agent must not read or traverse `plans/archive/`; an explicit
  restore is the only operation that returns a named archived plan to the active
  area before normal reading resumes.

## Constraints and non-goals

- Non-goal: any status change or deletion. Archive is purely organizational.

## Product Knowledge grounding

- `invariants` (`wrapper/contracts/invariants.yaml`), `conventions`
  (`context/CONVENTIONS.md`).

## Tasks

1. **ARC-001** — exact archive/restore moves + plan-organization lock.
2. **ARC-002** — `cc-archive` skill.

## Acceptance & verification

- Move preserves status/evidence; only explicit restore returns a plan; archive
  never traversed in normal orientation.
- `sh test/archive/test-archive.sh`.

## Assumptions, open questions, risks

- Risk: an agent reading archived plans as active context — prevented by
  INV-ARCHIVE-02.

## Expected commits and delivery notes

Filesystem moves within the workspace; no repository delivery involved.
