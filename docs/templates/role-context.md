---
type: Role
title: "<Business or project role>"
description: "<One sentence describing this role's perspective.>"
status: draft
slug: "<role-slug>"
role_type: business
owners: []
domains: []
related_workflows: []
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

# <Business or project role>

## Role definition

Who this business or project role is, what perspective it brings, and the
boundary of its responsibility. Business/project roles are Product Knowledge;
agent execution roles such as coordinator, worker, and verifier belong under
`agents/` and are not generated here.

## Primary outcomes

The outcomes this role is trying to achieve across the product.

## Product surfaces

The screens, menus, APIs, events, or channels this role uses, when supported by
the selected evidence.

## Cross-domain perspective

Explain how this role connects the canonical domains listed in `domains:`.
Describe the role's journey and handoffs without copying domain behavior,
business rules, or workflow steps.

## Related domains and workflows

Link to the canonical domain documents and only the workflow pages relevant to
this role. The simple `domains` metadata list is the relationship model; do not
introduce contribution, action, consumption, or approval matrices.

## Role-specific behavior and limitations

Capture behavior or constraints that are genuinely specific to this role. Link
to the domain or workflow page that owns exact behavior.

## Provenance

List every selected source and repository document actually read, why it was
relevant, and its revision or freshness signal. Raw source text stays in
`sources/`.

## Acceptance notes

Record proposed material, accepted decisions, unknowns, and contradictions. A
generated role remains `status: draft` until a human accepts it.
