---
id: BAD-ENUM
plan_id: compact-plan
title: Invalid task enum
status: executing
description: This task uses a runtime value outside the task projection.
repository: context-circuit
area: planning-contract
implementation_scope:
  - docs/
test_scope:
  - task fixture
verification_commands:
  - git diff --check
acceptance_criteria:
  - The task status must use the plan projection.
expected_evidence:
  - A blocked enum validation result.
dependencies: []
stop_conditions:
  - Stop on an invalid enum.
---

# Invalid task enum
