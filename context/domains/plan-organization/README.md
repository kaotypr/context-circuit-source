---
kind: domain
status: accepted
title: Plan organization
slug: plan-organization
owners: []
sources: []
source_revisions:
  - wrapper: HEAD
    commit: 4b8ac0b
    basis: current-wrapper
generated_at: 2026-08-24T00:00:00Z
review_date: 2026-11-24
freshness: accepted-from-current-wrapper
assumptions:
  - Archive and restore never validate plan status or execution state.
unknowns: []
contradictions: []
acceptance:
  state: accepted
  accepted_at: 2026-08-24
  accepted_by: maintainer
workflows:
  - .context-circuit/docs/getting-started.md
---

# Plan organization

## Summary

Setting a plan aside (archive) or bringing it back (restore) without changing
its status. Route "archive plan `<id>`" and "restore plan `<id>`" here. Owned by
the `cc-archive` skill.

## Scope

Inside: the `plans/<id>/` ↔ `plans/archive/<id>/` move under a short
organization lock, resolution by plan id only, active-index add/drop, and the
rule that the normal agent never traverses `plans/archive/`.

Outside: completion, delivery, execution cleanup, and any status/verification
validation.

## Behavior

Archiving moves `plans/<id>/` to `plans/archive/<id>/`, removes the row from
`plans/INDEX.md`, and preserves plan files, status, and runtime evidence. It
performs no plan-status or execution validation and implies no completion,
delivery, cleanup, or deletion (INV-ARCHIVE-01). The normal agent must not read
or traverse `plans/archive/`; an explicit restore is the only operation that
returns a named archived plan to the active plan area before normal reading
resumes (INV-ARCHIVE-02).

Plans are resolved by id, never by title fragment or repository name, under a
short-lived plan-organization lock. On a failed move or index update, the source
directory and index are left unchanged. An already-started execution continues
from its immutable snapshot.

## Workflows

- Archive/restore: `.context-circuit/docs/getting-started.md`

## Interfaces

- Human requests: "Archive plan `<id>`", "Restore plan `<id>`"
- Locations: `plans/<id>/`, `plans/archive/<id>/`, `plans/INDEX.md`

## Constraints and edge cases

Archive is not completion and not deletion; runtime evidence is preserved.
Restore is the only re-entry into the archived area.

## Implementation references

- `.agents/skills/cc-archive/SKILL.md`
- `.context-circuit/wrapper/runtime/engine.sh`: `cc_plan_org_lock`, `cc_plan_org_unlock`,
  `cc_plan_archive`, `cc_plan_restore`, `cc_plan_index_remove`,
  `cc_plan_index_upsert`
- `.context-circuit/wrapper/contracts/invariants.yaml`: INV-ARCHIVE-01, INV-ARCHIVE-02

## Verification

`sh test/acceptance.sh` (archive suite).

## Provenance

Authored from the current wrapper at HEAD `4b8ac0b`. Raw `sources/` was not
scanned.

## Acceptance notes

Accepted 2026-08-24 from proposal `0014-domain-plan-organization`.
Updated 2026-08-27 by explicit maintainer request to use the visible
`plans/archive/` location.
