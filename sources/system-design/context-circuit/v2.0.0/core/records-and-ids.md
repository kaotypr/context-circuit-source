# Records and IDs

Intents, plans, identity allocation, dates, and the concurrency limits v2 states
rather than hides.

## Two record kinds, one purpose each

| Kind | ID shape | Written | Holds |
| --- | --- | --- | --- |
| Intent | `iNNN-slug.md` in `intent/` | Before detailed code investigation | Goal, non-goals, constraints, observable success criteria, rough repository scope, approval |
| Plan | `pNNNN-slug.md` in `plans/` | After approval, from real code | Task order, dependencies, risks, expected checks, progress, completion |

The order is the point. An intent authored *before* reading target code states an
outcome in the project's own language, which is the only form a human can
meaningfully approve. A plan authored *after* reading the code is grounded, and
it is earned by the approval rather than gated again.

### Plan frontmatter

```yaml
id: p0001
created_by: maya
intent: i001
repositories:
  - api
  - web
depends_on:
  - p0002
completed: 2026-09-15
```

`repositories` is plural by design (D6). v1 required one plan per repository,
which fragmented a single coherent change across records that had to be
reassembled by hand. v2 lets one readable Markdown plan cover several
repositories, and splits into separate plans when *execution or delivery* wants
them separate — not because the git layout does.

`completed` is the machine-readable half of completion; the human-readable
`## Completion — YYYY-MM-DD` section is written in the same operation. Dependency
ordering reads the date, so ordering never has to parse prose.

## Identity

**Workspace-global, prefixed, permanently reserved.** `i001-add-billing`,
`p0001-billing-api`. Intent IDs carry a minimum of three digits and plan IDs
four; numbering expands beyond that width rather than wrapping. The `p` prefix
exists so a plan cannot be mistaken for an intent in a reference, a branch name,
or conversation.

**Optional allocation bands.** v1 divided the numeric range between members and
*required* it, so every ID carried a member's fingerprint whether or not anyone
needed one. v2 keeps the mechanism and drops the requirement:

```yaml
members:
  maya:
    name: Maya
    band: 1
  alex:
    name: Alex
    band: 2
```

A member holding band N allocates intents from `N*100` and plans from `N*1000` —
band 1 writes `i100`, `p1000`, `p1001`; band 2 writes `i200`, `p2000`. Blocks are
`N*width … N*width+width-1`, so distinct bands can never overlap and band 1 still
lands on the kind's minimum digit width. A member **without** a band allocates
from the numbers no band claims, which is why a solo workspace still counts from
`i001` and `p0001` and pays nothing for the feature.

Three properties keep a band from becoming a namespace:

1. **The ID stays global.** `p1000` is a workspace-global plan ID that happens to
   have come from a block. Nothing parses a member out of it.
2. **A band decides which number comes next, never which numbers were right.**
   Assigning, changing, or clearing a band renumbers nothing and releases no
   reservation, and a reserved number inside a band is still skipped.
3. **Attribution is unchanged.** `created_by` remains the only member metadata.
   There is no assignee, owner, reviewer, or member namespace — attribution
   answers "who wrote this record", not "whose work is this."

Two members on one band share a range, which defeats the only thing a band is
for, so the roster is **refused** until that is resolved rather than allocating
from an ambiguous state. An exhausted band errors naming the member, the band,
and its range.

**Reservations are permanent.** A reserved ID is never reused, including after
archival or deletion. Allocation inventories filenames including archives,
without reading archived content. `check` reports an ID present in a record but
missing from the permanent ledger, and vice versa.

## Allocation and the honest concurrency limit

Within one workspace directory, a portable OS file lock serializes cooperating
commands, and single-file writes are atomic replacements. That is genuine
coordination, and it is the only coordination v2 claims.

Across separate Git clones, **bands are the offline mechanism**: two members
holding distinct bands cannot choose the same number, however long they work
apart. That is the reason to assign them in a team workspace, and the reason
dropping them from the first v2 draft was a mistake.

What bands do not cover is stated rather than glossed:

> Bands prevent collisions only between members who actually hold distinct ones,
> against a current roster. Unbanded members working in separate clones, and any
> clone whose roster is stale, can still allocate the same number. Synchronize
> the shared workspace before allocating. Do not claim distributed collision
> prevention beyond what bands give.

Two people who each allocate `p0007` offline from a stale roster still have a
conflict only they can resolve, and resolving it means fixing every reference and
preserving both reservations. v2 declines a central allocation service — a
dependency it will not take — and describes the remaining limit instead.

## Dates

Every date in a workspace file is an **ISO 8601 calendar date, `YYYY-MM-DD`, in
UTC** — in frontmatter, in approval lines, in progress notes, and in prose, all
reading identically. This is enforced for executable-written dates and stated as
a rule for hand-written ones. The reason is retrieval: a catalog entry's
`reviewed 2026-02-04` and a plan's `completed` field must sort and compare
without a parser, and localized formats in shared files break both.

## A record is written at its gates

Only the two gates append. `record complete` appends `## Completion — <date>`
and stamps `completed_at`; `record approve` appends `## Approval — <date>` and
stamps `approved_at`. Nothing else writes to a record after it is created.

Direct Markdown editing by a person is equally legitimate. The commands exist for
consistency and date discipline, not to own the file.

Through 2.0.0-rc.5 a third command, `record note`, appended `## Update — <date>`,
and it was removed in 2.0.0-rc.6. What an agent put there was the result of
running the plan — the same summary completion writes — parked in the body
because completion is a person's request and had not come yet. Two places
recording one thing disagree the moment either is written first, and the parked
copy was always the one a later reader hit before the gate.

## A plan is an intention, not a history

A plan carries the approach, the tasks and their order, and the risks and checks.
Running it does not edit it. The result of a run reaches the person, and reaches
the record only when they request completion.

This is the difference between a plan record in v2 and an execution record in v1.
The v1 record was constructed to be evidence and was reliably staler than the
diff. The v2 plan does not compete with the diff at all: it says what was
intended, the diff says what happened, and on resume the real branches and diffs
are read first (P1) because a record was never a claim about them.

Progress, remaining work, and PR references are not recorded against the plan
while work is in flight. They are what the session reports to the person, who is
the one deciding whether the work is done.

## Archival

Archival is optional ordinary file organization, done on request. ID-based lookup
discovers records under their `intent/` and `plans/` directories including
archives; listing archives requires `record list --archived`, so normal listings
stay short. Moving a record keeps its reservation and requires fixing handwritten
relative links. There is no archival lifecycle, status, or gate.

## What `check` verifies

`check` is a diagnostic over the record graph and its surroundings:

- invalid or duplicate reserved IDs; duplicate record IDs;
- a record ID missing from the permanent ledger;
- an unknown `created_by` member;
- a plan with a missing intent reference, or not linked back from its intent;
- a plan with no repositories, or naming an unknown repository;
- a broken plan link from an intent;
- dependency cycles;
- relationships referencing unknown repositories;
- local binding problems;
- worktree associations needing attention, or naming a missing plan;
- knowledge-boundary and catalog issues (see
  [knowledge-circuit.md](./knowledge-circuit.md)).

It reports and exits. Nothing in the product waits on it, and it is not a
substitute for code review (D18).
