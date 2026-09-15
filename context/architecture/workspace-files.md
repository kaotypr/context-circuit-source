# Workspace files and safe editing

A workspace is readable YAML and Markdown. What separates the files is whether
they travel with the project or belong to one machine.

## Shared versus local

**Shared** (committed): the workspace record carrying schema version, name,
purpose, logical repository IDs with their default base branches, and recorded
relationships; the member roster with each member's optional allocation band;
the permanent ID ledger; the per-host role settings; intent and plan records;
and the knowledge tree.

**Local** (ignored): the active member on this machine; the map from logical
repository ID to this machine's checkout path; optional plan-to-worktree
associations; the lock file; and the working copies themselves.

The split is what lets one workspace travel. Cloning it onto a second machine is
selecting an existing member and connecting the existing repository IDs to local
paths — never a second initialization. Each execution environment (native
Windows, a container, a remote host) gets its own checkout, its own installed
executable, and its own bindings; installed dependencies and uncommitted work do
not migrate with the shared files.

A local binding holds a path only, including `.` when the workspace directory is
itself a bound repository. A new repository can be registered before its first
commit, though worktree preparation needs one. Setting a default base records
the intended branch without creating or resetting it.

## Editing guarantees

- **Atomic single-file replacement** plus a portable operating-system lock
  protect cooperating commands in one directory. The lock is released on process
  exit including abnormal exit; the empty lock file remains.
- **Structured edits are document edits** through the pinned YAML library, so
  surrounding comments and ordering survive. Presentation may normalize;
  byte-for-byte preservation is not promised.
- **Unknown fields and duplicate mapping keys in structured control records are
  errors.** User prose stays ordinary Markdown with no schema registry.
- **Multi-file operations may leave partial output** on interruption, and a
  number is reserved before its file is created and never rolled back. The
  command reports the path so the work can be resumed or repaired — silently
  reusing a reserved number would be worse.

## Dates

A record's instants — `created_at`, and the `approved_at` or `completed_at` a gate
stamps — are canonical ISO 8601 UTC timestamps, `2026-09-15T10:53:00Z`. An event
records when it happened rather than a word saying that it did, so nothing can
disagree with it, and two decisions on one day stay distinct. Every other date is
an ISO 8601 calendar date, `YYYY-MM-DD`, in UTC.

Both are written bare in structured records and unadorned in prose. The reason is
retrieval and comparison: a catalog entry's reviewed date and a plan's completion
instant must sort without a parser, and a localized format in a shared file
breaks both.

Owner: `context-circuit-source@internal/workspace/store.go` and
`workspace.go`; the shipped description is
`context-circuit-source@product/docs/workspace.md`.
