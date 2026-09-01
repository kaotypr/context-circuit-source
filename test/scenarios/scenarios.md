# v1.0 worked scenarios (fixture)

Four end-to-end walkthroughs plus the knowledge-debt path, each a human-simulated
fixture the driver (`test-scenarios.sh`) plays through the engine and asserts. The
human sees plain language (`conversations.md`); the steps below name the internal
mechanism the driver exercises.

## Scenario A — a Standard feature, one repository

1. Intent (Gate 1): `cc-intent` drafts goal, non-goals, criteria, scope
   `checkout-service:[src/checkout]`, tier `standard`. The spec adversary adds a
   latency criterion. Approval freezes the contract digest.
2. Plan: `cc-plan` derives one plan; the envelope check returns WITHIN — no second
   approval.
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
   isolated working copy; the human watches live. No verifier, no adversary.
2. The work turns out to be real. Promote in place: attach an intent, the adversary
   runs, raise the tier to Standard, author a plan of record. A candidate now
   exists and the independent verifier spawns. From here it is Scenario A.

Invariant: no plan file exists for the un-promoted Explore work.

## Scenario C — two stacked plans, one pull request

1. Two dependent plans (`0012` API, `0013` consumer), both in `checkout-service`,
   execute independently.
2. At delivery the human groups them into one change set (one pull request).
3. The change-set candidate is computed once over the combined tip set; the
   independent verifier runs once against it. One acceptance, one delivery.

## Scenario D — envelope drift re-gates

1. Mid-build, the plan turns out to also need `payments-lib`, a repository not in
   the approved scope. The envelope check returns EXCEEDS (new repository).
2. Execution is held and re-gated to the human.
3. The human widens the intent (a new decision), which re-freezes the contract
   digest and re-runs the adversary; only then does the plan proceed WITHIN.

## Scenario E — knowledge debt blocks the next change

1. A delivered/completed plan left a pending reconciliation-debt marker.
2. The next plan in the same knowledge scope is blocked at grounding
   (Standard/Critical) until the human reconciles or explicitly defers.
