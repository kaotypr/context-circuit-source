---
name: cc-gather-context
description: Resolve the minimum authoritative, plan-scoped context without mutating repositories or external systems.
---

# Gather context

Read `AGENTS.md`, `WORKFLOW.md`, `workspace.yaml`, relevant Product Knowledge pages, repository-local instructions, and the exact source requested by the human. Separate sourced facts, assumptions, contradictions, open questions, acceptance evidence, and risks. If the workspace has no useful source, plan, or repository context, suggest `$cc-idea-brief` rather than fabricating missing context. Do not copy credentials or unrelated conversation history into plan prompts.

Remain read-only. Do not create branches or worktrees, edit files, publish tasks, update statuses, or create runtime lifecycle records.
