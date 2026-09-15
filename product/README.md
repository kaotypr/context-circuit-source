# Context Circuit

A shared workspace for coordinating AI-assisted changes across one or more Git
repositories, for solo developers and teams. Speak to your coding agent in
ordinary language. A Go executable maintains the files and working copies.

## Start a workspace

Clone this workspace template and open it in your coding agent. Ask it to install
Context Circuit CLI using the included `cc-cli` skill. The CLI is a separate
product with its own releases; the skill selects the current execution platform,
verifies the download, and installs without administrator access. Ask the same
skill to update or roll back the CLI later. Workspace records remain unchanged.

Installed Git is required for repository operations. No Go or Python setup is
needed. Initialize the cloned blank template, or ask your agent to do it:

```sh
context-circuit-cli --workspace . init --name Acme \
  --purpose 'Billing software' --member maya --member-name Maya
```

Open the resulting workspace in your coding agent. Its shared instruction is
`AGENTS.md`; Claude and Cursor entry files point to it. Workspace data is readable
YAML and Markdown. Share the workspace through Git if useful; local checkouts,
member selection, and worktrees are ignored. Each OS, remote host, or container has its own CLI and local bindings.
The CLI can also export a blank seed or initialize a compatible blank template.

> Connect `../billing-api` as api, with main as its base branch here.
> Clone our web repository here and connect it as web, based on main.
> Initialize a documentation repository.
> The web repository consumes the API from api.

Existing checkouts, cloned repositories, new repositories, and using the workspace
root as a repository are supported. Shared definitions record logical repository
IDs with their URL and default branch. Each machine binds those IDs to its own
paths and to the branch it starts work from.

## Shared project knowledge

> Gather the billing rules from `sources/billing-requirements.md`.

The agent writes durable architecture, conventions, decisions, terminology, and
domain knowledge into `context/`. Notes and an index are optional. Relevant
knowledge is retrieved selectively; sources remain separate passive evidence.

## Build across repositories

> Add recurring billing to the API and web app.

The agent writes an intent defining the desired outcome and success criteria, then
waits. Approve or refine it. Approval is what sends the agent into the code to
create linked Markdown plans, which can span repositories or be split with
dependencies. It shows you the plans and waits again. Reading them is optional;
nothing is built until you ask for execution.

> Execute all plans of the billing intent.

For several plans at once, the agent derives the dependency waves, shows you
whether overlapping them or chaining them costs less, and then runs to completion
unattended — preparing worktrees, merging a dependent plan's base when needed,
dispatching workers, and recording progress. It stops and preserves everything on
a failed check or a decision it should not make alone.

Executing a plan prepares a worktree per repository first, unless you ask for the
change to be made directly in a checkout you already have. The executable prepares
the Git working copies; it reuses ignored node_modules and .env files using
filesystem CoW when available, with independent-copy fallback. The agent handles
any remaining setup and implements the change. A dispatcher skill supports explorer, planner, worker, and manually
requested reviewer subagents, with configurable per-host model/effort settings.
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

Members share plans. `created_by` records attribution; no assignee, owner, or
reviewer appears on plans. IDs look like `i001-add-billing`, `p0001-billing-api`,
and `p0002-billing-web`. Reservations survive archival and deletion.

Give each member of a team an **allocation band** and they can create records in
separate clones without ever choosing the same number:

> Give Maya band 1 and Alex band 2.

Maya then allocates `i100`, `p1000`, `p1001`; Alex allocates `i200`, `p2000`. A
solo workspace needs no band and keeps counting from `i001` and `p0001`. Bands
help only where they are actually assigned: unbanded members in separate clones,
or a clone with a stale roster, still have to synchronize and resolve competing
allocations. No distributed allocation service is included.

See [workspace files](.context-circuit/docs/workspace.md),
[commands](.context-circuit/docs/commands.md),
[worktrees](.context-circuit/docs/worktrees.md), [subagents](.context-circuit/docs/agents.md), and
[working records](.context-circuit/docs/working.md).
V2 supports fresh workspaces; it does not overwrite or automatically migrate v1 data.
