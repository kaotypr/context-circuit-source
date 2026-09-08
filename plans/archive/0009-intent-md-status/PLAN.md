# 0009 — Keep INTENT.md status truthful after a feasible tracer

Plan ID: `0009-intent-md-status`  
Intent: `i007-intent-page-status`  
Status: `draft`

## Original request and coverage

`INTENT.md` must not show a status that is wrong after the tracer and feasibility check have finished. Settled: keep the status line and update it; do not remove it.

## Objective and desired behavior

After a feasible tracer with no follow-up questions, the `INTENT.md` status line is updated to match that fact. `contract.yaml` remains the machine authority for approval.

## Constraints and non-goals

- Keep the `INTENT.md` status line; do not remove it.
- Do not replace `contract.yaml` as the approval authority.
- Do not add a second human approval of plans.
- Preserve unrelated and dirty work; do not commit, publish, or deliver as part of this plan.

## Product Knowledge grounding

- `intent` — owns `INTENT.md` and its status line.
- `tracing` — owns the feasible finding that status must reflect.

Grounding summary: derived from the approved i007 contract without a tracer, at explicit human request.

## Repositories and source evidence

`context-circuit-source`. Expected surface: `cc-intent`, `cc-trace`, coordinator, runtime status sync, intent template, intent tests.

## Tasks

1. **IMS-001** — keep `INTENT.md` status in sync after a feasible tracer. Acceptance: `IMS-AC-001`–`IMS-AC-003`. Verification: `IMS-VT-001`.
2. **IMS-002** — prove `INTENT.md` cannot lie. Acceptance: `IMS-AC-001`. Verification: `IMS-VT-002`.

## Acceptance criteria

- `IMS-AC-001` — After a feasible tracer with no follow-up questions, `INTENT.md` does not keep a stale status.
- `IMS-AC-002` — The status line is kept and updated, not removed.
- `IMS-AC-003` — `contract.yaml` remains the approval authority.

## Verification

- `IMS-VT-001` — `rg -n 'INTENT.md|status|feasible|human_status|intent-approve' .agents/skills/cc-intent/SKILL.md .agents/skills/cc-trace/SKILL.md wrapper/runtime/engine.sh docs/templates/intent.md`
- `IMS-VT-002` — `sh test/acceptance.sh`

## Assumptions, open questions, and risks

- Assumption: tracing for this maintainer planning pass was skipped at human request.
- Open question: none remaining at the intent level.
- Risk: runtime currently syncs status only on approve (`draft` → `approved`); a post-feasible line may need a new owned sync, not a skill-only note.

## Expected commits and delivery notes

One source-only Conventional Commit is expected. Delivery remains a separate explicit action.

## Expected Product Knowledge impact

In-place update expected on `intent`. Review at completion.
