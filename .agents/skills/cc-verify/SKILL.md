---
name: cc-verify
description: Run the independent, read-only verifier against the latest worker commits for an execution.
---

The verifier is a separate actor from the worker. It receives the immutable plan
snapshot, the repository map, the latest worker commit for each affected
repository, the worker handoff (as a claim, not evidence), the canonical
verification commands and evidence identifiers, and read-only worktree access.

It must:

- inspect the latest commit in every affected repository;
- replay the promised acceptance and verification evidence;
- check the evidence proves the required layer;
- check that repository and path scope was respected;
- report `passed`, `failed`, or `blocked` with evidence references;
- write only its own verifier result and handoff.

It must not modify product files, repair the worker's implementation, change
plan status, treat a worker claim as independent evidence, or downgrade an
evidence requirement because a host lacks a capability. Only `passed` satisfies
verification. Record the outcome with the runtime
`verifier-result-record <execution-dir> <attempt> passed|failed|blocked`.
`<attempt>` is `current_attempt` (`3` and `003` are the same attempt; a
different number is refused). The runtime rejects any product write and any
changed branch tip.

If an independent read-only verifier cannot be created, report `blocked`; never
self-verify. Full role behavior: `agents/verifier.md`.
