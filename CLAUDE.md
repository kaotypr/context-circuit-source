# Claude Code adapter

@AGENTS.md

Follow the shared `AGENTS.md` and `WORKFLOW.md` contract. Enter or resume work
through the same Context Circuit root or child session, including the same
receipt, lease, worktree, and handoff checks. A Claude Task/subagent maps only
to the bounded worker or independent read-only verifier packet and records
provider-neutral `host_evidence`.

Claude permission prompts, print mode, authentication, memory, MCP settings,
and transcripts remain host-local. They do not replace a human gate or enter
workspace state. If Task/subagent creation is unavailable, report
`host-blocked` and keep the route read-only; never self-verify.

## Optional: slash-invocation of skills

Product skills ship at `.agents/skills/<name>/SKILL.md` and the coordinator
resolves them by path (INV-SKILL-01). Claude Code does not discover
`.agents/skills/`, so to also invoke a product skill directly in this source
checkout — for example `/cc-execute plan 0014` — create per-skill symlinks under
`.claude/skills/` once:

    mkdir -p .claude/skills && for d in .agents/skills/*/; do ln -s "../../$d" ".claude/skills/$(basename "$d")"; done

In this source repository `.claude/` is a tracked maintainer surface (it already
holds source-only assets such as `.claude/skills/cc-test-case/` and
`.claude/agents/cc-human-simulator.md`), not shipped state; it is never part of
`context-circuit-template`. The symlinks are a regeneratable host convenience —
re-run the command after a change that adds or renames a product skill, and
remove dangling links for any that were dropped. This grants no route, role, or
authority the read-as-procedure path does not already carry.
