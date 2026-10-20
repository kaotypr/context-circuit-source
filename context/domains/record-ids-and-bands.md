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
happened, and is written from the code itself *after* approval, once a person
asks for planning. One plan may
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

A reserved number is **not reused while it stays reserved**, and archiving a
record or removing its file by hand leaves it reserved. Allocation inventories
filenames including archives without reading archived content, and the
diagnostic reports a number present in a record but reserved in no ID file. The
reverse — a reservation with no record — is not a fault: it is what an
interrupted create leaves behind, and what a hand-removed record leaves too.

Reservations live in one **ID file** per member, named after the member, and a
member's record creation appends to that member's file alone. Bands keep
two members from choosing the same number; separate files keep them from
editing the same lines, which Git would report as a conflict however different
the numbers. A number counts as taken when any member's file reserves it, so a
file may hold numbers outside its member's current band — old authorless
records, or ones made before bands changed — and that is valid. A member ID that
is a Windows device name is refused, because it is a file name on every clone.

A workspace written before this layout is at schema 2 and keeps every
reservation in one shared file. The first changing command run once setup is
complete **splits** it:

- each number goes to the file of its record's author;
- one with no author or no record goes to the member selected on that machine,
  and once setup is complete there always is one;
- the shared file is removed and the schema moves to 3, which older CLIs refuse,
  so the split is a cut-over every member upgrades through.

A shared file merged back later from a clone that had not pulled the split is
split again, skipping numbers a member file already holds. Reading commands and
a deletion preview never split, because neither may change anything.

Two clones that each split before pulling the other give the same authorless
numbers to two different members. That race is accepted: nothing is handed out
twice because allocation reads every file, the diagnostic reports each number
reserved in two files, and a person removes the extra entry.

## Releasing numbers

Deletion is the one deliberate way a number goes back. It removes every intent
and plan one member created, active and archived; or exactly the records a
person names by ID, on a roster of any size; or every record, which only a
roster of exactly one member may do.
It takes their numbers out of every ID file holding them, so the next allocation
hands the lowest free one in the author's band out again.

Records an earlier version wrote without an author are reached by deleting
everything or by naming them, since nothing selects them by member. Naming an ID
also frees a number an interrupted create reserved without writing a file.

That exists for work that should never have been recorded, or no longer should
be, where keeping the gap would only preserve noise:

- a trial run, or an intent the person changed their mind about;
- a member who has used up their band's numbers on abandoned work;
- the records a member leaving the team leaves behind.

Selecting by author or by name is what keeps it from reaching into work nobody
chose, and the refusals keep the selection honest. It deletes nothing while a
record outside the selection would be left naming a deleted one, or while a
worktree on this machine was prepared for a deleted plan. A released number is
about to mean something else, and a stale pointer to it would then be silently
wrong. Naming records is not access control, so any member's record can be
named.

The one link edited rather than refused is a named plan's entry in its own
intent's plan list, since nearly every plan made from an intent is listed there.
Any other link is resolved by naming that record too, and links inside the named
set do not block.

Its first run is a preview: the person sees the records, their authors and
whether they are approved or completed, and the numbers, and agrees to that list
before the confirmed run deletes anything. Branches and worktree directories Git
still holds for a deleted plan's number are reported and left alone, because a
plan later given that number meets them. Worktree preparation will not resume
one of them for the new plan without being told to: an association remembers
when the plan it was made for was created, and a branch whose commits predate
the plan record is taken as earlier work.

The plan-list edits come first, then the record files — plans before the
intents they point to, and a plan folder's `plan.md` after everything else in
it — and the ID files last. Interrupted anywhere, every selected number is still
reserved and every partly removed plan is still a record, so no number is handed
out twice and naming the IDs still reserved finishes it.

The error names that call, also when an ID file could not be rewritten after
every record file was gone. Deleting by member cannot reach a removed record's
number again, which is why the finishing call is by ID, or deleting everything
again. A plan folder that already lost its `plan.md` is no record any more, so
it is named for removal by hand.

What a released number already reached is not recalled. A branch name, a pull
request, or a conversation that carried it now names two records, and a clone
that has not pulled the deletion still holds the originals.

## Who is allocating

Every workspace has a roster, a person working alone included: they are a team
of one, set up once inside initialization, which writes them to the roster with a
band and selects them on that machine. Each other machine says which member is
working, because the author decides the band and the band decides the number.

Setup is complete when the roster can be read, lists this machine's selected
member, and every member on it holds a band. Until then every changing command
is refused naming the step that settles it. The steps themselves stay allowed —
initializing, adding a member, selecting one, and giving one a band — and so do
orientation and the diagnostic, so a person can always see why they are asked.

The agent applies the same rule by reading the roster itself, because editing a
record, implementing, and committing never pass through the CLI, and an older
CLI read an empty roster as needing nobody. A record an earlier version wrote
without an author keeps its missing author and stays valid.

The answer is taken as given. A person could name someone else's member, and
nothing checks: the roster attributes records and divides numbers, and it was
never access control. Building verification into it would cost every member a
ceremony to prevent a mistake nobody gains from.

## Allocation bands

Every member holds exactly one band: a numeric block that member allocates from
alone. A member holding band N takes intents from `N*100` and plans from
`N*1000`, so blocks can never overlap and band 1 still lands on each kind's
minimum digit width. Band 0 is a band like any other, covering intent numbers 1
to 99 and plan numbers 1 to 999 — never number zero — so it is stored explicitly
rather than read from an absent value. A member with no band at all
exists only in a roster an earlier version wrote, and leaves setup incomplete
until they are given one.

Allocation skips every number already taken inside the band, reserved or held by
a record file, so old records sitting in a member's range are never handed out
again. That is what makes picking any unused band safe, band 0 or one holding old
records included.

Suggestions are stricter, because a suggested band should start empty. They:

- begin at band 1 and never offer band 0;
- never offer a band another member holds;
- never offer one whose range already holds a taken number.

Initialization offers band 1; adding a member, or banding one an earlier roster
left without, offers the next such band, and the person confirms it or picks
another.

Bands exist because two people allocating from separate clones is not rare in
the workspace this product is for, and resolving a collision afterwards means
renaming records and fixing every reference to them.

Three properties keep a band from becoming a namespace:

1. **The number stays global.** Nothing parses a member out of it.
2. **A band decides which number comes next, never which numbers were right.**
   Assigning or changing one renumbers nothing and releases no reservation, and
   a reserved number inside a band is still skipped.
3. **Attribution is unchanged.** The author field remains the only member
   metadata — no assignee, owner, reviewer, or member namespace. It answers who
   wrote the record, not whose work it is.

Two members on one band share a range, which defeats the one thing a band is
for, so the roster is refused until that is resolved rather than allocating from
an ambiguous state. An exhausted band errors naming the member, the band, and
its range.

Two clones may band the same members at once, and that race is accepted rather
than coordinated. When their choices clash — a roster left with merge conflict
markers, two members on one band, or one member given two bands — the diagnostic
reports it as needing a person, and nothing picks a winner.

## The honest limit

Within one directory a file lock serializes cooperating commands and single-file
writes are atomic. Across clones, bands are the offline mechanism — two members
holding distinct bands cannot choose the same number however long they work
apart. What bands do not cover is stated rather than glossed: two clones of the
same member, and any clone whose roster is stale, can still collide, so
the shared workspace is synchronized before allocating. There is no distributed
allocation service, and none is claimed.

Owner:

- `context-circuit-source@internal/workspace/records.go` `Store.allocate`,
  `takenNumbers`, and `bandWidth` — how a number is chosen and which band it
  falls in.
- `context-circuit-source@internal/workspace/ids.go` `reservedIDs`, `reserve`,
  and `Store.SplitIDLedger` — where reservations are read and written, and how a
  shared ID file is split into member files.
- `context-circuit-source@internal/workspace/bands.go` `Store.SuggestBands` —
  which bands are offered.
- `context-circuit-source@internal/workspace/workspace.go` `Members` — the
  roster validation the bands are derived from, and `Store.Setup`, which decides
  whether a change may be made and names what is outstanding.
- `context-circuit-source@internal/workspace/deletion.go` `DeleteRecords` —
  selection by author or by ID, the refusals, unlinking a named plan from its
  intent, removing files before releasing numbers from the ID files, and
  `Store.unfinished` and `Store.unreleased`, which name the IDs still reserved
  and the call that finishes an interrupted deletion; `Store.leftovers` reports the branches and worktrees a deleted
  plan's number still has.
