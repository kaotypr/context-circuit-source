---
kind: domain
status: accepted
title: Repository grounding
slug: repository-grounding
owners: []
sources: []
source_revisions:
  - wrapper: HEAD
    commit: d95bbd7
    basis: current-wrapper
generated_at: 2026-09-10T00:00:00Z
review_date: 2026-12-10
freshness: accepted-from-current-wrapper
assumptions:
  - Repository guidance is referenced from the worktree, never captured into Product Knowledge.
unknowns: []
contradictions: []
acceptance:
  state: accepted
  accepted_at: 2026-08-27
  accepted_by: maintainer
workflows:
  - .context-circuit/docs/getting-started.md
---

# Repository grounding

## Summary

Grounding the worker in the *target repository's own* agent guidance — scanned
live from the execution worktree and delivered through a generated brief, never
hand-authored into each worker prompt. Part of **execution setup**; owned by the
runtime grounding/hardening/brief functions and the shipped `worker-brief.md`
template, applied by the `cc-execute` and `cc-run-stack` skills and honored by the
`.context-circuit/agents/worker.md` worker.

This is **not** the post-approval **planner**
([tracing and feasibility](../tracing/README.md)). Both phases read the repository,
but they differ: the planner runs *after approval* to judge *whether
the change is buildable* (file map, risks, done-checks, tier signal) and writes
the **plan files**; repository grounding runs *at execution setup, inside the
worktree* to establish *how to write code in this repository* and produces the
**worker brief**. Different timing, subject, and output.

## Scope

Inside: the deterministic discovery scan and its manifest, worktree hardening, the
assembled worker brief and its preflight, the worker's read-and-honor obligation,
precedence, and the repository-friction feedback loop.

Outside: what a worker changes and where (the plan and its scope), verification
([verification](../verification/README.md)), and base selection for a dependent
plan ([run-stack](../run-stack/README.md)).

## Behavior

- **Discovery emits data, not prompt text (INV-GROUND-01, INV-RUNTIME-01).**
  During worktree preparation the runtime scans the worktree for the repository's
  own agent guidance — `AGENTS.md`, `CLAUDE.md`, `.cursor/rules/`,
  `.github/copilot-instructions.md`, and `.agents/skills/*/SKILL.md` (with each
  skill's description) — and records a grounding manifest (`files`, `skills`,
  `environment`) as execution evidence. Scanning the *execution base* means
  bootstrapping propagates through a stack.
- **Worktree hardening.** Before a worker is attached, the runtime overlays the
  bound checkout's gitignored paths into the fresh worktree with copy-on-write or
  a real copy, never a symlink; untracked paths that are not ignored stay out.
  It then detects a toolchain. A copied install tree is kept only when it matches
  this worktree's lockfile; otherwise the runtime runs that toolchain's frozen
  install. The grounding-manifest schema owns the resulting environment value:
  preparation emits `ready` only after these steps succeed, while a setup failure
  blocks execution rather than reaching a worker.
- **Deliver, not author (INV-GROUND-03).** A fixed shipped template
  (`.context-circuit/wrapper/runtime/worker-brief.md`, shipped beside the nested
  runtime) is filled by deterministic
  slot substitution from the manifest and the plan; the coordinator adds only a
  one-line task focus and delivers it verbatim. A preflight refuses a brief that
  omits or leaves unfilled the required repository-grounding section.
- **Precedence (INV-GROUND-02).** The brief's scope and safety rules are
  authoritative on *what* and *where*; the repository's guidance is authoritative
  on *how* to write code there, within that scope. On a scope or safety conflict
  the worker stops and reports; repository guidance never overrides a CC rule.
- **Reference, not capture.** No per-repo profile and no `plan.yaml` field. When
  the repository's guidance did not cover something the worker needed, it returns
  a `repository_friction` note in the handoff, reconciled into a proposal on the
  repository's own agent docs — never a CC-side profile.

## Workflows

- Execute a plan grounded in the repo's own conventions: `.context-circuit/docs/getting-started.md`

## Interfaces

- Records: `.runtime/executions/<plan>/<exec>/grounding/<repo>.yaml` (manifest);
  the assembled brief `brief-<repo>.md`; `repository_friction` in the handoff
- Template: `.context-circuit/wrapper/runtime/worker-brief.md` (shipped, beside the runtime)

## Data

The grounding manifest (`files[]`, `skills[]{name, description}`, `environment`)
per affected repository, recorded per execution for reproducibility.

## Constraints and edge cases

A doc-less repository yields an empty manifest and an honest brief that says so; a
greenfield repository reports `no-toolchain` and does not block (a scaffold plan
may author the repo's agent docs, which later stacked plans then discover).

## Implementation references

- `.context-circuit/wrapper/adapters/worker-brief.md`; `.agents/skills/cc-execute/SKILL.md`,
  `.agents/skills/cc-run-stack/SKILL.md`; `.context-circuit/agents/worker.md`, `.context-circuit/agents/coordinator.md`
- `.context-circuit/wrapper/runtime/engine.sh`: `cc_overlay_ignored`,
  `cc_provision_worktree`, `cc_toolchain_matches`, `cc_toolchain_install`,
  `cc_discover_repo_grounding`, `cc_harden_worktree`, `cc_grounding_directive`,
  `cc_worker_brief_assemble`, `cc_brief_preflight`, `cc_skill_desc`
- `.context-circuit/wrapper/contracts/schemas/grounding-manifest.yaml`, `execution.yaml` (grounding
  record), `worker-handoff.yaml` (`repository_friction`)
- `.context-circuit/wrapper/contracts/invariants.yaml`: INV-GROUND-01, INV-GROUND-02, INV-GROUND-03

## Verification

`sh test/grounding/test-grounding.sh` and `sh test/acceptance.sh`
(repository-grounding suite).

## Acceptance notes

Accepted 2026-08-27.
Updated 2026-08-29: the shipped brief moved out of the workspace root to
`.context-circuit/wrapper/runtime/worker-brief.md`; source of truth
(`.context-circuit/wrapper/adapters/worker-brief.md`) and INV-GROUND semantics
unchanged.
Updated 2026-09-10: fresh worktrees receive the ignored-path overlay and
lockfile-correct provisioning before grounding; setup failures block rather than
asking a worker to compensate.
