---
name: cc-whats-next
description: Read-only recommendation of one dependency-ready task from approved plans.
---

# Whats next

Run `node .agents/bin/cc.mjs whats-next`. It ignores archived plans, considers only approved plans with unfinished tasks, requires declared plan dependencies to be `done`, and recommends one whole plan to run continuously. It lists the remaining tasks, explains the recommendation, shows the plan and Product Knowledge references, and warns when a recorded Product Knowledge source changed.

When no plan is executable, it recommends reviewing a draft plan or creating one. It never claims work, changes statuses, publishes tasks, or creates lifecycle records.
