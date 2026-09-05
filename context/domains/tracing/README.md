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
  - Tracing runs after Gate 1 (intent approval) and before any plan is written.
  - The feasibility check is coordinator judgment, not an engine verb.
unknowns: []
contradictions: []
acceptance:
  state: accepted
  accepted_at: 2026-09-04
  accepted_by: maintainer
workflows:
  - docs/planning.md
---

# Tracing and feasibility

## Summary

The bridge between an approved intent and its plans. On intent approval the
coordinator spawns the **tracer** — one read-only child per repository in scope —
which reads the *real code* and reports a **trace manifest**. The coordinator then
runs a **feasibility check** over those findings and decides whether to derive
plans, stop, or ask the human a scope question (INV-INTENT-02). Route the
post-approval "read the code and see if this is buildable" step here. Owned by the
`cc-trace` skill, the `agents/tracer.md` role, and the feasibility-check rule in
`invariants.yaml`.

## Scope

Inside: the tracer role and its read-only spawning (one child per repository, in
parallel, no lease, no worktree), the trace-manifest contents, bounded freshness
reuse, and the coordinator's feasibility judgment and its three outcomes.

Outside: drafting/approving the intent itself ([intent and Gate 1](../intent/README.md)),
how an authorized plan is written and reviewed ([plan review](../plan-review/README.md)),
the tier ladder ([assurance](../assurance/README.md)), and the *execution-time*
grounding scan that briefs the worker ([repository grounding](../repository-grounding/README.md)) —
a different phase (see the disambiguation below).

## Behavior

**Tracing (INV-INTENT-02).** When the human approves an intent (Gate 1), the
coordinator spawns the tracer: **one read-only child per repository** named in the
intent's scope, in parallel, with **no lease and no worktree** (it never writes).
Each tracer reads the real code and reports a `trace-manifest.yaml`:

- a **file / call-site map** of where the change lands;
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

The manifest is durable grounding evidence at `intent/<id>/trace/<repo>.yaml`,
reused later with a **bounded freshness check** against current code rather than
re-traced from zero.

**Feasibility check (INV-INTENT-02).** Before any plan is written the coordinator
judges the tracer's findings — this is **coordinator judgment, not an engine
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

After feasibility, the coordinator ratifies the manifest's `task_partition` by
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
inter-plan edge carries its own reason. This uses the existing contracts and
does not add a second plan-approval gate.

## Tracing is not the execution-time grounding scan

Two different phases both "read the repository"; keep them distinct:

- **Tracing** (this page) runs **after approval, before planning**. Its subject is
  *whether and how the change is buildable* — file/call-site map, risks,
  done-checks, tier signal, feasibility. It produces the **trace manifest**.
- **Repository grounding** ([repository grounding](../repository-grounding/README.md))
  runs **at execution setup**, inside the worker's worktree. Its subject is *how to
  write code in this repository* — the repo's own `AGENTS.md`/`CLAUDE.md`/skills —
  and it produces the **worker brief**.

They have different timing, different subjects, and different outputs. Neither
replaces the other.

## Interfaces

- Trigger: intent approval (Gate 1) automatically spawns the tracer
- Role: `agents/tracer.md` (read-only, one child per repository, never talks to the human)
- Record: `intent/<id>/trace/<repo>.yaml` (`trace-manifest.yaml`, `schema_version` 1)
- Coordinator outputs: set tier + derive plans · stop-and-explain · scope question to the human

## Data

The trace manifest per repository: `intent`, `repository`, `file_map`,
`task_partition`, `risks`, `done_checks`, `tier_signal`, `feasible`,
`out_of_scope_reach`, `open_questions`, `completeness_proof`, `recorded_at`,
`freshness_checked_at`.

## Constraints and edge cases

The tracer is strictly read-only — no lease, no worktree, no product write. Its
done-checks and any criterion it surfaces live in the manifest and the derived
plan; they are **never** folded into the frozen `contract_digest` (they are earned
post-approval and must not silently re-enter Gate 1). Intent detail
(`intent/<id>/detail/`) is pre-approval authoring evidence for later planning, not a tracing output and not a second gate; tracing remains the first code read.
Whether post-approval
done-checks should ever bind into the frozen contract is an open question recorded
in the design's `risks-and-open-questions`; current behavior keeps them out. If a
required child cannot be created read-only, the route stays `host-blocked` and
read-only rather than self-tracing.

## Implementation references

- `.agents/skills/cc-trace/SKILL.md`, `agents/tracer.md`
- `wrapper/contracts/schemas/trace-manifest.yaml`
- `wrapper/contracts/invariants.yaml`: INV-INTENT-02 (owner map: `tracer_role`, `trace_manifest`, `feasibility_check`)
- `wrapper/contracts/schemas/intent-contract.yaml` (the `trace_manifest` reference block)

## Verification

`sh test/acceptance.sh` (feasibility suite — pins authorization plus the
deliberate absence of a scope gate; tracing/manifest contract assertions).

## Provenance

Authored from the current wrapper for Context Circuit v1.0. This phase replaced
the earlier design's spec-adversary + automated scope-envelope check: the adversary
and envelope are gone; post-approval tracing plus a coordinator feasibility check
took their place. Grounds on the shipped contract; raw `sources/` was not scanned.

## Acceptance notes

Accepted 2026-09-04 (maintainer) to close a knowledge gap: v1.0 added post-approval
tracing and the feasibility check, but no domain page owned them.
