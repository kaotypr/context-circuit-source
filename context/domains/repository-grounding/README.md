---
kind: domain
status: accepted
title: Repository grounding
slug: repository-grounding
owners: []
sources: []
source_revisions:
  - wrapper: HEAD
    commit: 2c2adab
    basis: current-wrapper
generated_at: 2026-08-27T00:00:00Z
review_date: 2026-11-27
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
  - docs/getting-started.md
---

# Repository grounding

## Summary

Grounding the worker in the *target repository's own* agent guidance — discovered
live from the execution worktree and delivered through a generated brief, never
hand-authored into each worker prompt. Part of execution setup; owned by the
runtime discovery/hardening/brief functions and the shipped `worker-brief.md`
template, applied by the `cc-execute` and `cc-run-stack` skills and honored by the
`agents/worker.md` worker.

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
- **Worktree hardening.** The runtime detects the toolchain and reports
  `environment: ready | no-toolchain`, so CC-induced execution-environment
  workarounds are eliminated rather than documented (full dependency provisioning
  is a later phase).
- **Deliver, not author (INV-GROUND-03).** A fixed shipped template
  (`wrapper/runtime/worker-brief.md`, shipped beside the runtime) is filled by deterministic
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

- Execute a plan grounded in the repo's own conventions: `docs/getting-started.md`

## Interfaces

- Records: `.runtime/executions/<plan>/<exec>/grounding/<repo>.yaml` (manifest);
  the assembled brief `brief-<repo>.md`; `repository_friction` in the handoff
- Template: `wrapper/runtime/worker-brief.md` (shipped, beside the runtime)

## Data

The grounding manifest (`files[]`, `skills[]{name, description}`, `environment`)
per affected repository, recorded per execution for reproducibility.

## Constraints and edge cases

A doc-less repository yields an empty manifest and an honest brief that says so; a
greenfield repository reports `no-toolchain` and does not block (a scaffold plan
may author the repo's agent docs, which later stacked plans then discover).

## Implementation references

- `wrapper/adapters/worker-brief.md`; `.agents/skills/cc-execute/SKILL.md`,
  `.agents/skills/cc-run-stack/SKILL.md`; `agents/worker.md`, `agents/coordinator.md`
- `wrapper/runtime/engine.sh`: `cc_discover_repo_grounding`, `cc_harden_worktree`,
  `cc_grounding_directive`, `cc_worker_brief_assemble`, `cc_brief_preflight`,
  `cc_skill_desc`
- `wrapper/contracts/schemas/grounding-manifest.yaml`, `execution.yaml` (grounding
  record), `worker-handoff.yaml` (`repository_friction`)
- `wrapper/contracts/invariants.yaml`: INV-GROUND-01, INV-GROUND-02, INV-GROUND-03

## Verification

`sh test/acceptance.sh` (repository-grounding suite); a section in
`agent-harness/test-template-runtime.sh` proves the shipped template carries it;
live scenario `agent-harness/scenarios/11-repo-grounding` (grade.sh PASS incl.
`file_grounded` + `grounding_manifest_recorded`, human-simulator pass).

## Provenance

Authored from the current wrapper at HEAD `cb84870` (implementation `e3e95cd`).
Design source `sources/system-design/context-circuit/v0.6/repository-grounding/`
was named by the accepting request.

## Acceptance notes

Accepted 2026-08-27 from proposal `0017-domain-repository-grounding`.
Updated 2026-08-29 from proposal `0026-repository-grounding-brief-location` (v0.6.1
writer-brief-placement): the shipped brief moved out of the workspace root to
`wrapper/runtime/worker-brief.md`; source of truth (`wrapper/adapters/worker-brief.md`)
and INV-GROUND semantics unchanged. Implementation `2c2adab`.
