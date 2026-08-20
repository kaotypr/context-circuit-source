# Session handoff

- session_id: sess-legacy-runtime
- parent_session_id: sess-root-runtime
- root_session_id: sess-root-runtime
- status: blocked
- plan: plans/context-circuit-plans/0019-agent-context-runtime-documents
- task: ACR-0001
- worktree: .runtime/worktrees/context-circuit/agent-context-runtime-documents

## Objective and scope

- Historical Markdown-only handoff retained for compatibility.

## Evidence inspected

- The legacy reader can identify the recorded blocker.

## Changed files

- None.

## Tests and verification

- Legacy fallback fixture.

## Decisions and assumptions

- Do not rewrite active historical runtime.

## Questions, blockers, and limitations

- The structured handoff is unavailable.

## Recommended next action

- Resume after a structured handoff is written by the owning session.
