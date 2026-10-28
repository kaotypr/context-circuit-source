---
product: cli
kind: added
pr: 176
---

A stable version installs as `context-circuit-cli`, as before. A variant
version, `X.Y.Z-<label>.<N>`, installs as `<label>-context-circuit-cli`, so it
never shares a command with the public one; both keep their versions in the same
store. Both installers take `CONTEXT_CIRCUIT_BIN_DIR` as their default command
directory, and the launcher reads it too, so set it when installing with
`--bin-dir`.
