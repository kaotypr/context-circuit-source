# Record identity and allocation bands

Two record kinds, numbered so that two clones never collide.

## Two record kinds

An **intent** states the outcome a change should produce — goal, non-goals,
constraints, observable success criteria, rough repository scope — and is
written *before* detailed code investigation, because an outcome in the
project's own language is the only form a person can meaningfully approve.

That sentence decides the language question too. A member may record the
language their intents and plans are written in, because an outcome someone
skims in a second language is an approval gate in name only. Knowledge does not
follow: a note outlives the member who wrote it, anchors to code, and is found
by substring search that a mixed-language catalog breaks, so `context/` is
English whatever a member records. An approval or completion note keeps the
words the person used, because it is evidence of what someone said rather than
prose about it.

The boundary that makes this safe lives in the dispatch brief, not here. A
record is quoted into the brief as authoritative, and a worker works inside a
repository worktree under that repository's instructions and never reads the
workspace's — so the brief names the author's language and says that code,
comments, identifiers, and commit messages are English. Names are quoted and
never translated in either direction: the failure worth preventing is not a
Bahasa comment, which review catches, but a domain term rendered into English in
an identifier, where it silently disagrees with the glossary that named it.

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

A reserved number is **not reused while it stays in the ledger**, and archiving
a record or removing its file by hand leaves it there. Allocation inventories
filenames including archives without reading archived content, and the
diagnostic reports a number present in a record but missing from the permanent
ledger. The reverse — a reservation with no record — is not a fault: it is what
an interrupted create leaves behind, and what a hand-removed record leaves too.

## Releasing numbers

Deletion is the one deliberate way a number goes back. It removes every intent
and plan one member created, active and archived, or in a solo workspace every
record, and takes their numbers out of the ledger so the next allocation hands
them out again. That exists for work that should never have been recorded — a
trial run, a member's abandoned direction — where keeping the gap would only
preserve noise.

Selecting by author is what keeps it from reaching into anyone else's work, and
the refusals keep the selection honest. It deletes nothing while another
member's record would be left naming a deleted one, or while a worktree on this
machine was prepared for a deleted plan, because a released number is about to
mean something else and a stale pointer to it would then be silently wrong. Its
first run is a preview: the person sees the records and numbers and agrees to
that list before the confirmed run deletes anything.

The ledger is updated before the files go. Interrupted between the two, the
files still stand and allocation still counts their filenames as taken, so no
number is handed out twice, and running the same deletion again finishes it.

What a released number already reached is not recalled. A branch name, a pull
request, or a conversation that carried it now names two records, and a clone
that has not pulled the deletion still holds the originals.

## Who is allocating

A workspace whose roster is empty is a solo one: nobody selects an identity,
records carry no author, and allocation takes the lowest free numbers. Once the
roster lists anyone, each machine says which member is working, because the
author decides the band and the band decides the number. Until it does, record
creation is refused and orientation reports the identity as outstanding; the
agent declines changing requests and asks who the person is, offering an
existing member or adding a new one. A record written before the roster existed
keeps its missing author and stays valid.

The answer is taken as given. A person could name someone else's member, and
nothing checks: the roster attributes records and divides numbers, and it was
never access control. Building verification into it would cost every member a
ceremony to prevent a mistake nobody gains from.

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
  roster validation the bands are derived from, and `ActiveMember`, which tells
  a solo workspace from an unidentified machine in a team one.
- `context-circuit-source@internal/workspace/deletion.go` `DeleteRecords` —
  selection by author, the refusals, and releasing numbers from the ledger.
