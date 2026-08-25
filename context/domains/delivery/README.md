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

Opening a pull request, merging, pushing, or publishing — each a separate,
explicit human action never implied by a prior success. Route "open a pull
request for `<id>`", merge, and push requests here. Owned by the `cc-deliver`
skill; the delivery boundary is owned by `wrapper/adapters/WORKFLOW.md`.

## Scope

Inside: the pull-request source/target model and the block-don't-infer behavior
for missing branches or remotes.

Outside: completion ([completion](../completion/README.md)) and the execution
branch creation itself ([plan-execution](../plan-execution/README.md)).

## Behavior

Pull-request creation, merge, push, publication, deployment, archive, and
cleanup are separate human-requested actions. A pull request uses each execution
branch `cc/<plan-id>/<repo-id>` as source and the repository's recorded
`anchor_branch` as the default target; it never substitutes `default_branch` or
silently follows a moving remote (INV-DELIVER-01). Delivery blocks and reports
when the source branch, configured provider or remote, or target branch is
unavailable, rather than inferring a remote or pushing silently (INV-DELIVER-02).

None of these is implied by worker success, verifier success, or plan
completion. The engine's delivery function is report-only: it produces the
per-repository pull-request source and default target and never pushes, merges,
or opens pull requests itself.

## Workflows

- Open a pull request as a separate step: `docs/getting-started.md`

## Interfaces

- Human request: "Open a pull request for `<id>`"
- Branch model: source `cc/<plan-id>/<repo-id>`, target = recorded `anchor_branch`

## Constraints and edge cases

An unpublished source branch, an unavailable provider/remote, or a
missing/renamed anchor branch blocks delivery and asks for an explicit human
decision. Cleanup preserves dirty or unpushed work unless removal is explicitly
requested; a failed execution is never cleaned up as a side effect.

## Implementation references

- `.agents/skills/cc-deliver/SKILL.md`
- `wrapper/runtime/engine.sh`: `cc_delivery_targets` (read-only report)
- `wrapper/adapters/WORKFLOW.md` (delivery-boundary owner per `invariants.yaml`)
- `wrapper/contracts/invariants.yaml`: INV-DELIVER-01, INV-DELIVER-02

## Verification

`sh test/acceptance.sh` (delivery suite).

## Provenance

Authored from the current wrapper at HEAD `4b8ac0b`. Raw `sources/` was not
scanned.

## Acceptance notes

Accepted 2026-08-24 from proposal `0015-domain-delivery`.
