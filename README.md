# Context Circuit v2

Context Circuit is a shared workspace station for AI-assisted development across
one or more Git repositories. Solo developers and teams share project knowledge,
repository relationships, members, intents, and plans.

The separately released Context Circuit CLI manages workspace files, global IDs,
repository bindings, CoW worktrees, dependency-ordered execution, knowledge
consistency, and native subagent role settings. The workspace ships skills to
install/update the CLI and dispatch subagents. Your coding agent handles
understanding, planning, implementation, environment setup, ordinary checks, and
explicitly requested review and delivery.

## User journey

1. Clone the workspace template, ask its cc-cli skill to install the CLI, and
   initialize a named workspace with a purpose and first member.
2. Connect existing repositories, clone them, or initialize new repositories.
   The workspace root itself can be a repository. Record default base branches
   and relationships; keep concrete checkout paths local to each machine.
3. Gather durable knowledge from named sources into `context/` notes, catalogued
   one unwrapped entry each in `context/INDEX.md`, with project vocabulary in
   `context/glossary.md`. A note describes the project and anchors to repository
   paths; it never names a record or a file of raw evidence.
4. Describe a change; the agent writes an intent and stops. Approving it is what
   sends the agent into the code to create linked Markdown plans, which it
   presents and stops on again. Reading them is optional and there is no plan
   approval gate, but nothing is implemented until execution is requested.
5. Ask to execute. That prepares a worktree per repository — unless you ask to
   work directly in a bound checkout — then implements and runs normal tests,
   lint, and builds. For several plans at once,
   `record order` derives dependency waves or a linear chain and reports the
   cost of each; the agent confirms the shape once and runs it to completion,
   preparing worktrees, performing integration merges, dispatching workers, and
   stopping with all work preserved on a failed check or a decision it should
   not make alone.
6. Request delivery and independent code review separately. Review reports
   findings without modifying code. It never starts automatically during
   execution, and never blocks a pull request, delivery, or completion.
7. Explicitly mark plans done and reconcile relevant durable project knowledge.
   Completion returns the catalog entries scoped to those plans' repositories as
   candidates to judge; where meaning changed, the note and its entry move
   together.

Member attribution is created_by only. Plan IDs are workspace-global, never
member namespaces. The p prefix distinguishes plans (p0001) from intents (i001).
Numbers are allocated automatically; an optional per-member allocation band
divides the range so members holding distinct bands never choose the same
number, however long they work in separate clones.

`context-circuit-cli check` is an explicitly invoked diagnostic over records,
local bindings, dependency cycles, stale worktree associations, and knowledge
consistency — a note that crosses the durable-only boundary, a catalog entry
naming a missing note, a note no entry lists. It is not an execution gate and
nothing waits on it.

## Build and use

Go 1.25+ is needed by contributors. Users need the executable for their platform
and installed Git; no Python, Go toolchain, or YAML package installation is
needed.

```sh
go build -o /tmp/context-circuit-cli ./cmd/context-circuit
/tmp/context-circuit-cli --workspace /tmp/acme-new init \
  --name Acme --purpose 'Billing software' --member maya --member-name Maya
```

Native binary archives are built for macOS, Linux, and Windows on amd64/arm64.
With no output argument each build clean-rebuilds its own directory under
`dist/` (`dist/workspace-<version>` and `dist/cli-<version>`), so repeated runs
replace rather than fail and neither build removes the other's assets. An
explicit output directory must be new; builds never replace existing output
there:

```sh
sh scripts/build-dist.sh
sh scripts/build-cli.sh
sh scripts/build-dist.sh v2.0.0-dev /tmp/cc-v2-workspace
sh scripts/build-cli.sh 2.0.0-dev /tmp/cc-v2-cli
sh scripts/check-release.sh
```

Workspace publication uses VERSION and publishes the template repository's v*
releases. CLI publication uses CLI_VERSION and this repository's
context-circuit-cli-v* tags, which build here and publish their assets on the
template repository and its GitLab mirror, so installing needs no access to this
checkout. The two products have separate release workflows and package inventories.
Each workspace pins the CLI version it expects and versions install side by
side, so workspaces pinning different versions coexist on one machine. The CLI
also embeds a seed as a convenience for new workspaces; updates do not rewrite
existing workspaces. See [CLI product](CLI.md).

Pass a new directory to `check-release.sh` to retain both sets of checked
assets.

See [the product guide](product/README.md),
[workspace files](product/docs/workspace.md),
[commands](product/docs/commands.md),
[working records](product/docs/working.md),
[subagents](product/docs/agents.md),
[worktree responsibilities](product/docs/worktrees.md), and
[source workflow](WORKFLOW.md).

Workspace data remains readable YAML and Markdown. Local locking coordinates
edits in one directory, and allocation bands keep members holding distinct bands
out of each other's numbers offline. Neither is a distributed allocation
service: unbanded members in separate clones, and any clone whose roster is
stale, must still synchronize the workspace and resolve competing allocations
before sharing new IDs. Initialization is for fresh workspaces; existing v1
workspaces are not migrated automatically.

The YAML dependency is goccy/go-yaml, pinned in go.mod. A portable file-locking
library supplies the small operating-system-specific locking primitive. Release
assets include dependency licenses. Product history and maintainer data never
ship.
