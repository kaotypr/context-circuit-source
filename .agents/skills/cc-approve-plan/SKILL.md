---
name: cc-approve-plan
description: Approve a draft plan after human confirmation without starting execution.
---

# Approve a plan

Use this skill when the user asks to approve a named draft plan. `cc-approve-plan`
is the named plan-approval skill. It satisfies the human `plan-approval` gate.
It does not execute the plan.

## Route reads

Use the `planning.review_approval` subroute in
`docs/agent-workspace-workflow.md#route-read-manifests`. The plan-approval
confirmation is the local action-time guard: missing plan contracts or
unresolved contradictions remain blockers, and this skill never starts
execution.

## Preflight

Read the named plan's canonical `plan.yaml`, included tasks, declared
dependencies, and any latest `cc-review-plan` outcome. Confirm:

- the plan identifier is known and not contradictory;
- required plans-contract metadata, tasks, and acceptance criteria are present;
- remaining risks, open questions, and unresolved contradictions are visible.

A prior `cc-review-plan` run is not required. Still preflight, surface remaining
risks and open questions in the confirmation, and stop for revision when the
artifact is incomplete or blocked by a contradiction the human has not resolved.

For a newly generated bundle, confirm that the compact `plan.yaml`,
`overview.md`, and `tasks/` set is complete. Before any task can project to
`ready`, the active host must validate each new or repaired task's initial
front matter with the approved deterministic YAML/schema capability: delimiters,
YAML, required fields, enums, nested shapes, and the
`context-circuit.task` schema. Repair and rerun all checks on failure. Existing
historical bundles may contain optional specialist companions and remain
readable without migration; those companions do not replace `plan.yaml`
authority.

This validation is an approval-readiness precondition, not a caller
responsibility: the approval path invokes it for every new or repaired task
before projecting any task to `ready`. If one task fails, leave the canonical
plan status and every task projection unchanged and report the readiness block.

Report what will change and what will not:

- canonical `plan.yaml` status will become `approved`;
- included task status will become `ready` in one idempotent bulk operation;
- no lease is claimed, no worktree is created, and `cc-run-plan` does not start;
- runtime ownership does not change.

## Confirm

Change `draft` → `approved` only after explicit human confirmation in the
current session. A complete document, passing checks, or a positive review does
not approve the plan.

If the human declines, the plan stays `draft` and execution remains unstarted.

## On confirmed success

1. Set plan status to `approved`.
2. Reconcile included task status to `ready` in one idempotent bulk operation.
   Rewrite only a stale task `status` field. Preserve other task metadata.
3. Leave execution unstarted: no lease claim, no worktree creation, and no
   `cc-run-plan` side effect.

Schema validation is only a readiness precondition. It never substitutes for
the human plan-approval gate, execution evidence, independent verification, or
the status-change gate.

## Refuse or observe

- Already `done`: refuse. Do not invent or rewrite status.
- Unknown or contradictory identifier: refuse.
- Missing required metadata, tasks, or acceptance criteria from the plans
  contract: refuse as incomplete. Stop for revision.
- Unresolved contradiction: stop for revision; do not approve.
- Already `approved`: report that fact and recommend `cc-run-plan` without
  rewriting status.

## Safety

Must not merge, push, publish, deploy, or change runtime ownership. A child
worker or verifier cannot approve a plan. Only the root coordinator, after
current-session human confirmation, may write canonical plan status.
