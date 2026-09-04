# context-circuit-source workflow

This checkout is `context-circuit-source`, the maintainer source repository for
Context Circuit. Implementation proceeds directly from the maintainer design set
under `sources/system-design/context-circuit/` and the source-only ordering
guides under `plans/context-circuit-plans/`. The source is not bound to a
version; it tracks the current design and may lead the version last published to
`context-circuit-template`.

Source changes are made directly in the current active branch. Do not create an
execution or pairing branch, worktree, lease, execution record, or product
candidate for maintainer implementation. The source-only plan files describe
scope and implementation order; they do not invoke the product workspace
lifecycle.

Each phase owns its files and runs the semantic tests that exist at that point.
Preserve the dependency order:
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
