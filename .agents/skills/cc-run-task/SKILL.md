---
name: cc-run-task
description: Execute one explicitly selected approved plan continuously in its plan-scoped domain worktree.
---

# Run plan

Read the exact selected plan, its Product Knowledge references, repository instructions, all plan documents, and all task acceptance criteria. The plan must be `approved`, must have unfinished tasks, must use one repository domain, and its repository base must be clean when the plan worktree is first created.

Before execution, offer the human the explicit option to publish the plan and its tasks with `$cc-publish-plan`. If publication is requested, complete that action before running the plan.

Run `node .agents/bin/cc.mjs run-task --plan <reference>`. The command creates or reuses one isolated branch and worktree for the plan within its repository domain, then writes one plan-level prompt. The agent must continue through all unfinished tasks in dependency order without pausing for per-task human review. Return one human-readable plan summary, changed files, tests run, questions, blockers, worktree, branch, and review commands.

After the agent finishes, a human may explicitly request `$cc-review-plan` once, mark tasks and the plan `done`, or archive the plan. Never infer status from agent output, tests, Git, or publication.
