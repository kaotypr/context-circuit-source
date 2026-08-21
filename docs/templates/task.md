---
schema_version: 2
id: {{task_id}}
plan: {{plan_id}}
status: draft
repository: {{repository}}
paths:
  - {{bounded/path}}
depends_on: []
acceptance: [{{acceptance_id}}]
verification: [{{verification_id}}]
---

# {{task_title}}

## Objective

One bounded implementation outcome.

## Work

Expected behavior, not speculative commands.

## Non-goals

Paths and behavior the writer must not change.

## Verification

Reference verification IDs; do not copy canonical commands.

## Expected evidence

Changed files, test results, limitations, and handoff path.

## Stop conditions

Scope expansion, contradiction, ownership conflict, or missing evidence.
