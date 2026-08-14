# Overview

Rebuild the command-centric operating model as a host-neutral Agent Workspace Workflow driven only by durable instructions and filesystem state. A root session enters or resumes the workspace, coordinates multiple independent plans, and may delegate bounded work to child sessions. Runtime records provide identity, handoffs, leases, recovery evidence, and worktree ownership without a global current-session. The replacement has no Node/JavaScript command layer or generated runtime bundle.

Source
workflow-specification: docs/agent-workspace-workflow.md

Assumptions
- The workspace filesystem is the shared coordination surface and is available to root and child sessions in the same workspace.
- A session identity is explicit and durable; no implicit process identity or global current-session pointer is trusted.
- A session may coordinate multiple plans, but a plan has at most one active writing owner at a time.
- Repository work is isolated in a dedicated worktree per plan worker, while verifiers operate read-only against a selected worktree.
- Human approval remains the authority for plan approval, exceptional lease recovery, conflict resolution, and completion.
- The implementation is limited to documentation, pure agent instructions, filesystem conventions, fixtures, and acceptance checks; it must not add a central service, command runtime, or credentials.

Open questions
- Which host-specific session-entry metadata is available in each supported agent host?
- Which filesystem capabilities should a host document when exclusive directory creation is unavailable?
