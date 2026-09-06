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
| 0001-conversation-spec-library | Establish the conversation-spec library as the living harness specification | draft | — | 0001-conversation-spec-library/ |
| 0002-intent-id-width-migration | Change intent identifiers to three digits | draft | — | 0002-intent-id-width-migration/ |
| 0003-decomposition-rule | Define the trace-to-plan decomposition rule | draft | — | 0003-decomposition-rule/ |
| 0004-handoff-representation | Carry decomposition through the plan handoff | draft | 0003-decomposition-rule | 0004-handoff-representation/ |
| 0005-decomposition-acceptance | Prove one-plan and stacked-plan derivation | draft | 0004-handoff-representation | 0005-decomposition-acceptance/ |
| 0006-intent-detail | Intent detail via the system-design skill | draft | — | 0006-intent-detail/ |
| 0007-drop-code-review-graph | Drop the leftover .code-review-graph folder | draft | — | 0007-drop-code-review-graph/ |
| 0008-approval-then-trace | Spawn the tracer immediately after intent approval | draft | — | 0008-approval-then-trace/ |
| 0009-intent-md-status | Keep INTENT.md status truthful after a feasible tracer | draft | — | 0009-intent-md-status/ |
| 0010-plans-after-trace | Write plans as soon as the tracer is feasible | draft | 0008-approval-then-trace | 0010-plans-after-trace/ |
| 0011-one-repo-per-plan | One plan, one repository, one worker | draft | — | 0011-one-repo-per-plan/ |
| 0012-explore-session-path | Human-named Explore worktrees under .runtime/explore | draft | — | 0012-explore-session-path/ |
| 0013-cleanup-explore-worktrees | Runtime cleanup also removes Explore worktrees | draft | 0012-explore-session-path | 0013-cleanup-explore-worktrees/ |
| 0014-pair-role-tiering | Honor role-tiering.local.yaml in Explore | draft | — | 0014-pair-role-tiering/ |
| 0015-tracer-role-tiering | Support tracer in role-tiering.local.yaml | draft | — | 0015-tracer-role-tiering/ |
| 0016-in-place-knowledge-updates | Update Product Knowledge in place when a plan is marked done | draft | — | 0016-in-place-knowledge-updates/ |
