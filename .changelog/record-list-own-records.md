---
product: cli
kind: changed
upgrade: Scripts that relied on `record list` listing every record pass `--all`.
pr: 209
---

`record list` now lists the active member's own records. `--member ID` lists
another member's, and `--all` every record, including records with no author.
With no member selected it lists everything and says so on stderr. `status`
opens with `my_work`: the active member's intents awaiting approval, approved
intents without plans, plans not completed, and `unarchived_finished`, the count
of their finished records not archived.
