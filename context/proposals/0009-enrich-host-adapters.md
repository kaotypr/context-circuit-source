---
id: 0009-enrich-host-adapters
target_context_unit: context/domains/host-adapters/README.md
operation: change
statement: >
  Enrich host-adapters with the shipped adapter duties from design chapter 06/07:
  the verifier read-only host-capability-or-block rule, the adapter's must-not
  rules (bypass approval, mark done after a verification pass, self-verify), and
  its positive runtime-record duties.
evidence_refs:
  - wrapper/adapters/AGENTS.md
  - wrapper/adapters/CLAUDE.md
  - agents/coordinator.md
  - wrapper/contracts/invariants.yaml   # INV-HOST-01
affected_repositories:
  - context-circuit-source
affected_commits:
  - 4b8ac0b
confidence: high
status: review-needed
---

# Enrich host-adapters with adapter duties and must-not rules

Additions (shipped):
- Verifier read-only capability requirement: the adapter must provide actual
  read-only capability for the verifier child, or report blocked; a "do not
  edit" prompt is not enough.
- Adapter must-not rules: must not bypass a missing approval, must not change the
  plan to done after a verification pass, must not self-verify when the verifier
  child is unavailable.
- Adapter positive duties: pass the execution brief, ask the runtime to record
  commits and evidence, coordinate repair attempts, return child results.

## Acceptance action

Fold into Behavior/Constraints alongside the existing INV-HOST-01 evidence-only
narrative.
