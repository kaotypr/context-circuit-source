# Run-stack — schema and versioning

Each change names its canonical owner; this scope specifies intent, the owner
files carry the rule. No rule is duplicated (v0.5 principle 5.7).

| Owner | Change | Version impact |
| --- | --- | --- |
| `wrapper/contracts/schemas/plan.yaml` | add optional `plan_dependencies` (list of `{id, reason}`); acyclic; entries reference existing plan IDs | instances using it are `schema_version: 2` |
| `wrapper/contracts/schemas/execution.yaml` | add optional `based_on` to the per-repository record; broaden `base_commit` to allow a runtime-authored integration merge | stays `schema_version: 1` |
| `wrapper/contracts/schemas/lease.yaml` **(new)** | path-lease records under `.runtime/locks/paths/` — holder plan, repository, path region, acquired-at, released-at | born `schema_version: 1` |
| `wrapper/contracts/invariants.yaml` | add INV-CONCURRENCY-01 (path leases), INV-CONCURRENCY-02 (base + base-ref rule), extend the delivery rule with the drift guard; add owner-map entries for path leases and base selection (owner: `engine.sh`) and the run-stack action (owner: WORKFLOW / coordinator) | — |
| `wrapper/adapters/WORKFLOW.md` | add the conversational action "Execute plans X…Z / run the ready stack" beside the singular "Execute plan X" | — |
| `wrapper/runtime/engine.sh` | add deterministic functions: acquire/check/release path leases, build the integration base, evaluate readiness, detect base drift | — |
| `wrapper/manifest.yaml` | `runtime_version: 0.5.0 → 0.6.0`; `accepted_schema_versions.plan: [1, 2]`; bump `template_version` (minor, pre-1.0) at publish | — |

## Why plan.yaml bumps to 2 but execution.yaml does not

`plan_dependencies` is **load-bearing**: an older engine that ignored it would
schedule a plan without its dependencies and run it out of order — a silent
correctness failure, not a harmless unknown field. Bumping the instance
`schema_version` to `2` makes a v0.5 engine **refuse** such a plan instead of
mis-scheduling it. Only plans that actually use the field carry `schema_version:
2`; every existing plan stays `1` and is read identically by both engines. The
manifest's `accepted_schema_versions.plan: [1, 2]` admits both.

`based_on` and the broadened `base_commit` live in `execution.yaml`, which is
private runtime evidence written and read by the same engine version that created
it. There is no cross-version reader to mislead, so it stays `schema_version: 1`.
