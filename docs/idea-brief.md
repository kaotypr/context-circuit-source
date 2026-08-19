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

## Generated envelope

New writers emit the Context Circuit profile envelope below. `status` is the
OKF lifecycle value; `acceptance` is a separate human gate and is not inferred
from schema validation.

```markdown
---
type: Idea Brief
title: Example
description: A concise statement of the idea suitable for retrieval.
status: draft
acceptance:
  state: pending
  gate: idea-brief-acceptance
---
```

The artifact is an OKF concept only when its user-selected location is inside
an explicitly declared bundle. A path under `sources/` does not declare the
whole source inbox as a bundle. Without that declaration, validate the
artifact schema and keep the OKF checks `not-applicable`.

## Suggested structure

```markdown
---
type: Idea Brief
title: Example
description: A concise statement of the idea suitable for retrieval.
status: draft
acceptance:
  state: pending
  gate: idea-brief-acceptance
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

New output omits `kind`, unsupported sections, and generation instructions.
Compatibility readers may normalize a legacy `kind` only when `type` is absent;
conflicting `kind` and `type` values fail validation. Unknown top-level and
nested extension fields are preserved during a refresh.

## Generation readiness

Before presenting an Idea Brief as ready, the host's
`deterministic-yaml-schema-validation` capability must:

1. extract the initial front matter at line 1 through the first later standalone
   `---`;
2. parse that block as one deterministic YAML document;
3. validate the Idea Brief artifact schema and, only for an explicitly declared
   bundle, the base OKF and Context Circuit profile; and
4. rerun every applicable check after repairing a hard failure.

Advisory link, freshness, and copy findings remain separate from hard results.
If the host cannot provide the capability, the artifact remains not-ready.

## Next action

After human acceptance, the user may stop, refresh context, create a PRD,
create a plan, or begin a small implementation. No later artifact is required
by this contract.
