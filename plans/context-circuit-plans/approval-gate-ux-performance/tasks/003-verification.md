---
schema_version: 2
id: AGF-003
plan: approval-gate-ux-performance
status: done
repository: context-circuit-source
paths:
  - test/approval/test-approval-flow.sh
  - test/routing/fixtures.yaml
  - test/lifecycle/test-lifecycle.sh
  - test/gates/test-gates.sh
  - test/ownership/test-ownership.sh
  - test/acceptance.sh
depends_on: [AGF-001, AGF-002]
acceptance: [AGF-AC-01, AGF-AC-02, AGF-AC-03, AGF-AC-04, AGF-AC-05, AGF-AC-06]
verification: [AGF-VT-01, AGF-VT-02, AGF-VT-03, AGF-VT-04, AGF-VT-05]
expected_evidence:
  - Dedicated two-turn approval fixtures covering no-op pre-confirmation, exact confirmation, invalid lifecycle, formatting preservation, product-source follow-on commit card, and wrapped-workspace absence of that card.
  - Existing routing, lifecycle, gate, host, ownership, security, recovery, and release suites passing.
  - Independent review handoff with changed paths, acceptance mapping, limitations, and remaining human gates.
stop_conditions:
  - Tests pass while the diff contains unrelated files or body/formatting changes.
  - Verification relies on the implementing session's claim without reproducing the bounded evidence.
  - Any existing safety, ownership, or release boundary regresses.
---

# Add regression coverage and complete independent acceptance

## Objective

Turn the observed approval-session failure mode, and the product-source
dirty-base stop discovered at run time, into durable acceptance evidence
without baking in a flaky wall-clock threshold.

## Work

Add a focused approval-flow suite and wire it into semantic acceptance. Exercise
the initial card, exact confirmation, invalid or stale confirmation, status
projection, unchanged task bodies/endings, explicit separation from execution,
product-source presentation of the existing commit card immediately after
confirm approval, wrapped-workspace omission of that card, and continued
`MAINTAINER_APPROVAL_COMMIT_REQUIRED` / `DIRTY_BASE_BLOCKED` classification.
Re-run the existing suites that own routing, lifecycle, gates, ownership, host
behavior, security, recovery, and release boundaries. Produce an independent
verifier handoff rather than treating test output as plan completion.

## Non-goals

Do not mark the plan done, commit the maintainer approval projection during
verification, treat Run approved plan as the owner of the commit card, or test
external provider availability as a prerequisite for filesystem evidence.

## Verification

Use AGF-VT-01 through AGF-VT-05.

## Expected evidence

Focused test output, complete semantic acceptance output, changed-file and
diff-scope inspection, and an independent handoff identifying any remaining
latency limitations.

## Stop conditions

Stop on any regression in gate authorization, lifecycle projection, ownership,
release exclusion, or independent reproducibility.
