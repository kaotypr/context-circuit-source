---
name: cc-publish-plan
description: Optionally publish an approved plan and its tasks through a provider adapter before execution.
---

Offer publication as an explicit option after plan approval and before execution. If the human chooses it, use `node .agents/bin/cc.mjs publish-plan --plan <reference> --provider <name> --references <urls.json>`. Preserve plan and task IDs and store only returned current external URLs in YAML. Publishing does not change statuses, create activity or publication event records, start execution, monitor external status, or synchronize completion back into the plan.
