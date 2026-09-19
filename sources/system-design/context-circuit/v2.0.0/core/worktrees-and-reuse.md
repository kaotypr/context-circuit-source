# Worktrees and copy-on-write reuse

Isolation that costs little enough to actually use.

## Optional, and chosen by the agent

Worktrees are not mandatory. The agent picks a working strategy from repository
instructions, existing work, task needs, and stated preference, recommending
isolation for risky or parallel work. The executable then performs the Git
mechanics. Routine preparation covered by the user's request does not get an
extra permission gate.

| Concern | Executable | Agent |
| --- | --- | --- |
| Repository selection | Resolve logical ID to a local checkout | Decide which repositories are affected |
| Existing work | Report branch, changes, worktree inventory | Preserve unrelated work; reuse or isolate |
| Starting point | Resolve the selected ref to a commit | Choose base or dependency branch; fetch if needed |
| Branch and path | Validate names and collisions; create or reuse | Honor explicit choices and conventions |
| Result | Return path, branch, HEAD, start commit, association | Use the returned directory for all later work |
| Environment | Reuse selected ignored runtime entries | Finish setup the report left open |
| Recovery | List, inspect, move, repair through Git | Compare notes to real changes; decide |
| Cleanup | Remove one worktree when invoked | Hold authorization; handle branch deletion separately |

## Defaults and overrides

For a plan: branch `cc/<plan-id>/<repository-id>`, path
`.worktrees/<plan-id>/<repository-id>` under the workspace. Explicit `--branch`,
`--start`, and `--path` override. Preparation without a plan works with an
explicit branch and path.

An unspecified start uses the recorded local base branch if available, otherwise
its `origin` tracking branch. The base is **this machine's** binding, recorded in
`repositories.local.yaml` since 2.0.0-rc.7, so one machine can work off a release
branch while another stays on the default. **No fetch happens implicitly**, and
fetching implies neither rebase nor reset. An explicit start may be another
plan's branch or a commit, resolved before the worktree is created. `--start`
applies to new branches; a reused branch keeps its commits.

**Preparation derives no order.** A plan that records a dependency in the same
repository is refused without `--start`, rather than silently branching from a
base that does not contain its predecessor — the failure 2.0.0-rc.9 closed. The
refusal names the predecessor branch to pass, or says which predecessor has not
been implemented and committed yet. An explicit start is always honored,
including one that is not the predecessor.

## The preservation rules

These are the rules that make isolation safe to automate, and every one of them
is a refusal:

- **Preparation never resets an existing branch.** Never force-checkout, reset,
  stash, or overwrite unrelated files.
- **A branch or path already holding work is inspected**, then resumed or
  replaced with a new location — never silently taken over. `--reuse` is how the
  caller explicitly selects an existing branch or adopts an existing association.
- **Repeating preparation for the same association returns its actual Git
  state**, rather than recreating anything.
- **Removal preserves dirty, untracked, and ignored files** unless `--discard`
  carries explicit authorization to dispose of them. `--discard` is never added
  merely to get past a failure.
- **Removal never touches the primary checkout or a locked worktree**, and a
  successful removal keeps the branch. Branch deletion is separate ordinary Git.
- **Marking a plan done removes neither worktree nor branch.** No automatic
  cleanup or pruning follows failure or completion.
- **Base drift triggers nothing.** No automatic rebase, review, or repair loop.

## Git is the inventory

Local plan-to-worktree associations in `.context-circuit/local/worktrees.yaml`
are a convenience for resuming work. They do not prove a path exists or still
holds the same branch. Resume with `worktree list` and `worktree inspect`, then
inspect the real diff (P1). `worktree repair` reconnects a manually moved working
copy; `worktree move` moves one through Git. `check` reports stale associations
as diagnostics.

## Copy-on-write environment reuse

### The problem it solves

The most common reason people avoid worktrees is that a fresh working copy means
a fresh `node_modules` and a missing `.env`. The cost is paid every time
isolation would have helped, which means isolation is used when it is least
convenient and skipped when it matters.

### The mechanism

After Git creates or selects the worktree, preparation discovers ignored
`node_modules`, `.env`, and `.env.*` entries in the **bound local checkout** and
clones them using the filesystem's own cloning:

| Platform | Mechanism |
| --- | --- |
| macOS | `clonefile` |
| Linux | `FICLONE` |
| Windows | Aligned ReFS block cloning |

`--copy-mode` selects the policy: `auto` (default) attempts CoW and falls back to
independent copies; `required` fails if a regular file cannot be cloned; `copy`
forces copies; `off` disables reuse entirely. Windows files without complete
cloneable clusters take the copy fallback.

Repeatable `--copy-path PATH` selects additional ignored runtime entries. Only
repository-relative paths ignored in **both** source and target are eligible, and
discovery does not descend into arbitrary ignored parent directories — an entry
hidden by its parent must be selected explicitly.

### The safety properties

- **Staged and published only when complete.** A failure preserves the Git
  worktree and previously completed entries; rerunning resumes.
- **Existing destination files and directories are preserved**, including changes
  made during previous work. Reuse never overwrites.
- **Internal symlinks are translated** to the new checkout; external symlinks and
  special files are rejected. Hard links to the source are never created.
- **Control files are unreachable.** A reuse path containing `.git` or
  `.context-circuit`, or any non-repository-relative path, is rejected.
- **Environment contents are opaque.** `.env` files are copied only when Git
  ignores the target, never overwrite an existing worktree environment, and are
  never printed, put in a prompt, or stored in a shared record. The report names
  entries and counts cloned/copied files.

### The dependency-input check

Before reusing `node_modules`, the executable compares tracked
`package.json`, `package-lock.json`, `npm-shrinkwrap.json`, `yarn.lock`,
`pnpm-lock.yaml`, `pnpm-workspace.yaml`, `bun.lock`, `bun.lockb`, `.nvmrc`, and
`.node-version` between source and target. Any difference — including a file
present on one side only — skips dependencies with the reason *"dependency
inputs differ; prepare dependencies for this branch."*

### What reuse does not prove

The honest limits, stated in the shipped documentation rather than discovered
later:

- Matching inputs do not prove the source installation is **complete**.
- They do not prove its runtime or Node ABI is **compatible** with the target.
- The tree is **not an atomic snapshot**; a package install running during the
  copy produces an inconsistent result, so source installs should be stopped.
- Local files do **not** make native dependencies portable across OS,
  architecture, Node ABI, or container. A new environment needs its own.
- CoW does not start services, migrate databases, or make local ports unique.

The correct posture is the one the docs state: reuse a working local environment,
read the reuse report before application setup, and run the project's normal
checks. Reused dependencies need no reinstall *solely because the worktree is
new* — but a skipped entry is a reported fact the agent must act on.

## Across repositories

For a change spanning repositories, each working copy is prepared independently.
The agent orders the tasks and checks combined behavior. There is no
cross-repository transaction, and a failure in one repository does not discard
work in another.
