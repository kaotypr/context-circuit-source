# Shared project knowledge

This optional catalog points to durable architecture, conventions, decisions,
terminology, domain rules, and repository relationships. Retrieve relevant
entries instead of scanning every note. Keep task progress in the plan records.

## Writing an entry

One note is one entry, and an entry is one unwrapped line. Retrieval matches
whole lines, so a wrapped entry returns a fragment carrying no link. The line
holds everything needed to choose the note without opening it:

```
- [Invoice lifecycle](domains/billing/invoice-lifecycle.md) {REPOSITORY} — when an invoice is voided rather than credited · invoices, invoicing, dunning, proration · reviewed 2026-02-04
```

Title and relative link, then the repositories the note applies to in braces,
then the question the note answers in a few words, then the terms a reader
would actually search for, then the date the note was last confirmed against
the code. A real entry names real repositories; the placeholder above cannot
collide with one, because a repository ID is always lowercase. Matching is
plain case-insensitive substring: a line holding `invoices` already answers a
search for `invoice`, and the braces stop `{api}` from also matching
`{api-gateway}`. Choose terms that separate a note from its neighbours, since a
word carried by every entry narrows nothing.

Group entries under headings and order them within a group. This catalog stays
small enough to read in full; searching it is the fallback once it is not. A
note absent from here is reachable only by someone who already knows its
filename, and an entry naming a note that is not there is a confident miss, so
the entry and the note are written in the same edit. The `check` diagnostic
reports either half when it is missing.

## Where a note lives

A note goes in a directory named for the concern it belongs to, and its entry
sits under a heading for that same concern, so the directory listing and this
catalog tell one story rather than two. `architecture/`, `domains/`, `product/`,
`operations/`, and `references/` are a usual starting set; prefer the names this
project's own concerns suggest, and add a directory when a concern earns one
rather than in advance.

`INDEX.md` and `glossary.md` stay at the top level, because they describe the
whole project rather than one part of it. Everything else earns a directory. A
flat `context/` is what a project looks like before it has concerns to separate,
and it stops being readable well before this catalog does — which is the failure
this layout exists to postpone, since a reader who cannot guess where a note
lives is back to scanning.

## Durable content only

A note describes the project, not the machinery that produced it. Never name a
plan record, an intent record, or a file of raw supplied evidence. Those are
archived and rewritten while knowledge is meant to outlast them, and a recorded
evidence path becomes a standing invitation to read material that is supposed
to stay passive. Anchor to code instead, with the logical repository ID from
`workspace.yaml`, exact or patterned: `api@internal/billing/dunning/`, or
`web@src/features/<feature>/`.

An anchor earns its place by saying what it anchors. Write one into the sentence
that explains it — *the format layer in `api@internal/billing/format/` owns the
schema* — which is where an anchor belongs whenever a sentence can carry it.

A note covering several surfaces may close with an `Owner:` block naming where
its subject lives, one bullet per anchor, each saying what that path owns:

```
Owner:

- `api@internal/billing/dunning/` — the retry schedule and its wind-down.
- `api@internal/billing/invoice.go` `Void` — the one path that voids rather
  than credits.
```

Use that heading and no other, so a reader learns the shape once and finds it in
every note that has one. A run of paths separated by commas is the shape to
avoid: it reads as a string rather than a list, nothing says which path answers
which question, and a reader who cannot tell them apart opens all of them, which
is the scanning a note exists to prevent.

Link freely to other notes here. Which record or which evidence produced a note
belongs in that plan record. The `check` diagnostic reports any line in a note
that crosses this boundary.

## Knowledge units

- [Glossary](glossary.md) — project vocabulary and the code identifiers implementing it · glossary, terminology, vocabulary, term, jargon, naming
