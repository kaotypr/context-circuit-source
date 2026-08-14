# Development of Context Circuit

Context Circuit is being refactored from a command-oriented wrapper into an
instruction- and filesystem-driven AI-agent workspace.

The canonical design sources are:

- AGENTS.md for safety and agent behavior;
- WORKFLOW.md for session workflow and human gates;
- docs/agent-workspace-workflow.md for the full workflow contract;
- context/INDEX.md for context navigation;
- agents/ for root, implementation, and verification role instructions.

Host adapters under .agents/, .codex/, and .claude/ should remain thin. The
current Node and JavaScript command layer is transitional migration material.
The target design does not require a central database or a large command
runtime. Any deterministic helper must remain minimal and preserve the
filesystem as the source of runtime truth.

Implementation work should be verified against the acceptance scenarios in the
Agent Workspace Workflow document. The old command layer remains until the
evidence checklist in docs/legacy-cleanup-gate.md passes and a separate
cleanup plan is explicitly requested.
