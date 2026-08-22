---
schema_version: 2
id: MHS-004
plan: multi-host-agent-support
status: ready
repository: context-circuit-source
paths:
  - wrapper/adapters/AGENTS.md
  - wrapper/adapters/CLAUDE.md
  - docs/host-capabilities.md
  - docs/configuration.md
  - docs/getting-started.md
  - test/hosts/test-host-adapters.sh
  - test/security/test-boundaries.sh
depends_on: [MHS-001]
acceptance: [MHS-AC-01, MHS-AC-02, MHS-AC-05, MHS-AC-06]
verification: [MHS-VT-03, MHS-VT-04, MHS-VT-07]
expected_evidence:
  - Cursor CLI instruction-discovery fixtures for AGENTS.md and CLAUDE.md.
  - Interactive, print, resume, permission, and missing-child behavior evidence.
  - Proof that host-local Cursor configuration is optional and not committed.
stop_conditions:
  - A Cursor rule file becomes a parallel lifecycle or policy authority.
  - Print mode, force mode, or host-local permission settings bypass a Context Circuit gate.
  - Cursor authentication or MCP payloads are written to workspace state.
---

# Complete the Cursor Agent CLI adapter and permission boundary

## Objective

Make Cursor Agent CLI a first-class host adapter for Context Circuit without
requiring a deprecated `.cursorrules` file or introducing a second policy
surface.

## Work

Verify the released root `AGENTS.md` and `CLAUDE.md` surfaces are discovered by
Cursor CLI. Document the mapping for interactive sessions, `-p` print mode,
and `resume`; treat the host's project permission file as optional,
host-local, and subordinate to Context Circuit's human gates. Map any
available Task/subagent primitive to the same bounded writer and verifier
packets, and emit host-blocked when the verifier child cannot be created.

Keep the default smoke prompt read-only. If a mutating print-mode probe is
needed, run it only in a disposable fixture after explicit confirmation and
record only bounded outcome evidence.

## Non-goals

Do not add a `.cursorrules` compatibility layer, embed Cursor auth, configure
MCP servers, or treat `--force` as authorization. Do not require Cursor CLI or
network access in offline CI.

## Verification

Use MHS-VT-03, MHS-VT-04, and MHS-VT-07.

## Expected evidence

Static Cursor loading checks, command-mode fixtures, permission and
host-blocked outcomes, path/credential boundary results, and an optional
version-only live smoke handoff.

## Stop conditions

Stop on duplicate canonical policy, unsafe print-mode writes, auth/config
leakage, missing child evidence, or any request to make provider availability
an execution prerequisite.

