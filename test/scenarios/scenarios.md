# v1.0 worked scenarios (fixture)

Four end-to-end walkthroughs plus the knowledge-debt path, each a human-simulated
fixture the driver (`test-scenarios.sh`) plays through the engine and asserts. The
human sees plain language (`conversations.md`); the steps below name the internal
mechanism the driver exercises.

## Scenario A — a Standard feature, one repository

1. Intent (Gate 1): `cc-intent` drafts goal, non-goals, outcome criteria, a coarse
   scope `checkout-service:[src/checkout]`, tier `standard`. Approval freezes the
   contract digest and confirms the ask; a read-only tracer then reads the real code
   and the coordinator runs the feasibility check.
2. Plan: `cc-plan` derives one plan from the trace manifest with two ordered embedded
   tasks sharing the same bounded execution and verification boundary. It is
   authorized by the approved intent — no second approval and no automated scope
   gate.
3. Execute: one worker commits; the independent verifier passes, bound to the
   candidate.
4. Accept: the human accepts the candidate (first-class, candidate-bound).
5. Deliver (Gate 2): the human authorizes the pull request; delivery is recorded.
6. Completion is inferred from acceptance + delivery (Standard); it emits a
   reconciliation-debt marker.

Human touch points: approve the intent, accept the candidate, authorize delivery,
resolve reconciliation.

## Scenario B — Explore then promote

1. Explore: no intent, no plan, no candidate — coordinator + one worker in an
   isolated working copy; the human watches live. No verifier, no tracer.
2. The work turns out to be real. Promote in place: attach an intent, raise the tier
   to Standard so a tracer reads the code, author a plan of record from its manifest.
   A candidate now exists and the independent verifier spawns. From here it is
   Scenario A.

Invariant: no plan file exists for the un-promoted Explore work.

## Scenario C — two stacked plans, one pull request

1. Two dependent plans (API first, consumer second), both in `checkout-service` and
   under one intent, execute independently. The dependency is recorded between the
   plans, not as a second approval decision.
2. At delivery the human groups them into one change set (one pull request).
3. Each plan keeps its own worker/verifier lifecycle. The change-set candidate is
   computed once over the combined tip set; the independent verifier runs once
   against it. One acceptance, one delivery.

## Scenario D — a criteria change re-gates

1. After approval, the human edits the intent's criteria. The frozen contract digest
   no longer matches, so the derived plan's authorization fails (CRITERIA_CHANGED) —
   the change re-enters Gate 1.
2. Execution is held; it never proceeds on the stale approval.
3. The human re-approves the changed criteria (a new decision that re-freezes the
   digest); only then is the plan authorized again. (A reach beyond the coarse scope
   is not gated here — it is a feasibility question the coordinator surfaces, and
   scope-safety is settled at delivery, Gate 2.)

## Scenario E — knowledge debt blocks the next change

1. A delivered/completed plan left a pending reconciliation-debt marker.
2. The next plan in the same knowledge scope is blocked at grounding
   (Standard/Critical) until the human reconciles or explicitly defers.

## Decomposition rejection cases

The one-plan fixture is invalid if it is split solely to match the assurance tier;
the bounded tasks belong together. The stacked fixture is invalid if its independent
partitions are combined into one lifecycle or if the dependency edge and its reason
are omitted. These are decomposition failures, not new approval or runtime gates.
