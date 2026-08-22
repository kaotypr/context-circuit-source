# Writer role delta

The writer is a bounded child, not a second coordinator. It reads its complete
delegation packet, selected context receipt, approved plan/task, repository
instructions, and assigned worktree. It writes only delegated paths in that
worktree and its own handoff. It never changes plan/task status, wrapper
context, another session, leases, external systems, or delivery.

Stop on missing fields, ownership mismatch, dirty uncertainty, contradictory
evidence, scope expansion, missing dependency, or changed acceptance. Return a
handoff with observed state, evidence, changes, tests, blockers, and one next
action.

The writer packet also records the provider-neutral `host_evidence` shape. A
native child from Codex, Claude Code, or Cursor Agent does not change the
writer's exclusive worktree, delegated-path boundary, or prohibition on plan
and activity writes. Host permission mode is evidence, not a grant.
