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
  - scripts/release-artifact.sh
  - scripts/release-manifest.txt
  - test/acceptance.sh
depends_on: [IPR-001, IPR-002]
acceptance: [IPR-AC-01, IPR-AC-02, IPR-AC-05, IPR-AC-06]
verification: [IPR-VT-01, IPR-VT-02, IPR-VT-03, IPR-VT-04, IPR-VT-05]
expected_evidence:
  - Routing fixtures for named review, unnamed review, and non-stolen verification phrasing.
  - Lifecycle assertion that review remains read-only.
  - Host-adapter assertion that cc-review-plan is a discovery adapter.
  - Release allowlists accept cc-review-plan.
  - Complete semantic acceptance output.
stop_conditions:
  - Offline tests invoke a live host question UI or provider.
  - The new skill is treated as an unexpected legacy name.
  - Review fixtures gain mutation or gate authorization.
---

# Lock routing, release, and read-only review behavior in tests

## Objective

Prove that named-plan review is a specific read-only route, that the new
skill ships as an expected adapter, and that optional question prompts never
become a required host child or a hidden gate.

## Work

Add or refine routing fixtures for `Review plan <id>`, nearby walkthrough or
risk phrasing, and unnamed `Review the plan`. Assert named cases stay
`review-plan` with `authorization: read-only`, and unnamed cases stay
`clarify-target` or an equivalent focused question.

Keep a lifecycle contains assertion that the Review Card remains read-only
and still names risks and human decisions. Extend host-adapter checks so
`cc-review-plan` exists and cites read-only review without becoming a second
router. Update `scripts/release-manifest.txt`, `scripts/release-artifact.sh`,
and `test/release/test-release.sh` so `cc-review-plan` is an allowed shipped
skill.

Run the complete semantic suite. Do not add a live host-question smoke that
offline CI would have to invoke.

## Non-goals

Do not change ownership, delivery, or cleanup gates. Do not require
wall-clock timing or a provider call to prove the optional prompt.

## Verification

Use IPR-VT-01 through IPR-VT-05.

## Expected evidence

Fixture and allowlist diffs, passing routing/lifecycle/host/release tests,
and complete `test/acceptance.sh` output with no live question-prompt
invocation.

## Stop conditions

Stop if tests require a provider, if the skill fails the release allowlist,
or if review fixtures imply mutation or authorization.
