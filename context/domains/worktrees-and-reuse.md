# Worktrees and environment reuse

Isolation cheap enough to actually use.

## Optional, and chosen by the agent

Worktrees are not mandatory. The agent picks a working strategy from repository
instructions, existing work, task needs, and stated preference, recommending
isolation for risky or parallel work; the executable performs the Git mechanics.
Routine preparation covered by the request gets no extra permission gate.

For a plan the defaults are a branch and path named from the plan and repository
IDs under the workspace. Explicit branch, start, and path values override.
Preparation without a plan works with an explicit branch and path.

An unspecified start uses the recorded local base branch if available, otherwise
its remote tracking branch. **No fetch happens implicitly**, and fetching
implies neither rebase nor reset. An explicit start may be another plan's branch
or a commit, resolved before the worktree is created; it applies to new branches
only, since a reused branch keeps its commits.

## The preservation rules

Every one is a refusal, and together they are what make isolation safe to
automate:

- Preparation never resets an existing branch, force-checks-out, stashes, or
  overwrites unrelated files.
- A branch or path already holding work is inspected and then resumed, or a new
  location is chosen — never silently taken over. An explicit reuse flag is how
  a caller selects an existing branch or adopts an existing association.
- Repeating preparation for the same association returns its actual Git state
  rather than recreating anything.
- Removal preserves dirty, untracked, and ignored files unless an explicit
  discard flag carries authorization to dispose of them, and that flag is never
  added merely to get past a failure.
- Removal never touches the primary checkout or a locked worktree, and keeps the
  branch. Branch deletion is separate ordinary Git.
- Marking work complete removes neither worktree nor branch, and no automatic
  cleanup or pruning follows failure or completion.
- Base drift triggers nothing: no automatic rebase, review, or repair loop.

## Git is the inventory

Local plan-to-worktree associations are a convenience for resuming work. They do
not prove a path exists or still holds the same branch. Resume by listing and
inspecting through Git, then read the real diff. Repair reconnects a manually
moved working copy; move relocates one through Git. The diagnostic reports stale
associations.

## Copy-on-write environment reuse

The most common reason people avoid worktrees is that a fresh working copy means
a fresh dependency install and a missing environment file, so isolation gets
used when it is least convenient and skipped when it matters.

After Git creates or selects the worktree, preparation discovers ignored
dependency directories and environment files in the bound checkout and clones
them with the filesystem's own mechanism — `clonefile` on macOS, `FICLONE` on
Linux, aligned block cloning on Windows ReFS. The default mode falls back to
independent copies; other modes require cloning, force copying, or disable reuse
entirely. Additional ignored runtime entries are selected explicitly, must be
repository-relative and ignored in both source and target, and discovery does
not descend into arbitrary ignored parent directories.

Safety properties: entries are staged and published only when complete, so a
failure preserves the worktree and finished entries and a rerun resumes;
existing destination files are preserved and never overwritten; internal
symlinks are translated while external symlinks and special files are rejected;
hard links to the source are never created; paths naming control directories are
refused; and environment contents are copied opaquely — never printed, never
placed in a prompt, never stored in a shared record.

Before reusing dependencies the executable compares the tracked package
manifests, lockfiles, workspace definitions, and runtime version files between
source and target. Any difference, including a file present on one side only,
skips dependencies with a reason.

## What reuse does not prove

Matching inputs do not prove the source installation is complete, or that its
runtime and native modules are compatible with the target. The tree is not an
atomic snapshot, so a package install running during the copy produces an
inconsistent result. Local files do not make native dependencies portable across
operating system, architecture, runtime ABI, or container. Cloning starts no
services, migrates no databases, and makes no local port unique. The correct
posture is to reuse a working local environment, read the reuse report before
application setup, and run the project's normal checks.

Owner: `context-circuit-source@internal/workspace/worktrees.go` and
`reuse.go`; platform cloning in `context-circuit-source@internal/cow/`.
