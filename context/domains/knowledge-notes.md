# Durable notes and the retrieval catalog

The knowledge tree is the product's reason to exist, so its rules are enforced
rather than encouraged.

## A note describes the project, not the machinery

A note may never name a plan record, an intent record, or a file of raw supplied
evidence. Two failures motivate this, both observed in practice:

1. **Rot.** Those records are archived, renamed, and rewritten while knowledge
   is meant to outlast them. A citation to one is a dead reference the moment it
   is filed away, and a reader cannot tell a stale citation from a live one.
2. **Reopening passive material.** A recorded evidence path becomes a standing
   instruction to read material that is supposed to stay passive.

The durable anchor is a repository path written with the logical repository ID,
exact or patterned — `api@internal/billing/dunning/`,
`web@src/features/<feature>/` — never a local checkout path, because those are
per-machine. Links between notes are free. Which record or which evidence
produced a note belongs in that record.

The diagnostic reports any line in the tree that crosses this boundary, by line
number, naming the match.

## One note is one unwrapped catalog entry

An entry carries its title and relative link, the repositories it applies to in
braces, the question the note answers, the terms a reader would actually search
for, and the date it was last confirmed against the code.

The entry must be **one unwrapped line**, because retrieval matches whole lines
by case-insensitive substring: a wrapped entry returns a fragment carrying no
link, a match nobody can follow. The primitive matching is deliberate — one
small file needs no index, no embedding store, and no service, and it degrades
gracefully to reading the catalog, which stays short enough to read.

Braces do real work: a brace-delimited ID cannot partially match a longer one,
and a real repository ID is always lowercase so the shipped placeholder in the
seed cannot collide with one. Search terms should separate a note from its
neighbours; a word every entry carries narrows nothing.

**Both halves are written in the same edit.** An entry naming an absent note is
a confident miss; a note no entry names is reachable only by someone who already
knows its filename. The diagnostic reports either half, and skips fenced
examples in the catalog since those teach the shape rather than cataloguing
anything. An absent catalog is not an error — retrieval then falls back to
filenames and search terms.

## The glossary

A term is recorded the first time its meaning has to be asked for, naming the
code identifier whenever it differs from the word the project says out loud —
that mapping is the reason the file exists, since a reader searching the code
for the business word finds nothing without it. Rows stay alphabetical so two
members adding terms from separate clones conflict on one row rather than one
section. A meaning is one sentence; a term needing more becomes its own note
with the row pointing there.

## Gathering and reconciliation

Gathering writes live notes and their catalog entries in place, in one pass.
There is no proposal sidecar, no acceptance record, and no knowledge lifecycle —
a separate acceptance step only produces a queue of knowledge waiting to become
knowledge.

Explicit completion is where the circuit closes. It returns the catalog entries
whose braced repositories intersect the completed work's repositories: a
**candidate set to judge, never a list to rewrite**. Most completions change no
durable concept, and recording that is the normal outcome rather than a skipped
step — a mechanism implying every completion should produce an edit would
produce edits. Where meaning did change, the note and its entry move in one edit
and the reviewed date advances. Several completions together reconcile once
across the set.

The circuit is deliberately non-blocking: later work may begin while an earlier
reconciliation is outstanding. Absence is visible debt, not a hidden gate, since
blocking new work on documentation hygiene trains people to record that nothing
changed.

Owner:

- `context-circuit-source@internal/workspace/inspect.go` `knowledgeIssues` and
  `catalogIssues` — the durable-content boundary and catalog consistency.
- `context-circuit-source@internal/workspace/records.go` `KnowledgeCandidates` —
  the entries offered for reconciliation at completion.
- `context-circuit-source@template/context/INDEX.md` — the catalog and note
  conventions shipped to a workspace.
