---
product: cli
kind: changed
upgrade: A workspace on schema 2 is refused changes by a 2.5.0 CLI unless it runs the pinned CLI. Upgrade the template first, then let the pinned CLI make the split; one person on a team makes the split, commits and pushes it, and everyone else pulls before changing anything.
pr: 179
---

When no hand-off happens, because the pin is not installed beside the running
CLI or the hand-off was skipped, a change is refused, with nothing changed, if
it would move the workspace to schema 3, `init` included, or if the running CLI
and the pin are on different release lines. There is no override:
`CC_CLI_VERSION=local` never lifts a refusal. The refusal and the warning name
the workspace's launcher, and either where the pinned CLI is installed or the
installer call that installs it, from the registry `cli_registry` records; a
token is named, never printed. Read-only commands, previews such as `record
delete` without `--confirm`, and other changes on the same release line only
warn, as before.
