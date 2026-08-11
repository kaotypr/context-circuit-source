# Decisions

## 2026-08-12 — Version 0.2.0 operating decisions

- Use the team-owned product name Context Circuit and the host-neutral
  `.agents/bin/cc.mjs` command. Keep personal names out of product identifiers.
- Ship a neutral template with no placeholder repository. Interactive
  initialization gathers human choices, presents exact actions, and can safely
  initialize the configured wrapper and product Git roots after explicit commit
  authorization.
- Make the configured wrapper state its first commit. Give a new product
  repository an empty base commit instead of invented application content.
- Track intentional empty directories with `.gitkeep`, except `repositories/`,
  which initialization creates only when needed.
- Keep detailed protocols, host proofs, examples, and maintainer documentation
  in the source repository. Ship only human-facing getting-started, usage,
  configuration, and command-reference guides in the wrapper.
- Treat the generated `context-circuit-0.2.0.tar.gz` as the authoritative release
  artifact. Exclude macOS metadata and redundant version markers, and retain the
  generated directory only as an inspectable staging view.
- Support personal and team wrappers; recommend team mode. Team wrapper changes
  use pull requests; explicitly configured solo mode may use direct commits.
- Support Codex and Claude Code with one canonical workflow and thin adapters.
- Keep plans and activity providers optional. Direct requests receive local
  `ADHOC-*` IDs without creating a task database.
- When exclusive claiming is unavailable, continue only after explicit selection
  and record the duplicate-effort risk; isolated Git state prevents corruption
  but does not promise exclusive ownership.
- Keep `whats-next` read-only and separate from selection, publication, claiming,
  lifecycle mutation, and execution.
- Require explicit plan approval. Material changes revoke approval and increment
  the version while preserving stable work IDs; non-material repairs do not.
- Publish tasks only after explicit authorization and confirmed duplicate-safe
  discovery. Never infer external success.
- Use one branch/worktree and fresh scoped worker per repository, followed by a
  different read-only verifier. Resolve contracts before dependent work.
- Require humans to merge code and invoke closeout after merge or deliberate
  abandonment. Never merge or deploy automatically.
- Preserve runtime state until contribution evidence is durable and cleanup is
  proven safe. Cleanup never deletes branches or runtime evidence.
- Support ignored clones and tracked submodules; recommend ignored clones.
- Ship a dependency-free bundle and exclude maintainer tooling, fixtures, tests,
  package metadata, and development planning material from the template.

Append new durable decisions with date, context, and rationale.
