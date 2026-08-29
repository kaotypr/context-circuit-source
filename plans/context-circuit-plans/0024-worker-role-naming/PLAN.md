# 0024 — Unify the implementing execution role name on worker

- **Plan ID:** `0024-worker-role-naming`
- **Status:** done
- **Repository:** `context-circuit-source`
- **Depends on:** `0005-approval-and-execution`, `0015-repository-grounding`
- **Owns:** the role's *name* (worker), and the brief's name

## Original request

Retroactive plan for the v0.6.1 `worker-role-naming` scope, authored as if from
an empty repo. Source design:
`sources/system-design/context-circuit/v0.6.1/worker-role-naming/design.md`.

## Objective and desired behavior

- Give the implementing execution role **one name — `worker`** — everywhere it
  appears, retiring `writer` as a second name for the same role. One word for one
  concept, so a reader never learns that "writer" and "worker" mean the same
  thing.
- `worker` is already the dominant, user-facing term (the schema is
  `worker-handoff.yaml`; the conversational voice is "the worker"). Rename
  *toward* worker: it moves the fewest concepts and breaks no established word.
- The brief follows the role: "writer brief" → "worker brief" in concept, file
  (`writer-brief.md` → `worker-brief.md`), engine function
  (`cc_writer_brief_assemble` → `cc_worker_brief_assemble`), and command
  (`writer-brief-assemble` → `worker-brief-assemble`).

## Constraints and non-goals

- **Rename, do not re-scope.** No behavior, authority, or contract semantics
  change; rule **IDs and semantics are unchanged** — only the noun.
- **No core contract bump and no new owner**; the invariants owner map is edited
  in place (one path value, one term).
- **No retroactive history rewrite**: completed maintainer plans and prior
  version designs that say "writer" are historical record and are not edited.
- **Non-goal:** the brief's *location* — owned by `0022-worker-brief-placement`.
  If both land, the result is `wrapper/runtime/worker-brief.md`.

## Product Knowledge grounding

- `invariants` (`wrapper/contracts/invariants.yaml`) — owner-map values and
  INV-OWN / INV-GROUND prose nouns.
- `terminology` (`context/TERMINOLOGY.md`) — the glossary must define `worker`
  directly.

## Tasks

1. **RENAME-001** — rename the role file, brief, engine function, and command;
   reword invariant prose; fix glossaries; update tests and harness.

## Acceptance & verification

- No live product surface names the implementing role "writer"; the brief is the
  worker brief; rule IDs/semantics unchanged; no core contract bump.
- `sh test/acceptance.sh`, `sh test/contracts/test-contracts.sh`.

## Assumptions, open questions, risks

- A residual "writer" in any live surface reintroduces the split; the acceptance
  is a clean grep across live product surfaces.
- Overlaps `0022-worker-brief-placement` on the brief file; the two converge on
  `wrapper/runtime/worker-brief.md`.

## Expected commits and delivery notes

A file rename plus wording edits across roles, contracts, runtime, adapters,
skills, docs, and tests; no rule added, split, or moved to a new owner.
