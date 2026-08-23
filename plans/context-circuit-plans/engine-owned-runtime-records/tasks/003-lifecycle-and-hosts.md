---
schema_version: 2
id: ERR-003
plan: engine-owned-runtime-records
status: ready
repository: context-circuit-source
paths:
  - wrapper/adapters/WORKFLOW.md
  - .agents/skills/cc-entry/SKILL.md
  - .agents/skills/cc-execute/SKILL.md
  - .agents/skills/cc-verify/SKILL.md
  - agents/coordinator.md
  - docs/runtime-contract.md
  - docs/host-capabilities.md
  - test/hosts/test-host-adapters.sh
  - test/recovery/test-recovery.sh
depends_on: [ERR-002]
acceptance: [ERR-AC-03, ERR-AC-04, ERR-AC-05, ERR-AC-06]
verification: [ERR-VT-04, ERR-VT-05]
expected_evidence:
  - Root, writer, and verifier flows consume the same engine-generated graph.
  - Resume rejects changed primary evidence and preserves interrupted authority.
  - Every host launches the bounded projection or reports host-blocked.
stop_conditions:
  - A skill or role becomes a second runtime-record implementation.
  - A missing host primitive triggers self-verification or role downgrade.
  - Resume reconstructs records from examples instead of current evidence.
---

# Integrate lifecycle, resume, and host adapters

## Objective

Make coordinator, writer, verifier, resume, and host mappings consume the same
generated records without copying constructor policy into conversational files.

## Work

Update thin adapters and role guidance to invoke constructors, transport the
bounded launch projection, consume the packet and handoff skeleton, and validate
current primary evidence on resume and completion. Preserve host-blocked and
filesystem-only outcomes exactly.

## Non-goals

Do not add host-specific routes, auto-invoke providers, self-verify, or allow a
resume request to steal ownership.

## Verification

Use ERR-VT-04 and ERR-VT-05.

## Expected evidence

Cross-host root/writer/verifier fixtures, missing-child outcomes, changed
receipt and interrupted-session results, and guidance diffs that cite owners.

## Stop conditions

Stop if adapters duplicate constructors, a host mapping broadens permissions,
or resume can proceed from stale or foreign evidence.
