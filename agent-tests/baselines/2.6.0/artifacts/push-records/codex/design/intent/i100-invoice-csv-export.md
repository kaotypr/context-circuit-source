---
id: i100
created_by: maya
created_at: 2026-11-12T09:47:10Z
approved_at: 2026-11-12T09:47:26Z
---

# Invoice CSV export

## Goal

Customers can download their invoices from the api as a CSV file.

## Non-goals

Voiding a draft and crediting an issued invoice stay as they are. No export format besides CSV and the existing plain-text summary is added.

## Constraints

The change is in the api repository, in internal/invoice, beside the plain-text summary. The module stays free of third-party dependencies.

## Success criteria

- A caller can write invoices as CSV, with the invoice number, the customer, the status, and the total in cents.
- `go test ./...` passes in the api repository.

## Repository scope

The api repository. The work is confined to internal/invoice.

## Open questions

None known.

## Approval — 2026-11-12T09:47:26Z

Approved the invoice CSV export outcome as written, with no additional conditions or requests.
