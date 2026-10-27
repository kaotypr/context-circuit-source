# Intents, plans, knowledge, and review

For an intent-led change, the user approves the intended outcome, asks for it to
be planned — later, or in the same request — and later asks for a plan to be
executed. A developer may request a grounded plan directly for
an outcome they have already specified, with no intent. The agent interprets the project and writes meaningful
content between them. Go supplies IDs, record structures, references, and
reliable file operations. It does not authenticate consent or implement a plan.

A third decision is theirs as well: to skip both. A request to make a change
directly bypasses intent, planning, and worktree preparation, and the change is
made in the bound checkout on its current branch, uncommitted until they
authorize a commit. No record is written and no ID is allocated, so the diff and
the knowledge it touches are the only trace; `context find --repo ID` returns the
catalog entries claiming that repository and names the reconciliation the same
way completion does for a plan. Nothing else relaxes, and a change that turns out
to need an outcome nobody approved returns to the path below.

An intent has a global `id`, `created_by`, a `created_at` instant, linked
`plans`, and an `approved_at` instant once a person approves it. Its Markdown body contains goal, non-goals,
constraints, success criteria, rough repository scope, and open questions.
Approval records the instant and appends the person's own words under a heading
carrying it; the field records that decision and never establishes one, so
nothing reads it as consent. Renewed approval moves the field and appends another
section, keeping every instant and wording. Intent creation happens before detailed code
investigation and is grounded in existing knowledge, retrieved first and read against
what the request was understood to mean; what the knowledge contradicts or cannot
settle becomes an open question rather than an assumption.

Open questions are a numbered list — bold question, italic `_Answer:_` line beneath
once settled — so a person can answer by number. Numbers are stable: an answered
question keeps its number and is never deleted, and a later question takes the next
unused one. The unnumbered empty-state line stays when there are none. The agent
writes and answers them in the Markdown; no command records an answer, because only
approval is a gate. A question arising during planning that bears on the outcome
returns to the intent instead of being settled inside a plan.

The intent is presented and then waited on: the request that prompted it is not
approval of it, and the recorded approval is a person's answer to the written intent.

A plan's metadata looks like:

```yaml
id: p0001
created_by: maya
intent: i001 # omit for a standalone plan
repositories:
  - api
  - web
depends_on:
  - p0002
```

Dependencies are optional, and `record order` derives waves, start references,
and integration merges from them. Separate repository plans may be useful; no
one-plan-per-repository constraint applies. IDs are global and prefixed, with a minimum of
three intent digits and four plan digits; numbering expands beyond that width.
Every member holds a band and allocates from that band's block alone — intents
from `band*100` and plans from `band*1000`, so band 1 writes `i100` and `p1000`,
and band 0 writes `i001` and `p0001` — skipping every number already taken. Plans use
`created_by` only for member data; a band changes numbering, never ownership.

Approval alone stops with a recommendation to plan. The agent writes the plan from
real code once the person asks for planning of an approved intent, and links it
without asking for approval again. Record task order, useful risks, expected checks,
and optional affected paths. Paths are descriptive, not an enforcement contract.
Changes to the intended outcome reopen approval; an additional file within the
approved outcome does not automatically do so.

New plans live in `plans/pNNNN-slug/plan.md`. Add optional supporting files in
the same folder only when useful. Link each from the entry's Details section and
assign it in `required_files.shared` or `required_files.by_repository.<repo>`.
These paths are relative to the plan folder; each worker's brief lists only its
assigned files as paths relative to the stated absolute workspace root. Existing
single-file plans remain readable without migration. A direct planning request
authorizes investigation and drafting, then the agent presents the whole folder
for review. Material revisions are presented again before implementation.

Plans are presented and may be revised before execution. There is no
plan approval field or command, but implementation waits for a separate request to execute
them — the decision about when work starts, made once the complete plans are visible. That
request prepares a worktree per repository unless the user asks to work directly
in a bound checkout.

The agent runs appropriate tests/lint/builds and interprets failures. A plan
record holds the plan; the result of running it reaches the person, and
`record complete` — which a person requests — carries it into the record with
`completed_at`. Between planning and that request the record is not edited:
running a plan does not edit the plan. Resume by resolving a record ID, then
examining its current working copy and diff, because a record describes what was
intended and never proves what was implemented.

Independent review is manually requested, usually after PR creation and optionally
after delivery as an audit. An independent reviewer receives the diff, relevant
code, and success criteria; it returns findings and limitations without modifying
code or posting externally unless requested. If independence is unavailable,
describe that limit. Running implementation tests is distinct from independent
review. Findings do not trigger automatic fixes or block delivery mechanically.

Commit, push, PR, merge, deployment, and external publishing remain explicit
host tool actions. The CLI supplies repository and branch information. On an
explicit completion request, `record complete` appends a short note and records
`completed_at` in frontmatter, and returns the catalog entries scoped to that
plan's repositories. Those entries are the candidates for reconcile; the agent
reads the ones whose meaning the plan could have changed, edits each note
together with its catalog entry, and moves the reviewed date. A completion that
changed no durable concept reconciles nothing and says so in its note, which is
the ordinary case. Several plans marked complete together are reconciled once
across the set. `completed_at` is what dependency ordering reads to release a plan's
dependents, so complete a plan only when it actually landed. Delivery alone does
not mark the plan done or trigger cleanup.

Knowledge notes capture durable concepts, not task logs. Gather from named sources
on request. `context find` searches only the optional index, matching whole lines by
case-insensitive substring; the agent chooses and reads relevant notes. It can edit
knowledge and index files directly. There are no mandatory catalogs, acceptance
records, or provider integrations.

One note is one unwrapped index entry holding its link, the repositories it applies
to, the question it answers, its search terms, and the date it was last confirmed.
`context/glossary.md` is a table of project vocabulary mapped to the code identifiers
that implement it. A note anchors to repository paths written with the logical
repository ID, collected in one `Owner:` block at its end rather than scattered
through its sentences; it never names a plan record, an intent record, or a file
under `sources/`, because those are archived while the knowledge outlasts them and a
recorded evidence path would reopen material that must stay passive. `check` reports
the lines in `context/` that cross that boundary.

A record may sit directly in `intent/` or `plans/`, or in one level of group
folder below it, such as `plans/alex/p1010-thing/`; a group is ordinary file
organization and says nothing about who owns a record. `_archive/` may sit above
or below a group: `plans/_archive/`, `plans/<group>/_archive/`, and
`plans/_archive/<group>/`. Every command that reads records finds them in any of
these places. Folder names starting with `_` are kept for the CLI, and
`_archive/` is the only one it gives meaning to: a record under it is archived.
A top-level `archive/` is the legacy archive folder from before 2.5.0: it still
reads as archived and is not usable as a group name, and the first changing
command renames it to `_archive/`. An `archive/` inside a group is an ordinary
folder.

Archiving is ordinary file organization, and each member archives their own.
`record archive` moves the active member's finished records, a completed plan or
an approved intent whose plans are all completed, under the `_archive/` beside
them, keeping a group's records in that group's `_archive/`; `--id` archives
exactly the named records, whoever wrote them, and marks any unfinished one.
After `record complete`, the agent offers it, and a no changes nothing. ID-based
lookups find a record wherever it sits, archives included, and no ID file
changes, so every number stays reserved. `record list` shows the active member's
active records; `--member ID` shows another member's, `--all` everyone's, and
`--archived` adds archived ones. Allocation inventories filenames, including
archives, without reading archived content. `record archive` rewrites inline
relative links into and out of the records it moves and reports any other link
the move breaks; a record moved by hand keeps its reservation, and its
handwritten links are fixed by hand.

Deletion is not archival. `record delete --member ID` removes one member's intents
and plans, archived ones included, `record delete --all` removes every record
while the roster lists exactly one member, and `record delete --id ID`, repeatable,
removes the named records whoever wrote them, or frees a reserved ID no record
file holds. Each releases the deleted IDs from every ID file under
`.context-circuit/ids/` holding them, so allocation hands the lowest free number
in the author's band out again. A preview without `--confirm` comes first,
showing each record's author and whether it is approved or completed, and
deletes nothing; the person agrees to its list before `--confirm` runs. Records
outside the selection that point at a deleted one block the deletion rather
than being edited or removed with it. The one exception is by ID: a deleted plan
is taken out of its own intent's `plans:` list. Any other link is resolved by
naming that record too, and the refusal lists the IDs to add. Once deleted, a number that reached a
branch name, a pull request, or a conversation may later mean a different record,
so delete before an ID has travelled where possible.
