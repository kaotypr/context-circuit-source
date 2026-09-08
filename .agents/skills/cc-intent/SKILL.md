---
name: cc-intent
description: Author an intent — the first-class decision for one change (goal, non-goals, constraints, outcome-level acceptance criteria, a coarse optional scope, tier) — from the plain ask without reading the code, and take the single upstream human approval that freezes it and spawns the planner.
---

An intent is the one thing a human actually decides before code exists: **what
"correct" means** (Context Circuit v1.0, Mechanism 1). Approving an intent is Gate 1
— the single upstream human gate. Approval also confirms you understood the plain
ask, which is what lets the planner read the real code next (`cc-trace`); plans then
derive from the approved intent automatically, with no separate per-plan approval and
no automated scope gate (INV-INTENT-01, INV-INTENT-02, INV-APPROVE-01).

## Author the intent

Retrieve relevant Product Knowledge by the request's concepts via
`context/INDEX.md`, read only the selected units, and read only request-named
source files (a `sources/` doc or a `sources/system-design/` design may ground the
intent; `sources/` stays passive — INV-SEC-02). Allocate a stable intent id with
the runtime `intent-allocate-id` (form `i<NNN>-slug`, its own never-reused
sequence, distinct from plan ids). Before `intent-allocate-id`, resolve local
member identity with `member-band-resolve` (or `member-identity-read`). If
identity is missing (`MEMBER_IDENTITY_MISSING`), stop and use `cc-workspace` for
the one-time roster-member choice — never prompt for a block number or numeric
range. Allocation itself remains the workspace-wide sequence until band-scoped
allocation is in force.

Write two files with a strict division of audience. `intent/<id>/INTENT.md` is
**what the human reads** — plain, short; it reassures them that the blurry thing
they asked for was understood and shows what they get. Author it from
`.context-circuit/docs/templates/intent.md`, which fixes the five human-facing sections (Intention,
Expectations, The plans, How carefully this is checked, Open questions);
`.context-circuit/docs/templates/intent.example.md` shows the voice. Keep technical terms as they
are; do not rename them into something that does not mean the same thing. Do not
restate the product's mechanics or add out-of-scope, history, or
assurance-rationale sections there.

`intent/<id>/contract.yaml` is **what the agent reads** — the detailed,
machine-checkable record; the human is not expected to open it. It (schema
`.context-circuit/wrapper/contracts/schemas/intent-contract.yaml`) carries:

- `goal` — one paragraph of what a correct change achieves;
- `non_goals` and `constraints` — explicit exclusions and limits;
- `acceptance_criteria` — each `{id, statement}` at the **outcome level**: what must
  be true, in terms a human can approve ("no direct localStorage access remains").
  These are the definition of correct. Do **not** try to make them executable here —
  the runnable check that proves each one is earned against the real code by the
  planner after approval and carried into the plan (`cc-trace`). You cannot read the
  code yet, so you cannot author a machine-precise criterion; state the honest
  human-level target.
- `scope` — **coarse and optional**: `repositories: [{id, paths[]}]`, and it may be
  empty. It is not an enforced boundary — the planner reports where the change actually
  lands, the feasibility check surfaces a required change beyond it, and scope-safety
  is settled at delivery (Gate 2). Its one deterministic use is as a transparent input
  to tier signals. Note a rough boundary if the human gave one; do not invent paths.
- `tier` — `explore | standard | critical`, proposed from transparent risk signals
  (INV-ASSURE-01; the runtime `tier-classify` reports them): more than one repository, security/secrets, money,
  data migration, production/deploy, irreversibility, or novelty push higher; a single
  reversible well-covered change is Explore-eligible. Default anything uncertain to
  **standard**; a human may raise it. Fail upward. It is provisional — the planner's
  findings may raise it before plans are written.

Grounding, missing detail, and contradictions follow INV-PLAN-04: an unresolved
gap is an explicit open question, never a silently chosen decision. You draft the
intent from the plain ask and existing Product Knowledge only — you do **not** read
the codebase here; that is the planner's job after approval.

### Optional fuller write-up

A human may **request** a fuller by-concern write-up of this one change, or you
**recommend** it when the draft has several concerns — several numbered plans,
several independent outcomes, or the human already described several parts of one
change. A small single-outcome intent stays short. Skipping the write-up does not block approval and is never required.

Do **not** duplicate the three-tier rubric here. Honor the request or recommendation
by invoking `cc-system-design` to author `intent/<id>/detail/` (never `design/`).
That write-up is extra files beside the intent, not a sixth `INTENT.md` section and
not internal filenames on `INTENT.md`. Draft it without reading the codebase.
It is not part of `contract_digest` and has **no separate approval**.

A large multi-topic picture the human would review and ship as genuinely separate
decisions still lives one level up in `sources/system-design/` and spawns one
intent per concern. Intent detail is the shape of *this* one decision.

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
`intent-approve`, which flips `draft → approved`, synchronizes the human-facing
`INTENT.md` status line, and **freezes** `contract_digest`
(the frozen identity of the criteria). Approval is also your confirmation that you
understood the ask: it is what lets the planner read the real code next. Immediately
after approval, hand off to `cc-trace` — spawn one planner child per
repository in scope, collect each finding, and run the feasibility check —
**before any plan is written**. `contract.yaml` status `approved` is
Gate 1 only; it cannot skip or stand in for the planner. After a feasible planner
with no intent-level questions, update the `INTENT.md` status line with
`intent-human-status . <id> "approved, look complete, feasible"` so it cannot stay
stale, then `cc-plan` publishes the plan or plans the child already wrote **in that same turn**. After
approval the human is not asked to approve a plan.

Treat approval-to-plan as one continuation. Spawn the planner immediately after
freeze. Do not load plan templates or read the target repository first. A feasible
Standard/Critical approval publishes the child's plan files in the same turn; no
additional human action is inserted.

Changing any criteria-bearing field after approval is a new decision: it breaks
the frozen digest, so the plan's authorization fails until re-approved, and it voids
prior candidate evidence (INV-CANDIDATE-01). Take approval again on the changed
criteria.

Questions discovered after approval follow the same boundary: do not edit around an
intent-level decision or hide it in a plan; update the intent, re-enter Gate 1, and
only then derive a plan. A plan-level question does not reopen approval.

## One intent, one or more plans

An intent is one *decision*; a plan is one *execution* in **one repository**. The
relationship is 1:N — a small change is one plan, a larger change several stacked
plans that each name the intent. An intent whose scope covers two repositories
derives at least two plans, one per repository, with `plan_dependencies` when one
depends on the other. Genuinely separate decisions the human would
review and ship independently are **separate intents** (which may declare
`intent_dependencies` to be ordered), not one giant intent. A large multi-topic
picture lives one level up in `sources/system-design/` and spawns one intent per
concern. Several concerns of **one** decision may instead be written out as
optional intent detail (`intent/<id>/detail/`) so the human sees the shape before
they approve — that is not a second product-level design and not a second gate.

## Runtime actions — invoke, never read the engine

Invoke from the workspace directory (`.context-circuit/wrapper/adapters/AGENTS.md` → Runtime owns
the invoke-not-read boundary):

- `sh .context-circuit/wrapper/runtime/engine.sh member-band-resolve .` — require
  a resolved roster member before `intent-allocate-id`; never pass a numeric range.
- `sh .context-circuit/wrapper/runtime/engine.sh intent-allocate-id . <slug>` — next `i<NNN>-slug`.
- `sh .context-circuit/wrapper/runtime/engine.sh intent-validate . intent/<id>` — structure + fields.
- `sh .context-circuit/wrapper/runtime/engine.sh intent-approve . <id>` — Gate 1; freezes the digest.
- `sh .context-circuit/wrapper/runtime/engine.sh intent-human-status . <id> "<phrase>"` — rewrite
  the `INTENT.md` `_Status` line after a feasible planner; never changes
  `contract.yaml`.
- `sh .context-circuit/wrapper/runtime/engine.sh intent-archive . <id>` / `intent-restore . <id>` —
  status-blind organization (INV-ARCHIVE-01/02), the same as plans.

## Report — plain language only

Keep technical terms as they are; do not rename them into something that does
not mean the same thing. Do not dump runtime commands, digest hashes, or engine
invocations. "Here's what I understand you want to build… here's what you'll
have when it's done… this is Standard risk, so it gets an independent check.
Approve this and I'll look at the code and write the breakdown, or tell me what to change." Approving an
intent is a real decision; present it as one, not a rubber stamp.

## Boundaries

Authoring or the Gate 1 approve act never creates a plan, executes, verifies,
completes, or delivers, and never reads the codebase. It writes only under
`intent/<id>/` (`INTENT.md`, `contract.yaml`, and optional `detail/`). The planner
that reads the code runs only after approval (`cc-trace`) and may send an
intent-level question back here before a plan is published. After a feasible planner with no
intent-level questions, `cc-plan` publishes the child's plans in that same turn;
writing those plans does not start execution. Approval is conversational, never
a confirmation card or hidden token, and there is no second approval of the detail.
