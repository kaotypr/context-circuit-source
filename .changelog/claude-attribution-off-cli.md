---
product: cli
kind: added
pr: 208
---

`agent setup`, and `init` through it, now turns Claude Code's commit and pull
request attribution off for the workspace: it adds `"attribution": {"commit":
"", "pr": ""}` to the gitignored `.claude/settings.local.json` only when the key
is absent, keeps every other key, the file's formatting and any value a person
set, and refuses a file that is not one JSON object. Cursor accepts this setting
only in a person's own or a team admin's configuration, and Codex has none, so
neither gets a file.
