---
name: cc-publish
description: Publish Context Circuit data to an external system (ClickUp, Jira, GitHub, Notion, …) as self-contained, one-way, idempotent records. A manually-triggered publication, orthogonal to the core workflow — never a plan/approve/execute/verify/deliver step, and it never changes plan status. The first kind is publishing a plan and its tasks.
---

## When to use

Only on an explicit human request to reflect Context Circuit data outward — "publish
plan 0023 to ClickUp", "push this plan to our Jira", "mirror 0023 to the team's
Notion". A configured pipeline is a **publication** (INV-EXTERNAL-01/02/03); this
skill publishes it according to its `kind`. The first kind is **`plan`** (a plan and
its tasks); this file specifies that kind.

Publishing is **orthogonal to the core workflow**. It is not a phase and is not
triggered by one: approving, executing, verifying, delivering, or completing a plan
neither runs this nor is affected by it. You may publish a plan in any state —
`draft`, `approved`, `done`, or never. Run it only when the human asks, every time;
there is no automatic, scheduled, or event trigger.

Note the vocabulary: in this product "publish" means sending data to an external
system. Git delivery is "push" / "open a pull request" (`cc-deliver`) and never
"publish"; the two never share a word.

## What it is allowed to do

- **Read** only the workspace artifacts the publication's `config.yaml` declares
  (`reads:`, normally `plans` and `tasks` for the `plan` kind).
- **Create or update** work items in the external system, through the **host or MCP
  provider tools** — never through the runtime engine, which has no provider action
  (INV-RUNTIME-01).
- **Write** its own record: the authoritative record under
  `publication/<name>/published/<plan-id>.yaml`.

It must **never** modify `plan.yaml`, task files, plan status, or any other core
Context Circuit state, and it never delivers, merges, pushes, or completes anything.
Nothing is written under `plans/`.

## Setup (first use)

If no publication is configured for the request, create one — the `publication/`
folder is created on first use, never shipped empty:

1. Write `publication/<name>/config.yaml` per
   `wrapper/contracts/schemas/publication-config.yaml`: a `<name>` like
   `plans-clickup`, `kind: plan`, `direction: export`, `trigger: manual`, the
   `provider`, a bounded `reads:` list, and a **credential-free** `target_ref` (a
   list name, project key, repo id — never a token).
2. Confirm the provider credential lives at the host / MCP layer, not in any
   workspace file (INV-SEC-01). If the host has no usable provider capability,
   report `host-blocked` and stop; create nothing external.

## Preflight the provider

Before mapping, discover the target's real capabilities: read its configured
**statuses** (workspace/list/workflow-defined) and map the plan's status through the
config's `status:` map; and check which structures the adapter can create. Where the
adapter reaches the provider through a surface with no checklist-create call, realize
`acceptance`/`verification` as a description checklist (`acceptance_as:
description-checklist`).

## Publish (the `plan` kind)

Read the plan and its tasks, and the existing record if one exists. Then realize the
fixed mapping through the host/MCP tools:

| Context Circuit | external work item |
| --- | --- |
| the plan | one first-class work item (ClickUp Task, Jira Story, GitHub Issue, Notion database item) |
| each task | one child item under it (Subtask / Sub-task / sub-issue / sub-item) |
| a task's `acceptance` + `verification` | a checklist on that child item |
| a task's `depends_on` | a **dependency link** between sibling child items — never nesting |

- **Containment vs dependency.** Plan→task is containment. `depends_on` is a
  "blocked by" link between siblings; never encode it as nesting, even where the
  provider makes nesting cheap.
- **Type is the human's.** Create plain structural items; do not set a task type.
- **Status is one-way.** Set each item's status from `plan.yaml` at publish — mapped
  through the config's `status:` map onto the destination's discovered status set —
  and never read it back. A human's later edit in the tracker is cosmetic drift and
  never returns to the workspace (INV-PLAN-01).

## External text is self-contained (INV-EXTERNAL-03)

Everything a reader sees in the tracker must make sense to someone who has never
heard of Context Circuit. Write the *work*, never the plumbing.

- **Parent (the plan):** title `[<plan-number>] <plan title>`, where `<plan-number>`
  is the four-digit prefix of the plan id and `<plan title>` is the plan's human
  `title:` field — so `[0100] Add CSV export to the reports page`, **not**
  `[0100] csv-export-reports`. The id slug appears nowhere in external text.
  Description is the plan's objective and relevant context in plain prose.
- **Child (a task):** title `<task-number> · <task title>` (e.g. `001 · …`).
  Description is the task's intended behavior in plain prose, followed by
  acceptance/verification as **plain human criteria** — the statement only, with the
  internal `AC-…`/`VR-…` ids stripped.
- **Never put in any external field** a workspace file name (`plan.yaml`,
  `config.yaml`), a workspace-internal path (`plans/`, `context/`, `publication/`),
  a workspace-internal id (acceptance/verification ids), or internal
  mechanism/vocabulary ("mirror", "source of truth", "Context Circuit plan/task").
  The only allowed cross-reference is a stable plan or task id, normally in the
  title.

If a plan's own acceptance text happens to name an internal path, publish the
human-meaningful statement and drop the path; you are describing the outcome, not
the workspace.

## Idempotent re-run

On re-publish, compare each task's current digest to the `synced_digest` recorded in
the publication's record: update the **same** external item in place, and skip
unchanged tasks. Never create a second item for a plan or task that already has a
mapped id. Leave provider fields the record does not own (a human's manual edits)
untouched.

## Record and report

After the provider calls return the external ids, write the record to
`publication/<name>/published/<plan-id>.yaml` per
`wrapper/contracts/schemas/publication-record.yaml`. Everything about the
publication lives under `publication/`; nothing is written under `plans/`. If a
cross-plan or reverse view is asked for, produce it on demand by scanning the
publication's `published/` records.

Report in plain project language: which plan was reflected where, how many items were
created versus updated versus skipped, and the link to the external parent. Never
expose internal mechanism or runtime paths.

## Boundaries

Manual invocation only; export only; one-way and non-authoritative. Never mutate
`plan.yaml` or plan status, never trigger or be triggered by a workflow phase, never
put a credential or provider payload in a workspace file, never write under `plans/`,
and never invoke a runtime provider action (there is none). If the host cannot reach
the provider, report `host-blocked` and change nothing. Inbound flow (external →
plan) is out of scope; a plan is only ever authored through the normal planning gate.
