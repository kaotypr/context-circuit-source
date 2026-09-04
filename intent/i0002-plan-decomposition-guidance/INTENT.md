# Intention — i0002

_Status: draft, waiting for your approval._

## Intention

Make the handoff from tracing to planning explicit about the boundary between a
task and a plan.

Context Circuit should still use one plan containing several tasks when the work
is one bounded change with one worker lifecycle and one independent verification
boundary. When the traced work contains partitions that can be executed and
verified independently, or that have meaningful dependency or failure surfaces,
the coordinator should derive multiple stacked plans instead. The choice should
be explainable, grounded in the trace, and independent of the assurance tier.

## Expectations

- The current guidance clearly explains when several tasks belong in one plan and
  when they should become separate plans.
- A simple or normal Standard-tier change can remain a single plan with embedded
  tasks; plan count is not used as a proxy for risk tier.
- A larger change can produce multiple plans with explicit ordering, bounded
  scope, and a distinct worker/verifier lifecycle for each plan.
- The coordinator can explain the decomposition in human terms, including why it
  chose one plan or several.
- The change does not introduce a second plan-approval gate or alter ownership of
  delivery, candidates, verification, or scope safety.
- The behavior is covered by acceptance checks for both the one-plan and
  multi-plan cases.

## The plans

The work can be completed in these stages:

1. **Define the decomposition rule** — make the task-versus-plan boundary and
   the one-plan/multi-plan decision criteria explicit in the tracing and planning
   guidance.
2. **Carry the decision through the handoff** — update the tracer-to-planner
   representation and coordinator explanation so stacked plans are derived when
   the traced work requires them, while preserving the existing single-plan path.
3. **Prove both paths** — add or update deterministic checks and conversation-level
   coverage for bounded one-plan work and independently bounded stacked plans,
   then align the surrounding documentation.

## How carefully this is checked

**Standard**

This changes cross-cutting planning guidance and its supporting checks, so an
independent verifier should confirm that the decomposition is grounded in the
trace, that both paths work, and that no extra approval or authority is added.

## Open questions

None at the human-intent level. The exact representation of the decomposition
rationale can be chosen during tracing and planning, provided the rationale is
durable and reviewable.
