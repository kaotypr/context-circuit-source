# Capability-based activity lifecycle

Activity integration is optional and provider-neutral. The wrapper stores
semantic policy and runtime evidence, while the active Codex or Claude Code
session uses an already-authorized MCP tool or local CLI. No provider SDK,
credential field, or background process is included.

## Configure policy

Declare capabilities and ordered actions in `workspace.yaml`. Every action has
a stable ID, capability, required/optional/manual policy, and exact description:

```yaml
activity:
  provider: example
  access: auto
  required_capabilities: [read-tasks, update-status]
  optional_capabilities: [assign-task, timers]
  lifecycle:
    task.starting:
      - id: refresh-and-claim
        capability: read-tasks
        policy: required
        description: Refresh the task and confirm it is not owned by another contributor.
      - id: set-in-progress
        capability: update-status
        policy: required
        description: Set the external task status to In Progress.
      - id: start-timer
        capability: timers
        policy: optional
        description: Start the external task timer.
```

Required-action failure stops remaining actions. Optional failure is recorded
as a warning and allows unrelated work to continue. Manual actions return their
description verbatim and remain unfinished until confirmed.

## Prepare and record

Pass only capabilities verified in the current host session:

```bash
node .agents/bin/cc.mjs prepare-lifecycle \
  --run-id <run-id> \
  --event task.completed \
  --available update-status
```

The resulting `.runtime/runs/<run-id>/activity/<event>.json` lists pending,
skipped, failed, and manual actions with stable idempotency keys. The host then
performs a pending action and records the confirmed response:

```bash
node .agents/bin/cc.mjs record-lifecycle-action \
  --run-id <run-id> \
  --event task.completed \
  --action set-completed \
  --status completed \
  --evidence "Status transition confirmed by the session tool." \
  --reference TASK-42
```

Identical terminal retries are idempotent; conflicting retries are rejected.
Evidence and references must be single-line and credential-free.

## Start and closeout gates

When configured starting actions are pending, `run-task` writes preflight
evidence but creates no branch or worktree. After actions complete,
`node .agents/bin/cc.mjs run-task --resume-run <run-id>` revalidates repository cleanliness and the
base commit before creating isolated work.

Configured completion or cancellation actions must be completed or safely
skipped before `finish-work` creates a contribution. Failure preserves the run
and worktree for recovery. With `activity.provider: none`, semantic events are
recorded as skipped and planless execution remains fully available.
