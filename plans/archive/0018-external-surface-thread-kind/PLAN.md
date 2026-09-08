# 0018 — External surface: the thread publication kind

- **Plan ID:** `0018-external-surface-thread-kind`
- **Status:** done
- **Repository:** `context-circuit-source`
- **Depends on:** `0017-external-surface-plan-kind`
- **Owns:** the thread publication kind (reuses INV-EXTERNAL-01/02/03)

## Original request

Retroactive plan for the external-surface thread kind, built as if from an empty
repo. Source design:
`sources/system-design/context-circuit/v0.6/external-surface/thread.md`.

## Objective and desired behavior

- A second publication kind on the same backbone: publish a plan's open
  questions (INV-PLAN-04) to a chat discussion — a parent message
  (`[thread] [<plan-number>] <title> — open questions`) plus one fully-described
  threaded reply per question.
- Discussion-safe idempotency: it edits only its own messages, appends for new
  questions, and never touches or deletes a human reply or any other message.
- Slack-first, `ts`-based record; other chat providers share the shape.

## Constraints and non-goals

- Non-goal: any change to the plan kind or the core workflow.

## Product Knowledge grounding

- `invariants` (`wrapper/contracts/invariants.yaml`), `design-deltas`
  (`context/DESIGN-DELTAS.md`).

## Tasks

1. **THR-001** — thread record schema + the thread kind in `cc-publish`.

## Acceptance & verification

- One parent + one reply per open question; re-publish edits only its own
  messages, never a human reply.
- `sh test/external-surface/test-external-surface.sh`.

## Assumptions, open questions, risks

- Risk: clobbering a human's discussion reply — prevented by the
  own-messages-only idempotency rule.

## Expected commits and delivery notes

Additive to 0017; no core contract change.
