---
name: cc-archive-plan
description: Move one exact plan to its archive location without changing its status.
---

Use `node .agents/bin/cc.mjs archive-plan <reference>` only after the human explicitly names the plan. Draft, approved, and done plans may all be archived. The action does not inspect task statuses, tests, merges, or external issues. Archived plans live at `archives/plans/<repository-key>-plans/<number>-<slug>/` and are ignored by `whats-next`.
