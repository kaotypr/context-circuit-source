# 0010 — Write plans as soon as the tracer is feasible

Plan ID: `0010-plans-after-trace`  
Intent: `i008-plans-after-trace`  
Status: `draft`

## Original request and coverage

Once the intent is approved and the tracer plus feasibility check find it feasible, the coordinator writes the plans in that same turn. Execution stays a separate ask. This plan depends on `0008-approval-then-trace`.

## Objective and desired behavior

After approval and a feasible tracer with no follow-up questions, the coordinator writes the plans without you asking again. If the look is not feasible or an intent-level question remains, no plans are written. Writing plans does not start execution.

## Constraints and non-goals

- Depends on `0008-approval-then-trace` so the tracer actually runs first.
- Do not start execution, verification, completion, or delivery from this handoff.
- Preserve unrelated and dirty work; do not commit, publish, or deliver as part of this plan.

## Product Knowledge grounding

- `tracing` — feasible finding triggers plan derivation.
- `plan-authorization` — automatic derivation, no second gate.
- `planning` — writing `PLAN.md` and `plan.yaml` from the trace manifest.

Grounding summary: derived from the approved i008 contract without a tracer, at explicit human request.

## Repositories and source evidence

`context-circuit-source`. Expected surface: `cc-intent`, `cc-trace`, `cc-plan`, coordinator, tracing and plan-authorization domains, tests.

## Tasks

1. **PAT-001** — write plans as soon as the tracer is feasible. Acceptance: `PAT-AC-001`–`PAT-AC-003`. Verification: `PAT-VT-001`.
2. **PAT-002** — prove the handoff. Acceptance: `PAT-AC-001`. Verification: `PAT-VT-002`.

## Acceptance criteria

- `PAT-AC-001` — Same-turn plan derivation after a feasible tracer.
- `PAT-AC-002` — No plans when not feasible or an intent-level question remains.
- `PAT-AC-003` — Writing plans does not start execution.

## Verification

- `PAT-VT-001` — `rg -n 'feasible|derive|cc-plan|same turn|before any plan' .agents/skills/cc-intent/SKILL.md .agents/skills/cc-plan/SKILL.md .agents/skills/cc-trace/SKILL.md agents/coordinator.md`
- `PAT-VT-002` — `sh test/acceptance.sh`

## Assumptions, open questions, and risks

- Assumption: tracing for this maintainer planning pass was skipped at human request.
- Open question: none remaining at the intent level.
- Risk: a host-blocked tracer must still not invent plans.

## Expected commits and delivery notes

One source-only Conventional Commit is expected. Delivery remains a separate explicit action.

## Expected Product Knowledge impact

In-place updates expected on `tracing` and `plan-authorization`. Review at completion.
