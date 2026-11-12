## Verdict

- **Plain-text summary: kept.** `TextSummary` and its test are unchanged from main.
- **Totals stay in cents: kept.** `Total()` is unchanged, and the new CSV writer emits the total as an integer count of cents.
- **"The plain-text summary still lists every invoice": unclear.** `TextSummary` renders one invoice per call, and neither main nor the branch has a function that lists several. The verifier can't judge this behavior from the code, so the plan's wording needs a look.

## Findings

1. `internal/invoice/export.go` on the branch keeps `TextSummary(i Invoice) string` with the same format as main. The branch adds `CSV` and new imports to that file. `TestTextSummary` in `export_test.go` is identical to main's.
2. `internal/invoice/invoice.go` is the same as main's. The new `CSV` function writes `strconv.FormatInt(item.Total(), 10)`, and `TestCSV` asserts the row `7,Ada,issued,100`, which is cents.
3. The CSV header is `number,customer,status,cents`, with one row per invoice. `TestCSV` covers a single issued invoice, as the plan asks.
4. Neither version of `TextSummary` iterates over invoices. If the plan meant a plain-text listing across invoices, nothing on either side provides one. The person should say which function they meant.

## Not checked

- **Git diff and branch state.** The verifier had read-only file tools, so it compared file contents rather than running `git diff main...`. It also assumed `repositories/api` is checked out on `main` without confirming it.
- **Tests and build.** `go test ./...` was not run, so the plan's own check is unverified.
- **Line-by-line comparison.** `invoice_test.go` was read on the branch only, not compared with main.
- **Error paths.** Write failures and empty input in `CSV` are not covered by any test, and the verifier did not judge them beyond reading.

The verifier changed nothing in the worktree. Running `go test ./...` in the `api` repository would close the unverified test gap, and the plan's Parity wording is the thing to settle with the person.