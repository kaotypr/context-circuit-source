# External-surface — the `plan` kind

Continues from [design.md](./design.md); the config and record shapes it uses are
in [configuration-and-records.md](./configuration-and-records.md). This file
specifies the **first publication kind**: reflecting a Context Circuit plan and its
tasks onto an external task tracker (ClickUp, Jira, GitHub, Notion, …).

## What it does

On `/cc-publish <plan-id>` (the `cc-publish` skill, publishing the `plan` kind
against a configured publication such as `plans-clickup`), the skill reads the plan
and its tasks and creates — or, on re-run, updates — the matching work items in the
external tracker. It is **one-way, non-authoritative,
self-contained, and idempotent**: `plan.yaml` stays the source of truth, published
text reads as ordinary project work, and re-running reflects the current state onto
the *same* external items rather than duplicating them.

## Preflight the provider

Before mapping, discover the target's real capabilities: read its configured
**statuses** (workspace/list/workflow-defined — a ClickUp list's statuses, a Jira
workflow's states) and map the plan's status through the config's `status:` map;
and check which structures the adapter can create. Where the adapter reaches the
provider through a surface with no checklist-create call, realize
`acceptance`/`verification` as a description checklist (`acceptance_as:
description-checklist`).

## Granularity: plan → work-item, task → child-item

Context Circuit has exactly **two structural levels** — a plan, and a flat list of
tasks under it — which lines up 1:1 with the "work item over child item" level every
tracker has:

| Context Circuit | ClickUp | Jira | GitHub | Notion |
| --- | --- | --- | --- | --- |
| **Plan** | Task | Story | Issue | database item (page) |
| **Task** | Subtask | Sub-task | sub-issue | sub-item |
| Task `acceptance` + `verification` | checklist | description checklist | task-list | to-do checklist in the page |

The provider column is open-ended: a tracker is a valid `plan`-kind provider
whenever it offers a first-class work item with a child level and a status. The plan
maps to a first-class work item, not a container (a ClickUp List / Jira Epic) — that
is what makes the mapping generalize.

## Containment vs dependency — do not conflate them

Trackers carry **two** different relationships; use each for a different thing:

- **containment** (subtask-of / sub-issue-of) ← the **plan → task** edge.
- **dependency** (blocks / is-blocked-by) ← a task's **`depends_on`**, expressed as
  a native dependency *link between sibling task-items*.

Task 002 depending on 001 does not make 002 a child of 001; they are siblings under
the plan, joined by a "blocked by" link. Never encode `depends_on` as nesting, even
where the provider makes nesting cheap. Inter-plan `plan_dependencies` likewise map
to work-item-level dependency links, never nesting.

## Type is left to the human

The adapter creates plain structural items and does **not** set a task type. Custom
task types are plan-gated (ClickUp Business; Jira issue types are workspace-defined),
so setting one via API breaks where it does not exist. A human relabels the type
afterward. Note the asymmetry: a ClickUp subtask re-types freely, but a Jira
sub-task cannot be re-typed into a Story by dropdown — so on Jira the "adjust the
type later" affordance applies to the plan item, not the children.

## External text is self-contained (INV-EXTERNAL-03)

Everything published to the tracker must read as ordinary project work to someone
who has never heard of Context Circuit. The parent's description is the plan's
**objective and context** in plain prose; each child's is the task's **intended
outcome** plus acceptance/verification as plain criteria. External fields never
carry a workspace file name, a workspace-internal path, a workspace-internal id, or
internal mechanism/vocabulary.

Per field — internal framing on the left, the self-contained form to publish on the
right:

| Instead of | Write |
| --- | --- |
| a workspace file as the reference ("`plan.yaml` is the source of truth") | omit it — the reader does not need the mechanism |
| an internal id ("AC-001a — …") | the criterion alone ("The configuration carries no secret") |
| a mechanism note ("reflected by the publish adapter") | omit it |
| the id slug in the title | the human title, prefixed with `[<plan-number>]` for mapping |

The adapter never rewrites the *meaning* of a plan's acceptance text; it strips the
internal id and any workspace path and publishes the human-meaningful statement.

### Title convention

The one allowed cross-reference is a stable plan (or task) id, normally in the
title, so a human can map what went where:

- **Parent** — `[<plan-number>] <plan title>`, where `<plan-number>` is the
  four-digit prefix of the plan id and `<plan title>` is the plan's human `title:`
  field. So `[0100] Add CSV export to the reports page`, **not**
  `[0100] csv-export-reports`. The id slug appears nowhere in external text.
- **Child** — `<task-number> · <task title>` (e.g. `001 · …`).

## One-way and non-authoritative

Status flows **out** of `plan.yaml` on publish — mapped through the config's
`status:` map onto the target's discovered status set — and never back. The external
item has its own status field, set at publish time; a human editing it in the
tracker is cosmetic drift that never returns to the workspace (INV-PLAN-01). The
surface *publishes*; it does not *sync* — a two-way sync would promise a
bidirectional relationship the design refuses.

## Idempotent re-run

On re-publish, compare each task's current digest to the `synced_digest` recorded in
the publication's record: update the **same** external item in place, and skip
unchanged tasks. Never create a second item for a plan or task that already has a
mapped id. Leave provider fields the record does not own (a human's manual edits)
untouched.

## Provider wrinkles the adapter resolves

The provider-neutral shape above is identical everywhere; each provider adapter
absorbs its own constraints:

- **GitHub — issues are repo-scoped.** A Context Circuit task may span multiple
  repositories, but a GitHub issue lives in one repo. The GitHub adapter publishes
  into a designated tracker-home repo (or org-level GitHub Projects); ClickUp and
  Jira never face this because a Task / Story is not repo-bound.
- **Jira — sub-tasks are rigid.** One nesting level only, and a fixed issue-type
  category.
- **GitHub — no first-class "blocks."** `depends_on` maps to an issue reference or a
  Projects relationship field rather than a native dependency link.
- **Notion — schema-flexible, so it needs more config.** Everything is database
  properties; `target_ref` must name the database and which properties carry status
  (and, if used, type). Its upside: containment (sub-item relation) and dependency
  (a "Blocked by" relation property) are two distinct relations, so the
  containment-vs-dependency rule is expressed natively.
- **ClickUp — checklists depend on the adapter, not the platform.** ClickUp has
  native checklists, but an adapter may reach it through a surface with no
  checklist-create call, in which case `acceptance`/`verification` realize as a
  description checklist (`acceptance_as: description-checklist`).
- **Status sets are provider-defined — discover them first.** Statuses are
  workspace/list/workflow-defined, so the adapter reads the target's configured
  statuses and maps the plan's draft/approved/done through the config's `status:`
  map. There is no universal status vocabulary to assume.

## Worked trace

1. **Configure (first use).** User: *"set up publishing plans to our ClickUp
   Engineering list."* The adapter creates `publication/plans-clickup/` with a
   `config.yaml` (`kind: plan`, `provider: clickup`,
   `target_ref.list: Engineering`); the ClickUp token stays in the host / MCP layer.
2. **Publish.** User: *"publish plan 0023 to ClickUp."* The adapter reads
   `plans/0023-.../` (declared `reads: [plans, tasks]`), creates one Task for the
   plan (`[0023] …`) and one Subtask per task with self-contained descriptions and
   acceptance/verification as checklists, links `depends_on` as ClickUp
   dependencies, sets each status from `plan.yaml`, and writes the record to
   `publication/plans-clickup/published/0023-....yaml`.
3. **Re-publish after edits.** User changes two tasks, then: *"publish 0023
   again."* The adapter compares each task's digest to `synced_digest`, updates the
   two changed Subtasks in place, skips the rest, and leaves manual tracker edits to
   non-mapped fields untouched. No duplicates.
4. **Core workflow, meanwhile, is untouched.** Approving, executing, verifying, or
   delivering 0023 neither triggers nor is affected by any of the above; the plan
   could equally have been published while still `draft`, or never.
