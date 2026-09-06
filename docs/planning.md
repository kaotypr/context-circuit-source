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

When `intent/<id>/detail/` exists, planning uses it as the confirmed shape of what
to build so plans and tasks follow those topics. It does not replace the
post-approval read of the real code or the trace manifest.

## Choosing one plan or a stack

The trace's `task_partition` is a proposal. After feasibility, the coordinator
ratifies it using the real execution and verification boundaries:

- Keep a bounded Standard change in **one plan with embedded tasks** when it has one
  worker lifecycle, one independent verification boundary, and one change surface
  in **one repository**.
  Use task `depends_on` to order those tasks. Several tasks are not, by themselves,
  a reason to create several plans.
- Use **multiple stacked plans** when partitions can be independently executed or
  independently verified, have meaningful dependency edges, have distinct
  failure surfaces that should be isolated, **or when the intent covers two
  repositories**. Each plan names exactly one repository, an acyclic
  `plan_dependencies` list with edge reasons, and its own
  worker/verifier lifecycle at the applicable tier.

The coordinator records the overall decomposition rationale in
`context_grounding.decisions` and the readable `PLAN.md`, while
`plan_dependencies.reason` explains each inter-plan edge. Task `depends_on` and
`plan_dependencies` are different: the first orders tasks inside one plan; the
second orders separate plans. Do not split solely to match an assurance tier, and
do not combine genuinely independent execution, verification, dependency, or
failure boundaries. Plan count is not an assurance tier, and this decision never
introduces a separate plan-approval gate: one approved intent may authorize one
plan or several stacked plans.

## Detail

Plans are as detailed as necessary to preserve the request and make execution
unambiguous. Every task names the plan's single repository,
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
since approval). Standard and Critical `done` are an explicit mark-done;
Explore is planless. Execution status (running, verifying, repairing, verified,
failed, blocked) is runtime evidence and never replaces plan status.

See `docs/templates/plan.yaml`, `docs/templates/plan.md`, and
`docs/templates/task.md` for the shapes, and
`wrapper/contracts/schemas/plan.yaml` for the canonical field contract.
