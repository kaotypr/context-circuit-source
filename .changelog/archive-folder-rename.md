---
product: cli
kind: changed
upgrade: Every member installs CLI 2.5.0 before pulling a workspace whose `archive/` was renamed to `_archive/`, because a 2.4.x CLI does not read `_archive/`. One person runs the first changing command, commits and pushes the moved records, and everyone else pulls before changing anything.
pr: 209
---

The archive folder is now `_archive/`. The first changing command renames an
existing top-level `intent/archive/` and `plans/archive/` to `_archive/`, keeps
every record in it archived and every ID reserved, and reports the move as a
`notice:` on stderr. The schema stays at 3, so a top-level `archive/` is a
reserved legacy name: it reads as archived and cannot be a group, while an
`archive/` inside a group is an ordinary folder.
