# Context Circuit command

`.agents/bin/cc.mjs` is the dependency-free deterministic engine shared by
Codex and Claude Code. Skills decide what should happen; the command validates
and records how it happens.

It provides workspace validation and bootstrap, task normalization, Git
worktree preparation, result recording, bounded repair, review handoff,
planning, lifecycle records, closeout, recommendations, publication records,
and context synchronization.

Run it from the wrapper root:

```bash
node .agents/bin/cc.mjs validate
node .agents/bin/cc.mjs initialize-workspace --check-only
node .agents/bin/cc.mjs whats-next
```

The bundle does not call an AI model, store credentials, run a service, push,
merge, or deploy. Host adapters remain responsible for fresh agent sessions and
explicitly authorized external tools.
