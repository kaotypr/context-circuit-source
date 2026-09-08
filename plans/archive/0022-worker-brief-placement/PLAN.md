# 0022 — Relocate the runtime worker-brief template

- **Plan ID:** `0022-worker-brief-placement`
- **Status:** done
- **Repository:** `context-circuit-source`
- **Depends on:** `0015-repository-grounding`, `0010-release-assembly-and-publication`
- **Owns:** the promoted brief template's *location* (not its name)

## Original request

Retroactive plan for the v0.6.1 `writer-brief-placement` scope, authored as if
from an empty repo. Source design:
`sources/system-design/context-circuit/v0.6.1/writer-brief-placement/design.md`.
(The brief's name is `worker-brief` after `0024-worker-role-naming`; this plan
owns only where the promoted copy lives.)

## Objective and desired behavior

- Keep the shipped `context-circuit-template` workspace root clean of internal
  runtime machinery. The brief template is a runtime-only artifact a user should
  never open; promoted to the bare root beside the real entry docs
  (`README.md`, `AGENTS.md`, `CLAUDE.md`, `WORKFLOW.md`) it reads as clutter.
- Relocate the promoted copy to `wrapper/runtime/`, beside `engine.sh`, the
  runtime code that consumes it. Re-point the engine lookup and the release
  manifest to the new location.
- The `writer-brief-placement`/`worker-brief` **source of truth is unchanged**:
  it stays under `wrapper/adapters/`; only the promoted, shipped runtime copy's
  location moves.

## Constraints and non-goals

- The engine prefers the new `wrapper/runtime/` copy and **falls back** to
  `wrapper/adapters/` for source checkouts (which have no promoted root copy).
- **No user-visible grounding behavior change**; the assembled brief, preflight,
  and grounding execution are identical.
- **No core contract bump.**
- **Non-goal:** the brief's *name* — owned by `0024-worker-role-naming`.

## Product Knowledge grounding

- `invariants` (`wrapper/contracts/invariants.yaml`) — INV-GROUND ownership of
  the brief source of truth, unchanged.

## Tasks

1. **BRIEF-001** — move the promoted brief to `wrapper/runtime/`; re-point the
   engine lookup and the manifest; reword root-path prose.

## Acceptance & verification

- Shipped brief required at `wrapper/runtime/`, not the root; engine prefers it
  with the `wrapper/adapters/` fallback intact; grounding unchanged.
- `sh test/release/test-release.sh`, `sh test/grounding/test-grounding.sh`.

## Assumptions, open questions, risks

- Overlaps `0024-worker-role-naming` on the same file (that plan renames it to
  `worker-brief.md`); the two converge on `wrapper/runtime/worker-brief.md`.
- Harness grounding cases assert behavior, not the template's on-disk location;
  confirm none pins the old root path.

## Expected commits and delivery notes

Assembler + manifest + engine-lookup edits and minor prose; source of truth and
INV-GROUND semantics untouched.
