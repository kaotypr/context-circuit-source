# Mechanism 4 — Closed knowledge loop

Fixes pain 8: forgetting to reconcile Product Knowledge after a merge, which
silently grounds the next plan on stale knowledge. This is the one pain that is a
*correctness* bug, not just friction — so it gets the firmest fix.

## The problem it solves

Today, completion *generates* a reconciliation reference
(`context-impact.yaml`, tied to INV-COMPLETE-02 / INV-KNOWLEDGE-02), but **acting**
on it is manual: reconciliation yields proposals under `context/proposals/`, and a
human must accept them (INV-KNOWLEDGE-02, correctly, keeps knowledge acceptance an
explicit human decision). Nothing prevents the *next* plan from being authored and
grounded while merged work sits unreconciled. So a plan can be grounded on Product
Knowledge that the last merge already made stale — exactly pain 8.

vNext keeps the human acceptance gate (it is right) but makes the **debt visible
and blocking** so it cannot be silently skipped.

## The mechanism: reconciliation debt

- On completion/delivery of a candidate, the engine emits a **reconciliation-debt
  marker** for the affected knowledge scope (new verb behavior; the marker records
  the candidate, the changed paths/commits, and the Product Knowledge units the
  plan relied on).
- A new verb `knowledge-debt <root>` lists every delivered candidate whose
  reconciliation is not yet resolved.
- A new verb `knowledge-reconciled <candidate>` clears the marker once the impact
  proposals have been generated and either accepted or explicitly deferred by a
  human.

## The gate: grounding blocks on debt

New behavior in `cc-plan`'s grounding preflight (owned via INV-KNOWLEDGE-02
reworked and INV-COMPLETE-02 reworked): before a new plan is grounded, it runs
`knowledge-debt`. If unreconciled delivered work exists in the knowledge scope the
new plan would rely on:

- **Critical/Standard tiers: block.** The coordinator must reconcile (generate the
  impact proposals, let the human accept or explicitly defer) before grounding
  proceeds. Forgetting becomes structurally impossible.
- **Explore tier: loud warn.** Fluidity is preserved for throwaway work, but the
  human is told, in plain language, "you are grounding on knowledge that merged
  work has not yet updated."

The human still *accepts* knowledge — the gate does not auto-accept anything
(INV-KNOWLEDGE-02 preserved). It only refuses to let the debt be **forgotten**.

## Why "defer" is a first-class outcome

Not every merge changes Product Knowledge, and a human must be able to say "no
durable update needed" cheaply. So `knowledge-reconciled` accepts an explicit
**defer/no-update** resolution that clears the debt without inventing a proposal —
matching today's `context-impact` `status: no-update-needed`. The point is that
*someone decided*, not that every merge produces a knowledge edit.

## Relationship to the candidate (M2)

The debt marker is keyed to the **candidate**, so it inherits candidate honesty: a
post-delivery correction that produces a new candidate carries its own debt, and a
voided candidate's debt is voided with it. Reconciliation is thus always tied to
the exact delivered result, never to a floating "the plan" whose code has since
moved.

## Optional: the production feedback edge

A "closed production loop" — where a shipped change keeps a live link to its
criteria and a production failure reopens the claim — is a natural extension.
vNext treats it as an **optional extension**, not core: the
reconciliation-debt machinery already gives the hook (a delivered candidate with a
durable link to its intent's criteria), so an inbound signal — via the existing
publication inbound-as-counter-evidence path (INV-EXTERNAL-02) — could later lower
confidence on the relied-on knowledge units or raise a new intent. It stays out of
the core because it depends on external wiring the product does not require, and it
must never become an automatic mutation of workspace state.

## Skills and engine

- **New verbs:** reconciliation-debt marker emission (on completion/delivery),
  `knowledge-debt`, `knowledge-reconciled`.
- **Changed:** `cc-plan` grounding preflight consults `knowledge-debt`;
  `cc-complete`/delivery emit the marker.
- **Unchanged:** the `context/` durable-knowledge format, the retrieval index
  (INV-KNOWLEDGE-01), the proposal/acceptance mechanic, and the human acceptance
  gate (INV-KNOWLEDGE-02) — vNext makes the loop *closed*, not *automatic*.
