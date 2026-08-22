---
schema_version: 2
id: VEL-002
plan: verification-evidence-layers
status: draft
repository: context-circuit-source
paths:
  - wrapper/runtime/engine.sh
  - .agents/skills/cc-verify/SKILL.md
  - agents/verifier.md
  - docs/planning.md
  - docs/runtime-contract.md
  - test/runtime/test-runtime.sh
depends_on: [VEL-001]
acceptance: [VEL-AC-02, VEL-AC-03, VEL-AC-04, VEL-AC-05]
verification: [VEL-VT-02]
expected_evidence:
  - Runtime comparison rejects weaker or missing observed layers.
  - Verifier handoff reports passed, failed, blocked, and waived criteria separately.
  - Writer claims never substitute for independently observed evidence.
stop_conditions:
  - Verifier guidance becomes a second evidence-policy owner.
  - A waiver or unavailable capability is represented as pass.
  - The verifier repairs implementation or writes another session's runtime.
---

# Enforce required layers in verifier evidence

## Objective

Make independent verification compare required and observed evidence and
refuse false equivalence between lower-layer checks and promised behavior.

## Work

Implement host-neutral validation and thin verifier guidance. Carry evidence
layers through the generated packet and handoff, compare each criterion, and
report distinct pass, fail, blocked, and waived outcomes. Preserve read-only
verification and treat writer evidence as a claim.

## Non-goals

Do not invoke browser tooling, repair implementation, approve waivers, or
duplicate comparison policy in the skill or role file.

## Verification

Use VEL-VT-02.

## Expected evidence

Schema-versus-browser, build-versus-interaction, sufficient-process,
unavailable-host, waiver, writer-claim, and verifier-read-only results.

## Stop conditions

Stop if weaker evidence can pass, outcomes collapse into a single success
state, or verifier behavior crosses write or repair boundaries.
