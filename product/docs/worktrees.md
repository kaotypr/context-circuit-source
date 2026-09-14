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
generated files, databases, and development services. The executable does not run
arbitrary setup hooks or copy ignored files from another working copy. Credentials
stay with the host's normal environment and authentication tools.

For changes across repositories, prepare each needed working copy independently.
The agent orders tasks and checks combined behavior; a failure in one repository
does not discard work in another. There is no cross-repository transaction.

`worktree remove` preserves dirty, untracked, and ignored files unless the caller
explicitly uses `--discard` with authorization to dispose of them. It never removes
the primary checkout or a locked worktree. A successful removal keeps the branch;
branch deletion is a separate ordinary Git action. Marking a plan done removes
neither worktrees nor branches. Base drift does not trigger automatic rebasing,
review, or repair loops.
