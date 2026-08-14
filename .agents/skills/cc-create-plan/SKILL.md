---
name: cc-create-plan
description: Create a reviewable numbered draft plan and structured task breakdown from a source.
---

# Create plan

If the workspace has no useful source documents, plans, or registered repository and the request is still an early idea, use `$cc-idea-brief` first. After the Idea Brief is human-confirmed and useful facts are imported into Product Knowledge, read the source, relevant Product Knowledge, repository instructions, and Git status. Distinguish intended behavior from currently implemented behavior. Create the full numbered plan structure under `plans/<repository-key>-plans/` with `plan.yaml`, the seven Markdown documents, `tasks/README.md`, and one Markdown file per task.

Use small, understandable YAML and Markdown. Include only useful fields: plan identity, source, repositories, dependencies, connections, Product Knowledge references, implementation scope, test scope, verification commands, acceptance criteria, and risks. New plans and tasks start as `draft`.

Use `node .agents/bin/cc.mjs create-plan --input <plan-request.json>`, inspect the generated files, then run `validate-plan` and `validate`. Approval is a separate explicit human action: `set-plan-state --plan <reference> --status approved`. Do not publish, execute, or infer status from tests, Git, or agent output.
