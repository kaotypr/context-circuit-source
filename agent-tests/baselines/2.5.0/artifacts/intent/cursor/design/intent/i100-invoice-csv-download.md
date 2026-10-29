---
id: i100
created_by: maya
created_at: 2026-10-29T06:03:00Z
---

# Download invoices as CSV

## Goal

A customer can download their invoices from the API as a CSV file.

## Non-goals

- Changing when an invoice is voided or credited.
- A download surface outside the API.

## Constraints

Recorded product knowledge says invoices export only as a plain-text summary (`context/domains/invoices.md`). That rule and a CSV download do not both hold as written. Which export a customer gets is open question 1.

## Success criteria

- A customer can request their invoices from the API and receive a CSV file of them.
- The plain-text summary behaves as question 1 settles: either it remains available beside the CSV, or the CSV is the export and the summary is no longer the only one.

## Repository scope

- `api` — invoice export lives there today.

## Open questions

1. **The project records that invoices export only as a plain-text summary. Should CSV replace that summary, or sit beside it?**
   _Answer:_

2. **Which invoices can a customer download: issued invoices only, or drafts, voided drafts, and credit notes as well?**
   _Answer:_
