---
id: FIXTURE-0001
plan_id: document-contract-foundation
title: Invalid task status
status: executing
description: A task status outside the lifecycle projection must fail its schema.
repository: context-circuit
area: document-contracts
implementation_scope:
  - docs/document-system.md
test_scope:
  - schema fixture
verification_commands:
  - git diff --check
acceptance_criteria:
  - The task schema rejects runtime status values.
expected_evidence:
  - validation result
dependencies: []
stop_conditions:
  - scope expands
---

# Expected schema failure

`executing` is a runtime/session state, not a task lifecycle projection.
