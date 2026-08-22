---
schema_version: 2
id: ERR-001
plan: engine-owned-runtime-records
status: draft
repository: context-circuit-source
paths:
  - wrapper/contracts/invariants.yaml
  - wrapper/contracts/schemas/session.yaml
  - wrapper/contracts/schemas/context-receipt.yaml
  - wrapper/contracts/schemas/delegation.yaml
  - wrapper/contracts/schemas/handoff.yaml
  - wrapper/contracts/schemas/completion.yaml
  - wrapper/manifest.yaml
  - test/contracts/test-contracts.sh
depends_on: []
acceptance: [ERR-AC-01, ERR-AC-02, ERR-AC-03, ERR-AC-04, ERR-AC-07]
verification: [ERR-VT-01]
expected_evidence:
  - Constructor inputs and generated fields are explicit for every record.
  - Ownership graph and launch-validation relationships have one canonical owner.
  - Legacy compatibility cannot grant ownership or authorization.
stop_conditions:
  - A record field is supplied by an unvalidated prompt when the engine can derive it.
  - Validation evidence becomes a human gate or execution authorization.
  - Provider payloads, transcripts, credentials, or host auth state enter records.
---

# Define atomic runtime record operations

## Objective

Specify engine-owned constructor inputs, generated metadata, atomic publication,
ownership relationships, and compatibility behavior for the complete runtime
record set.

## Work

Extend existing schemas and invariants rather than adding a parallel record
layer. Distinguish caller-provided intent from engine-derived timestamps,
versions, revisions, bytes, digests, ancestry, and ownership. Define staged,
valid, published, interrupted, and legacy-readable outcomes.

## Non-goals

Do not implement constructors, launch children, change human gates, or persist
host/provider state beyond the bounded `host_evidence` shape.

## Verification

Use ERR-VT-01.

## Expected evidence

Schema and ownership-map diffs, complete/failing record fixtures, compatibility
rules, and proof that validation evidence cannot authorize work.

## Stop conditions

Stop if fields remain placeholder-capable, ownership has multiple authorities,
or compatibility silently upgrades legacy evidence into current authority.
