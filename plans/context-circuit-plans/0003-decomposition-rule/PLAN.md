# 0003 — Define the trace-to-plan decomposition rule

Plan ID: `0003-decomposition-rule`  
Intent: `i002-plan-decomposition-guidance`  
Status: `draft`

## Original request and coverage

Create the first i002 plan: make the boundary between a task and a plan explicit. This covers the rule that chooses one plan with embedded tasks versus stacked plans, while retaining the approved intent's tier and authority boundaries.

## Objective and desired behavior

Tracing and planning guidance states that one bounded worker/verifier lifecycle and one bounded change surface produce one plan with tasks. Independently executable or independently verifiable partitions, meaningful dependency edges, or distinct failure surfaces justify multiple stacked plans. The rationale is explainable in human terms and grounded in the trace.

## Constraints and non-goals

- Preserve the single intent approval gate and intent-derived plan authorization.
- Keep plan count independent from Explore/Standard/Critical tiering.
- Do not change runtime execution, delivery, candidate, verification, or scope-safety ownership.
- Do not redesign the separate i001 conversation-library work.
- Do not modify `sources/`, commit, publish, or deliver as part of this plan.

## Product Knowledge grounding

- `tracing` (`context/domains/tracing/README.md`) — owns the post-approval trace manifest, feasibility judgment, task partition, and plan derivation boundary.
- `plan-authorization` (`context/domains/plan-authorization/README.md`) — establishes intent-derived authorization and the absence of plan approval.
- `planning` (`docs/planning.md`) — defines plan derivation, grounding, identifiers, and lifecycle status.
- `intent` (`context/domains/intent/README.md`) — grounds Gate 1 and the approved intent as the upstream decision.

Grounding summary: the current trace has no out-of-scope reach and no unresolved intent-level question. Its plan-resolution question is carried forward: use the readable plan and coordinator handoff as the durable rationale location, without requiring a new trace schema field.

## Repositories and source evidence

`context-circuit-source` is the sole repository. The trace maps the change to `.agents/skills/cc-trace/`, `.agents/skills/cc-plan/`, `agents/tracer.md`, `context/domains/tracing/`, `context/domains/plan-authorization/`, and `docs/planning.md`; the relevant sites are the existing task-partition, feasibility, handoff, and no-plan-approval guidance.

## Tasks

1. **DCR-001** (context-circuit-source; paths: `.agents/skills/cc-trace/`, `.agents/skills/cc-plan/`, `agents/tracer.md`, `context/domains/tracing/`, `context/domains/plan-authorization/`, `docs/planning.md`; depends on: none) — define and align the one-plan/stacked-plan rule, including durable rationale and human explanation. Acceptance: `DCR-AC-001`, `DCR-AC-002`. Verification: `DCR-VT-001`, `DCR-VT-002`.

## Acceptance criteria

- `DCR-AC-001` — Guidance explicitly names bounded execution and verification boundaries, meaningful dependencies, and distinct failure surfaces as the signals for one plan versus stacked plans.
- `DCR-AC-002` — Guidance preserves one Standard plan with embedded tasks and states that plan count is not a proxy for assurance tier or a new approval gate.

## Verification

- `DCR-VT-001` — `rg -n -i 'one plan|multiple plans|stacked plans|bounded (execution|change|repository/path)|independent (execution|verification)|failure surface|dependency' .agents/skills/cc-trace/SKILL.md .agents/skills/cc-plan/SKILL.md agents/tracer.md context/domains/tracing/README.md docs/planning.md`; evidence contains both outcomes and all decision signals.
- `DCR-VT-002` — `rg -n -i 'one worker|embedded tasks|no separate plan.*approval|assurance tier|human.*explain|rationale' .agents/skills/cc-plan/SKILL.md docs/planning.md context/domains/plan-authorization/README.md`; evidence preserves lifecycle and authority boundaries.

## Assumptions, open questions, and risks

- Assumption: rationale is required in the derived readable plan/coordinator explanation and need not become a new trace-manifest field.
- Open question: none; the trace's representation question is resolved at plan level as permitted by the approved intent.
- Risk: guidance could conflate intra-plan task dependencies with inter-plan `plan_dependencies`; later plans must keep those concepts distinct.

## Expected commits and delivery notes

One source-only Conventional Commit is expected for this repository, covering i002 implementation as appropriate. Delivery, merge, push, and publication remain separate explicit actions.

## Expected Product Knowledge impact

No durable Product Knowledge update is expected from this plan. Existing tracing, planning, and authorization pages should be reviewed at completion.

