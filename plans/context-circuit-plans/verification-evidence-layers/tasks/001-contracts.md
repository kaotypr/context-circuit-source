---
schema_version: 2
id: VEL-001
plan: verification-evidence-layers
status: ready
repository: context-circuit-source
paths:
  - wrapper/contracts/invariants.yaml
  - wrapper/contracts/schemas/plan.yaml
  - wrapper/contracts/schemas/task.yaml
  - wrapper/contracts/schemas/delegation.yaml
  - wrapper/contracts/schemas/handoff.yaml
  - wrapper/contracts/schemas/completion.yaml
  - wrapper/manifest.yaml
  - test/contracts/test-contracts.sh
depends_on: []
acceptance: [VEL-AC-01, VEL-AC-02, VEL-AC-03, VEL-AC-04, VEL-AC-06]
verification: [VEL-VT-01]
expected_evidence:
  - One canonical evidence-layer vocabulary and comparison model.
  - Schema fixtures for matching, non-matching, missing, blocked, and waived evidence.
  - Legacy compatibility rules preserve completed historical evidence.
stop_conditions:
  - Implicit cross-layer inference or a global strength hierarchy misrepresents evidence layers.
  - Waiver metadata grants authorization or counts as verifier success.
  - Migration invents required or observed evidence for historical work.
---

# Define evidence-layer contracts and compatibility

## Objective

Define required and observed evidence in canonical plan and runtime schemas,
including safe legacy and waiver behavior.

## Work

Specify evidence layers and exact-match validation across acceptance,
verification, task, delegation, handoff, and completion records.
Define pass, fail, blocked, and waived outcomes. Require mappings for new and
unfinished work while preserving completed historical records unchanged.

## Non-goals

Do not implement verifier logic, choose a browser provider, approve waivers, or
rewrite existing completed plan evidence.

## Verification

Use VEL-VT-01.

## Expected evidence

Schema, invariant, and manifest diffs plus fixtures covering layer
comparability, weak evidence, missing capability, waiver, and legacy records.

## Stop conditions

Stop if evidence comparison is unsound, waiver ownership is ambiguous, or
compatibility creates evidence that was never observed.
