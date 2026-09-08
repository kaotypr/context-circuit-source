---
kind: domain
status: proposed
title: "<Domain title>"
slug: "<domain-slug>"
owners: []
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
  accepted_at:
  accepted_by:
workflows: []
---

# <Domain title>

## Summary

What this bounded project area covers and when an agent should route a request
here. Keep this page short; exact behavior belongs in its workflow pages.

## Scope

What is inside and outside this domain. Name the vocabulary or ownership
boundary that makes this a useful domain rather than an arbitrary category.

## Behavior

Current behavior supported by the selected evidence. Omit this section when the
request does not establish behavior yet.

## Workflows

Link only to workflow pages that were supported by the selected evidence. A
workflow page owns its exact steps, rules, and variations.

## Interfaces

Relevant screens, commands, APIs, events, or other product surfaces. Omit
unsupported categories rather than filling them with placeholders.

## Data

Relevant entities, fields, lifecycle, or ownership facts when the evidence
supports them.

## Constraints and edge cases

Known constraints, exceptions, failure modes, or unresolved boundaries.

## Implementation references

Shipped contract, adapter, or skill paths that implement this domain, when
verified. These references do not replace the domain's product behavior.
Invariant IDs and shipped paths remain nameable (INV-KNOWLEDGE-03).

## Verification

Tests, checks, examples, or review questions that support this knowledge.

## Knowledge-change recording

When this domain changes what is now true about the product, add a
`DECISIONS.md` entry with:

- **Decision:** what is now true
- **Rationale:** why
- **Consequence:** what follows for the product

Do not name a particular plan, intent file, or sources file, and do not name a
`sources/` path. Path patterns that explain product structure (for example
`intent/<id>/`) remain allowed.

## Provenance

Optional retrieval pointers to shipped contract or adapter paths and invariant
IDs. Do not name a particular plan, intent file, or sources file, and do not
name a `sources/` path. Raw source text is never copied into this document.

## Acceptance notes

Record what remains proposed, what a human accepted, and any contradiction that
requires review. A generated document remains `status: proposed` until the
human acceptance gate is satisfied.

<!--
Omit sections and metadata categories that the selected evidence does not
support. Preserve `assumptions`, `unknowns`, `contradictions`, `freshness`, and
`acceptance` even when their lists are empty so uncertainty remains visible.
-->
