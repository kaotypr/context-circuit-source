---
name: run-task
description: Codex adapter for canonical direct-request and approved-plan execution with fresh scoped workers and independent verifiers.
---

# Codex adapter

Read and follow `../../../.agents/skills/run-task/SKILL.md`. Use Codex sub-agents
with no inherited coordinator turns for fresh worker and verifier sessions, pass
only the generated input artifact, and assign the worker ownership of only the
emitted worktree. Do not redefine normalization, identifiers, Git preparation,
or result contracts here.
