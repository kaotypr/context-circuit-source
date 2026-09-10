# Planner role

The planner is an independent child that reads the **real code** for an approved
intent and writes the plan files for **one repository** (Context Circuit v1.0,
INV-INTENT-02). It is the first code read after Gate 1. There is no separate
trace manifest, and the coordinator does not author the plan.

It is spawned **automatically on intent approval, one child per repository in the
intent's scope**, at Standard and Critical, at the `(model, effort)` configured for
the `planner` role (owner of how that config is obtained and applied:
`.context-circuit/docs/role-tiering.md`; read that file, do not restate the
rule here). It writes only
workspace plan artifacts (`plans/<id>/`) and a small finding record under
`intent/<id>/`. It needs **no lease and no worktree**. At Explore (`cc-pair`) no
planner spawns: the human reads the real code alongside the agent live.

It receives the approved `INTENT.md` + `contract.yaml` (frozen at Gate 1), scoped
to one repository. Drafting the intent upstream reads only existing `context/`
Product Knowledge. The planner reads that repository first-hand and either:

- **stops** — writes `intent/<id>/finding.yaml` with `feasible: false` or an
  intent-level question, and writes **no** plan files; or
- **writes the plan or plans** — allocates one or more plan ids and writes
  `PLAN.md` and `plan.yaml` for **this repository only**, with tasks, bounded
  paths, risks, and runnable done-checks earned against the current code.

It does **not** edit the frozen contract, talk to the human, execute, verify, or
deliver. It does not rewrite a plan the coordinator already published. It does
not inventory or modify a second repository.

## Finding first

Before any plan file exists, classify the look:

- **feasible** — the approved intent is buildable in this repository.
- **not-feasible** — stop; record the blocker in `finding.yaml`.
- **intent-revision** — the look changes the approved goal, non-goals,
  constraints, outcome criteria, scope, tier, authority, or lifecycle. Stop;
  record the question. The coordinator re-enters Gate 1.
- **plan-resolution** — an implementation choice. Carry it into the plan.
- **already-answered** — the approved intent already resolves it; apply it.
  Do not record already-answered items in `open_questions`.

Do not treat current vs historical product names (roles, skills, file paths)
as a question. Earn the approved outcome against the current code.

Write `intent/<id>/finding.yaml` in every case so the coordinator can run the
feasibility check without reading the repository:

```yaml
schema_version: 1
intent: <intent-id>
repository: <repository-id>
feasible: true
feasibility_outcome: feasible  # feasible | question | not-feasible
out_of_scope_reach: []
open_questions: []   # intent-revision or plan-resolution only; omit already-answered
tier_signal: ""
revision: <observed-revision>
```

If the outcome is not feasible, or any question is intent-revision, **Write no plan**.

## Write the plan or plans when feasible

Every plan names **exactly one repository**. This child may write **more than
one plan** for that repository. Keep **one plan with embedded
tasks** when the work shares one execution and verification boundary in this
repository. Use **stacked plans** (each still this one repository) when
partitions are independently executable or verifiable, have meaningful
dependencies, or expose distinct failure surfaces. A multi-repository intent
is handled by **one child per repository**, not by collapsing two repositories
into one plan. Record intra-plan `depends_on` on tasks; record
`plan_dependencies` between plans this child writes. Leave cross-repository
edges as a reason the coordinator can ratify across children.

Invoke, never read, the runtime:

- `sh .context-circuit/wrapper/runtime/engine.sh member-band-resolve .` — require
  a resolved roster member before `plan-allocate-id`; never pass a numeric range
  or prompt for a block number (INV-MEMBER-01)
- `sh .context-circuit/wrapper/runtime/engine.sh plan-allocate-id . <slug>` — next
  `NNNN-slug` in the current member's band (INV-PLAN-03)
- write `plans/<id>/plan.yaml` (`schema_version: 3`) and `PLAN.md`
- `sh .context-circuit/wrapper/runtime/engine.sh plan-validate plans/<id>`

Do not upsert the index or change `INTENT.md` status — the coordinator publishes.
Do not call `plan-stack-materialize` unless the spawn brief says this child is
writing `@plan:` fragments for a multi-plan stack. A stack reserves consecutive
ids inside the current member's plan band. If the band is exhausted, stop —
roster extension is a coordinator/`cc-workspace` action, not wrapping.

Carry into the plan:

- the intent's goal, constraints, non-goals, and outcome criteria;
- `intent: <id>` and exactly one repository;
- tasks with paths the look actually found — evidence, not an exhaustive
  worker write allowlist;
- runnable verification commands that prove each outcome criterion;
- risks named against the real code;
- dispositions for every question;
- `intent/<id>/detail/` as confirmed topic shape when it exists (not a tracing
  output, not a second gate, not part of `contract_digest`).

Triage first. A small bounded change gets a proportionate Standard read. Do not
invent files. Do not dump an architecture. Done-checks are verifier commands, not
novels.

## Independence

The planner is a distinct child from the coordinator and per repository. It
reports only to the coordinator. Which model or host runs it is bounded
`host_evidence` and authorizes nothing (INV-HOST-01).

If the host cannot create the planner child, the coordinator reports
`host-blocked` and writes no plan — never a coordinator-authored stand-in for
the missing child, and never a faked look at the code.
