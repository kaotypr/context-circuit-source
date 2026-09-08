# Plan 0028 — Band-scoped intent and plan id allocation

**Intent:** i021-id-number-blocks  
**Repository:** context-circuit-source  
**Tier:** Standard  
**Status:** draft  
**Depends on:** 0027-member-roster

## Objective

`intent-allocate-id`, `plan-allocate-id`, and `plan-stack-materialize` allocate
only inside the current member's bands. Next id advances past the highest ever used
in that band (including archived). Stacks reserve consecutive in-band plan ids.
Exhaustion and prefix collision fail closed. Existing ids and `cc/<plan-id>/<repo>`
branches stay unchanged.

## Grounding (HEAD 7db954f)

Allocation today uses workspace-wide high-water marks in `cc_intent_allocate_id`,
`cc_plan_allocate_id`, and `cc_plan_stack_materialize` (via `cc_plan_max_sequence`).
Archive does not recycle ids. Branches use `cc/<plan-id>/<repo>` only.

Plan 0027 supplies roster and identity resolution; this plan rewires allocation.

## Decisions

**Stack plan 2 of 3** — allocation depends on roster helpers; product docs and
invariant wording follow in 0029 once behavior is fixed.

- Band-scoped max helpers shared across intent, plan, and materialization paths.
- Workspace-wide prefix uniqueness enforced on every new allocation.
- Test fixtures gain roster/identity seeding; existing suites must stay green.

## Tasks

### BA-001 — Band-scoped intent allocation

Rewire `cc_intent_allocate_id` with band high-water scan and parallel-member tests.

**Done when:** two members allocate intents without shared prefixes; archived ids not reused.

### BA-002 — Band-scoped plan allocation and stacks

Rewire `cc_plan_allocate_id` and stack materialization for in-band consecutive reservation.

**Done when:** stack reserves consecutive plan ids in one band; plan and latency suites pass.

### BA-003 — Collision and exhaustion guards

Reject overlapping bands, duplicate prefixes, and band overflow without wrap.

**Done when:** band allocation acceptance suite passes; full semantic acceptance green.

## Risks

- Forgetting materialization `allocation.tsv` entries in band max would recycle reserved ids.
- Latency materialization fixtures may need roster bootstrap to keep smoke tests passing.

## Acceptance criteria mapping

| Intent criterion | Task evidence |
| --- | --- |
| AC-BAND-01 | BA-001 parallel intent; BA-002 parallel plan |
| AC-BAND-02 | BA-002 id form and branch unchanged |
| AC-BAND-05 | BA-001 archived intent never reused |
| AC-BAND-06 | BA-002 stack consecutive reservation |
| AC-BAND-07 | BA-002 branch form preserved |
| AC-BAND-08 | BA-003 existing ids untouched |
| AC-BAND-09 | BA-003 overlap and exhaustion guards |
