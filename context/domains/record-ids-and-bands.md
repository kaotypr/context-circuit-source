# Record identity and allocation bands

## Two record kinds

An **intent** states the outcome a change should produce — goal, non-goals,
constraints, observable success criteria, rough repository scope — and is
written *before* detailed code investigation, because an outcome in the
project's own language is the only form a person can meaningfully approve.

A **plan** says how an approved intent maps to real code and what actually
happened, and is written *after* approval from the code itself. One plan may
cover several repositories; it is split when execution or delivery wants it
separate, not because the repository layout does.

A plan's structured header carries its ID, its author, its parent intent, its
repositories, its optional dependencies, and — once complete — a completion
date. That date is the machine-readable half of completion, written in the same
operation as the human-readable completion section, so ordering never parses
prose.

## Identity

Numbers are workspace-global and prefixed, with a minimum of three intent digits
and four plan digits, expanding beyond that width rather than wrapping. The
prefix exists so an intent cannot be mistaken for a plan in a reference, a
branch name, or conversation.

A reserved number is **never reused**, including after archival or deletion.
Allocation inventories filenames including archives without reading archived
content, and the diagnostic reports a number present in a record but missing
from the permanent ledger, or the reverse.

## Allocation bands

A member may hold an optional band: a numeric block that member allocates from
alone. A member holding band N takes intents from `N*100` and plans from
`N*1000`, so blocks can never overlap and band 1 still lands on each kind's
minimum digit width. A member **without** a band takes the numbers no band
claims, which is why a solo workspace keeps counting from the lowest numbers and
pays nothing for the feature. Unbanded allocation steps over every declared
block, or a band would stop preventing anything.

Bands exist because two people allocating from separate clones is not rare in
the workspace this product is for, and resolving a collision afterwards means
renaming records and fixing every reference to them.

Three properties keep a band from becoming a namespace:

1. **The number stays global.** Nothing parses a member out of it.
2. **A band decides which number comes next, never which numbers were right.**
   Assigning, changing, or clearing one renumbers nothing and releases no
   reservation, and a reserved number inside a band is still skipped.
3. **Attribution is unchanged.** The author field remains the only member
   metadata — no assignee, owner, reviewer, or member namespace. It answers who
   wrote the record, not whose work it is.

Two members on one band share a range, which defeats the one thing a band is
for, so the roster is refused until that is resolved rather than allocating from
an ambiguous state. An exhausted band errors naming the member, the band, and
its range.

## The honest limit

Within one directory a file lock serializes cooperating commands and single-file
writes are atomic. Across clones, bands are the offline mechanism — two members
holding distinct bands cannot choose the same number however long they work
apart. What bands do not cover is stated rather than glossed: unbanded members
in separate clones, and any clone whose roster is stale, can still collide, so
the shared workspace is synchronized before allocating. There is no distributed
allocation service, and none is claimed.

Owner:

- `context-circuit-source@internal/workspace/records.go` `Store.allocate` and
  `bandWidth` — how a number is chosen and which band it falls in.
- `context-circuit-source@internal/workspace/workspace.go` `Members` — the
  roster validation the bands are derived from.
