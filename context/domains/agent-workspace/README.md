---
kind: domain
status: proposed
title: "Agent Workspace"
slug: agent-workspace
owners: []
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
  - "agent-workspace" is the smallest useful umbrella domain for the source's end-to-end workspace lifecycle.
  - Detailed repository-specific behavior remains owned by the registered repository and its local documents.
  - The source's acceptance scenarios describe intended verification coverage; they are not evidence that every scenario has been executed.
unknowns:
  - Which artifacts should be created automatically and which should always be proposed first.
  - The minimum useful Product Knowledge set for a new workspace.
  - How multiple active plans should be presented without exposing internal session records.
  - The smallest team collaboration model and the minimum host capability contract.
  - Which optional integrations and completion or publication actions should remain human-gated.
contradictions: []
acceptance:
  state: pending
  accepted_at:
  accepted_by:
workflows:
  - context/domains/agent-workspace/workflows/session-lifecycle.md
---

# Agent Workspace

## Summary

The Agent Workspace domain covers the conversational, durable, and resumable
lifecycle of work in Context Circuit. Route requests here when they concern
starting or resuming a session, shaping intent, grounding work in Product
Knowledge, creating or executing a plan, verifying results, or leaving a
handoff.

## Scope

Inside this domain are the workspace entry contract, source-grounded context,
Idea Brief and PRD handoffs, human-reviewed plans, bounded execution,
verification, durable runtime handoffs, repository boundaries, and meaningful
human gates.

Outside this domain are the implementation details of a product repository,
the contents of raw sources themselves, provider-specific integration state,
and the business or project perspectives represented by Role Knowledge. Those
areas retain their own authority and are linked rather than copied here.

## Behavior

The intended user experience starts with a natural-language request. The agent
reads the minimum instructions, relevant context, active plans, repository
state, and handoffs; identifies the current route; reuses existing artifacts;
states assumptions and unresolved questions; and recommends the shortest safe
next action.

The workspace progressively moves from intent to a suitable artifact, durable
context, a proposed plan, bounded execution, verification, and a durable
handoff. Routine inspection, context loading, task sequencing, local
verification, and handoff creation do not require separate confirmation.

The current architecture baseline is instruction- and filesystem-driven:
Product Knowledge lives under `context/`, raw sources remain a passive
`sources/` inbox, product artifacts are separate from accepted context, plans
are human-reviewed, and runtime session state is kept under `.runtime/`.

## Workflows

The selected evidence establishes the [session lifecycle](workflows/session-lifecycle.md),
which owns the ordered entry, grounding, planning, execution, verification,
and handoff behavior.

## Interfaces

- Natural-language agent sessions, with “start or resume work in this
  workspace” as the primary entry.
- Discoverable skills for session entry, initialization, Idea Briefs, PRDs,
  context gathering, planning, execution, review, next-action routing, and
  optional workspace configuration.
- Reviewable Markdown and YAML artifacts in `sources/`, `context/`,
  `contributions/`, `plans/`, and `.runtime/`.
- Registered repositories and their configured active branches.

## Data

The domain distinguishes raw sources, accepted Project Knowledge, proposed or
accepted Idea Briefs and PRDs, human-reviewed plans and tasks, and private
runtime records such as sessions, leases, worktrees, prompts, and handoffs.
`workspace.yaml` records registered repositories and their roles. These layers
have different authority and must not become duplicate copies of one another.

## Constraints and edge cases

- A workspace with no repository is valid and may begin with an Idea Brief or
  source intake.
- Raw source files are passive and are read only for a request-scoped,
  source-based activity.
- Small changes should remain lightweight; not every request needs every
  artifact type.
- Human judgment is required for accepting or materially revising artifacts,
  approving plans or scope changes, resolving contradictions or ownership
  conflicts, publishing, merging, deploying, and marking intended work done.
- Multi-repository work must not write to unrelated repositories.
- Delegated work needs explicit ownership, scope, permissions, worktree, and
  handoff expectations.
- Optional integrations must remain opt-in and must not become a core
  dependency or store credentials in Product Knowledge or runtime records.
- Contradictory sources, failed checks, changed scope, and unsafe repository
  state must remain visible with a next action rather than being silently
  improvised around.

## Implementation references

- `repositories/context-circuit/context/PROJECT.md` — product and workspace
  boundary.
- `repositories/context-circuit/context/ARCHITECTURE.md` — filesystem-driven
  architecture and authority layers.
- `repositories/context-circuit/context/CONVENTIONS.md` — artifact and
  workflow conventions.
- `AGENTS.md` and `WORKFLOW.md` — wrapper session and safety contract.
- `repositories/context-circuit/.agents/skills/cc-gather-context/SKILL.md` —
  request-scoped Domain and Role Knowledge generation contract.

## Verification

The source provides acceptance scenarios covering natural-language entry,
artifact progression, source-grounded context, solo and team execution,
resumable handoffs, multi-repository boundaries, optional integrations, and
progressive structure. Reviewers should verify the relevant scenario against
the current implementation and distinguish documented intent from executed
checks.

## Provenance

- `sources/context-circuit-product-direction.md` — approved product direction,
  target users, lifecycle, information model, human gates, acceptance
  scenarios, non-goals, and open decisions. Revision:
  `sha256:dd1db2f0fb5f2e4c5fcc41af303760f7a6a8f5b7440a0c6755de8a3437223b58`.
- `repositories/context-circuit/context/PROJECT.md` — current product identity
  and repository/workspace boundary. Revision:
  `sha256:6c1694a5c9f5ad180da24517c1070b806f8ac481681653d230e44e323121c539`.
- `repositories/context-circuit/context/ARCHITECTURE.md` — verified
  instruction, filesystem, Product Knowledge, plan, and runtime boundaries.
  Revision:
  `sha256:e1232761db5efea7b964c615038755d70a1b666851de71299a571c950db6fc6e`.
- `repositories/context-circuit/context/CONVENTIONS.md` — repository and
  documentation conventions. Revision:
  `sha256:05c22bf6b958ef9171c14683d430883b202293994b74f016f7f31a3a5b38c443`.
- `repositories/context-circuit/context/DECISIONS.md` — accepted agent-workspace
  and pure filesystem implementation decisions. Revision:
  `sha256:5de8b173ec6b40ccfb48dc0d4afc74f7c3dc60221987ad69f8f92cefb8f140ce`.

## Acceptance notes

This page is a generated proposal with pending human acceptance. The umbrella
domain slug and the set of open decisions are explicit assumptions from the
source; no contradiction was found among the selected evidence. Acceptance
should confirm that `agent-workspace` is the right routing boundary and that
the linked lifecycle workflow is sufficiently specific.
