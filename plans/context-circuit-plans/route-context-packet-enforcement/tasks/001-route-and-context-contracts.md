---
schema_version: 2
id: RCP-001
plan: route-context-packet-enforcement
status: done
repository: context-circuit-source
paths:
  - wrapper/contracts/routes.yaml
  - wrapper/contracts/context-sets.yaml
  - wrapper/contracts/tier0.yaml
  - wrapper/contracts/invariants.yaml
  - wrapper/contracts/schemas/context-receipt.yaml
  - wrapper/manifest.yaml
  - test/routing/fixtures.yaml
  - test/contracts/test-contracts.sh
depends_on: []
acceptance: [RCP-AC-01, RCP-AC-02, RCP-AC-04, RCP-AC-05]
verification: [RCP-VT-01, RCP-VT-02]
expected_evidence:
  - A complete route-to-registered-context-set mapping.
  - Exact path, conditional evidence, prohibition, and budget fields per packet.
  - Fixtures for initialization, drafting, review, approval, execution, and recovery.
stop_conditions:
  - A probe, phase, action, and context-set id remain implicitly interchangeable.
  - A packet permits broad sources, sibling workspaces, or unrelated plans.
  - Context metadata can authorize a gate or competing route.
---

# Align routes with registered context packets

## Objective

Make the route contract select exactly one valid packet with explicit paths,
conditional evidence, prohibitions, and a hard byte budget.

## Work

Define the mapping from every Stage A/Stage B result to a registered context
set. Complete initialization and plan-draft packet contracts, normalize
product-source selections, and extend receipts with the selected packet
identity and measured references without treating receipts as authorization.

## Non-goals

Do not implement loading, add a second router, generate runtime child records,
or broaden source access.

## Verification

Use RCP-VT-01 and RCP-VT-02.

## Expected evidence

Contract diffs and routing fixtures proving every emitted set exists and that
nearby named, unnamed, future, read-only, gated, and blocked intents remain
distinct.

## Stop conditions

Stop if route and packet identifiers remain ambiguous, if packet selection can
grant authorization, or if a bounded allowlist cannot be expressed.
