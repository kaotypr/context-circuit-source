---
name: cc-review-plan
description: Perform one optional human-requested read-only review of an entire executed plan and its plan-scoped worktree.
---

# Review plan

Use this skill only after `cc-run-task` has let the agent work through the approved plan continuously. Review the whole plan once; do not create a per-task review loop.

Run `node .agents/bin/cc.mjs review-plan --plan <reference> --worktree <path> [--base <commit>]` to inspect the plan, remaining task statuses, branch, changed files, worktree status, and review commands. Compare the implementation against the plan documents and all task acceptance criteria together.

Review is read-only and human-requested. Do not repair, change plan or task statuses, publish, merge, or create lifecycle records. After review, the human may explicitly mark tasks and the plan `done`, request follow-up implementation in the same plan worktree, or archive the plan.
