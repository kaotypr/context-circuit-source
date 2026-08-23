---
schema_version: 2
id: MHS-001
plan: multi-host-agent-support
status: done
repository: context-circuit-source
paths:
  - wrapper/contracts/invariants.yaml
  - wrapper/contracts/routes.yaml
  - wrapper/contracts/schemas/session.yaml
  - wrapper/contracts/schemas/delegation.yaml
  - wrapper/contracts/schemas/handoff.yaml
  - docs/host-capabilities.md
  - docs/configuration.md
  - test/contracts/test-contracts.sh
  - test/routing/fixtures.yaml
  - test/behavior-matrix/matrix.yaml
depends_on: []
acceptance: [MHS-AC-01, MHS-AC-02]
verification: [MHS-VT-01, MHS-VT-02, MHS-VT-07]
expected_evidence:
  - One provider-neutral host capability shape and supported-host matrix.
  - Session/delegation/handoff ownership updates with no lifecycle duplication.
  - Contract and routing fixtures for host success, unavailable, and blocked states.
stop_conditions:
  - A host-specific rule would become a second router, lifecycle, or authorization owner.
  - A proposed field would persist credentials, provider payloads, or opaque transcripts.
  - Host capability is treated as authorization or missing child support is silently bypassed.
---

# Define the cross-host capability contract and acceptance matrix

## Objective

Define the provider-neutral host identity and capability evidence needed to
route Codex CLI, Claude Code, and Cursor Agent CLI through the same Context
Circuit contracts.

## Work

Specify the bounded host fields needed in session, delegation, and handoff
evidence: host identifier, observed version when available, instruction
surface, root/child capability, verifier capability, resume capability,
permission mode, and offline fallback. Keep the canonical route and lifecycle
owners unchanged. Update the host matrix, configuration guidance, route
fixtures, behavior-matrix labels, and contract checks so `host-blocked`,
`child-blocked`, `disabled`, `denied`, and `unavailable` remain distinguishable.

## Non-goals

Do not implement host commands, invoke providers, add an SDK, or encode host
credentials. Do not create a scheduler or make host configuration an
authorization grant.

## Verification

Use MHS-VT-01, MHS-VT-02, and MHS-VT-07.

## Expected evidence

Changed owner files, schema inventory, cross-host matrix, positive and
negative routing fixtures, and a proof that no forbidden credential-shaped
field or provider payload is accepted.

## Stop conditions

Stop on conflicting owner rules, ambiguous host identity, scope expansion into
provider integration, or any request to persist authentication material.

