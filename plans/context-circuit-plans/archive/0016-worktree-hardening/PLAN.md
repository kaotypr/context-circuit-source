# 0016 — Worktree hardening and repository-friction feedback

- **Plan ID:** `0016-worktree-hardening`
- **Status:** done
- **Repository:** `context-circuit-source`
- **Depends on:** `0015-repository-grounding`, `0007-completion-and-reconciliation`
- **Owns:** worktree hardening + friction feedback (extends INV-GROUND-01/03)

## Original request

Retroactive plan for worktree hardening and the friction loop, built as if from
an empty repo. Source design:
`sources/system-design/context-circuit/v0.6/repository-grounding/worktree-hardening.md`.

## Objective and desired behavior

- Harden the isolated execution worktree so environment workarounds (symlinked
  deps, dep-verify flags, hook bypasses) become unnecessary and are never
  hand-documented in the brief.
- Detect the toolchain from lockfiles/manifests, provision dependencies by
  reusing the bound checkout's installed deps, handle commit hooks in the
  isolated context, and set the manifest `environment` field
  (`ready | no-toolchain`).
- Route worker "guidance did not cover X" as a `repository_friction` handoff
  field that becomes a proposal on the repository's own agent docs — never core
  Product Knowledge.

## Constraints and non-goals

- Non-goal: adding a `plan.yaml` field or per-repo profile.

## Product Knowledge grounding

- `invariants` (`wrapper/contracts/invariants.yaml`), `design-deltas`
  (`context/DESIGN-DELTAS.md`).

## Tasks

1. **WTH-001** — toolchain detection, provisioning, environment state.
2. **WTH-002** — repository-friction feedback loop.

## Acceptance & verification

- Hardened worktree needs no hand-authored workaround; environment field honest;
  friction becomes a repo-docs proposal.
- `sh test/grounding/test-grounding.sh`.

## Assumptions, open questions, risks

- Risk: an ecosystem without a supported toolchain — reported honestly as
  `no-toolchain` rather than a silent workaround.

## Expected commits and delivery notes

The heaviest grounding piece; lands after 0015 in the same v0.6 release.
