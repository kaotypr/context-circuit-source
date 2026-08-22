---
schema_version: 2
id: WPE-001
plan: workspace-projection-and-gate-effects
status: draft
repository: context-circuit-source
paths:
  - wrapper/contracts/invariants.yaml
  - wrapper/contracts/routes.yaml
  - wrapper/contracts/schemas/workspace.yaml
  - wrapper/manifest.yaml
depends_on: []
acceptance: [WPE-AC-01, WPE-AC-03, WPE-AC-04, WPE-AC-05]
verification: [WPE-VT-01, WPE-VT-03]
expected_evidence:
  - One canonical owner for workspace projections and gate effects.
  - Schema fixtures for complete, incomplete, mismatched, and undeclared effects.
  - Stable failure results that do not authorize a route or gate.
stop_conditions:
  - Projection or effect policy is duplicated in a skill, role, or host adapter.
  - Machine-specific repository paths enter shared workspace identity.
  - An effect declaration can grant authorization or bypass confirmation.
---

# Define projection and gate-effect contracts

## Objective

Define the minimum canonical fields and ownership needed to derive workspace
summaries and describe gate effects without creating another router or policy
surface.

## Work

Extend the existing workspace, route, invariant, and manifest owners. Specify
which `workspace.yaml` fields project into each required summary, how mismatch
and incomplete confirmation are reported, and how immediate and later effects
are represented. Require delegated effects to be a subset of approved intent.

## Non-goals

Do not implement rendering, mutate workspace identity, create repositories,
change lifecycle status, or place policy in conversational adapters.

## Verification

Use WPE-VT-01 and WPE-VT-03.

## Expected evidence

Contract diffs, complete and failing fixtures, ownership-map coverage, and
tests showing that effect metadata is descriptive rather than authorization.

## Stop conditions

Stop if the design requires a competing route owner, stores local paths in
shared identity, or cannot distinguish current confirmation from later action.
