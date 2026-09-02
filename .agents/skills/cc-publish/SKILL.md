---
name: cc-publish
description: Publish Context Circuit data to an external system (ClickUp, Jira, GitHub, Notion, Slack, …) as self-contained, one-way, idempotent records. A manually-triggered publication, orthogonal to the core workflow — never an intent/plan/execute/verify/deliver step, and it never changes plan status. It publishes either a plan and its tasks (kind plan), or a plan's open questions as a chat discussion thread (kind thread).
---

## When to use

Only on an explicit human request to reflect Context Circuit data outward — "publish
plan 0023 to ClickUp", "push this plan to our Jira", "mirror 0023 to the team's
Notion", "open a discussion thread for 0023's open questions in Slack". A configured
pipeline is a **publication** (INV-EXTERNAL-01/02/03); this skill publishes it
according to its `kind`: **`plan`** (a plan and its tasks → a tracker) or **`thread`**
(a plan's open questions → a chat discussion).

Publishing is **orthogonal to the core workflow**. It is not a phase and is not
triggered by one: approving an intent, executing, verifying, delivering, or
completing a plan neither runs this nor is affected by it. You may publish a plan in
any state — `draft`, `done`, or never. Run it only when the human asks, every time;
there is no automatic, scheduled, or event trigger.

Note the vocabulary: in this product "publish" means sending data to an external
system. Git delivery is "push" / "open a pull request" (`cc-deliver`) and never
"publish"; the two never share a word.

## What it is allowed to do

- **Read** only the workspace artifacts the publication's `config.yaml` declares
  (`reads:` — `plans`/`tasks` for the `plan` kind; a plan's open questions for the
  `thread` kind).
- **Create or update** work items in the external system, through the **host or MCP
  provider tools** — never through the runtime engine, which has no provider action
  (INV-RUNTIME-01).
- **Write** its own files under its own `publication/<name>/` folder: the
  `field-intent/<plan-id>.yaml` field intent (below) and the authoritative record under
  `publication/<name>/published/<plan-id>.yaml`.
- **Read the provider's current field values — display only, and only when the
  config opts in** (`preview.drift_read: true`). This is used solely to show a
  "tracker now" column in a preview diff; nothing read is ever persisted to any
  workspace file, and a publish still pushes the desired value (INV-EXTERNAL-02).

It must **never** modify `plan.yaml`, task files, plan status, or any other core
Context Circuit state, and it never delivers, merges, pushes, or completes anything.
Nothing is written under `plans/`; the publication writes only within its own
`publication/<name>/` folder.

## Setup (first use)

If no publication is configured for the request, create one — the `publication/`
folder is created on first use, never shipped empty:

1. Write `publication/<name>/config.yaml` per
   `wrapper/contracts/schemas/publication-config.yaml`: a `<name>` like
   `plans-clickup`, `kind: plan`, `direction: export`, `trigger: manual`, the
   `provider`, a bounded `reads:` list, and a **credential-free** `target_ref` (a
   list name, project key, repo id — never a token). Optionally set
   `preview.drift_read: true` to let the preview show a display-only "tracker now"
   column (off by default).
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

## Authoring: instructions

All external text you write honors the publication's `instructions` (optional free
text), for every kind:

- **The language to author in** — e.g. "write in Bahasa Indonesia"; default is the
  source plan's language. Author every title, description, checklist item, and
  message in that language; never translate ids or `target_ref`, and never change any
  workspace text (INV-PLAN-01).
- **Tone, phrasing, and term handling** — e.g. "everyday conversational tone; keep
  technical terms in English."
- **Schedule policy** — workday hours, weekends skipped, a start date. This is
  *policy* and stays in `instructions:`. The concrete resolved per-plan values it
  implies (this plan starts here, is estimated at that) do **not** stay in the prose;
  they live in the structured `field-intent/` layer below.

Instructions guide wording and *policy* only: they never make you leak internals
(INV-EXTERNAL-03), read beyond `reads`, write back to the workspace, change the
mapping, or relax a boundary.

## Field intent and estimates (the `field-intent/` layer)

Publishable provider fields that are **not derivable from the plan itself** — the
schedule dates and the time estimate, pushed one-way as **best-effort estimates** —
live in a user-owned, structured, diffable file per plan, per
`wrapper/contracts/schemas/publication-field-intent.yaml`:

```
publication/<name>/field-intent/<plan-id>.yaml
```

- **Derived once, then owned.** On first publish or preview with no
  `field-intent/<plan-id>.yaml`, derive a **first draft** from the `instructions:` schedule
  policy (workday hours, weekends skipped, start date) and the plan set, and write it
  here — the only workspace write the preview makes, and only under the publication's
  own `field-intent/`, never under `plans/`. From then on this file is authoritative for
  intent; the derivation **never silently overwrites a human edit**.
- **What it holds, and never restates.** Only the field intent — plan-level and
  optional per-task `start_date`, `due_date`, `estimate_minutes`. Everything the plan
  already owns (title, tasks, `depends_on`, status) stays read from the plan as-found;
  `field-intent/` never restates it, holds no external id or url (those live in the record),
  and holds no credential (INV-SEC-01).

**The estimate unit.** Canonical on disk is `estimate_minutes`, an **integer** —
exact (no `2.5h` float), diff-stable (one representation per value), portable to every
provider. The friendly `"2h 30m"` / `"15m"` form is **input and display only**:

- **`h` and `m` units only.** `d`/`w` are rejected — a "day" is scheduling policy
  (its hours are config-relative) and belongs in `instructions:`, not a duration.
- **One canonical spelling when rendered:** largest unit first, no zero units, single
  space — `2h 30m`, `15m`, `1h`. `90m` typed in renders back as `1h 30m`.
- **Parse on input, store minutes, format on output.** A malformed input (`"2 hrs"`,
  bare `"90"`, `"1d"`) is rejected at the I/O boundary and never stored — because the
  stored form is always minutes, a bad edit cannot become bad state.

**Provider conversion from `estimate_minutes`:** ClickUp `time_estimate` (ms) =
`minutes × 60000`; Jira `timetracking.originalEstimate` = the `"2h 30m"` string
directly; GitHub Projects / Notion number field = `minutes ÷ 60` (hours) or minutes
as-is per the field.

## Preview and consult (the `plan` kind)

Changing an external field is **consult-first**: the default path is a preview, not a
blind push. `cc-publish` has a **preview (dry-run) mode** — it does everything the
publish path does *except* the provider writes. It is a **mode of this skill**, not a
new skill, route, or authority (INV-SKILL-01), and it is manual like every publication
action; nothing auto-advances from preview to publish (INV-EXTERNAL-01).

1. **Load the three local layers** — the plan as-found (`plans/<plan>/`), the intent
   (`field-intent/<plan-id>.yaml`), and the last snapshot (`published/<plan-id>.yaml`
   `fields:`). If no `field-intent/` file exists, derive a first draft and write it (the
   only write the preview makes, only under the publication's own `field-intent/`).
2. **Render a plain-language diff.** For each plan and task, show intended vs
   last-pushed values in the human format, marking added / changed / unchanged:

   ```
   [0021] Context references
     due date   2026-09-14 → 2026-09-15   (changed)
     estimate   2h → 3h                    (changed)
     001 · …    estimate 1h                (new — never published)
     002 · …    estimate 45m               (unchanged)
   ```

   The default diff is **intent vs last-published snapshot** — entirely local data, so
   the common "let me review the dates before I re-publish" case touches the provider
   **zero times**.
3. **Converse and edit.** Discuss; edits land in `field-intent/<plan-id>.yaml`; re-preview
   until it reads right. No provider call has happened yet.
4. **Publish on an explicit go.** Only "publish" pushes through the host/MCP tools and
   then refreshes the `fields:` snapshot in the record.

**Optional display-only drift read.** The snapshot answers "what did *we* send," not
"did a human change it in the tracker since." When the config opts in
(`preview.drift_read: true`), the preview may read the provider's *current* field
values and add a third column — `desired … | we sent … | tracker now …`. This read is
strictly bounded: **display only** (nothing read is written to `field-intent/`, the record,
`plan.yaml`, or any workspace file), **never authority** (a publish still pushes the
desired value; tracker drift stays cosmetic), read-only through the same
host/MCP tools, and **off by default**. If the config does not opt in, never read the
provider in preview.

## Publish (the `plan` kind)

Read the plan and its tasks, the intent (`field-intent/<plan-id>.yaml`), and the existing
record if one exists. Then realize the fixed mapping through the host/MCP tools, and
push the intent field values (dates, estimate — converted from `estimate_minutes` per
the provider table above):

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

On re-publish, compare against the record on **both** axes — content and fields:

- **Content:** each task's current digest vs the recorded `synced_digest`.
- **Fields:** the intent values vs the recorded `fields:` snapshot.

Update the **same** external item in place, and skip a task only when it is unchanged
on **both** — so a task whose content is unchanged but whose desired estimate moved
from 60→90 is still a **field-only update**, not a skip. This closes the prior blind
spot where a date/estimate change was invisible because only content was hashed. Never
create a second item for a plan or task that already has a mapped id. Leave provider
fields neither the intent nor the record owns (a human's manual edits) untouched.

## Publish (the `thread` kind)

For a `kind: thread` publication, read the plan's **open questions** (from its
readable plan/task sections) and open a discussion:

- **Parent message** — headline `[thread] [<plan-number>] <plan title> — open
  questions`, then one line of context. `[thread]` marks the kind; `[<plan-number>]`
  is the only cross-reference.
- **One reply per question, fully described** — each open question is a single
  threaded reply that states the question *and* the context needed to answer it, so a
  reader can engage with just that message. One question, one message.
- **Self-contained** — same INV-EXTERNAL-03 rule: no workspace file, path, id, or
  internal mechanism; the plan id appears only in the parent headline.
- **Discussion-safe idempotency** — the thread is for asking, not mirroring state. On
  re-run, edit only *your own* question messages in place and append a reply for a
  newly-added question; **never** edit or delete a human's reply, and **never delete a
  message** (a no-longer-open question is left in place, optionally with a short
  "resolved" reply).

Write the record per `wrapper/contracts/schemas/publication-thread-record.yaml` —
the thread `parent_ts` and each question's `reply_ts`, keyed by a stable per-question
`ref`.

## Record and report

After the provider calls return the external ids, write the record to
`publication/<name>/published/<plan-id>.yaml` per the kind's record schema
(`publication-record.yaml` for `plan`, `publication-thread-record.yaml` for
`thread`). For the `plan` kind, refresh the per-item **`fields:` snapshot** with the
values just pushed, so the next preview diffs locally. The snapshot records *what we
sent*, not what the provider now holds. Everything about the publication lives under
`publication/`; nothing is written under `plans/`. If a cross-plan or reverse view is
asked for, produce it on demand by scanning the publication's `published/` records.

Report in plain project language: which plan was reflected where, how many items were
created versus updated versus skipped — **noting field-only updates** (e.g. "0021: 2
tasks re-dated, content unchanged") so a schedule change reads clearly — and the link
to the external parent. A **preview** run instead reports the diff and states plainly
that **nothing was pushed**. Never expose internal mechanism or runtime paths.

## Boundaries

Manual invocation only; export only; one-way and non-authoritative. Preview and
publish are distinct manual actions; preview never auto-advances to a push. Never
mutate `plan.yaml` or plan status, never trigger or be triggered by a workflow phase,
never put a credential or provider payload in a workspace file, and never invoke a
runtime provider action (there is none), and never write under `plans/`. The
publication writes only within its own `publication/<name>/` folder (config,
`field-intent/`, `published/`).
The only permitted provider read is the **display-only** drift read in preview when
the config opts in: nothing read is ever persisted to a workspace file, and it never
becomes intent or authority. If the host cannot reach the provider, report
`host-blocked` and change nothing. Inbound flow (external → plan) is out of scope; a
plan is only ever authored through the normal planning gate.
