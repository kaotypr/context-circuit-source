---
schema_version: 2
id: MHS-003
plan: multi-host-agent-support
status: done
repository: context-circuit-source
paths:
  - CLAUDE.md
  - wrapper/adapters/CLAUDE.md
  - scripts/release-artifact.sh
  - scripts/release-manifest.txt
  - wrapper/migrations/upgrade.sh
  - docs/host-capabilities.md
  - docs/getting-started.md
  - test/hosts/test-host-adapters.sh
  - test/release/test-release.sh
depends_on: [MHS-001]
acceptance: [MHS-AC-01, MHS-AC-02, MHS-AC-04, MHS-AC-06, MHS-AC-07]
verification: [MHS-VT-03, MHS-VT-08, MHS-VT-09]
expected_evidence:
  - Source and released CLAUDE.md adapters with a shared-instruction import.
  - Claude root/child/verifier and resume fixtures with host-blocked fallback.
  - Release inventory and migration evidence for the new shipped adapter.
stop_conditions:
  - Claude-specific instructions duplicate or contradict AGENTS.md policy.
  - Release assembly would omit CLAUDE.md or include maintainer/runtime state.
  - Authentication, MCP configuration, or provider transcripts enter workspace files.
---

# Complete the Claude Code adapter and release surface

## Objective

Make Claude Code load the same Context Circuit instructions and role contract
while preserving a thin, explicit host adapter in both the source checkout and
the released workspace.

## Work

Update the source `CLAUDE.md` as needed for maintainer development and add a
wrapper-owned `wrapper/adapters/CLAUDE.md` that imports the shipped
`AGENTS.md`. Add only Claude-specific notes for root/child session mapping,
Task/subagent availability, permission prompts, and resume evidence. Extend
release assembly, manifest checks, migration scope, and onboarding so the
released root includes the adapter and Claude can confirm it loaded it.

Verify that `claude`, print mode, and resume behavior all re-enter the same
human-facing workflow. A missing Task/subagent primitive must block the
bounded route rather than let the root self-verify.

## Non-goals

Do not copy Claude settings, `.mcp.json`, auth state, memory, or transcripts
into the workspace. Do not make `CLAUDE.md` a second policy owner or require
Claude network access for semantic CI.

## Verification

Use MHS-VT-03, MHS-VT-08, and MHS-VT-09.

## Expected evidence

Import and loading fixtures, root/child/verifier mapping, versioned artifact
inventory, upgrade/rollback preservation results, and an optional read-only
Claude smoke result containing no prompt or auth payload.

## Stop conditions

Stop if the import bridge is not portable, host instructions become
contradictory, release staging includes local Claude state, or a provider
failure blocks the core filesystem workflow.

