---
product: cli
kind: changed
authors: "@hallode"
pr: 170
---

The CLI's code is restructured with no change in behaviour: each command is
declared once, and `check`, `record order`, `record create`, and `record
dependencies` read the records once per command rather than once per lookup.
