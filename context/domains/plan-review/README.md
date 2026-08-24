---
kind: domain
status: accepted
title: Planning and plan review
slug: plan-review
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
  - The cc-plan skill owns both plan authoring and plan review.
unknowns: []
contradictions: []
acceptance:
  state: accepted
  accepted_at: 2026-08-24
  accepted_by: maintainer
workflows:
  - docs/planning.md
  - docs/plan-review.md
---

# Planning and plan review

## Summary

Authoring a grounded plan from a request, and conversationally reviewing a draft
plan without approving or executing it. Both are owned by the `cc-plan` skill.
Route "create a plan for …" and "review plan `<id>`" here.

## Scope

Inside: grounded plan creation, stable plan-id allocation, plan/task structure,
the readable `PLAN.md` and canonical `plan.yaml`, and non-executing review.

Outside: approval (see [plan-approval](../plan-approval/README.md)), execution
(see [plan-execution](../plan-execution/README.md)), and any status change.

## Behavior

Creating a plan does not approve or execute it. `cc-plan` drafts against
route-selected Product Knowledge (`context/INDEX.md`), grounding the plan in
available evidence; missing or contradictory information becomes an explicit
open question, assumption, or risk rather than an invented decision
(INV-PLAN-04). Every task names the repository or repositories it may change and
bounded paths or an explicit repository-wide scope, with explicit dependencies
(INV-PLAN-02). New plans use stable ids `NNNN-<kebab-slug>`; the four-digit
sequence is the next after the highest ever allocated and is never reused
(INV-PLAN-03).

`plan.yaml` owns human plan status (`draft`, `approved`, `done`); task status is
a synchronized projection and never a second lifecycle authority (INV-PLAN-01).

Reviewing a named plan is a read-only discussion: it walks the original request,
objective/constraints/non-goals, grounding evidence, repository and task
mapping, task order and dependencies, acceptance and verification ids,
assumptions/open questions/risks, and delivery effects. It never changes plan
status, approves, or executes. Resolving questions may update the draft;
approval stays a separate explicit request. When the host exposes a native
question prompt (`AskUserQuestion`, `AskQuestion`, `request_user_input`), focused
choices may be offered there, but a missing or failed prompt is not
`host-blocked` and question transcripts are not recorded.

## Workflows

- Plan authoring: `docs/planning.md`
- Plan review: `docs/plan-review.md`

## Interfaces

- Canonical machine plan: `plan.yaml`; readable plan: `PLAN.md`
- Retrieval catalog for grounding: `context/INDEX.md`
- Active plan index: `plans/INDEX.md`

## Constraints and edge cases

Read only the selected context units, not the whole directory. Review writes
nothing — no `plan.yaml`, task status, lease, runtime, or Git. There is no
eighth shipped skill for review; `cc-plan` carries it.

## Implementation references

- `.agents/skills/cc-plan/SKILL.md`
- `wrapper/runtime/engine.sh`: `cc_plan_allocate_id`, `cc_plan_validate`,
  `cc_plan_status`
- `wrapper/contracts/schemas/plan.yaml`, `wrapper/contracts/schemas/task.yaml`
- `wrapper/contracts/invariants.yaml`: INV-PLAN-01, INV-PLAN-02, INV-PLAN-03,
  INV-PLAN-04

## Provenance

Re-grounded on the current wrapper at HEAD `4b8ac0b`. The previous-version
`interactive-plan-review` plan that seeded this page was deleted in `4b8ac0b`;
its provenance was retired, along with the previous-version route tokens
`review-plan` / `plan-review` / `clarify-target`, which are absent from the
current engine. Raw `sources/` was not scanned.

## Acceptance notes

Accepted 2026-08-24. Broadened from the former review-only `plan-review` draft to
cover plan authoring as well, since the single `cc-plan` skill owns both. Slug
kept as `plan-review`.
