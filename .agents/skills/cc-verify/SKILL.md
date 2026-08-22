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

The verifier records host identity and observed capability in
`host_evidence`, but host capability never replaces independent read-only
permissions. Codex, Claude Code, and Cursor Agent child mechanisms are
interchangeable adapters around the same verifier packet. If the required child
cannot be created, record `host-blocked` and do not verify in the root session.
