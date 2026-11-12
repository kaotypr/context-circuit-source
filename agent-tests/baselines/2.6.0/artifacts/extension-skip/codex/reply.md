Implemented p1000 on the prepared `cc/p1000/api` branch. `WriteCSV` writes a header and invoice rows with totals in cents; the test covers an issued invoice. The plain-text summary remains unchanged.

`go test ./...` passed. Commit: `ba3dab4 feat(invoice): add CSV export`. The worktree is clean.

### Parity

Skipped: the plan lists no behaviors to compare.