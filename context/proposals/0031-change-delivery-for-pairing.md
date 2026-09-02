---
id: 0031-change-delivery-for-pairing
target_context_unit: context/domains/delivery/README.md
operation: change
statement: >-
  A closed pairing branch may be delivered only by a separate explicit action and
  must be labeled human-supervised, not independently verified. Its target is the
  connected repository's recorded base branch; if the current base tip is not
  contained in the pairing branch, delivery blocks rather than silently rebasing.
evidence_refs:
  - sources/system-design/context-circuit/v0.7.0/direct-collaboration/session-and-isolation.md
  - wrapper/contracts/invariants.yaml#INV-PAIR-01
  - .agents/skills/cc-deliver/SKILL.md
  - test/pairing/test-pairing.sh
affected_repositories:
  - context-circuit-source
confidence: high
status: review-needed
---

## Proposed knowledge

Extend delivery with the pairing source/target and drift-blocking path. Preserve
the plan delivery drift guard unchanged; it remains the only path that rebases and
re-verifies automatically.
