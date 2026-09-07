---
kind: domain
status: accepted
title: Source release and upgrade
slug: source-release-and-upgrade
owners: []
sources: []
source_revisions:
  - wrapper: HEAD
    commit: 6921a31
    basis: current-wrapper
generated_at: 2026-08-24T00:00:00Z
review_date: 2026-11-24
freshness: accepted-from-current-wrapper
assumptions:
  - This maintainer checkout assembles a distributable product it is not itself.
unknowns: []
contradictions: []
acceptance:
  state: accepted
  accepted_at: 2026-08-24
  accepted_by: maintainer
workflows:
  - .context-circuit/docs/release.md
---

# Source release and upgrade

## Summary

How the maintainer source assembles the distributable product and how an upgrade
preserves a user's workspace. Route release-assembly, ships-vs-never-ships,
template-seed, and upgrade/migration questions here. This is distinct from the
runtime `delivery` domain (a user's pull request/merge/push for their own repos).

## Scope

Inside: the three named identities; template assembly and its exclusion
boundary; the mandatory template `.gitignore`; the upgrade preservation set and
migration-needed reporting; release-as-gate; the agent-harness isolation
contract.

Outside: the runtime lifecycle (plan/execute/verify/complete) and a user's
delivery actions.

## Behavior

Three identities stay distinct: `context-circuit-source` (this maintainer
checkout), `context-circuit-template` (the distributable universal project
workspace — the wrapper a user works in), and the instantiated project
workspace.

Assembly ships the template-owned files plus the blank `template/` seed and
never ships: the source design, maintainer plans and logs, source-only tests and
evidence, source `.runtime/`, `agent-harness/`, credentials, local bindings,
connected repositories, project Product Knowledge, customer plans, or
archived-plan contents.

The template's root `.gitignore` must exclude at least `/repositories/`,
`/repositories.local.yaml`, `/.runtime/`, and `.code-review-graph/`.

An upgrade may replace template-owned files, runtime code, contracts, and role
guidance, but must preserve workspace-owned files: project identity, accepted
Product Knowledge, sources and provenance, plans and their status, local
bindings, connected repositories, runtime evidence, and active worktrees. When a
template change alters the meaning of a plan, context, or runtime record, the
upgrade reports migration-needed and preserves the old state.

The template artifact's version is its `template_version` (`.context-circuit/wrapper/manifest.yaml`),
never the product's internal `runtime_version`. The dev build
`scripts/build-dist.sh` derives its default version from that field — matching the
published archive name `scripts/publish-template.sh` produces
(`context-circuit-v<template_version>`) — carries no hardcoded literal, and fails
loudly on an empty read. Each run is a guarded clean rebuild: it removes and
recreates the output directory (refusing an empty path, `/`, or the source root) so
runs replace rather than accumulate. `scripts/release-artifact.sh` is unchanged —
it still takes a validated positional version and keeps its strict no-overwrite
guard.

Release is a gate: a source change is not a product-template change until release
assembly includes it and the template acceptance checks pass.

The `agent-harness/` assembles or selects the same `context-circuit-template`
artifact that would be distributed and must not import the source repository's
Product Knowledge, plans, `.runtime/`, or implementation state.

## Interfaces

- Release boundary and shipped/never-ship sets: `.context-circuit/wrapper/manifest.yaml`,
  `scripts/release-manifest.txt`
- Assembler: `scripts/release-artifact.sh`, `scripts/build-dist.sh`
- Migration guidance: `.context-circuit/wrapper/migrations/README.md`
- Template seed: `template/`

## Constraints and edge cases

Credentials are never stored or shipped. The distributable artifact is never
published, merged, deployed, or pushed by assembly. Rollback replaces only
template-owned files and preserves workspace-owned state.

## Implementation references

- `.context-circuit/wrapper/manifest.yaml`, `scripts/release-manifest.txt`,
  `scripts/release-artifact.sh`, `scripts/build-dist.sh`
- `.context-circuit/wrapper/migrations/README.md`
- `WORKFLOW.md`, `.context-circuit/docs/release.md`
- root `.gitignore` and `template/.gitignore`

## Provenance

Authored from the current wrapper and release tooling at HEAD `4b8ac0b`; design
chapter 09 was the checklist. Raw `sources/` was not otherwise scanned.

## Acceptance notes

Accepted 2026-08-24 from proposal `0011-domain-source-release-and-upgrade`.
Updated 2026-08-29 from proposal `0028-dist-build-version-identity` (v0.6.1
dist-build-version): the template artifact's version is its `template_version`, and
`build-dist.sh` derives its default from that field and clean-rebuilds its output.
Implementation `10f864b`, `6921a31`.
