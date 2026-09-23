---
name: cc-plan
description: Create Context Circuit plan folders from an approved intent or a direct request for a grounded plan.
---

# Write grounded plan folders

Approval authorizes planning, and planning alone. `record approve` returns
`planning_required` naming this work, because approval reporting success is the
decision recorded, not the request finished.

Write a plan in the language recorded for the active member in `members.yaml`,
and in English where none is. A plan is quoted whole into the brief a worker
receives, and that brief tells the worker that what goes into the repository is
English regardless. Identifiers and the project's own domain vocabulary keep
their form whatever language the prose is in, and the slug stays a lowercase
ASCII slug. So do the headings: a plan's structure is English whatever its prose
is, for the reason the slug is.

Compose in that language rather than translating into it — let the sentences
follow that language rather than keeping English clause order, prefer the verb
where English would nominalize, and keep technical vocabulary in the form
engineers there actually say. Where the member records a `tone` in
`members.yaml`, follow it; it overrides that guidance.
`.agents/skills/cc-intent/SKILL.md` carries the same rules with worked examples.

For an approved intent, inspect real code and create linked plan folders:

```sh
context-circuit-cli --workspace <root> record create --kind plan --slug SLUG --title TITLE --intent i001 --repo api
```

For a person's explicit request for a detailed plan of an already specified
outcome, do the same grounding and omit `--intent`. Do not invent an intent.
Neither planning request authorizes execution. Present the complete plan and
wait for a separate request to execute it.

## How many plans

One plan folder may cover one or several repositories. Use separate
plans when that helps execution or delivery — work that lands in different
repositories at different times, or that a person may want to run in stages.

Where plans depend on each other, record it, so the order can be derived rather
than guessed:

```sh
context-circuit-cli --workspace <root> record dependencies --id p0001 --depends-on p0002
```

## What a plan holds

`plans/pNNNN-slug/plan.md` holds the outcome, repository scope, approach, task
order, optional dependencies, useful risks, expected checks, and a Details
section. Add concern files such as `api-contract.md` or `data-model.md` only when
the extra detail helps a person inspect the proposal. Link every supporting file
from Details and assign it in frontmatter, relative to the plan folder:

```yaml
required_files:
  shared: [data-model.md]
  by_repository:
    api: [api-contract.md]
    web: [ui-structure.md]
```

Omit `required_files` when there are no supporting files. Every regular file
inside the folder besides `plan.md` must be assigned. Paths stay inside the
folder, do not pass through symlinks, and refer to regular files. Use Markdown
for text details; diagrams and assets can also live inside the folder. The
entry remains a readable map of all details, including those assigned to other
repositories. Present its links and material choices to the person, then revise
the folder as requested. Re-present any material change made after presentation.

Explain and record newly needed files or repositories within the approved
outcome. Obtain renewed approval only when the intended outcome or its success
criteria materially change; a question that bears on the outcome returns to the
intent record.

Revise the plan during review and present material changes again. Once execution
starts, running a plan does not edit it, so a later reader can see what was
intended and read the diff for what happened.

## When a planner produced it

Transcribe the planner's approach, its tasks and order, and its risks and
checks. Its verdict and its plan shape are answers to you, not plan content:
once the records exist the shape *is* the records, so a plan saying it is a
single plan states what its own existence already says.

A `not-feasible` return is a complete answer — write no plan, and report it with
its reasons. A return saying the outcome or success criteria must change goes
back for renewed approval. The coordinator writes every record either way.
