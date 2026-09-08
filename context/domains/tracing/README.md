---
kind: domain
status: accepted
title: Tracing and feasibility
slug: tracing
owners: []
sources: []
source_revisions:
  - wrapper: HEAD
    commit: current
    basis: current-wrapper
generated_at: 2026-09-04T00:00:00Z
review_date: 2026-12-04
freshness: accepted-from-current-wrapper
assumptions:
  - The planner runs after Gate 1 (intent approval) and writes the plan when feasible.
  - The feasibility check is coordinator judgment, not an engine verb.
unknowns: []
contradictions: []
acceptance:
  state: accepted
  accepted_at: 2026-09-04
  accepted_by: maintainer
workflows:
  - .context-circuit/docs/planning.md
---

# Tracing and feasibility

## Summary

The bridge between an approved intent and its plans. On intent approval the
coordinator spawns the **planner** — one child per repository in scope —
which reads the *real code* and writes that repository's plan when the look is
feasible. The coordinator then
runs a **feasibility check** over the finding and decides whether to publish
those plans, stop, or ask the human a scope question (INV-INTENT-02). Route the
post-approval "read the code and see if this is buildable" step here. Owned by the
`cc-trace` skill, the `.context-circuit/agents/planner.md` role, and the feasibility-check rule in
`invariants.yaml`.

## Scope

Inside: the planner role and its spawning (one child per repository, in
parallel, no lease, no worktree), the finding it records, and the coordinator's
feasibility judgment and its three outcomes.

Outside: drafting/approving the intent itself ([intent and Gate 1](../intent/README.md)),
how an authorized plan is written and reviewed ([plan review](../plan-review/README.md)),
the tier ladder ([assurance](../assurance/README.md)), and the *execution-time*
grounding scan that briefs the worker ([repository grounding](../repository-grounding/README.md)) —
a different phase (see the disambiguation below).

## Behavior

**Planner (INV-INTENT-02).** When the human approves an intent (Gate 1), the
coordinator spawns the planner: **one child per repository** named in the
intent's scope, in parallel, with **no lease and no worktree** (it writes only
workspace plan files). Each planner reads the real code and, when the look is
feasible, writes that repository's `PLAN.md` and `plan.yaml`:

- tasks and bounded **paths** of where the change lands;
- concrete **risks**;
- the **runnable done-checks** that prove each outcome criterion against the real
  code (earned here, *after* approval — never authored on the intent, and **not**
  folded into the frozen `contract_digest`);
- a **tier signal** (a fail-upward input to consequence tiering — see
  [assurance](../assurance/README.md));
- a **feasibility** finding, and any **out-of-scope reach** (a change that would
  have to modify a repository or area outside the intent's coarse scope);
- **open questions** for the human, which the coordinator relays with a disposition:
  intent-level (stop and revisit Gate 1), plan-level (carry into the plan), or
  already answered (apply without asking again).

The finding is recorded at `intent/<id>/finding.yaml`. The coordinator does not
read the target repository and does not rewrite the plan the planner wrote.

**Feasibility check (INV-INTENT-02).** Before any plan is rewritten the coordinator
judges the planner's finding — this is **coordinator judgment, not an engine
verb**, and a **quality gate, not a safety gate**:

- **Feasible** → it sets the consequence tier, updates the human-facing
  `INTENT.md` status so it matches a completed feasible look, and derives the
  plan(s) **in the same turn**, with no second human gate.
- **Intent revision required** → it does not write a plan; it updates the intent's
  human-facing questions and re-enters Gate 1 when the approved decision changes.
- **Not feasible** → it **stops**, explains the blocker, and hands the decision to
  the human.
- **Out-of-scope reach** → a required change that must modify a repository or area
  beyond the intent's coarse scope is **surfaced to the human** only when the plain
  request did not already authorize that area; a changed approved scope is an
  intent revision before planning.

The feasibility check **fails upward**: when it cannot conclude the change is
buildable or an intent-level question remains unresolved, it stops rather than
proceed. Because it is a quality gate, it never
substitutes for delivery: concrete **scope-safety is settled at Gate 2**
([delivery](../delivery/README.md)), not here.

### Ratifying one plan or stacked plans

After feasibility, the coordinator ratifies the planner's task partition by
looking at the actual execution boundary. Every derived plan names **exactly one
repository**. A bounded change with one worker lifecycle and one independent
verification boundary in one repository stays **one plan with embedded tasks**,
even when the tasks are ordered. Independent execution or verification
boundaries, meaningful dependencies between partitions, distinct failure
surfaces, **or a multi-repository intent** justify **multiple stacked plans**.
An intent whose scope covers two repositories produces at least two plans, one
per repository, with `plan_dependencies` when order matters. Each stacked plan
keeps a bounded path mapping and its own worker/verifier lifecycle.

The partition is not chosen from plan count or assurance tier. Standard work in
one repository may remain one plan, and several plans may share the same
approved intent. Task `depends_on` expresses order inside one plan;
`plan_dependencies` expresses the acyclic ordering between stacked plans. The
coordinator records the human-readable decomposition rationale in the derived
plan's `context_grounding.decisions` and readable `PLAN.md`, while each
inter-plan edge carries its own reason. Anchors and initial task paths are advisory
evidence, not a worker write allowlist. Necessary intent-consistent expansion in the
same repository is recorded and independently checked over the complete diff; a
second-repository or approved-decision expansion stops as a coordinator finding.
This uses the existing contracts and does not add a second plan-approval gate.

## Tracing is not the execution-time grounding scan

Two different phases both "read the repository"; keep them distinct:

- **Planning look** (this page) runs **after approval**. Its subject is
  *whether and how the change is buildable* — file/call-site map, risks,
  done-checks, tier signal, feasibility. It produces the **plan files** and
  `intent/<id>/finding.yaml`.
- **Repository grounding** ([repository grounding](../repository-grounding/README.md))
  runs **at execution setup**, inside the worker's worktree. Its subject is *how to
  write code in this repository* — the repo's own `AGENTS.md`/`CLAUDE.md`/skills —
  and it produces the **worker brief**.

They have different timing, different subjects, and different outputs. Neither
replaces the other.

## Interfaces

- Trigger: intent approval (Gate 1) automatically spawns the planner
- Role: `.context-circuit/agents/planner.md` (one child per repository, never talks to the human)
- Record: `intent/<id>/finding.yaml` plus, when feasible, `plans/<id>/PLAN.md` and `plan.yaml`
- Coordinator outputs: set tier + publish the child's plans · stop-and-explain · scope question to the human

## Data

The finding per repository: `intent`, `repository`, `feasible`,
`feasibility_outcome`, `out_of_scope_reach`, `open_questions`, `tier_signal`,
`revision`. The plan carries tasks, bounded paths, risks, and done-checks.

## Constraints and edge cases

The planner needs **no lease and no worktree** — it writes only workspace plan
files and the finding. Its done-checks and any criterion it surfaces live in the
plan; they are **never** folded into the frozen `contract_digest` (they are earned
post-approval and must not silently re-enter Gate 1). Intent detail
(`intent/<id>/detail/`) is pre-approval authoring evidence for later planning, not a tracing output and not a second gate; the planner remains the first code read.
Whether post-approval
done-checks should ever bind into the frozen contract is an open question recorded
in the design's `risks-and-open-questions`; current behavior keeps them out. If a
required child cannot be created, the route stays `host-blocked` and
the coordinator does not author a stand-in plan.

## Implementation references

- `.agents/skills/cc-trace/SKILL.md`, `.context-circuit/agents/planner.md`
- `.context-circuit/wrapper/contracts/invariants.yaml`: INV-INTENT-02 (owner map: `planner_role`, `feasibility_check`)
- `.context-circuit/wrapper/contracts/schemas/intent-contract.yaml` (the `planner_finding` reference block)

## Verification

`sh test/acceptance.sh` (feasibility suite — pins authorization plus the
deliberate absence of a scope gate; tracing/manifest contract assertions).

## Acceptance notes

Accepted 2026-09-04 (maintainer) to close a knowledge gap: v1.0 added post-approval
tracing and the feasibility check, but no domain page owned them.
