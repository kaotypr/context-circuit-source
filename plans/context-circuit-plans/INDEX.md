# context-circuit-plans — maintainer plan index

Implementation plans that build Context Circuit from an empty repository, in
dependency order. Every plan targets the `context-circuit-source` repository and
follows the current plan contract (`wrapper/contracts/schemas/plan.yaml`).
Inter-plan ordering is encoded in each `plan.yaml` `plan_dependencies` field
(schema_version 2); the table's order is one valid topological sort of that DAG.

Layers: **A** v0.5 core product (0001–0010) · **B** v0.5 test harness, source-only
(0011–0012) · **C** v0.6 product scopes (0013–0018) · **D** v0.6 skill + maintainer
tooling (0019–0020) · **E** v0.6.1 refinements (0021–0025).

| Plan ID | Title | Status | Depends on | Path |
| --- | --- | --- | --- | --- |
| 0001-contracts-runtime-foundation | Contracts, runtime foundation, and host adapters | done | — | plans/context-circuit-plans/0001-contracts-runtime-foundation/PLAN.md |
| 0002-workspace-repository-model | Workspace and repository model | done | 0001 | plans/context-circuit-plans/0002-workspace-repository-model/PLAN.md |
| 0003-context-lifecycle | Product Knowledge and context lifecycle | done | 0001 | plans/context-circuit-plans/0003-context-lifecycle/PLAN.md |
| 0004-planning | Grounded planning | done | 0002, 0003 | plans/context-circuit-plans/0004-planning/PLAN.md |
| 0005-approval-and-execution | Approval and multi-repository execution | done | 0004 | plans/context-circuit-plans/0005-approval-and-execution/PLAN.md |
| 0006-verification-and-repair | Independent verification and repair | done | 0005 | plans/context-circuit-plans/0006-verification-and-repair/PLAN.md |
| 0007-completion-and-reconciliation | Human completion and knowledge reconciliation | done | 0006, 0003 | plans/context-circuit-plans/0007-completion-and-reconciliation/PLAN.md |
| 0008-archive-and-restore | Archive and restore | done | 0004 | plans/context-circuit-plans/0008-archive-and-restore/PLAN.md |
| 0009-delivery | Delivery boundary | done | 0005 | plans/context-circuit-plans/0009-delivery/PLAN.md |
| 0010-release-assembly-and-publication | Source/template boundary, seed, release assembly and publication | done | 0001–0009 | plans/context-circuit-plans/0010-release-assembly-and-publication/PLAN.md |
| 0011-human-simulated-harness | Human-simulated test harness mechanism | done | 0010 | plans/context-circuit-plans/0011-human-simulated-harness/PLAN.md |
| 0012-scenario-library-and-host-matrix | Scenario library and host matrix | done | 0011 | plans/context-circuit-plans/0012-scenario-library-and-host-matrix/PLAN.md |
| 0013-run-stack-substrate | Run-stack substrate: plan dependencies, path leases, execution bases | done | 0004, 0005 | plans/context-circuit-plans/0013-run-stack-substrate/PLAN.md |
| 0014-run-stack-scheduling | Run-stack scheduling, run action, and delivery drift guard | done | 0013, 0006, 0009 | plans/context-circuit-plans/0014-run-stack-scheduling/PLAN.md |
| 0015-repository-grounding | Repository grounding: discovery, manifest, and generated writer brief | done | 0005, 0013 | plans/context-circuit-plans/0015-repository-grounding/PLAN.md |
| 0016-worktree-hardening | Worktree hardening and repository-friction feedback | done | 0015, 0007 | plans/context-circuit-plans/0016-worktree-hardening/PLAN.md |
| 0017-external-surface-plan-kind | External surface and the plan publication kind | done | 0004, 0009 | plans/context-circuit-plans/0017-external-surface-plan-kind/PLAN.md |
| 0018-external-surface-thread-kind | External surface: the thread publication kind | done | 0017 | plans/context-circuit-plans/0018-external-surface-thread-kind/PLAN.md |
| 0019-system-design-authoring-skill | cc-system-design authoring skill | done | 0010 | plans/context-circuit-plans/0019-system-design-authoring-skill/PLAN.md |
| 0020-efficiency-ledger | Harness efficiency ledger (dimension D) | done | 0011 | plans/context-circuit-plans/0020-efficiency-ledger/PLAN.md |
| 0021-context-references | External-service reference knowledge category (context/references/) | done | 0003 | plans/context-circuit-plans/0021-context-references/PLAN.md |
| 0022-worker-brief-placement | Relocate the runtime worker-brief template out of the workspace root | done | 0015, 0010 | plans/context-circuit-plans/0022-worker-brief-placement/PLAN.md |
| 0023-design-layout-grouping | Generalize the system-design middle path segment to a grouping dimension | done | 0019 | plans/context-circuit-plans/0023-design-layout-grouping/PLAN.md |
| 0024-worker-role-naming | Unify the implementing execution role name on worker | done | 0005, 0015 | plans/context-circuit-plans/0024-worker-role-naming/PLAN.md |
| 0025-dist-build-version | Version and clean the local dist build from the template source of truth | done | 0010 | plans/context-circuit-plans/0025-dist-build-version/PLAN.md |
