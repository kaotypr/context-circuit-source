---
name: cc-complete
description: Mark a verified plan done at explicit human request and reconcile the implementation against Product Knowledge in place.
---

Completion is gated by the intent's consequence tier (INV-COMPLETE-01), and
verification never marks a plan complete.

- **Standard and Critical — explicit mark-done.** On "mark plan X complete" or
  "mark X done", run `plan-complete . <plan-id>`. It refuses unless the tier
  floor is met: a candidate-bound independent pass that still describes the
  current candidate (`completion-ready`). Delivery, candidate acceptance, and
  verification do not mark the plan done.
- **Explore — planless.** Explore work is direct human-supervised collaboration;
  it has no plan-of-record completion path until the human promotes it.

On success `plan-complete` changes plan status `draft → done`, writes the
implementation completion record (execution, per-repository commits,
`human_completion: accepted`), and keeps branches and worktrees intact (no
merge or publication). If the current result is `failed` or `blocked`, or the
candidate is stale, report that plainly and do not change status.

Then, **if the plan affected Product Knowledge**, reconcile in place — the same
act as gathering context. Compare the plan's knowledge references, changed
paths, worker handoffs, and verifier evidence with the live context units.
Edit live `context/` files and keep `INDEX.md` consistent only when there is a
durable knowledge change. If the plan did not affect Product Knowledge, leave
context files as they are. Record optional reconciliation references with
`context-impact-record`. Write live context files only; there is no sidecar
staging path and no extra knowledge-acceptance gate.

In-place edits are a coordinator act after mark-done. The runtime does not
write or interpret Product Knowledge (INV-RUNTIME-01, INV-COMPLETE-02). A later
plan may start even if this update has not landed (INV-KNOWLEDGE-02).
