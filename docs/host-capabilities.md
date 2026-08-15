# Host capabilities

Codex, Claude Code, and Cursor Agent are host adapters around the same
Context Circuit contract. They expose the same capability names and preserve
the same intent, Product Knowledge, plan, runtime ownership, verification, and
human-gate semantics.

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
