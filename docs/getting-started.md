# Getting started as an AI agent

Context Circuit is entered through a host agent session. The user should be
able to say:

> Start or resume work in this workspace.

## 1. Enter and orient

Read AGENTS.md, WORKFLOW.md, workspace.yaml, context/INDEX.md, relevant
Product Knowledge, active runtime sessions, plans, repository instructions,
and observed Git state.

Identify whether this is:

- a new workspace or project;
- a new PRD or source;
- an unfinished root session;
- a delegated child session;
- an approved plan ready to execute;
- a verification, review, or blocked session.

Explain the selected route before consequential action.

## 2. Start with a PRD or source

When the user provides a PRD, source document, or repository evidence:

1. Read the source.
2. Separate stated requirements, current behavior, assumptions, unknowns, and
   evidence.
3. Draft concise Product Knowledge with source references.
4. Present the draft for human confirmation.
5. Do not silently replace accepted context.

If no useful source or repository exists, ask focused discovery questions before
creating Product Knowledge or a plan.

## 3. Create and approve work

When context is sufficient:

1. Propose one or more plans by coherent domain or repository boundary.
2. Include dependencies, acceptance criteria, implementation scope, test scope,
   and verification commands.
3. Present plans as drafts.
4. Wait for explicit human approval.
5. Do not infer approval from conversation tone, tests, or agent output.

## 4. Execute with sessions

After approval, the root session may:

- claim a plan execution;
- create or reuse its exclusive worktree;
- delegate research, implementation, or verification subagents;
- continue through dependency-ready tasks;
- iterate between implementation and verification within approved scope;
- pause when scope or assumptions materially change.

Every child session receives a bounded delegation packet and returns a structured
handoff to its parent.

## 5. Resume and finish

At the end of each session, record:

- session and parent identifiers;
- objective and current route;
- evidence and decisions;
- changed files and tests;
- blockers and open questions;
- next safe action.

The next session reads this handoff instead of relying on conversation history.
Human review and explicit status changes remain required before work is
considered complete.
