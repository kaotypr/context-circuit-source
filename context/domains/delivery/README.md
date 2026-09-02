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
  - docs/getting-started.md
---

# Delivery

## Summary

Opening a pull request, merging, or pushing — each a separate, explicit human
action never implied by a prior success. Route "open a pull request for `<id>`",
merge, and push requests here. Owned by the `cc-deliver` skill; the delivery
boundary is owned by `wrapper/adapters/WORKFLOW.md`.

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
completion. The engine's delivery function is report-only: it produces the
per-repository pull-request source and default target and never pushes, merges,
or opens pull requests itself.

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

## Workflows

- Open a pull request as a separate step: `docs/getting-started.md`

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
- `wrapper/runtime/engine.sh`: `cc_delivery_targets` (read-only report),
`cc_delivery_drift`, `cc_delivery_rebase` (delivery drift guard)
- `.agents/skills/cc-deliver/SKILL.md` (Drift guard section)
- `wrapper/adapters/WORKFLOW.md` (delivery-boundary owner per `invariants.yaml`)
- `wrapper/contracts/invariants.yaml`: INV-DELIVER-01 (with the drift-guard
  clause), INV-DELIVER-02

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
INV-DELIVER-01 and AC-16).
