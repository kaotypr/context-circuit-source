# Context Circuit

Context Circuit is an AI-agent workspace for understanding a project,
coordinating root and subagent sessions, executing approved work, and
continuing development across sessions.

The only user-facing interaction is an agent session:

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
- sources/ — user/team-owned raw inputs and authored Idea Brief/PRD artifacts.
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

The workspace is instruction- and filesystem-driven. There is no user-facing
command layer and no Node or JavaScript command runtime. Hosts enter through
the same agent instructions and use the filesystem records as the coordination
surface. Verify implementation work with `sh test/acceptance.sh`. The verifier
must inspect runtime files and instructions directly.

## Foundation skills

Natural-language requests discover the shipped host-neutral skills. These
names are capability contracts around the same filesystem workflow, not
commands the user types. The user-facing interaction remains starting or
resuming work in this workspace.

- `cc-session-entry` — enter or resume a root or child session; see
  `docs/agent-workspace-workflow.md` and `docs/runtime-contract.md`.
- `cc-initialize-workspace` — establish solo/team identity, repositories, roles,
  and default active branches, including a valid no-repository path; see
  `docs/getting-started.md`.
- `cc-idea-brief` — capture uncertain intent without forcing a PRD; see
  `docs/idea-brief.md`.
- `cc-create-prd` — create a source-grounded, human-reviewed PRD; see
  `docs/prd.md`.
- `cc-gather-context` — ground or refresh Domain and Role Knowledge; see
  `docs/product-knowledge.md`.
- `cc-create-plan` — draft a human-reviewed plan; see `plans/README.md` and
  `docs/planning.md`.
- `cc-review-plan` — read-only readiness review; see `docs/plan-review.md`.
- `cc-run-plan` — execute an approved plan through sessions and exclusive
  worktrees from the repository default or active branch; see
  `docs/planning.md`.
- `cc-run-stack` — execute a connected set of already-approved plans from
  parent frozen SHAs with a runtime `graph.yaml` and `progress.yaml`; see
  `docs/planning.md` and `docs/runtime-contract.md`.
- `cc-whats-next` — recommend the smallest safe next action from durable
  state; see `docs/agent-workspace-workflow.md`.
- `cc-configure-workspace` — optional later delivery, host, and integration
  configuration; see `docs/configuration.md`.
- `cc-approve-plan` — human plan-approval gate; see `docs/planning.md`.
- `cc-finish-plan` — human completion status-change gate; see
  `docs/planning.md` and `docs/getting-started.md`.
- `cc-cleanup-runtime` — human-gated `.runtime/` cleanup; see
  `docs/getting-started.md`.
