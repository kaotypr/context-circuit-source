# External-surface — configuration and records

Continues from [design.md](./design.md). This file specifies **where a publication
keeps its configuration and its output records**. It does not restate the
orthogonality principle (that is `design.md`) or the contract additions (that is
[contracts.md](./contracts.md)).

## Everything lives under `publication/`

A publication writes two kinds of file, and both live under its own folder — no
file is ever written under `plans/`:

```text
publication/                          # product root, created on first use
  <name>/                             # one folder per publication, e.g. plans-clickup
    config.yaml                       # the user's configuration (credential-free)
    .local.yaml                       # gitignored: non-secret instance ids, if any
    published/
      <plan-id>.yaml                  # one authoritative record per published plan
```

- `<name>` is a user-chosen label, conventionally `<subject>-<provider>` —
  `plans-clickup`, `plans-github`, `docs-clickup`, `thread-slack`. It is the
  publication's identity; `config.yaml` restates `kind` and `provider`
  authoritatively.
- **Created on first use.** The template ships nothing here. The folder appears the
  first time a user configures a publication (the adapter skill creates it), so no
  unexplained folder sits in a workspace the user never opted into, and the blank
  seed stays clean — the same way `plans/` ships as a near-empty scaffold rather
  than pre-populated machinery.

Keeping the records under the publication — not in the plan directory — matches how
the two lifecycles relate. The external items live on independently of the plan:
archiving a plan does not retire its tracker tasks, so the record that tracks them
should not be hidden away with the archived plan either. Each record carries the
stable, never-reused plan id (INV-PLAN-03), so it is self-describing whether that
plan is active or archived. And because nothing external is written under `plans/`,
an agent working there during execution never encounters a publication record.

## `config.yaml` — the publication's configuration

```yaml
schema_version: 1
publication: plans-clickup     # matches the folder name
kind: plan                     # what cc-publish publishes (plan | docs | thread | …)
provider: clickup              # clickup | jira | github | notion | slack | … (open-ended)
direction: export              # export-first; the only value in this scope
trigger: manual                # manual is the only allowed value
reads:                         # DECLARED, bounded artifact scope — no scanning
  - plans
  - tasks
target_ref:                    # credential-free logical destination
  list: "Engineering"
mapping:                       # kind-specific field map (see publish-plan.md)
  plan: task
  task: subtask
  acceptance_as: checklist     # or description-checklist when the adapter has no checklist-create call
  status:                      # plan status -> the target's discovered status set
    draft: todo
    approved: todo
    done: completed
```

- **`kind` selects the publish procedure.** The `cc-publish` skill publishes a
  publication according to its `kind` (`plan` here); `provider` selects the host/MCP
  adapter. A future `docs`/`thread` kind is new behavior in `cc-publish` plus a
  `config.yaml`, nothing in the core.
- **`reads` is a bounded, declared allow-list.** A publication reads only the
  artifact types it names, echoing the v0.5 discipline that the workspace never
  scans to fill a gap (INV-SEC-02 ethos). A `plan` publication reads `plans`/`tasks`;
  a `thread` publication might read open questions only.
- **`direction: export` and `trigger: manual`** are the only values this scope
  allows; they exist as fields so the contract can refuse anything else rather than
  silently doing something coupling or automatic.

## Credentials never enter the workspace

`config.yaml` is credential-free. Provider tokens live in the host Git configuration
or the host agent / MCP server and never enter any workspace file or record
(INV-SEC-01). `target_ref` holds only a *logical* destination (a list name, a
project key, a repo id, a channel) — never a secret, never a session token. Anything
genuinely instance-specific but non-secret (a numeric account id) may go in the
gitignored `.local.yaml`, mirroring how machine-specific binding lives in
`repositories.local.yaml`.

## The authoritative record (under the publication)

```yaml
# publication/plans-clickup/published/0023-add-oauth.yaml
plan: 0023-add-oauth
provider: clickup
plan_item: { id: "86abc", url: "https://…" }   # the plan's external work item
task_items:
  - { task: "001", id: "86def", url: "https://…", synced_digest: "sha256:…" }
  - { task: "002", id: "86ghi", url: "https://…", synced_digest: "sha256:…" }
```

- `plan.yaml` is **not** touched — the external item id is environment-specific
  projection state, not plan identity, and never lives on the plan (INV-PLAN-01).
- `synced_digest` (via the runtime's existing `cc_digest`) makes re-publish
  idempotent and lets an unchanged task be skipped: re-running compares the current
  artifact digest to the recorded one and updates the *same* external item instead
  of creating a duplicate.
- **Not an execution input.** The record lives under `publication/`, never under
  `plans/`, so no plan/approve/execute/verify/deliver step encounters it, and none
  updates external status — that would be the surface reaching into the core flow,
  which INV-EXTERNAL-01 forbids. Only an explicit manual invocation reads or writes
  it.

## Cross-plan lookup

The `published/` records are the link records; an idempotent re-publish needs only
*its own* plan's record, keyed by `synced_digest`. A cross-plan or reverse view —
"which plans are on this publication", "which plan is external item X" — is produced
on demand by scanning the publication's `published/` records, following the
principle that a derived view is a catalog, not a copy (INV-KNOWLEDGE-01).

## What writes what

The host / MCP layer makes the external API call and obtains the item ids; the
deterministic write of the record into the publication's `published/` folder reuses
the runtime's existing atomic-write path (`cc_atomic_write`). "Written by the
runtime" never implies "lives under `.runtime/`" — the path is an argument, and this
record's home is the publication's folder, whose lifecycle is independent of any
single execution or of a plan's archive state.
