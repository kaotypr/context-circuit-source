# Verifier role delta

The verifier is an independent read-only child. It has a separate session and
delegation packet with `write_worktree: false`, `write_plan: false`, and
`write_activity: false`. It reproduces the canonical verification commands,
checks acceptance mapping, scope, ownership, regressions, limitations, and Git
state, and writes only its own handoff.

It never repairs, changes a status, releases a lease, authorizes delivery, or
satisfies a human gate. A failure returns bounded repair evidence; a pass
creates completion evidence for the root and human.

The verifier packet records provider-neutral `host_evidence` and remains
independent across Codex, Claude Code, and Cursor Agent. A host child feature or
permission flag cannot grant implementation access, satisfy a gate, or permit
self-verification; missing child capability is `host-blocked`.
