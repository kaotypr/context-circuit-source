---
schema_version: 2
id: WPE-001
plan: workspace-projection-and-gate-effects
status: done
repository: context-circuit-source
paths:
  - wrapper/contracts/invariants.yaml
  - wrapper/contracts/routes.yaml
  - wrapper/contracts/schemas/workspace.yaml
  - wrapper/manifest.yaml
depends_on: []
acceptance: [WPE-AC-01, WPE-AC-02, WPE-AC-03, WPE-AC-04, WPE-AC-05]
verification: [WPE-VT-01, WPE-VT-03]
expected_evidence:
  - One canonical owner for identity-region fields, displayed defaults, and gate-effect identifiers.
  - Schema fixtures for matching, mismatched, missing-region, displayed-default, incomplete, and undeclared-effect cases.
  - Stable failure results that do not authorize a route or gate.
stop_conditions:
  - Projection or effect policy is duplicated in a skill, role, or host adapter.
  - Machine-specific repository paths enter shared workspace identity.
  - An effect declaration can grant authorization or bypass confirmation.
  - The contract derives or replaces Product Knowledge from workspace.yaml.
---

# Define identity-projection and gate-effect contracts

## Objective

Define the minimum canonical fields and ownership needed to keep
`workspace.yaml` as identity metadata, keep `context/` as Product Knowledge,
and describe gate effects without creating another router or policy surface.

## Work

Extend the existing workspace, route, invariant, and manifest owners.

Specify the bounded identity region in `WORKSPACE.md`, `PROJECT.md`, and
`INDEX.md`: one contract-owned delimiter per file, the `workspace.yaml`
fields that must appear there (name, mode, identity kind and status, logical
repository keys with optional credential-free URL and default branch), and
the stable `projection-mismatch` result when the region is missing or
disagrees. Authored Product Knowledge outside the region is not a projection
input or output.

Specify identity-acceptance displayed fields: mode, repositories or project
items, roles, and default branches. Omitted fields are proposed defaults on
the card. If `roles` is not already a schema field, add a canonical optional
field. Fields with no default remain incomplete.

Lock effect identifiers as descriptive metadata, not authorization:

- `workspace.accept_identity`
- `workspace.register_repository`
- `repository-bootstrap`
- `repository-create-empty` (reserved; no action or gate in this plan)
- `git.commit` (worktree only)
- `delivery.push`

`git.init` and `git.clone` are not `execute-plan` effects. Delegated effects
must be a subset of approved intent.

## Non-goals

Do not implement rendering, mutate workspace identity, create repositories,
activate create-empty, change lifecycle status, generate Product Knowledge,
or place policy in conversational adapters.

## Verification

Use WPE-VT-01 and WPE-VT-03.

## Expected evidence

Contract diffs, complete and failing fixtures, ownership-map coverage, and
tests showing that effect metadata is descriptive rather than authorization.

## Stop conditions

Stop if the design requires a competing route owner, stores local paths in
shared identity, derives Product Knowledge from identity metadata, or cannot
distinguish current confirmation from later action.
