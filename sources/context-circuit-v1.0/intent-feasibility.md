# The feasibility check

After the tracer reads the real code, the coordinator makes two decisions before it writes
any plan: **the feasibility check** (this file) and **tier classification**
(`intent-tier.md`). This one is a **quality** gate, not a safety gate — it stops the
system building something impossible or guessing blindly, but it is not what keeps a
dangerous change safe (that is tier). It **fails upward**: when it cannot conclude the
change is buildable, it stops and asks a human rather than proceed.

**Purpose.** Once the tracer has read the real code, decide whether the approved intent is
actually buildable and what it will take — before any plan is written. It is the first
point where first-hand knowledge of the code meets the human's goal.

## Why feasibility, not scope

The check this replaces used to ask "did the work stay inside the paths the human drew?"
But a lay human draws almost no paths: they state a goal, not file paths, so an intent's
scope is coarse or empty and a scope check has nothing to bite on — it almost always
passes. The question the tracer can actually answer is the useful one: *is this possible,
and what does it cost?*

Scope-safety — the guarantee that nothing ships beyond what the human meant — is settled
later, at delivery. **Gate 2** already makes the human see and authorize the exact diff
and repositories, and everything before delivery runs in isolated worktrees, so work that
drifts wider than expected is wasted effort caught at plan review or at Gate 2 — never
irreversible harm. This is a deliberate trade (`design.md`): v1.0 does not add an
automated scope gate; it leans on the delivery gate that already exists.

## When it runs

Once, after the tracer reports its manifest, before any plan exists. It is a coordinator
step in `cc-trace`, not an engine verb — feasibility is a reasoned judgment over
the tracer's findings, so there is no per-plan scope-containment verb to run
(`engine-and-seam.md`).

## Inputs

- The tracer's manifest: the file/call-site map, integration points, concrete risks, what
  the change would actually have to touch, and any blockers.
- The intent: its goal and outcome-level acceptance criteria, plus a bound repository or
  folder **if the human named one** (often there is none).

## Outcomes

- **Feasible** → the coordinator sets the tier (`intent-tier.md`), creates the plan(s)
  from the manifest, and auto-commits them. No human step.
- **Not feasible** → **stop.** The coordinator explains the cause and the specific blocker
  the tracer found, suggests what would make it work, and hands the decision back to the
  human. Cheap, because nothing has executed.
- **Feasible, but it must change a repository or area beyond a bound scope** → surfaced in
  the same report, as a question: "this is doable, but it also needs to change
  `payments-lib`, which you scoped out — include it, or should I find another way?" The
  trigger is the tracer finding the work must **modify** an out-of-scope repository or area;
  merely reading or calling one does not trigger it. When the human bound no repository,
  there is nothing to reach past — the tracer simply reports where the change lands.

## What it does NOT do (stated honestly)

- It does not guarantee the work stays within an approved scope — that is Gate 2's job at
  delivery, not an automated check here.
- It does not check that the code does the *right* thing — that is the acceptance
  criteria, sharpened by the tracer into executable checks (`tracing-and-grounding.md`).

## On determinism

Feasibility is a reasoned judgment over the tracer's findings, not a small provable rule,
so its reliability comes from the quality of tracing. The one deterministic part is the
bound-scope reach test above — *does the tracer report a required change to a repository or
area outside a bound scope?* — which the manifest carries as a plain flag for the
coordinator to surface. Everything safety-critical that used to lean on a deterministic
path comparison now lives in tier classification (`intent-tier.md`) and at Gate 2.
