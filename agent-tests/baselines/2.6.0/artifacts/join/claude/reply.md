Dana is set up on this machine, and `check` is clean.

- **Member:** Added `dana` (Dana Reyes) at band 2, the band the workspace suggested, with the role "backend developer". Dana is selected as this machine's active member.
- **Workspace checkout:** Bound to this checkout on `main`.
- **api repository:** Bound the existing `repositories/api` checkout on `main` with `repo connect`. I didn't clone it. It was already clean and up to date with `origin/main`.
- **Agent setup:** Wrote the `cc-*` role definitions for Claude Code, Codex, and Cursor, plus `.claude/settings.local.json`.
- **Verification:** `check` returns no issues. `status` shows Dana with no open plans or pending approvals.

Nothing is committed. `members.yaml` is modified in the working tree, with Dana added to the roster. Commit it when you're ready, since the roster is shared with the team.

Dana's language and tone are unset, so records will be written in English in the `colleague` register. `sh .context-circuit/cli member presets` lists the tone presets if you want a different register.