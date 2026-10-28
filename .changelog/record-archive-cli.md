---
product: cli
kind: added
pr: 209
---

`record archive` is new. Without `--id` it moves the active member's finished
records, a completed plan or an approved intent whose plans are all completed,
under the `_archive/` beside them; with repeatable `--id` it moves exactly the
named records, whoever wrote them, and marks any unfinished one. IDs stay
reserved, lookup by ID keeps working, inline relative links into and out of the
moved records are rewritten, and any other broken link is reported. It returns
`commit_recommended` like approval and completion. `record complete` returns
`archive_offer` while the active member has finished records not archived.
