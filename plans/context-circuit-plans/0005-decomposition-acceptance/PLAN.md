# 0005 — Prove one-plan and stacked-plan derivation

Plan ID: `0005-decomposition-acceptance`  
Intent: `i002-plan-decomposition-guidance`  
Status: `draft`  
Depends on: `0004-handoff-representation`

## Original request and coverage

Create the third i002 plan: add deterministic and conversation-level coverage for both valid decomposition outcomes and their authority boundaries.

## Objective and desired behavior

The semantic suite proves that bounded Standard work can remain one plan with multiple tasks, independently bounded partitions become multiple dependent plans, and incorrectly combining or splitting either case fails. It also proves that intent authorization, verification, delivery, candidate, and scope-safety ownership are unchanged.

## Constraints and non-goals

- Preserve existing acceptance-map coverage and credential-free/provider-neutral behavior.
- Use existing plan, run-stack, feasibility, scenario, and full-suite surfaces unless a dedicated check is necessary.
- Do not add a second approval gate or change tier semantics.
- Do not deliver or publish implementation changes.

## Product Knowledge grounding

- `tracing` (`context/domains/tracing/README.md`) — grounds trace-derived done checks and feasibility.
- `plan-authorization` (`context/domains/plan-authorization/README.md`) — grounds authority preservation.
- `planning` (`docs/planning.md`) — grounds plan status and derivation behavior.
- `assurance` (`context/domains/assurance/README.md`) — ensures plan count is not a tier proxy.
- `acceptance-surface` (`test/acceptance/criteria-map.yaml`) — owns suite mapping and semantic acceptance invariants.

## Repositories and source evidence

The sole repository is `context-circuit-source`. The trace maps existing evidence to `test/plans/test-plans.sh`, `test/run-stack/test-run-stack.sh`, `test/intent/test-feasibility.sh`, `test/scenarios/test-scenarios.sh`, `test/acceptance.sh`, and their scenario/map documentation. Scenario A represents one Standard plan and Scenario C represents independently executed plans combined into one delivered change set.

## Tasks

1. **DCA-001** (context-circuit-source; paths: `test/plans/`, `test/run-stack/`, `test/intent/`, `test/scenarios/`, `test/acceptance.sh`, `test/acceptance/criteria-map.yaml`, `docs/`; depends on: none) — add or align explicit one-plan/stacked-plan checks, wrong combine/split negatives, and authority-preservation documentation. Acceptance: `DCA-AC-001`, `DCA-AC-002`, `DCA-AC-003`. Verification: `DCA-VT-001`, `DCA-VT-002`, `DCA-VT-003`.

## Acceptance criteria

- `DCA-AC-001` — A bounded Standard case passes as one plan containing multiple tasks and fails if unnecessarily split by tier.
- `DCA-AC-002` — An independently executable/verifiable case passes as stacked dependent plans and fails if incorrectly combined without its dependency/lifecycle boundary.
- `DCA-AC-003` — The complete suite remains green and documents no second approval gate or ownership shift for verification, candidate, delivery, or scope safety.

## Verification

- `DCA-VT-001` — `sh test/plans/test-plans.sh && sh test/intent/test-feasibility.sh`; one-plan task representation and intent-derived authorization remain valid.
- `DCA-VT-002` — `sh test/run-stack/test-run-stack.sh && sh test/scenarios/test-scenarios.sh`; stacked dependencies, independent lifecycles, and both scenarios pass.
- `DCA-VT-003` — `sh test/acceptance.sh`; the full semantic suite and criteria-map cross-check pass.

## Assumptions, open questions, and risks

- Assumption: Scenario A and Scenario C can be strengthened or supplemented without changing their calibrated lifecycle meaning.
- Open question: none; the trace answered that change-set aggregation does not remove per-plan verification.
- Risk: the criteria map has older metadata in places; new coverage must not weaken credential-free and provider-neutral checks.

## Expected commits and delivery notes

One source-only Conventional Commit is expected for this repository, normally containing the complete i002 implementation. Delivery, merge, push, and publication remain separate explicit actions.

## Expected Product Knowledge impact

No durable Product Knowledge update is expected. Review tracing, planning, assurance, and acceptance-surface references at completion.

