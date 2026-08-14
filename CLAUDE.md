# Claude Code adapter

Follow AGENTS.md and WORKFLOW.md. The primary interface is the workspace-entry
agent session that starts or resumes work, including root and child sessions.
Claude commands are host adapters only and must not redefine the workflow or
make the legacy command layer the user-facing interface.
