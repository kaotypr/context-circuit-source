Implemented p1000 in the `api` worktree at `.worktrees/p1000/api`.

- Added `WriteCSV` with a header and invoice rows for number, customer, status, and total cents.
- Added a test for one issued invoice.
- `go test ./...` passed.
- Committed on `cc/p1000/api` as `38de72d` (`feat(invoice): add CSV export`). The worktree is clean.
- The after-implementation hook’s style check failed with **3 warnings in `internal/invoice`**. The hook said to report the failure; it returned no other instructions.