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

Retrieve relevant Product Knowledge by the request's concepts, domains,
repositories, decisions, and constraints using `context/INDEX.md`; read only the
selected units, not the whole directory. Read repository instructions and only
request-named source files. Ground the plan in the approved `contract.yaml` (the
outcome criteria), the trace manifest (the grounded map and done-checks), Product
Knowledge, and repository grounding.

Expand the approved intent into a detailed readable plan grounded in the manifest.
Preserve the intent's goal, criteria, constraints, non-goals, assumptions, open
questions, and risks — remove repetition, not meaning. Carry the tracer's executable
done-checks into the plan's verification, and its completeness proof for any "change
every X" obligation. Write `PLAN.md` (readable) and `plan.yaml` (canonical,
`schema_version: 3`) with:

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

One intent may yield **one or more** stacked plans, each naming the same intent.
Deriving a plan does not execute it.

## Grounding-debt preflight (INV-KNOWLEDGE-02, closed loop)

Before grounding a new plan, run `knowledge-debt-check . <plan-id>`. It consults
the reconciliation-debt markers that completion/delivery emitted (INV-COMPLETE-02).
If delivered work in this plan's knowledge scope (its repositories) is still
unreconciled:

- **Standard/Critical → it blocks.** Reconcile first — generate the impact
  proposals and let the human accept or explicitly defer them
  (`knowledge-reconciled . <candidate> reconciled|deferred`) — before grounding
  proceeds. Say it in plain language: "there's merged work in this area I haven't
  folded into what the project knows yet — reconcile it first, or mark it as no
  update needed?" Never auto-accept knowledge (INV-KNOWLEDGE-02); the human decides.
- **Explore has no plan preflight.** Use the planless `cc-pair` path instead of
  creating a plan that would bypass the Explore assurance model.

`clear` means no overlapping debt and grounding proceeds normally.

## Review a plan

`review plan X` is a non-executing, status-preserving discussion. Walk through
request coverage, task detail, repository/path mapping, dependencies, context
references, acceptance, verification, assumptions, open questions, and risks.
When the human resolves a question, corrects a requirement, changes scope, or
asks for more detail, update the draft plan content and continue. Review never
changes plan status or executes. Authorization comes from the approved intent
(re-checked at execution start); there is no separate plan-approval request.
