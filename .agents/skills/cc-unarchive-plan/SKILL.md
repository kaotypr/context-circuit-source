---
name: cc-unarchive-plan
description: Restore one archived plan to its active numbered-plan location with collision checks.
---

Use `node .agents/bin/cc.mjs unarchive-plan <archived-reference>`. The action preserves plan and task content and status, refuses active-path collisions, and does not infer or change readiness.
