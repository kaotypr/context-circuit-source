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
INV-APPROVE-01), and every plan that derives from that intent and stays within its
approved **scope envelope** is authorized to execute with no second gate
(INV-INTENT-02). Plan status is only `draft` or `done` — there is no intermediate
`approved` plan status. Route "run plan `<id>`" through the envelope check, with
no separate plan approval step.

## Scope

Inside: how a plan becomes authorized to run (the intent-envelope check), the
`draft → done` plan status, and the synchronized task-status projection.

Outside: the intent gate itself (Gate 1, see the intent skill), execution and the
worker loop (see [plan-execution](../plan-execution/README.md)), plan
authoring/review (see [plan-review](../plan-review/README.md)), and completion.

## Behavior

Authorization is a property of the intent, not an act on the plan. A plan derived
from an approved intent, whose repositories and paths stay within the intent's
approved scope envelope, may execute; the runtime re-runs the envelope check at
execution start (crown jewel 1) so scope drift discovered after planning re-gates
to a human rather than proceeding (INV-EXEC-01, INV-INTENT-02). A plan that exceeds
the envelope — a repository or path the intent never approved — is refused, never
rubber-stamped.

`plan.yaml` owns the human status; it is `draft` from creation until the plan
completes, then `done`. Included task projections move in sync and are never a
second authority (INV-PLAN-01). Execution status (running, verifying, repairing,
verified, failed, blocked) is runtime evidence and never replaces plan status: a
verified execution is evidence, not a `done` plan.

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

- Human request: "Run plan `<id>`" (authorized via the intent envelope, no separate approval)
- Canonical status: `plan.yaml` `status` (`draft`, `done`)
- Authorization preflight: `intent-envelope-check` (crown jewel 1)

## Constraints and edge cases

The envelope check fails upward: any indeterminate comparison (unresolvable path,
missing scope entry, unapproved or missing intent) re-gates rather than authorizing.
A plan whose intent is not approved cannot execute. An already-`done` plan is not
re-executed. Unrelated dirty working-tree files are preserved.

## Implementation references

- `.agents/skills/cc-plan/SKILL.md`, `.agents/skills/cc-execute/SKILL.md`
- `wrapper/runtime/engine.sh`: `cc_intent_envelope_check`, `cc_execution_begin`, `cc_plan_set_status`
- `wrapper/contracts/schemas/plan.yaml`, `wrapper/contracts/schemas/intent-contract.yaml`
- `wrapper/contracts/invariants.yaml`: INV-APPROVE-01, INV-INTENT-02, INV-PLAN-01, INV-EXEC-01

## Provenance

Re-grounded on the current wrapper for Context Circuit v1.0. Supersedes the earlier
`plan-approval` domain: v1.0 removed the plan-level approval gate and the
intermediate `approved` plan status, folding authorization into the intent gate
(Gate 1) plus the scope-envelope check. Raw `sources/` was not scanned.

## Acceptance notes

Accepted 2026-09-02. Replaces the retired `draft → approved` plan gate and the
`cc_plan_approve` engine verb (both removed in v1.0) with intent-derived
authorization. Extends the accepted "separate lifecycle gates" and "intent is the
single upstream gate" decisions in `context/DECISIONS.md`.
