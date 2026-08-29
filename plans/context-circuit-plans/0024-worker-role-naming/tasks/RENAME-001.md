---
id: RENAME-001
title: Rename the role and brief to worker across all live surfaces
repositories: [context-circuit-source]
paths: [agents/worker.md, agents/coordinator.md, agents/verifier.md, wrapper/contracts/invariants.yaml, wrapper/runtime/engine.sh, wrapper/adapters/worker-brief.md, .agents/skills/cc-execute/SKILL.md, .agents/skills/cc-run-stack/SKILL.md, scripts/release-artifact.sh, scripts/release-manifest.txt, context/TERMINOLOGY.md, docs/terminology.md]
depends_on: []
acceptance: [RENAME-AC-01, RENAME-AC-02]
verification: [RENAME-VT-01, RENAME-VT-02]
---

## Intended behavior

The implementing execution role is named `worker` on every live product surface,
and `writer` no longer appears as a synonym for it. The change is purely a
rename: what the role does, who may act, and every gate are unchanged; rule IDs
and semantics are untouched.

## Concrete change and affected surfaces

Rename the role file `agents/writer.md` → `agents/worker.md` and point the
invariants owner-map `worker_role:` value at it (the key was already `worker_role`).
Reword the split invariant prose to `worker` — INV-OWN-01/02 ("one-worker lock",
"at most one active worker") and INV-GROUND-01/02/03 ("the worker reads and
honors …") — leaving IDs and semantics unchanged. Rename the brief to the worker
brief: the concept in terminology, the file `writer-brief.md` → `worker-brief.md`
(source of truth under `wrapper/adapters/`), the engine function
`cc_writer_brief_assemble` → `cc_worker_brief_assemble`, and the
`writer-brief-assemble` command → `worker-brief-assemble`. Fix the glossaries
(`context/TERMINOLOGY.md`, `docs/terminology.md`) to define `worker` directly and
retitle the "writer brief" rows. Update the release manifest/assembler brief
path, the `cc-execute` / `cc-run-stack` assemble-command references, and any
`writer` prose in `agents/coordinator.md` / `agents/verifier.md`. Update test and
harness fixture paths, function names, and asserted strings that carry `writer`.

## Inputs and outputs

Inputs: the dominant `worker` term already in the schema and conversational
voice. Outputs: one name per concept across all live surfaces; the brief's
content, assembly, and preflight are unchanged.

## Context references

- sources/system-design/context-circuit/v0.6.1/worker-role-naming/design.md

## Risks and open questions

- A residual `writer` alias anywhere in a live surface reintroduces the split; a
  final grep across live product surfaces must come back clean.
- Overlaps the brief file with `0022-worker-brief-placement`; converge on
  `wrapper/runtime/worker-brief.md` (that plan owns the location, this the name).

## Stop conditions

- Do not change any rule ID, semantics, authority, or gate.
- Do not rewrite historical plans or prior version designs that say "writer".
