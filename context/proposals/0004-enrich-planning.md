---
id: 0004-enrich-planning
target_context_unit: context/domains/plan-review/README.md
operation: change
statement: >
  The planning-and-review page describes authoring generically. Add the shipped
  plan-creation brief, the readable PLAN.md contents, the plan.yaml field list
  (including repair limit, stop conditions, and expected Product Knowledge
  impact), and the request-fidelity "remove repetition, not meaning" rule.
evidence_refs:
  - .agents/skills/cc-plan/SKILL.md
  - wrapper/contracts/schemas/plan.yaml
  - wrapper/contracts/schemas/task.yaml
  - docs/planning.md
affected_repositories:
  - context-circuit-source
affected_commits:
  - 4b8ac0b
confidence: high
status: review-needed
---

# Enrich planning with the shipped plan brief and plan.yaml surface

Add:

- The plan-creation brief (request, objective, repositories, product_knowledge,
  context_grounding, sources, known_decisions, open_questions) that grounds scope
  before a plan is written.
- The readable PLAN.md contents: status/execution summary, repository ownership,
  scope + non-goals, acceptance, verification, risks, expected commits/delivery,
  human decisions, expected Product Knowledge impact.
- The plan.yaml field surface owned by its schema: allowed paths, declared
  effects, acceptance/verification ids, repair limit, plan-level stop conditions,
  and context units to reassess at completion.
- Request fidelity: a plan removes repetition, not meaning; detail is not
  discarded merely to shorten.

## Acceptance action

Fold into Behavior/Data; ground in `wrapper/contracts/schemas/plan.yaml` and
`docs/planning.md`. Note: `INV-PLAN-03` (never-reused/highest-ever) is a shipped
invariant not stated in the design — recorded as a delta (see `0013`), page keeps
the shipped rule.
