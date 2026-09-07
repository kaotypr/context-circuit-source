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
find dist/context-circuit-v1.0.0 -type f | sort
```

Pass a version explicitly when building another release, for example
`sh scripts/build-dist.sh v1.0.0`.

Pass a second output directory when testing without using the repository's
`dist/` folder.

The assembler stages the current source tree without `.git`, runtime, tests,
scripts, source inbox, plans, Product Knowledge, maintainer design material, old skill
adapters, the source-only `template/` directory, or credentials. It overlays
root adapters from `.context-circuit/wrapper/adapters/` and the uninitialized mutable seed from
`template/` at the artifact root. It reports source revision, dirty state,
manifest inventory, and the remaining publication gate.

The artifact is not published, merged, deployed, or pushed by this command.
Rollback replaces wrapper-owned files only and preserves workspace identity,
accepted context, sources, plans, runtime, dirty worktrees, and repositories.

The artifact root includes AGENTS.md, CLAUDE.md, CURSOR.md, WORKFLOW.md, and
README.md from .context-circuit/wrapper/adapters/. CLAUDE.md and CURSOR.md are thin host import
bridges; they contain no host authentication, MCP configuration, memory, or
transcript state. Cursor-local permission files and all provider state remain
outside the release inventory. Upgrade and rollback preserve local host
configuration, user workspace data, registered repositories, and runtime evidence.
