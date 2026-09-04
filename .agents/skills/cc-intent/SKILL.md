---
name: cc-intent
description: Author an intent — the first-class decision for one change (goal, non-goals, constraints, outcome-level acceptance criteria, a coarse optional scope, tier) — from the plain ask without reading the code, and take the single upstream human approval that freezes it and spawns the tracer.
---

An intent is the one thing a human actually decides before code exists: **what
"correct" means** (Context Circuit v1.0, Mechanism 1). Approving an intent is Gate 1
— the single upstream human gate. Approval also confirms you understood the plain
ask, which is what lets the tracer read the real code next (`cc-trace`); plans then
derive from the approved intent automatically, with no separate per-plan approval and
no automated scope gate (INV-INTENT-01, INV-INTENT-02, INV-APPROVE-01).

## Author the intent

Retrieve relevant Product Knowledge by the request's concepts via
`context/INDEX.md`, read only the selected units, and read only request-named
source files (a `sources/` doc or a `sources/system-design/` design may ground the
intent; `sources/` stays passive — INV-SEC-02). Allocate a stable intent id with
the runtime `intent-allocate-id` (form `i<NNNN>-slug`, its own never-reused
sequence, distinct from plan ids).

Write two files with a strict division of audience. `intent/<id>/INTENT.md` is
**what the human reads** — plain, short, no jargon; it reassures them that the
blurry thing they asked for was understood and shows what they get. Author it from
`docs/templates/intent.md`, which fixes the five human-facing sections (Intention,
Expectations, The plans, How carefully this is checked, Open questions) and the
plain-language rules; `docs/templates/intent.example.md` shows the voice. Do not
restate the product's mechanics or add out-of-scope, history, or
assurance-rationale sections there.

`intent/<id>/contract.yaml` is **what the agent reads** — the detailed,
machine-checkable record; the human is not expected to open it. It (schema
`wrapper/contracts/schemas/intent-contract.yaml`) carries:

- `goal` — one paragraph of what a correct change achieves;
- `non_goals` and `constraints` — explicit exclusions and limits;
- `acceptance_criteria` — each `{id, statement}` at the **outcome level**: what must
  be true, in terms a human can approve ("no direct localStorage access remains").
  These are the definition of correct. Do **not** try to make them executable here —
  the runnable check that proves each one is earned against the real code by the
  tracer after approval and carried into the plan (`cc-trace`). You cannot read the
  code yet, so you cannot author a machine-precise criterion; state the honest
  human-level target.
- `scope` — **coarse and optional**: `repositories: [{id, paths[]}]`, and it may be
  empty. It is not an enforced boundary — the tracer reports where the change actually
  lands, the feasibility check surfaces a required change beyond it, and scope-safety
  is settled at delivery (Gate 2). Its one deterministic use is as a transparent input
  to tier signals. Note a rough boundary if the human gave one; do not invent paths.
- `tier` — `explore | standard | critical`, proposed from transparent risk signals
  (INV-ASSURE-01; the runtime `tier-classify` reports them): more than one repository, security/secrets, money,
  data migration, production/deploy, irreversibility, or novelty push higher; a single
  reversible well-covered change is Explore-eligible. Default anything uncertain to
  **standard**; a human may raise it. Fail upward. It is provisional — the tracer's
  findings may raise it before plans are written.

Grounding, missing detail, and contradictions follow INV-PLAN-04: an unresolved
gap is an explicit open question, never a silently chosen decision. You draft the
intent from the plain ask and existing Product Knowledge only — you do **not** read
the codebase here; that is the tracer's job after approval.

### Keep questions across phases

The `Open questions` section is **phase-aware**, not a one-time draft field:

- At draft time, record every known unresolved human decision from the plain ask
  and Product Knowledge. `No known unresolved human decisions at draft time` is
  precise; it must not be used to claim that later code reading cannot surface a
  question.
- If the plain ask already answers a question (for example, it says to update
  every current reference), apply that decision and do not ask the human to repeat
  it or narrow the scope.
- After tracing, revisit the section and classify every newly surfaced question.
  A question that changes the goal, non-goals, constraints, acceptance criteria,
  scope, tier, authority, or lifecycle semantics is **intent-level**: update the
  intent and contract, re-enter Gate 1, and preserve the question and its answer
  in the human record. A question about implementation shape only is **plan-level**:
  carry it into the trace and plan without inventing a second human gate.
- A plan must never silently answer an intent-level question. If one remains open,
  stop before planning and return to the intent route. The approval gate is repeated
  only when the approved decision changes; plan-level resolution remains mechanical.

## Take the approval (Gate 1)

Present the intent as one readable thing: the goal in plain language, the non-goals,
the repositories it will likely touch, any open questions the human must settle, and
the tier with a one-line "why this tier". Ask for a single conversational decision —
no confirmation card, no token (INV-APPROVE-01). On a yes, run the runtime
`intent-approve`, which flips `draft → approved` and **freezes** `contract_digest`
(the frozen identity of the criteria). Approval is also your confirmation that you
understood the ask: it is what lets the tracer read the real code next. Immediately
after approval, hand off to `cc-trace` — spawn one read-only tracer child per
repository in scope, collect the manifests, and run the feasibility check on the
findings — before any plan is written. After approval the human is not asked to
approve a plan; `cc-plan` derives it from the trace manifest.

Changing any criteria-bearing field after approval is a new decision: it breaks
the frozen digest, so the plan's authorization fails until re-approved, and it voids
prior candidate evidence (INV-CANDIDATE-01). Take approval again on the changed
criteria.

Questions discovered after approval follow the same boundary: do not edit around an
intent-level decision or hide it in a plan; update the intent, re-enter Gate 1, and
only then derive a plan. A plan-level question does not reopen approval.

## One intent, one or more plans

An intent is one *decision*; a plan is one *execution*. The relationship is 1:N — a
small change is one plan, a larger change several stacked plans that each name the
intent. Genuinely separate decisions the human would
review and ship independently are **separate intents** (which may declare
`intent_dependencies` to be ordered), not one giant intent. A large multi-topic
picture lives one level up in `sources/system-design/` and spawns one intent per
concern.

## Runtime actions — invoke, never read the engine

Invoke from the workspace directory (`wrapper/adapters/AGENTS.md` → Runtime owns
the invoke-not-read boundary):

- `sh wrapper/runtime/engine.sh intent-allocate-id . <slug>` — next `i<NNNN>-slug`.
- `sh wrapper/runtime/engine.sh intent-validate . intent/<id>` — structure + fields.
- `sh wrapper/runtime/engine.sh intent-approve . <id>` — Gate 1; freezes the digest.
- `sh wrapper/runtime/engine.sh intent-archive . <id>` / `intent-restore . <id>` —
  status-blind organization (INV-ARCHIVE-01/02), the same as plans.

## Report — plain language only

Say the effect, never the mechanism (`docs/terminology.md`). "Here's what I
understand you want to build… here's what you'll have when it's done… this is
Standard risk, so it gets an independent check. Approve this and I'll look at the
real code and build it, or tell me what to change." Never expose the intent id,
contract file, digest, runtime commands, the tracer, or the feasibility check in
normal conversation — the human hears the goal and any open questions, not the
machinery. Approving an intent is a real decision; present it as one, not a rubber
stamp.

## Boundaries

Authoring or approving an intent never creates a plan, executes, verifies,
completes, or delivers, and never reads the codebase. It writes only under
`intent/<id>/`. The tracer that reads the code runs only after approval (`cc-trace`)
and may send an intent-level question back here before planning. Approval is
conversational, never a confirmation card or hidden token.
