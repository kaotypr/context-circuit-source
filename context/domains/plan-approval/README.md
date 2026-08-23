---
kind: domain
status: proposed
title: Plan approval and product-source commit
slug: plan-approval
owners: []
sources:
  - plans/context-circuit-plans/approval-gate-ux-performance/plan.yaml
  - plans/context-circuit-plans/approval-gate-ux-performance/PLAN.md
  - docs/gates.md
source_revisions:
  - plan: approval-gate-ux-performance
    working_tree_status: done
    HEAD_status: approved
    updated_at: 2026-08-22T18:00:00Z
generated_at: 2026-08-23T00:00:00Z
review_date: 2026-09-22
freshness: proposed-from-done-plans
assumptions:
  - Working-tree plan.yaml status done is the selected done-plan evidence.
unknowns: []
contradictions:
  - approval-gate-ux-performance is done in the working tree but still dirty
    versus HEAD, which has status approved. The uncommitted diff is the finish
    status projection (approved→done and task ready→done).
acceptance:
  state: pending
  accepted_at:
  accepted_by:
workflows:
  - docs/gates.md
  - docs/planning.md
---

# Plan approval and product-source commit

## Summary

Plan approval is a two-turn, status-only gate. On this product-source
checkout, confirmed approval may immediately present the existing maintainer
commit card. Execution remains a later explicit request. Route approval,
confirm-approval, and the follow-on commit-approved-plan card here.

## Scope

Inside: `Approve plan <id>`, `Confirm approval of plan <id>`, the canonical
status transition, task projections, and the product-source
`commit-approved-plan` card.

Outside: combining approval with commit or execution, auto-committing, creating
a worktree from a dirty product-source base, and treating
`MAINTAINER_APPROVAL_COMMIT_REQUIRED` as an execution exemption.

## Behavior

`Approve plan <id>` presents a current session-bound card and mutates nothing.
The exact confirmation is `Confirm approval of plan <id>`. That confirmation
calls one canonical transition: `plan.yaml` `draft` → `approved` and included
task projections `draft` → `ready`. Task bodies, file endings, and unrelated
working-tree changes are preserved.

The confirmed path does not create a worktree, claim a lease, start
implementation, or commit Git.

On an instantiated or wrapped workspace, the handoff says approval is complete
and `Run approved plan <id>` is the separate next request.

On `product-source`, if the dirty set is exactly that approval projection,
present the existing commit card in the same session:

- Confirmation: `Confirm commit of the approved plan state.`
- Instantiated or wrapped workspaces do not receive this card.
- Unrelated dirty files remain `DIRTY_BASE_BLOCKED`.
- Running before that commit reports `MAINTAINER_APPROVAL_COMMIT_REQUIRED`.

After that commit, `Run approved plan <id>` is the next explicit request.

## Workflows

- Approval card and maintainer commit card: `docs/gates.md`
- Lifecycle separation: `docs/planning.md`

## Constraints and edge cases

Eligibility is not authorization. Confirmation is bound to the displayed
target and current session. Vague “yes” is not approval. Dirty files outside
the exact status projection must be preserved.

## Implementation references

- `wrapper/runtime/engine.sh` `cc_transition_plan_status`,
  `cc_maintainer_approval_commit_required`
- `wrapper/contracts/invariants.yaml` INV-AUTH-01, INV-AUTH-02, INV-OWN-07
- `docs/gates.md`

## Verification

Done-plan verification IDs AGF-VT-01–AGF-VT-05.

## Provenance

Read only the selected plan `approval-gate-ux-performance` (CC-004) from the
working tree, plus `docs/gates.md` named by that plan. Raw `sources/` was not
scanned. HEAD at generation was `b7a11f3`.

## Acceptance notes

This page is proposed. It extends the accepted “separate lifecycle gates”
decision with the product-source commit-card sequencing; it does not replace
that decision. Human context acceptance is still required.
