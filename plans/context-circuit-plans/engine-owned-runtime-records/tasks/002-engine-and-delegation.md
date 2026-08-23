---
schema_version: 2
id: ERR-002
plan: engine-owned-runtime-records
status: draft
repository: context-circuit-source
paths:
  - wrapper/runtime/engine.sh
  - agents/writer.md
  - agents/verifier.md
  - test/runtime/test-runtime.sh
  - test/ownership/test-ownership.sh
  - test/security/test-boundaries.sh
depends_on: [ERR-001]
acceptance: [ERR-AC-01, ERR-AC-02, ERR-AC-03, ERR-AC-04, ERR-AC-05]
verification: [ERR-VT-02, ERR-VT-03, ERR-VT-06]
expected_evidence:
  - Atomic constructors emit complete validated records, child-start evidence,
    a commit marker, and a bounded handoff skeleton.
  - Ownership graph validation rejects stale, foreign, partial, and mismatched inputs.
  - Writer and verifier launch projections cannot expand packet scope.
stop_conditions:
  - A commit marker can be observed without every staged record validating.
  - The engine starts a provider or becomes a scheduler.
  - Child role guidance duplicates ownership policy instead of citing contracts.
---

# Generate and validate the ownership graph

## Objective

Implement constructors and graph validation while keeping provider transport
and authorization outside the engine.

## Work

Generate sessions, receipts, delegations, child-start evidence, handoff
skeletons, and completion records from validated inputs. Stage the complete
graph in one transaction directory, validate ancestry, lease, repository,
worktree, host evidence, receipt, scope, permissions, verification, and
handoff, then publish one commit marker only after all records validate. Readers
must treat unmarked or incomplete transactions as non-authoritative. Return a
bounded launch projection.

## Non-goals

Do not invoke a host, schedule work, create a second lifecycle, or let launch
text redefine packet scope.

## Verification

Use ERR-VT-02, ERR-VT-03, and ERR-VT-06.

## Expected evidence

Success, interruption, partial graph, stale digest, foreign owner, prompt
mismatch, writer/verifier permission, and handoff-template results.

## Stop conditions

Stop if the transaction boundary is not recoverable, launch projection can
expand scope, or host/provider concerns enter the runtime authority layer.
