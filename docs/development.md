# Development of Context Circuit

Context Circuit is an instruction- and filesystem-driven AI-agent workspace.

The canonical design sources are:

- AGENTS.md for safety and agent behavior;
- WORKFLOW.md for session workflow and human gates;
- docs/agent-workspace-workflow.md for the full workflow contract;
- context/INDEX.md for context navigation;
- agents/ for root, implementation, and verification role instructions.

The workspace has no Node or JavaScript command layer. The target design does
not require a central database, package manager, or command runtime. Preserve
the filesystem as the source of runtime truth.

Implementation work is verified against the acceptance scenarios in the Agent
Workspace Workflow document with `sh test/acceptance.sh`. The verifier must
inspect runtime files and instructions directly.
