---
name: run-task
description: Codex adapter for the canonical planless run-task workflow. Use when coordinating an explicit repository task through Codex with a fresh worker and independent verifier.
---

# Codex adapter

Read and follow `../../../.agents/skills/run-task/SKILL.md`. Use Codex sub-agents
with no inherited coordinator turns for fresh worker and verifier sessions, pass
only the generated input artifact, and assign the worker ownership of only the
emitted worktree. Do not redefine normalization, identifiers, Git preparation,
or result contracts here.
