---
name: context-next
description: Present the read-only next-action card derived from current workspace evidence.
---

Use when the human asks what is eligible or what should happen next. Delegate
route evaluation to `wrapper/runtime/engine.sh` and report observed state,
selected probe, eligibility, authorization, reason codes, blockers, and one
safe next request. Do not claim ownership or change status. The router and
invariant owners remain canonical.
