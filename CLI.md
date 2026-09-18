# Context Circuit CLI

A separately versioned Go executable for Context Circuit workspaces. It handles
workspace records, Git repositories, CoW worktree preparation, and subagent role
configuration. Coding agents provide planning, implementation, and orchestration.

Download a platform package from the product repository's
`cli-v<version>` release, or ask the workspace's `cc-cli` skill to
install or update it. The skill detects the execution environment, verifies the
package checksum, and installs in a user-writable directory. While the product
repository is private, a download needs a GitHub token with read access to it; the
installers accept one as `--token` / `-Token` or in `CONTEXT_CIRCUIT_TOKEN`. An
organization that mirrors releases into its own GitLab project installs from there
instead, naming it with `--gitlab-url` / `-GitlabUrl` or in
`CONTEXT_CIRCUIT_GITLAB_URL`. No Go or Python runtime is required. Git is required
for repository operations. Installation does not update workspace files.

Run `context-circuit-cli help` for commands. `context-circuit-cli version` reports the CLI
version; `.context-circuit/VERSION` in a workspace records its separate template
version. CLI v2 supports workspace schema 2. Each workspace pins its CLI version
in `.context-circuit/CLI_VERSION`; a compatible CLI update does not change it.
Versions install side by side under one version store, so workspaces pinning
different versions coexist on a machine. The CLI warns when it is run against a
workspace that pins another version.

The CLI also embeds a workspace seed for `init` and `template export`. An existing
blank schema-2 template can be initialized independently of its template version.
An initialized workspace is never automatically migrated or rewritten.

Release packages cover macOS, Linux, and Windows on amd64 and arm64. Git, Node,
package managers, and application services are separate environment dependencies.
CoW attempts native filesystem cloning and falls back to independent file copies
in auto mode. Reused dependencies must match the target branch and local runtime.
