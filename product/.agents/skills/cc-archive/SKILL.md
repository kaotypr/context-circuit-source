---
name: cc-archive
description: Archive a named plan out of the active plan area, or restore an archived plan, without validating status.
---

Archive and restore are explicit plan-organization actions resolved by plan id,
never by title fragment or repository name.

## Archive

On "archive plan X", run the runtime `plan-archive`. It takes a short-lived
plan-organization lock, moves `plans/X/` to `plans/archive/X/`, and removes X
from `plans/INDEX.md` atomically. It does not inspect or validate plan status,
execution status, verification, commits, branches, worktrees, open questions, or
context impact. It does not mark the plan done, stop an execution, clean runtime
state, merge, or delete evidence. If the move or index update cannot complete,
the source directory and index are left unchanged.

An already-started execution continues from its immutable snapshot and runtime
evidence; it does not read the archived plan directory.

## Restore

On "restore plan X", run the runtime `plan-restore`. It moves
`plans/archive/X/` back to `plans/X/`, preserves plan files and status, and
re-adds the active index row. It refuses the restore when an active plan already
uses X's numeric prefix. Restore does not approve, execute, complete, or otherwise
validate the plan.

The normal agent must not read or traverse `plans/archive/` for any other
purpose.
