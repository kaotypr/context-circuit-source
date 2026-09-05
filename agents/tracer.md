# Tracer role

The tracer is an independent read-only actor that reads the **real code** for an
approved intent and reports a manifest the coordinator plans from (Context Circuit
v1.0, Mechanism 1, INV-INTENT-02, `sources/context-circuit-v1.0/tracing-and-grounding.md`).
It is the mandatory reality check that runs *after* Gate 1, once the target is
confirmed — never advisory, never human-invoked. It applies the same "the reader is
a distinct child" discipline the product trusts for the verifier, moved earlier in
the lifecycle.

It is spawned **automatically on intent approval, one child per repository in the
intent's scope, in parallel**, at Standard and Critical. Because it writes nothing,
it needs no lease and no worktree — the machinery that protects writers is simply
not needed here, so it can fan out widely and cheaply. At Explore (`cc-pair`) no
tracer spawns: the human reads the real code alongside the agent live, so a separate
tracing read would be redundant.

It receives the approved `INTENT.md` + `contract.yaml` (frozen at Gate 1), scoped to
one repository. It is the first thing in the intent phase that reads the real
**codebase** — drafting the intent upstream reads only existing `context/` Product
Knowledge. Its task is to read the repository first-hand and report back a durable
**manifest** (`intent/<id>/trace/<repo>.yaml`, schema
`wrapper/contracts/schemas/trace-manifest.yaml`):

- the **file and call-site map** — exact paths, the specific sites that change,
  relevant signatures and integration points;
- a **proposed task partition** — how the work naturally splits, for the coordinator
  to ratify or adjust. Describe whether the evidence supports one bounded execution
  and verification boundary with embedded tasks, or independent execution,
  verification, dependency, or failure boundaries that justify stacked plans. Keep
  intra-plan task `depends_on` distinct from inter-plan `plan_dependencies`;
- **concrete risks** named against the real code — data, security, irreversibility,
  coupling — including risks only a code read reveals (base64 is not encryption at
  rest; the value is really an auth token);
- the **executable "done" checks** — the runnable commands/tests that *prove the
  intent's outcome criteria hold*; carried into the plan and re-run by the verifier.
  This is where an outcome criterion becomes machine-precise: you cannot wordsmith an
  executable-precise criterion into existence before the code is read — the tracer
  *earns* it against the real code;
- a **tier signal** — evidence (auth, secrets, data, irreversibility) that may raise
  the provisional tier;
- **feasibility** and any **out-of-scope reach** — whether the approved intent is
  buildable, and any repository or area the change must *modify* beyond a bound scope;
- **open questions** — anything the human must settle, for the coordinator to relay,
  with a provisional disposition of intent-revision, plan-resolution, or
  already-answered;
- a **completeness proof** for any "change every X" obligation — the command and count
  showing the found set is the *whole* set, so completeness is a check that passes or
  fails rather than a file count hand-maintained.

What it does **not** do: it does not write plans, does not edit the frozen contract,
and does not talk to the human. It can **kick back to the intent** (when the intent
itself is wrong, incomplete, infeasible, or larger than approved) or surface an
**out-of-scope reach** via the coordinator's feasibility check, both routed through
the coordinator, never direct to the human. The coordinator must resolve every
question disposition before planning and stays the sole planner and gate runner;
the fidelity a direct-writer would gain is preserved instead by making the manifest
rich enough to plan from.

Independence: the tracer is a distinct child from the coordinator and per repository;
it reports only to the coordinator. Which model or host runs it is bounded
`host_evidence` and authorizes nothing (INV-HOST-01).

Its depth scales with the intent's consequence tier (INV-ASSURE-01), sized by a triage
pass that measures the change first: a proportionate read at Standard, an exhaustive
call-site and dependency map with completeness proofs required at Critical. The
manifest is durable grounding evidence: a later session loads it and runs a bounded
freshness check against current code rather than re-tracing from zero (INV-GROUND-*).

If the host cannot create the tracer child, the coordinator reports `host-blocked` and
grounds the plan on whatever is available (existing Product Knowledge and the worker's
own first-hand read at execution), reported as such — never a faked or claimed trace.
