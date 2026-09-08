# 0020 — Harness efficiency ledger (dimension D)

- **Plan ID:** `0020-efficiency-ledger`
- **Status:** done
- **Repository:** `context-circuit-source`
- **Depends on:** `0011-human-simulated-harness`
- **Owns:** the harness efficiency ledger (source-only tooling)

## Original request

Retroactive plan for making the harness efficiency ledger real, built as if from
an empty repo. Source design:
`sources/system-design/context-circuit/v0.6/template-harness/`
(`budgets`, `evaluation`, `telemetry`).

## Objective and desired behavior

- Turn dimension D (the efficiency ledger) from declared-but-dormant into
  measured: the driver derives per-action usage telemetry from each coordinator
  turn's own `result` event (no extra model call), the budget units are pinned
  (output tokens, conversational turns, agent turns, cost; `context_peak` is a
  diagnostic only), and the grader compares observed usage to each case's
  `budgets`.
- Dimension D stays soft and warning-only; state, transcript, and the access
  audit's forbidden reads remain the hard gates. A run with no telemetry degrades
  D to "unavailable", never a failure.
- Source-only maintainer tooling; no product surface changes, nothing ships.

## Constraints and non-goals

- Non-goal: any product gate. Dimension D never changes the verdict.

## Product Knowledge grounding

- `invariants` (`wrapper/contracts/invariants.yaml`), `architecture`
  (`context/ARCHITECTURE.md`).

## Tasks

1. **EFL-001** — telemetry emission, budget units, grader compare, recalibration.

## Acceptance & verification

- Per-action ledger reported and soft; no-telemetry run degrades to unavailable.
- `sh template-harness/human/run-scenario.sh` (live driver; the deterministic
  `sh test/acceptance.sh` gate is unaffected).

## Assumptions, open questions, risks

- Additive TSV extension; non-recalibrated cases still grade softly against old
  numbers; independent of the three v0.6 product scopes.

## Expected commits and delivery notes

Source-only; excluded from the release artifact.
