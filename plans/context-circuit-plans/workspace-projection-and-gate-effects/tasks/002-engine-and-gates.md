---
schema_version: 2
id: WPE-002
plan: workspace-projection-and-gate-effects
status: draft
repository: context-circuit-source
paths:
  - wrapper/runtime/engine.sh
  - docs/gates.md
  - template/context/WORKSPACE.md
  - template/context/PROJECT.md
  - template/context/INDEX.md
  - test/gates/test-gates.sh
  - test/recovery/test-recovery.sh
depends_on: [WPE-001]
acceptance: [WPE-AC-01, WPE-AC-02, WPE-AC-03, WPE-AC-04, WPE-AC-05]
verification: [WPE-VT-02, WPE-VT-03, WPE-VT-04]
expected_evidence:
  - Atomic render and validation results for all three summaries.
  - Cards showing complete current and later effect sections.
  - Failure fixtures proving no partial projection or post-confirmation default invention.
stop_conditions:
  - Rendering overwrites user-owned content without an accepted format contract.
  - A partial write can leave workspace summaries disagreeing.
  - Card wording changes authorization or combines separate gates.
---

# Render and validate projections atomically

## Objective

Implement the contract through host-neutral engine operations and canonical
gate cards while preserving mutable workspace data and separate human gates.

## Work

Add bounded render and validation operations for the required summaries,
staging every output before publication. Integrate validation into entry and
write preflight. Update canonical cards to show current and later effects and
reject omitted values before confirmation. Keep repository-local paths only in
the ignored binding file.

## Non-goals

Do not create a repository, change the plan lifecycle, combine approval with
commit or execution, or introduce provider-specific UI state.

## Verification

Use WPE-VT-02, WPE-VT-03, and WPE-VT-04.

## Expected evidence

Atomic-update fixtures, mismatch failures, complete cards, interruption tests,
and a path audit showing only declared mutable seed and wrapper files changed.

## Stop conditions

Stop on ambiguous generated/user-owned boundaries, any partial publication,
or any gate presentation that implies authorization beyond the current card.
