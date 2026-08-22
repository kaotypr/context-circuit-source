---
schema_version: 2
id: AGF-001
plan: approval-gate-ux-performance
status: ready
repository: context-circuit-source
paths:
  - wrapper/runtime/engine.sh
  - wrapper/adapters/README.md
  - wrapper/adapters/WORKFLOW.md
  - docs/gates.md
  - docs/planning.md
  - test/routing/fixtures.yaml
  - test/routing/test-router.sh
depends_on: []
acceptance: [AGF-AC-01, AGF-AC-05]
verification: [AGF-VT-01, AGF-VT-04]
expected_evidence:
  - Approval card wording that explicitly separates confirmation from approval, commit, and execution.
  - Routing fixtures proving the initial request remains gated and the exact confirmation is authorized.
  - No policy duplication outside the canonical route and gate owners.
stop_conditions:
  - The change bypasses or weakens the session-bound confirmation gate.
  - The initial request mutates plan, task, runtime, Git, lease, or worktree state.
  - Host-specific instructions become a second route or lifecycle authority.
---

# Clarify the approval card and two-stage human interaction

## Objective

Make the first response understandable and bounded: it should show what would
change, state that nothing has changed yet, and provide the exact confirmation
needed for the current target and session.

## Work

Adjust the canonical card and the thin conversational guidance where needed.
Preserve the existing `present-approval-card` and `approve-plan` route split,
the separate maintainer commit gate, the separate execution trigger, and the
invariant ownership map. The first-turn card must not imply that confirmation
commits Git or starts `Run approved plan`. Add or refine fixtures for the
initial approval request, exact confirmation, and ambiguous approval language.

## Non-goals

Do not auto-approve, auto-commit, start execution, acquire a lease, create a
worktree, or make host adapters authoritative.

## Verification

Use AGF-VT-01 and AGF-VT-04.

## Expected evidence

Updated card/guidance text, routing fixture results, and a diff showing no
runtime or plan mutation from the pre-confirmation route.

## Stop conditions

Stop on any authorization regression, route divergence, or wording that could
make approval appear to mean commit or execution.
