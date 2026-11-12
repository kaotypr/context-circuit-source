---
id: i100
created_by: maya
created_at: 2026-11-11T16:35:10Z
---

# Invoice CSV download

## Goal

A customer can download their own invoices from the API as a CSV file a spreadsheet can open. The file lists those invoices and nothing that belongs to another customer.

## Non-goals

Creating, issuing, voiding, and crediting invoices stay as they are. This adds a download format and no new screen.

## Constraints

Invoices export today only as a plain-text summary. A draft is voided; an issued invoice is credited with a credit note. The CSV has to follow those rules rather than invent a third way to close an invoice.

## Success criteria

- A customer who calls the API receives a CSV of their invoices.
- The file opens in a spreadsheet as one row per invoice, with a header row naming the columns.
- Invoices that belong to another customer are absent from that file.
- The plain-text summary still downloads as it does today, unless question 1 is answered the other way.

## Repository scope

The API is the only repository this touches.

## Open questions

1. **Does the CSV sit beside the plain-text summary, or replace it?**
   _Answer:_

2. **Does each row carry the invoice as a whole, or one line item?**
   _Answer:_

3. **Which invoices are in the file: issued ones only, or drafts and voided drafts as well?**
   _Answer:_
