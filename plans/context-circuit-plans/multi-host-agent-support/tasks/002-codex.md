---
schema_version: 2
id: MHS-002
plan: multi-host-agent-support
status: ready
repository: context-circuit-source
paths:
  - AGENTS.md
  - wrapper/adapters/AGENTS.md
  - wrapper/adapters/WORKFLOW.md
  - .agents/skills/cc-entry/SKILL.md
  - .agents/skills/cc-execute/SKILL.md
  - .agents/skills/cc-verify/SKILL.md
  - agents/coordinator.md
  - agents/writer.md
  - agents/verifier.md
  - docs/host-capabilities.md
  - test/hosts/test-host-adapters.sh
depends_on: [MHS-001]
acceptance: [MHS-AC-01, MHS-AC-02, MHS-AC-03, MHS-AC-06]
verification: [MHS-VT-03, MHS-VT-05, MHS-VT-06]
expected_evidence:
  - Codex instruction-discovery and native-subagent mapping fixtures.
  - Root, writer, verifier, resume, and missing-child behavior evidence.
  - No new parallel role policy or host credential configuration.
stop_conditions:
  - Codex-specific behavior changes canonical route or role ownership.
  - A writer or verifier packet is incomplete, shared, or self-verified.
  - Live Codex access or sign-in becomes a required CI dependency.
---

# Complete the Codex CLI adapter and child-session mapping

## Objective

Make Context Circuit development predictable from Codex CLI by using its
project instruction discovery and native subagent capability as thin host
adapters around the existing filesystem workflow.

## Work

Align the source and shipped `AGENTS.md` surfaces with the `cc-*` skills and
role documents. Define how a Codex root session creates a bounded writer or
independent verifier delegation packet, how `/agent` or equivalent native
subagent activity is represented only as host evidence, and how unavailable
subagents produce the existing host-blocked result. Verify that Codex can
resume through the same receipt, Git, wrapper, lease, and worktree checks.

Use the host's existing permissions and authentication. The adapter must not
ask Codex to persist provider state or replace human approval cards with CLI
flags.

## Non-goals

Do not add a Codex SDK, custom provider integration, scheduler, or a second
writer/verifier role definition. Do not require a live Codex binary for the
offline acceptance suite.

## Verification

Use MHS-VT-03, MHS-VT-05, and MHS-VT-06.

## Expected evidence

Static instruction checks, Codex role-mapping fixtures, host-blocked and
offline results, context-budget measurements, and any explicitly authorized
read-only CLI smoke result with executable/version only.

## Stop conditions

Stop on a missing delegation field, a changed canonical owner, permission
weakening, a request to capture credentials, or any live-provider failure that
would block filesystem-only operation.

