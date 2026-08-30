# Intent-field layer, the estimate unit, and the record snapshot

Continues [design.md](./design.md). This file is the implementation-depth design
for where publishable field values live and how they are stored.

## The `intent/` layer

A new user-owned file per plan under a publication:

```
publication/<name>/intent/<plan-id>.yaml
```

It holds the **field intent** a publication contributes that is *not* derivable
from the plan itself — the schedule dates and the time estimate. Everything the
plan already owns (title, task list, `depends_on`, status) stays read from the plan
as-found; `intent/` never restates it.

```yaml
schema_version: 1
plan: 0021-context-references
fields:                       # plan-level provider fields (best-effort, one-way)
  start_date: 2026-09-14      # ISO 8601 date (or date-time); canonical on disk
  due_date: 2026-09-14
  estimate_minutes: 120       # canonical integer minutes
tasks:
  - task: CTXREF-001          # plan.yaml task id
    estimate_minutes: 60
    # start_date/due_date optional per task; omit to inherit nothing
```

Properties:

- **User-owned, structured, diffable.** Unlike the prose schedule in
  `instructions:`, this is data the consult session edits and the preview diffs.
- **Derived once, then owned.** On first publish with no `intent/<plan-id>.yaml`,
  `cc-publish` derives a **first draft** from the `instructions:` schedule policy
  (workday hours, weekends skipped, start date) and the plan set — the same
  derivation v0.6 does inline today — and writes it here as a starting point. From
  then on the file is authoritative for intent; the derivation never silently
  overwrites a human edit.
- **Credential-free** (INV-SEC-01) and **self-contained in spirit** — it holds only
  values, no provider payload, no token, no external id (external ids live in the
  record, never here).
- **Not a core-workflow input.** It lives under `publication/`; no plan/approve/
  execute/verify/deliver step reads or writes it (INV-EXTERNAL-01/02).

`config.yaml` `instructions:` keeps exactly its v0.6 job — language, tone, term
handling, and the **schedule policy** — but the concrete per-plan resolved values
move out of the prose and into `intent/`. The prose sentence "8h workdays, skip
weekends, start 2026-08-17" is policy and stays; "0021 2h, 0022 3h on 2026-09-14"
is resolved intent and becomes `intent/` data.

## The estimate unit

**Canonical on disk: `estimate_minutes`, an integer.**

- **Exact** — minutes is fine enough for any real estimate and integers avoid the
  `2.5h` float-rounding trap.
- **Diff-stable** — exactly one representation per value, so the preview diff lights
  up only on a real change, never on `"2h 30m"` vs `"150m"` vs `"2.5h"` noise.
- **Portable** — every provider conversion starts from minutes (table below).

**Human I/O and display: the `"2h 30m"` / `"15m"` format.** Parsed to minutes on
input (in the consult session and any hand edit), formatted from minutes for the
diff and the report. Rules:

- Units are **`h` and `m` only**. `d`/`w` are rejected — a "day" is scheduling
  policy (its hours are config-relative), and belongs in `instructions:`, not in a
  duration.
- One canonical spelling when rendered: largest unit first, no zero units, single
  space — `2h 30m`, `15m`, `1h`. `90m` on input renders back as `1h 30m`.
- A malformed input (`"2 hrs"`, bare `"90"`, `"1d"`) is **rejected at the I/O
  boundary**, never stored — because the stored form is always minutes, a bad edit
  cannot become bad state.

### Provider conversion (from `estimate_minutes`)

| Provider | Native estimate field | Conversion |
| --- | --- | --- |
| ClickUp | `time_estimate` (milliseconds) | `minutes × 60000` |
| Jira | `timetracking.originalEstimate` (string) | format as `"2h 30m"` — Jira consumes this form directly |
| GitHub Projects | numeric custom field | `minutes ÷ 60` (hours), or minutes as-is per the field |
| Notion | number property | `minutes ÷ 60`, or as-is per the property |

The `"2h 30m"` format doubles as the Jira interchange value — another reason it is a
real layer, not cosmetics — while ClickUp's millisecond field is reached by the
pure-integer conversion.

## The record snapshot

`publication-record.yaml` gains an optional per-item **`fields:`** block recording
the values **last actually pushed**:

```yaml
schema_version: 1
plan: 0021-context-references
provider: clickup
plan_item:
  id: "86eyt0wzf"
  url: "https://app.clickup.com/t/86eyt0wzf"
  fields:                       # NEW — last-published snapshot
    start_date: 2026-09-14
    due_date: 2026-09-14
    estimate_minutes: 120
task_items:
  - task: CTXREF-001
    id: "86eyt0x0p"
    url: "https://app.clickup.com/t/86eyt0x0p"
    synced_digest: "sha256:…"   # unchanged — digest of task content
    fields:                     # NEW
      estimate_minutes: 60
```

- **Purpose:** answer "what did we send?" and "what changed since?" from **local
  data**, with no provider read. The preview diffs `intent/` against this snapshot
  by default.
- **Idempotency covers fields now.** A re-run compares both the content
  `synced_digest` *and* the `fields:` snapshot: a task whose content is unchanged
  but whose desired estimate moved from 60→90 is updated (field-only update); a task
  unchanged on both is skipped. This removes the v0.6 blind spot where a
  date/estimate change was invisible to idempotency because only content was hashed.
- **Still non-authoritative and credential-free.** The snapshot is *what we sent*,
  not *what the provider now holds* — a human edit in the tracker is drift and never
  updates this snapshot except through the next publish. No token, payload, or
  workspace path (INV-EXTERNAL-02, INV-SEC-01).
- **Additive and backward-compatible.** `fields:` is optional; a v0.6 record with no
  `fields:` block is valid and simply carries no snapshot until its next publish.

## What each file owns (single source per fact)

| Fact | Owner |
| --- | --- |
| Plan title, tasks, `depends_on`, status | `plans/<plan>/` (read-only source) |
| Language, tone, schedule *policy* | `config.yaml` `instructions:` |
| Resolved per-plan dates and estimates (**intent**) | `intent/<plan-id>.yaml` |
| External ids/urls (**identity**) | `published/<plan-id>.yaml` |
| Field values **last pushed** (**snapshot**) | `published/<plan-id>.yaml` `fields:` |
| Field values the provider **now holds** | the provider (read display-only, never stored) |

No fact is written twice, and nothing under `publication/` is authoritative over
`plans/` or reachable by a core-workflow step.
