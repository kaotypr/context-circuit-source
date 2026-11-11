---
product: cli
kind: added
---

`repo sync` and `workspace sync` bring this machine's base branches up to date after a merge. `repo sync` fetches, then fast-forwards the base branch of each repository you name with `--id`, or of every bound one, and reports one status per repository: updated, current, ahead, diverged, blocked by uncommitted changes, checked out in another worktree, or why nothing moved: no local base branch, no remote or upstream, not bound here, or a failed fetch. It never switches the branch you have checked out, and never merges, rebases, resets, stashes, or discards. `workspace sync` does the same for the workspace's own repository, but only when it has no uncommitted file at all; once it moves, it says to run `check` and names a pinned CLI version that moved with it.
