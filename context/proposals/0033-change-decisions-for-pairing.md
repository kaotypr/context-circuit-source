---
id: 0033-change-decisions-for-pairing
target_context_unit: context/DECISIONS.md
operation: change
statement: >-
  Adopt direct collaboration as an orthogonal working mode: one repository, user
  plus coordinator plus existing worker, no verifier/lease/execution record, a
  fresh isolated branch and working copy, a light resumable pointer, output
  labeled human-supervised rather than verified, explicit commits, preserved
  closure, separate delivery, and drift-blocking rather than silent rebasing.
evidence_refs:
  - sources/system-design/context-circuit/v0.7.0/direct-collaboration/design.md
  - sources/system-design/context-circuit/v0.7.0/direct-collaboration/session-and-isolation.md
  - sources/system-design/context-circuit/v0.7.0/direct-collaboration/skill-and-contract.md
  - wrapper/contracts/invariants.yaml#INV-PAIR-01
affected_repositories:
  - context-circuit-source
confidence: high
status: review-needed
---

## Proposed knowledge

Record the fixed direct-collaboration decisions and the implementation-specific
choices confirmed by the maintainer: one repository per session, block on anchor
drift, and close only clean state while preserving the branch and working copy.
