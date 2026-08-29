---
id: 0026-repository-grounding-brief-location
target_context_unit: repository-grounding
operation: change
statement: >
  The shipped worker-brief template no longer sits at the workspace root; the
  v0.6.1 writer-brief-placement scope relocated the promoted copy to
  wrapper/runtime/worker-brief.md (beside engine.sh), and the engine lookup
  (cc_worker_brief_assemble) now prefers that path with a wrapper/adapters/
  source-checkout fallback. The repository-grounding domain still describes the
  brief as "promoted to the workspace root" / "worker-brief.md (workspace root,
  shipped)", which is now stale. Correct both references to the runtime location.
  Source of truth (wrapper/adapters/worker-brief.md) and INV-GROUND semantics are
  unchanged; only the promoted, shipped copy's location moved.
evidence_refs:
  - sources/system-design/context-circuit/v0.6.1/writer-brief-placement/design.md
  - context/domains/repository-grounding/README.md:62
  - context/domains/repository-grounding/README.md:83
  - wrapper/runtime/engine.sh:706  # cc_worker_brief_assemble
  - wrapper/runtime/engine.sh:709  # primary lookup: wrapper/runtime/worker-brief.md
  - wrapper/runtime/engine.sh:710  # fallback: wrapper/adapters/worker-brief.md
affected_repositories:
  - context-circuit
affected_commits:
  - 2c2adab  # refactor(worker): relocate the brief to the runtime and unify the role name on worker
related_plan: 0022-worker-brief-placement
confidence: high
status: review-needed
---

# Correct the worker-brief location in repository-grounding

## What is stale

The `repository-grounding` domain page describes the shipped brief template as
living at the workspace root:

- line 62 — "(`worker-brief.md`, promoted to the workspace root) is filled by
  deterministic …"
- line 83 — "Template: `worker-brief.md` (workspace root, shipped)"

## What changed (v0.6.1 · writer-brief-placement, plan 0022)

The workspace root is user-facing surface; the brief is internal runtime
machinery a user should never open. The promoted, shipped copy was moved out of
the root into an unobtrusive machinery location that already ships:
`wrapper/runtime/worker-brief.md`, beside `engine.sh` — the runtime code that
consumes it. `cc_worker_brief_assemble` now prefers that path and falls back to
`wrapper/adapters/worker-brief.md` for this source checkout. The dist artifact
and every current template-harness workspace carry the brief at
`wrapper/runtime/worker-brief.md`.

## Proposed change

Reword the two references so the shipped/promoted location reads
`wrapper/runtime/worker-brief.md` rather than "the workspace root". Leave the
source-of-truth statement (`wrapper/adapters/worker-brief.md`, INV-GROUND-01/03)
and all grounding behavior untouched — only the promoted copy's on-disk location
is corrected.

## Why not more

No behavior, authority, gate, or contract semantics changed; this is a
stale-location correction, not a re-scope. No core contract bump. On acceptance,
refresh the domain's `source_revisions` commit pin (currently `cb84870`) to the
revision carrying `2c2adab`.
