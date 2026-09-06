---
name: cc-plan
description: Create a detailed grounded plan from a request, or conversationally review a draft plan without approving or executing it.
---

## Derive a plan from an approved intent

A v1.0 plan is the **derivation of an approved intent**, not the thing the human
approves (Context Circuit v1.0, Mechanism 1). Before planning a writing change,
there must be an approved parent intent (`cc-intent`, Gate 1). If none exists,
author the intent first — do not create a plan that invents its own goal or scope.
Plan derivation is then automatic — the coordinator's next action **in the same
turn** after a feasible tracer with no intent-level questions, not a gate — so
the human hears "here's the breakdown," never "approve this plan." Writing the
plans does not start execution. If the look is not feasible or an intent-level
question remains, write no plan.

If the approved intent is **Explore**, do not create a plan of record. Explore is
the planless `cc-pair` path: the user supervises the work live, and there is no
independent verifier or completion record to derive.

A plan is derived from the **trace manifest** the tracer produced after approval
(`cc-trace`, `intent/<id>/trace/<repo>.yaml`), not from a blind read of `context/`:
the tracer has already read the real code and reported the file/call-site map, a
proposed task partition, concrete risks, and the runnable "done" checks that prove
each outcome criterion. The feasibility check ran on those findings before you got
here — so you are planning a change already judged buildable. If no trace manifest
exists yet (the tracer has not run), run `cc-trace` first rather than improvising a
code read inline.

When `intent/<id>/detail/` exists, use it as the **confirmed shape** of what to
build so plans and tasks follow those topics. It does this without replacing the
trace manifest or the post-approval read of the real code, is not a second approval, and is not part of `contract_digest`.

## Ratify the task partition

Treat the trace's `task_partition` as grounded evidence, then make one explicit
plan-boundary decision. Every derived plan names **exactly one repository**; a
plan that lists two or more is invalid. Keep a **one-plan** derivation when all
tasks share one bounded execution/change surface in that one repository, one
worker lifecycle, and one independent verification boundary. Put the tasks in that
plan and use task `depends_on` for their intra-plan order; several tasks do not
require several plans.

Derive **stacked plans** when the partitions are independently executable or
independently verifiable, when there are meaningful dependencies between them,
when their failure surfaces need to be isolated, **or when the intent's scope
covers two or more repositories**. Give each plan that one repository and bounded
paths, record an acyclic `plan_dependencies` edge with a reason, and preserve one
worker plus one independent verifier lifecycle for every Standard/Critical plan.
Refuse collapsing two repositories into one plan.

Record the overall rationale in `context_grounding.decisions` and mirror it in the
readable `PLAN.md`; use each `plan_dependencies.reason` for the specific edge. Do
not split work solely because it is Standard or Critical, and do not combine
partitions merely to make one plan when their execution, verification, dependency,
or failure boundaries are genuinely independent. Plan count is not an assurance
tier and does not create a separate plan-approval gate. One approved intent may
authorize one plan or a stack of plans.

Before writing a plan, inspect the manifest's `open_questions` and
`out_of_scope_reach` and record a disposition for each finding. An
`intent-revision` question, or a required scope change not already authorized by
the plain request, blocks planning: return to `cc-intent`, update the approved
decision, and take Gate 1 again when the contract changes. A `plan-resolution`
question is carried into the plan's assumptions, risks, or verification. An
`already-answered` finding is applied and recorded without asking the human to
repeat it. No unresolved intent-level question may be hidden in a plan's open
questions.

Retrieve relevant Product Knowledge by the request's concepts, domains,
repositories, decisions, and constraints using `context/INDEX.md`; read only the
selected units, not the whole directory. Read repository instructions and only
request-named source files. Ground the plan in the approved `contract.yaml` (the
outcome criteria), the trace manifest (the grounded map and done-checks), the
intent detail when it exists (the confirmed topic shape), Product Knowledge, and
repository grounding.

Expand the approved intent into a detailed readable plan grounded in the manifest.
Preserve the intent's goal, criteria, constraints, non-goals, assumptions, open
questions, and risks — remove repetition, not meaning. Carry the tracer's executable
done-checks into the plan's verification, and its completeness proof for any "change
every X" obligation. Write `PLAN.md` (readable) and `plan.yaml` (canonical,
`schema_version: 3`) with:

- a stable plan id `NNNN-<kebab-slug>` allocated by the runtime `plan-allocate-id`;
- `intent: i<NNN>-slug` — the required parent intent (INV-INTENT-02);
- original request and request coverage;
- objective, desired behavior, constraints, non-goals;
- Product Knowledge grounding (each reference: stable id, path, reason);
- source and repository evidence;
- explicit repository mapping for every task, with bounded paths and dependencies;
- implementation details, acceptance criteria, and verification (distinct ids);
- expected commits, assumptions, open questions, risks, delivery notes;
- expected Product Knowledge impact, or an explicit no-durable-impact statement.

Perform a context-grounding drift check and a request-fidelity check. Any
contradiction or missing detail becomes an explicit open question, assumption,
or risk — never a silently chosen implementation. Add the plan to `plans/INDEX.md`
and validate it with the runtime `plan-validate`.

## Authorization preflight (INV-INTENT-02)

Run `intent-authorized . <plan-id>` as a preflight. If it reports `authorized: yes`,
the plan derives from an approved intent whose criteria are unchanged — no separate
plan approval, and execution may follow (INV-EXEC-01). If it reports `authorized: no`
with `reason: CRITERIA_CHANGED`, the intent's criteria were edited after approval:
the change re-enters Gate 1 — take approval again on the changed criteria via
`cc-intent`. A missing or unapproved intent is likewise unauthorized; author or
approve the intent first. There is **no** automated scope gate: scope-safety is
settled at delivery (Gate 2), the tracer reports where the change actually lands, and
a required change beyond a bound scope was surfaced as a question by the feasibility
check in `cc-trace` before you planned. The authorization check also runs again at
execution start.

One intent may yield **one or more** stacked plans, each naming the same intent
and exactly one repository. Deriving a plan does not execute it.

## Grounding and unupdated Product Knowledge (INV-KNOWLEDGE-02)

A later plan may start even if Product Knowledge from a prior plan has not yet
been updated. Do not wait on knowledge debt before grounding. Explore has no
plan preflight — use the planless `cc-pair` path instead of creating a plan
that would bypass the Explore assurance model.

## Review a plan

`review plan X` is a non-executing, status-preserving discussion. Walk through
request coverage, task detail, repository/path mapping, dependencies, context
references, acceptance, verification, assumptions, open questions, and risks.
When the human resolves a question, corrects a requirement, changes scope, or
asks for more detail, update the draft plan content and continue. Review never
changes plan status or executes. Authorization comes from the approved intent
(re-checked at execution start); there is no separate plan-approval request.
