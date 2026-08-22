---
schema_version: 2
id: RB-004
plan: repository-bootstrap
status: ready
repository: context-circuit-source
paths:
  - wrapper/migrations/upgrade.sh
  - wrapper/migrations/README.md
  - scripts/release-artifact.sh
  - scripts/release-manifest.txt
  - docs/getting-started.md
  - docs/configuration.md
  - docs/runtime-contract.md
  - README.md
  - plans/README.md
depends_on: [RB-002, RB-003]
acceptance: [RB-AC-05, RB-AC-06]
verification: [RB-VT-04, RB-VT-05, RB-VT-07]
---

# Integrate runtime, migration, release, and documentation boundaries

## Objective

Make repository bindings resumable and upgrade-safe while keeping maintainer
plans, local bindings, and repository contents out of release artifacts.

## Work

Connect resolved bindings to isolated worktree preparation and runtime records.
Define upgrade and rollback preservation for local bindings and repository
work. Update the blank template, onboarding documentation, release manifest,
dist builder, and source layout documentation.

## Non-goals

Do not make local paths part of shared identity, or copy local repositories or
maintainer plans into a release artifact.

## Verification

Use RB-VT-04, RB-VT-05, and RB-VT-07.

## Expected evidence

Migration fixtures, clean release inventory, rebuilt versioned dist artifact, and updated
onboarding/runtime documentation.

## Stop conditions

Stop if upgrade or rollback can overwrite local bindings, dirty work, runtime
evidence, or repository contents.
