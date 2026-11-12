---
name: qa
description: Run the QA pass on a plan's worktree when a person asks for it.
---

# QA pass

Run this only when a person asks for a QA pass on a plan. It checks work that
is already committed; it starts neither planning nor execution.

1. Run `sh .context-circuit/cli extension hooks --stage qa-pass/qa --when before --plan PLAN`
   from the workspace root, and follow what it returns.
2. In the plan's worktree, run `go vet ./...` and `go test ./...`. Change
   nothing there, and commit nothing.
3. Run `sh .context-circuit/cli extension hooks --stage qa-pass/qa --when after --plan PLAN`,
   and follow what it returns.
4. Report whether vet and the tests passed.
