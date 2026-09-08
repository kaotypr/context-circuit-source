# 0004 — Carry decomposition through the plan handoff

Plan ID: `0004-handoff-representation`  
Intent: `i002-plan-decomposition-guidance`  
Status: `draft`  
Depends on: `0003-decomposition-rule`

## Original request and coverage

Create the second i002 plan: represent the decomposition decision durably in the trace-to-plan handoff and derived plan set.

## Objective and desired behavior

The coordinator ratifies `task_partition` into either one plan with task dependencies or multiple plans with existing `plan_dependencies`, bounded repository/path mappings, and a rationale. Each plan remains independently executable and verifiable at Standard, and the existing schema is reused rather than expanded without need.

## Constraints and non-goals

- Preserve the existing trace and plan schema ownership model.
- Keep intra-plan `depends_on` distinct from inter-plan `plan_dependencies`.
- Preserve one worker and one independent verifier per Standard plan.
- Do not add plan approval, runtime redesign, or delivery semantics.

## Product Knowledge grounding

- `trace-manifest` (`wrapper/contracts/schemas/trace-manifest.yaml`) — owns `task_partition`, feasibility, questions, and done checks.
- `plan` (`wrapper/contracts/schemas/plan.yaml`) — represents plan dependencies, bounded paths, task dependencies, and execution policy.
- `tracing` (`context/domains/tracing/README.md`) — defines the manifest consumer and feasibility boundary.
- `planning` (`docs/planning.md`) — defines schema-3 derived plans and no plan-level approval.
- `plan-authorization` (`context/domains/plan-authorization/README.md`) — preserves intent-derived authorization.

## Repositories and source evidence

The sole repository is `context-circuit-source`. The trace identifies the existing schema fields and cc-plan handoff as sufficient for stacked plans; no schema expansion is required by the approved intent. Runtime evidence already exists in `wrapper/runtime/engine.sh` and `test/run-stack/test-run-stack.sh` for dependency readiness and lifecycle behavior.

## Tasks

1. **HPR-001** (context-circuit-source; paths: `wrapper/contracts/schemas/trace-manifest.yaml`, `wrapper/contracts/schemas/plan.yaml`, `.agents/skills/cc-plan/`, `docs/planning.md`; depends on: none) — carry ratified decomposition and rationale into plan records using existing fields. Acceptance: `HPR-AC-001`, `HPR-AC-002`. Verification: `HPR-VT-001`, `HPR-VT-002`.

## Acceptance criteria

- `HPR-AC-001` — One-plan derivation can retain multiple ordered tasks, while stacked derivation records explicit acyclic `plan_dependencies` and bounded repository/path mappings.
- `HPR-AC-002` — Each stacked plan retains a distinct worker/verifier lifecycle and the rationale explains why partitions are independently executable or verifiable.

## Verification

- `HPR-VT-001` — `sh test/plans/test-plans.sh`; the existing plan suite passes for task dependencies, bounded mappings, and intent-derived authorization.
- `HPR-VT-002` — `sh test/run-stack/test-run-stack.sh`; dependent-plan readiness, failure containment, and same/cross-repository stacking pass.

## Assumptions, open questions, and risks

- Assumption: existing `plan_dependencies` and its `reason` field are the canonical machine-readable representation; readable plan text carries the full rationale.
- Open question: none; the trace marks schema expansion as unnecessary and the approved intent permits this resolution.
- Risk: combining independently delivered plans later must not remove per-plan verification; change-set aggregation is separate.

## Expected commits and delivery notes

One source-only Conventional Commit is expected for this repository, normally combined with adjacent i002 implementation work. Delivery remains a separate explicit action.

## Expected Product Knowledge impact

No durable Product Knowledge update is expected. Existing schema and run-stack knowledge should be reviewed for alignment at completion.

