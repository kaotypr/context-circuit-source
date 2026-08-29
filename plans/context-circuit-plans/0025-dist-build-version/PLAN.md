# 0025 — Version and clean the local dist build

- **Plan ID:** `0025-dist-build-version`
- **Status:** done
- **Repository:** `context-circuit-source`
- **Depends on:** `0010-release-assembly-and-publication`
- **Owns:** the dev dist wrapper's default version and clean-rebuild behavior

## Original request

Retroactive plan for the v0.6.1 `dist-build-version` scope, authored as if from
an empty repo. Source design:
`sources/system-design/context-circuit/v0.6.1/dist-build-version/design.md`.

## Objective and desired behavior

- `scripts/build-dist.sh` (the dev wrapper around `release-artifact.sh`) must
  name its artifact for the **template release version** it actually carries. The
  shipped artifact is the `context-circuit-template`, so its identity is
  `wrapper/manifest.yaml` `template_version` — not context-circuit's internal
  `runtime_version`, and not a hardcoded literal (it had drifted to `v0.5.0`).
  The default becomes `context-circuit-v<template_version>`, matching the
  published archive name in `publish-template.sh`.
- Each run **replaces** the dist output: remove and recreate `output_dir` before
  staging, so a re-run neither fails on `release-artifact.sh`'s no-overwrite
  guard nor leaves stale artifacts (the `context-circuit-v0.5.0` tree) beside the
  new one.

## Constraints and non-goals

- An explicit `[version]` argument still overrides; the command interface is
  unchanged.
- Fail loudly if `template_version` is empty (never fall back to a literal or a
  blank name); the derived value must pass `release-artifact.sh` validation.
- `release-artifact.sh` keeps its strict no-overwrite guard — only the dev
  wrapper cleans, so the publication path is never made to clobber.
- The clean is guarded against removing an empty path, `/`, or the source root.
- **Non-goal:** the `context-circuit-` artifact name prefix (separate question).

## Product Knowledge grounding

- `invariants` (`wrapper/contracts/invariants.yaml`) — the owner map names
  `wrapper/manifest.yaml` as the release-boundary/version-stamp owner.

## Tasks

1. **DIST-001** — default the version from `template_version`; clean the output
   dir with guards.

## Acceptance & verification

- Default names the artifact `context-circuit-v<template_version>`; explicit
  version overrides; empty `template_version` fails loudly; re-run replaces
  rather than fails; unsafe output dirs are refused.
- `sh scripts/build-dist.sh` (produces `context-circuit-v<template_version>`).

## Assumptions, open questions, risks

- The source's `template_version` is the last published one, so a local build
  names "the template as of the last release"; producing the *next* version's
  name stays `publish-template.sh`'s job.

## Expected commits and delivery notes

A single wrapper-script change (default derivation + guarded clean); no change to
`release-artifact.sh`, the manifest, or the publication scripts.
