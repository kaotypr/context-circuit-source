# Risks and open questions

Reported honestly. v1.0 moves safety from "gate every transition" to "gate the two
that matter," which concentrates the entire safety property into two deterministic
checks. If either is wrong, the design is worse than what it replaces. Everything
else here is secondary.

## The two crown-jewel risks

### 1. Intent-envelope drift detection (M1)

Moving the gate to intent is safe only if a plan or candidate that exceeds the
approved `scope` reliably re-gates to a human. If the envelope check is too loose
(or a path region comparison is subtly wrong), scope creep ships under an approval
the human never gave for that scope — exactly the accident the system exists to
prevent.

- **Mitigation:** the check is small, deterministic, and **fails upward** — any
  ambiguity (unrecognized path, repository not named, contract digest mismatch)
  re-gates rather than passes. It must be the most-tested component in the system.
- **Residual uncertainty:** "within scope" is a path/repo comparison, which is
  mechanical, but *intent* drift (the work technically stays in the paths yet does
  something the criteria did not anticipate) is not caught by a path check — that is
  what the spec adversary and the acceptance criteria are for. The envelope guards
  *where*; the criteria guard *what*. A weak criteria set plus in-scope paths can
  still ship the wrong thing. This residual is inherent — no acceptance contract
  fully anticipates intent.

### 2. Consequence tiering (M3)

Skipping the independent verifier at Explore is safe only if the tier is classified
deterministically and errs high. A Critical change mislabeled Explore ships with no
independent check.

- **Mitigation:** tier is set at intent time from transparent, auditable risk
  signals and is raisable by the human; the human's one residual duty is to glance
  at the chosen tier and bump it up. Tiering **fails upward**. The engine only
  enforces the floor, never lowers it.
- **Residual uncertainty:** the risk-signal → tier mapping is a judgment encoded as
  rules; a novel risk not covered by a signal could under-tier. The default for
  anything uncertain must be Standard (independent verifier on), not Explore.

Both crown jewels share a discipline: **when unsure, do more, not less.** That is
the whole safety argument for gate-outcomes, and it must be enforced, not hoped.

## Secondary risks

- **The spec adversary can become theatre.** A noisy adversary that fabricates
  implausible edge cases trains people to ignore it — and then it is the new
  rubber stamp. It must be calibrated to *real, reachable* failure, and its depth
  tiered (light at Explore, full at Critical). Measure how often its findings change
  a criterion; if that rate is low, it is miscalibrated.
- **Inferred completion can ship something a human would have paused on.** Removing
  the explicit done-flip at Explore/Standard moves a decision onto the acceptance +
  delivery signals. Mitigation: delivery (Gate 2) is still an explicit human act, so
  nothing reaches the outside world without a human; "done" is only the *internal*
  status. Critical keeps explicit completion.
- **Knowledge-debt blocking can nag or stall.** Too-aggressive blocking makes the
  next plan wait on reconciliation the human considers trivial; too-loose warning
  reproduces pain 8. Mitigation: block at Standard/Critical, warn at Explore, and
  make "no update needed" a cheap first-class resolution.
- **Candidate churn under active development.** Every commit voids prior evidence,
  which is correct but could feel like re-verification thrash on a fast-moving
  branch. Mitigation: verification is normally cheap because acceptance criteria are
  commands (re-running is not a model call); only human acceptance genuinely needs
  redoing, and that is the point.
- **Two-tier vocabulary load.** Intent vs plan vs candidate vs tier adds concepts.
  Mitigation: hold a hard concept budget ("delete one to add one"); several of these
  *remove* concepts the user juggles today (plan-approval, mark-done, gather-context
  chore). Net concept count should not rise.

## Open questions

- **Where does the envelope check draw the line on path regions** — exact prefix
  match (like the lease overlap rule) or something looser? Recommendation: reuse the
  existing `cc_region_overlap` semantics so envelope and lease share one definition.
- **Should Explore work ever auto-produce a candidate** (for knowledge feedback)
  without promotion, or only on promotion? Recommendation: only on promotion, to
  keep the fast path truly recordless.
- **Does promotion author a plan, or a plan-less candidate?** When an Explore session
  is promoted, its already-made commits need a candidate — but does a `plans/<id>`
  get created (a lightweight plan of record, authored partly retroactively), or does
  the execution/candidate form directly from the pairing branch with no plan?
  Recommendation: author a plan of record, so "every candidate has a plan behind it"
  stays a uniform rule (see `lifecycle-and-gates.md`). The invariant either way: no
  plan file for un-promoted Explore work.
- **Tier defaults per risk signal** — the exact signal→tier table. Recommendation:
  ship conservative (bias to Standard) and tune from real mis-tier incidents, the
  way `role-tiering` is host-local and adjustable.
- **Does the production feedback edge (M4 optional) belong in core at all?**
  Recommendation: no — keep it an opt-in extension over the publication inbound
  path, never an automatic state mutation.
- **Team acceptance semantics** — is one human acceptance record enough, or does a
  team want N? Recommendation: one by default, with the record naming the accepter;
  quorum is a later policy, not core.

## What would falsify the design

v1.0 is worth adopting only if, in real use: (a) moving the gate to intent
*reduces* felt ceremony without a scope-drift regression; (b) the spec adversary
changes criteria often enough to be worth its cost; and (c) tiering removes verifier
spawns on low-risk work without a mislabeled-Critical escaping unverified. If the
envelope or tiering checks prove hard to make reliable, the honest fallback is to
keep more gates — the current gate-every-transition model is the safe default this
design is betting against.
