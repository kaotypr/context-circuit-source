---
id: 0005-enrich-plan-approval
target_context_unit: context/domains/plan-approval/README.md
operation: change
statement: >
  The plan-approval page says "deterministic readiness checks pass" without
  listing them. Add the enumerated readiness checks the wrapper runs at approval.
evidence_refs:
  - .agents/skills/cc-execute/SKILL.md
  - wrapper/contracts/schemas/plan.yaml
  - wrapper/contracts/invariants.yaml   # INV-APPROVE-01
affected_repositories:
  - context-circuit-source
affected_commits:
  - 4b8ac0b
confidence: high
status: review-needed
---

# Enrich plan-approval with the enumerated readiness checks

Add the readiness checks that gate `draft -> approved`:

- objective represented accurately;
- repository mappings explicit;
- no task depends on an unknown repository;
- context conflicts resolved or accepted as risks;
- acceptance and verification testable;
- worker/verifier scope bounded;
- status is currently draft.

## Acceptance action

List these under Behavior as the content of "deterministic readiness checks".
Keep the rest of the page unchanged.
