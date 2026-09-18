---
name: cc-intent
description: Write a Context Circuit intent record — grounding it in existing knowledge, its body fields, and the numbered open questions a person answers.
---

# Write an intent

An intent says what outcome is wanted, before any detailed code investigation.
It is written for a person to read and approve.

Write it in the language recorded for the active member in `members.yaml`, and
in English where none is. Approval is the product's central gate, and a person
can only approve an outcome they actually understand; an intent they skim in a
second language is a gate in name only. Record their approval in the words they
used rather than a translation of them. Identifiers and the project's own domain
vocabulary keep their form whatever language the prose is in, and the slug stays
a lowercase ASCII slug.

## Writing in a language that is not English

A recorded language settles which words a record uses and nothing about how they
are put together. Composing in English and translating produces prose that is
grammatical and wrong in register — stiff, formal, and formal in a way nobody
chose, because words that are neutral in English land as bureaucratic once
carried across literally.

**The record's structure stays English whatever the prose is.** The headings —
`Goal`, `Non-goals`, `Constraints`, `Success criteria`, `Repository scope`,
`Open questions` — are the record's shape, like `created_by` and the slug, and
they are named in the language table for that reason. Translating one makes two
members produce differently-shaped records, and a translation that varies per
session is not a field name.

Then compose in that language rather than into it:

- **Let the sentence structure follow the target language.** Do not carry
  English clause order or sentence length across. A sentence that hinges on
  `, so` in English is often two sentences elsewhere.
- **Prefer the verb where English would nominalize.** "No change to how the
  current quarter is determined" is ordinary English and heavy almost everywhere
  else.
- **Technical vocabulary keeps the form engineers in that language actually
  say**, which is usually the English word. Translating `refresh`, `cache`, or
  `deploy` into a native equivalent produces manual-register prose, not clearer
  prose. Be consistent: half-translated technical vocabulary reads worse than
  either choice made whole.
- **Names are quoted, never translated**, in either direction, as always.

Where the member records a `tone` in `members.yaml`, follow it; it is that
team's own answer and it overrides the guidance above. A workspace writing
Bahasa Indonesia might record `semi-formal; keep technical terms in English`.

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
