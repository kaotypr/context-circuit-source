---
schema_version: 2
id: GRE-003
plan: greenfield-repository-create-empty
status: draft
repository: context-circuit-source
paths:
  - wrapper/migrations/README.md
  - wrapper/migrations/upgrade.sh
  - scripts/release-artifact.sh
  - scripts/release-manifest.txt
  - test/contracts/test-contracts.sh
  - test/upgrades/test-upgrades.sh
  - test/release/test-release.sh
  - test/acceptance.sh
depends_on: [GRE-002]
acceptance: [GRE-AC-05, GRE-AC-06]
verification: [GRE-VT-01, GRE-VT-06, GRE-VT-07, GRE-VT-08]
expected_evidence:
  - Upgrade and rollback preserve created repositories and ignored bindings.
  - Release assembly excludes repository contents and local machine paths.
  - Full semantic acceptance passes with existing bootstrap behavior unchanged.
stop_conditions:
  - Migration moves, rewrites, or deletes a registered repository.
  - Release assembly includes repositories.local.yaml or repository contents.
  - Existing bootstrap, binding, or worktree behavior regresses.
---

# Verify recovery, security, release, and compatibility

## Objective

Complete compatibility and release coverage for the new action without
changing existing repository ownership or artifact boundaries.

## Work

Extend migration, upgrade, release, and semantic fixtures. Prove created
repositories and local bindings survive wrapper replacement, stay excluded
from artifacts, and remain compatible with existing bootstrap and isolated
execution worktrees.

## Non-goals

Do not publish an artifact, migrate repository contents, or repair unrelated
repository and host failures.

## Verification

Use GRE-VT-01, GRE-VT-06, GRE-VT-07, and GRE-VT-08.

## Expected evidence

Upgrade/rollback round trips, release artifact inventory, existing bootstrap
regression results, and independent acceptance mapping.

## Stop conditions

Stop on repository loss, local-binding leakage, release contamination, or a
regression outside the declared plan scope.
