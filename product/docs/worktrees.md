# Worktree responsibilities

Worktrees are optional. The agent chooses a working strategy from the user's
preferences, repository instructions, existing work, and change dependencies.
The executable performs the corresponding Git mechanics.

| Concern | Executable | Agent |
| --- | --- | --- |
| Repository selection | Resolve logical ID and local checkout | Decide which repositories are affected |
| Existing work | Report branch, changes, and Git worktree inventory | Preserve unrelated work and choose reuse or isolation |
| Starting point | Resolve selected ref to a commit | Choose default base or dependency branch; fetch if needed |
| Branch and path | Validate names and collisions; create/reuse using Git | Honor explicit choices and repository conventions |
| Result | Return path, branch, HEAD, starting commit, association | Use the returned directory for work |
| Environment setup | Supply location and Git facts | Follow setup instructions; install dependencies and prepare services |
| Recovery | List, inspect, move, or repair through Git | Compare plan notes to real changes and decide next steps |
| Cleanup | Remove a selected worktree only when invoked | Obtain or reuse authorization; preserve work; handle branch deletion separately |

For a plan, the default branch is `cc/<plan-id>/<repository-id>` and the default
path is `.worktrees/<plan-id>/<repository-id>` under the workspace. Explicit
`--branch`, `--start`, and `--path` values override these choices. Preparation
without a plan is supported with an explicit branch and path.

An unspecified start uses the recorded local base branch if available, otherwise
its `origin` tracking branch. No fetch happens implicitly. An explicit start may
be another plan's branch or a commit. It is resolved before creating the worktree.
Preparation never resets an existing branch. Use `--reuse` when explicitly
selecting an existing branch or adopting another worktree association. Repeating
preparation for the same existing association returns its actual Git state.
`--start` applies to new branches; reused branches retain their existing commits.

Git's inventory is authoritative. Local plan associations help locate work but
do not prove that a path exists or still holds the same branch. Resume with
`worktree list` and `worktree inspect`; then inspect the real diff. `worktree repair`
can reconnect a manually moved working copy; `worktree move` moves it through Git.
No automatic cleanup or pruning occurs after failure or completion.

Successful preparation means a usable Git working copy. The agent still needs to
follow repository setup instructions for package managers, toolchains, submodules,
generated files, databases, and development services. The executable does not run arbitrary setup hooks. It reuses selected ignored
runtime files as described below; unrelated credentials remain in host facilities.

For changes across repositories, prepare each needed working copy independently.
The agent orders tasks and checks combined behavior; a failure in one repository
does not discard work in another. There is no cross-repository transaction.

`worktree remove` preserves dirty, untracked, and ignored files unless the caller
explicitly uses `--discard` with authorization to dispose of them. It never removes
the primary checkout or a locked worktree. A successful removal keeps the branch;
branch deletion is a separate ordinary Git action. Marking a plan done removes
neither worktrees nor branches. Base drift does not trigger automatic rebasing,
review, or repair loops.

## CoW environment reuse

After Git creates or selects a worktree, preparation discovers ignored
`node_modules`, `.env`, and `.env.*` entries from the bound local checkout. It
attempts native CoW cloning: clonefile on macOS, FICLONE on Linux, and aligned ReFS
block cloning on Windows. Filesystem and volume capabilities determine success.
The default `--copy-mode auto` falls back to independent copies, avoiding a package
reinstall when source dependencies are usable. `required` fails if a regular file
cannot be cloned; `copy` forces copies; `off` disables environment reuse.
Windows files without complete cloneable clusters use the copy fallback.

Select additional ignored runtime entries with repeatable `--copy-path PATH`.
Only repository-relative paths ignored in both source and target are eligible.
Discovery does not descend into arbitrary ignored parent directories: select an
entry explicitly if its parent hides it from Git's ignored-file inventory.

Each entry is staged and published only when complete. Existing destination files
or directories are preserved, including changes made during previous work. Internal
symlinks are translated to the new checkout; external symlinks and special files
are rejected. Hard links to the source are never created. A failure preserves the
Git worktree and previously completed entries; rerun preparation to resume.
The report names entries and counts cloned/copied files without printing content.

Before reusing node_modules, compare tracked package manifests, lockfiles, workspace
package definitions, and Node version files against the target working copy. A
mismatch skips dependencies and reports that setup is needed. Matching inputs do
not prove that the source installation is complete or that its runtime/Node ABI
is compatible: reuse a working local environment and run normal project checks.
Stop source package installs while copying; the whole tree is not an atomic
snapshot. A new OS/container needs its own compatible dependencies.

Environment files are copied opaquely, only when Git ignores the target, and never
overwrite an existing worktree environment. This avoids routinely recreating local
settings while keeping their contents out of prompts, logs, and shared records.
CoW does not start services, migrate databases, or make local ports unique.
