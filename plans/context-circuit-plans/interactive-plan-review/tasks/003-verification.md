---
schema_version: 2
id: IPR-003
plan: interactive-plan-review
status: draft
repository: context-circuit-source
paths:
  - test/routing/fixtures.yaml
  - test/routing/test-router.sh
  - test/lifecycle/test-lifecycle.sh
  - test/hosts/test-host-adapters.sh
  - test/release/test-release.sh
  - test/acceptance.sh
depends_on: [IPR-001, IPR-002]
acceptance: [IPR-AC-01, IPR-AC-02, IPR-AC-03, IPR-AC-04, IPR-AC-05, IPR-AC-06]
verification: [IPR-VT-01, IPR-VT-02, IPR-VT-03, IPR-VT-04, IPR-VT-05]
expected_evidence:
  - Routing fixtures covering named review, unnamed review, and nearby phrasing.
  - Read-only review and Review Card coverage in the lifecycle suite.
  - Host-adapter coverage that cc-plan remains the discovery adapter and question prompts are not a second router.
  - Release coverage that the seven shipped skills remain the allowlist and cc-review-plan is absent.
  - Complete semantic acceptance output and an independent verifier handoff.
stop_conditions:
  - Tests pass while a new skill directory is present in the artifact.
  - Offline CI requires a live host question UI.
  - Verification relies on the implementing session's claim without reproducing the bounded evidence.
---

# Lock routing, release, and read-only review behavior in tests

## Objective

Turn named-plan review, missing-target clarification, `cc-plan` discovery, and
optional question-prompt fallback into durable acceptance evidence.

## Work

Add or refine routing fixtures for named review, unnamed review, and nearby
phrasing so they keep `review-plan` or `clarify-target`. Keep review
read-only in the lifecycle suite and keep the Review Card’s risks and human
decisions. Assert in host and release coverage that `cc-plan` remains the
plan discovery adapter, question prompts do not become authorization, and
`cc-review-plan` is not a shipped skill. Re-run complete semantic acceptance.
Produce an independent verifier handoff.

Live host question UIs are optional and must not be required by offline CI.

## Non-goals

Do not mark the plan done, treat a question-prompt answer as a gate, or
expand the shipped skill allowlist.

## Verification

Use IPR-VT-01 through IPR-VT-05.

## Expected evidence

Fixture and suite output, unchanged seven-name allowlist, complete acceptance
output, and an independent handoff.

## Stop conditions

Stop on any regression in review authorization, skill allowlist, or
independent reproducibility.
