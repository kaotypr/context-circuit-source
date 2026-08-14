# Risks

- Filesystem atomicity and lock semantics differ across host platforms, so a lease convention may appear safe while still permitting races.
- The current command layer encodes assumptions about one active plan and may conflict with concurrent session semantics during migration.
- A runtime record can become stale or partially written after interruption, requiring conservative recovery and explicit human resolution.
- Overlapping plans can still edit the same product paths even when their worktrees are isolated; this must be detected and surfaced as a review conflict.
- Host agents may not expose identical subagent spawning, hook, or persistent-session capabilities, so instructions need a degraded but safe path.
- If acceptance tests depend on the old Node/JavaScript layer, they could falsely validate the replacement; the core scenarios must inspect filesystem contracts independently.
- The self-hosted repository’s current workspace.yaml is legacy-shaped, so schema migration must preserve repository ownership and avoid mutating the base checkout.
