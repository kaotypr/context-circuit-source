---
schema_version: 2
id: RB-003
plan: repository-bootstrap
status: done
repository: context-circuit-source
paths:
  - wrapper/runtime/engine.sh
  - wrapper/contracts/routes.yaml
  - docs/getting-started.md
  - docs/gates.md
depends_on: [RB-001, RB-002]
acceptance: [RB-AC-03, RB-AC-04]
verification: [RB-VT-02, RB-VT-06]
---

# Add explicit clone and bootstrap behavior

## Objective

Allow a new member to request repository cloning while preserving explicit
human authorization, offline operation, and credential boundaries.

## Work

Add a bootstrap capability and confirmation-card shape that lists repository,
canonical URL, selected remote, branch, destination, and existing-path checks.
Use host Git authentication without reading or persisting secrets. Record
success or failure as durable evidence and leave partial or failed operations
recoverable without overwriting existing work.

## Non-goals

Do not automatically clone on initialization, fetch remote changes, push,
merge, publish, deploy, or clean an existing directory.

## Verification

Use RB-VT-02 and RB-VT-06.

## Expected evidence

Routing fixtures, gate outcomes, successful local clone fixture, refusal and
offline fixtures, and no-credential-storage assertions.

## Stop conditions

Stop when the target is not explicit, credentials are requested for storage,
or a provider failure would block filesystem-only operation.
