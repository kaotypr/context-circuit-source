# 0012 — Scenario library and host matrix

- **Plan ID:** `0012-scenario-library-and-host-matrix`
- **Status:** done
- **Repository:** `context-circuit-source`
- **Depends on:** `0011-human-simulated-harness`
- **Owns:** the harness scenario library and host bindings (source-only)

## Original request

Retroactive plan for the harness scenario library and host matrix, built as if
from an empty repo. Source design:
`sources/system-design/context-circuit/v0.5/template-harness/scenario-library.md`
and `host-matrix.md`.

## Objective and desired behavior

- A full case library exercising the lifecycle: orientation, connect, plan
  review/approval, refuse-unapproved-execution, approve-and-execute,
  repair-then-complete, archive/restore, delivery-boundary, and the
  verifier-unavailable → host-blocked path.
- Full-execution cases require nested worker/verifier spawning.
- Per-host bindings (Claude Code / Codex CLI / Cursor Agent) cover only what
  differs per host; an opt-in matrix runner iterates {scenario × host} and
  reports keyed by (scenario, host) with a trace-availability flag.
- The matrix runner stays opt-in and outside the deterministic acceptance gate.

## Constraints and non-goals

- Non-goal: the efficiency-ledger telemetry (0020).

## Product Knowledge grounding

- `invariants` (`wrapper/contracts/invariants.yaml`), `architecture`
  (`context/ARCHITECTURE.md`).

## Tasks

1. **SCN-001** — the full case library.
2. **SCN-002** — per-host bindings + matrix runner.

## Acceptance & verification

- Each case has a human + grader block mapped to invariants; matrix runner
  reports per (scenario, host) with trace availability.
- `sh template-harness/test-template-runtime.sh`.

## Assumptions, open questions, risks

- Risk: a host without a trace silently reporting a green ledger — prevented by
  the trace-availability degradation.

## Expected commits and delivery notes

Source-only; excluded from the release artifact.
