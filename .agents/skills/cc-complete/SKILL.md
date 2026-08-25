---
name: cc-complete
description: Mark a verified plan done at explicit human request and reconcile the implementation against Product Knowledge.
---

Completion is human-controlled. Verification never marks a plan complete.

On "mark plan X complete", run the runtime `plan-complete`. It refuses unless the
latest execution passed independent verification; if the current result is
`failed` or `blocked`, report that plainly and do not change status. On success
it changes plan status `approved → done`, writes the implementation completion
record (plan revision, execution, per-repository commits, verifier result, human
request), and keeps branches and worktrees intact (no merge or publication).

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
