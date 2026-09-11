# Plan 0044 — Recover real Codex, Claude Code, and Cursor session integration

**Intent:** i030-integrate-real-host-sessions  
**Repository:** agent-harness  
**Tier:** Standard  
**Status:** draft

## Objective

Implement the complete approved real-host outcome afresh from `agent-harness`
main. Stream genuine Codex, Claude Code, and Cursor protocols through
pre-persistence redaction, retain truthful native session evidence, and run one
scenario in complete independent per-host disposable worlds with meaningful
live acceptance.

## Grounding

The executable base is the clean registered repository at revision
`3f25ccb187ba320acde306e227d334845d67b392`. It has scenario resolution,
digest-verified artifact staging, synthetic Git fixtures, disposable run worlds,
manifests, finalization, retention, cleanup, and reproduction, but no host layer.

Plan 0043's failed execution is retained only as diagnostic evidence and is not
an implementation dependency. Its candidate buffered raw output before parsing,
inferred incomplete native facts, falsely graded a Claude authentication failure
as success, omitted usage and child topology, created directory shells rather
than complete matrix worlds/manifests, and let skipped live hosts count as proof.
This plan assumes none of those commits exists in the new worktree.

## Decisions

- **One recovery plan:** collection, redaction, grading, world ownership, and
  live proof share one trust boundary and require one independent verification.
- **Streaming is structural:** redact each stdout/stderr chunk or native record
  before any file, log, diagnostic, manifest, or retained evidence receives it.
- **Native evidence is never synthesized:** harness process identity cannot
  substitute for a native root, terminal, usage, child, parent, or role fact.
- **Failure beats optimistic completion:** authentication/error evidence and
  nonzero exits cannot be overridden by a superficially successful event.
- **Matrix entries are real runs:** each gets a built world, manifest, native
  root, evidence stream, finalization, retention decision, and result.
- **Requested live hosts cannot skip into proof:** every explicitly named host
  must satisfy native evidence requirements within realistic declared limits.

## Tasks

| id | title | depends on |
|---|---|---|
| RRH-001 | Define trustworthy streaming lifecycle and evidence semantics | — |
| RRH-002 | Implement native Codex, Claude Code, and Cursor adapters | RRH-001 |
| RRH-003 | Integrate independent selected-host and matrix run worlds | RRH-002 |
| RRH-004 | Establish deterministic and strict live acceptance | RRH-003 |

Canonical paths, acceptance statements, and runnable checks are in `plan.yaml`.

## Risks

- Native protocols vary by installed version; adapters must preserve native
  facts, reject contradictions, and report unsupported observations honestly.
- Provider output can carry secrets; malformed, stderr, exception, and timeout
  paths need the same before-write redaction guarantee.
- Loose process cleanup can harm unrelated sessions; ownership tests need
  unrelated sentinels and normal plus forced termination.
- Live checks consume time and usage and depend on device authentication; keep
  them explicit, synthetic, minimally prompted, and bounded, but never optional
  after the operator requests a host.

## Verification

Run the deterministic suite and compile checks, then the strict device command
for Codex, Claude Code, and Cursor. It must return nonzero unless every requested
host produces a fresh native root, trustworthy native terminal evidence, and all
scenario-required observable child facts within declared limits.
