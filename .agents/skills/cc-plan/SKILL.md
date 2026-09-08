---
name: cc-plan
description: Publish the plan a planner child already wrote, or conversationally review a draft plan without approving or executing it.
---

## Publish a plan from an approved intent

A v1.0 plan is the **derivation of an approved intent**, not the thing the human
approves (Context Circuit v1.0, Mechanism 1). Before a writing change is planned
there must be an approved parent intent (`cc-intent`, Gate 1). If none exists,
author the intent first — do not create a plan that invents its own goal or scope.

The planner child (`cc-trace`) already read the real code and, when the look was
feasible, wrote `PLAN.md` and `plan.yaml`. Plan publication is then automatic —
the coordinator's next action **in the same turn** after a feasible planner with
no intent-level questions, not a gate — so the human hears "here's the
breakdown," never "approve this plan." Writing or publishing those plans does not start execution.
If the look is not feasible or an intent-level question remains,
write no plan.

If the approved intent is **Explore**, do not create a plan of record. Explore is
the planless `cc-pair` path: the user supervises the work live, and there is no
independent verifier or completion record to derive.

**Do not author a second copy.** Do not read the target repository. Do not glob it.
Validate and index the files the planner wrote. If no plan and no finding exist
yet, run `cc-trace` rather than improvising a code read inline.

When `intent/<id>/detail/` exists, it is the **confirmed shape** of what to build so
plans and tasks follow those topics. It does this without replacing the
post-approval read of the real code, is not a second approval, and is not part of
`contract_digest`.

## Ratify the task partition

Every derived plan names **exactly one repository**; a plan that lists two or more
is invalid. Keep a **one-plan** derivation when all tasks share one bounded
execution/change surface in that one repository, one worker lifecycle, and one
independent verification boundary. Put the tasks in that plan and use task
`depends_on` for their intra-plan order; several tasks do not require several plans.

Derive **stacked plans** when the partitions are independently executable or
independently verifiable, when there are meaningful dependencies between them,
when their failure surfaces need to be isolated, **or when the intent's scope
covers two or more repositories**. Give each plan that one repository and bounded
paths, record an acyclic `plan_dependencies` edge with a reason, and preserve one
worker plus one independent verifier lifecycle for every Standard/Critical plan.
Refuse collapsing two repositories into one plan.

The coordinator may add a `plan_dependencies` reason across children. It must not
rewrite task bodies, paths, or checks the planner earned. An unresolved intent-level question
blocks publication: return to `cc-intent`. A
`plan-resolution` question is already in the plan. An `already-answered` finding
is applied without asking again.

Task paths are advisory grounding, not an exhaustive worker write allowlist. A
necessary intent-consistent path discovered in the same repository may proceed
when recorded in the worker handoff and checked over the complete diff. A
required second repository or approved-decision change stops as a coordinator
finding.

For one plan, run `plan-validate`, `intent-authorized`, and `plan-index-upsert`.
For a stack the planner staged as `@plan:` fragments, submit one
`plan-stack-materialize` request. It reserves a consecutive id range by
invocation identity, stages readable and canonical artifacts, validates
dependency/coverage/current-authorization/index state, and publishes all plans and index rows or none.
A successful same-invocation retry returns the original
ids. Writing or materializing plans does not execute them.

## Authorization preflight (INV-INTENT-02)

Run `intent-authorized . <plan-id>` as a preflight. If it reports `authorized: yes`,
the plan derives from an approved intent whose criteria are unchanged — no separate
plan approval, and execution may follow (INV-EXEC-01). If it reports `authorized: no`
with `reason: CRITERIA_CHANGED`, the intent's criteria were edited after approval:
the change re-enters Gate 1 — take approval again on the changed criteria via
`cc-intent`. A missing or unapproved intent is likewise unauthorized; author or
approve the intent first. There is **no** automated scope gate: scope-safety is
settled at delivery (Gate 2), the planner reports where the change actually lands, and
a required change beyond a bound scope was surfaced as a question by the feasibility
check in `cc-trace` before publication. The authorization check also runs again at
execution start.

One approved intent may yield **one or more** stacked plans, each naming the same intent
and exactly one repository. Deriving a plan does not execute it.

## Grounding and unupdated Product Knowledge (INV-KNOWLEDGE-02)

A later plan may start even if Product Knowledge from a prior plan has not yet
been updated. Do not wait on knowledge debt before grounding. Explore has no
plan preflight — use the planless `cc-pair` path instead of creating a plan
that would bypass the Explore assurance model.

## Review a plan

`review plan X` is a non-executing, status-preserving discussion. Walk through
request coverage, task detail, repository/path mapping, dependencies, context
references, acceptance, verification, assumptions, open questions, and risks.
When the human resolves a question, corrects a requirement, changes scope, or
asks for more detail, update the draft plan content and continue. Review never
changes plan status or executes. Authorization comes from the approved intent
(re-checked at execution start); there is no separate plan-approval request.
