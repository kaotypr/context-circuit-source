I did not implement p1000. The workspace's `contract` extension stopped the implement stage before any change, and its rule is that a person decides before the stage goes on.

**What ran**
- Read plan `p1000` (invoice CSV export in the `api` repo, linked to intent `i100`).
- `worktree prepare --repo api --plan p1000` created the worktree `.worktrees/p1000/api` on branch `cc/p1000/api`, starting from `main` at `d82563a`. No files were changed.
- `extension hooks --stage implement --when before --plan p1000 --repo api` ran the `contract` check. It failed with `on_failure: stop`: "2 of 14 orders-service behaviors differ from the recorded contract." The hook says it will not re-run until the person decides.

**Decision needed from you**
1. Look into the 2 contract mismatches (which behaviors differ, and whether the recorded contract or the orders-service code is the one that's wrong). Once that's settled, I can re-run the hook.
2. Tell me to proceed anyway. Only you can make that call, since the extension is set to stop.

**Other things to know**
- The plan's approach is still valid as far as I can tell. `api/internal/invoice/export.go` already holds `TextSummary`, and the new CSV writer would go in the same package. Its test is `export_test.go`. The repo's `AGENTS.md` requires `go test ./...` before committing and `type(scope): subject` commit messages.
- Once you decide, I'll make the change in the worktree, run the after-implement hook and `go test ./...`, and commit on `cc/p1000/api`. I won't push or open a PR without your say-so.