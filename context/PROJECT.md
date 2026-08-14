# Project

Context Circuit is an AI-agent workspace. It gives root agents and delegated
subagents the context, plans, runtime state, worktree isolation, and handoff
conventions needed to continue project development across sessions.

The workspace is human-controlled at consequential gates, but the user does
not need to operate an internal CLI. Agents enter the workspace through a
session workflow that reads authoritative context, routes the work, delegates
bounded scopes, executes approved plans, verifies results, and records the next
safe action.

Product repositories own their code and repository-local conventions. Context
Circuit owns the workspace protocol, durable Product Knowledge, plans, and
private runtime coordination state.
