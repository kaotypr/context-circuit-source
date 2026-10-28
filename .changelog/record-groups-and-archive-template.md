---
product: template
kind: added
upgrade: Every member installs CLI 2.5.0 before pulling a workspace whose `archive/` was renamed to `_archive/`, because a 2.4.x CLI does not read `_archive/`. One person runs the first changing command, commits and pushes the moved records, and everyone else pulls before changing anything.
pr: 209
---

`AGENTS.md`, the record and command references, and the `cc-complete` skill
describe grouped records and the archive: a record may sit in one level of group
folder with `_archive/` above or below it, `record list` and `status` show your
own work first, and `record archive` moves your own finished records under
`_archive/`. After completion the agent offers to archive, and a no changes
nothing. This workspace version pins CLI 2.5.0, which renames an existing
`intent/archive/` or `plans/archive/` to `_archive/` on its first changing
command; no record's content changes.
