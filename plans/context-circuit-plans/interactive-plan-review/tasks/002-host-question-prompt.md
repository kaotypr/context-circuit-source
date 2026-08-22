---
schema_version: 2
id: IPR-002
plan: interactive-plan-review
status: done
repository: context-circuit-source
paths:
  - docs/plan-review.md
  - docs/host-capabilities.md
  - docs/planning.md
  - .agents/skills/cc-plan/SKILL.md
  - agents/coordinator.md
  - wrapper/contracts/schemas/session.yaml
depends_on: [IPR-001]
acceptance: [IPR-AC-04, IPR-AC-05]
verification: [IPR-VT-02, IPR-VT-03]
expected_evidence:
  - Review Card still reports risks, open questions, and at most three focused human decisions.
  - Procedure text tells the agent to offer those decisions through the host question prompt when that primitive exists.
  - Missing, denied, or failed prompts fall back to the card and are not host-blocked.
  - Question answers never write plan.yaml, task status, or host_evidence transcripts.
stop_conditions:
  - A question-prompt answer is treated as approval, execution, or any other gate.
  - A missing prompt is classified as host-blocked or retried as a required child.
  - Question transcripts, answers, or tool payloads are recorded in host_evidence.
---

# Present review decisions through optional host question prompts

## Objective

After the Review Card, offer the same 0–3 decisions through the host native
question UI when that primitive exists, without making the prompt a gate or a
required child.

## Work

Keep the Review Card contract in `docs/plan-review.md`. Add procedure text so
that, when the card has one to three focused human decisions, the current
session may present them through the host primitive:

- Cursor Agent: `AskQuestion`
- Claude Code: `AskUserQuestion`
- Codex CLI: `request_user_input` when listed

A missing, denied, or failed prompt falls back to the card text. That is not
`host-blocked` and must not be retried as if it were a required writer or
verifier child. Put this instruction in `cc-plan` and the review docs, not in
a new skill.

A chosen next-action label such as `Revise this plan` or `Approve this plan`
may continue only into that route’s existing first card in the same session.
It cannot skip confirmation or write status.

Do not add `question_prompt_capability` to `host_evidence`.

## Non-goals

Do not require a live host question UI in CI. Do not store answers as durable
review artifacts.

## Verification

Use IPR-VT-02 and IPR-VT-03.

## Expected evidence

Updated Review Card / skill / coordinator wording, fallback language, and a
diff that does not add host_evidence fields or a new skill.

## Stop conditions

Stop if the prompt becomes authorization, a required child, or a recorded
transcript.
