# Intents, plans, knowledge, and review

The user makes two decisions: approving the intended outcome, and asking for a
plan to be executed. The agent interprets the project and writes meaningful
content between them. Go supplies IDs, record structures, references, and
reliable file operations. It does not authenticate consent or implement a plan.

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
intent: i001
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
A member holding a band allocates from that band's block alone — intents from
`band*100` and plans from `band*1000`, so band 1 writes `i100` and `p1000` — and
an unbanded member allocates from the numbers no band claims. Plans use
`created_by` only for member data; a band changes numbering, never ownership.

The agent writes the plan from real code once the intent is approved and links it,
without asking a second time. Record task order, useful risks, expected checks,
and optional affected paths. Paths are descriptive, not an enforcement contract.
Changes to the intended outcome reopen approval; an additional file within the
approved outcome does not automatically do so.

Plans are then presented and left alone. Reading them is optional, so there is no
plan approval gate, but implementation waits for a separate request to execute
them — the decision about when work starts, made once the plans are visible. That
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

Commit, push, PR, merge, deployment, and external publishing remain explicit host
tool actions. The executable supplies repository and branch information. On an
explicit completion request, `record complete` appends a short note and records
`completed_at` in frontmatter, and returns the catalog entries scoped to that
plan's repositories. Those entries are the candidates for reconcile; the agent reads
the ones whose meaning the plan could have changed, edits each note together with its
catalog entry, and moves the reviewed date. A completion that changed no durable
concept reconciles nothing and says so in its note, which is the ordinary case.
Several plans marked complete together are reconciled once across the set. That date
is what dependency ordering reads to release a plan's dependents, so complete a plan
only when it actually landed. Delivery alone does not mark the plan done or trigger
cleanup.

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

Archival is optional ordinary file organization. ID-based lookups discover records
under their intent/plans directories, including archives. Listing archives requires
`record list --archived`; normal listings read active top-level records. Allocation
inventories filenames, including archives, without reading archived content. Keep
reservations and fix any handwritten relative links when records move.
