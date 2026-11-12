---
id: i100
created_by: maya
created_at: 2026-11-12T09:41:01Z
---

# Let customers download invoices as CSV

## Goal

Customers can download their invoices as a CSV through the API. The CSV gives
customers invoice data they can use in a spreadsheet or accounting tool.

## Non-goals

This does not change how invoices are created, issued, voided, or credited. It
does not replace the existing plain-text invoice summary.

## Constraints

The existing invoice behavior says invoices export only as a plain-text summary.
This outcome adds a CSV download while preserving that existing export.

## Success criteria

- A customer can request an invoice CSV through the API and receive a usable CSV
  file.
- The CSV contains only invoice data the requesting customer is allowed to
  access.
- The existing plain-text summary remains available and unchanged.

## Repository scope

`api` — provide the customer-facing CSV download.

## Open questions

1. **Should one download contain all invoices the customer can access, or only
   one invoice at a time?** I'd recommend all accessible invoices, since that
   makes the export useful for accounting; tell me if you want per-invoice
   downloads instead.
2. **Should each invoice be one CSV row, with its line items represented in
   columns, or should each line item get its own row?** I'd recommend one row
   per invoice for a simple invoice list; tell me if accounting tools need
   line-item detail in this export.
