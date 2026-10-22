# Context Circuit CLI

A separately versioned Go CLI for Context Circuit workspaces. It handles
workspace records, Git repositories, CoW worktree preparation, and subagent role
configuration. Coding agents provide planning, implementation, and
orchestration.

Download a platform package from this repository's `cli-v<version>` release, or
ask the workspace's `cc-cli` skill to install or update it. The skill detects the
execution environment, verifies the package checksum, and installs in a
user-writable directory. Reading a release that is not public needs a GitHub token
with access to it, which the installers accept as `--token` / `-Token` or in
`CONTEXT_CIRCUIT_TOKEN`. No Go or Python runtime is required. Git is required
for repository operations. Installation does not update workspace files.

Run `context-circuit-cli help` for commands. `context-circuit-cli version` reports the CLI
version; `.context-circuit/VERSION` in a workspace records its separate template
version. CLI v2 supports workspace schema 2. Each workspace pins its CLI version
in `.context-circuit/CLI_VERSION`; a compatible CLI update does not change it.
Versions install side by side under one version store, so workspaces pinning
different versions coexist on a machine, and a variant build installs as
`<label>-context-circuit-cli`. A workspace runs its pinned version through its
launcher, `sh .context-circuit/cli` (or `.context-circuit/cli.ps1` in
PowerShell), which refuses and names the install step when that version is not
installed. A CLI called another way hands off to the pinned version when it is
installed in the same store, and otherwise warns, and refuses a change from
another release line or one that would migrate the workspace schema.

The CLI also embeds a workspace seed for `init` and `template export`. An existing
blank schema-2 template can be initialized independently of its template version.
An initialized workspace is never automatically migrated or rewritten.

Release packages cover macOS, Linux, and Windows on amd64 and arm64. Git, Node,
package managers, and application services are separate environment dependencies.
CoW attempts native filesystem cloning and falls back to independent file copies
in auto mode. Reused dependencies must match the target branch and local runtime.

## License

`context-circuit-cli` is licensed under the Apache License 2.0, and every release
archive carries that license beside the binary along with the third-party notices
for its dependencies and the Go runtime.

The workspace template is a separate product under a separate license — 0BSD,
which asks nothing of the repository it is copied into. It is released from its
own repository, so neither product's terms have to be inferred from the other's.
