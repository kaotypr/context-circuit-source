---
id: 0029-add-direct-collaboration
target_context_unit: context/domains/direct-collaboration/README.md
operation: add
statement: >-
  Context Circuit supports cc-pair as a one-repository direct-collaboration mode
  outside the plan lifecycle: user, coordinator, and one worker iterate in an
  isolated cc-pair branch and working copy; the user is the live oracle, no
  verifier/lease/execution record is created, only a light resumable pointer is
  kept, output is human-supervised rather than verified, and completion and
  delivery remain separate.
evidence_refs:
  - sources/system-design/context-circuit/v0.7.0/direct-collaboration/design.md
  - sources/system-design/context-circuit/v0.7.0/direct-collaboration/session-and-isolation.md
  - sources/system-design/context-circuit/v0.7.0/direct-collaboration/skill-and-contract.md
  - wrapper/contracts/invariants.yaml#INV-PAIR-01
  - .agents/skills/cc-pair/SKILL.md
  - test/pairing/test-pairing.sh
affected_repositories:
  - context-circuit-source
confidence: high
status: review-needed
---

## Proposed knowledge

Add a bounded direct-collaboration domain covering entry, the live loop,
single-repository isolation, resumability, convergence, human-supervised labeling,
host-blocked behavior, and separate delivery. The accepted page should reference
INV-PAIR-01 rather than duplicate its policy.
