---
name: cc-knowledge
description: Write and maintain Context Circuit live knowledge — what a context note may say, how it is shaped so a person can read it, how it anchors to code, and how its catalog entry stays consistent.
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

## Anchors belong in one place

Repository paths are the durable anchor, exact or patterned, written with the
logical repository ID — `api@internal/billing/dunning/`,
`web@src/features/<feature>/` — never a local checkout path.

They appear in one place: an `Owner:` block at the end of the note. The body
calls things by their project names — *the goodwill limit*, *the retry
schedule* — and the block says where each one lives.

```
Owner:

- `api@internal/billing/dunning/` — the retry schedule and its wind-down.
- `api@internal/billing/invoice.go` `Void` — the one path that voids rather
  than credits.
```

A path dropped into a sentence is most of what makes a note tiring to read: the
reader stops on a string they cannot act on, and the sentence grows longer to
explain why it is there. Collecting anchors at the end costs one stop and puts
every path where a reader can find it twice. Use that heading and no other; past
a handful of paths a two-column table reads better than bullets, and both are
the same block. Never a comma-separated run of paths, which names no question
any of them answers and sends the reader to open all of them.

`context/glossary.md` is the exception, because mapping a project word to the
identifier the code uses is the reason that file exists.

## Shape a note by what it is saying

Open with the title and one line saying what the note is for. Then `##`
headings, each phrased as the question that section answers, so the headings
read as the questions the note settles.

Match the form to the content:

| When the content is | Write it as |
| --- | --- |
| Several things sharing attributes — states, fields, limits, options, comparisons | a table |
| A flow, a state machine, a layering, a sequence between parts | a `mermaid` diagram |
| Rules or constraints that each stand alone | bullets |
| An ordered procedure, or a chain where each step causes the next | a numbered list |
| Why it is this way, what was rejected, what breaks without it | prose, which nothing else does as well |

A diagram earns its place when the relation between things is the fact worth
recording. One that redraws the list beside it is noise. Keep to `flowchart`,
`sequenceDiagram`, and `stateDiagram`, with short labels and no styling or theme
directives, because some readers see the source rather than a picture.

More than eight lines of unbroken prose usually has a list or a table hiding
inside it.

## An actor note

`actors/` holds one note per actor that deals with the project — a person in a
role, a team, or an external system — and what each one needs from it. These are
the notes an intent is grounded against, so they record what an actor needs
today rather than what the next change will give them.

Group the stories under the business flow they belong to, and head each group
with a link to the `domains/` note describing that flow:

```markdown
# Support agent

Answers customer contact after an order is placed, and is the only actor who can
return money to a customer without finance review.

## [Order lookup](../domains/order-lookup.md)

1. **As a support agent, I want to find an order from a phone number alone,** so
   that I can help a caller who never received a confirmation email.
2. **As a support agent, I want to see where the courier is,** so that "where is
   my order" is answered without contacting the courier.

## [Refunds](../domains/refunds.md)

3. **As a support agent, I want to refund up to the goodwill limit myself,** so
   that a small complaint closes in one call.
4. **As a support agent, I want a larger refund to reach finance with the case
   attached,** so that the customer does not restate it to a second person.

Owner:

- `api@internal/support/` — case handling and the goodwill limit.
- `console@src/routes/orders/` — the surface this actor works in.
```

Four rules hold that shape up:

- **The flow itself is described once, in `domains/`.** An actor note carries
  only that actor's stake in it. A flow retold in every participating actor's
  note becomes several copies that disagree within a quarter.
- **Stories are numbered once through the note**, not per group, so a number
  stays a handle worth citing while groups are renamed and reordered. A story
  that stops being true goes, and the numbers around it stay put.
- **A story says what the project does for that actor now.** A capability
  someone wants is an intent; a wish list kept here rots into a backlog nobody
  trusts, and takes the reviewed date's meaning with it.
- **Everything else is a note like any other**: anchors in `Owner:` only, and
  nothing about the workspace machinery.

An actor earns a note once the project has to ask who a request is for, or once
two actors want incompatible things from the same flow.

## Write it so a person can read it

A note is read by someone deciding something and by an agent about to change
code. Both want the answer first, in a shape they can scan.

- Lead with the answer, then the reason. A paragraph that reaches its point in
  the last sentence is read twice or not at all.
- One idea per paragraph, and sentences short enough to read once.
- Active voice, and the project's own vocabulary rather than invented synonyms.
- State what is true. Hedging reads as uncertainty about the project rather than
  about the sentence.
- Cut any sentence that restates its heading, and any word doing no work.
- Explain a term the first time it appears, or record it in the glossary and
  move on.

The test is whether a reader can answer the heading's question from that section
alone.

## One note, one catalog entry

A note is one unwrapped catalog entry carrying its own link, the repositories it
applies to, the question it answers, its search terms, and the date it was last
confirmed against the code, as `context/INDEX.md` describes. `check` reports an
entry that claims repositories without that date.

A note lives in a directory named for its concern and is catalogued under a
heading for the same one, so a reader can guess where it is. Each heading
carries one line under it saying what belongs there, written before its first
entry, so the next note has somewhere obvious to go; `check` reports a heading
missing that line, and a directory catalogued under two headings. Only the
catalog and the glossary sit at the top level. Link freely between notes.

Record a domain term in `context/glossary.md` the first time its meaning has to
be asked for, naming the code identifier when it differs from the project's own
word.

## Refining a note written earlier

A note gathered by an earlier version of this product, or in a hurry, is often a
wall of prose with paths scattered through it. Rewriting one into the shape above
is ordinary work, with one hard limit: restructuring is not a licence to add
facts.

- Nothing enters the note that was not already in it or confirmed in the code
  just now. A gap the old note left stays a gap.
- Confirm each anchor still exists before it survives the rewrite. Drop one that
  does not, and say so.
- Move the `reviewed` date only when the note was actually re-confirmed against
  the code. A rewrite that changed only the shape leaves the date where it is.
- Report what could not be verified instead of smoothing it into a confident
  sentence. Removing a claim nobody can confirm is a real outcome, not a loss.

The note and its catalog entry move together, as always.

## Reconciling after a plan completes

`record complete` returns the catalog entries scoped to that plan's repositories
and names the judging as `reconcile_required`. Judge each against what the plan
actually changed: edit the note and its catalog entry together and move its
reviewed date, or record in the completion note that it changed nothing. Most
completions change no durable knowledge, and recording that is the normal
outcome rather than a skipped step.
