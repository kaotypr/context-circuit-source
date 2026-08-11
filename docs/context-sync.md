# Context synchronization

`sync-context` turns selected contribution learnings into small, source-linked
updates to canonical wrapper context. It does not copy conversations, edit live
task state, or treat every implementation detail as durable knowledge.

Create a `context-sync-request` JSON document that lists validated contribution
paths and classifies each proposed learning as `durable-wrapper`,
`repository-local`, `one-off`, or `future-task`. Durable wrapper proposals name
one canonical `context/*.md` target and a concise proposed change. Repository
conventions and future tasks are routed as follow-ups instead of silently being
written as facts.

```bash
node .agents/bin/cc.mjs sync-context --request <request.json>
```

The command refuses a dirty wrapper and creates a dedicated wrapper branch and
worktree. Give a fresh curator only the emitted request, worktree, and allowed
paths. Every canonical update must cite its source contribution path. After the
curator commits, validate the review handoff:

```bash
node .agents/bin/cc.mjs prepare-context-review --sync-id <sync-id>
```

The handoff rejects uncommitted or out-of-scope changes. Team mode requires an
origin-backed wrapper pull request; solo `direct-commit` mode permits presenting
the commit for explicit human confirmation. Neither command pushes, opens or
merges a pull request, deletes evidence, or edits a product repository.
