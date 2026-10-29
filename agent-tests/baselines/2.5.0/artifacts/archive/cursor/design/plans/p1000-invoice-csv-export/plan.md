---
id: p1000
created_by: maya
created_at: 2026-10-29T05:51:23Z
intent: i100
repositories:
- api
completed_at: 2026-10-29T05:52:07Z
---

# Invoice CSV export

## Outcome and repository scope

The api repository gains a CSV export for invoices. The export lives in internal/invoice, next to the plain-text summary. No other repository changes.

## Approach

Add a CSV writer in the api repository's internal/invoice package. It writes a header and one row per invoice: number, customer, status, and total in cents. A test in the same package covers a single issued invoice. The plain-text summary stays.

## Tasks and order

- Read the invoice type in internal/invoice.
- Add the CSV writer and its test in internal/invoice.
- Run `go test ./...` in the api repository.

## Risks and checks

A CSV that omits the total would disagree with the plain-text summary. The check is `go test ./...` in the api repository.

## Details

This plan has no supporting file. The CSV writer and its test are the whole change, and both live in the api repository.

## Completion — 2026-10-29T05:52:07Z

p1000 has been merged into main in the api repo. Please mark it complete.
