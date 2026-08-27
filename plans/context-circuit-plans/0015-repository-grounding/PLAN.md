# 0015 — Repository grounding: discovery, manifest, and generated writer brief

- **Plan ID:** `0015-repository-grounding`
- **Status:** done
- **Repository:** `context-circuit-source`
- **Depends on:** `0005-approval-and-execution`, `0013-run-stack-substrate`
- **Owns (invariants):** INV-GROUND-01, INV-GROUND-02, INV-GROUND-03

## Original request

Retroactive plan for repository grounding, built as if from an empty repo.
Source design: `sources/system-design/context-circuit/v0.6/repository-grounding/`
(`design`, `discovery-and-manifest`, `schema-and-contracts`, `writer-brief-template`).

## Objective and desired behavior

- Before implementing, the writer reads and honors the target repository's own
  agent guidance (`AGENTS.md`, `CLAUDE.md`, `.cursor/rules`,
  `.github/copilot-instructions.md`, `.agents/skills/*/SKILL.md`), discovered
  live from the execution worktree by a deterministic runtime scan.
- The runtime emits **data** — a grounding manifest recorded as execution
  evidence — never prompt text.
- The writer brief's fact sections are a fixed template filled by deterministic
  slot substitution; the coordinator delivers the brief and authors only a
  one-line task focus. A preflight refuses a brief missing the grounding slot.
- Precedence: the brief's scope/safety rules win on *what/where*; the
  repository's guidance is authoritative on *how* within that scope.

## Constraints and non-goals

- Non-goal: worktree hardening and the friction feedback loop (0016). No
  `plan.yaml` field and no per-repo profile are introduced.

## Product Knowledge grounding

- `invariants` (`wrapper/contracts/invariants.yaml`), `design-deltas`
  (`context/DESIGN-DELTAS.md`).

## Tasks

1. **RGD-001** — live discovery + grounding manifest.
2. **RGD-002** — writer-brief template, slot-fill, preflight.
3. **RGD-003** — worker + coordinator grounding contract.

## Acceptance & verification

- Discovery emits a manifest not prompt text; empty manifest valid; missing
  grounding slot refused; repo guidance never overrides CC safety/scope.
- `sh test/grounding/test-grounding.sh`.

## Assumptions, open questions, risks

- Rides the coordinated v0.6 bump owned by 0013; introduces no schema-version
  change of its own.

## Expected commits and delivery notes

Delivered in the same v0.6 release as run-stack.
