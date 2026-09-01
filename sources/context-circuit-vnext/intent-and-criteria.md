# Mechanism 1 — Intent front door + spec adversary + envelope

Fixes pains 1 (no bigger-picture home), 3 and 4 (plan approval feels redundant).
Adds the single highest-leverage safety mechanism in the study.

## The problem it solves

Today the design→plan path is `sources/ → context-proposals → accepted Product
Knowledge → plan`. There is no first-class object for *the thing being built* and
— confirmed from `cc-system-design` — "there is no separate design-acceptance
gate." So the human's real decision (agreeing to the design) has no home, and the
system instead asks for its approval later, on the plan, which is a mechanical
elaboration the human already trusts. The gate is on the derivative.

## The object

A first-class `intent/` object, parallel to `plans/`:

```
intent/0007-checkout-retries/
  INTENT.md        # human-facing: the bigger picture, reviewable as one thing
  contract.yaml    # the frozen decision (schema below)
  adversary.md     # the independent critique of the criteria
```

`INTENT.md` is what pain 1 asks for: a place to *see* what will be built — goal,
shape, what is deliberately out of scope — before it fragments into plans and
tasks.

### `contract.yaml` (sketch)

```yaml
schema_version: 1
intent: 0007-checkout-retries
title: Retry failed checkout charges
goal: >
  A failed card charge at checkout is retried up to 3 times with backoff before
  the order is marked failed.
non_goals:
  - Changing the payment provider
  - Any change to refund logic
constraints:
  - No new runtime dependencies
  - Existing checkout latency budget unchanged
acceptance_criteria:            # each is executable or explicitly manual
  - id: ac-1
    statement: A transient charge failure triggers up to 3 retries with backoff.
    method: test                # test | command | build | static | manual
    surface: checkout-service
  - id: ac-2
    statement: A permanent failure (card declined) is NOT retried.
    method: test
done_when: >
  ac-1 and ac-2 pass on an independent check, and a human accepts the candidate.
scope:                          # the ENVELOPE — what plans derived from this may touch
  repositories:
    - id: checkout-service
      paths: [src/checkout/, test/checkout/]
tier: standard                  # explore | standard | critical  (see M3)
status: draft                   # draft | approved   (approval freezes a digest)
contract_digest:                # set on approval; identifies the frozen criteria
```

Rules that make it load-bearing:

- **Every criterion is executable or explicitly `manual` — no third option.** This
  makes the acceptance contract *the* definition of correct, and "done" (M2) is
  computed against it, not asserted.
- **Approval freezes `contract_digest`.** Changing the criteria after approval is a
  new decision that re-enters the gate, and (via M2) voids evidence bound to the
  old contract.
- **`scope` is the envelope.** It is the machine-checkable boundary that lets plan
  approval collapse safely (below).

## The spec adversary — challenge the criteria before any code

An independent role — same "writer ≠ checker" discipline the product already
trusts for the verifier, pointed at a different input — reads `contract.yaml`
*before approval* and tries to:

- **invent a case that satisfies every criterion yet is still wrong**, and
- name the missing edge / error / security / concurrency paths the criteria do not
  cover.

Its findings land in `adversary.md`. Criteria that cannot survive are rewritten or
logged as explicit open questions before the human approves. This is the one place
the entire corpus was under-solved: verification rigor downstream all certifies
work *against* criteria, so if the criteria are wrong, the rigor faithfully
certifies the wrong thing. The adversary is cheap (it reuses the independent-actor
machinery) and it is aimed at the least-checked, highest-leverage input in the
system.

The adversary's depth scales with tier (M3): light or skipped at Explore, full
battery at Critical.

## Plans become a derivation — no second gate

Once intent is approved, `cc-plan` derives tasks from the approved `contract.yaml`
+ Product Knowledge + repository grounding, and `plan.yaml` gains a required field:

```yaml
intent: 0007-checkout-retries
```

Because the human already approved the intent, **the plan carries no second
approval gate** (pains 3, 4). The approval did not disappear — it moved upstream to
where the decision actually lives (pain 1). What replaces per-plan approval is a
mechanical check, not a human one.

## One intent, one or more plans — and structuring larger work

An intent is one *decision*; a plan is one *execution*. The relationship is **1 : N**:

- A small change is one intent → one plan.
- A larger change is one intent → **several stacked plans** (for example an API plan
  and a consumer plan) that together satisfy the intent. Each names `intent: <id>`,
  each stays inside the intent's envelope, and together they may form one change set
  at delivery (`concurrency-and-candidate.md`).

**Plan derivation is automatic.** Once the intent is approved, `cc-plan` derives the
plan(s) as the coordinator's next action — grounding, task breakdown, envelope check —
with **no separate plan-approval gate**. Because the intent's approval also authorizes
execution within the envelope (INV-EXEC-01 reworked), work then proceeds; the human
hears "here's the breakdown, building now," not "approve this plan." (A human who wants
to lock the intent without building yet can say so — that separates approval from
execution, the one case INV-EXEC-01 still distinguishes.)

**Genuinely separate decisions are separate intents.** If a request spans topics the
human would review and ship independently, each is its own intent — not one giant
intent. Intents may declare **dependencies on other intents** (mirroring inter-plan
`plan_dependencies`, INV-PLAN-05) so a multi-topic initiative can be ordered.

**The big multi-topic picture lives above intents, in `sources/system-design/`.** When
the human wants a detailed picture reviewed as a whole and split by topic, that is the
existing three-tier design layout (README → design → one file per concern). A system
design then *spawns multiple intents* — typically one per concern — each approved on
its own:

```
sources/system-design/   the big picture, by concern, reviewed whole
        │  spawns
        ▼
intent/                  one coherent decision per topic, approved individually
        │  derives
        ▼
plans/                   one or more execution plans per intent
```

So "a detailed intent I can review as a big picture" is answered two ways: the
`INTENT.md` can be as detailed as the human wants, and a genuinely large, multi-topic
picture belongs one level up in `sources/system-design/`, which feeds the intents.
`sources/` stays passive; the intent carries the authority
(`file-and-folder-structure.md`).

## The envelope check — the crown-jewel safety mechanism

Moving the gate up is safe **only** if a plan (or, later, a candidate) that exceeds
the approved scope reliably re-gates to a human. This is the price of gate-outcomes,
and vNext pays it explicitly:

- A new verb `intent-envelope-check <plan>` compares the plan's declared
  repositories and path regions against the parent intent's `scope`.
- If the plan stays within the envelope → it is authorized to execute by the
  intent's approval; no human step.
- If the plan **exceeds** the envelope (a new repository, a path region outside the
  declared scope, or a criteria change) → it is **held and re-gated to a human**,
  who either widens the intent (a new decision, re-frozen digest, adversary re-run)
  or narrows the plan.

This check must **fail upward**: an ambiguous or unresolvable comparison re-gates
rather than passes. It is deliberately small and deterministic so it can be the
most-tested component in the system (see `risks-and-open-questions.md`).

## Skills

- **New `cc-intent`** — author `INTENT.md` + `contract.yaml`, run the spec
  adversary, and perform the single upstream approval (`intent-approve`, which
  freezes `contract_digest`).
- **`cc-plan` changed** — requires an approved parent intent; sets `plan.intent`;
  runs `intent-envelope-check` as a preflight; never asks for a separate plan
  approval within the envelope.

## What stays the same

Grounding (INV-GROUND-*), the worker brief template, repository binding, and the
whole execution mechanic are untouched. The intent object sits *before* them and
feeds the criteria the candidate (M2) will be proven against.
