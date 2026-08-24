# Terminology

Accepted glossary for the Context Circuit source project. The source project's
domain is Context Circuit itself, so this page accepts the product vocabulary
shipped to workspaces as `docs/terminology.md`. Term *meaning* and the
*user-facing translation* are owned there, and their exact authority is settled
by the runtime contracts under `wrapper/contracts/`; this page records that the
source project accepts that vocabulary as Product Knowledge and makes it
retrievable.

## Product and runtime terms

| Term | Meaning |
| --- | --- |
| Wrapper | Accepted synonym for the universal project workspace — the Context Circuit product. The directory `wrapper/` holds its shipped layer (runtime, contracts, adapters). |
| Product Knowledge | Accepted, agent-oriented understanding of the project, stored as indexed, human-readable units under `context/`. |
| Context unit / context index | One knowledge unit; the retrieval catalog (`context/INDEX.md`) mapping concepts and aliases to units. |
| Plan / plan status | Human-reviewed intent for an outcome; its human-owned `draft`/`approved`/`done` state. |
| Execution | One runtime attempt to implement an approved plan, with one worker and one independent verifier. |
| Worker / verifier | The single writer role for an execution; the independent read-only role that checks the worker's latest commits. |
| Repair attempt | A new worker commit plus a new independent check after a failed verification. |
| Completion | The human decision to mark a plan done after a verified execution. |
| Archive / restore | Setting a plan aside, or bringing it back, without changing its status. |
| Delivery | Opening a pull request, merging, pushing, or publishing — always a separate, explicit action. |
| Connected repository | A repository registered in the workspace and resolved to a local checkout. |
| Host-blocked | A state where the environment cannot run a required step, so the coordinator reports it and preserves the work rather than faking it. |

## User-facing translation

The coordinator reports actions by their effect and never exposes internal
mechanism unless a user asks for diagnostics. The canonical internal → plain
mapping is the "say the effect, not the mechanism" table in `docs/terminology.md`
(worktree, anchor branch, binding, execution branch, verifier, worker,
host-blocked, delivery, archive/restore, and internal file names), whose
authority is settled by the runtime contracts under `wrapper/contracts/`.

## Provenance

- `docs/terminology.md` — the shipped product glossary and internal → user-facing
  translation table; the coordinator references it, and its authority is settled
  by the runtime contracts under `wrapper/contracts/`.
- `agents/coordinator.md` — the role that applies the translation.
- Design material under `sources/` is historical maintainer input, not a live
  owner, and is read only when a request names it.

Accepted from proposal `0001-terminology-glossary` on 2026-08-24.
