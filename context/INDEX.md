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
schema* — and where several are collected together, give each its own bullet
carrying that same explanation. A run of paths separated by commas is the one
shape to avoid: it reads as a string rather than a list, nothing says which path
answers which question, and a reader who cannot tell them apart opens all of
them, which is the scanning a note exists to prevent.

Link freely to other notes here. Which record or which evidence produced a note
belongs in that plan record. The `check` diagnostic reports any line in a note
that crosses this boundary.

## Knowledge units

- [Glossary](glossary.md) — project vocabulary and the code identifiers implementing it · glossary, terminology, vocabulary, term, jargon, naming

### The product

- [What Context Circuit v2 is](product/what-v2-is.md) {context-circuit-source} — what the product is for and which two decisions stay human · product, purpose, knowledge circuit, thesis, gates, approval, coordination · reviewed 2026-09-15
- [Two products and their version lines](product/two-products-and-versioning.md) {context-circuit-source} — why the template and the executable release separately and how a workspace pins one · version, release, pinning, version store, install, seed, migration · reviewed 2026-09-15

### Architecture

- [The executable and agent seam](architecture/executable-and-agent-seam.md) {context-circuit-source} — which work belongs to the Go executable and which to the agent · seam, boundary, executable, refusals, model-blind, judgment, mechanism · reviewed 2026-09-15
- [Command surface and output contract](architecture/command-surface.md) {context-circuit-source} — what the executable exposes and how its output and failures behave · cli, command, json, yaml, exit status, flags, help, documentation parity · reviewed 2026-09-15
- [Workspace files and safe editing](architecture/workspace-files.md) {context-circuit-source} — which records are shared, which are machine-local, and how edits stay safe · file contract, shared, local binding, lock, atomic, document edit, schema · reviewed 2026-09-15

### Coordination domains

- [Record identity and allocation bands](domains/record-ids-and-bands.md) {context-circuit-source} — how intent and plan numbers are chosen and kept unique across clones · id, allocation, band, reservation, ledger, member, offline, collision · reviewed 2026-09-15
- [Grounding an intent and its open questions](domains/intent-grounding.md) {context-circuit-source} — what an intent is written against and how unsettled decisions reach the person · intent, grounding, open questions, numbered, answer, assumption, contradiction · reviewed 2026-09-15
- [Durable notes and the retrieval catalog](domains/knowledge-notes.md) {context-circuit-source} — what belongs in a project note and how a reader finds it again · knowledge, note, catalog, index, glossary, reconcile, durable, boundary · reviewed 2026-09-15
- [Worktrees and environment reuse](domains/worktrees-and-reuse.md) {context-circuit-source} — how isolated working copies are prepared and what is carried into them · worktree, branch, isolation, execution, copy-on-write, clone, dependency, environment · reviewed 2026-09-15
- [Deriving stacked plan order](domains/plan-ordering.md) {context-circuit-source} — how several plans of one intent are sequenced, merged, and released · order, waves, chain, dependency, integration merge, concurrency, stop · reviewed 2026-09-15
- [Subagent roles and host settings](domains/subagent-roles.md) {context-circuit-source} — what each role may do, how a brief is composed, and how model and effort reach the coding host · role, explorer, planner, worker, reviewer, model, effort, dispatch, brief, prompt, integration, review · reviewed 2026-09-15
- [Authorization boundaries](domains/authorization.md) {context-circuit-source} — which actions need a person and what one authorization covers · approval, execution request, authorization, delivery, completion, cleanup, scope, consent · reviewed 2026-09-15

### Maintaining this source

- [Source layout and what ships](maintenance/source-layout.md) {context-circuit-source} — which component owns what, and what never leaves this checkout · layout, ownership, manifest, ships, embed, seed, packaging, history · reviewed 2026-09-15
- [Validation and release checks](maintenance/validation-and-release.md) {context-circuit-source} — how a change here is proven before it is published, and what stays unproven · test, vet, format, ci, release check, build output, publish, unverified · reviewed 2026-09-15
