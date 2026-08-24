---
kind: domain
status: accepted
title: Plan approval
slug: plan-approval
owners: []
sources: []
source_revisions:
  - wrapper: HEAD
    commit: 4b8ac0b
    basis: current-wrapper
generated_at: 2026-08-24T00:00:00Z
review_date: 2026-11-24
freshness: accepted-from-current-wrapper
assumptions:
  - Approval is a single explicit conversational gate with no confirmation card or token.
unknowns: []
contradictions: []
acceptance:
  state: accepted
  accepted_at: 2026-08-24
  accepted_by: maintainer
workflows:
  - docs/planning.md
---

# Plan approval

## Summary

Approval is an explicit conversational human gate that changes plan status from
`draft` to `approved` after deterministic readiness checks. It is not a
confirmation card and carries no hidden confirmation token. Route "approve plan
`<id>`" here. Execution is a separate authorization.

## Scope

Inside: the `draft → approved` transition and its readiness checks, and the
synchronized task-status projection.

Outside: execution and the writer loop (see
[plan-execution](../plan-execution/README.md)), plan authoring/review (see
[plan-review](../plan-review/README.md)), and completion.

## Behavior

Approval only changes plan status `draft → approved` after deterministic
readiness checks pass; it creates no worktree, claims no lease, starts no
implementation, and commits no Git (INV-APPROVE-01). Eligibility is not
authorization: a vague "yes" is not approval, and the decision is an explicit
conversational act bound to the named plan.

`plan.yaml` owns the human status; included task projections move in sync and
are never a second authority (INV-PLAN-01).

Only an approved plan may execute, and execution is a separate authorization
from approval unless one request explicitly asks for both — "approve plan `<id>`
and execute it" is two sequential explicit actions (INV-EXEC-01).

On this `product-source` checkout, committing the approved plan state is an
ordinary source-only commit that an explicit user request may authorize (see
`AGENTS.md`); it is not a gated card and there is no
`MAINTAINER_APPROVAL_COMMIT_REQUIRED` token. Registered product repositories
still require their own delivery and publication gates.

## Workflows

- Lifecycle separation and the approval gate: `docs/planning.md`

## Interfaces

- Human request: "Approve plan `<id>`"
- Canonical status: `plan.yaml` `status`

## Constraints and edge cases

Approval never runs deterministic readiness against a nonexistent plan, never
combines silently with execution, and never persists as a reusable token.
Unrelated dirty working-tree files are preserved.

## Readiness checks

The deterministic readiness checks that gate `draft -> approved`:

- the objective is represented accurately;
- repository mappings are explicit;
- no task depends on an unknown repository;
- context conflicts are resolved or accepted as risks;
- acceptance and verification are testable;
- worker and verifier scope are bounded;
- the plan status is currently draft.

## Implementation references

- `.agents/skills/cc-plan/SKILL.md`, `.agents/skills/cc-execute/SKILL.md`
- `wrapper/runtime/engine.sh`: `cc_plan_approve`, `cc_plan_set_status`
- `wrapper/contracts/schemas/plan.yaml`
- `wrapper/contracts/invariants.yaml`: INV-APPROVE-01, INV-PLAN-01, INV-EXEC-01

## Provenance

Re-grounded on the current wrapper at HEAD `4b8ac0b`. The previous-version
`approval-gate-ux-performance` plan that seeded this page was deleted in
`4b8ac0b`; its provenance was retired. Raw `sources/` was not scanned.

## Acceptance notes

Accepted 2026-08-24. Substantively corrected: the earlier draft's two-turn
confirmation card, `commit-approved-plan` card, `MAINTAINER_APPROVAL_COMMIT_REQUIRED`,
and `DIRTY_BASE_BLOCKED` machinery, and the engine functions
`cc_transition_plan_status` / `cc_maintainer_approval_commit_required`, are all
absent from the current wrapper and contradicted INV-APPROVE-01 ("no confirmation
card and no hidden confirmation token"). Extends the accepted "separate lifecycle
gates" decision in `context/DECISIONS.md`.
