# Context synchronization

Context Circuit has three related workflows with distinct boundaries:

- `cc-gather-context` resolves the minimum context for one task and remains
  entirely read-only. Its findings are temporary unless another reviewed
  workflow records them.
- `cc-import-context` performs bounded, deterministic discovery against one
  named registered repository. It writes an append-only, source-cited
  contribution in the wrapper but never modifies the source repository.
- `cc-sync-context` classifies tracked contributions and prepares scoped changes
  to canonical wrapper context in an isolated worktree.

Use import when repository evidence should bootstrap or refresh durable context;
use gather when the evidence is needed only to understand the current task.

## Import registered repository evidence

Create an ignored request that matches
`.agents/contracts/import-context-request.schema.json`:

```json
{
  "contract_version": 1,
  "repository": "product",
  "authorize_contribution_write": true,
  "limits": {
    "max_files": 50,
    "max_file_bytes": 131072,
    "max_total_bytes": 1048576
  }
}
```

The authorization covers only creation of a new contribution snapshot. Run the
canonical skill, or invoke its deterministic discovery step from the clean
wrapper root:

```bash
node .agents/bin/cc.mjs import-context --request <import-context-request.json>
```

Discovery checks that the repository is registered at its Git root and that both
the wrapper and source repository are clean. It reads root instructions, common
documentation folders, nested repository instruction files, and structural
signals in a fixed order. A repository does not need a `context/` directory;
when `context/*.md` exists it is recorded last as high-trust evidence. Limits
bound the number of files and excerpt bytes. Each evidence item records the exact
repository-relative path and source commit. Unknowns must remain explicit, and
high trust does not permit unsupported inference or blind copying.

The command writes:

- `contributions/import-context/<repository>/<timestamp>-import-context-<repository>.md`,
  an append-only snapshot for review and curation; and
- `.runtime/import-context/<timestamp>-<repository>/manifest.json`, ignored
  evidence containing the bounds, source commit, and discovered items.

It also returns structured handoffs for `sync-context --request` and
`prepare-context-review --sync-id`. The source repository remains unchanged.
Work from a dedicated non-default wrapper review branch. Review the snapshot for
source accuracy and secrets, then commit only that file locally after exact
authorization. Synchronization starts only after the contribution is tracked at
`HEAD` and the wrapper is clean; it never hides or discards unrelated changes.
The synchronization worktree starts from this commit, so the eventual review
against the configured wrapper target includes both the append-only evidence and
the canonical-context changes. A local contribution commit does not authorize a
push, pull request, or merge.

## Curate tracked contributions

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

For imported evidence, merge or append the smallest supported durable facts into
the existing canonical files. Every changed target must cite the import
contribution path. Overlap, uncertainty, repository-local conventions, one-offs,
and future work remain visible in the review handoff instead of being silently
overwritten or promoted to current facts.
