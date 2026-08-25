# Claude Code adapter

@AGENTS.md

Follow the shared `AGENTS.md` and `WORKFLOW.md` contract. Claude Code is a thin
transport. The root Claude session is the coordinator (see `agents/coordinator.md`)
for all normal conversation: orientation, context gathering, planning, review,
approval interpretation, completion, archive/restore, and delivery discussion.

A Claude Task/subagent maps only to the single bounded worker (`agents/writer.md`)
or the independent read-only verifier (`agents/verifier.md`) for one execution.
Record provider-neutral `host_evidence` for the child; a host permission flag is
an observation, not authorization.

Claude permission prompts, print mode, authentication, memory, and MCP settings
are host-local. They never replace a human approval or completion gate and never
enter workspace state. If Task/subagent creation is unavailable, report
`host-blocked` and keep the route read-only; never self-verify.

## Optional: slash-invocation of skills

Product skills ship only at `.agents/skills/<name>/SKILL.md`, and the coordinator
resolves them by path (see INV-SKILL-01). Claude Code does not discover
`.agents/skills/`, so to also invoke a skill directly — for example
`/cc-execute plan 0078` — create per-skill symlinks under `.claude/skills/` once:

    mkdir -p .claude/skills && for d in .agents/skills/*/; do ln -s "../../$d" ".claude/skills/$(basename "$d")"; done

`.claude/` is host-local: it is never part of workspace or shipped state and an
upgrade neither creates nor preserves it. Re-run the command after an upgrade
that adds or renames a skill; remove any dangling links for skills an upgrade
dropped. This is a host convenience only — it grants no route, role, or authority
that the read-as-procedure path does not already carry.
