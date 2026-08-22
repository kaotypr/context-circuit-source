---
schema_version: 2
id: GRE-002
plan: greenfield-repository-create-empty
status: draft
repository: context-circuit-source
paths:
  - wrapper/runtime/engine.sh
  - wrapper/adapters/WORKFLOW.md
  - docs/configuration.md
  - docs/gates.md
  - docs/getting-started.md
  - test/ownership/test-ownership.sh
  - test/recovery/test-recovery.sh
  - test/security/test-boundaries.sh
depends_on: [GRE-001]
acceptance: [GRE-AC-02, GRE-AC-03, GRE-AC-04, GRE-AC-05]
verification: [GRE-VT-03, GRE-VT-04, GRE-VT-05]
expected_evidence:
  - Confirmed creation stages, validates, binds, and projects one repository.
  - Failure injection leaves no partial authoritative binding or shared identity.
  - Execute continues to use only clean bound sources and isolated worktrees.
stop_conditions:
  - The operation scans for or substitutes an undeclared destination.
  - Existing data can be overwritten, reset, stashed, or deleted.
  - Credentials or provider payloads enter workspace or runtime records.
---

# Create and bind an empty repository safely

## Objective

Implement the confirmed operation using existing path, binding, identity, and
atomic-publication primitives.

## Work

Validate the explicit destination and branch, reject existing or unsafe paths,
stage an empty Git repository, verify its identity and cleanliness, write the
host-local binding, and register shared identity through the canonical
projection operation. Define interruption behavior so partial staging never
becomes authoritative.

## Non-goals

Do not create a remote repository, fetch, push, initialize during execute,
change host Git credentials, or modify a bound source after registration.

## Verification

Use GRE-VT-03, GRE-VT-04, and GRE-VT-05.

## Expected evidence

Success, interruption, existing-path, traversal, symlink, incomplete-field,
credential, clean-source, and isolated-worktree results.

## Stop conditions

Stop on any overwrite path, unbounded filesystem search, partial authoritative
state, credential access, or execute-time provisioning.
