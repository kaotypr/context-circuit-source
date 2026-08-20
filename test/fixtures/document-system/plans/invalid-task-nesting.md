---
id: BAD-NESTING
plan_id: compact-plan
title: Invalid task nesting
status: draft
description: This task uses a scalar where the schema requires a list.
repository: context-circuit
area: planning-contract
implementation_scope: docs/
test_scope:
  - task fixture
verification_commands:
  - git diff --check
acceptance_criteria:
  - The task has the expected nested list shapes.
expected_evidence:
  - A blocked nesting validation result.
dependencies: []
stop_conditions:
  - Stop on an invalid nested shape.
---

# Invalid task nesting
