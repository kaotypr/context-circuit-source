---
product: template
kind: added
pr: 178
---

`.claude/settings.json` allows both launcher commands, so Claude Code runs them
without a prompt. A workspace that already has the file merges the two entries
into it entry by entry, keeping every entry it added, rather than replacing it.
