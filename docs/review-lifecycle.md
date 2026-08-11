# Verification, repair, and review preparation

The planless workflow uses a fresh read-only verifier after every completed
worker. Recording a verifier result confirms that the worktree is clean, still
uses the isolated branch, and has the exact commit recorded by the worker.

## Bounded repair

When verification fails, prepare one scoped repair attempt:

```bash
node .agents/bin/cc.mjs prepare-repair --run-id <run-id> --repository frontend
```

The command reads the current verifier findings, checks the worktree and recorded
commit, increments the manifest attempt counter, and writes new worker and
verifier inputs with unique evidence paths. Give the repair-worker input to a
fresh session. Its result describes the complete base-to-head commit history and
changed-file set, not only the newest commit.

Record the repair worker and the new independent verifier with the existing
recorder:

```bash
node .agents/bin/cc.mjs record-result --run-id <run-id> --repository frontend --stage worker-result
node .agents/bin/cc.mjs record-result --run-id <run-id> --repository frontend --stage verifier-result
```

Do not record `worker-started` for a prepared repair: `prepare-repair` already
advances the run to `running`. If another failure reaches
`workflow.maximum_repair_attempts`, the next preparation call records
`repair-exhausted`, changes the run to `blocked`, and emits no new writable input.
Preserve the branch, worktree, and evidence for human recovery.

## Draft pull-request handoff

After independent verification passes, prepare review metadata:

```bash
node .agents/bin/cc.mjs prepare-review --run-id <run-id> --repository frontend
```

The generated document validates against
`.agents/contracts/review-preparation.schema.json`. It contains the base and head
branches and commits, complete changed-file list, title, body, and worker and
verifier evidence. It records only the Git remote name, so credential-bearing
remote URLs cannot enter runtime evidence.

The handoff is `ready` when an `origin` remote exists and `blocked` with an
explicit explanation when it does not. Preparation is idempotent for an
unchanged verified commit and performs no external mutation.

With explicit authorization, the coordinator may push only the prepared branch
and use an authenticated host-native `gh` or `glab` command to open a draft pull
request from the prepared title and body. Record the confirmed response:

```bash
node .agents/bin/cc.mjs record-review-publication \
  --run-id <run-id> --repository frontend --status published --tool gh \
  --pull-request <reference> --evidence "gh confirmed the draft pull request"
```

The recorder revalidates the clean worktree and exact prepared head, rejects
credential-bearing or conflicting evidence, and is idempotent. It never invokes
the remote tool itself. Missing authorization/tooling returns a manual handoff;
failure must not be represented as a published pull request. Merge and cleanup
remain human-gated.
