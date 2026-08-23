---
schema_version: 2
id: RCP-002
plan: route-context-packet-enforcement
status: ready
repository: context-circuit-source
paths:
  - wrapper/runtime/engine.sh
  - wrapper/adapters/AGENTS.md
  - wrapper/adapters/CLAUDE.md
  - wrapper/adapters/WORKFLOW.md
  - .agents/skills/cc-entry/SKILL.md
  - test/context-budget/test-budgets.sh
  - test/hosts/test-host-adapters.sh
  - test/security/test-boundaries.sh
depends_on: [RCP-001]
acceptance: [RCP-AC-03, RCP-AC-04, RCP-AC-05]
verification: [RCP-VT-03, RCP-VT-04, RCP-VT-05]
expected_evidence:
  - One host-neutral loader result with selected paths, actual bytes, and digests.
  - Over-budget and out-of-set requests fail before context silently broadens.
  - Host adapters invoke shared primitives without copying route policy.
stop_conditions:
  - Enforcement depends on credentials, provider payloads, or a host-specific router.
  - Direct host reads are represented as sandboxed when the host cannot enforce them.
  - A stale receipt is allowed to authorize a write.
---

# Load and account for bounded packets

## Objective

Turn packet declarations into measured, host-neutral loading and receipt
evidence without misrepresenting host capabilities.

## Work

Add engine operations that resolve allowed packet paths, validate selected
evidence, measure bytes, enforce ceilings, and emit receipt references and
digests. Keep adapters thin: invoke the operation, load only returned evidence,
and report an honest blocked result when required enforcement is unavailable.

## Non-goals

Do not embed a host SDK, record prompts or transcripts, generate child runtime
records, or make host permission mode an authorization grant.

## Verification

Use RCP-VT-03, RCP-VT-04, and RCP-VT-05.

## Expected evidence

Measured success, overrun, undeclared-path, stale-receipt, host-unavailable,
source-boundary, and cross-host adapter fixtures.

## Stop conditions

Stop if path or byte enforcement is advisory only, if host limitations are
reported as success, or if receipt evidence can replace primary evidence.
