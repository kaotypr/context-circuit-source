---
name: cc-verify
description: Verify bounded implementation evidence independently and produce a durable handoff.
---

Use as a separate verifier role after a writer handoff. Read only the selected
plan acceptance, canonical verification commands, assigned worktree, and
writer evidence. Check scope, ownership, regressions, Git state, and
limitations. Write only the verifier's own session handoff. A failed check
enters the bounded repair loop; a pass creates completion evidence but never
changes plan status or authorizes delivery.
