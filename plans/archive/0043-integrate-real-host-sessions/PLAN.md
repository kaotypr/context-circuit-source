# Plan 0043 — Integrate real Codex, Claude Code, and Cursor sessions

**Intent:** i030-integrate-real-host-sessions  
**Repository:** agent-harness  
**Tier:** Standard  
**Status:** draft

## Objective

Extend Agent Harness with a shared real-host lifecycle and genuine Codex,
Claude Code, and Cursor drivers. Each driver launches a new native root session
inside a disposable world, preserves observable native child topology and host
differences, reaches a native terminal state or declared harness limit, and
collects safe evidence. The same scenario can target one host or an explicit
three-host compatibility matrix.

## Grounding

The target is clean at revision
`3f25ccb187ba320acde306e227d334845d67b392`. The Python 3.12 package already
resolves and digests declarative scenarios, records prompt identity and host
facts, builds isolated run worlds, atomically finalizes manifests, and safely
retains or cleans them. It has no driver protocol, process ownership layer,
native event collector, host matrix command, or real-host tests; `README.md`
explicitly says host drivers are deferred.

## Decisions

- **One plan:** the externally meaningful result is a matrix-capable lifecycle.
  The shared contract, all three adapters, evidence boundary, orchestration, and
  end-to-end proof need one execution and independent-verification boundary.
- **Protocol around the foundation:** typed driver and evidence records extend
  the existing resolved-scenario and manifest seams without mixing host-specific
  process behavior into fixture construction.
- **Ambient authentication:** adapters use supported host interfaces and the
  host's existing account, passing only the prepared project, unchanged prompt,
  declared limits, and supported native flags. Credentials and configuration
  trees never enter the run.
- **Honest deterministic tests:** captured protocol records test parsing only.
  They can never satisfy the opt-in test that proves a new native session ran.
- **Independent matrix entries:** one resolved scenario expands to a distinct
  workspace and native root per host, with per-host reporting before summary.

## Tasks

| id | title | depends on |
|---|---|---|
| RH-001 | Define the shared real-host lifecycle and evidence boundary | — |
| RH-002 | Implement genuine native host drivers | RH-001 |
| RH-003 | Orchestrate selected-host and matrix runs | RH-002 |
| RH-004 | Prove the real-host story end to end | RH-003 |

Canonical paths, changes, acceptance statements, and runnable checks are in
`plan.yaml`.

## Risks

- Host CLIs and structured-event formats evolve independently. Each adapter must
  retain native version and capability evidence and fail honestly when an
  expected interface is unavailable.
- Child topology is not equally observable on every host. Missing evidence must
  remain unavailable or make a required assertion inconclusive; a separate
  process or inferred model label is not a child.
- Ambient authentication and provider events may contain secrets. Redaction must
  happen on the stream before persistence, with tests for configured secret and
  private-payload patterns.
- Timeout cleanup can affect unrelated work if process ownership is loose.
  Termination must remain rooted in the session launched for the run and record
  a harness-limit reason rather than success.
- Real-host tests consume time and possibly usage. They are explicit opt-in
  acceptance checks, while the default suite deterministically verifies
  contracts, parsing, orchestration, and safety.

## Verification

Run the complete unit suite and compile check, then explicitly opt into the
real-host acceptance test with configured Codex, Claude Code, and Cursor
installations. The acceptance story proves a new native root and terminal state
per host, required native child linkage when supported, independent matrix
workspaces, honest unavailable evidence, limit handling, and pre-persistence
redaction.
