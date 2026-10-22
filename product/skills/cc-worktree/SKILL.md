---
name: cc-worktree
description: Prepare a Context Circuit worktree for a plan's repositories, choose its starting point, and read the environment reuse report before application setup.
---

# Prepare a worktree

Preparing a worktree for each repository the plan names is the first step of
execution, not a judgment call, and the execution request already covers it. The
one exception is a person asking to work directly in a bound checkout: work
there, and preserve everything it already holds.

```sh
sh <root>/.context-circuit/cli worktree prepare --repo api --plan p0001
```

Preparing again for a plan whose worktree already exists resumes it without a
flag. `--reuse` takes an existing branch or worktree the CLI did not prepare for
this plan record, so pass it only after inspecting what a refusal names: a
deleted plan's number can come back for a new plan, and its old branch carries
the old plan's commits.

## Choosing the starting point

The CLI resolves repository bindings, reports Git state, resolves the selected
starting point, and creates or reuses worktrees. Select the recorded base branch
— reported as `base_branch` beside the checked-out branch — or an appropriate
dependency branch.

Preparation defaults to the recorded base branch and derives no order, so a plan
that depends on another needs its start named explicitly:

```sh
sh <root>/.context-circuit/cli worktree prepare --repo api --plan p0002 --start cc/p0001/api
```

It refuses rather than guessing, and the refusal names the branch to pass. Where
several predecessors touch the repository it names the base and the integration
merges the extras need; `record order` reports the same shape for the whole set.
A refusal saying a predecessor's work is not there yet means that plan has not
been implemented and committed — do that first. The predecessor's work is there
to start from because implementation commits before it is reported; never copy
files between worktrees to stand in for that ancestry.

Fetch only when needed and covered by the task; fetching does not imply rebasing
or resetting local work. Honor explicit branch names and paths. Use the returned
real worktree location for all later work, not the one you asked for.

The result's `start` is what the worktree actually started from, and
`base_branch` stays where its work is delivered. A `sync_required` there means
the start was behind its fetched remote: pass it on to the person, with its
remedy, rather than building on that start unremarked.

If a branch or path already holds work, inspect and resume it or select a new
location. Never silently force checkout, reset, stash, or overwrite unrelated
files.

## Environment reuse

Preparation attempts filesystem copy-on-write for ignored `node_modules` and
`.env` files from the selected local checkout, with an independent-copy
fallback. Additional ignored runtime paths can be selected explicitly with
`--copy-path`.

Reuse never overwrites existing worktree entries, and skips dependencies when
package inputs differ. Environment contents are copied opaquely: never print
them, put them in prompts, or store them in shared records. Do not copy
unrelated credential stores or host configuration.

## Before application setup

Inspect the reuse report first. Reused dependencies need no reinstall solely
because the worktree is new.

Read repository setup instructions for missing or incompatible dependencies,
toolchains, submodules, or services. Local files do not make native dependencies
portable across operating system, architecture, Node ABI, or container
environments. Report skipped entries and any remaining setup.

`.context-circuit/docs/worktrees.md` carries recovery, move, repair, and cleanup
mechanics.
