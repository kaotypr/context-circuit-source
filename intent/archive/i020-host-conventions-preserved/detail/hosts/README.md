# Per-host native surfaces

What each supported host actually reads in a project, and what this intent
puts there. Open the host you care about; [design.md](../design.md) has the
shared decisions.

## Reading order

1. [claude.md](claude.md) — `.claude/` agents, rules, skills, settings.
2. [codex.md](codex.md) — `.codex/` agents and config, root `AGENTS.md`, `.agents/skills`.
3. [cursor.md](cursor.md) — `.cursor/` rules, agents, skills, hooks.
