---
id: DIST-001
title: Default the version from template_version and clean the output dir
repositories: [context-circuit-source]
paths: [scripts/build-dist.sh]
depends_on: []
acceptance: [DIST-AC-01, DIST-AC-02]
verification: [DIST-VT-01]
---

## Intended behavior

Run with no version argument, `scripts/build-dist.sh` names the local dist
artifact for the template release version it carries —
`context-circuit-v<template_version>` — the same name the published archive uses.
Each run yields a fresh dist tree: the prior output is replaced, not accumulated,
and a re-run does not fail on the assembler's no-overwrite guard.

## Concrete change and affected surfaces

In `scripts/build-dist.sh`: replace the hardcoded `version=${1:-v0.5.0}` default
with a `v`-prefixed read of `wrapper/manifest.yaml` `template_version` (the same
`sed -n 's/^template_version:...' | head -n1` shape `publish-template.sh` uses),
preserving the explicit `$1` override and the second positional (`output-dir`);
fail loudly if the read is empty. Before staging, remove and recreate
`output_dir` so the run replaces prior output, guarded to refuse an empty path,
`/`, or the source root. `release-artifact.sh`, the manifest, and the publication
scripts are unchanged.

## Inputs and outputs

Inputs: `wrapper/manifest.yaml` `template_version` (the release identity of the
template artifact). Outputs: a freshly built, correctly named artifact
(`context-circuit-v<template_version>`) with any prior dist output replaced. The
context-circuit `runtime_version` is not used to name the template artifact.

## Context references

- sources/system-design/context-circuit/v0.6.1/dist-build-version/design.md

## Risks and open questions

- `output_dir` is overridable, so the pre-build clean must guard against
  catastrophic targets rather than deleting them; the guard fails loudly.
- The derived version must satisfy `release-artifact.sh` validation; a
  `v`-prefixed semver, including a prerelease, does.

## Stop conditions

- Do not change `release-artifact.sh`, the manifest, or the publication scripts.
- Do not touch the `context-circuit-` name prefix (separate question).
