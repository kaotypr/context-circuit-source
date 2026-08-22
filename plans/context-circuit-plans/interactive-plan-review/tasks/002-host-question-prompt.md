---
schema_version: 2
id: IPR-002
plan: interactive-plan-review
status: draft
repository: context-circuit-source
paths:
  - docs/plan-review.md
  - docs/planning.md
  - docs/host-capabilities.md
  - .agents/skills/cc-review-plan/SKILL.md
  - agents/coordinator.md
  - wrapper/adapters/README.md
depends_on: [IPR-001]
acceptance: [IPR-AC-04, IPR-AC-05]
verification: [IPR-VT-02, IPR-VT-03, IPR-VT-05]
expected_evidence:
  - Review Card still lists risks, open questions, and 0-3 human decisions.
  - Docs name optional host question prompts and the conversational fallback.
  - Guidance forbids treating a selected option as gate confirmation.
  - No host_evidence field records question transcripts or tool payloads.
stop_conditions:
  - A question-prompt answer changes plan or task status.
  - A selected Approve label skips present-approval-card and becomes approve-plan.
  - Missing question prompt is classified host-blocked.
  - Question transcripts are stored in workspace or host_evidence records.
---

# Present review decisions through optional host question prompts

## Objective

After a named-plan review, offer the same 0–3 focused decisions through the
host’s optional native question prompt when that primitive is listed for the
turn, and keep the Review Card as the durable fallback.

## Work

Keep the Review Card contract. Make `Human decisions` the source of the
optional prompt: at most three questions, each with concrete options plus the
host’s built-in Other path when it exists. Present the card first, then the
prompt, so a host without the primitive still has a complete read-only
result.

Document the observed optional primitives without making them a second
router:

- Cursor Agent: `AskQuestion`
- Claude Code: `AskUserQuestion`
- Codex CLI: `request_user_input` when the host lists it

If the primitive is absent, denied, mode-gated, or fails, fall back to the
card text. That fallback is `filesystem-only` presentation, not
`host-blocked`. Do not add `question_prompt_capability` to `host_evidence` in
this plan.

If a decision is a next-action label, the only allowed continuation is the
existing first card for that route in the same session. `Approve this plan`
still shows the approval card and still requires the exact current
confirmation. Review answers must not revise the bundle; a later explicit
revise request remains a separate write.

## Non-goals

Do not require live host UI in offline tests. Do not persist answers as a
review artifact. Do not fold approval, execution, or plan mutation into the
prompt.

## Verification

Use IPR-VT-02, IPR-VT-03, and IPR-VT-05.

## Expected evidence

Updated Review Card and host-capability wording, adapter guidance that names
the optional primitives and the fallback, and no new transcript-bearing
fields.

## Stop conditions

Stop if answers authorize a gate, if a missing prompt blocks the route, or if
question payloads enter workspace records.
