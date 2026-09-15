---
name: cc-plan
description: Create Context Circuit plan records from an approved intent — what a plan holds, how many to write, and what stays out of them.
---

# Write the plans an approved intent authorizes

Approval authorizes planning, and planning alone. `record approve` returns
`planning_required` naming this work, because approval reporting success is the
decision recorded, not the request finished.

Without asking again, inspect real code and create the linked plan records:

```sh
context-circuit-cli --workspace <root> record create --kind plan --slug SLUG --title TITLE --intent i001 --repo api
```

## How many plans

One readable Markdown plan may cover one or several repositories. Use separate
plans when that helps execution or delivery — work that lands in different
repositories at different times, or that a person may want to run in stages.

Where plans depend on each other, record it, so the order can be derived rather
than guessed:

```sh
context-circuit-cli --workspace <root> record dependencies --id p0001 --depends-on p0002
```

## What a plan holds

Approach, task order, optional dependencies, useful risks, and the checks the
work is expected to pass. Detailed paths are descriptive planning information,
not hard enforcement gates.

Explain and record newly needed files or repositories within the approved
outcome. Obtain renewed approval only when the intended outcome or its success
criteria materially change; a question that bears on the outcome returns to the
intent record.

A plan is written when the work is planned, and then left alone. Running a plan
does not edit it, so a later reader can see what was intended and read the diff
for what happened.

## When a planner produced it

Transcribe the planner's approach, its tasks and order, and its risks and
checks. Its verdict and its plan shape are answers to you, not plan content:
once the records exist the shape *is* the records, so a plan saying it is a
single plan states what its own existence already says.

A `not-feasible` return is a complete answer — write no plan, and report it with
its reasons. A return saying the outcome or success criteria must change goes
back for renewed approval. The coordinator writes every record either way.
