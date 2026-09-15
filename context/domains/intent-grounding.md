# Grounding an intent and its open questions

An intent is written before detailed code investigation, which leaves accepted
project knowledge as the only thing it can be grounded in. Retrieval comes first,
through [the retrieval catalog](knowledge-notes.md), and what the request was
understood to mean is then read against what the project already records.

This is the inbound half of the knowledge circuit that
[the product](../product/what-v2-is.md) exists for. The outbound half — knowledge
reconciled at completion — was mechanized from the start; grounding was left to a
clause until it earned a shape of its own, because an intent written past a
recorded decision costs a spent approval and a reopened gate.

## What the reading produces

A reading the knowledge contradicts, and one it cannot settle, are the same thing
from the intent's side: a decision the agent is not positioned to make. Filling
either with a quiet assumption is what the gate exists to prevent, since the
person then approves an outcome whose real question never reached them.

So they become open questions rather than prose the reader has to notice. An
implementation question is not one of these — it belongs to a plan, where code is
actually being read — and a question the request already answers is applied
without asking again.

## Numbered, stable, and answered in place

Questions are a numbered list: the question in bold, an italic `_Answer:_` line
beneath it once settled. The number is a handle, so a person answers "1: yes,
everyone goes through the queue" instead of restating the question. That only
works if numbers hold still: an answered question keeps its number and is never
deleted, and a later question takes the next unused number rather than
renumbering what somebody already answered by number.

With no questions, the unnumbered empty-state line stays as written; a
placeholder numbered item would claim a decision is pending when none is.

Answers are written into the record by the agent, not by a command. Only a gate
needs an operation that records a decision and its date — see
[authorization boundaries](authorization.md) — and answering a question is not a
gate but the material a person approves against. Numbering stays human-facing
and never becomes a machine identifier.

A question that surfaces later, during planning, returns to the intent when it
bears on the outcome; a plan never silently settles one. Whether answering it
reopens approval is the ordinary test — renewed approval only when the outcome
or its success criteria materially change.

Owner:

- `context-circuit-source@product/AGENTS.md.in` — what an intent is written
  against and how a question reaches the person.
- `context-circuit-source@internal/workspace/records.go` `Store.CreateRecord` —
  the scaffold that seeds the sections.
