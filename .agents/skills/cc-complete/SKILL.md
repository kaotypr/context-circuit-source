---
name: cc-complete
description: Mark a plan done at explicit human request and reconcile the implementation against Product Knowledge in place.
---

A plan becomes done only when a human asks to mark it done. Verification,
candidate acceptance, and delivery never mark a plan complete.

- **Standard and Critical — explicit mark-done.** On "mark plan X complete" or
  "mark X done", or the same ask naming several plans, run `plan-complete .
  <plan-id>` for each named plan. There is no look at that plan's work or
  evidence, and no unreadiness refusal. A plan that was never built, failed a
  check, or has no evidence still becomes done when asked. Delivery, candidate
  acceptance, and verification do not mark the plan done.
- **Explore — planless.** Explore work is direct human-supervised collaboration;
  it has no plan-of-record completion path until the human promotes it.

On success `plan-complete` changes plan status `draft → done` and keeps
branches and worktrees intact (no merge or publication). When execution
evidence files exist it may also write the implementation completion record
(`human_completion: accepted`); missing files do not block the status flip.

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
