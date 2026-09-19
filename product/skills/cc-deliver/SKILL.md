---
name: cc-deliver
description: Deliver committed plan work as pull or merge requests against each repository's base branch, whether the work is one plan, one bound checkout, or a run of several plans.
---

# Deliver plan work

"Deliver the changes", "deliver all changes", and "open a PR" are requests to
push committed work and open a pull or merge request for it. This skill says
which branch each request opens from, and what it targets.

It carries no authorization. The request authorizes push and request creation;
**merging is separate and still needs a person**, as do deployment and
publication. Delivering does not mark a plan completed.

## Which branches to open

Delivery is per repository. A plan has one worktree and one branch in each
repository it names, so a plan can produce more than one request.

**Work in a bound checkout.** Direct execution commits to the branch that
checkout is already on. Deliver that branch; there is no `cc/*` branch involved.

**One plan.** Deliver `cc/<plan>/<repo>` for each repository the plan names.

**Several plans run together, or "deliver all changes".** Ask the CLI rather
than reading the dependency graph yourself:

```sh
context-circuit-cli --workspace <root> --json record order --intent i001
```

Each plan carries `deliver`: the repositories whose chain ends at that plan.
Open one request for every mark, from `cc/<plan>/<repo>`. A chain end's branch
already contains the plans it was prepared from, so those plans need no request
of their own — delivering them too would duplicate commits.

Do not deliver the last wave. Waves answer readiness, not chain ends: a plan
nothing depends on sits in wave 1 and still ends its chain, so delivering by wave
drops it silently.

Skip a branch holding no commits beyond its base. A plan can name a repository it
turned out not to change, and an empty request wastes a review.

## Why the marks are per repository

A predecessor that does not share a repository contributes nothing to that
repository's branch, so one plan can end a chain in one repository and sit
mid-chain in another:

```
p0003  repositories: web, sdk   depends_on: p0001   deliver: sdk
p0006  repositories: web        depends_on: p0003   deliver: web
```

`p0003` is delivered for `sdk`, where nothing follows it, and not for `web`,
where `cc/p0006/web` already carries it. One repository can also hold two chain
ends and take two requests. Deliver exactly the marks, no more and no fewer.

## Where to target

The base branch is this machine's recorded base for that repository, falling back
to the repository's default branch when the binding records none:

```sh
context-circuit-cli --workspace <root> --json worktree inspect --repo api --path <worktree>
context-circuit-cli --workspace <root> --json repo inspect --id api
```

Read `base_branch` from the worktree you are delivering, or from the repository
when the work is in its bound checkout. `workspace.yaml` records a
`default_branch`. That is the repository's default and not what this machine
delivers to; answering from it targets the wrong branch convincingly. Use
`base_branch`, or a branch the person names.

## Opening the request

Use ordinary Git and provider tools. The CLI never delivers.

```sh
git -C <worktree> push --set-upstream origin cc/p0001/api
gh pr create --base <base_branch> --head cc/p0001/api --title "..." --body "..."
```

Confirm the work is committed and pushed first. Uncommitted work is not delivered
by a request, and a worker that reported without committing left nothing to open.

Report each request with its URL, the branch it opened from, the branch it
targets, and the plans it carries. Where two chain ends share a predecessor, that
predecessor's commits appear in both requests and merging either shrinks the
other — say so rather than letting a reviewer meet the same commits twice.

Report drift or conflicts against the base plainly. Do not rebase, force-push, or
re-verify automatically to clear them.

## If the person wants one request per plan

Within a single run of plans that depend on each other and are delivered
together, some teams review one request per plan, each targeting its
predecessor's branch so it shows that plan's diff alone. Do that when asked, and
say the requests must merge in order.

This applies only inside one such run. A plan from a separate intent, or one
whose predecessor has already merged, has no predecessor branch to target: it is
delivered against the base branch like any other. Most delivery is this case —
one plan, finished and merged, then the next.
