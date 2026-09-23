# Plan folder and command contract

## Folder shape

```text
plans/
  p0001-billing/
    plan.md               # record frontmatter and complete reading map
    api-contract.md       # optional detail
    data-model.md         # optional detail
    ui-structure.md       # optional detail
    assets/                # optional illustrations linked from a detail
  p0002-legacy.md         # existing v2.1 record, still readable
```

Only `plan.md` has record frontmatter. A new plan's frontmatter keeps the
current fields: `id`, `created_by`, `created_at`, nonempty `repositories`,
optional `depends_on`, and later `completed_at`. It adds `required_files`,
with separate `shared` and `by_repository` fields holding paths relative to
the plan folder. This avoids colliding with a valid repository ID such as `all`.
`intent` is present only when
the plan was derived from an approved intent. A standalone plan omits it; an
empty string, `null`, or invented sentinel ID is not the standalone form.

```yaml
---
id: p0001
created_by: maya
created_at: 2026-09-23T10:00:00Z
repositories:
  - api
  - web
depends_on: []
required_files:
  shared:
    - data-model.md
  by_repository:
    api:
      - api-contract.md
    web:
      - ui-structure.md
---
```

`plan.md` names the requested outcome in its own words, the repository scope,
the approach, ordered tasks, risks, and expected checks. It links each supporting
file from a `Details` section and identifies what question that file answers.
The reader can review the whole proposal by starting there and following those
links. Every supporting file belongs to at least one `required_files` group;
`shared` is for context every worker needs. The same file may be assigned to
several repositories. With no supporting files, omit `required_files`. Paths
may include a subdirectory such as `assets/flow.svg`, but cannot be absolute,
escape the folder with `..`, or point through a symlink. The CLI validates
that each listed path exists as a regular file and each repository key appears
in the plan's `repositories` list. Supporting files have no runtime status.
An unassigned regular file in the plan folder is a diagnostic finding: a
reviewer may see it while no worker is told to read it. Repeat paths within one
group are refused, and a file assigned through both `shared` and a repository is
deduplicated in that worker's brief.
Text details should be Markdown; Markdown can carry Mermaid, tables, JSON/YAML
examples, and field definitions. Linked illustrations may be kept under
`assets/`. Do not put secrets or machine-local data in any of them.

## Identity and lookup

The permanent ID ledger remains unchanged. A plan ID resolves to exactly one
entry: either a legacy `plans/pNNNN-slug.md` or a folder
`plans/pNNNN-slug/plan.md`, including under an archive. A duplicate across
these forms is an error. Supporting files are not records and are never
allocated IDs. `record list`, `record show`, `record dependencies`,
`record complete`, worktree operations, dispatch, and diagnostics resolve the
entry by ID through the same API; they do not infer a record from an arbitrary
Markdown file in a plan folder.

Plan creation retains `--kind plan`, `--slug`, `--title`, and one or more
`--repo`. `--intent` becomes optional for plan creation: when supplied, it must
resolve to an intent and the plan ID is appended to that intent's `plans` list;
when omitted, no intent lookup or backlink write occurs. `--intent` continues
to be irrelevant to intent creation. The CLI reserves an ID, creates the folder
and `plan.md`, and returns the entry path and plan ID. It does not generate
empty detail files or determine how many the plan needs.

`record show --id pNNNN` returns the entry text as it does today. Its
frontmatter and `Details` section expose the supporting files and their reading
assignments without concatenating all detail into a misleading single
document. `record list` returns the entry path; a new folder path may be
exposed separately for navigation. Existing JSON field names and legacy-file
results stay stable where possible.

`record order --intent iNNN` selects only plans linked to that intent. A new
repeatable `--plan pNNNN` selects exactly the named plans, with no requirement
that they share an intent. The two filters are mutually exclusive. With neither
filter, `record order` includes all active standalone and linked plans.
Dependencies outside an explicit selection are reported as satisfied when
completed, or as blocked when unfinished; selecting a plan never silently
authorizes or adds its predecessors to the run. Wave, start, integration, and
per-repository delivery marks are derived over the selected run, so a delivery
of several standalone plans uses the same selected IDs as their execution.
For one plan, `record order --plan pNNNN` gives its start and delivery facts
without pulling unrelated workspace plans into the result. The order
computation needs no synthetic intent. `record complete` appends its note and
`completed_at` to `plan.md` for a folder plan and to the existing file for a
legacy plan.

The order result names its exact selected plan IDs, even when `--intent`
produced the selection, so the coordinator can show the run's scope and reuse
it for delivery. An omitted filter is explicitly reported as an all-active
selection rather than looking like an empty intent.

## Safe contents

The CLI's existing workspace path rules apply to plan entries. Folder lookup
must reject symlinks and ambiguous ID matches. A dispatch resolves only the
regular `required_files` assigned to its repository and `shared`; it does not
follow links out to a local checkout or read credentials. It reports a missing,
unreadable, or nonregular assigned file rather than silently omitting material
the person may have reviewed. `check` reports a missing `plan.md`, a frontmatter
ID that disagrees with its folder name, duplicate IDs, broken intent backlinks,
an unknown repository key, unassigned regular files, and unsafe or missing
required paths. It also checks that each required path has a relative Markdown
link in the entry's `Details` section, so the person's reading map and the
worker's reading map do not drift apart. It does not fetch external links or
treat them as supporting files. `check` remains a diagnostic, not an execution
gate.
