---
name: cc-intent
description: Write a Context Circuit intent record — grounding it in existing knowledge, its body fields, and the numbered open questions a person answers.
---

# Write an intent

An intent says what outcome is wanted, before any detailed code investigation.
It is written for a person to read and approve.

## Ground it first

Retrieve the bearing knowledge through `context/INDEX.md` and read the request's
apparent meaning against what the project already records. What that reading
contradicts, or cannot settle, is a question for the person — not a gap to fill
with a quiet assumption.

Where the project records the actors it serves, read the stories of the ones
this request touches. A request cutting across a standing need is a question for
the person, and naming the actor an outcome serves tells the next reader who it
is for.

Create the record through the executable, which reserves the ID:

```sh
context-circuit-cli --workspace <root> record create --kind intent --slug SLUG --title TITLE
```

Then write its body from that grounding: goal, non-goals, constraints,
observable success criteria, and rough repository scope.

## Open questions

Each decision a person must settle goes under `## Open questions`: a numbered
list, the question in bold, an italic `_Answer:_` line beneath once settled.

The number is the handle a person answers by, so it holds still. An answered
question keeps its number and is never deleted, and a later one takes the next
unused number. With no questions the unnumbered empty-state line stays as
written; a placeholder would invent a pending decision.

A question the request already answers is applied without asking. An
implementation question belongs in the plan, not here.

## What belongs here and what does not

A question that bears on the intended outcome returns to this record even after
planning has started; a plan never silently settles one. A newly needed file or
repository *within* the approved outcome is explained and recorded rather than
re-gated.

Instants — `created_at` and the `approved_at` a gate stamps — are canonical ISO
8601 UTC timestamps, `2026-09-15T10:53:00Z`. Every other date, written by hand or
through the executable, is an ISO 8601 calendar date, `YYYY-MM-DD`. No other date
format belongs in a workspace file.

Intent and plan IDs are workspace-global. Use only `created_by` for
member-related metadata; there are no assignee, owner, reviewer, or member
namespaces. Never reuse a reserved ID, including after archival or deletion.

`.context-circuit/docs/working.md` describes the record structures in full.
