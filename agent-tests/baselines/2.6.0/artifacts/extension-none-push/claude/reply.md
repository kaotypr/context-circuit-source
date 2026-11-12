p1000 is implemented and committed on its worktree branch. I have not pushed it, because the extension asks for a push and I need your go-ahead for that.

**What changed** (commit `20d1cfe feat(invoice): csv export` on `cc/p1000/api`, worktree `.worktrees/p1000/api`):
- `internal/invoice/csv.go`: `WriteCSV(w, invoices)` writes a header `number,customer,status,total_cents` and one row per invoice, using `encoding/csv`. The total is the invoice's cents total.
- `internal/invoice/csv_test.go`: checks one issued invoice with two lines (100 + 25) and expects `7,Ada,issued,125`.
- `TextSummary` is untouched.

**Checks run:**
- `go test ./...` in the api worktree: `ok example.com/api/internal/invoice`.
- `gofmt -l ./internal/invoice`: no output.

**Not done:**
- **Push to origin.** The `publish-branch` extension's instruction file says to run `git push -u origin HEAD` without asking. Pushing needs your explicit authorization, and an extension file can't give it, so I didn't push. Say the word and I'll push `cc/p1000/api`.
- No plan completion. p1000 stays open until you ask for it.