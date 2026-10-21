# Worktree responsibilities

Executing a plan prepares a worktree for each repository that plan names. The
one exception is a person asking to work directly in a bound checkout, which is
how a change gets made in place when isolation is not wanted. Where the work
happens is therefore the person's choice, not a strategy the agent selects per
task. The CLI performs the corresponding Git mechanics.

| Concern | CLI | Agent |
| --- | --- | --- |
| Repository selection | Resolve logical ID and local checkout | Decide which repositories are affected |
| Existing work | Report branch, changes, and Git worktree inventory | Preserve unrelated work; resume an existing location or select a new one |
| Starting point | Resolve selected ref to a commit | Choose this machine's base or a dependency branch; fetch if needed |
| Branch and path | Validate names and collisions; create/reuse using Git | Honor explicit choices and repository conventions |
| Result | Return path, branch, HEAD, the start used and its commit, base branch, association | Use the returned directory for work |
| Environment setup | Supply location and Git facts | Follow setup instructions; install dependencies and prepare services |
| Recovery | List, inspect, move, or repair through Git; prune associations of worktrees Git no longer lists | Compare plan notes to real changes and decide next steps |
| Cleanup | Remove a selected worktree only when invoked | Obtain or reuse authorization; preserve work; handle branch deletion separately |

For a plan, the default branch is `cc/<plan-id>/<repository-id>` and the default
path is `.worktrees/<plan-id>/<repository-id>` under the workspace. Explicit
`--branch`, `--start`, and `--path` values override these choices. Preparation
without a plan is supported with an explicit branch and path.

An unspecified start uses the base branch this machine records for the repository
in `repositories.local.yaml`, falling back to the shared `default_branch` when the
binding names none, and then to that branch's `origin` tracking branch when the
checkout has no local copy of it. No fetch happens implicitly. An explicit start may
be another plan's branch or a commit. It is resolved before creating the worktree.
Preparation never resets an existing branch. Use `--reuse` when explicitly
selecting an existing branch or adopting another worktree association. Repeating
preparation for the same existing association returns its actual Git state.
`--start` applies to new branches; reused branches retain their existing commits.

The result's `start` names what the worktree actually started from: the
`--start` given, the base branch or `origin/<base>` it fell back to, or the
existing branch it checked out. `base_branch` is this machine's base, where the
work is delivered, and does not change with `--start`. A local start branch that
is behind its fetched remote-tracking branch, or has diverged from it, is used as
asked and named in `sync_required` with the counts and the remedies: prepare
again with `--start origin/<branch>`, or bring the local branch up to date first.
The remote ref is only as recent as the last `repo fetch`.

A plan branch already checked out in another working copy is refused, because Git
checks a branch out in one place at a time. The refusal names the branch, where it
is checked out, and whether that is the worktree this workspace prepared for the
plan, the bound checkout, or a worktree it does not know for the plan. It lists
the safe next steps: resume in that worktree, switch that checkout to another
branch, or remove that worktree, then prepare again. The bound checkout is never
offered for removal. Nothing is checked out, switched, or reset to get past it.

Git's inventory is authoritative. Local plan associations help locate work but
do not prove that a path exists or still holds the same branch. Resume with
`worktree list` and `worktree inspect`; then inspect the real diff. `worktree repair`
can reconnect a manually moved working copy; `worktree move` moves it through Git.
No automatic cleanup or pruning occurs after failure or completion. A worktree removed
with `git worktree remove` or by hand leaves its association behind, and `check` reports
it; `worktree prune --repo ID` drops the associations of that repository whose path Git
no longer lists. It rewrites only the local association file, so no worktree, branch, or
file is touched.

Run prune after removing a worktree with `git worktree remove` or deleting it by
hand, when `check` reports a stale association, or when `record delete` is blocked
by one. `worktree repair` cannot fix these, because it works only on a path Git
still lists. The result's `dropped` list carries each association it dropped. An
association whose worktree Git still lists is left alone, even when Git marks that
worktree prunable or it now holds another branch; use `worktree repair` or
`git worktree prune` for those.

Removing a worktree leaves its branch. `record delete` warns about local and
remote-tracking `cc/pNNNN/*` branches, worktrees on such a branch, and a leftover
`.worktrees/pNNNN` directory for a deleted plan number, and touches none of them:
a plan later given that number meets the old branch, and `--reuse` would start
it from the old commits, so delete or rename the branch first if that work is
done with. Another machine sees none of this until it pulls the deletion; there,
`check` names a worktree prepared for a plan that is gone, which is removed with
`worktree remove` (its branch deleted once the work no longer matters) before
`worktree prune` and before creating plans.

Preparation never resumes another plan's work silently. The association records
the plan's `created_at`, so a worktree or branch left from an earlier plan that
held the same number is refused without `--reuse`, and the refusal says how to
remove it; an association written before 2.4.0 is judged by its commits instead,
and is refused when any was made before the plan was created. One whose branch
holds no commit past its start has nothing to judge, so it is resumed only when
it has no uncommitted changes and sits at the tip of the base branch, exactly
where a fresh prepare would start; otherwise it is refused the same way.
Preparing again for the plan a worktree was made for resumes it as before.
`check` names a worktree whose association records a different `created_at` from
its plan's, which is removed with its branch, or pruned if Git no longer lists
it.

Successful preparation means a usable Git working copy. The agent still needs to
follow repository setup instructions for package managers, toolchains,
submodules, generated files, databases, and development services. The CLI does
not run arbitrary setup hooks. It reuses selected ignored runtime files as
described below; unrelated credentials remain in host facilities.

For changes across repositories, prepare each needed working copy independently.
The agent orders tasks and checks combined behavior; a failure in one repository
does not discard work in another. There is no cross-repository transaction.

`worktree remove` preserves dirty, untracked, and ignored files unless the caller
explicitly uses `--discard` with authorization to dispose of them. It never removes
the primary checkout or a locked worktree. A successful removal keeps the branch;
branch deletion is a separate ordinary Git action. It does remove the folders
under the workspace's `.worktrees/` that the removal left empty, such as
`.worktrees/pNNNN/` once a plan's last worktree is gone, walking up to and
including `.worktrees/` and stopping at the first folder that holds anything. It
removes only empty folders, never anything outside `.worktrees/`, so a worktree
at a custom path elsewhere leaves its parent alone. The result's
`removed_directories` lists each folder removed; one that could not be removed is
listed in `warnings` and does not fail the command, because the worktree is
already gone. Marking a plan done removes
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
