---
name: cc-knowledge
description: Write and maintain Context Circuit live knowledge — what a context note may say, how it anchors to code, and how its catalog entry stays consistent.
---

# Write live knowledge

On a request to gather knowledge, synthesize durable concepts into live context
notes and maintain the catalog entries beside them. There is no separate
knowledge acceptance lifecycle. Keep raw evidence separate from accepted
knowledge, and keep task progress and temporary results in plans. The executable
can locate catalog entries; semantic interpretation and the writing itself are
yours.

Run the diagnostic after editing notes — it reports any line that crosses the
boundaries below:

```sh
context-circuit-cli --workspace <root> check
```

## A note describes the project, not the machinery

Never name a plan record, an intent record, or a file under `sources/` inside a
note. Records are archived while knowledge outlasts them, and a recorded evidence
path becomes a standing instruction to read material that must stay passive.
Which record or evidence produced a note belongs in that plan.

## Anchors

Repository paths are the durable anchor, exact or patterned, written with the
logical repository ID — `api@internal/billing/dunning/`,
`web@src/features/<feature>/` — never a local checkout path.

An anchor says what it anchors. Write it into the sentence that explains it.
Where a note covering several surfaces closes with a block naming where its
subject lives, head it `Owner:` and give each anchor its own bullet carrying
that explanation. Never a comma-separated run of paths: that names no question
any of them answers and sends the reader to open all of them.

## One note, one catalog entry

A note is one unwrapped catalog entry carrying its own link, the repositories it
applies to, the question it answers, its search terms, and the date it was last
confirmed against the code, as `context/INDEX.md` describes. `check` reports an
entry that claims repositories without that date.

A note lives in a directory named for its concern and is catalogued under a
heading for the same one, so a reader can guess where it is. Only the catalog and
the glossary sit at the top level. Link freely between notes.

Record a domain term in `context/glossary.md` the first time its meaning has to
be asked for, naming the code identifier when it differs from the project's own
word.

## Reconciling after a plan completes

`record complete` returns the catalog entries scoped to that plan's repositories
and names the judging as `reconcile_required`. Judge each against what the plan
actually changed: edit the note and its catalog entry together and move its
reviewed date, or record in the completion note that it changed nothing. Most
completions change no durable knowledge, and recording that is the normal
outcome rather than a skipped step.
