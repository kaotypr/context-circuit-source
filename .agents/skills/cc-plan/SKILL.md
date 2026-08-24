---
name: cc-plan
description: Create a detailed grounded plan from a request, or conversationally review a draft plan without approving or executing it.
---

## Create a plan

Retrieve relevant Product Knowledge by the request's concepts, domains,
repositories, decisions, and constraints using `context/INDEX.md`; read only the
selected units, not the whole directory. Read repository instructions and only
request-named source files.

Expand the request into a detailed readable plan. Preserve the request's
requirements, desired behavior, constraints, non-goals, assumptions, open
questions, and risks — remove repetition, not meaning. Write `PLAN.md` (readable)
and `plan.yaml` (canonical) with:

- a stable plan id `NNNN-<kebab-slug>` allocated by the runtime `plan-allocate-id`;
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
and validate it with the runtime `plan-validate`. Creating a plan does not
approve or execute it.

## Review a plan

`review plan X` is a non-executing, status-preserving discussion. Walk through
request coverage, task detail, repository/path mapping, dependencies, context
references, acceptance, verification, assumptions, open questions, and risks.
When the human resolves a question, corrects a requirement, changes scope, or
asks for more detail, update the draft plan content and continue. Review never
changes plan status, approves, or executes. Approval is a separate explicit
request handled during execution.
