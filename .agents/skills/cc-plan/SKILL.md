---
name: cc-plan
description: Create a detailed grounded plan from a request, or conversationally review a draft plan without approving or executing it.
---

## Derive a plan from an approved intent

A v1.0 plan is the **derivation of an approved intent**, not the thing the human
approves (Context Circuit v1.0, Mechanism 1). Before planning a writing change,
there must be an approved parent intent (`cc-intent`, Gate 1). If none exists,
author the intent first — do not create a plan that invents its own goal or scope.
Plan derivation is then automatic — the coordinator's next action, not a gate — so
the human hears "here's the breakdown, building now," never "approve this plan."

Retrieve relevant Product Knowledge by the request's concepts, domains,
repositories, decisions, and constraints using `context/INDEX.md`; read only the
selected units, not the whole directory. Read repository instructions and only
request-named source files. Ground the plan in the approved `contract.yaml` (the
criteria and scope), Product Knowledge, and repository grounding.

Expand the approved intent into a detailed readable plan. Preserve the intent's
goal, criteria, constraints, non-goals, assumptions, open questions, and risks —
remove repetition, not meaning. Write `PLAN.md` (readable) and `plan.yaml`
(canonical, `schema_version: 3`) with:

- a stable plan id `NNNN-<kebab-slug>` allocated by the runtime `plan-allocate-id`;
- `intent: i<NNNN>-slug` — the required parent intent (INV-INTENT-02);
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

## Envelope preflight (INV-INTENT-02, crown jewel 1)

Run `intent-envelope-check . <plan-id>` as a preflight. If it reports
`within`, the plan is authorized by the approved intent — no separate plan
approval, and execution may follow (INV-EXEC-01). If it reports `exceeds` (a new
repository, a path region outside the intent's scope, or a criteria change), the
plan is **held and re-gated to a human** in plain language: offer to widen the
intent (a new decision — re-freeze the digest and re-run the adversary on any
changed criteria via `cc-intent`) or to narrow the plan back inside scope. Never
proceed on the original approval, and never widen scope on your own. The check
also runs again at execution start.

One intent may yield **one or more** stacked plans, each naming the same intent
and each staying inside the envelope. Deriving a plan does not execute it; the
grounding-debt preflight (`cc-plan` at Standard/Critical) still applies.

## Review a plan

`review plan X` is a non-executing, status-preserving discussion. Walk through
request coverage, task detail, repository/path mapping, dependencies, context
references, acceptance, verification, assumptions, open questions, and risks.
When the human resolves a question, corrects a requirement, changes scope, or
asks for more detail, update the draft plan content and continue. Review never
changes plan status, approves, or executes. Approval is a separate explicit
request handled during execution.
