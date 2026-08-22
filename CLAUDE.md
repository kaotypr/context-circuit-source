# Claude Code adapter

@AGENTS.md

Follow the shared `AGENTS.md` and `WORKFLOW.md` contract. Enter or resume work
through the same Context Circuit root or child session, including the same
receipt, lease, worktree, and handoff checks. A Claude Task/subagent maps only
to the bounded writer or independent read-only verifier packet and records
provider-neutral `host_evidence`.

Claude permission prompts, print mode, authentication, memory, MCP settings,
and transcripts remain host-local. They do not replace a human gate or enter
workspace state. If Task/subagent creation is unavailable, report
`host-blocked` and keep the route read-only; never self-verify.
