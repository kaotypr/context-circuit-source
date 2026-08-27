# External-surface — configuration and records

Continues from [design.md](./design.md). This file specifies **where the surface
keeps its configuration and its output records**, and the lifecycle rules that keep
those records correct. It does not restate the orthogonality principle (that is
`design.md`) or the contract additions (that is [contracts.md](./contracts.md)).

## Two homes, split by lifecycle

The surface writes two very different things, and they must not share a home:

| Thing | Lifecycle | Home |
| --- | --- | --- |
| **Target configuration** (which external system, how fields map) | workspace-scoped, long-lived | `external-targets/<target-id>/` |
| **Sync state + derived catalog** (cursors, last run, reverse index) | target-scoped, regenerable | `external-targets/<target-id>/state/` |
| **The authoritative plan↔external-item mapping** (the idempotency key) | **plan-scoped**, must live and die with its plan | inside the plan directory |

The reason the mapping is not under `external-targets/` is lifecycle coupling: the
mapping is the idempotency key, and it must travel with its plan through archive and
restore. Archiving a plan moves the whole `plans/<id>/` directory (INV-ARCHIVE-01);
anything inside it is carried along for free, with no new archive logic. A mapping
in a separate top-level folder would be orphaned the moment its plan is archived —
still pointing at live external items whose plan has left the active area. So the
**authoritative** mapping lives with the plan; the central folder holds only config,
sync state, and a **derived** catalog.

## `external-targets/` — the user-owned config folder

```text
external-targets/                     # created on first use, never shipped empty
  <target-id>/
    external-target.yaml              # what this target is (credential-free)
    .local.yaml                       # gitignored: account/instance ids, if any
    state/
      catalog.md                      # DERIVED reverse index (target → plans, item → plan)
      sync.yaml                       # last-run time, cursors
```

**Created on first use.** The template ships nothing here. The folder appears the
first time a user configures a target (the adapter skill creates it), so no
unexplained folder sits in a workspace the user never opted into, and the blank
seed stays clean — the same way `plans/` ships as a near-empty scaffold rather than
pre-populated machinery.

**`external-target.yaml`** — the declaration, credential-free:

```yaml
target: team-clickup          # stable local id for this target
kind: publish-plan            # names the adapter skill that runs it
direction: export             # export-first; the only value in this scope
provider: clickup             # clickup | jira | github | notion | … (open-ended)
reads:                        # DECLARED, bounded artifact scope — no scanning
  - plans
  - tasks
target_ref:                   # credential-free logical destination
  list: "Engineering"
mapping:                      # kind-specific field map (see publish-plan.md)
  plan: task
  task: subtask
  acceptance_as: checklist
trigger: manual               # manual is the only allowed value
```

- **`reads` is a bounded, declared allow-list.** A target may read only the artifact
  types it names, echoing the v0.5 discipline that the workspace never scans to fill
  a gap (INV-SEC-02 ethos). A `publish-plan` target reads `plans`/`tasks`; a future
  `notify-chat` target might read completion records only.
- **`direction: export` and `trigger: manual`** are the only values this scope
  allows; they exist as fields so the contract can refuse anything else rather than
  silently doing something coupling or automatic.

## Credentials never enter the workspace

`external-target.yaml` is credential-free. Provider tokens live in the host Git
configuration or the host agent / MCP server and never enter any workspace file or
record (INV-SEC-01). `target_ref` holds only a *logical* destination (a list name,
a project key, a repo id) — never a secret, never a session token. Anything
genuinely instance-specific but non-secret (a numeric account id) may go in the
gitignored `.local.yaml`, mirroring how machine-specific binding lives in
`repositories.local.yaml`.

## The authoritative mapping (in the plan directory)

```text
plans/<plan-id>/
  plan.yaml                                   # canonical identity + status (unchanged, pure)
  PLAN.md
  tasks/NNN-<slug>.md
  external/<target-id>.yaml                   # AUTHORITATIVE mapping — the idempotency key
```

```yaml
# plans/0023-.../external/team-clickup.yaml
target: team-clickup
provider: clickup
plan_item: { id: "86abc", url: "https://…" }   # the plan's external work item
task_items:
  - { task: "001", id: "86def", url: "https://…", synced_digest: "sha256:…" }
  - { task: "002", id: "86ghi", url: "https://…", synced_digest: "sha256:…" }
```

- `plan.yaml` is **not** touched — the mapping is a sibling file, so `plan.yaml`
  stays canonical-identity-only and portable (INV-PLAN-01). The external item id is
  environment-specific projection state, not plan identity, and must not live on the
  plan itself.
- `synced_digest` (via the runtime's existing `cc_digest`) makes re-publish
  idempotent and lets an unchanged task be skipped: re-running compares the current
  artifact digest to the recorded one and updates the *same* external item instead
  of creating a duplicate.

## The derived catalog is never a source of truth

`external-targets/<id>/state/catalog.md` is a convenience index (target → plans,
external id → plan) regenerable by scanning the per-plan `external/` files — exactly
the relationship `plans/INDEX.md` has to the plan directories, and the same
principle as the context index being a catalog, not a copy (INV-KNOWLEDGE-01). If
the catalog and a per-plan mapping ever disagree, the per-plan file wins. Archiving
a plan simply drops its row from the catalog on the next regeneration; no orphaned
authoritative record is ever left behind.

## What writes what

The host / MCP layer makes the external API call and obtains the item ids; the
deterministic write of the mapping into the plan directory reuses the runtime's
existing atomic-write path (`cc_atomic_write`). "Written by the runtime" never
implies "lives under `.runtime/`" — the path is an argument, and this record's home
is the plan directory because its lifecycle is the plan's, not an execution's.
