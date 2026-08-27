# Proposal 0021 — change: DECISIONS.md (v0.6 accepted decisions)

- id: 0021-change-decisions
- target_context_unit: context/DECISIONS.md
- operation: change
- confidence: high
- status: review-needed
- related_plan: —
- affected_commits: [74eb510, e3e95cd, dac0b76, 660b6a7, 343dedb, 048eb67]
- evidence_refs:
  - sources/system-design/context-circuit/v0.6/
  - wrapper/manifest.yaml (runtime_version 0.6.0; accepted_schema_versions.plan [1, 2])
  - wrapper/contracts/invariants.yaml (INV-PLAN-05, INV-CONCURRENCY-01/02, INV-GROUND-01/02/03, INV-DELIVER-01 drift clause)

## Statement

Append the v0.6 decisions to `context/DECISIONS.md` (dated 2026-08-27):

1. **Coordinated v0.6 contract bump.** `runtime_version` 0.5.0 → 0.6.0; plan
   `accepted_schema_versions` becomes `[1, 2]`. Rationale: `plan_dependencies` is
   load-bearing, so a v0.5 engine must refuse a `schema_version: 2` plan rather
   than schedule it dependency-blind (INV-PLAN-05); `execution.yaml` stays schema 1
   (private, single-version reader). Consequence: only plans using
   `plan_dependencies` stamp 2; existing plans stay 1.

2. **Concurrency is orchestration, not authority (run-stack).** Executing a set of
   approved plans changes only order and overlap; approval, verification,
   completion, and delivery gates are unchanged. Conflicts are *prevented*
   (dependencies order waves, path leases serialize file overlaps, a dependent's
   base already contains its prerequisites), not resolved after the fact. The
   runtime detects readiness/leases/bases deterministically; the coordinator
   decides how many ready plans to launch — no scheduler heuristic in the runtime
   (INV-RUNTIME-01, INV-CONCURRENCY-01/02).

3. **Delivery drift guard.** A plan whose recorded base diverged from the anchor
   tip is rebased and re-verified before its pull request (INV-DELIVER-01,
   extended). The only merge the runtime authors is the integration *base*, never a
   delivery merge.

4. **Repository grounding: reference, not capture.** The writer honors the target
   repository's own agent guidance, discovered live from the worktree as data
   (INV-GROUND-01); precedence is CC scope/safety on *what/where*, repo guidance on
   *how* within that scope (INV-GROUND-02); the brief is a template filled from the
   manifest and delivered, never authored (INV-GROUND-03). No per-repo profile and
   no `plan.yaml` field.

5. **A system design is a source, not a lifecycle stage.** The v0.6
   system-design scope ships only the `cc-system-design` authoring skill; a system
   design lives under `sources/system-design/` with no status, acceptance gate, or
   runtime record, and feeds Product Knowledge and plans through the existing flow.

6. **Efficiency ledger made real (maintainer tooling).** The template-harness
   dimension D now measures per-action usage from the runner's own result and
   compares it to case budgets (units fixed: `max_tokens` = output tokens,
   `max_turns` = conversational turns), staying soft (never gates). No product
   surface change.

## Change

Append the six dated decision blocks above to `context/DECISIONS.md`, each in the
existing Decision / Rationale / Consequence form, noting acceptance from this
proposal.
