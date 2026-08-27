# 0002 — Workspace and repository model

- **Plan ID:** `0002-workspace-repository-model`
- **Status:** done
- **Repository:** `context-circuit-source`
- **Depends on:** `0001-contracts-runtime-foundation`
- **Owns (invariants):** INV-REPO-01/02/03/04 (leans on INV-SEC-01)

## Original request

Retroactive plan for the workspace/repository increment, built as if from an
empty repo. Source design: `sources/system-design/context-circuit/v0.5/core/03-workspace-and-repositories.md`.

## Objective

Make the workspace an agent-oriented coordination layer over one or more Git
repositories, with portable identity and fail-closed, host-local bindings.

## Desired behavior

- **Portable identity** (`workspace.yaml`): logical repo key, optional
  credential-free URL, optional `default_branch`. No paths, no credentials.
- **Local binding** (`repositories.local.yaml`, gitignored): host path plus the
  required `anchor_branch` (the execution base and default PR target).
- `default_branch` is clone/setup guidance only — never silently the anchor.
- The workspace root, when a Git repo, binds as reserved id `workspace` at `.`;
  cloned/initialized repos land under the gitignored `repositories/`.
- Binding resolution is explicit and bounded; anything missing, ambiguous,
  non-Git, traversing, or identity-mismatched fails closed without scanning.

## Scope / deliverables

- Schemas: `workspace.yaml`, `repositories-local.yaml`.
- Runtime: binding resolution/preflight, clone/`git init` bootstrap into
  `repositories/`, deterministic worktree preparation from the anchor tip.
- Skill: `.agents/skills/cc-workspace/SKILL.md` (orient, register, clone, init).

## Non-goals

- Planning or execution over these repos (later plans).
- Any credential handling in workspace files.

## Product Knowledge grounding

- `architecture` (`context/ARCHITECTURE.md`), `conventions`
  (`context/CONVENTIONS.md`), `invariants` (`wrapper/contracts/invariants.yaml`).

## Tasks

1. **WRK-001** — workspace + local-binding schemas.
2. **WRK-002** — binding resolution, bootstrap, worktree prep in the runtime.
3. **WRK-003** — `cc-workspace` orientation and registration skill.

## Acceptance & verification

- No path/credential in portable identity; anchor required; fail-closed binding;
  anchor checkout never written.
- `sh test/contracts/test-contracts.sh`, `sh test/repositories/test-repositories.sh`.

## Assumptions / risks

- Assumes Git is available on the host. Risk: a similarly-named path being
  substituted — prevented by identity-match + fail-closed resolution.

## Delivery notes

Provides the repository substrate that planning and execution map onto.
