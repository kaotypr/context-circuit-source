# Context Circuit v2

Context Circuit is a shared workspace station for AI-assisted development across
one or more Git repositories. Solo developers and teams share project knowledge,
repository relationships, members, intents, and plans.

A small Go executable manages workspace files, global IDs, repository bindings,
and worktrees. Your coding agent handles understanding, planning, implementation,
environment setup, ordinary checks, and explicitly requested review and delivery.

## User journey

1. Initialize a named workspace with a purpose and first member.
2. Connect existing repositories, clone them, or initialize new repositories.
   The workspace root itself can be a repository. Record default base branches
   and relationships; keep concrete checkout paths local to each machine.
3. Describe a change and approve its intent. The agent then inspects the code,
   creates linked Markdown plans, and proceeds without a separate plan approval.
4. Prepare worktrees when useful, or work in selected checkouts. Implement across
   repositories in dependency order and run normal tests, lint, and builds.
5. Request delivery and independent code review separately. Review reports findings
   without modifying code. It never starts automatically during execution.
6. Explicitly mark plans done and reconcile relevant durable project knowledge.

Member attribution is created_by only. Plan IDs are workspace-global, never
member namespaces. The p prefix distinguishes plans (p0001) from intents (i001).
Numbers are allocated automatically without member-specific ranges.

## Build and use

Go 1.25+ is needed by contributors. Users need the executable for their platform
and installed Git; no Python, Go toolchain, or YAML package installation is needed.

```sh
go build -o /tmp/context-circuit ./cmd/context-circuit
/tmp/context-circuit --workspace /tmp/acme-new init \
  --name Acme --purpose 'Billing software' --member maya --member-name Maya
```

Native binary archives are built for macOS, Linux, and Windows on amd64/arm64.
Use a new output directory; builds never replace existing output:

```sh
sh scripts/build-dist.sh v2.0.0-dev /tmp/cc-v2-new-dist
sh scripts/check-release.sh
```

Pass a new directory to `check-release.sh` to retain the checked release assets.

See [the product guide](product/README.md), [commands](product/docs/commands.md),
[worktree responsibilities](product/docs/worktrees.md), and [source workflow](WORKFLOW.md).
Workspace data remains readable YAML and Markdown. Local locking coordinates
edits in one directory; independent Git clones must synchronize the workspace
and resolve competing allocations before sharing new IDs. Initialization is for
fresh workspaces; existing v1 workspaces are not migrated automatically.

The YAML dependency is goccy/go-yaml, pinned in go.mod. A portable file-locking
library supplies the small operating-system-specific locking primitive. Release
assets include dependency licenses. Product history and maintainer data never ship.
