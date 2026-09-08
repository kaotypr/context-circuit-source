# Risks and open questions

Reported honestly. v1.0 moves safety from "gate every transition" to a lighter model,
which concentrates the safety property into **one** automated check — consequence
tiering — supported by the delivery gate and the feasibility check. If tiering is wrong,
the design is worse than what it replaces. Everything else here is secondary.

## The safety-critical risk (and where scope-safety actually lives)

### 1. Consequence tiering (M3) — the one safety-critical check

Skipping the independent verifier at Explore is safe only if the tier is classified
deterministically and errs high. A Critical change mislabeled Explore ships with no
independent check — the single worst failure the system can have.

- **Mitigation:** tier is set at intent time from transparent, auditable risk signals,
  re-checked with the tracer's findings, and raisable by the human; the human's one
  residual duty is to glance at the chosen tier and bump it up. Tiering **fails upward**.
  The engine only enforces the floor, never lowers it. It must be the most-tested
  component in the system.
- **Residual uncertainty:** the risk-signal → tier mapping is a judgment encoded as
  rules; a novel risk not covered by a signal could under-tier. The default for anything
  uncertain must be Standard (independent verifier on), not Explore.

### 2. Scope-safety rests on the delivery gate, not an automated check

v1.0 deliberately does **not** enforce scope with an automated gate. A lay human names
little or no scope, so there is rarely a tight boundary to enforce before delivery; the
concrete scope is authorized by the human at **Gate 2**, where they see the exact diff
and repositories, and all pre-delivery work is sandboxed in isolated worktrees.

- **Mitigation:** Gate 2 is an explicit human act that already exists and is unchanged
  (INV-DELIVER-01); nothing reaches the outside world without the human seeing its actual
  scope. The feasibility check additionally surfaces, before planning, any required
  change that would *modify* a repository beyond a bound scope.
- **Residual uncertainty:** scope-safety now depends on the human actually reading the
  delivery — the same way it depends on them reading any pull request. A human who
  rubber-stamps delivery gets weaker scope protection, but that is already true of the
  delivery gate, and pre-delivery drift is wasted effort, not irreversible harm.
- **Feasibility is a judgment, not a proof.** The feasibility check is a reasoned read of
  the tracer's findings, so its reliability comes from the tracer's quality; a shallow
  the tracer could call something buildable that is not. That is a *quality* risk (wasted
  effort, a surprise at plan review), not a safety risk. And *intent* drift — work that
  stays in-scope yet does something the criteria did not anticipate — is caught by the
  acceptance criteria, not by feasibility; a weak criteria set can still ship the wrong
  thing, an inherent limit no acceptance contract fully removes.

The discipline that carries all of this: **when unsure, do more, not less.** Tiering
errs high, feasibility stops rather than guesses, and delivery is always an explicit
human act. That must be enforced, not hoped.

## Secondary risks

- **Dropping the spec adversary is a deliberate tradeoff, stated honestly.** An
  earlier version of this design included an independent spec adversary — a
  context-free role that read only the intent contract and tried to find ways to
  *satisfy every criterion and still be wrong*. That is a real defense against
  gameable or under-specified criteria, and removing it means v1.0 no longer has a
  dedicated pre-build check for that specific failure mode. **The tracer is not a
  replacement for it** — the tracer is a different, context-*full* check (what does
  the real code actually do, and what does it risk?), not an adversarial attack on
  the human's wording, and it runs only *after* the human has already committed to
  the goal at Gate 1, not before. The two catch different things: an adversary
  would have caught a criterion that is technically satisfiable while missing the
  point; the tracer catches a criterion that collides with the real codebase. This
  was a deliberate choice to cut cost (a mandatory extra role and pass on every
  Standard/Critical intent) and shift the saved effort toward grounding the plan
  itself in real code. If gameable criteria turn out to be a live problem in
  practice, reintroducing a criteria-adversary step remains the honest fallback —
  not something the tracer quietly covers for.
- **Tracing compute is spent after the human has already committed.** Because
  the tracer only spawns on approval, its cost lands on an intent the human has
  already agreed to plainly. So a deep surprise (the localStorage
  encode-vs-encrypt example in `tracing-and-grounding.md`) reaches the human only
  at **plan review**, potentially re-opening an intent they thought was settled at
  Gate 1. The cost of moving the gate upstream is that some rework now happens
  after a "yes," not before it.
- **The tracer's completeness is fallible.** A grep-based file/call-site map can
  miss a dynamically-constructed key, a reflectively-called method, or a string
  built at runtime. This is mitigated by the executable completeness check (a
  command that proves the found set is the whole set, for a "change every X"
  obligation) and by the worker's execution-time confirmation read
  (INV-GROUND-01) — but neither is a formal guarantee, just two independent nets
  under a fallible first pass.
- **Plan review is informal, not a gate.** A plan carries no approval status, so
  nothing structurally forces a human to look at it before triggering execution.
  The natural backstop is that execution is itself a separate human-triggered
  action, so nothing runs unreviewed — but a human who skims past the plan and
  immediately asks for execution gets no forced pause the way Gate 1 forces one.
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

- **How does the tracer decide a change "modifies a repo beyond a bound scope"** — the
  feasibility check's one deterministic sub-test? Recommendation: keep it a plain flag in
  the manifest (a required write to a repository/area not in the bound scope), surfaced by
  the coordinator — not a path-region algorithm, since scope-safety is settled at delivery.
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
*reduces* felt ceremony without a scope-safety regression at delivery; (b) the tracer, spawned
automatically on every approval, catches real feasibility/risk problems often
enough — surfaced at plan review or kicked back to the intent — to justify its
post-approval cost, without plan review becoming a rubber stamp that gets skimmed
past out of habit; and (c) tiering removes verifier spawns on low-risk work without
a mislabeled-Critical escaping unverified. If tier classification proves hard to make reliable,
the honest fallback is to keep more gates — the current gate-every-transition model
is the safe default this design is betting against.
