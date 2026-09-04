# Mechanism 1 — Intent front door + discovery + envelope

Fixes pains 1 (no bigger-picture home), 3 and 4 (plan approval feels redundant).
Adds the single highest-leverage safety mechanism in the study, plus the discovery
phase that makes plans detailed instead of guessed.

## The problem it solves

Today the design→plan path is `sources/ → context-proposals → accepted Product
Knowledge → plan`. There is no first-class object for *the thing being built* and
— confirmed from `cc-system-design` — "there is no separate design-acceptance
gate." So the human's real decision (agreeing to the design) has no home, and the
system instead asks for its approval later, on the plan, which is a mechanical
elaboration the human already trusts. The gate is on the derivative.

There is a second, quieter problem. A lay human cannot describe the codebase — they
state a plain goal and may know nothing about the files. For a plan to be *detailed*,
something has to read the real code. Today nothing reliably does before the plan is
written: the coordinator plans from curated `context/` plus whatever files were
named, and the first first-hand read of the repository happens at execution, in the
worker — after the plan and scope are already fixed. When a detailed plan does appear
(observed in real sessions) it is because a discovery/design phase *happened* to run
first, driven by the human, not because the lifecycle guaranteed it. v1.0 makes that
phase first-class: **discovery.**

## The object

A first-class `intent/` object, parallel to `plans/`:

```
intent/i0007-checkout-retries/
  INTENT.md        # human-facing: the bigger picture, reviewable as one thing
  contract.yaml    # the frozen human decision (schema below)
```

`INTENT.md` is what pain 1 asks for: a place to *see* what will be built — goal,
shape, what is deliberately out of scope — before it fragments into plans and tasks.

### `contract.yaml` — the frozen *human decision*, not a machine checklist

The intent is approved **before** discovery reads the code, so `contract.yaml` can
only hold what a human can actually decide without seeing the codebase: the goal, the
boundaries, and *what "done" means at the outcome level*. The executable, per-file
checks are not authored here — they cannot be, because nobody has read the code yet.
Discovery produces them later, into the plan (below, and `discovery-and-grounding.md`).

```yaml
schema_version: 2
intent: i0007-checkout-retries
title: Retry failed checkout charges
goal: >
  A failed card charge at checkout is retried a few times before the order is
  marked failed — transient failures recover, real declines do not.
non_goals:
  - Changing the payment provider
  - Any change to refund logic
constraints:
  - No new runtime dependencies
  - Existing checkout latency budget unchanged
acceptance_criteria:            # OUTCOME level — what must be true, not how it is measured
  - A transient charge failure recovers without the order being marked failed.
  - A permanent decline is never retried.
scope:                          # the ENVELOPE — coarse: candidate repositories + a boundary
  repositories:
    - id: checkout-service
      paths: [src/checkout/]    # the region discovery findings are checked against
tier: standard                  # PROVISIONAL — discovery may push it up (see M3)
status: draft                   # draft | approved   (approval freezes a digest)
contract_digest:                # set on approval; identifies the frozen decision
```

Rules that make it load-bearing:

- **Acceptance criteria are stated at the outcome level.** They say *what must be
  true*, in terms a human can approve — not the executable command that measures it.
  The runnable checks (the grep that proves no direct access remains, the test ids)
  are a **discovery output**, carried in the plan, not frozen here. The contract sets
  the goalpost; discovery builds the measuring tape.
- **Approval freezes `contract_digest`.** Changing the frozen decision after approval
  is a new decision that re-enters the gate, and (via M2) voids evidence bound to the
  old contract.
- **`scope` is the envelope — and it is coarse on purpose.** At approval the human
  can only name candidate repositories and a rough boundary; the exact paths are not
  known yet. The envelope is the boundary that **discovery findings and plans are both
  checked against** (below). Keep it as tight as the human honestly intends.
- **`tier` is provisional.** It is the human's first guess at how carefully to check.
  Discovery can raise it (a change that turns out to touch auth/secrets/data is
  Critical whatever the human first thought); the raise surfaces at plan review.

## The flow — draft, approve, discover, plan, review

The order is the point. Discovery is expensive and must aim at a *confirmed* target,
so it runs **after** the human approves — and the approval doubles as the human
confirming the coordinator understood the plain ask.

```mermaid
flowchart TD
  A["Human states a plain goal"] --> B["Coordinator drafts the intent<br/>reads existing context/ knowledge, NOT the code"]
  B --> C{"Gate 1 — human approves?<br/>(also confirms: understood correctly)"}
  C -->|approve, freeze digest| D["Discovery spawns<br/>one read-only child per repo, parallel"]
  C -->|correct me| B
  D --> E["Discovery reports back a manifest<br/>files, call-sites, risks, executable done-checks, tier signal, open questions"]
  E -->|intent itself is wrong| B
  E -->|grounded| F["Coordinator creates the plan(s)<br/>from the manifest — no plan gate"]
  F --> G["Human may review the plans<br/>(a plan is just a plan — no status)"]
  G --> H["Human asks for execution"]
```

Step by step:

1. **Draft — grounded in existing knowledge, not the code.** The coordinator
   interprets the plain ask into `INTENT.md` + `contract.yaml`: goal, non-goals,
   constraints, outcome-level acceptance, candidate repositories, a provisional tier.
   It reads the **existing `context/` Product Knowledge** as reference — what the
   project already knows about itself — but it does **not** read the codebase. On a
   brand-new project with no `context/` yet, drafting works from the human's request
   alone. **This is the intent/plan divide:** the intent is
   grounded in *existing knowledge* (context files, if any); the plan is grounded in
   *fresh discovery* (the real code, read after approval). Reading the codebase at
   intent time is discovery's job, not the intent's.
2. **Gate 1 — the one upstream human decision.** The human approves the plain intent.
   This *is* the confirmation that the coordinator understood them, so discovery never
   fires on an unconfirmed guess. Approval freezes `contract_digest`.
3. **Discovery (spawned on approval).** One read-only discovery child per repository,
   in parallel, reads the real code and **reports back** to the coordinator — a rich,
   durable manifest: the file/call-site map, integration points, concrete risks, a
   proposed task partition, the **executable "done" checks** that prove the outcome
   criteria, a tier signal, and any open questions. Discovery does *not* write plans.
   Full mechanism in `discovery-and-grounding.md`.
4. **Coordinator creates the plan(s).** From the manifest, the coordinator authors the
   plan(s) — it stays the sole author and gate-runner. The executable criteria and
   exact paths discovery found land in the plan, not back in the frozen contract.
5. **Plan review — informal, no gate.** A plan is just a plan; it has no approval
   status. The human can read the plans and course-correct, and the natural checkpoint
   is that **execution is a separate human-triggered action** — nothing runs until the
   human asks, so review always has its moment. This is where a grounded surprise
   (e.g. "you said encode, but for auth tokens that is not protection") reaches the
   human, before any code is written.

Two feedback edges keep it honest:

- **Discovery can kick back to the intent.** Usually it feeds the plan. But if it
  finds the *intent itself* is wrong — the goal is infeasible, or means something
  bigger than approved — that cannot become a plan; it re-opens the intent for a new
  decision. Cheap, because nothing has executed.
- **Tier can be raised by discovery**, surfaced at plan review, before execution.

## Plans become a derivation — no second gate

Once the intent is approved and discovery has reported, the coordinator creates the
plan(s), and `plan.yaml` gains a required field:

```yaml
intent: i0007-checkout-retries
```

Because the human already approved the intent, **the plan carries no second approval
gate** (pains 3, 4). The approval did not disappear — it moved upstream to where the
decision actually lives (pain 1). A plan is a work product, not a gated object: the
human reviews it if they wish, and triggers execution when ready.

## One intent, one or more plans — and structuring larger work

An intent is one *decision*; a plan is one *execution*. The relationship is **1 : N**:

- A small change is one intent → one plan.
- A larger change is one intent → **several stacked plans** (for example an API plan
  and a consumer plan) that together satisfy the intent. Each names `intent: <id>`,
  each stays inside the intent's envelope, and together they may form one change set
  at delivery (`concurrency-and-candidate.md`). Discovery fan-out is per repository,
  so a multi-repo intent naturally yields per-repo plans.

**Plan creation is automatic after discovery.** The human hears "here's the
breakdown" — not "approve this plan." The intent's approval authorizes execution
within the envelope (INV-EXEC-01 reworked); execution is then a human-triggered
action, the one place approval and execution are distinguished.

**Genuinely separate decisions are separate intents.** If a request spans topics the
human would review and ship independently, each is its own intent — not one giant
intent. Intents may declare **dependencies on other intents** (mirroring inter-plan
`plan_dependencies`, INV-PLAN-05) so a multi-topic initiative can be ordered.

**The big multi-topic picture lives above intents, in `sources/system-design/`.** A
system design spawns multiple intents — typically one per concern — each approved on
its own:

```
sources/system-design/   the big picture, by concern, reviewed whole
        │  spawns
        ▼
intent/                  one coherent decision per topic, approved individually
        │  discovery grounds, then derives
        ▼
plans/                   one or more execution plans per intent
```

`sources/` stays passive; the intent carries the authority
(`file-and-folder-structure.md`).

## The envelope check — the crown-jewel safety mechanism

Moving the gate up is safe **only** if work that exceeds the approved scope reliably
re-gates to a human. In v1.0 the envelope is checked at **two** points, not one:

- **Against discovery findings.** When discovery reads the code and finds the change
  actually reaches *outside* the approved boundary (a path region the human did not
  envisage, or another repository entirely), that is **held and re-gated** before any
  plan is written: "this is bigger than you approved."
- **Against each plan.** The verb `intent-envelope-check <plan>` compares the plan's
  declared repositories and path regions against the parent intent's `scope`. Within
  → authorized by the intent's approval, no human step. Exceeds (a new repository, an
  out-of-scope region, or a change to the frozen decision) → **held and re-gated**,
  and the human either widens the intent (a new decision, re-frozen digest) or narrows
  the plan.

This check must **fail upward**: an ambiguous or unresolvable comparison re-gates
rather than passes. It is deliberately small and deterministic so it can be the
most-tested component in the system (see `risks-and-open-questions.md`).

## Skills

- **New `cc-intent`** — author `INTENT.md` + `contract.yaml` (the plain human
  decision) and perform the single upstream approval (`intent-approve`, which freezes
  `contract_digest`). It reads the existing `context/` Product Knowledge as reference,
  but not the codebase — that is discovery's job, after approval.
- **New `cc-discover`** — on approval, spawn the per-repo discovery children, collect
  their manifests, and check the findings against the envelope
  (`discovery-and-grounding.md`).
- **`cc-plan` changed** — creates the plan(s) from the discovery manifest (not from a
  blind read of `context/`); sets `plan.intent`; runs `intent-envelope-check`; never
  asks for a separate plan approval within the envelope.

## What stays the same

Grounding (INV-GROUND-*), the worker brief template, repository binding, and the
whole execution mechanic are untouched. The intent object sits *before* them; the
worker's execution-time repository read (INV-GROUND-01) now *confirms* discovery's
findings rather than being the first contact with the code. The intent feeds the
outcome criteria the candidate (M2) is proven against; discovery feeds the executable
checks that measure them.
