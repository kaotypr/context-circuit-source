# Cursor Agent maintainer adapter

@AGENTS.md

Follow root `AGENTS.md` and `WORKFLOW.md` for direct source development.
Read and edit files on the active branch; use the source-only
`.agents/skills/cc-source-develop/SKILL.md` when useful. Do not launch the
product coordinator or invoke the CC lifecycle, role-tiering, worktrees,
leases, or runtime execution records for maintainer work.

Product host integrations under `product/` are packaging material. Release
assembly installs them at workspace-root paths in generated workspaces.
Existing source-only host tools remain available for explicitly requested
maintainer testing and never ship. Preserve the shared safety and explicit
delivery requirements, including the prohibition on agent attribution.
