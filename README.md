# Context Circuit

Context Circuit is a team-owned, reusable project wrapper for coordinating
AI-assisted delivery across one or more Git repositories. It keeps workflow
policy, durable context, plans, verification evidence, and host-neutral agent
behavior separate from product code.

The wrapper works with Codex and Claude Code. Plans and activity integrations
are optional; review, merge, deployment, and destructive cleanup remain human
decisions.

## Get started

Start with the [wrapper getting-started guide](docs/getting-started.md). Then use
the [human workflow guide](docs/using-the-wrapper.md) for day-to-day delivery,
the [configuration reference](docs/configuration.md) for workspace policy, the
[command reference](docs/command-reference.md) for troubleshooting or advanced
direct use, and the [Product Knowledge guide](docs/product-knowledge.md) for the
optional business-context capability.

The downloadable wrapper requires Git and Node.js 22 or newer. It includes the
standalone `.agents/bin/cc.mjs` command and does not require npm installation,
TypeScript, a database, a background service, or a provider SDK.

## Core workflow

1. `$cc-configure-workspace` configures fresh or existing wrappers and safely uses
   an internal bootstrap phase only for first-time Git and base commits.
2. `$cc-create-plan` optionally records reviewed delivery intent.
3. `$cc-whats-next` recommends one source-backed action without changing state.
4. `$cc-run-task` prepares isolated worktrees for scoped workers and independent
   verifiers.
5. `$cc-finish-work` records durable outcomes after merge or abandonment.
6. `$cc-sync-context` curates reusable learning through a reviewable wrapper change.

Canonical behavior lives under `.agents/`. `.codex/` and `.claude/` contain thin
host adapters only. `$cc-initialize-workspace` remains a state-detecting compatibility
alias.
