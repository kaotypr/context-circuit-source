---
id: 0030-change-host-adapters-for-pairing
target_context_unit: context/domains/host-adapters/README.md
operation: change
statement: >-
  A native child may map to the existing worker role for direct collaboration as
  well as plan execution; pairing never launches a verifier, and inability to
  create the worker child yields host-blocked with no coordinator write fallback.
evidence_refs:
  - sources/system-design/context-circuit/v0.7.0/direct-collaboration/design.md
  - sources/system-design/context-circuit/v0.7.0/direct-collaboration/skill-and-contract.md
  - wrapper/adapters/AGENTS.md
  - wrapper/adapters/CLAUDE.md
  - agents/coordinator.md
affected_repositories:
  - context-circuit-source
confidence: high
status: review-needed
---

## Proposed knowledge

Extend the accepted host-adapter domain without changing INV-HOST-01: child
capability remains bounded evidence, and missing required capability fails closed.
