# 0017 — External surface and the plan publication kind

- **Plan ID:** `0017-external-surface-plan-kind`
- **Status:** done
- **Repository:** `context-circuit-source`
- **Depends on:** `0004-planning`, `0009-delivery`
- **Owns (invariants):** INV-EXTERNAL-01, INV-EXTERNAL-02, INV-EXTERNAL-03

## Original request

Retroactive plan for the external surface and its first (plan) kind, built as if
from an empty repo. Source design:
`sources/system-design/context-circuit/v0.6/external-surface/`
(`design`, `contracts`, `configuration-and-records`, `publish-plan`).

## Objective and desired behavior

- An opt-in, manually-triggered, config-driven way to publish workspace data
  (plans, tasks) to external systems, fully orthogonal to the core workflow. A
  publication is a peer command, never a phase; nothing runs except an explicit
  human invocation of `cc-publish`.
- A publication reads only the artifact types its `config.yaml` names and writes
  only its own records under `publication/<name>/published/`. Nothing is written
  under `plans/`; no core state or human plan status is mutated; it is
  export-only and credential-free.
- Every external artifact is self-contained — no workspace file name, internal
  path, internal id, or internal vocabulary leaks; only a stable plan/task id
  may appear (as `[<plan-number>]`).
- Provider/network calls go through the host or MCP layer, never the runtime.

## Constraints and non-goals

- Non-goal: the thread kind (0018). This plan ships the backbone + plan kind.
- A wording-only edit vacates "publish" from the git-delivery side.

## Product Knowledge grounding

- `invariants` (`wrapper/contracts/invariants.yaml`), `design-deltas`
  (`context/DESIGN-DELTAS.md`).

## Tasks

1. **EXT-001** — isolation invariants + publish vocabulary split.
2. **EXT-002** — config + record schemas + `publication/` layout.
3. **EXT-003** — `cc-publish` skill (the plan kind).

## Acceptance & verification

- No core phase touches a publication; reads only configured types; idempotent
  re-publish; self-contained external items.
- `sh test/external-surface/test-external-surface.sh`.

## Assumptions, open questions, risks

- No core contract change; a workspace configuring no publication is a v0.5
  workspace plus the availability of `cc-publish`.

## Expected commits and delivery notes

Independent of run-stack / repository-grounding; no `runtime_version` bump.
