Implemented `p1000` in the prepared worktree at `.worktrees/p1000/api`.

- Added `WriteCSV` in [csv.go]($RUN/ws/.worktrees/p1000/api/internal/invoice/csv.go), writing invoice number, customer, status, and total in cents.
- Added an exact-output test in [csv_test.go]($RUN/ws/.worktrees/p1000/api/internal/invoice/csv_test.go).
- `go test ./...` passed.
- Committed as `8b13fde` (`feat(invoice): add CSV export`). The worktree is clean.
- The after-implementation hook ran; no extensions apply. I did not push or mark the plan complete.