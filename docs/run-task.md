# Plan execution

`run-task` is the explicit execution gate for one approved plan. It reads the plan, Product Knowledge references, repository instructions, plan documents, and all unfinished tasks. It requires one repository domain and a clean repository base when the plan worktree is first created, then creates or reuses one isolated branch/worktree and writes a plan-level Markdown prompt under `.runtime/plans/`.

```sh
node .agents/bin/cc.mjs run-task \
  --plan plans/api-plans/0010-checkout
```

The handoff includes the plan domain, worktree, branch, base commit, prompt, changed-file review commands, and every task's acceptance/test expectations. The agent works through unfinished tasks in dependency order in the same worktree without asking for human review between tasks.

The agent returns one ordinary plan-level summary. No machine-readable result artifact or automatic status change is required. `$cc-review-plan` is the one optional whole-plan review after execution.
