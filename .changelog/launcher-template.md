---
product: template
kind: added
upgrade: A 2.5.0 CLI refuses changes to a workspace still at schema 2 unless the pinned CLI runs them. Bring in this template first, then let the pinned CLI make the ID split through the launcher; one person on a team makes the split, commits and pushes it, and everyone else pulls before changing anything.
pr: 177
---

Every workspace now carries a launcher: `.context-circuit/cli`, run as `sh
.context-circuit/cli COMMAND …`, and `.context-circuit/cli.ps1` for Windows
PowerShell or cmd. It runs exactly the version `.context-circuit/CLI_VERSION`
pins from the installers' version store, never a CLI found on PATH, and when
that version is not installed it refuses and names the install step. It refuses
a `--workspace` that names another workspace. It looks in the installers'
default directory, or in `CONTEXT_CIRCUIT_BIN_DIR` when that is set.
`.context-circuit/.gitattributes` keeps `cli` at LF and `cli.ps1` at CRLF line
endings on every checkout.
