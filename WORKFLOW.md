# Context Circuit source workflow

Implementation proceeds directly from
`sources/context-circuit-design/DESIGN-SPEC.md` in FINAL-001 through FINAL-009.
This source workflow is not the product lifecycle and does not require plan
approval, run-plan, run-stack, finish-plan, or cleanup ceremonies.

Each phase owns its files, updates `test/baselines/implementation-log.md`, and
runs the semantic tests that exist at that point. Preserve the dependency order:
baseline → wrapper/template → contracts/router → human plan artifacts →
runtime/roles → gates/operations → compatibility → semantic suites → final A/B.

The released product's lifecycle is specified by `wrapper/adapters/WORKFLOW.md`
and owned by `wrapper/contracts/`. The source workflow only governs safe
maintainer changes to that product.

Before changing a rule, identify its canonical owner and update the semantic
fixture that proves it. Keep changes offline and credential-free. Release
assembly is staged and inspectable; publication, deployment, merge, and
destructive cleanup remain outside this task.
