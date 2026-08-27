# Template harness — v0.6 (efficiency ledger, made real)

The v0.6 delta on the maintainer human-simulated test harness. It does **one
thing**: turn the already-specified but never-measured **dimension D (efficiency
ledger)** into a real, evaluated ledger — by emitting per-action usage telemetry
from the runtime's own run output, defining the units a budget is stated in, and
having the grader actually compare observed usage to each case's `budgets`.

v0.6 is a **delta on v0.5** — read the base first:
[../../v0.5/template-harness/](../../v0.5/template-harness/), especially
[harness-and-evaluation.md](../../v0.5/template-harness/harness-and-evaluation.md)
§4.2 and [roles/grader.md](../../v0.5/template-harness/roles/grader.md) §D, which
this scope changes.

## Reading order

1. [design.md](./design.md) — the problem (D declared but unmeasured; units
   undefined; a turn-unit mismatch), the shape of the fix, principles, and fixed
   decisions.
2. [telemetry.md](./telemetry.md) — where usage comes from (the host runner's own
   run result), how the driver derives and writes it, and per-action attribution.
3. [budgets.md](./budgets.md) — the unit definitions, the case `budgets` schema
   delta, and how budgets are recalibrated from measured runs.
4. [evaluation.md](./evaluation.md) — the dimension-D grader delta (a real numeric
   compare, still soft), the worked ledgers for cases 10 and 11, and acceptance.

## Scope and authority

This is **maintainer test tooling only** and changes **no product surface**: no
`wrapper/` contract, no invariant, no `engine.sh`, no shipped file. The release
manifest already excludes `test/` and `template-harness/`, so nothing here ships
in `context-circuit-template`. It is independent of the three v0.6 product scopes
(run-stack, system-design-stage, repository-grounding) and needs none of their
contract bump; it rides the v0.6 version umbrella only for versioning tidiness.

Where it observes product behavior, the product design and
`wrapper/contracts/invariants.yaml` remain the authority for what "correct" means
(v0.5 template-harness "Authority").
