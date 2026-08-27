# context-circuit-source workflow

Implementation proceeds directly from the maintainer design set under
`sources/system-design/context-circuit/` in its implementation order. The source
is not bound to a version; it tracks the current design and may lead the version
last published to context-circuit-template.
This source workflow is not the product lifecycle and does not require plan
approval, run-plan, run-stack, finish-plan, or cleanup ceremonies.

Each phase owns its files, updates `test/baselines/implementation-log.md`, and
runs the semantic tests that exist at that point. Preserve the dependency order:
workspace/plan contract → Product Knowledge → runtime reduction → conversational
adapter → multi-repository execution → independent verification → repair loop →
recovery/delivery → semantic verification. Scoped increments layer on that core
in the order recorded in `plans/context-circuit-plans/INDEX.md`.

The released product's lifecycle is specified by `wrapper/adapters/WORKFLOW.md`
and owned by `wrapper/contracts/`. The source workflow only governs safe
maintainer changes to that product.

Before changing a rule, identify its canonical owner and update the semantic
fixture that proves it. Keep changes offline and credential-free. Release
assembly is staged and inspectable; publication, deployment, merge, and
destructive cleanup remain outside this task.
