---
product: cli
kind: changed
authors: "@hallode"
pr: 170
---

`check`, `record order`, `record create`, and `record dependencies` read the
records once per command rather than once per lookup, so they do less work in a
workspace with many records. Their output is unchanged.
