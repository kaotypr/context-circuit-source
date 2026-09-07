---
kind: domain
status: accepted
title: Delivery
slug: delivery
owners: []
sources: []
source_revisions:
  - wrapper: HEAD
    commit: 4b8ac0b
    basis: current-wrapper
generated_at: 2026-08-24T00:00:00Z
review_date: 2026-11-24
freshness: accepted-from-current-wrapper
assumptions:
  - Delivery is never implied by worker success, verifier success, or completion.
unknowns: []
contradictions: []
acceptance:
  state: accepted
  accepted_at: 2026-08-24
  accepted_by: maintainer
workflows:
  - .context-circuit/docs/getting-started.md
---

# Delivery

## Summary

Opening a pull request, merging, or pushing — each a separate, explicit human
action never implied by a prior success. Delivery is **Gate 2**, the second (and
only other) human gate after intent approval: because v1.0 has **no automated
scope gate** upstream (the scope-envelope check was removed), delivery is where
**scope-safety is actually settled** — the human sees and authorizes the exact
diff and repositories before anything lands, while all pre-delivery work stays
sandboxed in isolated worktrees (INV-INTENT-02, INV-DELIVER-01). Route "open a pull
request for `<id>`", merge, and push requests here. Owned by the `cc-deliver`
skill; the delivery boundary is owned by `.context-circuit/wrapper/adapters/WORKFLOW.md`.

In this product "publish"/"publication" names sending data to an external system
([external-surface](../external-surface/README.md)); git delivery is "push" / "open
a pull request" and never "publish" — the two never share a word, so no qualifier is
needed.

## Scope

Inside: the pull-request source/target model and the block-don't-infer behavior
for missing branches or remotes.

Outside: completion ([completion](../completion/README.md)) and the execution
branch creation itself ([plan-execution](../plan-execution/README.md)).

## Behavior

Pull-request creation, merge, push, deployment, archive, and
cleanup are separate human-requested actions. A pull request uses each execution
branch `cc/<plan-id>/<repo-id>` as source and the repository's recorded
`base_branch` as the default target; it never substitutes `default_branch` or
silently follows a moving remote (INV-DELIVER-01). Delivery blocks and reports
when the source branch, configured provider or remote, or target branch is
unavailable, rather than inferring a remote or pushing silently (INV-DELIVER-02).

None of these is implied by worker success, verifier success, or plan
completion. Delivery does not mark a plan done and does not start Product
Knowledge reconcile (INV-COMPLETE-01, INV-COMPLETE-02). The engine's delivery
function is report-only: it produces the
per-repository pull-request source and default target and never pushes, merges,
or opens pull requests itself.

A pull request is per covering tip, not per plan and not per repository. Named
plans partition by covering execution branch: stacked dependents that already
nest are one change set and one pull request; sibling stacks in the same
repository are several pull requests (`cc/2/A` and `cc/3/A` when 2 and 3 both
depend on 1), not zero. Four plans of which three stack in A and one lives in
B are two candidates and two pull requests. Plans in different repositories
never form one change set and never require a fresh combined verifier
(INV-DELIVER-01). A change set is only same-repository stacked plans that
already share one covering execution branch; delivery records that tip map and
does not spawn a verifier.

**Drift guard (INV-DELIVER-01 extended).** When a plan is delivered and its
recorded base has diverged from the current `base_branch` tip (because a sibling
plan already merged), the plan is rebased onto the current tip and re-verified
before its pull request opens — a plan is never merged from a base that no longer
reflects the branch it will land on. `cc_delivery_drift` reports divergence
(read-only); `cc_delivery_rebase` rebases the execution branch onto the tip and
flags re-verification, and a rebase conflict is reported as blocked
(`DELIVERY_REBASE_CONFLICT`) with the work preserved. The only merge the runtime
authors anywhere is the integration *base* on a plan's own branch (run-stack,
INV-CONCURRENCY-02) — never a delivery merge.

**Pairing delivery blocks; it never rebases (INV-PAIR-01).** A closed
[direct-collaboration](../direct-collaboration/README.md) branch
(`cc-pair/<session>`, the Explore tier) may be delivered only by a separate
explicit action and is labeled **human-supervised, not independently verified** —
there is no verifier evidence to carry. Its target is the connected repository's
recorded `base_branch`. If the current base tip is not contained in the pairing
branch at delivery time, delivery **blocks**: the work must be brought forward in
a new human-supervised pairing session, not silently rebased. This is the
deliberate opposite of the plan drift guard above — the plan path rebases and
*re-verifies* because it has a verifier; a pairing branch has none, so
auto-rebasing would ship unreviewed drift. The plan delivery drift guard remains
the only path that rebases and re-verifies automatically.

## Workflows

- Open a pull request as a separate step: `.context-circuit/docs/getting-started.md`

## Interfaces

- Human request: "Open a pull request for `<id>`"
- Branch model: source `cc/<plan-id>/<repo-id>`, target = recorded `base_branch`

## Constraints and edge cases

An unpublished source branch, an unavailable provider/remote, or a
missing/renamed base branch blocks delivery and asks for an explicit human
decision. Cleanup preserves dirty or unpushed work unless removal is explicitly
requested; a failed execution is never cleaned up as a side effect.

## Implementation references

- `.agents/skills/cc-deliver/SKILL.md`
- `.context-circuit/wrapper/runtime/engine.sh`: `cc_delivery_targets` (read-only report),
`cc_delivery_drift`, `cc_delivery_rebase` (delivery drift guard),
`cc_change_set_partition` / `cc_change_set_prepare` / `cc_change_set_candidate`
(`CHANGE_SET_CROSS_REPO`, `CHANGE_SET_NO_SINGLE_TIP`,
`CHANGE_SET_DELIVERY_HAS_NO_VERIFIER`)
- `.agents/skills/cc-deliver/SKILL.md` (one pull request per covering tip; Drift guard section)
- `.context-circuit/wrapper/adapters/WORKFLOW.md` (delivery-boundary owner per `invariants.yaml`)
- `.context-circuit/wrapper/contracts/invariants.yaml`: INV-DELIVER-01 (same-repository change sets
  and the drift-guard clause), INV-DELIVER-02, INV-PAIR-01 (pairing delivery blocks on base drift)
- `.agents/skills/cc-pair/SKILL.md` (pairing delivery is a separate `cc-deliver`
  action, human-supervised)

## Verification

`sh test/acceptance.sh` (delivery suite).

## Provenance

Authored from the current wrapper at HEAD `4b8ac0b`. Raw `sources/` was not
scanned.

## Acceptance notes

Accepted 2026-08-24 from proposal `0015-domain-delivery`. Extended 2026-08-27 from
proposal `0019-change-delivery` (the v0.6 drift guard), and 2026-08-28 from proposal
`0024-change-delivery` (vacate "publish"/"publication" from git delivery; the word is
reserved for the [external surface](../external-surface/README.md), matching
INV-DELIVER-01 and AC-16). Extended 2026-09-03 from proposal
`0031-change-delivery-for-pairing`: pairing-branch delivery is a separate,
human-supervised action that blocks on base drift rather than rebasing (the plan
drift guard is unchanged and remains the only auto-rebase path). Re-grounded
2026-09-04 for Context Circuit v1.0: delivery is named as Gate 2, the point where
scope-safety is settled now that there is no automated scope gate upstream
(mechanics unchanged). Extended 2026-09-06: a change set is only same-repository
plans that will ship as one pull request; different repositories deliver
independently as separate covering-tip pull requests and never share an
integration candidate. Extended 2026-09-06:
delivery does not spawn a verifier; a same-repository change set is the member
tip map and the covering execution branch. Extended 2026-09-06: named plans
partition by repository covering tip, so a three-deep stack in one repository
plus one plan in another is two pull requests, not four. Sibling stacks in
one repository are several pull requests, not zero. Extended 2026-09-06:
delivery records Gate 2 only — it does not mark members done and does not
start in-place knowledge reconcile.
