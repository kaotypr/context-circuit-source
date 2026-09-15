# Two products and their version lines

One checkout assembles two separately released products.

| Product | Version file | Tags | Contents |
| --- | --- | --- | --- |
| Workspace template | `context-circuit-source@VERSION` | `v*` on the published template repository | shared instruction, docs, skills, blank seed |
| Executable | `context-circuit-source@CLI_VERSION` | `context-circuit-cli-v*` here, assets published on the template repository | Go binaries for six platform targets |

They change for different reasons and at different rates. A wording fix in the
shared instruction should not force anyone to reinstall a binary, and a
filesystem-cloning fix should not force a workspace to adopt new instructions.
Coupling them would make every release the larger of the two. Changing either
version file does not publish; publication is explicitly invoked.

## Per-workspace pinning

A workspace records the executable version it expects, and that record is shared
— it travels with the workspace through Git. Versions install side by side in a
version store, one directory per version and platform, and the shared command on
PATH is a link into one of them. Installing a version for one workspace repoints
that shared command and must not change which version another workspace runs.

Resolution is therefore: read the workspace's pinned version, locate the store,
run that exact path, and install it first if it is absent. When the executable
notices it is running against a workspace pinning a different version it warns
on the error stream and continues, leaving structured output untouched. The
warning is framed as a caller mistake — the workspace recorded what it wanted
and something ran the wrong binary.

Owner: `context-circuit-source@internal/cli/cli.go` `pinnedVersionWarning`, and
the installer under `context-circuit-source@product/skills/cc-cli/`.

The executable is tagged and built in this checkout, where its source is, so a
released binary stays traceable to the commit that produced it. Its assets are
published on the template repository and its GitLab mirror instead, so installing
the CLI never needs read access to this maintainer checkout. The tag created on
those destinations is a distribution marker rather than a build reference, and the
release notes carry the source commit.

## Installation properties

Installation targets the **execution environment**, not the user's desktop: a
container or remote Linux host installs Linux packages. An organization that
mirrors releases into its own GitLab project installs from that project's generic
package registry, which the workspace records as `cli_registry` so one member
configures it for everyone; the credential reaches only the registry in use. It needs no
administrator access, verifies a checksum and the executable's own reported
version before switching the shared command, leaves the current command
untouched on any failure, and supports offline install from a trusted release.
Rollback is reinstalling the older version, which is still in the store.
Installation never initializes, migrates, or rewrites a workspace.

## The embedded seed

The executable embeds a blank workspace used by initialization and by export.
It is a convenience, not a second source of truth: the exact source-to-output
mapping is `context-circuit-source@scripts/release-manifest.txt`, the embedding
reads only what that manifest names, and the release check verifies the embedded
inventory against it so the seed cannot drift from the shipped template by
accident. An initialized workspace is never automatically rewritten.
