# direct-collaboration

The v0.7.0 scope for **`cc-pair` — direct interactive collaboration**. When a
repo/project is bound, the user can just *work on it* with the agent: three actors
— **user, coordinator, worker** — in a loop, with **no verifier, no lease, no
execution records.** The user shares intent, the coordinator interprets ordinary
language and tells the worker what to do, the worker does it, and the user judges
the result live until "that's it."

It is **orthogonal to the plan lifecycle** — not a plan, not an execution, and
neither gates nor is gated by plan/approve/execute/verify/deliver. It can be
entered anytime a repo is bound, by intent or `/cc-pair`, or offered as a suggested
next step after a plan or stack execution completes.

This scope **supersedes the earlier `ui-refinement` framing** (and its
`interactive-acceptance` restatement), which mistakenly modeled interactive work as
a plan acceptance criterion verified by a final independent pass. Interactive work
is a way of *working*, not a way of *accepting a plan*.

## Reading order

1. [design.md](./design.md) — the normative overview: the one capability, where it
   sits (outside the lifecycle), why there is no verifier and why that is not a
   hole, the three actors, isolation, what it is *not*, the principles, and the
   fixed decisions. Stop here for a review-level understanding.
2. [session-and-isolation.md](./session-and-isolation.md) — where a session writes:
   own `cc-pair/<session>` branch + worktree from a base commit (never the active
   branch, never a plan branch in place), base selection, the light resumable
   pointer, `host-blocked`, convergence, and separate delivery.
3. [skill-and-contract.md](./skill-and-contract.md) — the `cc-pair` skill,
   triggering, its relation to the other skills, and the contract delta
   (`INV-PAIR-01`, a light session-pointer schema, the `.runtime/pairing/` folder;
   **no plan-contract change**).

## Authority

This scope drafts source only. It adds one skill (`cc-pair`) and one invariant
(`INV-PAIR-01`), a light session-pointer schema, and a `.runtime/pairing/` folder;
it changes **no plan contract**. The owner map in
`wrapper/contracts/invariants.yaml` settles the final ID on acceptance. It grants
no route, role, or authority by itself (INV-SKILL-01).
