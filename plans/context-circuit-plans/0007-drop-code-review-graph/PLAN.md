# 0007 — Drop the leftover .code-review-graph folder

Plan ID: `0007-drop-code-review-graph`  
Intent: `i005-drop-code-review-graph`  
Status: `draft`

## Original request and coverage

A workspace started from the Context Circuit template should not keep a leftover `.code-review-graph` folder. This plan covers removing it from what ships and proving lifecycle behavior is unchanged.

## Objective and desired behavior

A new workspace from the template does not include `.code-review-graph`. A workspace that already has it no longer keeps it. Intent approval, tracing, planning, execution, and delivery are unchanged.

## Constraints and non-goals

- Do not change how an intent is approved, traced, planned, executed, or delivered.
- Do not redesign code review, graph tooling, or Product Knowledge retrieval.
- Preserve unrelated and dirty work; do not commit, publish, or deliver as part of this plan.

## Product Knowledge grounding

- `source-release-and-upgrade` (`context/domains/source-release-and-upgrade/README.md`) — owns what the template ships.
- `intent` (`context/domains/intent/README.md`) — lifecycle must stay unchanged.

Grounding summary: derived from the approved i005 contract without a tracer, at explicit human request. Feasible as a template/ignore cleanup.

## Repositories and source evidence

`context-circuit-source` is the sole repository. Expected surface: `template/`, assembly/ignore under `wrapper/` and `docs/`, and tests.

## Tasks

1. **CRG-001** (context-circuit-source; paths: `template/`, `wrapper/`, `docs/`; depends on: none) — stop keeping `.code-review-graph`. Acceptance: `CRG-AC-001`, `CRG-AC-002`. Verification: `CRG-VT-001`.
2. **CRG-002** (context-circuit-source; paths: `test/`; depends on: CRG-001) — prove a fresh workspace is clean and lifecycle still works. Acceptance: `CRG-AC-003`. Verification: `CRG-VT-002`.

## Acceptance criteria

- `CRG-AC-001` — A newly instantiated workspace from the template does not include `.code-review-graph`.
- `CRG-AC-002` — A workspace that already has it no longer keeps it after this change lands.
- `CRG-AC-003` — Intent approval, tracing, planning, execution, verification, and delivery continue to work as they do today.

## Verification

- `CRG-VT-001` — `! find template wrapper -name '.code-review-graph' -o -path '*/.code-review-graph/*' | grep -q .`
- `CRG-VT-002` — `sh test/acceptance.sh`

## Assumptions, open questions, and risks

- Assumption: tracing was skipped at human request; paths may be refined during implementation.
- Open question: none remaining at the intent level.
- Risk: a generated leftover could return unless ignore or assembly also covers it.

## Expected commits and delivery notes

One source-only Conventional Commit is expected. Delivery, merge, push, and publication remain separate explicit actions.

## Expected Product Knowledge impact

No new context unit is required. Review source-release-and-upgrade at completion if the ignore boundary changes.
