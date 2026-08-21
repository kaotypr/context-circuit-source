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
