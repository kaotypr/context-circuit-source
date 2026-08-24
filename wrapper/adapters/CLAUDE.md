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
