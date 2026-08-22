---
schema_version: 2
id: GRE-001
plan: greenfield-repository-create-empty
status: draft
repository: context-circuit-source
paths:
  - wrapper/contracts/invariants.yaml
  - wrapper/contracts/routes.yaml
  - wrapper/contracts/context-sets.yaml
  - wrapper/contracts/schemas/workspace.yaml
  - docs/gates.md
  - test/routing/fixtures.yaml
  - test/routing/test-router.sh
depends_on: []
acceptance: [GRE-AC-01, GRE-AC-03, GRE-AC-04]
verification: [GRE-VT-01, GRE-VT-02]
expected_evidence:
  - One distinct action, context mapping, human gate, and stable blocker set.
  - Card fixtures with complete target, branch, existing-path, and effect fields.
  - Execute fixtures that never select repository creation.
stop_conditions:
  - The action overloads repository-bootstrap or execute-plan.
  - Card fields are invented after confirmation.
  - Route metadata becomes a second authorization mechanism.
---

# Define the create-empty action and confirmation contract

## Objective

Specify greenfield repository creation as a separate, current human gate with
complete inputs, effects, and safe failure outcomes.

## Work

Extend existing route, invariant, context, workspace, and card owners. Define
the action identifier, exact confirmation, destination and branch fields,
existing-path check, immediate effects, later execution effects, and blockers.
Keep execute unable to select the creation operation.

## Non-goals

Do not implement filesystem creation, change clone/bootstrap semantics, create
a remote repository, or treat route eligibility as confirmation.

## Verification

Use GRE-VT-01 and GRE-VT-02.

## Expected evidence

Contract and fixture diffs showing a distinct action, exact card, safe
clarification for missing fields, and fail-closed execute preflight.

## Stop conditions

Stop if creation can occur through execute, existing bootstrap is weakened, or
the card cannot fully determine the resulting shared and local state.
