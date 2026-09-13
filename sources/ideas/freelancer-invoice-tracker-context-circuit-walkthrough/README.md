# Freelancer Invoice Tracker — A Natural-Language Context Circuit Walkthrough

This idea describes a small product that helps freelancers manage clients,
create invoices, export them as PDFs, record payments, and see which invoices
are overdue.

The package shows how someone can take the idea through the complete Context
Circuit workflow using ordinary conversation. The reader is assumed to know
nothing about Context Circuit. This is passive source material, not an active
intent, plan, or implementation record.

## Files

1. [prompts.md](prompts.md) contains the exact prompts to give, in order,
   including the two plain-language questions that bridge the idea into system
   design and implementation.
2. [implementation.md](implementation.md) explains what each prompt should cause
   inside an instantiated Context Circuit workspace.

## Starting point

This is a completely new project:

- no product repositories exist;
- no application code or tests exist;
- no READMEs or technical documentation exist; and
- no previous product behavior can be inspected.

The first facts come from the person describing the product. Repository creation
is a later, explicit setup action.

## Intended product shape

- `invoice-tracker-api`: clients, invoices, payment records, PDF generation,
  reminder scheduling, and automated tests.
- `invoice-tracker-web`: invoice editor, dashboard, settings, and browser tests.
- A freelancer can manage clients and reusable business details.
- An invoice has numbered line items, tax, currency, due date, and notes.
- Invoices can be saved as drafts, issued, downloaded as PDFs, marked paid, and
  identified as overdue.
- The first release records payments manually; it does not process money.
- Email reminders require confirmation before they are sent.

These are facts supplied by the example prompts. The workspace should not invent
additional product behavior.

## Use

Start with prompt 1 in the empty directory that will hold the Context Circuit
workspace. Replace angle-bracket placeholders with real names, provider targets,
and pull-request numbers. Give one prompt at a time and review the result before
continuing.
