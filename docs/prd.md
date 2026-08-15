# PRD contract

A PRD is a human-reviewed product definition for work that needs durable
requirements and acceptance criteria. It may be created directly from a clear
request or from an accepted Idea Brief.

## Artifact home

Store drafts and accepted PRDs at:

`contributions/prds/<slug>.md`

The artifact remains `status: draft` until the user accepts it. Acceptance does
not approve an implementation plan, merge, publication, or deployment.

## Suggested structure

```markdown
---
kind: prd
status: draft
title: Example
---

# Example

## Problem and desired outcome

The user problem, intended outcome, and why it matters.

## Users and scenarios

Users, primary flows, edge cases, and success signals.

## Requirements

Numbered, testable requirements grounded in accepted intent and evidence.

## Non-goals

Explicitly excluded behavior and scope.

## Acceptance criteria

Observable conditions for human review and later verification.

## Constraints and dependencies

Repository, architecture, policy, timing, or integration boundaries.

## Assumptions and open questions

Unconfirmed statements and questions that may change the plan.

## Product Knowledge references

Links to accepted context used by this PRD.

## Provenance

Selected source paths, why each was read, and revision or freshness when known.
Raw source text stays in `sources/`.
```

Requirements, evidence, assumptions, proposals, and unresolved questions must
not be conflated.

## Next action

After human acceptance, offer a context refresh, a draft plan, or stopping with
the PRD. A plan remains a separate artifact and human approval gate.
