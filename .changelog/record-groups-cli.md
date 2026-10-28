---
product: cli
kind: added
pr: 209
---

Records may sit in one level of group folder under `intent/` or `plans/`, such
as `plans/alex/<plan folder>/`, with `_archive/` above or below that group, and
every command that reads records finds them there. A group is ordinary file
organization and says nothing about who owns a record. Folder names starting
with `_` are kept for the CLI, and `_archive/` is the only one it gives meaning
to. `check` reports a `reserved folder name`, a `folder named like a record
holds records`, a `record file outside the folders records are read from`, and a
`legacy archive folder awaits its rename to _archive/`.
