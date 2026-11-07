---
name: cc-writing
description: Write intents, plans, notes, and the chat around them so a person can read and decide on them, in the member's language and tone.
---

# Write for a person to read

This skill covers what an agent writes for a person to read and decide on:
intents, plans, the chat message that presents them, open questions, and
approval and completion notes. It stops where that writing stops. When you move
on to a `context/` note, the catalog, the glossary, or a file inside a
repository, neither this skill nor the member's tone applies any more: knowledge
follows `.agents/skills/cc-knowledge/SKILL.md` and is English, always, and a
repository follows its own conventions.

Reading this skill and its reference files is part of writing, and it starts no
investigation. An intent is still written from the request and the knowledge in
`context/`, before any detailed code investigation, as
`.agents/skills/cc-intent/SKILL.md` says. The code is read once the intent is
approved and planning is asked for.

People find agent-written records hard to read even when every fact in them is
right, and length is rarely why. They have to read word by word because the
ideas arrive as a list with nothing joining them, actions hide inside nouns,
identifiers stand in for names, and in another language the sentences were built
in English and carried across. The rules below fix the shape, and a member's tone
sets the register on top of them. A tone never replaces a rule.

## Which files to open

Open only the reference files this member needs, because a tone you are not
writing in is noise that leaks its words into the one you are. Read the member's
`language` and `tone` from `members.yaml`.

The language picks a folder. With no language recorded, or English, open
`.agents/skills/cc-writing/references/en/language.md`. For Indonesian, Bahasa
Indonesia, or Bahasa, open `.agents/skills/cc-writing/references/id/language.md`.
Any other language has no folder yet, so use the rules here together with the
section on languages without a folder below. When the chat and the record are in
different languages, open the folder for the language you are writing at that
moment.

The tone picks one file in that folder, and applies to the record and the chat
alike; the tone file says how it shows in each. A recorded tone that starts
with a preset's name and a colon, such as `casual: ...`, is that preset, even if
its wording is older than the table below. With no tone, write as `colleague`.
Open the matching file, for example
`.agents/skills/cc-writing/references/id/casual.md`; the others are
`.agents/skills/cc-writing/references/en/colleague.md`, `en/casual.md`,
`en/formal.md`, `id/colleague.md`, and `id/formal.md`. `non-technical` changes
what the reader is shown rather than how it sounds, so it opens `colleague.md`
for the wording and `.agents/skills/cc-writing/references/non-technical.md` for
what to show. A tone line that names no preset is the team's own: follow it as written,
on top of `colleague`.

| Preset | Line recorded in `members.yaml` |
| --- | --- |
| `colleague` | `colleague: semi-formal, the way you'd explain it to a teammate; technical terms as the team says them` |
| `casual` | `casual: chat register, short and relaxed; still complete sentences` |
| `formal` | `formal: for managers or people outside the team; no slang or contractions, still plain` |
| `non-technical` | `non-technical: for readers who don't read code; goal, non-goals and questions say what people will see or can do, and the rest stays technical for the team` |

## The rules

Open with the point. The first two sentences say what will change and what the
reader needs to do, so someone who stops there can still answer. The reason and
the detail follow.

Write connected sentences. The reader needs to know how the points relate, and
words like because, so, but, and instead carry that. A list removes them, so use
one only for things that really are separate and parallel: success criteria,
ordered steps, or options to choose from. In a chat message, use one list at
most.

A record is read twice, once to decide and again to check a detail. So in a
record, a section that describes several separate things, such as the parts of
a page or what changes in each area, opens with one sentence that says what they
add up to, then lists them as points. A paragraph that holds five parts in one
run of sentences is the hardest thing in a record to read. The connected-sentence
rule still holds inside each point and in the sentences around the list.

Put a person or a program in front of every action. Write "the CLI checks the
branch", not "validation of the branch is performed". Nouns built from verbs,
and endings like ", ensuring that ...", make text slow to read even when the
reader understands it in the end.

Call things what the reader calls them. A record ID, a file path, or a field name
may follow a plain description, as in "the UI fixes (p1055)", but it never
stands in for one and never opens a sentence. With every ID and path deleted, the
text should still make sense. For another record, take the description from its
title rather than inventing one.

Explain only what this reader is missing. Don't define a word the reader used
themselves, and don't explain the project back to the person who owns it.

Say each thing once. A table that repeats the prose, a summary that repeats the
body, and a heading that repeats its paragraph all make the text harder to
follow.

Keep the evidence out of the decision. An intent says the outcome and what is
still open. What the code looks like, and surveys of how others do it, belong in
the plan or in sources, reduced to a sentence where the reason needs them.

Leave out the habits that read as machine-made: bold on phrase after phrase,
labelled blocks such as "Summary:" or "Assessment:", lists of three by reflex,
"not X but Y", filler adjectives like key, crucial, and robust, and a closing
menu of things you could do next. A record's own format is not one of them: an
open question stays a bold question, as the skill that writes the record says,
and a short bold label opening each point of a list, such as the name of a page
section, helps the reader scan it. What reads as machine-made is bold scattered
through running prose.

## Asking for a decision

Say what you found or propose, what still works as it is, and why you need an
answer before going on. Then give each question as one short sentence in the
reader's own terms, followed by your recommendation and its reason in a
sentence. End with what happens in the meantime or by default, so nobody is
blocked waiting. When something isn't done yet, say how the reader moves it
forward rather than only that it hasn't happened: "The plan isn't written yet.
Once you approve this, ask me and I'll write it", not "Nothing gets planned until
then."

A question that bears on the outcome is the only kind worth asking. Settle
anything the request already answers, and leave implementation choices to the
plan.

## Approval and completion notes

An approval or completion note states the person's decision faithfully, in the
record's language: what they approved or asked for, their answers to open
questions, and any condition they attached. Faithfulness comes before tone. The
tone may shape how the note sounds, but it never adds to, softens, or
reinterprets what was decided.

## Languages without a folder

A recorded language settles which words are used and nothing about how they are
put together. Compose in that language rather than translating into it. Let the
sentence structure follow that language instead of carrying English clause order
and sentence length across; a sentence that hinges on ", so" in English is often
two sentences elsewhere. Prefer the verb where English would nominalize. Keep
technical vocabulary in the form engineers in that language actually say, which
is usually the English word, and keep that choice consistent through the text.
Names are quoted, never translated, in either direction. A record's headings
stay English whatever its prose is, as `AGENTS.md` says.

## The check before showing it

Read it aloud as if you were telling a capable colleague across the desk. If they
would never say a sentence that way, or in that order, rewrite it. Then read the
first two sentences alone and ask whether they carry the point, and search the
text for the habits listed above.
