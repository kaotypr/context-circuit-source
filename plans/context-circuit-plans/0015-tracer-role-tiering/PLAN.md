# 0015 — Support tracer in role-tiering.local.yaml

Plan ID: `0015-tracer-role-tiering`  
Intent: `i013-tracer-role-tiering`  
Status: `draft`

## Original request and coverage

When you ask the coordinator to set up `role-tiering.local.yaml` including tracer, that tracer config is written into the file and used when the tracer is spawned. Today the coordinator creates the file but reports tracer tiering unsupported and omits the entry.

## Objective and desired behavior

A `tracer` entry you ask for is written into `role-tiering.local.yaml`. A spawned tracer uses that `(model, effort)`. The coordinator does not report that tracer tiering is unsupported. Worker and verifier stay configurable.

## Constraints and non-goals

- Do not change the tracer's read-only role, spawn rules, or feasibility check.
- Do not let model or effort select a consequence tier or grant a gate.
- Adapter defaults still apply when no tracer entry is present.
- Explore worker honor is `0014-pair-role-tiering`, not this plan.
- Preserve unrelated and dirty work; do not commit, publish, or deliver as part of this plan.

## Product Knowledge grounding

- `tracing` — tracer spawn.
- `host-adapters` — host evidence, never a gate.
- `role-tiering` (`docs/role-tiering.md`) — currently worker and verifier only.

Grounding summary: derived from the approved i013 contract without a tracer, at explicit human request.

## Repositories and source evidence

`context-circuit-source`. Expected surface: `docs/role-tiering.md`, adapters, `cc-trace`, `cc-intent`, coordinator, tests.

## Tasks

1. **TRT-001** — allow tracer in `role-tiering.local.yaml`. Acceptance: `TRT-AC-001`, `TRT-AC-003`, `TRT-AC-004`. Verification: `TRT-VT-001`.
2. **TRT-002** — honor it when spawning the tracer. Acceptance: `TRT-AC-002`. Verification: `TRT-VT-002`.
3. **TRT-003** — prove it is no longer refused. Acceptance: `TRT-AC-001`. Verification: `TRT-VT-003`.

## Acceptance criteria

- `TRT-AC-001` — Requested tracer config is written into `role-tiering.local.yaml`.
- `TRT-AC-002` — Spawned tracer uses that `(model, effort)`.
- `TRT-AC-003` — Coordinator does not report tracer tiering as unsupported.
- `TRT-AC-004` — Tracer meaning and independence are unchanged.

## Verification

- `TRT-VT-001` — `rg -n 'tracer' docs/role-tiering.md && ! rg -n 'Tracer tiering isn’t supported|Tracer tiering isn.t supported' docs/role-tiering.md .agents/skills/ docs/ wrapper/adapters/`
- `TRT-VT-002` — `rg -n 'role-tiering|Model & effort|tracer' .agents/skills/cc-trace/SKILL.md .agents/skills/cc-intent/SKILL.md agents/coordinator.md docs/role-tiering.md`
- `TRT-VT-003` — `sh test/acceptance.sh`

## Assumptions, open questions, and risks

- Assumption: tracing for this maintainer planning pass was skipped at human request.
- Open question: none remaining at the intent level.
- Risk: the current "only worker and verifier" sentence in `docs/role-tiering.md` is the contract the coordinator quoted; it must change in the owning doc, not only in a skill.

## Expected commits and delivery notes

One source-only Conventional Commit is expected. Delivery remains a separate explicit action.

## Expected Product Knowledge impact

In-place update expected on `tracing` and `docs/role-tiering.md`. Review at completion.
