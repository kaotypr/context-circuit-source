# Run a task

`run-task` turns an explicit direct request or selected approved-plan work into
scoped, isolated work. It prepares artifacts and Git worktrees; it does
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

## Prepare approved-plan work

Pass a request file whose `source.kind` is `plan`. Its `source.reference` must
resolve beneath `context/plans/` and include the current `plan_version` and
`approved_digest`. Select one stable `work_id`, or list a dependency-ordered set:

```json
{
  "contract_version": 1,
  "source": {
    "kind": "plan",
    "reference": "context/plans/reset-flow",
    "plan_version": 1,
    "approved_digest": "sha256:<64 lowercase hexadecimal characters>"
  },
  "request": "Implement the approved reset flow",
  "work_ids": ["RESET-001", "RESET-010"],
  "dependency_evidence": [],
  "acceptance_criteria": ["The approved reset flow is implemented and verified."],
  "repositories": [{
    "name": "frontend",
    "depends_on": [],
    "scope": ["src/App.tsx"],
    "test_scope": ["src/App.test.tsx"],
    "test_policy": "required",
    "verification_commands": ["npm test"],
    "acceptance_criteria": ["Reset returns the count to zero."]
  }]
}
```

Every dependency outside the selected set requires a `dependency_evidence`
entry with confirmed completion evidence. The command validates the plan files,
approval, selected IDs, order, repository mapping, scope, test expectations,
and acceptance criteria before allocating runtime state or creating worktrees.
The task brief retains the selected plan IDs and approval evidence. The runtime
manifest keeps each selected item pending until its repository verifier records
an outcome, so later or unselected dependencies are never fabricated as passed.

The existing direct-request forms still allocate an `ADHOC-*` work ID and remain
planless.

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
