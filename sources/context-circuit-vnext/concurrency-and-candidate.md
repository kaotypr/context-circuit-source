# Concurrency × candidate — the trickiest interaction

This file resolves how the candidate (M2) composes with Context Circuit's existing
concurrency mechanics — path leases (INV-CONCURRENCY-01) and base selection with
integration merge (INV-CONCURRENCY-02). It is the interaction most likely to be got
wrong, so it is specified separately, and its unresolved edges are named honestly.

## The problem

M2 says "several stacked plans that converge to one pull request become one candidate
→ one verification." But the engine today works per **plan/execution**: each plan gets
its own lease, its own base (anchor / predecessor / integration merge), its own
worktree, and its own verifier. So "one candidate for N plans" needs a defined
**delivery unit** and a defined **verification point**, or it contradicts the
per-plan machinery.

## The delivery unit: a change set

Define a **change set** — the set of plan execution branches that will land in the
base branch as **one pull request**. A change set is:

- **derived, not new state**: it is the set of plans a human chooses to deliver
  together (today already expressed by delivering a run-stack's plans into one PR),
  plus their same-repo dependency closure.
- the unit the **candidate** is computed over: the candidate digest ranges over the
  *combined* per-repo tip set of the change set + their bases + the shared/So
  relevant contract digest(s).
- the unit **human acceptance** and **delivery** (Gate 2) attach to.

A single plan delivered alone is simply a change set of one — the common case, and
identical to today.

## The verification point: the integration tip

vNext reuses INV-CONCURRENCY-02's integration base, inverted for delivery:

- During execution, each plan still runs independently with its own lease, base, and
  worktree — **unchanged**. Per-plan attempts and repairs are unchanged.
- For a multi-plan change set, the **candidate is computed over the integration tip**
  — the composite of the change set's branches on their integration base, which the
  engine already knows how to author (`cc_base_prepare`, integration merge, with
  `BASE_UNBUILDABLE` on conflict).
- The independent verifier (at Standard/Critical) runs **once against that integration
  candidate**, not once per plan. That is the "one verification for one PR" (pain 6).

So per-plan execution is unchanged; what changes is that **evidence and acceptance
bind to the change-set candidate**, verified on the integration tip, rather than to
each plan's isolated attempt.

## Leases are unchanged and still per-plan

Path leases stay exactly as they are (per `(repo, path-region)`, held until delivery,
descendant-exempt). The change set does not weaken leasing:

- Plans in a change set that are dependency-related are already lease-exempt of each
  other (descendant rule); plans that are unrelated but co-delivered still each hold
  their own non-overlapping leases.
- Delivery still releases leases per the existing rule (held until delivered).

The candidate is a *view for evidence*, not a new lock. It grants no ownership
(consistent with "ownership grants exclusivity, never permission").

## Tiering interacts cleanly

- The change set's tier is the **max** tier among its member intents (fail upward): a
  change set containing one Critical intent is Critical.
- Explore work does not join a multi-plan change set (Explore is single-repo,
  recordless, human-supervised); it produces a candidate only on promotion, at which
  point it has a tier ≥ Standard.

## Drift and re-candidacy

The existing drift guard (`cc_delivery_rebase`) already rebases and forces
re-verification. Under vNext this is the general candidate rule: a rebase changes the
tip set → new candidate → prior verification and acceptance void → re-verify and
re-accept. No special case; the drift guard is the candidate rule applied at delivery.

## Honest open edges

- **What defines "one PR" precisely?** A human grouping (deliver these together) vs. a
  derived closure. Recommendation: explicit human grouping, defaulting to the
  same-repo dependency closure — never an automatic "merge everything ready," which
  would silently widen a change set. This must be a human choice at delivery.
- **Mixed contract digests.** A change set spanning multiple intents has multiple
  contract digests; the candidate incorporates all of them, and *any* of their
  criteria changing voids the candidate. This is correct but can cause churn on wide
  change sets — a reason to keep change sets small.
- **Integration-unbuildable at candidate time.** If the change set's integration tip
  will not build cleanly (`BASE_UNBUILDABLE`), there is no candidate to verify — the
  change set is blocked (not a worker failure), exactly as INV-CONCURRENCY-02 already
  treats an unbuildable base. The human splits or reorders. vNext adds no new
  atomicity claim across repositories (kept honesty; see `preserved-core.md`).
- **Per-plan vs change-set verification choice.** For a Standard change set a team may
  prefer to verify each plan *and* the integration; vNext requires only the
  integration verification but does not forbid per-plan checks. Recommendation: verify
  the integration candidate (the thing that ships); per-plan checks are optional
  belt-and-suspenders, not the gate.

## Summary

Execution stays per-plan and unchanged. The **change set** is the delivery unit; the
**integration tip** is where its single candidate is verified and accepted; **leases
and base selection are untouched**. The candidate is a binding for evidence, layered
on top of the existing concurrency mechanics, not a replacement for them.
