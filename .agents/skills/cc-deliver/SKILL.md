---
name: cc-deliver
description: Handle separate delivery actions (pull request, merge, push) with explicit targets; never implied by verification or completion.
---

Merge, push, pull-request creation, and deployment are separate
human-requested actions. None is implied by worker success, verifier success, or
plan completion. The runtime never performs them and never interprets
verification as merge authorization. Delivery is **Gate 2** — the second and final
human gate (INV-DELIVER-01); it is never implied by the acceptance of a candidate.

## One pull request per covering tip

A pull request is per covering tip, not per plan and not per repository. When
delivering several plans, `change-set-partition . <plan> ...` groups them by
covering execution branch.

- Same-repository stacked plans that already nest (3 depends on 2 depends on 1)
  are one change set, one candidate, and one pull request from the covering
  branch (plan 3's `cc/<plan>/<repo>`).
- Sibling stacks in the same repository are several covering tips. Plan 1 with
  dependents 2 and 3 in repo A is two pull requests (`cc/2/A` and `cc/3/A`),
  not one and not zero.
- Each other repository is its own covering tip (or tips). Four plans of which
  three stack in repo A and one lives in repo B are two candidates and two pull
  requests, not four and not one.

Do not pass members from different repositories to `change-set-prepare` or
`change-set-candidate` (`CHANGE_SET_CROSS_REPO`). Partition first, then prepare
one change set per ready group. Do not spawn a fresh combined verifier and do
not invent a combined candidate across repositories. Cross-repo dependencies
are ordering gates only (INV-CONCURRENCY-02).

## Change set — one pull request, one repository

When several stacked plans **in the same repository** converge to a single pull
request, they form one **change set**. The candidate is the live member tip map
(`change-set-candidate . <plan> <plan> ...`). Delivery does **not** spawn a
verifier and does **not** author a merge. Each member's already-bound
independent pass is the floor.

1. `change-set-prepare . <plan> <plan> ...` records that tip-map candidate and
   names the one **covering** execution branch (the tip that already contains
   every other member). If the members do not share a repository it reports
   `CHANGE_SET_CROSS_REPO` — partition and prepare per covering tip instead of
   falling back to one pull request per plan. If no single execution branch
   contains the others it reports `CHANGE_SET_NO_SINGLE_TIP`; partition already
   splits those sibling tips into separate pull requests. Do not merge them at
   delivery.
2. Do not call `change-set-verifier-prepare` or `change-set-verifier-record`
   (`CHANGE_SET_DELIVERY_HAS_NO_VERIFIER`). Change sets are for plan-bearing
   Standard/Critical work; planless Explore work uses `cc-pair`.
3. The human accepts once: `change-set-accept . <change-set-id> <who>`.
4. `change-set-ready . <change-set-id>` enforces the tier floor: the tip map is
   unchanged, a candidate-bound acceptance, and each member still has its own
   candidate-bound independent pass. A member that moves makes the candidate
   stale and re-gates.
5. Open the one pull request from the reported covering execution branch
   (Gate 2), then `change-set-complete . <change-set-id>` marks **every member**
   done and emits one change-set-bound reconciliation-debt marker — accept once,
   and the whole set completes.

A single Standard/Critical plan delivered alone is a change set of one, identical
to `change-set-candidate . <plan>` and `candidate-current . <plan>`; it may complete
through this same path, or the single-plan path below. Explore work is planless
and has no delivery record in the plan lifecycle.

## Single-plan delivery and inferred completion

Acceptance is keyed to the candidate: a single plan uses its own execution candidate;
a change set uses the member tip-map candidate (above). For a single plan delivered on
its own, after the human authorizes and you open the pull request (Gate 2), record the
delivery with `delivery-record . <plan-id>` — the delivery signal, bound to the
current candidate; it performs no git action itself. At Standard, this lets
completion be **inferred** from candidate acceptance + delivery
(`completion-infer . <plan-id>`), rather than a manual "mark done"; at Critical,
completion stays an explicit human act (`plan-complete`). A post-delivery change
yields a new candidate and re-gates rather than completing stale work.

## Pull request

On an explicit request to open a pull request for an implemented plan:

- for a change set, the source is the covering execution branch reported by
  partition / prepare (`cc/<covering-plan>/<repository-id>`);
- for a single plan delivered alone, the source is that plan's execution
  branch `cc/<plan-id>/<repository-id>`;
- the default target is that repository's recorded `base_branch`;
- an alternative target must be named explicitly;
- never substitute `default_branch` and never silently follow a moving or
  renamed remote branch;
- one pull request per covering tip. Stacked same-repository plans share that
  one pull request; sibling stacks in one repository are several pull requests;
  two repositories are at least two pull requests.

Before opening, confirm the execution branch and its commits are available to
the configured repository remote or provider and that the requested target
branch exists. Do not discover an unrelated remote, silently push, or silently
substitute a target. If the source branch is unpublished, the provider/remote is
unavailable, or the recorded base branch is missing or renamed, report the
delivery as blocked and ask for an explicit human decision.

## Direct-collaboration branch

On an explicit request to open a pull request for a finished pairing session,
invoke `pair-delivery-targets . <session>`. The reported pairing branch is the
source and the connected repository's recorded `base_branch` is the target.
Describe the work as human-supervised and not independently verified.

If the session is still active, its worktree is dirty, or the base tip has
moved outside the pairing branch, delivery is blocked. Do not run the plan drift
rebase path: start a new direct-collaboration session from the current base branch and
have the worker bring the change forward under live human supervision. Opening
the pull request, pushing its source branch, and cleanup remain separate explicit
actions (INV-PAIR-01).

## Drift guard (v0.6)

Before opening a pull request, check `delivery-drift`: when a sibling plan has
already merged and advanced the recorded `base_branch`, the plan's base has
diverged from the tip it will land on. If drift is detected, run `delivery-rebase`
to rebase the execution branch onto the current base tip, then re-verify (one
fresh independent verifier pass) before the pull request opens — a plan is never
merged from a base that no longer reflects its target branch. A rebase conflict
(`DELIVERY_REBASE_CONFLICT`) is reported as blocked, with the work preserved, for
an explicit human decision (INV-DELIVER-01).

## Merge / push / cleanup

Report the effect and the branches involved; take the action only on an explicit
request. Cleanup must inspect dirty or unpushed work and preserve it unless the
human explicitly requests its removal. A failed execution is never cleaned up as
a side effect of reporting failure.
