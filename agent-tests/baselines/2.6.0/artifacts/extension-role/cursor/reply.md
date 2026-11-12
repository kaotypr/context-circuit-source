I'll have the parity verifier check plan p1000's branch. First I'm reading how dispatch works and what that role expects.The parity verifier checked `cc/p1000/api` against `main`. Both behaviors in the plan's Parity section are still in place.

## Verdict

kept

## Findings

- The plain-text summary still lists every invoice. `TextSummary` in `internal/invoice/export.go` is unchanged: it still formats one invoice as `Invoice %d for %s: %d cents (%s)` using number, customer, `Total()`, and status. `TestTextSummary` in `internal/invoice/export_test.go` still checks that line. The new `CSV` writer is additive and walks the slice it is given.
- Totals stay in cents. `Invoice.Total` in `internal/invoice/invoice.go` still sums `Line.Cents`. `TextSummary` still prints that sum with the word `cents`. `CSV` writes the same `Total()` under the `cents` column. `TestCSV` expects `7,Ada,issued,100` for a line of 100 cents, and `TestTotalSumsLines` still expects 125 from 100 plus 25.

## Not checked

Nothing further. Both parity behaviors are visible in the export and invoice sources and their tests. Tests were not executed.

Context Circuit plugin 0.5.0 is available (this session has 0.2.2). Update with `/plugin marketplace update [private]-internal`, then `/plugin update context-circuit@[private]-internal`.