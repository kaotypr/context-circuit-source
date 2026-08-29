---
id: BRIEF-001
title: Move the promoted brief to wrapper/runtime and re-point lookup and manifest
repositories: [context-circuit-source]
paths: [scripts/release-artifact.sh, scripts/release-manifest.txt, wrapper/runtime/engine.sh, agents/coordinator.md, .agents/skills/cc-execute/SKILL.md, .agents/skills/cc-run-stack/SKILL.md]
depends_on: []
acceptance: [BRIEF-AC-01, BRIEF-AC-02]
verification: [BRIEF-VT-01, BRIEF-VT-02]
---

## Intended behavior

The shipped template workspace root stops carrying the internal runtime brief
template. The promoted copy is placed under `wrapper/runtime/`, beside the engine
that consumes it, so the user-facing root holds only the recognizable entry docs.
Repository grounding is unaffected: the same template is assembled, filled, and
preflighted; only where the unfilled template lives changes.

## Concrete change and affected surfaces

In the assembler (`scripts/release-artifact.sh`) copy the brief from
`wrapper/adapters/` to `wrapper/runtime/` in the staged tree instead of the
root. In the manifest (`scripts/release-manifest.txt`) mark the brief required at
`wrapper/runtime/`, keeping the leak-assertion on `wrapper/adapters/`. In the
engine (`wrapper/runtime/engine.sh`, the brief-assemble lookup) prefer the new
`wrapper/runtime/` location, then fall back to `wrapper/adapters/` for source
checkouts. Reword any prose that implied a root path
(`agents/coordinator.md`, the `cc-execute` / `cc-run-stack` skills); the
assemble command interface is unchanged.

## Inputs and outputs

Inputs: the existing promotion, manifest, and lookup wiring. Outputs: a brief
that ships under `wrapper/runtime/` with the source-checkout fallback preserved.
Source of truth (`wrapper/adapters/`) and INV-GROUND semantics are unchanged.

## Context references

- sources/system-design/context-circuit/v0.6.1/writer-brief-placement/design.md

## Risks and open questions

- A harness assertion that inadvertently pins the old root path would break;
  confirm the grounding cases assert content and manifest, not on-disk location.

## Stop conditions

- Do not move or rename the source of truth beyond what `0024-worker-role-naming`
  owns; do not change the assemble command interface or grounding behavior.
