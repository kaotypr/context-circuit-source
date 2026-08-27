# External-surface — the `publish-plan` kind

Continues from [design.md](./design.md); the config and record shapes it uses are
in [configuration-and-records.md](./configuration-and-records.md). This file
specifies the **first target kind**: reflecting a Context Circuit plan and its tasks
onto an external task tracker (ClickUp, Jira, or GitHub).

## What it does

On `/publish-plan <plan-id>` (against a configured `kind: publish-plan` target),
the adapter reads the plan and its tasks and creates — or, on re-run, updates — the
matching work items in the external tracker (ClickUp, Jira, GitHub, Notion, … — the
provider is a config field, not part of the kind). It is **one-way,
non-authoritative, and idempotent**: `plan.yaml` stays the source of truth, and
re-running reflects the current state onto the *same* external items rather than
duplicating them (via the `synced_digest` mapping in the plan directory).

## Granularity: plan → work-item, task → child-item

Context Circuit has exactly **two structural levels** — a plan, and a flat list of
tasks under it — which lines up 1:1 with the "work item over child item" level every
tracker has:

| Context Circuit | ClickUp | Jira | GitHub | Notion |
| --- | --- | --- | --- | --- |
| **Plan** | Task | Story | Issue | database item (page) |
| **Task** | Subtask | Sub-task | sub-issue | sub-item |
| Task `acceptance` + `verification` | checklist | description checklist | task-list | to-do checklist in the page |

The provider column is open-ended: a tracker is a valid `publish-plan` provider
whenever it offers a first-class work item with a child level and a status — which
is why Notion (a database of items with sub-items and a Status property) qualifies
just as ClickUp/Jira/GitHub do, and why the list is illustrative, not closed.

The plan maps to a **first-class work item one level up**, not to a container (a
ClickUp List / Jira Epic). Making the plan a real work item is why the mapping
generalizes cleanly across providers; the *type label* (Story / Epic / Bug) is a
human concern applied afterward, not something the adapter sets — see
[Type is left to the human](#type-is-left-to-the-human).

## Containment vs dependency — do not conflate them

Trackers carry **two** different relationships; Context Circuit uses each for a
different thing, and they must not be mixed:

- **containment** (subtask-of / sub-issue-of) ← the **plan → task** edge, the one
  structural level.
- **dependency** (blocks / is-blocked-by) ← a task's **`depends_on`**, expressed as
  a native dependency *link between sibling task-items*.

Task 002 depending on 001 does **not** make 002 a child of 001; they are siblings
under the plan, joined by a "blocked by" link. This matters most on providers whose
nesting is cheap (ClickUp subtasks, GitHub sub-issues): the deeper nesting they
offer is capacity left unfilled, never a place to encode `depends_on`. Inter-plan
`plan_dependencies` likewise map to work-item-level dependency links, never nesting.

## Acceptance and verification → checklists

Because tasks occupy the *child-item* level, their `acceptance` / `verification`
ids become a **checklist on the child item**, not a further nesting level (which
Jira forbids on sub-tasks and ClickUp nests poorly). This keeps the mapping within
every provider's reliable structure.

## Type is left to the human

The adapter creates plain structural items (work item + child item) and does **not**
set a task type. Reason: custom task types are plan-gated (ClickUp Business; Jira
issue types are workspace-defined), so setting one via API breaks on workspaces that
lack it. Creating the universal primitive and letting a human re-label the type
afterward is portable. Note the asymmetry the adapter must not paper over: a ClickUp
subtask re-types freely, but a **Jira sub-task cannot be re-typed into a Story with a
dropdown** (it is a distinct issue-type category) — so on Jira the "adjust the type
later" affordance applies to the plan item, not the children. The adapter promises
only what its provider can deliver.

## One-way and non-authoritative

Status flows **out** of `plan.yaml` on publish and never back. The external item has
its own status field; the adapter sets it from `plan.yaml` at publish time, and a
human editing it in the tracker is cosmetic drift that never returns to the
workspace (INV-PLAN-01). This is why the kind is `publish-plan`, not `sync-plan`:
"sync" would promise a bidirectional relationship the design explicitly refuses.

## Provider wrinkles the adapter resolves

The provider-neutral shape above is identical everywhere; each provider adapter
absorbs its own constraints, which is the whole reason the kind is provider-neutral
and the provider is a config field:

- **GitHub — issues are repo-scoped.** A Context Circuit task may span multiple
  repositories, but a GitHub issue lives in exactly one repo. The GitHub adapter
  resolves this by publishing into a designated **tracker-home repo** (or org-level
  GitHub Projects); ClickUp and Jira never face this because a Task / Story is not
  repo-bound.
- **Jira — sub-tasks are rigid.** One nesting level only, and a fixed issue-type
  category (see above).
- **GitHub — no first-class "blocks."** `depends_on` maps to an issue reference or a
  Projects relationship field rather than a native dependency link.
- **Notion — schema-flexible, so it needs more config.** There are no fixed issue
  types or a fixed Status field; everything is database properties. `target_ref`
  must name the database and which properties carry status (and, if used, type), so
  a Notion target's `external-target.yaml` is richer than a ClickUp one. Its upside:
  containment (sub-item parent relation) and dependency (a "Blocked by" relation
  property) are two *distinct* relations, so the containment-vs-dependency rule is
  expressed natively with no risk of conflation.

## The "publish" wording guardrail

`publish-plan` reuses the word "publish," which the v0.5 delivery boundary already
uses for **git** publication (push the branch, open the PR — a core-workflow action,
INV-DELIVER-01 / INV-RUNTIME-01). The two are disambiguated by **object**: this kind
publishes a *plan* to a tracker; delivery publishes a *branch / pull request* in git
and never "publishes a plan." To keep that clean, every core-delivery mention stays
qualified as *git* publication, so bare "publish" + "plan" unambiguously means this
kind. This guardrail is restated as a contract note in [contracts.md](./contracts.md).

## Worked trace

1. **Configure (first use).** User: *"set up publishing plans to our ClickUp
   Engineering list."* The adapter creates `external-targets/team-clickup/` with an
   `external-target.yaml` (`kind: publish-plan`, `provider: clickup`,
   `target_ref.list: Engineering`); the ClickUp token stays in the host / MCP layer.
2. **Publish.** User: *"publish plan 0023 to ClickUp."* The adapter reads
   `plans/0023-.../` (declared `reads: [plans, tasks]`), creates one Task for the
   plan and one Subtask per task with acceptance/verification as checklists, links
   `depends_on` as ClickUp dependencies, sets each status from `plan.yaml`, and
   writes the mapping to `plans/0023-.../external/team-clickup.yaml`.
3. **Re-publish after edits.** User changes two tasks, then: *"publish 0023
   again."* The adapter compares each task's digest to `synced_digest`, updates the
   two changed Subtasks in place, skips the rest, and leaves a human's manual ClickUp
   edits to non-mapped fields untouched. No duplicates.
4. **Core workflow, meanwhile, is untouched.** Approving, executing, verifying, or
   delivering 0023 neither triggers nor is affected by any of the above; the plan
   could equally have been published while still `draft`, or never.
