---
kind: domain
status: proposed
title: Named-plan review
slug: plan-review
owners: []
sources:
  - plans/context-circuit-plans/interactive-plan-review/plan.yaml
  - plans/context-circuit-plans/interactive-plan-review/PLAN.md
source_revisions:
  - plan: interactive-plan-review
    status: done
    updated_at: 2026-08-22T18:12:12Z
generated_at: 2026-08-23T00:00:00Z
review_date: 2026-09-22
freshness: proposed-from-done-plans
assumptions:
  - Canonical lifecycle status is plan.yaml, not PLAN.md.
unknowns: []
contradictions:
  - interactive-plan-review PLAN.md still says Status: draft while plan.yaml is done.
acceptance:
  state: pending
  accepted_at:
  accepted_by:
workflows:
  - docs/plan-review.md
---

# Named-plan review

## Summary

A human names one plan and gets a read-only Review Card. Optional host
question prompts may repeat at most three focused choices. Route
`Review plan <id>` and equivalent named-plan review requests here.

## Scope

Inside: named-plan review routing, the `plan-review` context set, Review Card
contents, missing-target clarification, and optional host question prompts.

Outside: mutating plan or task status, treating a prompt answer as a gate,
adding `cc-review-plan`, inspecting an unrelated plan, and expanding into
`sources/` or other plan bundles.

## Behavior

`Review plan <id>` and equivalent named phrasing route to `review-plan` with
`authorization: read-only`. Hosts discover this through `cc-plan`; there is
no eighth shipped skill.

A review request with no usable plan id routes to `clarify-target` and does
not inspect an unrelated bundle.

The Review Card reports outcome, summary, approval scope, non-effects, tasks,
acceptance mapping, risks, contradictions, and at most three focused human
decisions. When the current host exposes an optional native question prompt,
those decisions are also offered there:

| Host | Optional primitive | Fallback |
| --- | --- | --- |
| Cursor Agent | `AskQuestion` | Review Card text |
| Claude Code | `AskUserQuestion` | Review Card text |
| Codex CLI | `request_user_input` when listed | Review Card text |

A missing, denied, or failed prompt is not `host-blocked` and must not retry
as a required child. Question transcripts are not recorded in `host_evidence`.

Choosing a next-action label such as `Approve this plan` may continue only
into that route's existing first card in the same session. It cannot skip
confirmation or change status.

## Workflows

- Review Card owner: `docs/plan-review.md`
- Approval confirmation, if chosen later: `docs/gates.md` via `cc-gates`

## Constraints and edge cases

Review never writes `plan.yaml`, task status, leases, runtime, or Git.
`cc-plan` does not keep a second approval procedure.

## Implementation references

- `.agents/skills/cc-plan/SKILL.md`
- `wrapper/runtime/engine.sh` `plan-review` / `review-plan`
- Shipped skill allowlist remains the existing seven adapters

## Verification

Done-plan verification IDs IPR-VT-01–IPR-VT-05.

## Provenance

Read only the selected done plan `interactive-plan-review` (CC-005). Raw
`sources/` was not scanned. HEAD at generation was `b7a11f3`.

## Acceptance notes

This page is proposed. Human context acceptance is still required.
