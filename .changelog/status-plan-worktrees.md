---
product: cli
kind: added
---

`status` now lists, under each of your open plans in `my_work`, the worktrees this machine prepared for it: the repository, the worktree's path, and its branch. A plan's branch lives in the repository it names rather than in the workspace's own, so an agent asked to check a plan's work could look in the wrong repository and conclude the plan was never built.
