---
kind: domain
status: accepted
title: Verification
slug: verification
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
  - The verifier is a separate actor from the worker, read-only for product files.
unknowns: []
contradictions: []
acceptance:
  state: accepted
  accepted_at: 2026-08-24
  accepted_by: maintainer
workflows:
  - .context-circuit/docs/getting-started.md
---

# Verification

## Summary

An independent, read-only check of the latest worker commits. Verification is
the only authority for "verified"; the worker never self-verifies. Route the
verify step and its outcome-reporting here. Owned by the `cc-verify` skill and
the `.context-circuit/agents/verifier.md` role.

## Scope

Inside: independent verification of the latest commit of every affected
repository, the outcome vocabulary, worker-claim-is-not-evidence, write and
tip-change rejection, and the host-blocked fallback.

Outside: repairing the implementation, changing plan status, and the worker loop
([plan-execution](../plan-execution/README.md)).

## Behavior

Independent verification applies at **Standard and Critical** only
(INV-ASSURE-01): at those tiers one independent verifier checks the current
**candidate** of every affected repository, strictly read-only with respect to
product files, and its result **binds to that candidate** (INV-CANDIDATE-01) — a
new commit or a criteria change yields a new candidate and voids the prior pass. At
**Explore** there is no independent verifier and the result is never labeled
"verified" (see [assurance](../assurance/README.md)). The verifier never repairs,
modifies product files, changes plan status, or self-verifies (INV-VERIFY-01). If
the host cannot create an independent verifier with read-only capability, the
execution is blocked; the worker or coordinator must not self-verify as a
substitute (INV-VERIFY-02).

The verifier receives the worker handoff as a claim, not as evidence. Its
outcome is `passed`, `failed`, `blocked`, or `waived`; only `passed` satisfies
verification. Recording rejects any product write (`VERIFIER_WRITE_REJECTED`) or
changed branch tip (`VERIFIER_MODIFIED_PRODUCT`). Each verifier rejection —
including the initial implementation — increments the worker-failure counter
toward the three-failure limit (INV-REPAIR-01). A verifier must not downgrade an
evidence requirement because a host lacks a capability.

## Workflows

- Execute → verify → repair loop: `.context-circuit/docs/getting-started.md`

## Interfaces

- Outcome tokens: `passed`, `failed`, `blocked`, `waived`; user-facing phrasing
  for an environment limitation: `host-blocked`
- Record: `verifier.yaml` per attempt (`read_only: true`)

## Constraints and edge cases

`blocked`/`waived` are non-passing and never a false success. The verifier writes
only its own result and handoff, never product files.

## Verifier report and read-only capability

The verifier runs each canonical verification command plus the bounded additional
checks needed to establish the declared acceptance. For each acceptance and
verification id it reports the observed evidence, the passed/failed/blocked
outcome, the exact failure and a repair recommendation, and the repository and
commit checked.

The host adapter must provide actual read-only capability for the verifier child;
if it cannot guarantee it, the adapter reports blocked. A prompt saying "do not
edit" is not enough.

Not represented here: the design requires an evidence layer per acceptance
criterion, but the shipped verifier does not implement it. That divergence is
logged in `context/DESIGN-DELTAS.md`; this page describes the shipped verifier.

## Implementation references

- `.agents/skills/cc-verify/SKILL.md`, `.context-circuit/agents/verifier.md`
- `.context-circuit/wrapper/runtime/engine.sh`: `cc_verifier_prepare`, `cc_verifier_result_record`,
  `cc_repair_allowed`
- `.context-circuit/wrapper/contracts/schemas/verifier-result.yaml`
- `.context-circuit/wrapper/contracts/invariants.yaml`: INV-VERIFY-01, INV-VERIFY-02, INV-REPAIR-01

## Verification

`sh test/acceptance.sh` (execution/verification coverage).

## Acceptance notes

Accepted 2026-08-24. Re-grounded 2026-09-04 for Context Circuit v1.0: the
verifier floor is tier-conditional (Standard/Critical only; Explore has none)
and results bind to the candidate (INV-ASSURE-01, INV-CANDIDATE-01). The
verifier mechanism itself is unchanged.
