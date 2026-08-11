# Context Circuit

Context Circuit is a team-owned, reusable project wrapper for coordinating
AI-assisted delivery across one or more Git repositories. It keeps workflow
policy, durable context, plans, verification evidence, and host-neutral agent
behavior separate from product code.

The wrapper works with Codex and Claude Code. Plans and activity integrations
are optional; review, merge, deployment, and destructive cleanup remain human
decisions.

## Choose your path

- **Use Context Circuit for a project:** read the
  [wrapper getting-started guide](docs/getting-started.md).
- **Develop or package Context Circuit itself:** read the
  [maintainer guide](docs/development.md).
- **Understand the deterministic command:** read
  [the wrapper command reference](docs/command.md).

The downloadable wrapper requires Git and Node.js 22 or newer. It includes the
standalone `.agents/bin/cc.mjs` command and does not require npm installation,
TypeScript, a database, a background service, or a provider SDK.

## Core workflow

1. `$initialize-workspace` interviews the human and bootstraps new, cloned,
   existing, or submodule repositories safely.
2. `$create-plan` optionally records reviewed delivery intent.
3. `$whats-next` recommends one source-backed action without changing state.
4. `$run-task` prepares isolated worktrees for scoped workers and independent
   verifiers.
5. `$finish-work` records durable outcomes after merge or abandonment.
6. `$sync-context` curates reusable learning through a reviewable wrapper change.

Canonical behavior lives under `.agents/`. `.codex/` and `.claude/` contain thin
host adapters only.
