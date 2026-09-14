# Shared project knowledge

This optional catalog points to durable architecture, conventions, decisions,
terminology, domain rules, and repository relationships. Retrieve relevant
entries instead of scanning every note. Keep task progress in the plan records.

## Writing an entry

One note is one entry, and an entry is one unwrapped line. Retrieval matches
whole lines, so a wrapped entry returns a fragment carrying no link. The line
holds everything needed to choose the note without opening it:

- [Invoice lifecycle](domains/billing/invoice-lifecycle.md) {api} {web} — when an invoice is voided rather than credited · invoices, invoicing, dunning, proration · reviewed 2026-02-04

Title and relative link, then the repositories the note applies to in braces,
then the question the note answers in a few words, then the terms a reader
would actually search for, then the date the note was last confirmed against
the code. Matching is plain case-insensitive substring: a line holding
`invoices` already answers a search for `invoice`, and the braces stop `{api}`
from also matching `{api-gateway}`. Choose terms that separate a note from its
neighbours, since a word carried by every entry narrows nothing.

Group entries under headings and order them within a group. This catalog stays
small enough to read in full; searching it is the fallback once it is not. A
note absent from here is reachable only by someone who already knows its
filename, so the entry and the note are written in the same edit.

## Durable content only

A note describes the project, not the machinery that produced it. Never name a
plan record, an intent record, or a file of raw supplied evidence. Those are
archived and rewritten while knowledge is meant to outlast them, and a recorded
evidence path becomes a standing invitation to read material that is supposed
to stay passive. Anchor to code instead, with the logical repository ID from
`workspace.yaml`, exact or patterned: `api@internal/billing/dunning/`, or
`web@src/features/<feature>/`. Link freely to other notes here. Which record or
which evidence produced a note belongs in that plan record. The `check`
diagnostic reports any line in a note that crosses this boundary.

## Knowledge units

No notes recorded yet.
