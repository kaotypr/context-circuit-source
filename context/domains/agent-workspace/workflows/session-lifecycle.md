---
kind: workflow
status: proposed
title: "Agent Session Lifecycle"
slug: session-lifecycle
domain: context/domains/agent-workspace/README.md
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
  - The source's single end-to-end loop is the primary workflow to preserve here.
unknowns:
  - How the agent should present multiple active plans when more than one is resumable.
  - Which host limitations require a degraded workflow.
contradictions: []
acceptance:
  state: pending
  accepted_at:
  accepted_by:
---

# Agent Session Lifecycle

## Trigger

The user starts or resumes work in natural language, or explicitly asks for a
workspace action such as initialization, context gathering, planning,
execution, review, or configuration.

## Steps

1. Read the workspace instructions and only the relevant Product Knowledge,
   plans, repository state, and durable handoffs.
2. Identify the current request, route, and the smallest useful next artifact
   or action. Reuse existing artifacts before creating new ones.
3. State material assumptions, unresolved questions, current blockers, and the
   accepted scope.
4. Capture or revise the appropriate Idea Brief, PRD, or context proposal when
   the request needs a durable artifact. Keep drafts visibly unaccepted.
5. Ground the work in explicitly selected source material and current project
   knowledge. Separate requested intent, current implementation, decisions,
   inferences, unknowns, and contradictions.
6. Propose a bounded plan with repositories, dependencies, non-goals,
   acceptance criteria, verification, risks, and human gates when
   implementation is needed.
7. After plan approval, execute within the accepted scope, directly or through
   bounded delegated sessions with explicit ownership and isolated worktrees.
8. Verify the intended outcome, report evidence and limitations, and stop for
   human review at publication, merge, deployment, ownership, scope, or
   completion gates.
9. Leave a concise durable handoff so a later session can continue without
   replaying the conversation.

## Variations and stop conditions

- A no-repository workspace can continue with an Idea Brief or source intake.
- A small change can remain lightweight and skip unnecessary artifact types.
- A failed check, contradictory source, unsafe repository state, or ambiguous
  ownership stops the workflow with a visible explanation and next action.
- Optional integrations may assist only when explicitly configured and must
  retain a safe filesystem fallback.

## Verification

Verification checks the approved intended outcome rather than only command
success. The source acceptance scenarios are the current review guide for
entry, grounding, execution, handoff, multi-repository safety, and optional
capabilities.

## Provenance

This workflow is derived from the approved source and the selected repository
architecture and conventions listed in
`context/domains/agent-workspace/README.md`. Raw source text remains in
`sources/`.

## Acceptance notes

The workflow remains proposed pending human confirmation of the sequence,
variation rules, and the two unresolved presentation and host-capability
questions.
