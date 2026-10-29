I'll execute plan p1000. First I'm reading the plan and the workspace rules that govern how a plan is carried out.The worktree is on `cc/p1000/api`. I'll add the CSV writer beside the plain-text summary and cover one issued invoice with a test.p1000 is implemented on `cc/p1000/api` in the api worktree at `.worktrees/p1000/api`.

`internal/invoice` now has `WriteCSV`, which writes a header and one row per invoice: number, customer, status, and total in cents. The plain-text summary is unchanged. `TestWriteCSVIssuedInvoice` covers one issued invoice (number 7, customer Ada, total 100 cents).

`go test ./...` passed. The change is committed as `436301b` with message `feat(invoice): add CSV export`. The plan record is still open.

The context-circuit plugin in this session is 0.2.4 (installed 0.2.2). Update it with `/plugin marketplace update [private]-internal`, then `/plugin update context-circuit@[private]-internal`.