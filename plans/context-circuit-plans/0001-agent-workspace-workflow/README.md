# Implement the Agent Workspace Workflow

Rebuild the command-centric operating model as a host-neutral Agent Workspace Workflow driven only by durable instructions and filesystem state. A root session enters or resumes the workspace, coordinates multiple independent plans, and may delegate bounded work to child sessions. Runtime records provide identity, handoffs, leases, recovery evidence, and worktree ownership without a global current-session. The Node/JavaScript command layer, generated bundle, and command-specific host adapters are removed by this plan; they are not part of the replacement workflow.

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
