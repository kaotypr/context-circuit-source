Implemented p1000 in the API worktree and committed it as `3f4f5fb` (`feat(invoice): add CSV export`). Added `WriteCSV`, which writes invoice number, customer, status, and total in cents, plus a test for an issued invoice.

`go test ./...` passed, and the worktree is clean. The post-implementation hook returned a parity instruction; I checked both listed behaviors against the code and test.

### Parity

- Plain-text summary still lists every invoice: **kept**.
- Totals stay in cents: **kept**.