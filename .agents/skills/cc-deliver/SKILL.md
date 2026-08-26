---
name: cc-deliver
description: Handle separate delivery actions (pull request, merge, push) with explicit targets; never implied by verification or completion.
---

Merge, push, pull-request creation, publication, and deployment are separate
human-requested actions. None is implied by worker success, verifier success, or
plan completion. The runtime never performs them and never interprets
verification as merge authorization.

## Pull request

On an explicit request to open a pull request for an implemented plan:

- the source is each repository's execution branch `cc/<plan-id>/<repository-id>`;
- the default target is that repository's recorded `anchor_branch`;
- an alternative target must be named explicitly;
- never substitute `default_branch` and never silently follow a moving or
  renamed remote branch;
- a multi-repository plan may produce one pull request per affected repository.

Before opening, confirm the execution branch and its commits are available to
the configured repository remote or provider and that the requested target
branch exists. Do not discover an unrelated remote, silently push, or silently
substitute a target. If the source branch is unpublished, the provider/remote is
unavailable, or the recorded anchor branch is missing or renamed, report the
delivery as blocked and ask for an explicit human decision.

## Drift guard (v0.6)

Before opening a pull request, check `delivery-drift`: when a sibling plan has
already merged and advanced the recorded `anchor_branch`, the plan's base has
diverged from the tip it will land on. If drift is detected, run `delivery-rebase`
to rebase the execution branch onto the current anchor tip, then re-verify (one
fresh independent verifier pass) before the pull request opens — a plan is never
merged from a base that no longer reflects its target branch. A rebase conflict
(`DELIVERY_REBASE_CONFLICT`) is reported as blocked, with the work preserved, for
an explicit human decision (INV-DELIVER-01).

## Merge / push / cleanup

Report the effect and the branches involved; take the action only on an explicit
request. Cleanup must inspect dirty or unpushed work and preserve it unless the
human explicitly requests its removal. A failed execution is never cleaned up as
a side effect of reporting failure.
