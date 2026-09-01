# Context Circuit vNext — an intent-gated, candidate-proven evolution

**Status: source material only.** This is a design study. It has no status, no
authority, and changes nothing on its own. It is not a plan, not accepted Product
Knowledge, not a release, and not permission to change the product. Read it as an
argument to test, not an answer to adopt. The working name "vNext" is a
design-line label, not a rebrand of the product; the product stays Context
Circuit unless a separate, explicit vocabulary decision says otherwise.

## What this is

A proposal for **major changes to Context Circuit's trust core**, keeping its
mechanics. Unlike a ground-up rewrite, this study was written *after a full
line-by-line read of the current core* — `wrapper/runtime/engine.sh` (1,975
lines), `wrapper/contracts/invariants.yaml`, the host adapters, and the
skill/engine coupling map. Its central finding is that **the current core is
good** — it is already a small, model-blind, deterministic mechanics kernel — so
the right move is to **evolve the policy layer and preserve the mechanics**, not
to start over.

## The thesis in one sentence

> Context Circuit today gates the *state transitions* of a plan; vNext gates the
> two things a human actually decides — **what "correct" means** (the intent) and
> **what ships** (delivery) — and makes everything between them mechanical,
> consequence-tiered, and self-invalidating.

## Where it came from

A grounded read of Context Circuit's own core, to separate **hard-won correctness
worth preserving** from **policy worth changing**. vNext keeps what the read
showed is load-bearing and adopts only the changes that fit Context Circuit's
single-host, human-in-the-loop, provider-neutral shape. It deliberately declines
three otherwise-appealing ideas — fencing tokens, a full event-sourced rewrite,
and a context/token compiler — because they solve problems this product does not
yet have (see `comparison-matrix.md`).

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

**The argument**
- `design.md` — the normative overview. A reviewer can stop here.
- `preserved-core.md` — what vNext keeps unchanged, and why each piece is
  load-bearing. Grounded in the engine read; this is what separates this study
  from a ground-up rewrite.

**The four mechanisms**
- `intent-and-criteria.md` — Mechanism 1: the intent front door + spec adversary +
  the envelope check. (Pains 1, 3, 4.)
- `candidate-and-evidence.md` — Mechanism 2: candidate identity + mechanical
  evidence staleness. (Pains 6, 7.)
- `tiered-assurance.md` — Mechanism 3: consequence-tiered verification, with
  pairing folded in as the lowest tier. (Pains 5, and cc-pair.)
- `closed-knowledge-loop.md` — Mechanism 4: reconciliation debt that blocks
  grounding. (Pain 8.)

**How it works**
- `roles-and-spawning.md` — the four agent roles and the per-tier spawning model
  (and whether Explore spawns an agent).
- `lifecycle-and-gates.md` — the resulting gate model: two human gates, inferred
  completion.
- `crown-jewels.md` — the two safety-critical checks specified: the envelope check
  and tier classification.
- `concurrency-and-candidate.md` — how the candidate composes with leases, base
  selection, and integration merge; the change-set delivery unit.
- `human-experience.md` — the conversational surface: the two gates in plain
  language, re-gating, and the transparency questions.
- `scenarios.md` — four worked end-to-end walkthroughs (also acceptance fixtures).
- `diagrams.md` — mermaid views: the full flow, roles-per-tier, object states, and
  edge cases.
- `conversations.md` — expected human↔agent dialogues per scenario (human-sim
  fixtures).

**The specification**
- `data-model.md` — the new schemas (intent-contract, candidate, acceptance, debt)
  and their fields.
- `file-and-folder-structure.md` — the installed workspace layout, `sources/`
  organization, and every folder/file convention.
- `invariant-deltas.md` — exact `INV-*` changes: new / reworked / kept, by owner.
- `engine-and-seam.md` — new and changed engine verbs, and how the stable CLI seam
  keeps two-thirds of the periphery unchanged.

**Adoption**
- `migration-and-build-order.md` — the strangler-fig sequence, grounded in the
  coupling map.
- `risks-and-open-questions.md` — the two crown-jewel checks and honest tradeoffs.
- `glossary.md` — the new vocabulary, internal→human translations, and the concept
  budget.
- `comparison-matrix.md` — vNext vs. current Context Circuit, and the ideas it
  adopts and declines.
