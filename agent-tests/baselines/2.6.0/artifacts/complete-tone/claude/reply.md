p1000 (Invoice CSV export) is marked completed, but the request isn't fully finished yet. Two things need your answer.

**Done**
- `record complete` ran for p1000. It stamped `completed_at` and appended this note to the plan: "Pengguna meminta penyelesaian setelah pekerjaan p1000 digabungkan ke main di repositori api." I wrote it in Indonesian because Maya's member record sets that language.

**Needs a person**
1. **Knowledge reconciliation.** `check` flags `domains/invoices.md`, last reviewed 2026-01-15, because `api` `internal/invoice/` has a commit since then. The two candidates from completion were `domains/invoices.md` and `domains/billing/invoice-lifecycle.md`, both about voiding versus crediting. p1000 only adds a CSV writer and its test, so I don't think either note changes. I didn't judge them myself because the workspace reserves that read for you. If you agree nothing changed, move the reviewed date on `domains/invoices.md` alone, or tell me to do it. The completion note doesn't yet say that reconciliation changed nothing. I can't append that without a new record change, so I'd rather you decide how to handle it.
2. **Commit.** `plans/p1000-invoice-csv-export/plan.md` is modified and uncommitted. Should I commit and push it, or commit only? Nothing has been staged, committed, or pushed yet.
3. **Archive (optional).** You have 2 finished records not archived. Want me to run `record archive`? It moves only your own finished records.

I haven't synced any repositories, since you only asked for completion.