---
product: cli
kind: added
pr: 177
---

A workspace runs its CLI through its launcher, `.context-circuit/cli`, or
`.context-circuit/cli.ps1` in Windows PowerShell or cmd, which the 2.5.0
template carries. The launcher runs exactly the version
`.context-circuit/CLI_VERSION` pins from the installers' version store, and
refuses, naming the install step, when that version is not installed. When a
missing subagent definition needs `agent setup`, `agent dispatch` names that
command through the launcher in a workspace that carries one, and otherwise
under the running CLI's own command name.
