# Project

Context Circuit is an AI-agent workspace. It gives root agents and delegated
subagents the context, plans, runtime state, worktree isolation, and handoff
conventions needed to continue project development across sessions.

The workspace is human-controlled at consequential gates, but the user does
not need to operate an internal CLI. Agents enter the workspace through a
session workflow that reads authoritative context, routes the work, delegates
bounded scopes, executes approved plans, verifies results, and records the next
safe action. The core discoverable set includes `cc-approve-plan`,
`cc-finish-plan`, and `cc-cleanup-runtime` for the named approval, completion,
and runtime-cleanup gates.

Approved-plan execution through `cc-run-plan` directs a writer child and a
verifier child. Solo or team workspace identity does not change that topology.
`solo-local` and `team-review` remain delivery policies.

Product repositories own their code and repository-local conventions. Context
Circuit owns the workspace protocol, durable Product Knowledge, plans, and
private runtime coordination state.

The workspace can begin with no repository. Initialization captures project
identity and repository roles when they exist, while source intake and
conversation-first Idea Brief or PRD work can begin from an idea alone.
