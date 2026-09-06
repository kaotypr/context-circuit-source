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
