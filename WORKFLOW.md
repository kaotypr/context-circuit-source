# context-circuit-source workflow

Implementation proceeds directly from
`sources/context-circuit-v0.5-design/context-circuit-v0.5-design.md` in its
implementation order.
This source workflow is not the product lifecycle and does not require plan
approval, run-plan, run-stack, finish-plan, or cleanup ceremonies.

Each phase owns its files, updates `test/baselines/implementation-log.md`, and
runs the semantic tests that exist at that point. Preserve the dependency order:
workspace/plan contract → Product Knowledge → runtime reduction → conversational
adapter → multi-repository execution → independent verification → repair loop →
recovery/delivery → semantic verification.

The released product's lifecycle is specified by `wrapper/adapters/WORKFLOW.md`
and owned by `wrapper/contracts/`. The source workflow only governs safe
maintainer changes to that product.

Before changing a rule, identify its canonical owner and update the semantic
fixture that proves it. Keep changes offline and credential-free. Release
assembly is staged and inspectable; publication, deployment, merge, and
destructive cleanup remain outside this task.
