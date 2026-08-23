# Wrapper migration and rollback

`wrapper/manifest.yaml` owns wrapper version and schema support. The classifier
has three outcomes: `compatible`, `migration-needed`, and `blocked`, plus the
read-only `legacy-unknown` label for records without a version.

Compatible updates may replace wrapper-owned prose, additive invariant IDs,
optional fields, and explanatory docs. Required field changes or renamed
meanings are migration-needed and need a current human confirmation. Unknown or
contradictory live state blocks. No migration rewrites accepted context, plans,
leases, dirty worktrees, repositories, or runtime evidence merely to make them
look current.
It never rewrites accepted context or canonical runtime evidence.

Rollback restores the previous wrapper-owned files and manifest. It never
reverts `workspace.yaml`, `context/`, `sources/`, `plans/`, `.runtime/`, or
repositories. Upgrade fixtures cover a legacy live lease, dirty worktree,
pending handoff, approved plan, and interrupted stack.
