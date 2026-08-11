---
name: publish-plan-tasks
description: Publish work items from an explicitly approved numbered plan into a configured activity system without duplicates. Use only when a human explicitly asks to review or create external tasks from approved plan work and confirmed read/create capabilities are available or a manual fallback is requested.
---

# Publish plan tasks

1. Read `AGENTS.md`, `WORKFLOW.md`, `workspace.yaml`, and the complete selected plan. Require explicit human publication intent, approved metadata with a valid digest, and configured `read-tasks` and `create-tasks` access.
2. Search the destination read-only for every stable work ID. Normalize only confirmed matches into `plan-publication-discovery`; never infer a mapping from title similarity alone.
3. Run `node .agents/bin/cc.mjs prepare-plan-publication --plan <plan-id> --discovery <json>`. Present the ordered tasks, hierarchy, dependencies, destination, existing mappings, create actions, and idempotency keys. Stop for human review before external writes.
4. Create one proposed task at a time through the authorized session MCP/CLI, preserving the stable work ID in its description or metadata. Follow record order so parents and dependencies exist first.
5. After a confirmed response, run `node .agents/bin/cc.mjs record-plan-publication --plan <plan-id> --work <work-id> --status created --reference <external-ref> --evidence <single-line evidence>`. On failure, record `--status failed`, stop the batch, and report confirmed mappings plus untouched work.
6. Never claim creation without a confirmed external reference. Identical retries are idempotent; conflicting mappings stop. Confirmed references are written to the work breakdown as a non-material approved-plan update.

Do not publish draft plans, create duplicate tasks, store credentials, start implementation, claim tasks, or continue automatically after partial failure.
