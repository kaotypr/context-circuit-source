# 0009 — Delivery boundary

- **Plan ID:** `0009-delivery`
- **Status:** done
- **Repository:** `context-circuit-source`
- **Depends on:** `0005-approval-and-execution`
- **Owns (invariants):** INV-DELIVER-01, INV-DELIVER-02

## Original request

Retroactive plan for the delivery increment, built as if from an empty repo.
Source design: `sources/system-design/context-circuit/v0.5/core/03-workspace-and-repositories.md`
(delivery) and `06-runtime-engine.md`.

## Objective and desired behavior

- Pull-request creation, merge, push, deployment, archive, and cleanup are
  separate human-requested actions.
- A pull request uses each execution branch as source and the repository's
  recorded `anchor_branch` as the default target; it never substitutes
  `default_branch` or silently follows a moving remote.
- Delivery blocks and reports when the source branch, configured provider or
  remote, or target branch is unavailable, rather than inferring a remote.

## Constraints and non-goals

- Non-goal: the v0.6 delivery drift guard (added by 0014). This plan is the
  single-plan delivery boundary.

## Product Knowledge grounding

- `invariants` (`wrapper/contracts/invariants.yaml`), `decisions`
  (`context/DECISIONS.md`).

## Tasks

1. **DEL-001** — delivery boundary + anchor-targeted PR + block-not-infer.
2. **DEL-002** — `cc-deliver` skill.

## Acceptance & verification

- PR targets anchor; unavailable source/provider/target blocks; each verb is a
  distinct explicit action.
- `sh test/delivery/test-delivery.sh`.

## Assumptions, open questions, risks

- Risk: silently pushing to a moving remote — prevented by INV-DELIVER-02.

## Expected commits and delivery notes

This plan defines the delivery boundary itself; it performs no delivery.
