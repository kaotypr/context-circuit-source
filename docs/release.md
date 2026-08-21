# Maintainer release

This source-only document describes staged artifact assembly. It is excluded
from the released workspace.

Run from the intended source checkout:

```text
sh test/acceptance.sh
sh scripts/release-artifact.sh /tmp/cc-release-stage /tmp/cc-release-out v1.0.0
```

The assembler stages the current source tree without `.git`, runtime, tests,
scripts, source inbox, plans, Product Knowledge, maintainer plan, old skill
adapters, or credentials. It overlays root adapters from `wrapper/adapters/`
and the uninitialized mutable seed from `template/`. It reports source revision,
dirty state, manifest inventory, and the remaining publication gate.

The artifact is not published, merged, deployed, or pushed by this command.
Rollback replaces wrapper-owned files only and preserves workspace identity,
accepted context, sources, plans, runtime, dirty worktrees, and repositories.
