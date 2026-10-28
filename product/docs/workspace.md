# Workspace files

The CLI carries the blank workspace as embedded assets. Initialization creates
the shared instruction and folders alongside these small records:

| File | Contents | Shared? |
| --- | --- | --- |
| .context-circuit/VERSION | Workspace template version | Yes |
| .context-circuit/CLI_VERSION | CLI version this workspace pins; installed side by side | Yes |
| .context-circuit/role-tiering.yaml | Per-host role model and effort preferences | Yes |
| .context-circuit/role-tiering.local.yaml | This machine's overrides of those preferences | No |
| workspace.yaml | Version, name, purpose, the workspace's own repository, repository IDs with their URL and default branch, borrowed knowledge repositories, an optional knowledge review threshold, relationships | Yes |
| members.yaml | Member ID to display name, allocation band, optional record language and tone, and optional informational roles; lists at least one member | Yes |
| .context-circuit/ids/MEMBER.yaml | One member's permanent intent and plan ID reservations | Yes |
| intent/iNNN-slug.md | Intent content, created_by, created_at, approved_at, approval note, linked plans | Yes |
| plans/pNNNN-slug/plan.md | New plan entry, repositories, dependencies, optional intent and required_files, created_by, created_at, completed_at | Yes |
| plans/pNNNN-slug/* | Optional supporting design files assigned in the plan entry | Yes |
| plans/pNNNN-slug.md | Legacy single-file plan, still readable and editable | Yes |
| intent/, plans/ group and _archive/ folders | One optional level of group folder, with `_archive/` above or below it; a record under `_archive/` is archived, and a top-level `archive/` is the legacy archive the first changing command renames | Yes |
| context/ | Optional durable project notes and catalog | Yes |
| member.local.yaml | Active member ID on this machine; written by init, and by member use on every other machine | No |
| repositories.local.yaml | This machine's workspace checkout, repository ID to its checkout path and base branch, and each borrowed knowledge repository's checkout path | No |
| .context-circuit/local/worktrees.yaml | Optional plan/worktree associations | No |
| .context-circuit/local/write.lock | OS-managed edit lock | No |
| repositories/, .worktrees/ | Local Git working copies | No |

Example shared repository definition:

```yaml
version: 3
name: Acme
purpose: Billing software
workspace_repository:
  url: git@github.com:acme/acme-workspace.git
  default_branch: main
repositories:
  api:
    url: https://git.example.com/acme/api.git
    default_branch: main
  web:
    url: https://git.example.com/acme/web.git
    default_branch: main
relationships:
  - from: web
    to: api
    description: Consumes the billing API
```

Use a separate checkout for each execution environment (native Windows, WSL,
remote server, or container). Shared files synchronize through Git; local paths,
worktrees, installed dependencies, and uncommitted work do not migrate automatically.
A local binding holds this machine's checkout path and the branch that machine
starts work from:

```yaml
workspace:
  path: .
  base_branch: main
bindings:
  api:
    path: ../api
    base_branch: release/24.4
```

The base branch is local because it differs per machine: one checkout follows a
release branch while another stays on the default. Worktrees branch from it, and
a binding that records none falls back to the repository's shared
`default_branch`. `repo base` changes this machine's base; `repo remote` changes
the shared URL and default branch.

The two are easy to confuse, and confusing them sends a pull request to the wrong
branch, so `status`, `repo inspect`, `worktree prepare`, and `worktree inspect`
all report `base_branch` beside the branch currently checked out. Read the base
from one of those rather than from `workspace.yaml`, whose `default_branch` is
the repository's default and not what this machine delivers to.

## The workspace's own repository

The shared records travel through the Git repository carrying the workspace
itself, so a member reads plans from whatever commit that checkout sits on.
`workspace connect` describes it: where it lives and which branch it defaults to
go into `workspace_repository`, and this machine's checkout root and base branch
into the `workspace` binding. `workspace base` changes this machine's base and
`workspace remote` changes the shared record, the same split `repo base` and
`repo remote` make. Once described, `status` reports the workspace's branch,
base, and uncommitted records beside every repository's, and `check` stops asking
for it.

`--path` defaults to the workspace root and needs naming only when the workspace
sits inside a larger repository, where it names the containing root. The URL is
read from that checkout's `origin` unless `--url` gives one.

It is deliberately not an entry in `repositories`. Every consumer of that map
reads an entry as somewhere work happens — a plan names repositories, a
relationship joins two, a worktree is cut from one — and a worktree of the
workspace would duplicate the records the CLI is reading. `repo connect` refuses
the workspace's own checkout and names `workspace connect` instead.

No command commits, pushes, or merges the workspace. Members share one branch of
it; a member who keeps the workspace on another branch records that in the local
binding, where it describes that machine and nobody else.

## Member roles

A member may carry `roles`, a list of free text saying what they contribute,
such as `[maintainer, docs]`, shown wherever the roster is shown. Roles are
informational and decide nothing: approval, allocation, deletion, and dispatch
behave the same with or without them, and the roster is still not access
control. `member roles` sets or clears them, `member add --role` and
`init --member-role` record them with the member. Each role is one nonempty line
of text, like a member's name, with no list of allowed values; a role given
twice is kept once.

Roles are written only at schema 3. A CLI too old to know the field reads the
roster strictly, so roles wait until the workspace is at a schema that CLI has
already refused to work in. `check` reports roles found in a schema 2 roster.

## Language

What gets written down follows the artifact rather than the conversation. A
member may record the language their intents and plans are written in; unset
means English.

| What is written | In what language |
| --- | --- |
| An intent or a plan | its author's recorded language, English where none is |
| An approval or completion note | exactly what the person said |
| `context/` notes, the catalog, the glossary | English, always |
| Anything written inside a repository | that repository's convention, English where it states none |
| A record's headings and field names | English, always |
| IDs, slugs, repository IDs, branch names, dates | unchanged |

Knowledge is English because a note outlives the member who wrote it and anchors
to code, and because the catalog is searched by substring, which a
mixed-language index quietly breaks. An intent is in its author's language
because a person can only approve an outcome they understand.

A member may also record a `tone` beside their language: the register that
language is written in, as one line of prose — `semi-formal; keep technical
terms in English` — quoted into the instruction that writes a record rather than
parsed. It exists because naming a language decides the words and not how they
are put together, and an agent composing in English and translating produces
prose that is grammatical and formal in a way nobody chose. Unset leaves the
general composition guidance in `cc-intent` to decide.

`agent dispatch` names the record author's language in the brief and states that
what goes into the repository is English. That boundary lives in the brief
because a worker works inside a repository worktree, under that repository's
instructions, and never reads the workspace's.

Names are quoted, never translated, in either direction: domain vocabulary keeps
the project's own form inside an English note, and code identifiers keep the
code's form inside a record written in another language. Slugs stay lowercase
ASCII whatever the title says, so a title in another script is transliterated.

## When a note stops being true

A catalog entry carries the date its note was last confirmed against the code.
`check` reads that date against the code the note itself points at, through the
anchors in its `Owner:` block: a note with commits under its anchors since it was
last confirmed is reported, with how many.

The comparison is against the code rather than the calendar because age is not
evidence. A note whose anchors nobody has touched is not stale however old its
date is, and reporting it would produce a list that never reaches zero and
teaches people to scroll past findings. Commits on the day of the review do not
count: a note confirmed that day was confirmed against them.

This is the second way a note comes back for judgment, and it covers what the
first cannot. Completion returns the entries scoped to a plan's repositories,
which only ever fires for work done through this workspace; a merge by somebody
else, a commit from before the workspace existed, or a hotfix pushed directly
reaches nothing. The finding is always a person's to settle — re-read the note,
then edit it with its entry or move the date alone if nothing it says changed.

A note that anchors to no code has no evidence to read. Calendar age is all there
is for it, and it is reported only where a workspace asks:

```yaml
knowledge_review_days: 180
```

Unset, nothing is reported for those notes. It is shared, because how long a fact
may go unconfirmed is a project judgment rather than one machine's.

## Borrowed knowledge

`context/` holds knowledge this workspace owns: it describes this project, and
completing a plan brings the notes that plan changed back for judgment. An
organization's knowledge center works the other way. It is shared by many
workspaces, it is changed through its own repository, and no plan completed here
will ever invalidate it. So it is mounted rather than copied:

```yaml
knowledge_repositories:
  core-service-knowledge:
    url: git@git.example.com:platform/core-service-knowledge.git
    default_branch: main
    index: index.md
```

Copying is the failure this prevents. A note duplicated into `context/` goes
stale silently while still reading as current, and nothing in the workspace can
tell that it has. Mounting keeps one copy, upstream, and reads it from there.

It sits beside `repositories` for the reason `workspace_repository` does: every
consumer of that map treats an entry as somewhere work happens, and a plan, a
relationship, or a worktree cut from a knowledge repository is never what the
caller meant. A borrowed repository records no base branch, because work never
starts in it. Both maps share one ID namespace, so a note's anchor means one
thing; connecting an ID already used on the other side is refused.

`index` names the file that maps its contents, the way `context/INDEX.md` maps
this workspace's own notes. A borrowed repository chooses its own entry point,
so the value is recorded rather than assumed; `knowledge connect` detects the
usual spellings and `check` reports the gap when it cannot.

The checkout path is local, like every other checkout path. `knowledge clone`
puts it under `repositories/<id>`, beside the working copies, which the shipped
`.gitignore` already excludes; `--path` keeps it anywhere else. The record is
what separates knowledge from work, not the directory, and the shared ID
namespace is what makes one directory safe.

**Read-only here is enforced rather than asked for.** Obtaining one disables its
push URL, and every sync reasserts that, so a push fails instead of
half-succeeding. `knowledge sync` fetches and then fast-forwards only when the
checkout is clean and on the shared branch. It never merges, rebases, resets, or
discards: a checkout holding local work is reported and left exactly as it is,
because moving that work is a decision for whoever made it. Improvements go
upstream through that repository's own review, from a separate checkout of it.

`context find` reads borrowed indexes beside this workspace's catalog and marks
which side each match came from, so an entry to edit is never confused with one
to raise upstream. It reports any index it could not read rather than narrowing
the result in silence. Nothing validates borrowed content the way `check`
validates `context/`: those rules exist so a person here can fix what they
break, and neither half holds for a repository this workspace does not own.
Completion never names a borrowed entry as knowledge to reconcile, because
nobody here can move it.

## Joining a workspace

A clone carries the shared records and none of this machine's state. `check`
names what is missing: an active member, the workspace binding, and a binding
per repository. Until a member is selected, `status` reports `identity_required`,
every changing command is refused, and the agent declines changing requests;
reading still works. Add the member to the roster with a band if nobody has,
select them, run
`workspace connect`, obtain each repository, and run `agent setup` to write the
role definitions, which are gitignored and so never arrive with the clone.

Working copies belong under `repositories/<id>`, which the shipped `.gitignore`
excludes. `repo clone` takes no `--url` for an ID the workspace already
describes — the shared record carries it, and nothing shared is rewritten, which
is what obtaining a described repository should mean. `repo init` is refused for
such an ID, because an ID already naming a repository is not one to start empty.

After cloning a shared workspace onto another machine, use `member use`,
`workspace connect`, and `repo clone` or `repo connect` with the existing IDs. Do not initialize it again. `repo connect`
writes the shared record only when the ID is new, taking the URL from the
checkout's `origin` unless `--url` names one. A new repository can be registered
before its first commit; worktree preparation needs a commit. Setting a base
records the intended branch without creating or resetting it.

A record's instants — `created_at`, `approved_at`, `completed_at` — are canonical
ISO 8601 UTC timestamps, `2026-09-15T10:53:00Z`; every other date is an ISO 8601
calendar date, `YYYY-MM-DD`, in UTC. Both are written bare in YAML and unadorned
in prose, so an instant in frontmatter and the same instant heading its note read
identically.

YAML changes use goccy/go-yaml document edits to retain surrounding comments and
ordering. Presentation can normalize; byte-for-byte formatting preservation is
not promised. Unknown fields in structured control records and duplicate mapping
keys produce errors. User notes remain ordinary Markdown with no schema registry.

Atomic single-file replacements and a portable OS lock protect cooperating
commands in one directory. Multi-file operations may leave useful partial output
on interruption; IDs are reserved before file creation and are not rolled back.
`record delete` releases IDs deliberately, from every ID file holding them,
only after it has removed the record files — plans before intents, and each
`plan.md` after the rest of its folder — so an interrupted deletion leaves every
selected number reserved and every partly removed plan still a record. `--all`
empties every ID file and keeps the files. `--id` first takes each deleted plan
out of its own intent's `plans:` list. A deletion interrupted at any point is
finished by running `record delete --id` with the IDs still reserved, which the
error names, or the same `--id` or `--all` call again; a plan folder already
emptied of its `plan.md` is named for removal by hand. Deleting another
member's records edits their ID file, so if they create a record before pulling
the deletion, Git reports a conflict in that file; keeping the union of both
sides is safe, since a released number left reserved is only a gap.
Check the reported path and resume or repair the named files. Locks are released
when the process exits, including abnormal exit; the empty local lock file remains.

Every member holds an allocation band, which divides the numeric range so
members can allocate offline without colliding. A member holding band N takes
intents from `N*100` and plans from `N*1000`. Band 0 is a band like any other,
written as `band: 0`, and covers `i001`-`i099` and `p0001`-`p0999`. Allocation
skips every number already taken inside the band. Bands must be distinct, band 0
included, and a roster with two members on one band is refused until it is
resolved. A member with no `band` key, which only a roster written by an earlier
version holds, leaves setup incomplete: every change is refused until that
member is given one. Changing a band leaves existing records and their
reservations untouched.

```yaml
members:
  maya:
    name: Maya
    band: 1
  alex:
    name: Alex
    band: 2
```

## ID files

Each member's reservations live in their own file,
`.context-circuit/ids/MEMBER.yaml`, holding an `intents` list and a `plans`
list. Creating a record appends to its author's file and to no other, so once
every clone has the per-member layout, two members who create records on
separate clones change different files. The one-time split below is the
exception: it writes every author's file at once. Allocation treats a
number as taken when any member's file reserves it or any record file holds it,
archived ones included.

A member's file may hold numbers outside that member's current band, such as
records an earlier version wrote without an author, or ones made before bands
changed. That is valid, and `check` does not report it. A member ID that is a
Windows device name (`con`, `nul`, `com1`, and the like) is refused, because no
clone on Windows could check its file out.

A workspace an earlier CLI wrote is at schema 2, with every reservation in one
shared `.context-circuit/ids.yaml`. The first command that changes something,
once a member is selected and every member holds a band, splits that file: each
reservation goes to the file of the member its record names in `created_by`,
and one with no author, or whose record file is gone, goes to the member
selected on this machine. The shared file is removed and `workspace.yaml` moves
to `version: 3`. CLIs before 2.4.0 refuse a schema 3 workspace with "requires a
newer CLI", so every member needs 2.4.0 from then on. The setup step that
completes setup — giving the last member a band, adding the last member with
one, or selecting a member — makes the split itself, right after it succeeds.
Reading commands, `repo fetch`, which changes only Git refs in a bound checkout,
a `record delete` preview, and a command that refuses its own arguments never
split. The split's notice names every shared file the command
changed, `members.yaml` included when the command changed it, and those are the
files to commit.

Because the split writes a file for every author, only one person on a team
upgrades first. That person makes the first changing command, which splits the
IDs and moves the workspace to schema 3, then commits and pushes it. Everyone
else pulls that commit before making any change of their own. Teammates who
upgrade at the same time and each make a change before pulling will each have
written the same member files with different content, and Git reports a conflict
in `.context-circuit/ids/*.yaml`. That conflict is resolved by keeping the union
of both sides' entries in each member file, then running `check`, which reports
any number reserved twice; a number reserved twice is the person's to settle.

A member ID that is a Windows device name, which a 2.3.0 roster could hold,
stops the split: it names the member, or the author and their records, and
changes nothing. `status` and `check` report such a member beforehand. No command
renames a member, so rename it by hand in `members.yaml`, in the `created_by` of
every record they authored, and in `member.local.yaml` on each machine that
selected it.

If the shared file comes back, typically merged from a teammate's clone that had
not pulled the split, the next changing command splits it again: numbers a member
file already holds are skipped, new ones are moved by the same rules, and the
file is removed again. If Git reports a conflict in it, resolve it first,
keeping every entry. If Git reports the shared file renamed to two different
member files, keep both member files and remove the shared file.

Two clones that each split before pulling the other also give the same
authorless numbers to two different members. Nothing is handed out twice,
because allocation reads every file. `check` reports each number reserved in
more than one file, and a person removes the extra entry by hand.

Bands are the only offline collision prevention here, and they only cover members
who hold distinct ones against a current roster. Before allocating from separate
Git clones, synchronize workspace commits: a local file lock is not distributed
coordination, and a stale roster can hand two people the same band. Resolve
competing IDs and every reference before sharing either conflicting record;
preserve the reservation history.

An explicitly invoked `check` reports local bindings, broken record links,
duplicate IDs, dependency cycles, and stale worktree associations. It is diagnostic
information, not an execution gate or a substitute for code review.
