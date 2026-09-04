# Planning

A plan is the execution contract. It has a readable `PLAN.md` and a canonical
`plan.yaml`; both describe the same plan, and `plan.yaml` owns identifiers, status,
repository mappings, dependencies, and verification ids.

## A plan derives from an approved intent (Context Circuit v1.0)

A v1.0 plan is the **derivation of an approved intent** (`.agents/skills/cc-intent`),
not the thing the human approves. Before planning a writing change there must be an
approved parent intent, which holds the goal, outcome-level acceptance criteria, a
coarse optional scope, and tier. On approval a read-only tracer reads the real code
and reports a manifest; the plan is derived from that manifest and names the intent
(`intent: i<NNN>-slug`, `schema_version: 3`). The derives-from-an-approved-intent
authorization runs as a preflight and again at execution start — a criteria change
after approval re-enters Gate 1 — but there is no automated scope gate: scope-safety
is settled at delivery (Gate 2), and a required change beyond a bound scope is
surfaced by the feasibility check before planning. Because the human already approved
the intent, the plan carries no second approval gate — every plan derives from an
approved intent, and there is no plan-level `approved` status.

## Grounding

Before drafting, retrieve the Product Knowledge that governs the work by concept,
domain, repository, decision, and constraint using `context/INDEX.md`. Read only
the selected units. Perform a context-grounding drift check and a
request-fidelity check: any contradiction or unresolved detail becomes an
explicit open question, assumption, or risk — never a silently chosen
implementation.

Questions are dispositioned by phase before planning. A trace finding that changes
the approved goal, scope, criteria, tier, authority, or lifecycle is an intent-level
question and sends the work back through Gate 1; an implementation-only question is
carried into the plan; a question already answered by the user's request is applied
without asking again. No plan may hide an unresolved intent-level question.

## Detail

Plans are as detailed as necessary to preserve the request and make execution
unambiguous. Every task names the repository or repositories it may change,
bounded paths or an explicit repository-wide scope, dependencies, the concrete
change, acceptance (desired result), verification (evidence), and stop
conditions. Detail means concrete behavior, not repeating the same explanation
in every task.

## Identifiers

New plans use stable ids `NNNN-<kebab-slug>`. The four-digit sequence is the next
after the highest ever allocated and is never reused, including after archive or
restore. The id is the canonical mention key in conversation, paths, branches,
and runtime records.

## Status

Plan status is a **projection**, not a second gate. A plan is `draft` from creation
until it completes, then `done` — there is no intermediate `approved` status.
Authorization follows automatically from the approved intent (criteria unchanged
since approval). Standard `done` is inferred from candidate acceptance + delivery;
Critical `done` is explicit. Execution status (running, verifying, repairing, verified,
failed, blocked) is runtime evidence and never replaces plan status.

See `docs/templates/plan.yaml`, `docs/templates/plan.md`, and
`docs/templates/task.md` for the shapes, and
`wrapper/contracts/schemas/plan.yaml` for the canonical field contract.
