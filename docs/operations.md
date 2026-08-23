# Operational contract

The system is filesystem-first and offline-capable. A host may provide child
creation, Git inspection, or an integration adapter, but those are optional.
Provider outcomes are `disabled`, `denied`, or `unavailable` and never become
canonical plan state.

The root coordinator owns the human request, route, lease, worktree, child
packets, evidence consolidation, and cards. A writer owns only assigned paths
in its exclusive worktree. A verifier is a separate read-only child and writes
only its own handoff. There is no scheduler, global current pointer, database,
or implicit publication.

## Host operations

Codex CLI reads the shared AGENTS.md surface. Claude Code loads the shipped
CLAUDE.md bridge, which imports AGENTS.md. Cursor Agent CLI may discover both
root surfaces and uses the same two-stage route in interactive, print, and
resume modes. These observations are recorded as host evidence only.

Native child support is optional. A writer still has the exclusive delegated
worktree and a verifier still has independent read-only permissions. A missing
required child is labeled host-blocked; a provider failure is labeled
disabled, denied, or unavailable and falls back to filesystem-only evidence.
Host-local permission files, authentication, MCP settings, and transcripts are
not workspace state.

Optional live smoke checks are opt-in and disposable. Record only the host,
executable, version, and bounded outcome: pass, unavailable, or host-blocked.
The offline semantic suite remains the release gate.
