---
product: cli
kind: added
---

`status` now lists, under each of your open plans in `my_work`, the worktrees this machine prepared for it: the repository, the worktree's path, and its branch. A plan's branch lives in the repository the plan names, not in the workspace's own, so this tells an agent where to look for a plan's work.
