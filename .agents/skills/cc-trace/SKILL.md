---
name: cc-trace
description: On intent approval, spawn one read-only tracer child per repository in scope to read the real code and report a manifest, then run the feasibility check on the findings so the coordinator can plan from a grounded basis. Read-only over the code; never edits the contract; never writes plans.
---

Tracing is the phase that turns an approved plain intent into a grounded basis for
planning (Context Circuit v1.0, Mechanism 1, INV-INTENT-02,
`sources/context-circuit-v1.0/tracing-and-grounding.md`). It is the answer to "where
does the agent learn, in detail, what it is going to build?" — a structured, recorded,
reusable phase, not an ad-hoc grep the coordinator improvises inline.

It runs **after Gate 1 and before any plan exists**. The human's approval confirmed
you understood the plain ask, so the tracer aims at a confirmed target and never fires
on an unconfirmed guess. At **Explore** no tracer spawns — the human reads the real
code alongside the agent live (`cc-pair`), so a separate tracing read would be
redundant. At **Standard** and **Critical** it is mandatory.

## Spawn the tracer — one read-only child per repository

For each repository in the intent's scope, spawn a **tracer child** (`agents/tracer.md`)
in parallel. Each reads *its* repository first-hand for this change. Because the tracer
writes nothing, it needs **no lease and no worktree** — the machinery that protects
writers is not needed, so the reads fan out widely and cheaply (this mirrors execution's
per-repo worker fan-out, moved earlier). If the intent bound no repository, trace the
repositories the change plainly implicates and report where it lands.

A triage pass inside tracing **sizes the change first**, so a one-file change gets a
glance and a cross-cutting change gets a full survey; depth scales with tier (a
proportionate read at Standard, exhaustive call-site and dependency mapping with
completeness proofs at Critical).

Each child reports a manifest, recorded at `intent/<id>/trace/<repo>.yaml` (schema
`wrapper/contracts/schemas/trace-manifest.yaml`) as durable grounding evidence:

- the **file/call-site map** — exact paths, the sites that change, signatures,
  integration points;
- a **task partition** — how the work naturally splits, for you to ratify or adjust;
- **concrete risks** named against the real code (data, security, irreversibility,
  coupling) — including risks only a code read reveals;
- the **executable "done" checks** — the runnable commands/tests that prove each
  outcome criterion holds (a `grep` proving zero direct call-sites remain); these are
  carried into the plan and re-run by the verifier, and are why an outcome criterion
  becomes machine-precise only now;
- a **tier signal** — evidence (auth, secrets, data, irreversibility) that may raise
  the provisional tier;
- **feasibility** and any **out-of-scope reach**;
- **open questions** for you to relay, each with a provisional disposition:
  `intent-revision`, `plan-resolution`, or `already-answered`;
- a **completeness proof** for any "change every X" obligation — the command and count
  showing the found set is the whole set.

### Surface the plan boundary for the coordinator

`task_partition` is evidence for the coordinator to ratify, not an instruction for
the tracer to turn into plans. The partition report should make the execution
boundary visible in human terms:

- Keep the work in **one plan with embedded tasks** when the tasks share one bounded
  execution/change surface, one worker lifecycle, and one independent verification
  boundary. Order those tasks with their intra-plan `depends_on` relationships.
- Use **multiple stacked plans** when partitions can be independently executed or
  independently verified, have meaningful dependencies between them, or expose
  distinct failure surfaces. Each partition must retain a bounded repository/path
  mapping and an acyclic inter-plan dependency reason.

Task count, repository count, and assurance tier are not substitutes for this
judgment. A Standard change may still be one plan with several tasks, and a larger
change may be several plans under the same approved intent. The coordinator records
why it kept or split the partition after feasibility, then explains that choice to
the human; the tracer only supplies the grounded evidence.

The tracer never writes plans, never edits the frozen contract, and never talks to the
human — it finds, records, and reports to you. Intent detail (`intent/<id>/detail/`),
when present, was authored before approval without reading the code. Tracing remains
the first code read. Do not author or edit detail here, and do not fold it into
`contract_digest`. A later plan may consume it as confirmed topic shape when it
exists.

### Classify every question before feasibility

The coordinator must disposition every trace question before planning:

- **Intent revision:** the finding changes the approved goal, non-goals,
  constraints, outcome criteria, scope, tier, authority, or lifecycle semantics.
  Stop planning, update the intent and its human-facing Open questions section,
  re-enter Gate 1, and trace again only as needed after approval.
- **Plan resolution:** the finding is an implementation choice that does not
  change the approved decision. Carry it into the derived plan's assumptions,
  risks, or verification; no additional human gate is introduced.
- **Already answered:** the plain request or approved intent already resolves the
  question. Apply that answer, record the rationale, and do not ask the human to
  repeat it.

An approved intent with an unresolved intent-level question is not ready for plan
derivation. This is a feasibility outcome, not a plan review detail. A scope reach
that the user's plain request already includes (such as an explicit "every current
reference" requirement) is `already-answered`; if the approved contract's coarse
scope must be widened to reflect it, update the contract and re-enter Gate 1 rather
than asking the user to restate the requirement.

## Freshness — pay for tracing once

If a manifest for this intent already exists, do not re-trace from zero: load it and run
a **bounded freshness check** against current code (have the files it describes changed
since it was written?), re-reading only what drifted and updating `freshness_checked_at`.
Captured design knowledge (`sources/`, `context/`) is reused the same way; the freshness
check is what keeps a reused snapshot from silently going stale.

## Run the feasibility check on the findings

Once the manifests are in, make the two coordinator judgments before writing any plan.
The **feasibility check** (`sources/context-circuit-v1.0/intent-feasibility.md`) is a
**quality** gate, not a safety gate, and it **fails upward** — when you cannot conclude
the change is buildable, stop and ask a human rather than proceed:

- **Feasible** → set the tier (re-check with the tracer's `tier_signal`; a raise
  surfaces at plan review) and hand `cc-plan` the grounded basis. No human step.
- **Intent revision required** → do not write a plan. Relay the specific finding,
  update the intent and its Open questions, and take Gate 1 again if the approved
  decision changes.
- **Not feasible** → **stop.** Explain the specific blocker the tracer found, suggest
  what would make it work, and hand the decision to the human. Cheap — nothing executed.
- **A required change that must *modify* a repository or area beyond a bound scope** →
  surface it as a question only when the plain request does not already authorize
  that area. If it changes the approved scope, treat the answer as an intent
  revision and re-enter Gate 1; merely reading or calling an out-of-scope area does
  not trigger this. When no scope was bound, there is nothing to reach past — just
  report where the change lands.

Feasibility is a reasoned judgment over the tracer's findings, not an engine verb; there
is **no** automated scope-containment check. Scope-safety is settled at delivery (Gate 2).

## Feedback edges

- **Kick back to the intent.** If the *intent itself* is wrong, incomplete, or larger
  than approved, re-open it via `cc-intent` rather than planning the wrong goal in
  detail. Nothing has executed, so this is cheap.
- **Tier raise.** A tier signal raises the provisional tier; the raise surfaces at plan
  review, never silently.

## Report — plain language only

The human never hears "the tracer" or "the manifest." Relay only what they must decide:
"I looked at the code — this touches X and Y; one thing I found is Z; here's a question
I need you to settle." Report an out-of-scope reach or an infeasible intent as a plain
choice, not a mechanism.

## Boundaries

Tracing reads code and writes only under `intent/<id>/trace/`. It never edits the
contract, never writes a plan, never executes, verifies, completes, or delivers, and
never talks to the human directly — the coordinator relays. If the host cannot create a
tracer child, report `host-blocked` and ground the plan on whatever is available
(existing Product Knowledge, the worker's own read at execution), reported as such —
never a faked trace, and the route stays read-only until the human decides.
