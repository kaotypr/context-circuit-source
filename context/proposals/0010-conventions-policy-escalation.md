---
id: 0010-conventions-policy-escalation
target_context_unit: context/CONVENTIONS.md
operation: change
statement: >
  Add two shipped conventions the design makes explicit: policy-change escalation
  (a change to repository identity, branch policy, security boundary, or
  execution behavior must update the owning workspace/contract file through that
  owner's action, not live only in a context page) and unequal evidence scope
  (repo shows what code does not what the product should do; a plan is intent not
  implementation; a verifier confirms behavior not a domain decision).
evidence_refs:
  - wrapper/contracts/invariants.yaml   # one-owner-per-rule; INV-KNOWLEDGE-02
  - docs/product-knowledge.md
affected_repositories:
  - context-circuit-source
affected_commits:
  - 4b8ac0b
confidence: medium
status: review-needed
---

# Add policy-change escalation and evidence-scope conventions

Add to CONVENTIONS.md:
- Policy-change escalation: an update that changes repository identity, branch
  policy, a security boundary, or execution behavior must also update the owning
  workspace or contract file through that owner's normal action. A policy change
  must not live only in a context page. (Extends the existing one-owner-per-rule
  convention.)
- Unequal evidence scope: repository evidence shows what the code does, not what
  the product should do; a plan is intent, not implementation; a verifier
  confirms behavior, not a domain decision.

## Acceptance action

Append these as two short convention sentences; keep the file terse.
