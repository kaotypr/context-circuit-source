---
schema_version: 2
id: VEL-003
plan: verification-evidence-layers
status: ready
repository: context-circuit-source
paths:
  - docs/host-capabilities.md
  - wrapper/migrations/README.md
  - wrapper/migrations/upgrade.sh
  - test/hosts/test-host-adapters.sh
  - test/security/test-boundaries.sh
  - test/upgrades/test-upgrades.sh
  - test/release/test-release.sh
  - test/acceptance.sh
depends_on: [VEL-002]
acceptance: [VEL-AC-04, VEL-AC-05, VEL-AC-06, VEL-AC-07]
verification: [VEL-VT-03, VEL-VT-04, VEL-VT-05, VEL-VT-06, VEL-VT-07]
expected_evidence:
  - Cross-host fixtures preserve limitation and waiver outcomes without false passes.
  - Upgrade and rollback keep legacy evidence readable and future checks explicit.
  - Release inspection and full semantic acceptance pass.
stop_conditions:
  - Tests require implicit external provider calls or credentials.
  - Legacy migration rewrites completed evidence.
  - Release artifacts contain runtime evidence, waivers, transcripts, or provider payloads.
---

# Verify host limitations, migration, and semantic behavior

## Objective

Prove evidence-layer enforcement remains provider-neutral, compatible,
credential-free, and complete across released workspaces.

## Work

Add cross-host limitation and waiver fixtures, security negatives, legacy
upgrade/rollback behavior, release exclusions, and full semantic acceptance.
Keep live provider checks optional and bounded to declared outcomes.

## Non-goals

Do not authenticate a provider, run implicit external checks, rewrite completed
evidence, publish artifacts, or repair unrelated failures.

## Verification

Use VEL-VT-03 through VEL-VT-07.

## Expected evidence

Host unavailable/blocked/waived results, credential and provider-data
negatives, migration round trips, artifact inventory, and independent
acceptance mapping.

## Stop conditions

Stop on false pass, external activity without a gate, historical evidence
mutation, release leakage, or a regression outside this plan.
