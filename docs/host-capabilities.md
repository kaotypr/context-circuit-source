# Host capabilities

Codex, Claude Code, and Cursor Agent are host adapters around the same
Context Circuit contract. They expose the same capability names and preserve
the same intent, Product Knowledge, plan, runtime ownership, verification, and
human-gate semantics.

Hosts discover the same shared skill catalog: `cc-session-entry`,
`cc-initialize-workspace`, `cc-idea-brief`, `cc-create-prd`,
`cc-gather-context`, `cc-create-plan`, `cc-review-plan`, `cc-run-plan`,
`cc-run-stack`, `cc-whats-next`, `cc-configure-workspace`, `cc-approve-plan`,
`cc-finish-plan`, and `cc-cleanup-runtime`. The table below is a
safe-fallback example for optional configuration, not a per-host command
matrix. Hosts may improve discovery or credential storage. They must not
define a second workflow.

| Host | Capability discovery | Safe fallback |
| --- | --- | --- |
| Codex | `cc-configure-workspace` through the agent skill catalog | Core conversational workflow |
| Claude Code | `cc-configure-workspace` through the agent skill catalog | Core conversational workflow |
| Cursor Agent | `cc-configure-workspace` through the agent skill catalog | Core conversational workflow |

Host mappings may improve discovery or expose native secure credential storage,
but they must not create a second user-facing command workflow. A host that
cannot provide a requested integration reports `unavailable` and continues
with the filesystem-only path.

Host support never changes:

- the portable `workspace.yaml` configuration shape;
- the exclusive plan lease and worktree rules;
- independent verification;
- human approval for plan, delivery, publication, deployment, completion, or
  ambiguous ownership; or
- the prohibition on credentials and external activity records in workspace
  files.

## Child-session primitives

The filesystem packet is the coordination record. The host supplies the child
execution context. Name primitives behaviorally so a later host rename does
not require a new command layer.

| Host | Child-session primitive |
| --- | --- |
| Cursor Agent | Task / subagent tool |
| Claude Code | subagent / Task tool |
| Codex | native child-agent or equivalent |

Cursor's Task/subagent tool is a valid child-session primitive. Do not treat
a Cursor session as having no child primitive.

If the host cannot spawn a child, report the missing host primitive to the
human and ask how to proceed. A missing primitive is not a reason to skip
children.
