# Implement the Agent Workspace Workflow

Replace the command-centric operating model with a host-neutral Agent Workspace Workflow driven by durable instructions and filesystem state. A root session enters or resumes the workspace, coordinates multiple independent plans, and may delegate bounded work to child sessions. Runtime records provide identity, handoffs, leases, recovery evidence, and worktree ownership without a global current-session. The existing Node/JavaScript command layer remains available as a compatibility adapter until the replacement passes the full acceptance suite; its removal is a separate future cleanup plan.

Plan metadata lives in [plan.yaml](./plan.yaml). The plan status changes only through an explicit human request.

## Documents

- [overview](overview.md)
- [requirements](requirements.md)
- [acceptance-criteria](acceptance-criteria.md)
- [solution](solution.md)
- [delivery](delivery.md)
- [verification](verification.md)
- [risks](risks.md)

## Tasks

See [tasks/README.md](./tasks/README.md).
