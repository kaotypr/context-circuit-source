---
id: FIXTURE-0002
plan_id: document-contract-foundation
title: Valid non-OKF task
status: ready
description: A task contract uses its own schema and lifecycle.
repository: context-circuit
area: document-contracts
implementation_scope:
  - docs/document-system.md
test_scope:
  - schema fixture
verification_commands:
  - git diff --check
acceptance_criteria:
  - The task remains outside OKF.
expected_evidence:
  - validation result
dependencies: []
stop_conditions:
  - scope expands
---

# Task

This task front matter is validated as a task, not as an OKF concept.
