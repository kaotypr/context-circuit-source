# 0008 — Spawn the tracer immediately after intent approval

Plan ID: `0008-approval-then-trace`  
Intent: `i006-approval-then-trace`  
Status: `draft`

## Original request and coverage

Intent approval must not skip the tracer and the feasibility check. Today `contract.yaml` can move to approved before the tracer is spawned. This plan covers the required order and proof that it cannot be skipped.

## Objective and desired behavior

After approval, the tracer always runs next (one child per repository). The feasibility check always runs after the tracer and before any plan is written. `contract.yaml` still becomes approved immediately (Gate 1); that stamp cannot skip those steps.

## Constraints and non-goals

- Approval still sets `contract.yaml` to approved immediately; the tracer is the required next step.
- Do not add a second human approval of plans.
- Tracing remains read-only; feasibility remains coordinator judgment.
- Preserve unrelated and dirty work; do not commit, publish, or deliver as part of this plan.

## Product Knowledge grounding

- `intent` — Gate 1, draft → approved.
- `tracing` — tracer spawn and feasibility before planning.
- `plan-authorization` — automatic plan derivation, no second gate.

Grounding summary: derived from the approved i006 contract without a tracer, at explicit human request. Settled: keep immediate approved; fix the skipped tracer.

## Repositories and source evidence

`context-circuit-source`. Expected surface: `cc-intent`, `cc-trace`, `cc-plan`, `agents/coordinator.md`, intent and tracing domain pages, intent tests.

## Tasks

1. **ATT-001** — put approval, tracer, and feasibility in the required order. Acceptance: `ATT-AC-001`–`ATT-AC-004`. Verification: `ATT-VT-001`.
2. **ATT-002** — prove the sequence cannot be skipped. Acceptance: `ATT-AC-001`. Verification: `ATT-VT-002`.

## Acceptance criteria

- `ATT-AC-001` — After approval, the tracer starts next before any plan is written.
- `ATT-AC-002` — Feasibility check runs after the tracer and before plans.
- `ATT-AC-003` — `contract.yaml` status cannot skip or fake those steps.
- `ATT-AC-004` — No second human approval of plans.

## Verification

- `ATT-VT-001` — `rg -n 'intent-approve|cc-trace|feasibility|before any plan' .agents/skills/cc-intent/SKILL.md .agents/skills/cc-trace/SKILL.md agents/coordinator.md`
- `ATT-VT-002` — `sh test/acceptance.sh`

## Assumptions, open questions, and risks

- Assumption: tracing for this maintainer planning pass was skipped at human request.
- Open question: none remaining at the intent level. Gate 1 still flips `contract.yaml` immediately.
- Risk: coordinator skill text without a conversation-level check could still skip the tracer in a live session.

## Expected commits and delivery notes

One source-only Conventional Commit is expected. Delivery remains a separate explicit action.

## Expected Product Knowledge impact

In-place updates expected on `intent` and `tracing`. Review at completion.
