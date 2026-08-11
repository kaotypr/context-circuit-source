# Run a task

`run-task` turns an explicit request into scoped, isolated work without requiring
a plan or activity integration. It prepares artifacts and Git worktrees; it does
not launch agents, push branches, open pull requests, merge, or deploy.

## Prepare one repository

```bash
node .agents/bin/cc.mjs run-task \
  --request "Add a reset button to the counter" \
  --repository frontend \
  --scope "src/App.tsx" \
  --test-scope "src/App.test.tsx" \
  --acceptance "A Reset button sets the displayed count to zero" \
  --verify "npm test"
```

Use `--request-file docs/examples/contract-first-run.json` for contract-first
multi-repository work. Only inputs marked `ready: true` may start. A dependent
input remains locked until all dependencies pass independent verification.

Before launching a worker, record its start:

```bash
node .agents/bin/cc.mjs record-result --run-id <run-id> --repository <name> --stage worker-started
```

## Codex

Invoke `$run-task`. Give a fresh worker sub-agent only the emitted worker-input
JSON and referenced instructions. The worker operates only in the emitted
worktree, runs the required checks, commits its changes, and writes its result.
Record `worker-result`, then give a different fresh, read-only verifier only the
verifier-input JSON. Record `verifier-result` after it writes its result.

## Claude Code

Invoke `/run-task`. Start a fresh Claude Code worker with only the emitted
worker-input JSON and referenced instructions. After recording `worker-result`,
start a separate fresh, read-only Claude Code verifier with only the emitted
verifier input, then record `verifier-result`.

For non-interactive sessions, put the positional prompt before Claude Code's
variadic `--add-dir`, use `--no-session-persistence`, and expose only the wrapper
and emitted worktree paths required by the input.

## Repair and review

On a failing verdict, run `prepare-repair` and give its emitted input to a fresh
repair worker. Use another fresh verifier afterward. Stop when the configured
repair limit is exhausted. On an all-passing run, use `prepare-review` for each
repository. Push or draft-PR publication still requires explicit authorization;
record only a confirmed publication response.

If any step is interrupted, preserve `.runtime/`, branches, and worktrees.
Never delete, reset, stash, or clean unrecorded work to recover a run.

