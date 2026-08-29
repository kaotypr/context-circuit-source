# dist-build-version — design

## Capability

Make the local dist build (`scripts/build-dist.sh`) name its artifact for the
**template release version** the artifact actually carries, by deriving its
default version from the single source of truth (`wrapper/manifest.yaml`
`template_version`) instead of a hardcoded literal that drifts.

## Problem

`scripts/build-dist.sh` is the developer/local convenience wrapper around
`scripts/release-artifact.sh`. Run with no version argument, it currently
produces `context-circuit-v0.5.0` — the symptom that prompted this scope, and
visible in the committed build output (`dist/context-circuit-v0.5.0/` and its
tarball). The cause is a stale hardcoded default:

```sh
version=${1:-v0.5.0}
```

But the real defect is deeper than a stale literal: it is the **wrong version
altogether**. The shipped artifact is the `context-circuit-template` — a
universal-project-workspace-template (`wrapper/manifest.yaml`
`artifact_kind`). Its release identity is `template_version`, not the
context-circuit product's own `runtime_version`. `wrapper/manifest.yaml` carries
both, and they mean different things:

- `runtime_version: 0.6.0` — context-circuit's **own runtime/source** version.
- `template_version: 0.0.1-alpha.2` — the **published template release**
  version; "Stamped at publish time; the context-circuit-template git tag is
  authoritative."

The dist build builds *the template*, so the artifact's version is the
**template** version. A hardcoded `v0.5.0` is neither — it is a third,
drifting copy that matches nothing.

The publication path already treats the template version as the release
identity: `scripts/publish-template.sh` takes the template `<version>` as its
argument, stamps `template_version: $version` into the shipped manifest, tags the
template repo `v$version`, and names the published archive
`context-circuit-v$version` (its final line copies
`dist/context-circuit-v$version.tar.gz`). The `runtime_version`-derived
`src_version` in that script is only an **internal staging name** for the
pre-stamp assembly, never the release identity. `build-dist.sh` should mirror the
*published* artifact name — `context-circuit-v<template_version>` — not the
staging label and not a literal.

## Principle

The version is declared once. `wrapper/manifest.yaml` owns the release boundary
and version stamps (`wrapper/contracts/invariants.yaml` owner map:
`release_boundary: wrapper/manifest.yaml`). The **template artifact's** identity
is its `template_version`; any build that produces the template and needs a
default version derives it from that field — no script restates the version as
its own literal, and no build labels the template with context-circuit's internal
`runtime_version`. An explicit argument may still override for ad-hoc builds.

## Fixed decisions

1. **Derive the default from `template_version`.** `build-dist.sh`'s default
   version comes from `wrapper/manifest.yaml` `template_version`, formatted with
   the `v` prefix the release naming already uses. With
   `template_version: 0.0.1-alpha.2` the default becomes `v0.0.1-alpha.2`,
   producing `context-circuit-v0.0.1-alpha.2` — the **same name the published
   archive carries** in `publish-template.sh` (`context-circuit-v$version`). Dev
   and publication now agree on what the artifact is called.
2. **Not `runtime_version`.** The context-circuit runtime version is an internal
   detail of the product, not the template's release identity; the dist build
   must not stamp it onto the template artifact. (It remains what
   `publish-template.sh` uses for its intermediate staging name only.)
3. **An explicit `[version]` argument still wins.** The command interface is
   unchanged — `sh scripts/build-dist.sh [version] [output-dir]`; only the
   *source of the default* changes.
4. **No hardcoded version literal remains** in `build-dist.sh`.
5. **Fail loudly on a missing manifest value.** If the `template_version` read
   yields empty (malformed or missing manifest), the build stops with an error
   rather than falling back to a literal or an empty/`context-circuit-v` name.
   The derived value must anyway pass `release-artifact.sh`'s invalid-version
   guard (no slash, space, leading `.`/`-`); a `v`-prefixed semver — including a
   prerelease like `v0.0.1-alpha.2` — satisfies it.
6. **`release-artifact.sh` is unchanged.** It still takes the version as a
   required, validated positional argument; only where `build-dist.sh` *sources
   its default* changes. No manifest edit, no change to the publication scripts.
7. **Clean rebuild — the dist output is replaced, not accumulated.** Every
   `build-dist.sh` run yields a fresh output directory: before building, it
   removes any prior output tree at `output_dir` (the default `dist/`, or `$2`)
   and recreates it. Without this, a second run either fails —
   `release-artifact.sh` refuses to overwrite an existing
   `output_dir/context-circuit-$version` (its "output artifact already exists"
   guard) — or, across versions, silently leaves stale artifacts (the
   `context-circuit-v0.5.0` tree) beside the new one. The clean is scoped to the
   *dev wrapper*; `release-artifact.sh` keeps its strict no-overwrite guard so
   the publication path is never made to clobber. The clean is **guarded**: it
   refuses to remove an empty path, `/`, or the source root, so a mistyped
   `output-dir` cannot wipe an unintended tree.

## Shape of the change (surfaces that move)

Each is owned elsewhere and changed through its owner's normal action:

- **Dist wrapper** (`scripts/build-dist.sh`): replace the
  `version=${1:-v0.5.0}` default with a default read from
  `wrapper/manifest.yaml` `template_version` (prefixed `v`, the same
  `sed -n 's/^template_version:...' | head -n1` shape `publish-template.sh`
  already uses for `runtime_version`), preserving the `$1` override; error if the
  read is empty. The `usage` line and the second positional (`output-dir`) are
  unchanged. It also **clears any prior `output_dir` and recreates it** before
  staging (decision 7), so the run replaces rather than accumulates or fails.

No other source changes. `dist/` is build output (in the manifest `never_ship`
set), not source: the stale `context-circuit-v0.5.0` tree already there is a
regenerable artifact — now removed automatically by the next clean rebuild — not
tracked design state this scope needs to migrate.

No core contract bump. No owner is added; the version owner
(`wrapper/manifest.yaml`) and INV semantics are unchanged.

## Constraints and edge cases

- **Prerelease versions must survive validation.** `release-artifact.sh` rejects
  `''|*/*|.*|-*|*' '*|*..*`. `v0.0.1-alpha.2` passes: it starts with `v` (not
  `-` or `.`), has no slash, space, or `..`, and the internal `-` is not leading.
  A blank read does *not* pass, which is why decision 5 fails rather than passing
  an empty string downstream.
- **The source's `template_version` is the last published one.** A local build
  therefore names the artifact for that release (plus whatever uncommitted source
  it contains) — the honest identity of "the template as of the last release."
  Producing the *next* version's name is `publish-template.sh`'s job (it takes an
  explicit version and stamps it); `build-dist.sh` is a pre-publish inspection
  build, so defaulting to the current `template_version` is correct. A user who
  wants a specific label still passes it as `$1`.
- **Source-checkout only.** `build-dist.sh` and `release-artifact.sh` already
  require a git source checkout; the manifest read is against that same
  `source_root`, so the derivation always has the file it needs.
- **The clean must not become a foot-gun.** Because `output_dir` defaults to
  `dist/` but is overridable as `$2`, the pre-build removal guards against
  catastrophic targets: it refuses an empty value, `/`, and the source root
  itself, failing loudly instead of deleting. A user pointing `$2` at a real
  directory still owns that choice, but the obvious mistakes cannot wipe a tree
  the build never created.
- **Publication path unaffected.** `publish-template.sh` computes its own
  versions and never calls `build-dist.sh`, so its behavior is untouched; this
  scope only brings the *dev* build's default into line with the published
  artifact's identity.

## How it feeds the rest — unchanged

This is source. It reaches Product Knowledge and a plan through the normal path:
the coordinator gathers context from this named design source, proposes the
change through the existing context-proposal path, a human accepts, and a plan
grounds in the result. No separate design-acceptance gate.
