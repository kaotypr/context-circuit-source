# Run-stack — worked examples and acceptance

These two traces are the canonical acceptance fixtures for the run-stack scope.
Both use the request *"execute approved plans 0001 until 0010"*.

## Multi-repository stack (three repositories)

Ten plans across a backend and two SPAs, with a shared library, an events file two
domains touch, and two SPA plans that both edit one menu file.

- **Wave 1** — three scaffolds run in parallel (different repositories).
- **Wave 2** — the shared money library, a serialization point everything imports.
- **Wave 3** — order and kpi run in parallel **inside the backend** on disjoint paths.
- **Wave 4** — commission (stacked on order), cms-order-ui (stacked on its
  scaffold, gated by order across repositories), and agent-order-flow run together.
- The two cms plans both declare the menu file; the path lease serializes the
  second behind the first though no dependency links them.
- commission fails after three attempts; evidence is preserved and only its
  dependent (commission-ui) is held; the other eight plans are unaffected.
- **Outcome:** eight verified and ready to deliver, one failed, one held.

## Single-repository stack (one modular monolith)

The same ten plans, all in one backend service. Parallelism comes only from
disjoint paths.

- **Wave 1** — the scaffold holds a repository-wide lease — exclusive.
- **Wave 2** — money and events run in parallel (disjoint `internal/pkg/` paths).
- **Wave 3** — order, partner, kpi run in parallel (three disjoint module
  directories); order uses an **integration base** merged from money + events.
- **Wave 4** — commission and claim run in parallel; claim uses an integration
  base from order + kpi; claim fails.
- The module-registry plan is the single writer of the shared wiring file and
  depends on every domain; because claim failed, the registry (and the e2e plan
  below it) is held — the single-repository blast radius (see
  [scheduling.md](./scheduling.md)).
- Delivery is one serialized rebase train into the single anchor branch.
- **Outcome:** seven verified, one failed, two held.

## Acceptance criteria

The run-stack scope is acceptable when all v0.5 acceptance criteria still hold and:

1. A single request executes a named set of approved plans; an unapproved plan in
   the set is refused without blocking the rest.
2. Plans with disjoint leases in the same repository execute concurrently, each in
   its own worktree.
3. Plans that overlap on a path region are serialized, even when no dependency
   links them, and neither produces a merge conflict.
4. A plan runs only after all its plan dependencies are verified.
5. A dependent with one same-repo predecessor is based on that predecessor's
   branch; with two or more, on a runtime-authored integration merge; with only
   cross-repo dependencies, on its own anchor tip.
6. The integration merge is authored by the runtime before the worker starts and
   is not counted as a worker attempt; a base that cannot be built cleanly is
   blocked, not failed.
7. A failed plan preserves evidence and holds only its descendants; unrelated
   verified plans are unaffected.
8. Every plan in a stack is still executed by one worker and one independent
   verifier, with the three-failure limit intact.
9. A plan whose base diverged before delivery is rebased onto the current anchor
   tip and re-verified before its pull request.
10. A plan using `plan_dependencies` is `schema_version: 2`; a v0.5-era plan
    without it stays `schema_version: 1`; the manifest accepts `[1, 2]`; a v0.5
    engine refuses a `schema_version: 2` plan.
11. Path leases and integration bases are produced by deterministic runtime
    functions; the runtime contains no scheduling heuristic and no provider launch logic.
12. The semantic acceptance suite covers both worked examples above: concurrent
    execution, lease serialization, single- and multi-predecessor bases, failure
    containment, held descendants, and the delivery drift guard.
