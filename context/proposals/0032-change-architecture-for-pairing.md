---
id: 0032-change-architecture-for-pairing
target_context_unit: context/ARCHITECTURE.md
operation: change
statement: >-
  Context Circuit has two orthogonal repository-changing modes: the existing
  plan execution path for unwatched work with independent verification, and
  cc-pair for live human-supervised work with a coordinator and one worker. The
  pairing mode does not alter or bypass any plan lifecycle gate.
evidence_refs:
  - sources/system-design/context-circuit/v0.7.0/direct-collaboration/design.md
  - wrapper/contracts/invariants.yaml#INV-PAIR-01
  - wrapper/adapters/WORKFLOW.md
  - .agents/skills/cc-pair/SKILL.md
affected_repositories:
  - context-circuit-source
confidence: high
status: review-needed
---

## Proposed knowledge

Add the direct-collaboration path beside—not inside—the core execution model and
retain the invariant catalog as the one-owner-per-rule authority.
