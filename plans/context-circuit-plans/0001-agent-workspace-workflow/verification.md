# Verification

- git diff --check
- node .agents/bin/cc.mjs validate
- node .agents/bin/cc.mjs validate-plan plans/context-circuit-plans/0001-agent-workspace-workflow
- The repository’s existing test command, with its current Node/JavaScript command-layer tests still passing
- The new fixture-driven Agent Workspace acceptance suite covering every scenario in docs/agent-workspace-workflow.md
- Independent read-only verifier review of runtime invariants, handoffs, lease contention, worktree isolation, and human-gate evidence
