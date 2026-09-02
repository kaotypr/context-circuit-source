---
name: cc-complete
description: Mark a verified plan done at explicit human request and reconcile the implementation against Product Knowledge.
---

Completion is gated by the intent's consequence tier (INV-COMPLETE-01), and
verification never marks a plan complete.

- **Standard — inferred.** Completion is a projection of "candidate accepted +
  delivered". After the human accepts the candidate and delivery is recorded
  (Gate 2), run `completion-infer . <plan-id>`; it sets `draft → done` when
  `completion-ready` passes and a delivery record still binds to the current
  candidate. A post-delivery change re-gates rather than completing stale work.
- **Explore — planless.** Explore work is direct human-supervised collaboration;
  it has no plan-of-record completion path until the human promotes it.
- **Critical — explicit.** An explicit human completion is still required: on "mark
  plan X complete", run `plan-complete`. It refuses unless the tier floor is met
  (a candidate-bound independent pass). Standard uses the inferred path; it has no
  second completion gate.

On success either path changes plan status `draft → done`, writes the
implementation completion record (execution, per-repository commits, whether the
completion was accepted or inferred), and keeps branches and worktrees intact (no
merge or publication). If the current result is `failed` or `blocked`, or the
candidate is stale, report that plainly and do not change status.

Then reconcile the implementation against Product Knowledge. Compare the actual
committed changes, changed paths, worker handoffs, and verifier evidence with the
context units that grounded the plan. Produce one of:

- no durable knowledge change;
- one or more context update proposals under `context/proposals/` (target unit,
  operation, statement, evidence, affected repositories/commits, confidence,
  conflict);
- a stale/conflict warning.

Record the reconciliation references with `context-impact-record` and index
pending impacts. Never silently accept a Product Knowledge change: accepting a
proposal is a separate explicit human decision ("accept the context update for
X"). The plan may remain done while a proposal is pending; surface a relevant
pending impact during future plan creation.

## Reconciliation debt (closed knowledge loop, INV-COMPLETE-02)

Completion emits a **reconciliation-debt marker** keyed to the accepted candidate.
Until it is resolved, the next Standard/Critical plan's grounding in the same
knowledge scope blocks — so reconciliation cannot be silently skipped. Explore is
planless and has no plan grounding preflight. When you have generated the proposals
and the human has accepted or explicitly deferred them, clear the marker with `knowledge-reconciled .
<candidate> reconciled` (or `deferred` for an explicit "no durable update needed").
`knowledge-debt .` lists every delivered candidate still awaiting reconciliation.
Clearing the marker is bookkeeping that the decision was made; it never accepts
Product Knowledge on the human's behalf.
