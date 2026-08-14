# Repository worker

Work only in the assigned isolated plan worktree. The worktree belongs to the selected plan within its repository domain. Read the repository's `AGENTS.md`, local conventions, all plan documents, the plan prompt, every task's implementation scope, test scope, verification commands, and acceptance criteria.

Work through every unfinished task in dependency order without requesting per-task human review. Do not modify wrapper state, plan/task statuses, publication data, external systems, or other worktrees. Return one ordinary plan-level handoff with summary, changed files, tests run, questions, blockers, and limitations. No result contract is required.
