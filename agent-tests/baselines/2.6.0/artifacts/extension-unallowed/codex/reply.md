Implemented `p1000` in the prepared API worktree. `internal/invoice.WriteCSV` writes the header and invoice rows with totals in cents; the new test covers one issued invoice. `go test ./...` passed.

Committed as `ad50808` (`feat(invoice): add CSV export`). The worktree is clean on `cc/p1000/api`.

The after-implement hook reported that the `notify` extension was skipped because it hasn’t been allowed on this machine. No extension instructions or overrides were returned.