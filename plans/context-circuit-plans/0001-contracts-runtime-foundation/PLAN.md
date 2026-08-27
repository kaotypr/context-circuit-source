# 0001 — Contracts, runtime foundation, and host adapters

- **Plan ID:** `0001-contracts-runtime-foundation`
- **Status:** done
- **Repository:** `context-circuit-source`
- **Depends on:** none (root of the build)
- **Owns (invariants):** INV-RUNTIME-01/02, INV-OWN-01, INV-SEC-01, INV-HOST-01, INV-SKILL-01, INV-COMMIT-01

## Original request

Retroactive plan reconstructing the foundation increment of Context Circuit as
if built from an empty repository. Source design:
`sources/system-design/context-circuit/v0.5/core/` (docs 01, 06, 08, 09, 16).

## Objective

Stand up the layer everything else sits on: a small, host-neutral, deterministic
runtime; the one-rule-one-owner invariant index and authority order; the safety
spine; the controlled vocabulary; and the host-adapter/role surface.

## Desired behavior

- Intelligence lives in roles and skills; determinism lives in the runtime. The
  runtime never interprets knowledge, writes plans, routes conversation, or ships.
- Every product rule has exactly one owning file. Other files cite the invariant
  id and never restate the rule (no parallel policy).
- The root session is always the coordinator; a task/subagent maps only to the
  single worker or the independent verifier.
- Host identity/version/capability/permission/provider status is bounded,
  provider-neutral evidence that authorizes nothing.

## Scope / deliverables

- `wrapper/contracts/invariants.yaml` — the catalog and the concern→owner map.
- `wrapper/adapters/WORKFLOW.md` — authority order + conversational contract.
- `wrapper/adapters/AGENTS.md` — safety spine + shared host-adapter surface.
- `wrapper/runtime/engine.sh` — deterministic skeleton: safe path/id checks,
  atomic writes, digests, the one-writer exclusive-create lock, exclusion list.
- `agents/{coordinator,writer,verifier}.md` — the three role boundaries.
- `docs/terminology.md` + `context/TERMINOLOGY.md` — plain-language projection.

## Non-goals

- Any lifecycle capability (planning, execution, verification, delivery) — those
  are later plans that reference these owners.
- Provider launch code, prompts, or any network/PR/merge/push logic in the runtime.

## Product Knowledge grounding

- `architecture` (`context/ARCHITECTURE.md`) — layered architecture.
- `decisions` (`context/DECISIONS.md`) — determinism/intelligence split.
- `terminology` (`context/TERMINOLOGY.md`) — controlled vocabulary.
- `invariants` (`wrapper/contracts/invariants.yaml`) — the index itself.

## Tasks

1. **FND-001** — Invariant index, authority order, safety spine, terminology.
2. **FND-002** — Deterministic runtime skeleton (paths/ids, atomic writes,
   digests, one-writer lock, exclusion list).
3. **FND-003** — Host adapters, three roles, skill-packet convention.

## Acceptance & verification

- Every rule resolves to one owner; host evidence authorizes nothing;
  runtime writes are atomic and carry no provider/PR logic; skills resolve by path.
- `sh test/contracts/test-contracts.sh`, `sh test/runtime/test-runtime.sh`,
  `sh test/security/test-boundaries.sh`.

## Assumptions / risks

- Assumes POSIX `sh`; the runtime must stay host-neutral.
- Risk: rule duplication drifting into role/skill files — mitigated by the owner
  map and the contracts test asserting single ownership.

## Delivery notes

Foundation for all subsequent plans; nothing ships or executes from this plan on
its own.
