# Idea Brief contract

An Idea Brief is the smallest durable artifact for an uncertain idea or an
early intent discussion. It preserves what the user means without pretending
that requirements or decisions are settled.

## Artifact home

Store drafts and accepted briefs at a user/team-selected path under `sources/`.
Context Circuit does not prescribe subdirectories or filenames. Record the
exact chosen path in the artifact's Provenance section.

The artifact remains `status: draft` until the user accepts it. Acceptance is a
human gate; it is not inferred from a complete-looking document or a passing
test.

## Suggested structure

```markdown
---
kind: idea-brief
status: draft
title: Example
---

# Example

## Intent

What the user is trying to accomplish, in their words.

## Desired outcome

What would be different if the idea succeeds.

## Users and context

Who is involved and what situation motivates the idea.

## Constraints

Known boundaries, repository context, timing, or policy constraints.

## Assumptions

Statements that still need confirmation.

## Open questions

Questions whose answers may change the next action.

## Product Knowledge references

Links to the relevant accepted context files.

## Provenance

Selected source paths, why each was read, and revision or freshness when known.
Raw source text stays in `sources/`.
```

Observed evidence, user intent, assumptions, proposals, and open questions
should remain visibly distinct.

## Next action

After human acceptance, the user may stop, refresh context, create a PRD,
create a plan, or begin a small implementation. No later artifact is required
by this contract.
