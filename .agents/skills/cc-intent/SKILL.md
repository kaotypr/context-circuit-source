---
name: cc-intent
description: Author an intent — the first-class decision for one change (goal, non-goals, constraints, acceptance criteria, scope envelope, tier) — have an independent spec adversary attack the criteria, and take the single upstream human approval that freezes it.
---

An intent is the one thing a human actually decides before code exists: **what
"correct" means** and **what scope is in bounds** (Context Circuit v1.0,
Mechanism 1). Approving an intent is Gate 1 — the single upstream human gate.
Plans then derive from it automatically within its scope envelope; there is no
separate per-plan approval (INV-INTENT-01, INV-INTENT-02, INV-APPROVE-01).

## Author the intent

Retrieve relevant Product Knowledge by the request's concepts via
`context/INDEX.md`, read only the selected units, and read only request-named
source files (a `sources/` doc or a `sources/system-design/` design may ground the
intent; `sources/` stays passive — INV-SEC-02). Allocate a stable intent id with
the runtime `intent-allocate-id` (form `i<NNNN>-slug`, its own never-reused
sequence, distinct from plan ids).

Write `intent/<id>/INTENT.md` (the human-facing bigger picture — goal, shape,
what is deliberately out of scope, reviewable as one thing) and
`intent/<id>/contract.yaml` (the machine record, schema
`wrapper/contracts/schemas/intent-contract.yaml`) with:

- `goal` — one paragraph of what a correct change achieves;
- `non_goals` and `constraints` — explicit exclusions and limits;
- `acceptance_criteria` — each `{id, statement, method, surface}` where `method`
  is `test | command | build | static | manual`. **Every criterion is executable
  or explicitly `manual` — no third option.** These are the definition of correct;
  "done" is computed against them, never asserted.
- `scope` — the **envelope**: `repositories: [{id, paths[]}]`. This is the
  machine-checkable boundary a derived plan may not exceed. Keep it as tight as the
  change honestly needs.
- `tier` — `explore | standard | critical`, proposed from transparent risk signals
  (`crown-jewels.md` / `docs/tiered-assurance.md`): more than one repository,
  security/secrets, money, data migration, production/deploy, irreversibility, or
  novelty push higher; a single reversible well-covered change is Explore-eligible.
  Default anything uncertain to **standard**; a human may raise it. Fail upward.

Grounding, missing detail, and contradictions follow INV-PLAN-04: an unresolved
gap is an explicit open question, never a silently chosen decision.

## Run the spec adversary

Before approval, spawn an independent **spec adversary** child (`agents/spec-adversary.md`)
with only the contract as input — never the implementation, which does not exist
yet. It tries to (a) satisfy every criterion and still be wrong and (b) name the
missing edge / error / security / concurrency / data-loss paths. Record its
findings in `intent/<id>/adversary.md` as `{severity, statement, suggested
criterion}` plus a verdict `criteria_sound: yes | needs-work`; the record also
includes the exact `contract_digest` of the criteria it challenged. Fold surviving
findings into new or revised criteria, or log them as explicit open questions
*before* the human approves. The adversary's depth scales with tier: light or a
single inline pass at Explore, a full battery at Critical. If the host cannot
create the adversary child, report `host-blocked` and present the criteria to the
human as unchallenged — never fake the pass. The runtime refuses approval when
`adversary.md` is absent, stale, or not sound.

## Take the approval (Gate 1)

Present the intent as one readable thing: the goal in plain language, the non-goals,
the repositories it will touch, the adversary's findings and what you changed, and
the tier with a one-line "why this tier". Ask for a single conversational decision
— no confirmation card, no token (INV-APPROVE-01). On a yes, run the runtime
`intent-approve`, which flips `draft → approved` and **freezes** `contract_digest`
(the frozen identity of the criteria). After approval the human is not asked to
approve a plan; `cc-plan` derives it.

Changing any criteria-bearing field after approval is a new decision: it breaks
the frozen digest, re-gates the envelope, and voids prior candidate evidence
(INV-CANDIDATE-01). Re-run the adversary on the changed criteria and take approval
again.

## One intent, one or more plans

An intent is one *decision*; a plan is one *execution*. The relationship is 1:N — a
small change is one plan, a larger change several stacked plans that each name the
intent and stay inside its envelope. Genuinely separate decisions the human would
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
understand you want to build… I had this checked adversarially first… this is
Standard risk, so it gets an independent check. Approve this and I'll build it, or
tell me what to change." Never expose the intent id, contract file, digest, or
runtime commands in normal conversation. Approving an intent is a real decision;
present it as one, not a rubber stamp.

## Boundaries

Authoring or approving an intent never creates a plan, executes, verifies,
completes, or delivers. It writes only under `intent/<id>/`. The spec adversary is
read-only over the contract and never edits it. Approval is conversational, never a
confirmation card or hidden token.
