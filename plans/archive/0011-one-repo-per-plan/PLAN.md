# 0011 — One plan, one repository, one worker

Plan ID: `0011-one-repo-per-plan`  
Intent: `i009-one-repo-per-plan`  
Status: `draft`

## Original request and coverage

One plan is work in one repository, done by one worker who only works in that repository. An intent that covers two repositories must become two plans. This plan covers the rule, the split, and the proof.

## Objective and desired behavior

Every plan names exactly one repository. The worker for that plan changes only that repository. An intent whose scope covers two repositories produces at least two plans. You still approve the intent once. Combining finished plans at delivery stays a later step.

## Constraints and non-goals

- Do not change Explore's existing one-repository pairing rule.
- Do not add a second human approval of plans.
- Do not make plan count a proxy for assurance tier.
- Delivery may still present more than one finished plan as one pull request.
- Preserve unrelated and dirty work; do not commit, publish, or deliver as part of this plan.

## Product Knowledge grounding

- `tracing` — task partition and stacked plans.
- `plan-execution` — current one-worker-across-repos model to tighten.
- `plan-authorization` — one intent may authorize several plans.
- `architecture` — current execution model statement.

Grounding summary: derived from the approved i009 contract without a tracer, at explicit human request. This is a deliberate change to INV-EXEC-02 / architecture.

## Repositories and source evidence

`context-circuit-source`. Expected surface: invariants, plan-execution and tracing domains, `cc-trace`, `cc-plan`, coordinator, tracer role, tests.

## Tasks

1. **ORP-001** — make the one-repository rule the rule. Acceptance: `ORP-AC-001`, `ORP-AC-002`. Verification: `ORP-VT-001`.
2. **ORP-002** — split a two-repository intent into two plans. Acceptance: `ORP-AC-003`–`ORP-AC-005`. Verification: `ORP-VT-002`.
3. **ORP-003** — prove the split. Acceptance: `ORP-AC-001`. Verification: `ORP-VT-003`.

## Acceptance criteria

- `ORP-AC-001` — Every derived plan names exactly one repository.
- `ORP-AC-002` — The worker for a plan changes only that repository.
- `ORP-AC-003` — A two-repository intent produces at least two plans.
- `ORP-AC-004` — No second human approval.
- `ORP-AC-005` — Delivery may still combine finished plans.

## Verification

- `ORP-VT-001` — `rg -n 'one repository|one worker|mapped repositories|INV-EXEC-02' wrapper/contracts/invariants.yaml context/domains/plan-execution/README.md context/ARCHITECTURE.md`
- `ORP-VT-002` — `rg -n 'one repository|stacked|plan_dependencies|task_partition' .agents/skills/cc-trace/SKILL.md .agents/skills/cc-plan/SKILL.md context/domains/tracing/README.md`
- `ORP-VT-003` — `sh test/acceptance.sh`

## Assumptions, open questions, and risks

- Assumption: tracing for this maintainer planning pass was skipped at human request.
- Open question: none remaining at the intent level.
- Risk: existing stacked-plan guidance (i002) must stay; this tightens the repo bound rather than undoing stacks.

## Expected commits and delivery notes

One source-only Conventional Commit is expected. Delivery remains a separate explicit action.

## Expected Product Knowledge impact

In-place updates expected on `tracing`, `plan-execution`, and `plan-authorization`. Review at completion.
