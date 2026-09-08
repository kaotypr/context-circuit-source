# 0014 — Honor role-tiering.local.yaml in Explore

Plan ID: `0014-pair-role-tiering`  
Intent: `i012-pair-role-tiering`  
Status: `draft`

## Original request and coverage

`role-tiering.local.yaml` must apply in Explore the same way it applies in `cc-execute`. The `cc-pair` skill should include the same Model & effort per role step.

## Objective and desired behavior

The Explore worker runs at the model and effort set for that role. `cc-pair` includes the same Model & effort step as `cc-execute`. This changes cost and speed only.

## Constraints and non-goals

- Do not change what Explore means, the one-repository rule, or the absence of a verifier.
- Do not let model or effort select a consequence tier or grant a gate.
- Adapter defaults still apply when no local file is present.
- Preserve unrelated and dirty work; do not commit, publish, or deliver as part of this plan.

## Product Knowledge grounding

- `direct-collaboration` — Explore worker spawn.
- `host-adapters` — host evidence, never a gate.
- `role-tiering` (`docs/role-tiering.md`) — file shape and defaults.

Grounding summary: derived from the approved i012 contract without a tracer, at explicit human request.

## Repositories and source evidence

`context-circuit-source`. Expected surface: `cc-pair`, `cc-execute`, `docs/role-tiering.md`, pairing tests.

## Tasks

1. **PRT-001** — honor `role-tiering.local.yaml` in Explore. Acceptance: `PRT-AC-001`, `PRT-AC-003`. Verification: `PRT-VT-001`.
2. **PRT-002** — add the Model & effort step to `cc-pair`. Acceptance: `PRT-AC-002`. Verification: `PRT-VT-002`.
3. **PRT-003** — prove the Explore worker follows it. Acceptance: `PRT-AC-001`. Verification: `PRT-VT-003`.

## Acceptance criteria

- `PRT-AC-001` — Explore worker uses the configured model and effort.
- `PRT-AC-002` — `cc-pair` has the same Model & effort step as `cc-execute`.
- `PRT-AC-003` — Explore meaning and verifier floor are unchanged.

## Verification

- `PRT-VT-001` — `rg -n 'role-tiering|Model & effort|worker' .agents/skills/cc-pair/SKILL.md docs/role-tiering.md`
- `PRT-VT-002` — `rg -n 'Model & effort per role' .agents/skills/cc-pair/SKILL.md .agents/skills/cc-execute/SKILL.md`
- `PRT-VT-003` — `sh test/acceptance.sh`

## Assumptions, open questions, and risks

- Assumption: tracing for this maintainer planning pass was skipped at human request.
- Open question: none remaining at the intent level.
- Risk: `docs/role-tiering.md` currently allows only worker and verifier; this plan must still honor those roles in Explore even before tracer support (i013).

## Expected commits and delivery notes

One source-only Conventional Commit is expected. Delivery remains a separate explicit action.

## Expected Product Knowledge impact

In-place update expected on `direct-collaboration`. Review at completion.
