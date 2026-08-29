---
id: 0028-dist-build-version-identity
target_context_unit: source-release-and-upgrade
operation: add
statement: >
  The local dist build (scripts/build-dist.sh) names its artifact for the
  template release identity, not context-circuit's internal runtime version. Run
  with no version argument it derives the default from wrapper/manifest.yaml
  template_version (v-prefixed), matching the published archive name
  publish-template.sh produces (context-circuit-v<template_version>); it no longer
  carries a hardcoded literal and fails loudly on an empty manifest read. Each run
  is a guarded clean rebuild — the output dir is removed and recreated (refusing an
  empty path, /, or the source root) so runs replace rather than accumulate or trip
  release-artifact.sh's no-overwrite guard. release-artifact.sh and the publication
  path are unchanged. This makes explicit, in the domain that already keeps the
  three identities distinct, that the template artifact's version is its
  template_version — never runtime_version.
evidence_refs:
  - sources/system-design/context-circuit/v0.6.1/dist-build-version/design.md
  - scripts/build-dist.sh
  - context/domains/source-release-and-upgrade/README.md:48  # "Three identities stay distinct"
  - context/domains/source-release-and-upgrade/README.md:80  # Assembler: build-dist.sh
affected_repositories:
  - context-circuit
affected_commits:
  - 10f864b  # fix(build): default build-dist.sh version from template_version
  - 6921a31  # feat(build): clean the dist output dir on each build-dist run
related_plan: 0025-dist-build-version
confidence: low
status: review-needed
---

# Record the dist build's version identity

## Context

The `source-release-and-upgrade` domain already keeps three identities distinct —
`context-circuit-source` (this checkout), `context-circuit-template` (the shipped
artifact), and an instantiated workspace — and names `scripts/build-dist.sh` as an
assembler. It does **not** currently state which version the dist build stamps.
That silence let a stale default drift.

## What changed (v0.6.1 · dist-build-version, plan 0025)

`build-dist.sh` previously defaulted to a hardcoded `v0.5.0`, producing
`context-circuit-v0.5.0` — a third version matching neither identity. The shipped
artifact **is** the template, so its release identity is `template_version`
(`0.0.1-alpha.2` at design time), not the product's own `runtime_version`
(`0.6.0`). The fix:

- Derive the default from `wrapper/manifest.yaml` `template_version`, `v`-prefixed
  — the same name `publish-template.sh` gives the published archive
  (`context-circuit-v$version`). Dev and publication now agree on the artifact
  name.
- Never stamp `runtime_version` onto the template artifact; an explicit
  `[version]` argument still overrides; no hardcoded literal remains.
- Fail loudly if the `template_version` read is empty.
- Clean rebuild: remove and recreate the output dir each run, guarded against an
  empty path, `/`, and the source root, so runs replace rather than accumulate
  stale trees or trip `release-artifact.sh`'s no-overwrite guard.

`release-artifact.sh` and `publish-template.sh` are unchanged; only where the dev
wrapper *sources its default* moved.

## Proposed change

Add one clarifying statement to the domain (near the identities/assembler prose):
the template artifact's version is its `template_version`; `build-dist.sh` derives
its default from that field and clean-rebuilds its output; `release-artifact.sh`
still takes a validated positional version and keeps its strict no-overwrite
guard.

## Why low confidence

This is close to implementation-detail altitude — a dev-convenience script's
default and cleanup, with no contract, gate, or owner change (no core contract
bump). It is proposed because it makes the *identity rule* the domain already
implies (template artifact ↔ `template_version`) explicit at the point where the
drift occurred. A reviewer may reasonably judge it below the domain's altitude and
`defer` it; the fix itself already stands in the script regardless.
