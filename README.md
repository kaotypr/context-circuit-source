# Context Circuit

Context Circuit is an AI-agent workspace for understanding a project,
coordinating root and subagent sessions, executing approved work, and
continuing development across sessions.

The user-facing interaction is an agent session:

> Start or resume work in this workspace.

The agent reads the workspace instructions and current runtime state, routes the
request, delegates bounded work when useful, works in isolated repository
worktrees, verifies results, and leaves a durable handoff. The user does not
need to operate an internal command-line interface.

## Workspace layers

- AGENTS.md — safety, authority, and agent behavior.
- WORKFLOW.md — session entry, delegation, human gates, and development loop.
- workspace.yaml — repositories, branches, and workspace configuration.
- context/ — durable, source-cited Product Knowledge and decisions.
- plans/ — human-reviewed intended work and task dependencies.
- .runtime/ — private sessions, leases, prompts, handoffs, and worktrees.
- Registered repositories — product code and repository-local conventions.

## Session model

A root session owns the human request and may coordinate multiple child
sessions. Each child receives an explicit objective, scope, permissions, plan or
task, context references, stop conditions, and handoff format.

Multiple sessions and plans may run concurrently. A plan has at most one active
writing owner, and each writable worktree is exclusive.

## Development loop

The normal loop is:

orient → gather → draft context → draft plan → human approval → execute →
verify → review or repair within scope → handoff → resume or close.

Plans and tasks remain human-controlled. Agents do not infer approval or
completion from tests, Git state, child output, or external systems.

## Documentation

Read the agreed workflow contract in
docs/agent-workspace-workflow.md, then use context/INDEX.md to find the
smallest relevant project context.

The current command layer and its documentation are transitional legacy
material. It may remain during migration, but it is not the intended user
interface and should not be used as the source of the new workflow design.
