# Getting started

In a released workspace, say “Start or resume work in this workspace.” The
agent reads the compact entry spine and reports identity, repositories,
runtime state, the selected probe, and one safe next action.

## New workspace

Until `workspace.yaml` records `identity.status: accepted`, requests that
would create or change anything are routed back to initialization. This
includes broad requests such as “help me build this”; the agent must not infer
identity, create a plan, or write implementation files. Read-only orientation
remains available.

Initialization records only mode, repositories or project items, roles, and
default branches. Zero repositories is valid. Identity acceptance is a human
gate. Later, selected evidence may produce an Idea Brief, PRD, accepted Product
Knowledge, or direct plan; no artifact is forced when the request is already
clear.

## Plan journey

Create a plan bundle with `plan.yaml`, human-facing `PLAN.md`, and task files.
Review is read-only. Approval changes `draft` → `approved` and tasks to
`ready`; it does not execute. A separate named run request creates the root
session, exclusive lease/worktree, writer delegation, and independent verifier.
After evidence is ready, a separate finish confirmation changes `approved` →
`done`. Runtime and worktrees remain until a later, separately gated cleanup.

## Interrupted work

Resume validates the session, receipt, handoff, wrapper version, primary
evidence, Git state, and ownership. A foreign live owner or ambiguous lease is
read-only and needs a human takeover decision. Dirty work is preserved.

## Human surfaces

Humans normally need only this README, a plan's `PLAN.md`, the conversational
next-action card, and the latest surfaced handoff. Machine records remain
inspectable but are not a manual editing workflow.
