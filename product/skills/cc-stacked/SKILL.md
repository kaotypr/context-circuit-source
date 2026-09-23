---
name: cc-stacked
description: Run several Context Circuit plans at once — deriving the dependency order, integrating each wave, resolving merges, and knowing where to stop.
---

# Run several plans at once

To execute several plans at once, derive the order instead of guessing it:

```sh
context-circuit-cli --workspace <root> --json record order --mode auto
```

For an explicitly selected run, pass each plan with repeatable `--plan ID`.
Use `--intent ID` for only that intent's linked plans. The filters are mutually
exclusive. The result names `plan_ids` and `selection`; preserve that exact set
through execution and delivery. An unfinished predecessor outside it is blocked,
never silently added.

`record order` reports dependency waves, each plan's starting reference per
repository, the integration merges a dependent plan needs, and which plans in one
wave share a repository. It runs nothing and reserves nothing.

## Two shapes

Waves overlap independent plans and pay an integration merge at each fan-in. A
linear chain stacks every plan on the previous one, needs no merge, and runs
strictly serially.

Present the recommendation with its cost — a single repository with any fan-in
usually favors the chain — and honor an explicit choice. The request to execute
the presented plans authorizes the selected run; use the recommended shape when
the person has not chosen one.

## Then run it

Run to completion without further prompting: prepare each worktree from the
reported start, perform any reported integration merge with ordinary Git,
dispatch a worker per plan, wait, and carry what each worker reports into
completion.

Completion is a later human request. Follow the selected order's waves and use
reported worker results and actual Git ancestry to decide what can proceed;
`record order` reads recorded `completed_at` when present but never marks a plan
complete itself.

## Integrate from what the worker reported

A worker reports the checks it ran and their real outcome, and that report is
what you integrate from. Do not re-read its diff or re-run its checks to satisfy
yourself: that repeats the expensive half of the work and is why delegation stops
paying.

Read the diff where integration needs it — a merge to resolve, or a report naming
a conflict, a failure, or an assumption — not as a routine audit. A report of
failing checks is information: report it to the person and leave the plan
unmarked and unedited rather than repairing in a loop.

None of that changes work this session implemented itself, which this session
checks as usual. The rule is about not repeating a delegate's reading and test
run, not about trusting your own.

Resolve a conflict from an integration merge directly when both sides are in the
person's selected run and their records and diffs are available. Preserve
both plans' behavior, then run the repository's checks — a resolved merge is not
trusted until they pass. Record the resolution in the plan so it can be audited.
A conflict against anything outside the run, or one where preserving both sides
is impossible, is a stop.

## Where to stop

Stop and report, preserving all work, on: failing checks after implementation or
after a resolved merge; a conflict outside the run or one needing a decision; a
worker that cannot complete or returns blocked; worktree preparation that
refuses; skipped dependency reuse whose fallback setup fails; or an order
reporting a cycle, an unknown repository, or a broken link.

Never unwind completed plans. Name which plans finished, which is stuck, and what
is held behind it.

## What confirming the run authorizes

For that run: worktree preparation and the local integration merges that assemble
a dependent plan's base. Implementation commits on a worktree's own branch are
ordinary execution and need no separate confirmation — every plan's work is
committed before it is reported, which is what gives a dependent plan something
to start from and an integration merge something to merge.

It does not authorize push, pull request creation, merging into a base branch,
deployment, or deletion — each of those is still a separate act needing a
person.
