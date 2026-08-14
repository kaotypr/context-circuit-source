# Decisions

## 2026-08-14 — Superseded command-centric baseline

The earlier command-centric design described Context Circuit as a human-driven
planning wrapper with no runtime session lifecycle. That baseline is retained
in Git history for reference but is no longer the target behavior.

## 2026-08-14 — Agent workspace model

- Context Circuit is an AI-agent workspace, not a CLI users must operate.
- A root session may coordinate child sessions and multiple plans.
- Runtime session state is private, filesystem-based, resumable, and separate
  from durable Product Knowledge and plan intent.
- One plan has at most one active writing owner, and each writable session has
  an exclusive worktree.
- Human approval remains required for consequential gates.
- Agents must report evidence, decisions, assumptions, blockers, and next
  actions separately.
- The canonical workflow is defined by the Agent Workspace Workflow document
  and its normative summary in WORKFLOW.md.

## 2026-08-14 — Pure agent-workspace implementation

- The replacement is rebuilt around agent instructions and inspectable
  filesystem conventions.
- Node/JavaScript command scripts, generated bundles, command-specific host
  adapters, and package-manager entry points are not part of the workflow.
- Runtime state remains under `.runtime/` and is preserved until a human
  explicitly chooses runtime cleanup.
