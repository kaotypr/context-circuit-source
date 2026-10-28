---
product: template
kind: added
pr: 208
---

`agent setup` turns Claude Code's commit attribution off in the gitignored
`.claude/settings.local.json`, and `.gitignore` carries that file. The orient
check names the setting, so a machine without it runs `agent setup`, which is
how existing workspaces receive it. The agents guide names Cursor's personal and
team controls; the written no-attribution rule stays the guard where neither is
set.
