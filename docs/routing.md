# Two-stage routing

The router is implemented once in `wrapper/runtime/engine.sh` and specified by
`wrapper/contracts/routes.yaml`. Skills and role files are adapters; they do
not maintain alternative route tables.

Stage A reads only Tier 0 and selects one probe. Stage B reads that probe and
emits exactly one normalized decision:

```yaml
intent: execute-plan
session_kind: root
phase: execution-preflight
probe: execution-preflight
eligibility: ready
capability: execute-plan
authorization: explicitly-requested
reason_codes: [ROUTE_SELECTED]
context_set: execution-preflight
human_gate: none
```

Eligibility and authorization are separate fields. Blockers have precedence;
an approved plan discovered during orientation produces a recommendation, not a
lease. A malformed, unsafe, stale, or foreign-owned state routes to recovery
before mutation.

The router never pretends to know unread state. Budget overruns use the reason
contract in `wrapper/contracts/context-sets.yaml`.
