---
id: COMPX-0001
plan_id: complex-plan
title: Validate complex bundle references
status: ready
description: Check that specialist companions reference canonical plan fields.
repository: context-circuit
area: planning-verification
implementation_scope:
  - docs/
  - test/
test_scope:
  - complex bundle fixture
verification_commands:
  - git diff --check
  - sh test/acceptance.sh
acceptance_criteria:
  - The task preserves the plan authority.
expected_evidence:
  - A passing complex bundle validation result.
dependencies:
  - document-contract-foundation
stop_conditions:
  - Stop if a companion copies lifecycle authority.
external_status:
  provider: fixture-provider
  value: in_progress
---

# Validate complex bundle references

Read `plan.yaml#acceptance_criteria` and the specialist rationale by reference.
