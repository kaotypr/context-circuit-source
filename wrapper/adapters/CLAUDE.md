# Claude Code adapter

@AGENTS.md

Read `WORKFLOW.md` after the shared instructions. Claude Code is a thin host
adapter: its root session uses the canonical two-stage router, and its
Task/subagent primitive maps to the existing bounded writer or independent
read-only verifier packet. Record only provider-neutral `host_evidence`.

Permission prompts, print mode, authentication, memory, MCP settings, and
transcripts remain host-local and never authorize a route, gate, lease, or role
change. If a required Task/subagent cannot be created, return `host-blocked`
and preserve the filesystem-only workflow; do not self-verify.
