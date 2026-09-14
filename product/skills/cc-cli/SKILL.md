---
name: cc-cli
description: Install or update the separately released Context Circuit CLI for the current execution OS and architecture, including when a workspace needs its CLI.
---

# Install or update the CLI

Use the execution environment's OS, not the user's desktop OS. WSL, containers,
and remote Linux hosts use Linux packages. Each environment installs independently.

For initial setup, read `.context-circuit/CLI_VERSION` for the recommended version.
For an explicit update request, resolve the requested version, or the latest stable
compatible `cli-v2.*` release from
[CLI releases](https://github.com/kaotypr/context-circuit-source/releases).
Use available GitHub tools or the releases API; exclude drafts/prereleases unless
requested. Do not confuse workspace `v*` releases with CLI `cli-v*` releases. Pass
an exact version to the installer; do not guess an unpublished version exists.

Run the bundled script, using absolute paths when outside this skill directory:

- macOS/Linux: `sh scripts/install.sh --version <version>`
- Windows: `powershell -NoProfile -ExecutionPolicy Bypass -File scripts/install.ps1 -Version <version>`

The Windows policy flag applies only to this installer process; it does not change
the user's saved execution policy. Honor any organization policy that blocks it.

Both install and update use the same script. `--bin-dir` / `-BinDir` selects a
user-writable command directory. No administrator access is needed. The installer
verifies SHA-256 and the executable's version before switching the command to the
new version. Old versions remain available for rollback by reinstalling that
version. Offline use accepts `--archive` and `--checksums` (PowerShell `-Archive`
and `-Checksums`) from the same trusted release.

Use the reported command path immediately. If its directory is outside PATH,
explain the directory to add; do not rewrite shell profiles or unrelated host
configuration. Verify `version`, then run `--workspace <root> check` only for an
initialized workspace. A blank template still needs `init` when the user requests
workspace initialization. Installation never initializes or migrates a workspace.

On updates, retain workspace records and template version. Report old/new CLI
versions and any compatibility issue. A checksum, platform, download, or version
failure leaves the current command unchanged; fix the cause before retrying.
