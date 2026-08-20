# Maintainer release packaging

This document is source-only. It is excluded from the published template
artifact and is not part of the user-facing workspace contract.

The published artifact is a clean cloneable workspace: agent instructions,
skills, context navigation, plans/sources layout, and user-facing docs.
Maintainer validation, GitHub automation, and this file stay in the source
repository.

## Destination and authorization

- Destination repository: `kaotypr/context-circuit-release`
- Destination branch: `main`
- Authorization secret: `secrets.CONTEXT_CIRCUIT_RELEASE_TOKEN`
- Human gate: `publication`

Publication stays behind that secret and the workspace publication gate. A
successful source validation or staged artifact is not permission to publish.
If the token is missing, the destination is unavailable, or the push is
denied, keep the staged artifact and deliver it manually.

## Version and tag contract

- Tag trigger: push a `v*` tag (for example `v0.4.0`).
- Manual trigger: GitHub `workflow_dispatch` with a `version` input.
- Local invocation: pass the same version string as the helper's third
  argument.

The helper names the artifact directory and optional archive with that
version: `context-circuit-<version>/` and `context-circuit-<version>.tar.gz`.

## Reproduce a release

Source validation always runs before staging.

### From a version tag

1. Confirm the source revision is the intended release.
2. Create and push an annotated or lightweight `v*` tag. Tag push is a human
   action.
3. The source workflow `.github/workflows/sync-context-circuit-release.yml`
   checks out that tag, runs `git diff --check` and `sh test/acceptance.sh`,
   stages the artifact, and smoke-checks the inventory.
4. If `CONTEXT_CIRCUIT_RELEASE_TOKEN` is present, the workflow checks out
   `kaotypr/context-circuit-release` at `main` and synchronizes the artifact.
5. If the token is absent or the destination push is denied, download or copy
   the retained artifact and publish it manually.

### Manual workflow invocation

Run **Synchronize Context Circuit release** with a version input. The same
validation, staging, inventory, and publication-or-fallback path applies.

### Local maintainer invocation

From a clean source checkout at the intended revision:

```
git diff --check
sh test/acceptance.sh
sh scripts/release-artifact.sh /tmp/cc-release-stage /tmp/cc-release-out v0.4.0
```

The helper stages tracked files only (`git archive` of `HEAD`). Dirty or
ignored working-tree files, including `.runtime/`, cannot enter the artifact.
It then removes the exclusion list, fails if a required file is missing or a
forbidden path remains, and prints version, source revision, inventory, and
destination.

Copy `/tmp/cc-release-out/context-circuit-v0.4.0/` onto
`kaotypr/context-circuit-release` `main` only after the human publication
gate.

## Artifact boundary

Required in the artifact (minimum): `AGENTS.md`, `WORKFLOW.md`, `CLAUDE.md`,
`workspace.yaml`, `README.md`, `.gitignore`, the starter `context/` files,
user-facing `docs/` contracts, `.agents/skills/`, `agents/`, `plans/README.md`,
and `sources/README.md`.

Excluded from the artifact: `test/`, `.github/`, `scripts/`, this file
(`docs/release.md`), the retired authored-artifact directory, `.runtime/`,
and any Node or package-manager files.

The helper fails instead of stripping package-manager files. The source
acceptance suite also forbids those files in the source tree.

The workflow does not use Node or a package manager.

## Document-system package checklist

The approved package-boundary decision
`dsr-0003-package-boundary-20260820t042053z` distinguishes two layers. The
maintainer source package is the complete evidence package and includes:

- `schemas/documents/` — the versioned OKF, profile, plan/task,
  instruction-skill, runtime, extraction, and validation-result schemas;
- `docs/templates/` — the current Product Knowledge templates;
- `.agents/skills/` and `agents/` — the shared host-neutral skills and local
  safety/ownership roles;
- `test/fixtures/document-system/` — attributed conformance and measurement
  fixtures;
- `context/index.md`, `context/domains/index.md`, and `context/roles/index.md`
  — the lowercase live routing indexes;
- `LICENSES/` — applicable third-party notices and license text.

The published template artifact remains a separate layer governed by the
existing release manifest. Its boundary excludes `test/`, `docs/release.md`,
`.github/`, and `scripts/`; it therefore does not contain the source-package
fixtures. Source-package completeness and published-boundary exclusion are
separate checks. This decision does not authorize publication or change the
release boundary.

## OKF attribution

The Context Circuit profile is derived from the pinned Open Knowledge Format
v0.2 evidence. Its upstream repository, revision, selected files, checksums,
and reasons are recorded in the accepted source provenance at
`sources/okf-v0.2/source.yaml`. The applicable Apache-2.0 terms and attribution
notice are retained in `LICENSES/OKF-v0.2-NOTICE.md` and
`LICENSES/Apache-2.0-OKF.txt`. Derived Context Circuit rules remain identified
as derived documentation; the attribution does not grant a new product or
publication gate.

## Remaining human gates

- **publication** — push or sync to `kaotypr/context-circuit-release` `main`,
  or an equivalent manual delivery. Requires the release token or an explicit
  maintainer copy.
- **merge / tag** — creating the source tag or merging the source change that
  the tag points at is a separate human action.
- **status-change** — this packaging work does not mark a plan done.

A missing token is not a packaging failure. The staged artifact remains the
manual fallback.
