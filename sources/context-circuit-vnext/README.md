# Context Circuit vNext — an intent-gated, candidate-proven evolution

**Status: source material only.** This is a design study. It has no status, no
authority, and changes nothing on its own. It is not a plan, not accepted Product
Knowledge, not a release, and not permission to change the product. Read it as an
argument to test, not an answer to adopt. The working name "vNext" is a
design-line label, not a rebrand of the product; the product stays Context
Circuit unless a separate, explicit vocabulary decision says otherwise.

## What this is

A proposal for **major changes to Context Circuit's trust core**, keeping its
mechanics. Unlike the two from-scratch studies filed beside it
(`claude-redesign-context-circuit/` "Proof Circuit" and
`codex-redesign-context-circuit/` "Proof Circuit"), this study was written *after
a full line-by-line read of the current core* — `wrapper/runtime/engine.sh`
(1,975 lines), `wrapper/contracts/invariants.yaml`, the host adapters, and the
skill/engine coupling map. Its central finding is that **the current core is
good** — it is already the small, model-blind, deterministic mechanics kernel
that every from-scratch redesign independently reinvented — so the right move is
to **evolve the policy layer and preserve the mechanics**, not to start over.

## The thesis in one sentence

> Context Circuit today gates the *state transitions* of a plan; vNext gates the
> two things a human actually decides — **what "correct" means** (the intent) and
> **what ships** (delivery) — and makes everything between them mechanical,
> consequence-tiered, and self-invalidating.

## Where it came from

This study is the product of two inputs:

1. A survey of six independent "design an AI-agent workspace" studies
   (`workspace-design-by-gpt-5.6-sol` "Durable Workspace Fabric",
   `collaborative-workspace` "Bench", `universal-workspace` "Lattice",
   `workspace-design-by-opus.5` "Drydock", and the two Proof Circuit redesigns).
   Across all six, a small cluster of ideas recurred independently — the
   strongest signal for what a workspace core should do.
2. A grounded read of Context Circuit's own core, to separate **hard-won
   correctness worth preserving** from **policy worth changing**.

vNext keeps what the read showed is load-bearing and adopts only the convergent
ideas that fit Context Circuit's single-host, human-in-the-loop, provider-neutral
shape. It deliberately declines the convergent ideas the two CC-specific
redesigns also declined (fencing tokens, a full event-sourced rewrite, a
context/token compiler) — they solve problems this product does not yet have, and
the CC-specific redesigners reached the same conclusion.

## The eight lived pains this design answers

vNext is anchored to real usage friction, not theory:

1. No place to see the *bigger picture* of what will be built before it fragments
   into plans and tasks — no "intent" home.
2. Design docs under `sources/` feed into plans, but the design is where the human
   actually decides.
3. Plans are trusted without re-reading, because they derive from an already-agreed
   design.
4. Re-approving the plan therefore feels like a wasted step.
5. Verification feels heavy for low-risk solo work — "do I even need a verifier?"
6. Multiple stacked plans spawn multiple verifiers for work that ends in **one**
   pull request.
7. Marking a plan "done" feels like unnecessary bookkeeping.
8. Forgetting to reconcile Product Knowledge after a merge silently grounds the
   next plan on stale knowledge.

Every mechanism below maps to one or more of these.

## Reading order

- `design.md` — the normative overview. A reviewer can stop here.
- `preserved-core.md` — what vNext keeps unchanged, and why each piece is
  load-bearing. Grounded in the engine read; this is what separates this study
  from the from-scratch ones.
- `intent-and-criteria.md` — Mechanism 1: the intent front door + spec adversary +
  the envelope check. (Pains 1, 3, 4.)
- `candidate-and-evidence.md` — Mechanism 2: candidate identity + mechanical
  evidence staleness. (Pains 6, 7.)
- `tiered-assurance.md` — Mechanism 3: consequence-tiered verification, with
  pairing folded in as the lowest tier. (Pains 5, and cc-pair.)
- `closed-knowledge-loop.md` — Mechanism 4: reconciliation debt that blocks
  grounding. (Pain 8.)
- `lifecycle-and-gates.md` — the resulting gate model: two human gates, inferred
  completion.
- `invariant-deltas.md` — exact `INV-*` changes: new / reworked / kept, by owner.
- `engine-and-seam.md` — new and changed engine verbs, and how the stable CLI seam
  keeps two-thirds of the periphery unchanged.
- `migration-and-build-order.md` — the strangler-fig sequence, grounded in the
  coupling map.
- `risks-and-open-questions.md` — the two crown-jewel checks and honest tradeoffs.
- `comparison-matrix.md` — vNext vs. current Context Circuit (keep / rework / add /
  remove).
