# 0011 — Human-simulated test harness mechanism

- **Plan ID:** `0011-human-simulated-harness`
- **Status:** done
- **Repository:** `context-circuit-source`
- **Depends on:** `0010-release-assembly-and-publication`
- **Owns:** the source-only harness mechanism (maintainer tooling, never shipped)

## Original request

Retroactive plan for the human-simulated harness mechanism, built as if from an
empty repo. Source design:
`sources/system-design/context-circuit/v0.5/template-harness/`.

## Objective and desired behavior

- Exercise the *assembled* `context-circuit-template` as a real, isolated
  workspace: a maintainer-only human-simulator (an ignorant lay user) talks to
  the real product coordinator, and a deterministic grader checks whether the
  product experience is correct.
- The coordinator-under-test is the real product coordinator loaded from the
  instantiated workspace, spawned with no awareness it is under test.
- The grader is mechanical (no model call): state post-conditions and transcript
  checks are hard gates; the access audit's forbidden reads gate; the efficiency
  ledger is soft. Every finding maps to an acceptance criterion and an invariant.
- Source-only maintainer tooling; nothing ships into the template.

## Constraints and non-goals

- Non-goal: the full scenario library and host matrix (0012); the efficiency
  ledger being made real (0020).

## Product Knowledge grounding

- `invariants` (`wrapper/contracts/invariants.yaml`), `architecture`
  (`context/ARCHITECTURE.md`).

## Tasks

1. **HRN-001** — human-simulator role + coordinator spawn contract.
2. **HRN-002** — scenario schema, runner, anchor case.
3. **HRN-003** — deterministic grader dimensions A–D.

## Acceptance & verification

- Coordinator spawned test-blind; a run captures transcript/trace/telemetry;
  state+transcript are hard gates, ledger soft.
- `sh template-harness/test-template-runtime.sh` (deterministic lab gate; the
  live human suite is agent-driven and opt-in).

## Assumptions, open questions, risks

- Risk: the human-simulator leaking test framing to the coordinator — prevented
  by the spawn contract giving it only the persona block.

## Expected commits and delivery notes

Source-only; excluded from the release artifact by the manifest.
