---
type: Workflow
title: "<Workflow title>"
description: "<One sentence describing the workflow outcome.>"
status: draft
slug: "<workflow-slug>"
domain: "<domain-slug>"
sources: []
source_revisions: []
generated_at:
review_date:
freshness: unknown
assumptions: []
unknowns: []
contradictions: []
acceptance:
  state: pending
  gate: product-knowledge-acceptance
  accepted_at:
  accepted_by:
---

# <Workflow title>

## Outcome

The user or system outcome established by the selected evidence.

## Trigger and participants

What starts the workflow and which participants are involved.

## Steps

The supported sequence, decision points, and variations. Keep exact workflow
behavior here rather than duplicating it in a domain index.

## Inputs and outputs

Relevant inputs, outputs, events, or handoffs when supported by evidence.

## Constraints and edge cases

Known boundaries, exceptions, failure modes, and unresolved questions.

## Verification

Tests, examples, or review questions that support this workflow knowledge.

## Provenance

List every selected source and repository document actually read, why it was
relevant, and its revision or freshness signal. Raw source text stays in
`sources/`.

## Acceptance notes

Record proposed material, accepted decisions, unknowns, and contradictions. A
generated workflow remains `status: draft` until the human acceptance gate is
satisfied.
