---
schema_version: 2
id: IPR-001
plan: interactive-plan-review
status: draft
repository: context-circuit-source
paths:
  - wrapper/runtime/engine.sh
  - wrapper/contracts/routes.yaml
  - .agents/skills/cc-review-plan/SKILL.md
  - .agents/skills/cc-plan/SKILL.md
  - docs/plan-review.md
  - docs/planning.md
  - docs/getting-started.md
  - docs/host-capabilities.md
  - agents/coordinator.md
  - test/routing/fixtures.yaml
depends_on: []
acceptance: [IPR-AC-01, IPR-AC-02, IPR-AC-03]
verification: [IPR-VT-01, IPR-VT-05]
expected_evidence:
  - Named-plan review phrasing routes to review-plan and stays read-only.
  - Unnamed review clarifies the target instead of selecting an unrelated plan.
  - cc-review-plan exists as a thin discovery adapter; cc-plan no longer owns review procedure.
  - Human-facing docs name Review plan <id> without requiring a skill name.
stop_conditions:
  - Review mutates plan, task, lease, runtime, or Git state.
  - A skill or role file becomes a second route or lifecycle owner.
  - Unnamed review inspects a guessed or unrelated plan.
---

# Make named-plan review a specific request and skill

## Objective

Give humans one obvious named-plan review request and give hosts a dedicated
skill to discover that path, without changing the canonical `review-plan`
action or making review write anything.

## Work

Keep `review-plan` as the Stage B action. Extend probe matching so named
review phrasing such as `Review plan <id>`, `Walk me through plan <id>`, and
questions about that plan’s risks or open decisions select `plan-review`.
When the request is a review but has no usable plan id, emit `clarify-target`
the same way unnamed approval already does.

Add `.agents/skills/cc-review-plan/SKILL.md` as a thin adapter: named-plan
review, read-only Review Card, cite `INV-AUTH-02`, do not approve, claim a
lease, or start children. Narrow `cc-plan` to draft and prepare so the two
skills do not carry parallel review procedures.

Update coordinator and docs so the human surface is `Review plan <id>`. Do
not add slash commands or `.cursor` policy files. Add routing fixtures for
named review, unnamed review, and at least one nearby phrase that must not
steal the verification probe.

## Non-goals

Do not implement the host question-prompt presentation in this task. Do not
change plan status, approval, execution, or context-set budgets.

## Verification

Use IPR-VT-01 and IPR-VT-05.

## Expected evidence

Updated probe/action mapping, new skill file, narrowed `cc-plan`, docs that
name the human request, and routing fixtures for named versus unnamed review.

## Stop conditions

Stop if review becomes a write, if unnamed review guesses a plan, or if the
new skill defines route or gate policy.
