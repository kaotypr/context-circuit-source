---
kind: domain
status: accepted
title: Intent and Gate 1
slug: intent
owners: []
sources: []
source_revisions:
  - wrapper: HEAD
    commit: current
    basis: current-wrapper
generated_at: 2026-09-04T00:00:00Z
review_date: 2026-12-04
freshness: accepted-from-current-wrapper
assumptions:
  - Intent approval is the single upstream human gate (Gate 1); delivery is the only other gate (Gate 2).
unknowns: []
contradictions: []
acceptance:
  state: accepted
  accepted_at: 2026-09-04
  accepted_by: maintainer
workflows:
  - .context-circuit/docs/planning.md
---

# Intent and Gate 1

## Summary

Every writing request is anchored to a first-class **intent** before any plan is
written, and the human approves that intent once — **Gate 1**, the single upstream
human gate (INV-APPROVE-01). The intent is what the human reads and approves; it
holds the goal, non-goals, constraints, **outcome-level** acceptance criteria, a
coarse and **optional** scope, and a provisional consequence tier (INV-INTENT-01).
Route "I want to build/change …", intent drafting, and intent approval here. Owned
by the `cc-intent` skill and the intent contract schema.

## Scope

Inside: the two-file intent (`INTENT.md` for the human, `contract.yaml` for the
machine record), what belongs in each, the outcome-criteria and coarse-scope
shape, the provisional tier, and the approval act that freezes the contract.

Outside: what happens *after* approval — reading the real code and judging
buildability ([tracing and feasibility](../tracing/README.md)), how a plan then
becomes authorized ([plan authorization](../plan-authorization/README.md)), the
tier ladder itself ([assurance](../assurance/README.md)), and delivery/Gate 2
([delivery](../delivery/README.md)).

## Behavior

An intent is drafted by the coordinator **without reading the codebase** — it
grounds only in existing `context/` Product Knowledge, reflecting back the plain
ask (INV-INTENT-01). Two artifacts hold it:

- **`INTENT.md`** is human-facing: plain, short, jargon-free — no ids, file names,
  digests, branch/model names, or runtime commands. It uses exactly five sections
  (Intention, Expectations, The plans, How carefully this is checked, Open
  questions) per `.context-circuit/docs/templates/intent.md`. The human is not expected to open the
  machine record. A fuller by-concern write-up, when present, lives as extra files
  under `intent/<id>/detail/` — not as a sixth section.
- **`contract.yaml`** is the machine record: goal, non-goals, constraints,
  outcome-level `acceptance_criteria`, a coarse optional `scope`, a provisional
  `tier`, and (after approval) a frozen `contract_digest`.
- **Optional `intent/<id>/detail/`** is a fuller by-concern write-up of this one
  intent, authored by `cc-system-design` when the human requests it or when
  `cc-intent` recommends it (several concerns). Skipping it does not block Gate 1.
  It is not part of `contract_digest`, has no status, and has no separate approval.

Acceptance criteria are stated at the **outcome** level — what must be true, not
how to test it. They are **not** frozen as executable checks here; the runnable
check that proves each criterion against the real code is *earned after approval*
by the planner ([tracing and feasibility](../tracing/README.md), INV-INTENT-02),
never authored on the intent.

Scope is **coarse and optional**. It is a hint, not a fence: there is no
automated scope gate anywhere in the pre-delivery path. Concrete scope-safety is
settled at **delivery (Gate 2)**, where the human sees and authorizes the exact
diff and repositories ([delivery](../delivery/README.md), INV-DELIVER-01).

**Approval is Gate 1** (INV-APPROVE-01): an explicit conversational human act that
moves the contract `draft → approved`, **freezes `contract_digest`** — the frozen
identity of the criteria the change is later proven against (INV-CANDIDATE-01) —
and confirms the coordinator understood the plain ask, which is what lets the
planner read the real code next (INV-INTENT-02). There is no confirmation card and
no hidden confirmation token; a vague "yes" is not an approval. Approving an intent
and asking to build in one turn is honored as two sequential explicit actions; a
human who wants to lock an intent without building yet may separate the two.

Approval does **not** trigger execution and carries **no** second gate for the
plans that derive from it — plan readiness after a feasible planner is automatic
(INV-INTENT-02), not a second human approval. Optional intent detail likewise has
**no separate approval**. That automatic derivation happens **in the same turn**
after a feasible look with no unresolved intent-level question; `contract.yaml`
status `approved` does not skip the planner. After a feasible planner the
human-facing `INTENT.md` status line is updated so it cannot stay stale;
`contract.yaml` remains the approval authority.

The human-facing **Open questions** section is phase-aware. At draft time it records
known unresolved decisions; after the planner look, newly discovered questions are classified
as intent-level, plan-level, or already answered. An intent-level question sends the
work back through Gate 1 when the approved decision changes. A plan-level question
is carried into planning, and an answer already present in the plain request is
applied without asking again.

## Interfaces

- Human request: "I want to build/change …" → a drafted intent; "write this change out by topic" → optional detail; "approve intent `<id>`" → Gate 1
- Human-facing file: `intent/<id>/INTENT.md` (five sections, plain language)
- Machine record: `intent/<id>/contract.yaml` (`schema_version` 2; `status` `draft`/`approved`; frozen `contract_digest`)
- Optional detail: `intent/<id>/detail/` (three-tier write-up; not digested; not a second gate)
- Intent ids: the distinct form `i<NNN>-slug`, their own never-reused sequence
  with a hard `i999` ceiling; allocation fails rather than inventing a fourth digit

## Data

The intent contract: `goal`, `non_goals`, `constraints`, `acceptance_criteria`
(outcome-level `{id, statement}`), optional coarse `scope`, provisional `tier`,
`status`, and the post-approval `contract_digest`. Legacy schema-1 intents (with
per-criterion `method`/`done_when`) are preserved as authored; new intents use
schema 2.

## Constraints and edge cases

The `INTENT.md` "Open questions" section is only for a genuine undecided question
the human must settle (a bold question, plus an italic answer once decided); it is
not a place for assurance rationale or accepted-risk notes — those belong in the
machine record. `None` means no known unresolved human decision **at the current
phase**, not that tracing cannot reveal one. A criteria, scope, tier, authority, or
lifecycle change after approval is an explicit edit to the contract that re-freezes
the digest and **re-enters Gate 1** (INV-CANDIDATE-01), voiding any candidate proven
against the old criteria. Implementation-only questions stay in the trace and plan.

## Implementation references

- `.agents/skills/cc-intent/SKILL.md`
- `.context-circuit/docs/templates/intent.md`, `.context-circuit/docs/templates/intent.example.md`
- `.context-circuit/wrapper/runtime/engine.sh`: `cc_intent_validate`, `cc_intent_approve`, `cc_intent_criteria_count`, `cc_intent_authorized`
- `.context-circuit/wrapper/contracts/schemas/intent-contract.yaml`
- `.context-circuit/wrapper/contracts/invariants.yaml`: INV-INTENT-01, INV-APPROVE-01 (owner map: `intent_gate`, `intent_contract`)

## Verification

`sh test/acceptance.sh` (intent suite; the feasibility suite pins the deliberate
absence of an automated scope gate).

## Provenance

Authored from the current wrapper for Context Circuit v1.0 (the two-gate,
tiered-assurance layer: INV-INTENT-01/02, INV-APPROVE-01, INV-ASSURE-01). Grounds
on the shipped contract; raw `sources/` was not scanned.

## Acceptance notes

Accepted 2026-09-04 (maintainer) to close a knowledge gap: v1.0 introduced the
intent front door and Gate 1 but no domain page owned it (authorization referenced
it only in passing). Supersedes nothing; complements
[plan authorization](../plan-authorization/README.md) (which owns how a plan
derives authorization from an approved intent) and
[tracing and feasibility](../tracing/README.md) (which owns what happens after
approval).
