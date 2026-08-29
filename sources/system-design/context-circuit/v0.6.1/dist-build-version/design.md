# dist-build-version — design

## Capability

Make the local dist build (`scripts/build-dist.sh`) name its artifact for the
version the source actually is, by deriving its default version from the single
source of truth (`wrapper/manifest.yaml` `runtime_version`) instead of a
hardcoded literal that drifts.

## Problem

`scripts/build-dist.sh` is the developer/local convenience wrapper around
`scripts/release-artifact.sh`. Run with no version argument, it currently
produces `context-circuit-v0.5.0` — the exact symptom that prompted this scope,
and visible in the committed build output (`dist/context-circuit-v0.5.0/` and its
tarball). The cause is a stale hardcoded default:

```sh
version=${1:-v0.5.0}
```

The source is well past v0.5.0 — `wrapper/manifest.yaml` declares
`runtime_version: 0.6.0` — so the default is wrong: the staged tree, the artifact
directory name (`context-circuit-$version`), and the `.tar.gz` all carry a
version the source is not.

The literal is a **second, drifting copy** of the version number. It is not kept
in step with the source of truth by anything; every version bump would need a
manual edit here, and that edit was missed. This is exactly the failure mode the
publication path already avoids: `scripts/publish-template.sh` derives its
source-versioned assembly from the manifest —

```sh
src_version=v$(sed -n 's/^runtime_version:[[:space:]]*//p' \
  "$source_root/wrapper/manifest.yaml" | head -n1)
```

— and passes *that* to `release-artifact.sh`. `build-dist.sh` is the only version
consumer that hardcodes instead of deriving, so it alone falls out of date.

## Principle

The version is declared once. `wrapper/manifest.yaml` is the owner of the release
boundary and version stamps (`wrapper/contracts/invariants.yaml` owner map:
`release_boundary: wrapper/manifest.yaml`); its `runtime_version` is the source
of truth for the source/runtime version. Any build that needs a *default* version
derives it from that owner — no script restates the version as its own literal.
An explicit argument may still override for ad-hoc builds.

## Fixed decisions

1. **Derive the default from the manifest.** `build-dist.sh`'s default version
   comes from `wrapper/manifest.yaml` `runtime_version`, formatted with the `v`
   prefix the artifact naming already uses (`release-artifact.sh`'s
   `context-circuit-$version`, `publish-template.sh`'s `src_version`). With
   `runtime_version: 0.6.0` the default becomes `v0.6.0`. This is the same read
   `publish-template.sh` already performs; the two paths now agree by
   construction. (The 3-number `runtime_version` also satisfies the
   [design-layout-grouping](../design-layout-grouping/) 3-number-semver
   direction, so the derived tag needs no reshaping.)
2. **An explicit `[version]` argument still wins.** The command interface is
   unchanged — `sh scripts/build-dist.sh [version] [output-dir]`; only the
   *source of the default* changes.
3. **No hardcoded version literal remains** in `build-dist.sh`.
4. **Fail loudly on a missing manifest value.** If the `runtime_version` read
   yields empty (malformed or missing manifest), the build stops with an error
   rather than falling back to a literal or an empty/`context-circuit-v` name.
   This matches `release-artifact.sh`'s existing invalid-version guard, which the
   derived value must anyway pass (no slash, space, leading `.`/`-`); `v` +
   semver satisfies it.
5. **`release-artifact.sh` is unchanged.** It still takes the version as a
   required, validated positional argument; only where `build-dist.sh` *sources
   its default* changes. No manifest edit, no change to the publication scripts.

## Shape of the change (surfaces that move)

Each is owned elsewhere and changed through its owner's normal action:

- **Dist wrapper** (`scripts/build-dist.sh`): replace the
  `version=${1:-v0.5.0}` default with a default read from
  `wrapper/manifest.yaml` `runtime_version` (prefixed `v`), preserving the `$1`
  override; error if the read is empty. The `usage` line and the second
  positional (`output-dir`) are unchanged.

No other source changes. `dist/` is build output (in the manifest `never_ship`
set), not source: the stale `context-circuit-v0.5.0` tree already there is a
regenerable artifact, cleaned by re-running the build after the fix, not tracked
design state this scope needs to migrate.

No core contract bump. No owner is added; the version owner
(`wrapper/manifest.yaml`) and INV semantics are unchanged.

## Constraints and edge cases

- **The derived value must pass `release-artifact.sh` validation**
  (`''|*/*|.*|-*|*' '*|*..*`). A `v`-prefixed 3-number semver does; a blank read
  does not, which is why decision 4 fails rather than passing an empty string
  downstream.
- **Source-checkout only.** `build-dist.sh` and `release-artifact.sh` already
  require a git source checkout (`git rev-parse --show-toplevel`); the manifest
  read is against that same `source_root`, so the derivation has the file it
  needs wherever the build can run at all.
- **Publication path unaffected.** `publish-template.sh` computes `src_version`
  itself and never calls `build-dist.sh`, so its behavior is untouched; this
  scope only brings the *dev* build into line with what publication already does.

## How it feeds the rest — unchanged

This is source. It reaches Product Knowledge and a plan through the normal path:
the coordinator gathers context from this named design source, proposes the
change through the existing context-proposal path, a human accepts, and a plan
grounds in the result. No separate design-acceptance gate.
