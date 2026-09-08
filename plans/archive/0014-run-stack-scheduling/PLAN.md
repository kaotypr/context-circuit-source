# 0014 — Run-stack scheduling, run action, and delivery drift guard

- **Plan ID:** `0014-run-stack-scheduling`
- **Status:** done
- **Repository:** `context-circuit-source`
- **Depends on:** `0013-run-stack-substrate`, `0006-verification-and-repair`, `0009-delivery`
- **Owns (invariants):** run_stack_action (WORKFLOW.md); extends INV-DELIVER-01 (drift guard)

## Original request

Retroactive plan for the run-stack scheduling and delivery drift guard, built as
if from an empty repo. Source design:
`sources/system-design/context-circuit/v0.6/run-stack/`
(`scheduling`, `delivery`, `design`, `examples`).

## Objective and desired behavior

- A plan becomes runnable when its dependencies are verified **and** its leases
  are free. The runtime detects readiness; the coordinator decides launch
  concurrency. Every plan is still separately approved, verified, completed, and
  delivered — run-stack adds no new authority.
- Failure containment: a failed plan holds only its descendants; unrelated
  siblings keep running.
- The run-stack conversational action executes only already-approved plans, via
  the `cc-run-stack` skill.
- Delivery drift guard: when a plan is delivered and its recorded base has
  diverged (a sibling already merged), the plan is rebased onto the current
  base tip and re-verified before its pull request opens.

## Constraints and non-goals

- Non-goal: repository grounding (0015–0016).

## Product Knowledge grounding

- `invariants` (`wrapper/contracts/invariants.yaml`), `design-deltas`
  (`context/DESIGN-DELTAS.md`).

## Tasks

1. **RSC-001** — readiness, scheduling loop, failure containment.
2. **RSC-002** — run-stack action + `cc-run-stack` skill.
3. **RSC-003** — delivery drift guard.

## Acceptance & verification

- Runs only when deps verified + leases free; failure blocks only descendants;
  run-stack executes only approved plans; drifted plan rebased + re-verified
  before PR.
- `sh test/run-stack/test-run-stack.sh`, `sh test/delivery/test-delivery.sh`.

## Assumptions, open questions, risks

- Risk: a scheduling heuristic leaking into the runtime — avoided by keeping
  detection in the runtime and launch decisions in the coordinator.

## Expected commits and delivery notes

Delivery order across a stack is dependency-ordered rebase trains.
