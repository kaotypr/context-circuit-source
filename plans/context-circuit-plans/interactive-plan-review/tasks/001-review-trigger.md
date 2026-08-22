---
schema_version: 2
id: IPR-001
plan: interactive-plan-review
status: ready
repository: context-circuit-source
paths:
  - wrapper/runtime/engine.sh
  - wrapper/contracts/routes.yaml
  - test/routing/fixtures.yaml
  - .agents/skills/cc-plan/SKILL.md
  - .agents/skills/cc-gates/SKILL.md
  - docs/plan-review.md
  - docs/planning.md
  - docs/getting-started.md
  - docs/host-capabilities.md
  - agents/coordinator.md
depends_on: []
acceptance: [IPR-AC-01, IPR-AC-02, IPR-AC-03, IPR-AC-06]
verification: [IPR-VT-01, IPR-VT-03, IPR-VT-04]
expected_evidence:
  - Named review phrasing routes to review-plan with authorization read-only.
  - A review request without a usable plan id routes to clarify-target and does not inspect an unrelated bundle.
  - cc-plan names read-only review as a discovery mode and does not add cc-review-plan.
  - Approval confirmation remains a gate owned by cc-gates and docs/gates.md.
stop_conditions:
  - The change adds .agents/skills/cc-review-plan or expands the shipped skill allowlist.
  - Named review becomes a mutating route or inspects a guessed plan.
  - cc-plan grows a second approval or execution procedure.
---

# Make named-plan review a specific request through cc-plan

## Objective

Make named-plan review a specific request that the engine can select, while
keeping `cc-plan` as the only plan discovery adapter.

## Work

Tighten Stage A so `Review plan <id>`, `Walk me through plan <id>`, and
equivalent named-plan phrasing select `plan-review` / `review-plan`. A review
request with no usable plan id must become `clarify-target` instead of
inspecting an arbitrary bundle, matching unnamed approval.

Update `cc-plan` so hosts can attach it to a review request: the description
and body must name read-only plan review, load the `plan-review` context set
when that route is selected, and forbid drafting, approving, claiming a
lease, or executing from a review. Leave approval confirmation to `cc-gates`
and `docs/gates.md`. Do not create `cc-review-plan`.

Cite the existing owners from coordinator and docs rather than copying
lifecycle policy into the skill.

## Non-goals

Do not merge `review-plan` into `draft-plan`. Do not add a skill, slash
command, or host-local policy file for review.

## Verification

Use IPR-VT-01, IPR-VT-03, and IPR-VT-04.

## Expected evidence

Routing fixtures for named and unnamed review, a `cc-plan` skill that names
read-only review, absence of `cc-review-plan`, and the unchanged seven-name
allowlist.

## Stop conditions

Stop if review mutates state, guesses a target, or introduces an eighth
shipped skill.
