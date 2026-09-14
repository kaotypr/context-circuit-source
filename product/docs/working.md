# Intents, plans, knowledge, and review

The user approves the intended outcome. The agent interprets the project and
writes meaningful content. Go supplies IDs, record structures, references, and
reliable file operations. It does not authenticate consent or implement a plan.

An intent has a global `id`, `created_by`, and linked `plans`. Its Markdown body
contains goal, non-goals, constraints, success criteria, rough repository scope,
and an approval note. Intent creation happens before detailed code investigation;
existing knowledge may inform it. Record actual approval, including approval
already supplied in the conversation for that exact outcome.

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

Dependencies are optional. Separate repository plans may be useful; no one-plan-
per-repository constraint applies. IDs are global and prefixed, with a minimum of
three intent digits and four plan digits; numbering expands beyond that width.
There are no member number ranges. Plans use `created_by` only for member data.

The agent writes the plan from real code after intent approval, links it, and
proceeds without a second approval. Record task order, useful risks, expected
checks, and optional affected paths. Paths are descriptive, not an enforcement
contract. Changes to the intended outcome reopen approval; an additional file
within the approved outcome does not automatically do so.

The agent runs appropriate tests/lint/builds, interprets failures, and records
observed progress in the same file. Resume by resolving a record ID, then examining
its current working copy and diff. A note is not proof of current implementation.
Use `record note` for append-only progress, or edit the Markdown directly.

Independent review is manually requested, usually after PR creation and optionally
after delivery as an audit. An independent reviewer receives the diff, relevant
code, and success criteria; it returns findings and limitations without modifying
code or posting externally unless requested. If independence is unavailable,
describe that limit. Running implementation tests is distinct from independent
review. Findings do not trigger automatic fixes or block delivery mechanically.

Commit, push, PR, merge, deployment, and external publishing remain explicit host
tool actions. The executable supplies repository and branch information. On an
explicit completion request, `record complete` appends a short note; the agent
then reconciles affected durable project knowledge and an existing context index.
Delivery alone does not mark the plan done or trigger cleanup.

Knowledge notes capture durable concepts, not task logs. Gather from named sources
on request. `context find` searches only the optional index; the agent chooses and
reads relevant notes. It can edit knowledge and index files directly. There are
no mandatory catalogs, acceptance records, or provider integrations.

Archival is optional ordinary file organization. ID-based lookups discover records
under their intent/plans directories, including archives. Listing archives requires
`record list --archived`; normal listings read active top-level records. Allocation
inventories filenames, including archives, without reading archived content. Keep
reservations and fix any handwritten relative links when records move.
