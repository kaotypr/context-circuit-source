---
name: cc-cli
description: Install or update the separately released Context Circuit CLI for the current execution OS and architecture, including when a workspace needs its CLI.
---

# Install or update the CLI

Use the execution environment's OS, not the user's desktop OS. WSL, containers,
and remote Linux hosts use Linux packages. Each environment installs independently.

Each workspace pins its own CLI version in `.context-circuit/CLI_VERSION`. Versions
install side by side, so several workspaces on one machine can pin different ones.
Read that file first and treat it as the version to run for that workspace.
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
configuration.

## Run the version a workspace pins

The shared `context-circuit-cli` command on PATH points at whichever version was
installed last, which is not necessarily the one a given workspace pins. For every
workspace operation, resolve the pinned version instead of assuming the shared
command is correct:

1. Read `.context-circuit/CLI_VERSION` in that workspace.
2. Locate the version store. The installer reports it as `version store`; it can
   also be recovered from the shared command, which is a symlink (a `.cmd`
   launcher on Windows) into `<store>/<version>-<os>-<arch>/`.
3. If `<store>/<version>-<os>-<arch>/context-circuit-cli` exists, run that path.
4. Otherwise install that exact version first, then run the reported
   `versioned command`.

Installing a version for one workspace repoints the shared command and must not
change which version another workspace runs. The CLI warns on its error stream
when it is run against a workspace that pins a different version; treat that
warning as an instruction to re-resolve, not as a workspace fault. Verify `version`, then run `--workspace <root> check` only for an
initialized workspace. A blank template still needs `init` when the user requests
workspace initialization. Installation never initializes or migrates a workspace.

On updates, retain workspace records and template version. Report old/new CLI
versions and any compatibility issue. A checksum, platform, download, or version
failure leaves the current command unchanged; fix the cause before retrying.
