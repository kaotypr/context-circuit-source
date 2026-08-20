---
name: cc-cleanup-runtime
description: Delete workspace runtime after inspecting worktrees and confirming dirty or unpushed work.
---

# Clean up runtime

Use this skill when the user asks to clear local execution state or to delete
`.runtime/`. `cc-cleanup-runtime` is the named runtime-cleanup skill. Cleanup
is workspace-wide for `.runtime/`; confirmation must say that. It does not
change plan or task status.

## Route reads

Use the `resume` manifest for runtime records and the `entry` manifest for
workspace identity in
`docs/agent-workspace-workflow.md#route-read-manifests`. Cleanup keeps its
local inspection and human-confirmation guard: a missing runtime record or
uncertain ownership is not permission to delete state.

## Inspect

Inspect every runtime worktree under `.runtime/worktrees/` for:

- uncommitted changes, including untracked files;
- commits not present on the worktree's upstream, or local-only commits when
  no upstream exists.

Classify each worktree as `clean`, `uncommitted`, `unpushed`, or both.

Also list live non-terminal sessions. Deleting `.runtime/` removes resume
records. Live sessions do not invent permission to discard Git work.

Reject path traversal, symlinks, and identifiers outside the runtime contract.

## Confirm

If any dirty or unpushed work exists, STOP with a confirmation list of affected
worktrees and the risk:

- uncommitted work is destructive if cleanup proceeds;
- unpushed commits are not deleted by `git worktree remove` because the branch
  remains in the parent repository; still confirm because the checkout is going
  away and the human may want to push first.

If the runtime is clean, still confirm a short summary of sessions and
worktrees that will be removed, then delete after the human chooses cleanup.

`--force` is allowed only after the human confirmed discarding the listed dirty
work.

Missing or already-empty `.runtime/` is a successful no-op.

## On confirmed cleanup

1. Remove each registered Git worktree with `git worktree remove`. Use
   `--force` only after dirty-work confirmation.
2. Delete the remaining `.runtime/` records.

## Safety

Must not modify the base repository checkout, delete Git branches, reset or
stash product repositories, or change plan or task status. A child worker or
verifier cannot clean runtime. Only the root coordinator, after the required
human confirmation, may delete `.runtime/`.
