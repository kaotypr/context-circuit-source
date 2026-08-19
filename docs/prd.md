# PRD contract

A PRD is a human-reviewed product definition for work that needs durable
requirements and acceptance criteria. It may be created directly from a clear
request or from an accepted Idea Brief.

## Artifact home

Store drafts and accepted PRDs at a user/team-selected path under `sources/`.
Context Circuit does not prescribe subdirectories or filenames. Record the
exact chosen path in the artifact's Provenance section.

The artifact remains `status: draft` until the user accepts it. Acceptance does
not approve an implementation plan, merge, publication, or deployment.

## Generated envelope

New writers emit the Context Circuit profile envelope below. `status` is the
OKF lifecycle value; `acceptance` is a separate human gate and does not approve
an implementation plan, publication, or deployment.

```markdown
---
type: Product Requirements
title: Example
description: A concise statement of the product definition suitable for retrieval.
status: draft
acceptance:
  state: pending
  gate: prd-acceptance
---
```

The artifact is an OKF concept only when its user-selected location is inside
an explicitly declared bundle. A path under `sources/` does not declare the
whole source inbox as a bundle. Without that declaration, validate the
artifact schema and keep the OKF checks `not-applicable`.

## Suggested structure

```markdown
---
type: Product Requirements
title: Example
description: A concise statement of the product definition suitable for retrieval.
status: draft
acceptance:
  state: pending
  gate: prd-acceptance
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

New output omits `kind`, unsupported sections, and generation instructions.
Compatibility readers may normalize a legacy `kind` only when `type` is absent;
conflicting `kind` and `type` values fail validation. Unknown top-level and
nested extension fields are preserved during a refresh.

## Generation readiness

Before presenting a PRD as ready, the host's
`deterministic-yaml-schema-validation` capability must:

1. extract the initial front matter at line 1 through the first later standalone
   `---`;
2. parse that block as one deterministic YAML document;
3. validate the PRD artifact schema and, only for an explicitly declared
   bundle, the base OKF and Context Circuit profile; and
4. rerun every applicable check after repairing a hard failure.

Advisory link, freshness, and copy findings remain separate from hard results.
If the host cannot provide the capability, the artifact remains not-ready.

## Next action

After human acceptance, offer a context refresh, a draft plan, or stopping with
the PRD. A plan remains a separate artifact and human approval gate.
