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
  - .context-circuit/docs/planning.md
  - .context-circuit/docs/plan-review.md
---

# Planning and plan review

## Summary

Authoring a grounded plan from a request, and conversationally reviewing a draft
plan without executing it. Both are owned by the `cc-plan` skill. Route "create a
plan for …" and "review plan `<id>`" here.

## Scope

Inside: grounded plan creation, stable plan-id allocation, plan/task structure,
the readable `PLAN.md` and canonical `plan.yaml`, and non-executing review.

Outside: authorization (see
[plan-authorization](../plan-authorization/README.md), derived from the parent
approved intent), execution (see [plan-execution](../plan-execution/README.md)),
and any status change.

## Behavior

Creating a plan does not execute it; it is authorized by its parent approved
intent, scope-free (approved intent + unchanged criteria). `cc-plan` drafts against
route-selected Product Knowledge (`context/INDEX.md`), grounding the plan in
available evidence; missing or contradictory information becomes an explicit
open question, assumption, or risk rather than an invented decision
(INV-PLAN-04). Every task names the repository or repositories it may change and
bounded paths or an explicit repository-wide scope, with explicit dependencies
(INV-PLAN-02). New plans use stable ids `NNNN-<kebab-slug>`; the next id is one
past the highest active plan in the current member's band. Archived plans are
ignored by allocation, and restore refuses an active numeric-prefix collision
(INV-PLAN-03, INV-ARCHIVE-02, INV-MEMBER-01). Member names do not appear in the id or in execution branches
(`cc/<plan-id>/<repository-id>`).

`plan.yaml` owns human plan status (`draft`, `done`); task status is a
synchronized projection and never a second lifecycle authority (INV-PLAN-01).

Reviewing a named plan is a read-only discussion: it walks the original request,
objective/constraints/non-goals, grounding evidence, repository and task
mapping, task order and dependencies, acceptance and verification ids,
assumptions/open questions/risks, and delivery effects. It never changes plan
status or executes. Resolving questions may update the draft; authorization comes
from the approved intent, not a separate plan-approval step. When the host exposes a native
question prompt (`AskUserQuestion`, `AskQuestion`, `request_user_input`), focused
choices may be offered there, but a missing or failed prompt is not
`host-blocked` and question transcripts are not recorded.

## Workflows

- Plan authoring: `.context-circuit/docs/planning.md`
- Plan review: `.context-circuit/docs/plan-review.md`

## Interfaces

- Canonical machine plan: `plan.yaml`; readable plan: `PLAN.md`
- Retrieval catalog for grounding: `context/INDEX.md`
- Active plan index: `plans/INDEX.md`

## Constraints and edge cases

Read only the selected context units, not the whole directory. Review writes
nothing — no `plan.yaml`, task status, lease, runtime, or Git. There is no
eighth shipped skill for review; `cc-plan` carries it.

## Plan brief, PLAN.md, and plan.yaml

Authoring starts from an internal plan-creation brief (request, objective,
repositories, product_knowledge, context_grounding, sources, known_decisions,
open_questions) that grounds scope before any plan is written.

The readable `PLAN.md` carries the status and execution summary, repository
ownership, scope and non-goals, acceptance, verification, risks, expected commits
and delivery, human decisions, and expected Product Knowledge impact. The
canonical `plan.yaml` additionally owns allowed paths, declared effects,
acceptance and verification ids, the repair limit, plan-level stop conditions,
and the context units to reassess at completion.

Request fidelity: a plan removes repetition, not meaning; detail is not discarded
merely to shorten. (Note: `INV-PLAN-03` is the band-scoped never-reused /
next-after-highest-in-band id rule; see `context/DESIGN-DELTAS.md`.)

## Implementation references

- `.agents/skills/cc-plan/SKILL.md`
- `.context-circuit/wrapper/runtime/engine.sh`: `cc_plan_allocate_id`, `cc_plan_validate`,
  `cc_plan_status`, `cc_member_band_resolve`
- `.context-circuit/wrapper/contracts/schemas/plan.yaml`, `.context-circuit/wrapper/contracts/schemas/task.yaml`,
  `.context-circuit/wrapper/contracts/schemas/members.yaml`, `.context-circuit/wrapper/contracts/schemas/member-local.yaml`
- `.context-circuit/wrapper/contracts/invariants.yaml`: INV-PLAN-01, INV-PLAN-02, INV-PLAN-03,
  INV-PLAN-04, INV-MEMBER-01

## Acceptance notes

Accepted 2026-08-24. Broadened from the former review-only `plan-review` draft to
cover plan authoring as well, since the single `cc-plan` skill owns both. Slug
kept as `plan-review`.
