# Overview

Replace the command-centric operating model with a host-neutral Agent Workspace Workflow driven by durable instructions and filesystem state. A root session enters or resumes the workspace, coordinates multiple independent plans, and may delegate bounded work to child sessions. Runtime records provide identity, handoffs, leases, recovery evidence, and worktree ownership without a global current-session. The existing Node/JavaScript command layer remains available as a compatibility adapter until the replacement passes the full acceptance suite; its removal is a separate future cleanup plan.

Source
workflow-specification: docs/agent-workspace-workflow.md

Assumptions
- The workspace filesystem is the shared coordination surface and is available to root and child sessions in the same workspace.
- A session identity is explicit and durable; no implicit process identity or global current-session pointer is trusted.
- A session may coordinate multiple plans, but a plan has at most one active writing owner at a time.
- Repository work is isolated in a dedicated worktree per plan worker, while verifiers operate read-only against a selected worktree.
- Human approval remains the authority for plan approval, exceptional lease recovery, conflict resolution, and completion.
- The old Node/JavaScript command layer is a temporary internal compatibility adapter, not the desired user-facing workflow.
- The implementation may add documentation, skills, fixtures, and minimal migration adapters, but it must not add a central service or credentials.

Open questions
- Which filesystem lease primitive should be the minimum portable contract: exclusive-create directory, exclusive-create file, or a host-provided lock?
- Which host-specific hooks, if any, should be thin adapters around the same session and handoff contract?
- Should the first acceptance harness use POSIX shell fixtures only, or also include a host-agent simulation matrix for non-POSIX environments?
- After replacement acceptance passes, should the legacy cleanup plan remove only the user-facing command layer first or all Node/JavaScript implementation files in one migration?
