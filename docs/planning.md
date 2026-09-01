# Planning

A plan is the execution contract. It has a readable `PLAN.md` and a canonical
`plan.yaml`; both describe the same plan, and `plan.yaml` owns identifiers, status,
repository mappings, dependencies, and verification ids.

## A plan derives from an approved intent (Context Circuit v1.0)

A v1.0 plan is the **derivation of an approved intent** (`.agents/skills/cc-intent`),
not the thing the human approves. Before planning a writing change there must be an
approved parent intent, which holds the goal, acceptance criteria, scope envelope,
and tier. The plan names it (`intent: i<NNNN>-slug`, `schema_version: 3`) and stays
inside its scope envelope; the envelope check runs as a preflight and again at
execution start, and a plan that exceeds the approved scope is re-gated to a human
rather than approved. Because the human already approved the intent, the plan
carries no second approval gate. A legacy plan without a parent intent keeps the
earlier explicit `draft → approved` plan gate.

## Grounding

Before drafting, retrieve the Product Knowledge that governs the work by concept,
domain, repository, decision, and constraint using `context/INDEX.md`. Read only
the selected units. Perform a context-grounding drift check and a
request-fidelity check: any contradiction or unresolved detail becomes an
explicit open question, assumption, or risk — never a silently chosen
implementation.

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

For a v1.0 plan, status is a **projection**: `approved` follows automatically from
the approved intent within its envelope, and `done` is inferred from candidate
acceptance + delivery (explicit at the Critical tier). For a legacy plan, status is
human-controlled: `draft → approved → done`. Execution status (running, verifying,
repairing, verified, failed, blocked) is runtime evidence and never replaces plan
status.

See `docs/templates/plan.yaml`, `docs/templates/plan.md`, and
`docs/templates/task.md` for the shapes, and
`wrapper/contracts/schemas/plan.yaml` for the canonical field contract.
