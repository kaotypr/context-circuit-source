# ui-refinement

The v0.7.0 scope for **interactive, human-gated refinement** — the mode Context
Circuit enters when the acceptance criterion is subjective and best judged by a
human looking at a rendered artifact ("refine the UI until it looks good"). Its
driving instance is frontend/UI layout; the mechanism generalizes.

## Reading order

1. [design.md](./design.md) — the normative overview: the one capability, the
   problem, principles, the central "no new role" decision, what changes relative
   to v0.6, and the shape of the whole. Stop here for a review-level understanding.
2. [roles-and-expertise.md](./roles-and-expertise.md) — why there is **no new
   role**: roles are authority-shaped, "expert frontend engineer" is domain
   competence, and where that competence actually attaches (the worker's brief and
   the coordinator's fluency).
3. [refinement-loop.md](./refinement-loop.md) — the interactive execution mode: the
   session-scoped worker, the turn cadence, the two rendering channels, the
   `host-blocked` fallback, and how the lease/worktree bracket the loop.
4. [acceptance-and-verification.md](./acceptance-and-verification.md) — the
   human-gated visual acceptance type, the single end verify, and the
   visual-regression baseline freeze that turns a subjective approval into a
   durable, observable gate.
5. [skill-and-schema.md](./skill-and-schema.md) — the `cc-refine` skill, its
   relation to `cc-plan` / `cc-execute` / `cc-verify` / `cc-deliver`, triggering,
   and the contract delta.

## Authority

This scope drafts source only. It adds one skill and one acceptance kind, and
proposes provisional `INV-REFINE-01/02`; the owner map in
`wrapper/contracts/invariants.yaml` settles the final IDs on acceptance. It grants
no route, role, or authority by itself (INV-SKILL-01).
