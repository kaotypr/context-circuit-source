---
name: cc-finish-work
description: Close a Context Circuit run after human-confirmed merge or deliberate abandonment. Use only when a human explicitly invokes closeout and wants an append-only contribution, lifecycle record, or safe worktree cleanup.
---

# Finish work

1. Confirm that the human explicitly invoked closeout. Read `AGENTS.md`,
   `WORKFLOW.md`, `workspace.yaml`, the selected
   runtime manifest, task brief, current Git state, and any review handoff.
2. Resolve the run ID, repository, author slug, and outcome. The outcome is
   `merged` only after human review and merge. First run the exact `confirm-merge`
   handoff emitted by review preparation; proceed only when it returns
   `closeout-ready` after proving the exact reviewed head and reported merge are
   reachable from the configured default target. `abandoned` requires the
   human's reason. Never infer merge or abandonment from a passing verifier or
   publication record alone.
3. Prepare closeout without cleanup first:

   When an activity provider is configured, first run `node .agents/bin/cc.mjs prepare-lifecycle --run-id <id> --event <task.completed|task.cancelled> [--available <capability>]`. Perform pending actions through the authorized session tool, record only confirmed results with `node .agents/bin/cc.mjs record-lifecycle-action`, and present exact manual fallbacks. A required failure or unfinished manual action blocks closeout without writing a contribution.

   ```bash
   node .agents/bin/cc.mjs finish-work --run-id <id> --repository <name> --outcome <merged|abandoned> --author <slug> [--reason <text>] [--pull-request <ref>] [--merge-commit <verified-sha>]
   ```

   Inspect the schema-valid closeout record and the new append-only contribution.
   Do not overwrite or edit another contribution. Redact secrets and summarize
   evidence; never copy raw conversations or runtime logs.
4. In team mode, save the contribution through the wrapper's normal reviewed Git
   workflow. In explicitly configured solo mode, follow its direct-commit policy.
   Do not push, open a wrapper pull request, or merge unless separately authorized.
5. Only after the contribution is tracked, committed, and clean may the human
   request worktree cleanup by rerunning the same command with `--cleanup`. The
   deterministic command must confirm the base and run worktrees are clean, the
   branch and HEAD still match, and commits are either reachable from the default
   branch or preserved by a remote ref. A supplied merge commit must be reachable
   from the configured default branch.
6. If cleanup is blocked, preserve the worktree, branch, runtime evidence, and
   ordered blocker checklist for recovery. Resolve it in order and use the exact
   shell-safe cleanup rerun printed as its final item. Never reset, clean,
   force-remove, delete a branch, discard unpushed commits, or claim full
   completion after a failed closeout.
7. Successful cleanup removes only the clean registered worktree. It preserves
   the branch and runtime evidence, records the run as closed, and leaves any
   later evidence pruning as a separate explicit operation.

Normal closeout never prunes `.runtime/`. Retention or pruning is not currently
implemented; do not turn `--cleanup` into wholesale evidence deletion.

Activity integrations remain optional. When none is configured, record the
semantic completion or cancellation hook as skipped rather than fabricating an
external status update. Never infer success from issuing a tool call; require a
confirmed response and preserve only its non-secret reference or concise evidence.
