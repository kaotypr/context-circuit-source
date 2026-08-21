# Host capabilities

Codex, Claude Code, and Cursor Agent are host adapters around the same
filesystem contract. They discover the natural-language skills in
`.agents/skills/` and may provide child-session primitives, secure credential
storage, or optional integration adapters. They do not define a second route or
lifecycle.

The shipped discovery names are `context-entry`, `context-next`, `context-plan`,
`context-execute`, `context-verify`, `context-gates`, and `context-upgrade`.
The human never needs to type these names.

| Host | Child primitive | Safe fallback |
| --- | --- | --- |
| Codex | native child agent or equivalent | filesystem-only workflow |
| Claude Code | subagent/Task equivalent | filesystem-only workflow |
| Cursor Agent | Task/subagent equivalent | filesystem-only workflow |

If a host cannot create a required writer or verifier child, the route is
host-blocked and remains read-only. The agent does not skip independent
verification or silently self-verify.
