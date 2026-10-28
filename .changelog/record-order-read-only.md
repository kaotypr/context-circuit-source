---
product: cli
kind: fixed
pr: 171
---

`record order` is read-only: it writes nothing and takes no lock. `repo fetch`
changes only Git refs in a bound checkout, so it no longer makes the ID split on
a workspace still at schema 2.
