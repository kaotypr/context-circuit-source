# Host surfaces

What stays at the workspace root so a host still finds Context Circuit after
the product home nests.

## Mechanism

Codex, Claude Code, and Cursor enter through files they already look for at
the workspace root: shared `AGENTS.md`, then host-specific adapters
(`CLAUDE.md` for Claude), and — for Codex — `.agents/skills/cc-*`.

Those files are **discovery**, not the product. After this change they remain
at the locations hosts already search. Each one points into
`.context-circuit/` (the nested wrapper adapters and the nested role files)
instead of describing a top-level `wrapper/`, `agents/`, or `docs/`.
Product docs live at `.context-circuit/docs`; hosts do not look for `docs/`
at the workspace root, so that folder does not stay behind as a front door.

They must not become a second copy of contracts, runtime, or role guidance.
One owner per rule still holds: the nested wrapper owns the rule; the root
file is a pointer.

## Interfaces

- Shared instruction surface: workspace-root `AGENTS.md` / `WORKFLOW.md`,
  sourced from nested `wrapper/adapters/`.
- Claude adapter: still the thin `CLAUDE.md` import of that shared surface.
- Role files: nested `.context-circuit/agents/` (coordinator, tracer, worker,
  verifier). Root files do not restate those roles.
- Product docs: nested `.context-circuit/docs/` (terminology, templates,
  getting-started). Not a host discovery surface.
- Skill discovery: `.agents/skills/cc-*` stays at the workspace root. It
  does not move under `.context-circuit`.

## Edge cases

- A host that only reads `AGENTS.md` at the repository root must still reach
  the coordinator and the runtime without knowing `.context-circuit` exists
  as a special name.
- Cursor does not require `.cursor/rules` for this behavior; a future scoped
  Cursor rule stays a thin adapter, same as today.
- `.agents/` stays at root, so skill paths hosts already know keep working.
  A leftover move of `.agents/` under `.context-circuit` is a bug.

## Out of this file

How an existing workspace is upgraded onto the nested home:
[upgrade.md](upgrade.md).
