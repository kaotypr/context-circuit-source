---
product: template
kind: fixed
pr: 209
---

`AGENTS.md` now says that `workspace.yaml`, `members.yaml`, `member.local.yaml`,
and `repositories.local.yaml` sit at the workspace root, never under
`.context-circuit/`. Agents looked for them there, found nothing, and either
searched the workspace or refused to go on for want of a roster.
