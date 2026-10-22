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
For an explicit update request, resolve the requested version, or the latest
release within this workspace's own release line, from the registry
`cli_registry` names (none means the public
[CLI releases](https://github.com/kaotypr/context-circuit-source/releases)).
The pin says which line:

- A pin `X.Y.Z` is the stable line. Only `cli-vX.Y.Z` releases without a
  suffix are candidates. Exclude drafts and prereleases.
- A pin `X.Y.Z-<label>.<N>` is that label's line, where the label is 1-20
  lowercase letters and not `alpha`, `beta`, `rc`, `dev`, `test`, `pre`,
  `preview` or `snapshot`. Only `cli-vX.Y.Z-<label>.<N>` releases with the
  same label are candidates.

Order candidates by X, Y, Z, then N, as numbers, never as text. Never offer a
release from another line as an update: moving lines is a template change a
person names. Use available GitHub tools or the releases API; on GitLab, list
releases through its API with the same token the installer uses. The workspace
template and the CLI release on the same repository: template releases are
tagged `v*` and CLI releases `cli-v*`. Pass an exact version to the installer;
do not guess that an unpublished version exists. A stable version installs the
command `context-circuit-cli`; a `<label>` version installs
`<label>-context-circuit-cli`. Both share one version store.

Run the bundled script, using absolute paths when outside this skill directory:

- macOS/Linux: `sh scripts/install.sh --version <version>`
- Windows: `powershell -NoProfile -ExecutionPolicy Bypass -File scripts/install.ps1 -Version <version>`

The Windows policy flag applies only to this installer process; it does not change
the user's saved execution policy. Honor any organization policy that blocks it.

## Choose which registry hosts this workspace's CLI

Read `cli_registry` in the workspace's `workspace.yaml` before installing. A
workspace published from somewhere other than the public repository records the
registry it came from:

```yaml
cli_registry:
  source: gitlab
  repository: acme/context-circuit-source
  api: https://gitlab.example.com/api/v4
```

Pass those three values as `--source` / `--repo` / `--api` (PowerShell `-Source`
/ `-Repo` / `-Api`; the environment variables `CONTEXT_CIRCUIT_SOURCE`,
`CONTEXT_CIRCUIT_REPO`, `CONTEXT_CIRCUIT_API` are read as fallbacks).

No `cli_registry` key means the public GitHub repository the installer already
defaults to; pass none of these and it resolves correctly. Never guess a
registry that is not written there — a wrong host is not a failed install, it is
a request for a CLI binary sent somewhere nobody chose.

## Authenticate while the registry is private

A private repository or project serves release assets only through its API, not
the plain download path (which answers 404 on GitHub, and does not exist at all
for a private GitLab project — gitlab installs always require a token). Supply a
credential with read access: `--token <value>` / `-Token <value>`, or the
environment variable `CONTEXT_CIRCUIT_TOKEN` (`GH_TOKEN` and `GITHUB_TOKEN` are
read as github-specific fallbacks). For github, a token makes the installer
resolve the release through the API and download each asset by id; without one
it uses the public download path, which works unchanged if the repository is
ever published. For gitlab, the token is a personal or project access token with
at least read access to the project, sent as `PRIVATE-TOKEN`.

Ask the user for the token or read it from an environment variable they already
set. Never write it into workspace files, a shell profile, a commit, or the
transcript, and do not pass it to anything but these installer scripts. A 404
while a token is set usually means the token lacks access to the repository, not
that the version is missing. Offline `--archive` / `-Archive` installation from a
trusted release needs no token.

Both install and update use the same script. `--bin-dir` / `-BinDir` selects a
user-writable command directory. `CONTEXT_CIRCUIT_BIN_DIR` sets the same
directory for the installers and the workspace launcher; set it when you install
with `--bin-dir`. No administrator access is needed. The installer verifies
SHA-256 and the CLI's version before switching the command to the new version.
Old versions remain available for rollback by reinstalling that version. Offline
use accepts `--archive` and `--checksums` (PowerShell `-Archive` and
`-Checksums`) from the same trusted release.

Check the reported versioned command's `version`, then run workspace commands
through the launcher. If the command directory is outside PATH, explain the
directory to add; do not rewrite shell profiles or unrelated host configuration.

## Run the version a workspace pins

Every workspace carries `.context-circuit/cli` and `.context-circuit/cli.ps1`.
They read the pin and run exactly that version from the version store. Call the
launcher for every workspace command. When it refuses because the pinned version
is missing, install exactly that version with this skill, then call the launcher
again.

A CLI called any other way hands off to the pinned version when that version is
installed in the same version store, with one note on its error stream. When it
is not installed, that CLI warns, and refuses a change that would migrate the
workspace schema or cross release lines. Treat the warning or refusal as the
wrong CLI. `CC_CLI_VERSION=local` runs the called binary for one command without
the hand-off, only when a person asks to try a specific build; it never lifts a
refusal.

Verify `version` through the launcher, then run `check` only for an initialized
workspace. A blank template still needs `init` when the user requests workspace
initialization. Installation never initializes or migrates a workspace.

## Recognize a CLI that is behind the template

The workspace's `AGENTS.md` wins whenever it and CLI output disagree. A running
version on a different release line is not behind; it is the wrong CLI for this
workspace: run the launcher. A CLI is
behind the template when the running version is older than the pin, or when its
output has a shape the template no longer expects — before 2.4.0, `record
approve` returned `planning_required` and told the agent to plan without asking.
Tell the person the CLI needs installing or upgrading to the pinned version,
offer to do it, and do not act on the older output's text meanwhile.

On updates, retain workspace records and template version. Report old/new CLI
versions and any compatibility issue. A checksum, platform, download, or version
failure leaves the current command unchanged; fix the cause before retrying.
