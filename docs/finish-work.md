# Finish a run safely

`finish-work` is invoked by a human after repository work is merged or when the
human deliberately abandons it. A merged closeout requires the versioned
`closeout-ready` merge confirmation emitted by `confirm-merge`; it does not infer
merge from tests, verification, or pull-request metadata.

## Prepare closeout

After the human merge, run the exact confirmation command from review preparation.
It verifies the configured default target, full reported merge commit, and exact
reviewed head without changing Git state. Then run its emitted closeout command
without cleanup first:

```bash
node .agents/bin/cc.mjs finish-work \
  --run-id <run-id> \
  --repository frontend \
  --outcome merged \
  --author <author-slug> \
  --pull-request <reference> \
  --merge-commit <commit>
```

Use `--outcome abandoned --reason "<human reason>"` for deliberate abandonment.
The command writes a machine-readable runtime closeout record and one Markdown
contribution under `contributions/general/` for a planless run. The contribution
uses `<UTC timestamp>-<author>-<slug>.md`, contains the task and run identity,
commits, pull requests, verification, deviations, risks, and candidate durable
learnings, and is never overwritten on a retry.

Preparation changes the runtime manifest to `closing`. With no activity provider,
the semantic `task.completed` or `task.cancelled` lifecycle event is recorded as
`skipped`; no external status is fabricated.

With a configured provider, prepare the completion or cancellation event first
as described in [activity lifecycle hooks](activity-lifecycle.md). Required or
unfinished manual actions block closeout before contribution creation; optional
capability absence records a warning and exact manual fallback.

## Save the contribution

Review the contribution and save it using the wrapper's configured policy. Team
mode uses the normal wrapper pull-request workflow; explicitly configured solo
mode may permit a direct commit. `finish-work` itself does not commit, push, open
a pull request, merge, or update an activity tool.

## Remove the worktree

After the contribution is tracked, committed, and clean, rerun the same closeout
command with `--cleanup`. Cleanup proceeds only when:

- the wrapper contribution is durably tracked and has no uncommitted changes;
- the base repository and run worktree are clean;
- the worktree still has the recorded branch and head commit;
- merged work is reachable from the configured default branch, directly or by a
  supplied merge commit; and
- abandoned commits are already merged or preserved by a remote ref.

A blocked cleanup writes the exact blockers and preserves everything. Successful
cleanup removes only the registered worktree, preserves its branch and runtime
evidence, and marks the manifest `closed`. Branch deletion and runtime-evidence
pruning require separate explicit future operations.

Blocked records return an ordered checklist. Resolve it from top to bottom:
save the contribution through the configured wrapper policy, make both Git
roots clean without discarding work, establish merge or remote-preservation
evidence, then run the exact shell-safe `finish-work --cleanup` command printed
as the final item. The rerun rechecks every condition.

## Retention boundary

Normal closeout is not a retention or pruning command. It retains the complete
`.runtime/runs/<run-id>/` evidence, product branch, contribution, and review and
merge records. A future pruning feature would need a separately invoked,
bounded, versioned contract with an explicit run or age selection, dry-run
inventory, and recoverable archive where practical. Until then, preserve closed
runtime evidence; do not bulk-delete `.runtime/` as part of `finish-work`.
