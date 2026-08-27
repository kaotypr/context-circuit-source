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
| Delivery | Opening a pull request, merging, or pushing — always a separate, explicit action. "Publish" is not a delivery word; it names the external surface. |
| Connected repository | A repository registered in the workspace and resolved to a local checkout. |
| Host-blocked | A state where the environment cannot run a required step, so the coordinator reports it and preserves the work rather than faking it. |
| Plan stack | A named set of approved plans executed in one run (v0.6 run-stack); ordered and overlapped safely with no new authority. |
| Plan dependency | Inter-plan ordering (`plan_dependencies`); declaring it makes a plan `schema_version: 2`. Distinct from a task's `depends_on`. |
| Path lease | A `(repository, path-region)` reservation extending the one-writer lock; overlapping plans serialize, disjoint ones run together. |
| Execution base / integration base | The commit a plan builds on: anchor tip, a predecessor branch (stack), or a runtime-authored integration merge (≥2 predecessors). |
| Drift guard | Rebase-onto-current-anchor-tip + re-verify before a pull request when the recorded base has diverged. |
| Repository grounding | The worker honoring the target repository's own agent guidance, discovered live from the worktree. |
| Grounding manifest / writer brief | The discovered guidance (files, skills, environment) for one execution; the assembled instructions handed to the worker. |
| System design | A structured source describing the shape of a change, authored via `cc-system-design`; a source, not a lifecycle stage. |
| Publication | A user-declared, manually-triggered pipeline (`publication/<subject>-<provider>/`, created on first use) that publishes Context Circuit data to an external system; orthogonal to the core workflow (v0.6 external-surface). |
| Publish / `cc-publish` | Sending workspace data outward via the manual `cc-publish` skill; reserved for the external surface. Git delivery is "push" / "open a pull request", never "publish" — the two never share a word. |
| Publication kind | What a publication publishes and how `cc-publish` realizes it: `plan` (a plan and its tasks → a tracker) or `thread` (a plan's open questions → a chat discussion). Record shapes are per-kind. |
| Publication `instructions` | Optional per-publication free-text authoring guidance (the language to author in, tone/phrasing, and which optional provider fields to enrich by estimation); guides wording and optional field values only, never a boundary. |

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

Accepted from proposal `0001-terminology-glossary` on 2026-08-24. Extended
2026-08-27 from proposal `0022-change-terminology` with the v0.6 terms (plan
stack, path lease, execution/integration base, drift guard, repository grounding,
grounding manifest, writer brief, system design), mirrored in `docs/terminology.md`.
Extended 2026-08-28 from proposal `0025-change-terminology`: the Delivery row drops
"publishing" (git delivery no longer uses the word), and the external-surface terms
(publication, publish/`cc-publish`, publication kind, publication `instructions`) are
added and mirrored in `docs/terminology.md`.
