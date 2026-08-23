---
schema_version: 2
id: RB-002
plan: repository-bootstrap
status: done
repository: context-circuit-source
paths:
  - wrapper/runtime/engine.sh
  - wrapper/contracts/routes.yaml
  - wrapper/contracts/context-sets.yaml
  - wrapper/contracts/schemas/workspace.yaml
depends_on: [RB-001]
acceptance: [RB-AC-02, RB-AC-04, RB-AC-05]
verification: [RB-VT-02, RB-VT-03, RB-VT-06]
---

# Resolve and validate repository bindings

## Objective

Resolve a logical repository to a safe external path, relative path, or
optional `repositories/<id>` path and validate its Git identity before use.

## Work

Implement host-neutral binding lookup and validation. Reject traversal,
unsafe symlinks, missing repositories, identity mismatches, and dirty source
checkouts where isolated execution requires a clean base. Preserve explicit
path evidence for resumable sessions.

## Non-goals

Do not infer a repository from broad filesystem scans or silently select a
different repository when a binding is missing or ambiguous.

## Verification

Use RB-VT-02, RB-VT-03, and RB-VT-06.

## Expected evidence

Positive and negative binding fixtures, path safety results, Git identity
checks, and source-preservation evidence.

## Stop conditions

Stop on ambiguous bindings, dirty ownership, path escape, or changed primary
evidence.
