# Phase 4 coverage cross-check — Execute & verify

How the nine phase-4 plots cover the existing `scenarios/` cases and the harness gaps
they surface. Descriptive only: every rule cited is owned by
`wrapper/contracts/invariants.yaml` or `test/acceptance/criteria-map.yaml`. This is
the largest phase and the first with a mix of full-execution rewires and light
conversational probes.

## Existing cases → plots

| existing case | plot | relationship |
| --- | --- | --- |
| `05-approve-and-execute` | `execute-standard-verify-not-complete` | **evolves** — same seed/budget |
| `09-verifier-unavailable-host-blocked` | `verifier-host-blocked` | **evolves** — same seed/fault/budget |
| `11-repo-grounding` | `repository-grounding` | **evolves** — same seed/budget |
| `10-run-approved-stack` | `run-stack-single-repo` | **evolves** — same 10-plan seed/budget |
| `12-run-multi-repo-stack` | `run-stack-multi-repo` | **evolves** — same 7-plan/2-repo seed/budget |
| `13-execution-latency` | `execution-tiering-hidden` | **evolves** — same seed/budget |
| *(none)* | `three-failure-stop` | **new** — conv. probe (AC-13) |
| *(none)* | `interrupted-recovery` | **new** — conv. probe (AC-19) |
| *(none)* | `lease-ownership-conflict` | **new** — conv. probe (INV-CONCURRENCY-01/OWN-01) |

The six rewires keep their AC/invariant mappings verbatim — those were already
correct — so unlike phases 1–2 there is no mapping-drift correction this phase.

## The three new conversational probes

Each exercises the *lay-user explanation* of a mechanic already owned
deterministically (the mechanic is not re-tested; the plain-language reporting is):

- **`three-failure-stop`** (AC-13, INV-REPAIR-01, INV-PRESERVE-01) — the coordinator
  reports the stop faithfully ("couldn't get it to pass"), never as success, and
  confirms preservation. Deterministic owner: `test/runtime/test-runtime.sh`.
- **`interrupted-recovery`** (AC-19, INV-PRESERVE-01, INV-RUNTIME-02) — reads the
  interrupted state accurately, offers resume-or-hold, no silent redo/loss.
  Deterministic owner: `test/runtime/test-runtime.sh`.
- **`lease-ownership-conflict`** (INV-CONCURRENCY-01, INV-OWN-01, INV-PRESERVE-01) —
  a second request to touch in-progress work is held read-only, never a silent steal.
  Deterministic owner: the engine tests.

All three are `conversation-only` (claude-p) over a seeded terminal/mid state, so they
stay light: they observe how the outcome is *explained*, not whether the mechanic
fires. Keeping them conversation-only (rather than costly forced-failure
full-execution runs) is deliberate.

## AC gap: `lease-ownership-conflict`

The path lease / one-worker ownership is a v0.6 mechanic owned by INV-CONCURRENCY-01 /
INV-OWN-01 and the engine tests; the v0.5-rooted `criteria-map.yaml` has no dedicated
lease AC. The plot maps `acceptance_criteria: []` honestly and leans on the
invariants — the same pattern used where a rule is invariant-owned without a discrete
AC. If the maintainer wants an AC handle, that is a criteria-map decision, not a
library one.

## Harness capability gaps surfaced this phase

Flagged inline in the affected `generated/*.case.yaml`. None blocks authoring.

1. **New seed states for the probes.**
   - `three-failures-stopped` — an execution seeded in its terminal three-failure
     state (extends the mid-flight seeding case 09 already does with
     `worker-committed`).
   - `executing` **with a live lease held** — an in-progress execution holding a lease
     on a path region (`lease-ownership-conflict`). Seeding an *active* lease is new.
   - `interrupted-recovery` reuses the existing `worker-committed` seed, framed as an
     interruption — no new seed needed there.
2. **New post-condition predicates.**
   - `execution_stopped_three_failures`, `failure_evidence_preserved` (three-failure-stop).
   - `no_silent_redo` (interrupted-recovery).
   - `second_execution_blocked_read_only`, `lease_not_stolen`, `running_work_preserved`
     (lease-ownership-conflict).
   These assert the *conversational* outcomes; the deterministic mechanics are already
   covered by the runtime/engine suites, so these predicates need only be light
   state/heuristic checks in the human-sim grader.
3. **Driver note.** The six rewires are full-execution (`cc-test-case`); the three
   probes are conversation-only (`claude-p`). Access-discipline (dimension C) is a hard
   gate for the conversation-only probes and warning-only for the full-execution
   rewires, per `scenarios/README.md`.
