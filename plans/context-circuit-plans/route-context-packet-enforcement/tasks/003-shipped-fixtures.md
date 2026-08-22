---
schema_version: 2
id: RCP-003
plan: route-context-packet-enforcement
status: draft
repository: context-circuit-source
paths:
  - docs/templates/plan.md
  - docs/templates/plan.yaml
  - docs/templates/task.md
  - docs/templates/prd.md
  - docs/planning.md
  - scripts/release-artifact.sh
  - scripts/release-manifest.txt
  - test/release/test-release.sh
depends_on: [RCP-001]
acceptance: [RCP-AC-02, RCP-AC-07]
verification: [RCP-VT-01, RCP-VT-07]
expected_evidence:
  - Released canonical examples cover plan, task, gate, receipt, and delegation shapes.
  - Plan drafting and initialization packets reference only shipped examples.
  - Artifact inspection excludes maintainer plans, customer workspaces, and runtime state.
stop_conditions:
  - A fixture copies a maintainer or customer plan as released product state.
  - Release assembly exposes sources, local bindings, credentials, or runtime evidence.
  - Shipped examples become a second schema or policy owner.
---

# Ship canonical packet fixtures

## Objective

Provide bounded released examples so initialization and plan drafting never
need tests, maintainer plans, or sibling workspaces to discover artifact shape.

## Work

Complete the canonical plan template and add only the minimal released examples
required by packet contracts. Wire them into release and upgrade inventories
and document that schemas remain authoritative when examples differ.

## Non-goals

Do not ship this maintainer plan stack, source evidence, tests, runtime state,
or a full example customer workspace.

## Verification

Use RCP-VT-01 and RCP-VT-07.

## Expected evidence

Packet references, artifact inventory, schema/example consistency checks, and
negative release assertions for maintainer and customer evidence.

## Stop conditions

Stop if an example becomes authoritative, release boundaries broaden, or
drafting still requires a non-shipped example.
