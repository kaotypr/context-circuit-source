---
kind: role
status: proposed
title: "Workspace User"
slug: workspace-user
role_type: project
owners: []
domains:
  - ../domains/agent-workspace/README.md
related_workflows:
  - ../domains/agent-workspace/workflows/session-lifecycle.md
sources:
  - sources/context-circuit-product-direction.md
  - repositories/context-circuit/context/PROJECT.md
  - repositories/context-circuit/context/ARCHITECTURE.md
  - repositories/context-circuit/context/CONVENTIONS.md
  - repositories/context-circuit/context/DECISIONS.md
source_revisions:
  - path: sources/context-circuit-product-direction.md
    revision: sha256:dd1db2f0fb5f2e4c5fcc41af303760f7a6a8f5b7440a0c6755de8a3437223b58
  - path: repositories/context-circuit/context/PROJECT.md
    revision: sha256:6c1694a5c9f5ad180da24517c1070b806f8ac481681653d230e44e323121c539
  - path: repositories/context-circuit/context/ARCHITECTURE.md
    revision: sha256:e1232761db5efea7b964c615038755d70a1b666851de71299a571c950db6fc6e
  - path: repositories/context-circuit/context/CONVENTIONS.md
    revision: sha256:05c22bf6b958ef9171c14683d430883b202293994b74f016f7f31a3a5b38c443
  - path: repositories/context-circuit/context/DECISIONS.md
    revision: sha256:5de8b173ec6b40ccfb48dc0d4afc74f7c3dc60221987ad69f8f92cefb8f140ce
generated_at: 2026-08-15T15:12:41+02:00
review_date:
freshness: current-at-generation
assumptions:
  - "workspace-user" is an umbrella project role covering the source's solo builder, technical lead or product owner, and engineering team perspectives.
  - The agent host or integration author is adjacent to this role and is not assigned human workspace-user responsibilities here.
unknowns:
  - Role-specific ownership and permissions for different team arrangements.
  - How team members should view and resolve multiple active plans without exposing runtime internals.
contradictions: []
acceptance:
  state: pending
  accepted_at:
  accepted_by:
---

# Workspace User

## Role definition

The Workspace User is the human project participant who uses Context Circuit to
carry work from intent to a reviewed, verifiable outcome. This umbrella view
covers a solo builder, a technical lead or product owner, and an engineering
team participant. It is a Product Knowledge role, not an agent execution role
or a technical access-control policy.

## Primary outcomes

- Preserve project understanding and decisions across sessions.
- Shape ambiguous intent into the smallest useful durable artifact.
- Review context and plan tradeoffs without reconstructing the project from
  chat history.
- Move approved work through bounded execution and independent verification.
- Retain a clear handoff and a human decision point for material scope,
  ownership, publication, merge, deployment, and completion changes.

## Product surfaces

The role primarily uses the conversational agent session and its discoverable
workspace skills. It encounters reviewable source, context, Idea Brief, PRD,
plan, task, repository, and handoff artifacts, while runtime records remain
implementation state rather than a user-facing workflow. Optional external
activity surfaces are available only when explicitly configured.

## Cross-domain perspective

The Workspace User enters through the [Agent Workspace domain](../domains/agent-workspace/README.md),
where the agent routes the request, reuses durable knowledge, and proposes the
next useful action. The user supplies intent and meaningful judgments as work
moves through context, planning, execution, verification, and handoff. Exact
workflow steps and domain rules remain owned by the linked domain and its
[session lifecycle workflow](../domains/agent-workspace/workflows/session-lifecycle.md).

## Related domains and workflows

- [Agent Workspace](../domains/agent-workspace/README.md) — the canonical
  bounded lifecycle for conversational entry, durable context, planning,
  execution, verification, and handoff.
- [Agent Session Lifecycle](../domains/agent-workspace/workflows/session-lifecycle.md)
  — the workflow through which this role moves across those surfaces.

## Role-specific behavior and limitations

The role should be able to express a request naturally and receive concise
orientation without learning internal paths or session identifiers. It may
accept, revise, or reject proposed artifacts and plans, and it remains the
decision-maker at meaningful human gates. It should not be asked to perform
routine context loading, task sequencing, local verification, or handoff
creation as separate ceremonies.

The source does not define a team permission model, role-specific approval
matrix, or exact behavior when several active plans compete for attention. This
page therefore does not infer those relationships; future role-specific pages
or domain decisions should address them when evidence is available.

## Provenance

- `sources/context-circuit-product-direction.md` — target users, primary user
  experience, human judgment boundaries, and acceptance scenarios. Revision:
  `sha256:dd1db2f0fb5f2e4c5fcc41af303760f7a6a8f5b7440a0c6755de8a3437223b58`.
- `repositories/context-circuit/context/PROJECT.md` — product and workspace
  identity. Revision:
  `sha256:6c1694a5c9f5ad180da24517c1070b806f8ac481681653d230e44e323121c539`.
- `repositories/context-circuit/context/ARCHITECTURE.md` — separation of
  user-facing product artifacts from runtime state. Revision:
  `sha256:e1232761db5efea7b964c615038755d70a1b666851de71299a571c950db6fc6e`.
- `repositories/context-circuit/context/CONVENTIONS.md` — artifact and
  workflow conventions. Revision:
  `sha256:05c22bf6b958ef9171c14683d430883b202293994b74f016f7f31a3a5b38c443`.
- `repositories/context-circuit/context/DECISIONS.md` — accepted human-control
  and agent-workspace decisions. Revision:
  `sha256:5de8b173ec6b40ccfb48dc0d4afc74f7c3dc60221987ad69f8f92cefb8f140ce`.

## Acceptance notes

This page is a generated proposal with pending human acceptance. Acceptance
should confirm that the umbrella `workspace-user` role is useful, that its
scope excludes agent execution roles, and that the linked Agent Workspace
domain is the correct canonical owner of behavior.
