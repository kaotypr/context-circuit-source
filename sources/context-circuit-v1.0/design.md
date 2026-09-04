# v1.0 design — normative overview

A reviewer can stop at this file. Everything else is depth.

## The one structural error v1.0 corrects

Context Circuit spends its heaviest machinery gating the **state transitions of a
plan**: approval (draft→approved), execution (separate authorization), completion
(explicit human), delivery (separate action). Each is a real, structurally
enforced gate. But the human's actual decision — "is this the right thing, built
correctly?" — happens *upstream*, when they agree to a design or intent. By the
time a plan exists, the meaningful decision is already made; the plan is a
mechanical elaboration of it. So Context Circuit gates the **derivative** (the
plan, the done-flip) and under-builds the **source** (the intent). The result is
ceremony that trains habituated approval — and habituated approval is worse than
no approval, because it looks like safety while being reflex.

v1.0 moves the gate to where the decision actually lives, and makes the rest
mechanical.

## The shape: two gates, four mechanisms, one preserved kernel

**Two human gates, and only two:**

1. **Intent** — the human approves *what "correct" means*: the goal, the
   non-goals, the constraints, the outcome-level acceptance criteria, and a coarse,
   optional scope. This is approved once, up front; the approval also confirms the
   coordinator understood the plain ask, since the tracer — the read of the real
   code — only spawns once that target is confirmed.
2. **Delivery** — the human authorizes *the irreversible act* (pull request,
   merge, push, deploy). Never implied by verification or completion.

Everything between those two gates is **mechanical, consequence-tiered, and
self-invalidating**, realized by four mechanisms:

- **M1 · Intent + tracing + feasibility (front door).** A first-class `intent/`
  object holds the bigger picture and a frozen acceptance contract. Approval spawns
  **the tracer** — one read-only child per repository, in parallel — which reads the
  real code and reports back a manifest (files, call-sites, risks, executable
  done-checks, a tier signal) that the coordinator plans from. The coordinator then
  runs a **feasibility check** on the findings: buildable → plans are derived with no
  second gate; not buildable → held and returned to the human; buildable but reaching
  beyond a bound scope → surfaced as a question before any plan is written. Scope-safety
  itself is settled at delivery (Gate 2), not by an automated gate here.
- **M2 · Candidate + evidence.** Evidence (independent check, human acceptance)
  binds to a **candidate** — a digest over the exact per-repository commit map,
  the selected bases, and the frozen contract revision. Any new commit or criteria
  change yields a new candidate and auto-voids prior evidence and acceptance.
  Several stacked plans that converge to one pull request become one candidate →
  one verification → one acceptance.
- **M3 · Tiered assurance.** Verification effort scales to consequence
  (Explore / Standard / Critical). Explore is human-supervised with no separate
  verifier spawn; Standard and Critical require an independent read-only verifier
  bound to the current candidate. "Never self-verify" stays absolute. Direct
  collaboration (`cc-pair`) is folded in as the Explore tier — no longer a separate
  mode, but the promotable bottom of one ladder.
- **M4 · Closed knowledge loop.** Completion or delivery emits a
  reconciliation-debt marker; the next plan's grounding preflight blocks or loudly
  warns while merged work remains unreconciled. Forgetting to update Product
  Knowledge becomes structurally hard.

**One preserved kernel.** The deterministic mechanics — leases, base selection,
worktree isolation, verifier read-only enforcement, delivery drift guard,
preserve-on-failure, atomic records — are kept as-is. They are already correct and
subtle; v1.0 rides on top of them through the same stable engine-action seam. See
`preserved-core.md`.

## Why not a from-scratch rebuild

The current `engine.sh` is not a tangle. It is disciplined, correctness-dense
POSIX shell whose own boundary (INV-RUNTIME-01) already matches the
"small, model-blind, deterministic mechanics kernel with all intelligence at
swappable edges" that a trust core wants. The changes v1.0 wants live in the
**invariants + skills** layer plus a
handful of new verbs — not in the kernel. A from-scratch design would spend most
of its effort re-earning correctness the kernel already has, and would arrive more
complex, not easier. Ease comes from **subtraction and moving the gate**, which
this evolution delivers without discarding the mechanics.

## The gate model, before and after

| Decision | Context Circuit today | v1.0 |
| --- | --- | --- |
| What "correct" means | one unassisted human glance folded into plan approval | **explicit intent gate + post-approval tracing** |
| Plan is ready | explicit human approval (INV-APPROVE-01) | automatic derivation once the intent is approved; scope settled at delivery |
| Execution may start | separate authorization (INV-EXEC-01) | authorized by the approved intent |
| Work is checked | one independent verifier, always (INV-VERIFY-01) | independent verifier at Standard/Critical; human-supervised at Explore |
| Plan is done | explicit human flip (INV-COMPLETE-01) | inferred from candidate acceptance + delivery at low tier; explicit at Critical |
| Knowledge is current | manual reconciliation, forgettable | reconciliation debt blocks the next grounding |
| It ships | separate human action (INV-DELIVER-01) | **unchanged — the second and final gate** |

## The one thing that must be right (and where scope-safety lives)

v1.0 moves safety from "gate every transition" to a lighter model, which concentrates
the safety property into **one** automated check:

- **Consequence tiering (`intent-tier.md`)** — skipping the independent verifier at
  Explore is safe only if the tier is classified deterministically and **fails upward**
  (when unsure, tier higher). A mis-low tier ships an unverified risky change — the
  single worst failure the system can have — so tiering is the smallest, most-tested
  component in the system.

Two supporting pieces carry the rest. The **feasibility check** (`intent-feasibility.md`)
is a quality gate — don't build the impossible, don't guess — not a safety gate.
**Scope-safety** is carried by **Gate 2 (delivery)**, a human gate that already exists:
the human sees and authorizes the exact diff and repositories, and everything before
delivery is sandboxed in isolated worktrees, so pre-delivery drift is wasted effort, not
irreversible harm. This is the deliberate trade — v1.0 does not add an automated scope
gate; it leans on the delivery gate.

## What v1.0 deliberately does NOT adopt

v1.0 declines three of the ideas it weighed:

- **Fencing tokens on leases** — cheap to add later; single-host CC does not hit
  the zombie-writer race today.
- **A full event-sourced engine rewrite** — the append-only-record discipline CC
  already has captures most of the benefit; a log-as-truth rewrite trades a real
  schema-evolution cost for auditability that is not the current bottleneck.
- **A deterministic context/token compiler** — worthwhile as invisible present
  quality, but not required to fix the eight pains; deferred.

Each is a defensible future addition, not part of the core evolution.
