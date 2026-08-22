# Maintainer release

This source-only document describes staged artifact assembly. It is excluded
from the released workspace.

Run from the intended source checkout:

```text
sh test/acceptance.sh
sh scripts/release-artifact.sh /tmp/cc-release-stage /tmp/cc-release-out v1.0.0
```

For a local build, the convenience command writes the versioned generated
directory and archive to the ignored `dist/` folder:

```text
sh scripts/build-dist.sh
find dist/context-circuit-v0.5.0 -type f | sort
```

Pass a version explicitly when building another release, for example
`sh scripts/build-dist.sh v0.6.0`.

Pass a second output directory when testing without using the repository's
`dist/` folder.

The assembler stages the current source tree without `.git`, runtime, tests,
scripts, source inbox, plans, Product Knowledge, maintainer design material, old skill
adapters, the source-only `template/` directory, or credentials. It overlays
root adapters from `wrapper/adapters/` and the uninitialized mutable seed from
`template/` at the artifact root. It reports source revision, dirty state,
manifest inventory, and the remaining publication gate.

The artifact is not published, merged, deployed, or pushed by this command.
Rollback replaces wrapper-owned files only and preserves workspace identity,
accepted context, sources, plans, runtime, dirty worktrees, and repositories.
