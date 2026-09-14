# Context Circuit

A shared workspace for coordinating AI-assisted changes across one or more Git
repositories, for solo developers and teams. Speak to your coding agent in
ordinary language. A Go executable maintains the files and working copies.

## Start a workspace

Download the executable for your operating system and architecture. Put
`context-circuit` on PATH (Windows: `context-circuit.exe`); installed Git is the only
external requirement for repository operations. No Go or Python setup is needed.
Initialize a fresh directory, or ask your agent to do it:

```sh
context-circuit --workspace ./acme init --name Acme \
  --purpose 'Billing software' --member maya --member-name Maya
```

Open the resulting workspace in your coding agent. Its shared instruction is
`AGENTS.md`; Claude and Cursor entry files point to it. Workspace data is readable
YAML and Markdown. Share the workspace through Git if useful; local checkouts,
member selection, and worktrees are ignored. The CLI can also export a blank seed
for inspection or initialize a seed exported by the same version.

> Connect `../billing-api` as api, with main as its base branch.
> Clone our web repository here and connect it as web, based on main.
> Initialize a documentation repository.
> The web repository consumes the API from api.

Existing checkouts, cloned repositories, new repositories, and using the workspace
root as a repository are supported. Shared definitions record logical repository
IDs and default base branches. Each machine binds those IDs to its own paths.

## Shared project knowledge

> Gather the billing rules from `sources/billing-requirements.md`.

The agent writes durable architecture, conventions, decisions, terminology, and
domain knowledge into `context/`. Notes and an index are optional. Relevant
knowledge is retrieved selectively; sources remain separate passive evidence.

## Build across repositories

> Add recurring billing to the API and web app.

The agent writes an intent defining the desired outcome and success criteria.
Approve or refine it. After approval, the agent investigates the code, creates
linked Markdown plans, and proceeds without separate plan approval. Plans can
span repositories or be split with dependencies.

Use worktrees when helpful. The executable prepares the Git working copies;
the agent follows each repository's setup instructions and implements the change.
Normal tests, linting, and builds remain part of implementation. Progress and
remaining work stay in the plan so a later session can resume from actual Git state.

## Review, deliver, and complete

> Open a PR for the API change.
> Independently review the PR.
> Fix those findings.
> Merge the API PR.
> Mark the billing plans done.

Independent verification is a manually requested read-only code review. It never
runs automatically during execution; you decide whether to fix findings or
proceed. Review after delivery is possible as an audit. Delivery is explicit and
can happen per repository. Marking a plan done is separate and can update relevant
durable knowledge. It does not delete a worktree or branch.

## Team records

Members share plans. `created_by` records attribution; no assignee, owner, reviewer,
or member-specific numbering appears on plans. IDs look like `i001-add-billing`,
`p0001-billing-api`, and `p0002-billing-web`. Reservations survive archival and
deletion. The helper coordinates simultaneous edits in one workspace directory;
separate Git clones must synchronize and resolve competing allocations before
sharing new record IDs. No distributed allocation service is included.

See [workspace files](.context-circuit/docs/workspace.md),
[commands](.context-circuit/docs/commands.md),
[worktrees](.context-circuit/docs/worktrees.md), and
[working records](.context-circuit/docs/working.md).
V2 supports fresh workspaces; it does not overwrite or automatically migrate v1 data.
