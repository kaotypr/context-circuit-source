---
kind: domain
status: accepted
title: Plan authorization
slug: plan-authorization
owners: []
sources: []
source_revisions:
  - wrapper: HEAD
    commit: current
    basis: current-wrapper
generated_at: 2026-09-02T00:00:00Z
review_date: 2026-12-02
freshness: accepted-from-current-wrapper
assumptions:
  - The single upstream human gate is on the intent (Gate 1); a plan carries no second approval.
unknowns: []
contradictions: []
acceptance:
  state: accepted
  accepted_at: 2026-09-02
  accepted_by: maintainer
workflows:
  - docs/planning.md
---

# Plan authorization

## Summary

A plan is not separately approved. Its authorization is derived from its parent
**intent**: the human approves the intent once (Gate 1, `draft → approved`,
INV-APPROVE-01), and every plan that derives from that intent, whose approved
criteria are unchanged, is authorized to execute with no second gate
(INV-INTENT-02). The authorization check is **scope-free** — it does not compare
repositories or paths; scope-safety is settled at delivery (Gate 2). Plan status
is only `draft` or `done` — there is no intermediate `approved` plan status. Route
"run plan `<id>`" through the `intent-authorized` check, with no separate plan
approval step.

## Scope

Inside: how a plan becomes authorized to run (the scope-free `intent-authorized`
check), the `draft → done` plan status, and the synchronized task-status projection.

Outside: the intent gate itself (Gate 1, see the intent skill), execution and the
worker loop (see [plan-execution](../plan-execution/README.md)), plan
authoring/review (see [plan-review](../plan-review/README.md)), and completion.

## Behavior

Authorization is a property of the intent, not an act on the plan. A plan derived
from an approved intent, whose frozen acceptance criteria are unchanged since
approval, may execute; the runtime re-runs the `intent-authorized` check at
execution start so that a criteria change after planning (the intent's
`contract_digest` no longer matches) re-gates to a human rather than proceeding
(INV-EXEC-01, INV-INTENT-02). A plan whose intent is unapproved, missing, or whose
criteria have moved is refused, never rubber-stamped. This check makes no
scope/path comparison — that safety is settled at delivery (Gate 2); the one
safety-critical automated check on the path here is consequence-tier
classification (INV-ASSURE-01).

`plan.yaml` owns the human status; it is `draft` from creation until the plan
completes, then `done`. Included task projections move in sync and are never a
second authority (INV-PLAN-01). Execution status (running, verifying, repairing,
verified, failed, blocked) is runtime evidence and never replaces plan status: a
verified execution is evidence, not a `done` plan.

One approved intent may authorize one plan with multiple embedded tasks or a stack
of plans. In the stacked case, `plan_dependencies` controls ordering and carries
the reason for each edge; it does not add a human decision. Every Standard/Critical
plan still has its own bounded worker and independent verifier lifecycle. Combining
the delivered plans into one change-set candidate is a separate delivery concern
and does not erase those per-plan execution records. Plan count is not a proxy for
assurance tier. Optional intent detail has no approval status of its own; plan
derivation remains automatic after Gate 1 and feasibility.

There is no confirmation card and no hidden confirmation token. A vague "yes" is
not an intent approval; the decision is an explicit conversational act bound to the
named intent, handled upstream (Gate 1). Approving an intent and asking to build in
one turn is honored as two sequential explicit actions; a human who wants to lock an
intent without building yet may still separate the two.

On this `product-source` checkout, committing plan state is an ordinary source-only
commit that an explicit user request may authorize (see `AGENTS.md`); it is not a
gated card. Registered product repositories still require their own delivery and
publication gates.

## Interfaces

- Human request: "Run plan `<id>`" (authorized via its approved intent, no separate approval)
- Canonical status: `plan.yaml` `status` (`draft`, `done`)
- Authorization preflight: `intent-authorized` (scope-free)

## Constraints and edge cases

The authorization check fails upward: an unapproved or missing intent, or
acceptance criteria changed since approval (digest mismatch), re-gates rather than
authorizing. A plan whose intent is not approved cannot execute. An already-`done`
plan is not re-executed. Unrelated dirty working-tree files are preserved.

## Implementation references

- `.agents/skills/cc-plan/SKILL.md`, `.agents/skills/cc-execute/SKILL.md`
- `wrapper/runtime/engine.sh`: `cc_intent_authorized`, `cc_execution_begin`, `cc_plan_set_status`
- `wrapper/contracts/schemas/plan.yaml`, `wrapper/contracts/schemas/intent-contract.yaml`
- `wrapper/contracts/invariants.yaml`: INV-APPROVE-01, INV-INTENT-02, INV-PLAN-01, INV-EXEC-01

## Provenance

Re-grounded on the current wrapper for Context Circuit v1.0. Supersedes the earlier
`plan-approval` domain: v1.0 removed the plan-level approval gate and the
intermediate `approved` plan status, folding authorization into the intent gate
(Gate 1). Authorization is scope-free (approved intent + unchanged criteria);
scope-safety is settled at delivery (Gate 2), not by an authorization-time scope
check. Raw `sources/` was not scanned.

## Acceptance notes

Accepted 2026-09-02. Replaces the retired `draft → approved` plan gate and the
`cc_plan_approve` engine verb (both removed in v1.0) with intent-derived
authorization. Extends the accepted "separate lifecycle gates" and "intent is the
single upstream gate" decisions in `context/DECISIONS.md`.
