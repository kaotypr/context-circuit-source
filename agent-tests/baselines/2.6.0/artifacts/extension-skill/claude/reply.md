The QA pass on p1000 passed: `go vet` and `go test` both succeeded in the worktree.

**Commands run**
- `go vet ./...` in `.worktrees/p1000/api` (branch `cc/p1000/api`): passed (exit 0)
- `go test ./...` in the same worktree: passed (exit 0), `ok example.com/api/internal/invoice`

I changed nothing in the worktree and committed nothing. The before and after extension hooks for `qa-pass/qa` returned no blocking instructions, and the after hook only asked for this command list.