---
id: i100
created_by: maya
created_at: 2026-10-29T06:01:54Z
---

# Download invoices as CSV

## Goal

Customers can download their invoices as CSV through the API.

## Non-goals

Changing invoice lifecycle rules or the existing plain-text summary is out of scope unless needed to make the CSV download available.

## Constraints

The project currently records that invoices have a plain-text summary. The CSV download must respect the API's existing customer access boundaries.

## Success criteria

- An API customer can request and download invoice data in CSV form.
- The returned data contains only invoices that customer is permitted to access.

## Repository scope

`api` — expose the CSV download through the API.

## Open questions

1. **Should CSV be offered alongside the existing plain-text summary, or replace it?**
   _Answer:_
