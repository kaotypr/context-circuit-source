---
product: both
kind: added
upgrade: Run `agent setup` after pulling a workspace whose extensions ship skills, and add the new `x-*` skill paths from the template's `.gitignore`.
---

Extensions can now ship skills. A manifest's `skills` names each skill's folder, which holds a `SKILL.md` with `name` and `description`, and `agent setup` installs each skill of an extension on for you as `x-<extension>-<skill>`: a copy under `.agents/skills/`, which Codex reads, and a one-line wrapper under `.claude/skills/` and `.cursor/skills/`. These are gitignored and regenerated on each machine, never replace a shipped skill, and are never touched by a template update; setup keeps a copy you edited and removes one whose extension is gone or off. A skill runs on your request, or when a core hook instruction names it, never on its own. A skill may declare a named stage, which other extensions hook as `before|after <extension>/<stage>` and the skill runs with `extension hooks --stage <extension>/<stage>`; a stage whose extension is not listed or is off for you is skipped with that reason, and named stages never join the core lifecycle. `status` lists `extension_skills` and `named_stages`, `extension list` names each extension's skills and stages, and `check` reports a skill that is not valid and a hook at a named stage no listed extension declares.
