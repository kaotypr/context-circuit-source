# Proposal 0018 — change domain: plan-execution (base-aware begin + plan schema)

- id: 0018-change-plan-execution
- target_context_unit: context/domains/plan-execution/README.md
- operation: change
- confidence: high
- status: review-needed
- related_plan: —
- affected_commits: [74eb510, e3e95cd]
- evidence_refs:
  - wrapper/runtime/engine.sh (cc_execution_begin base-aware branch; cc_base_prepare; grounding discovery in the per-repo loop)
  - wrapper/contracts/schemas/execution.yaml (based_on; grounding/<repo>.yaml)
  - wrapper/contracts/schemas/plan.yaml (schema_version [1, 2]; plan_dependencies)
  - wrapper/contracts/invariants.yaml (INV-CONCURRENCY-02, INV-PLAN-05, INV-GROUND-01)
  - .agents/skills/cc-execute/SKILL.md (grounded-brief step)

## Statement

The shipped plan-execution behavior gained two v0.6 extensions that the page
should record (shipped-is-truth):

1. **`cc_execution_begin` is base-aware for a plan with same-repo predecessors.**
   Such a plan's worktree is prepared by `cc_base_prepare` (a single predecessor
   branch, or a runtime-authored integration merge of several), and the repository
   record gains `based_on`. A plan with **no** dependency is unchanged: it keeps
   `cc_worktree_prepare` from the captured anchor tip (INV-EXEC-03 intact). Base
   selection and leases belong to the [run-stack](../run-stack/README.md) domain
   (proposal 0016); this page cross-references it.
2. **The execution brief is now a grounded, delivered brief.** Execution begin also
   discovers the target repository's own agent guidance and records a grounding
   manifest; the coordinator delivers a brief assembled by the runtime rather than
   free-composed (see [repository-grounding](../repository-grounding/README.md),
   proposal 0017). One worker / one independent verifier / three-failure limit are
   unchanged.

Also note the plan contract now accepts `schema_version` `[1, 2]` (2 required when
`plan_dependencies` is present, INV-PLAN-05); single-plan v0.5 plans stay
`schema_version: 1` and behave exactly as before.

## Change (edits to the plan-execution page on acceptance)

- In **Behavior / Execution brief**: note the base-aware begin for same-repo
  predecessors and the grounded brief; keep the no-dependency anchor-tip path as
  the default.
- In **Data**: `execution.yaml` repository records may carry `based_on`; a
  `grounding/<repo>.yaml` manifest is recorded per execution.
- In **Implementation references**: add `cc_base_prepare`,
  `cc_discover_repo_grounding`, `cc_writer_brief_assemble`, and invariants
  INV-CONCURRENCY-02, INV-PLAN-05, INV-GROUND-01.
- Add cross-references to the new `run-stack` and `repository-grounding` domains.
