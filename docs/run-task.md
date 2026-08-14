# Run one manual task

`run-task` is a rare human-selected escape hatch for one exact direct request or
one exact task from one approved plan. It prepares scoped isolated work; it does
not launch agents, push branches, open pull requests, merge, or deploy.
It is not the core approved-plan workflow and never advances plan lifecycle,
holistic review, publication, merge, or closeout state. Use `execute-plan` for
the complete numbered plan.

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

## Prepare one approved-plan task only

Pass a request file whose `source.kind` is `plan`. Its `source.reference` must
resolve beneath `context/plans/` and include the current `plan_version` and
`approved_digest`. Select exactly one stable, dependency-free `work_id`:

```json
{
  "contract_version": 1,
  "source": {
    "kind": "plan",
    "reference": "context/plans/reset-flow",
    "plan_version": 1,
    "approved_digest": "sha256:<64 lowercase hexadecimal characters>"
  },
  "work_ids": ["RESET-001"]
}
```

The approved work breakdown contains a canonical execution-contract JSON block
for every item: repository, implementation and test paths, test policy,
verification commands, and acceptance criteria. Because that block is numbered
plan material, approval covers it and material edits invalidate the digest.
`run-task` derives its brief from those fields and rejects caller overrides.
The repository comes only from the explicit execution contract and must be an
exact current key in `workspace.yaml`; descriptive `area` text is never used for
execution routing.

The conservative initial execution contract allows one dependency-free plan
item per run. It does not accept caller-supplied dependency-completion evidence;
dependent items remain blocked until a later workflow phase provides a trusted
local work-state projection. A verifier result therefore records exactly one
selected work-item outcome and cannot complete siblings or later dependencies.

The existing direct-request forms still allocate an `ADHOC-*` work ID and remain
planless.

Before launching a worker, record its start:

```bash
node .agents/bin/cc.mjs record-result --run-id <run-id> --repository <name> --stage worker-started
```

## Codex

Invoke `$cc-run-task`. Give a fresh worker sub-agent only the emitted worker-input
JSON and referenced instructions. The worker operates only in the emitted
worktree, runs the required checks, commits its changes, and writes its result.
Record `worker-result`, then give a different fresh, read-only verifier only the
verifier-input JSON. Record `verifier-result` after it writes its result.

## Claude Code

Invoke `/cc-run-task`. Start a fresh Claude Code worker with only the emitted
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
