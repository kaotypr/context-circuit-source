---
name: cc-deliver
description: Handle separate delivery actions (pull request, merge, push) with explicit targets; never implied by verification or completion.
---

Merge, push, pull-request creation, and deployment are separate
human-requested actions. None is implied by worker success, verifier success, or
plan completion. The runtime never performs them and never interprets
verification as merge authorization. Delivery is **Gate 2** — the second and final
human gate (INV-DELIVER-01); it is never implied by the acceptance of a candidate.

## Change set — one pull request, one candidate

When several stacked plans converge to a single pull request, they form one
**change set**, verified once against an integration tip (not once per plan):

1. `change-set-prepare . <plan> <plan> ...` builds the integration tip — every
   member branch merged onto the anchor, per repository — and records one change-set
   candidate. If the combined result will not build it reports `BASE_UNBUILDABLE`;
   report that as blocked, not a worker failure, for the human to split or reorder.
2. Spawn ONE independent verifier over the integration tip and record it with
   `change-set-verifier-record . <change-set-id> <passed|failed|blocked|waived>`
   (read-only; a tip that moved since prepare voids it). At Explore there is none.
3. The human accepts once: `change-set-accept . <change-set-id> <who>`.
4. `change-set-ready . <change-set-id>` enforces the tier floor for the whole set
   (the max tier across members): a candidate-bound acceptance, plus a candidate-bound
   independent pass at Standard/Critical. A member that moves after prepare makes the
   candidate stale and re-gates.
5. After you open the one pull request (Gate 2), `change-set-complete . <change-set-id>`
   marks **every member** done and emits each member's reconciliation-debt marker, from
   the single change-set acceptance — accept once, and the whole set completes.

A single plan delivered alone is a change set of one, identical to
`change-set-candidate . <plan>` and `candidate-current . <plan>`; it may complete
through this same path, or the single-plan path below.

## Single-plan delivery and inferred completion

Acceptance is keyed to the candidate: a single plan uses its own execution candidate;
a change set uses the integration candidate (above). For a single plan delivered on
its own, after the human authorizes and you open the pull request (Gate 2), record the
delivery with `delivery-record . <plan-id>` — the delivery signal, bound to the
current candidate; it performs no git action itself. At Explore/Standard, this lets
completion be **inferred** from candidate acceptance + delivery
(`completion-infer . <plan-id>`), rather than a manual "mark done"; at Critical,
completion stays an explicit human act (`plan-complete`). A post-delivery change
yields a new candidate and re-gates rather than completing stale work.

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

## Direct-collaboration branch

On an explicit request to open a pull request for a finished pairing session,
invoke `pair-delivery-targets . <session>`. The reported pairing branch is the
source and the connected repository's recorded `anchor_branch` is the target.
Describe the work as human-supervised and not independently verified.

If the session is still active, its worktree is dirty, or the anchor tip has
moved outside the pairing branch, delivery is blocked. Do not run the plan drift
rebase path: start a new direct-collaboration session from the current anchor and
have the worker bring the change forward under live human supervision. Opening
the pull request, pushing its source branch, and cleanup remain separate explicit
actions (INV-PAIR-01).

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
