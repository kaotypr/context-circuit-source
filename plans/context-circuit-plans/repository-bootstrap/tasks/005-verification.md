---
schema_version: 2
id: RB-005
plan: repository-bootstrap
status: draft
repository: context-circuit-source
paths:
  - test/contracts/test-contracts.sh
  - test/routing/test-router.sh
  - test/ownership/test-ownership.sh
  - test/recovery/test-recovery.sh
  - test/upgrades/test-upgrades.sh
  - test/security/test-boundaries.sh
  - test/release/test-release.sh
  - test/acceptance.sh
depends_on: [RB-001, RB-002, RB-003, RB-004]
acceptance: [RB-AC-01, RB-AC-02, RB-AC-03, RB-AC-04, RB-AC-05, RB-AC-06]
verification: [RB-VT-01, RB-VT-02, RB-VT-03, RB-VT-04, RB-VT-05, RB-VT-06, RB-VT-07, RB-VT-08]
---

# Verify the self-hosted workflow independently

## Objective

Exercise the repository-bootstrap plan through independent contract, routing,
ownership, recovery, security, release, and complete acceptance suites.

## Work

Add semantic fixtures for external and convenience paths, clone confirmation,
existing and dirty targets, offline fallback, resume, migration, and release
exclusion. Run the full acceptance suite and record results in the durable
implementation log.

## Non-goals

Do not mark the plan done from implementation evidence alone, and do not use
the verifier to repair implementation or runtime state.

## Verification

Use all RB-VT IDs, with RB-VT-08 as the final suite.

## Expected evidence

Independent verifier handoff, test output, generated artifact inventory,
budget measurements, migration results, known deviations, and unresolved
risks.

## Stop conditions

Stop on failed verification, missing child primitives, scope expansion, or
incomplete evidence.
