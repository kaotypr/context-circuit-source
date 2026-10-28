---
product: cli
kind: added
pr: 179
---

A CLI that is not the pinned version now hands off to it. When the pin is
installed in the same version store as the running binary, the CLI prints one
note on stderr and runs the pinned binary with the same arguments, setting
`CC_CLI_SWITCHED_FROM` to the version that handed off, so the pinned CLI never
hands off again. `CC_CLI_VERSION=local` keeps the running CLI for one command;
any other value is refused and nothing is run. A source build, or a binary
outside the version store, never hands off.
