---
kind: domain
status: accepted
title: External surface
slug: external-surface
owners: []
sources: []
source_revisions:
  - wrapper: HEAD
    commit: 6614841
    basis: current-wrapper
generated_at: 2026-08-28T00:00:00Z
review_date: 2026-11-28
freshness: accepted-from-current-wrapper
assumptions:
  - Publishing is orthogonal to the core workflow; it adds no authority and touches no plan status.
unknowns: []
contradictions: []
acceptance:
  state: accepted
  accepted_at: 2026-08-28
  accepted_by: maintainer
workflows:
  - .context-circuit/docs/getting-started.md
---

# External surface

## Summary

A manually-triggered, config-driven way to **publish** Context Circuit data to an
external system a team already uses — a task tracker, chat, or docs space. A
configured pipeline is a **publication**; the `cc-publish` skill runs it. Route
"publish plan `<id>` to ClickUp/Jira/GitHub/Notion", "mirror `<id>` to our tracker",
and "open a discussion thread for `<id>`'s open questions in Slack" here. Owned by
the three isolation invariants (INV-EXTERNAL-01/02/03) and the `cc-publish` skill
(INV-SKILL-01). It is **orthogonal to the core workflow**: it adds no runtime or
network code, no new authority, and no reference anywhere in the intent → plan →
execute → verify → deliver path.

In this product "publish"/"publication" names sending data to an external system;
git delivery is "push" / "open a pull request" ([delivery](../delivery/README.md))
and never "publish" — the two never share a word.

## Scope

Inside: the `publication/` folder and its `config.yaml`, the manual `cc-publish`
trigger, the isolation contract, the two implemented kinds (`plan`, `thread`), the
per-kind records, and the credential/host boundary.

Outside: git delivery ([delivery](../delivery/README.md)), and any inbound flow
(external → plan), which is out of scope; a plan is only ever authored through the
normal planning flow ([plan-review](../plan-review/README.md)) and authorized by
its parent approved intent.

## Behavior

- **Orthogonality (INV-EXTERNAL-01).** A publication is a peer command, never a
  phase, trigger, gate, dependency, or side effect. No workflow phase, runtime
  action, coordinator route, or role references, triggers, waits on, or is affected
  by it, and none runs except on an explicit human invocation. A plan in any state —
  `draft`, `done`, or never — can be published. There is no automatic,
  scheduled, or event trigger; nothing in the workflow triggers a publication, so
  there is nothing to hook.
- **Data boundary, export-first (INV-EXTERNAL-02).** A publication reads only the
  workspace artifact types its `config.yaml` names (`reads:`) and writes only its
  own records under `publication/<name>/published/`. **Nothing is written under
  `plans/`.** It never mutates `plan.yaml`, task files, or plan status
  (INV-PLAN-01); it is export-only. A record's lifecycle is independent of the
  plan's archive state, keyed to the stable plan id. Any future import must pass
  through the normal intent-approval gate (INV-APPROVE-01), never around it.
- **Self-contained external artifacts (INV-EXTERNAL-03).** Every external artifact
  (task, document, message) is understandable to a reader with no workspace access:
  plain-language description, and never a workspace file name, workspace-internal
  path, workspace-internal id, or internal Context Circuit mechanism/vocabulary. A
  stable plan or task identifier may appear (normally in a title, as
  `[<plan-number>]`) as the one allowed cross-reference — the external mirror of the
  product's existing "never expose internal mechanism" report hygiene.
- **The backbone.** (1) a user-owned `publication/` folder, one sub-folder per
  publication named `<subject>-<provider>` (e.g. `plans-clickup`, `thread-slack`),
  **created on first use**, never shipped empty; (2) the manual trigger — the
  `cc-publish` skill, invoked by name or `/cc-publish` and run as a read-as-procedure
  skill the same way `cc-deliver`/`cc-archive` are; (3) the isolation contract (the
  three invariants above).
- **Kinds.** `plan` — a plan → one first-class work item, each task → a child item,
  a task's `acceptance`+`verification` → a checklist, `depends_on` → a sibling
  dependency link (containment, never nesting). Status is set one-way at publish
  through the config `status:` map onto the provider's discovered status set and
  never read back; re-run is idempotent via a per-task `synced_digest` (the
  runtime's existing `cc_digest`), updating the same item and skipping unchanged
  tasks. `thread` — a plan's open questions → a chat discussion: a `[thread]`-prefixed
  parent plus one fully-described reply per question; discussion-safe (on re-run it
  edits only its own messages, appends a reply for a new question, and never edits or
  deletes a human's reply or any message). Record shapes are **per-kind**.
- **`instructions` (optional, per publication).** Free-text authoring guidance for
  every kind: the language to author in (default: the plan's language), tone /
  phrasing / term handling, and which optional provider fields to enrich by
  estimation (e.g. a time estimate per task, a target date range) written one-way as
  best-effort. It guides wording and optional field values only; it never overrides
  INV-EXTERNAL-03, expands `reads`, writes back to the workspace, changes the
  mapping, or relaxes a boundary.
- **Host carries provider weight; credential-free (INV-RUNTIME-01, INV-SEC-01).**
  External calls go through the host / MCP layer, never the runtime, which gains no
  provider or network action (it at most reuses its atomic-write for the record).
  `config.yaml` holds only a logical `target_ref` (a list name, project key, repo
  id, channel) — never a token. If the host has no usable provider capability,
  `cc-publish` reports `host-blocked` and creates nothing external.

## Workflows

- Publish a plan to a tracker as a separate step: `.context-circuit/docs/getting-started.md`

## Interfaces

- Human request: "Publish plan `<id>` to `<provider>`" / "open a discussion thread
  for `<id>`'s open questions"
- Config: `publication/<name>/config.yaml` (`kind`, `provider`, `direction: export`,
  `trigger: manual`, `reads`, `target_ref`, optional `instructions`, optional
  `mapping`)
- Records: `publication/<name>/published/<plan-id>.yaml` (per-kind shape)

## Data

Per-plan records under `publication/<name>/published/`: the `plan` kind records
`plan_item` + `task_items[{task, id, synced_digest}]`; the `thread` kind records
`thread{channel, parent_ts}` + `questions[{ref, reply_ts}]`. Both hold ids and
timestamps only — no credentials, no message text, no external id on `plan.yaml`.

## Constraints and edge cases

Export-only and one-way: a human's later edit in the tracker is cosmetic drift and
never returns to the workspace. `direction` other than `export` and `trigger` other
than `manual` are refused, not silently honored. The `plan` kind never encodes
`depends_on` as nesting even where a provider makes nesting cheap. A `thread` never
deletes a message; a no-longer-open question is left in place.

## Implementation references

- `.agents/skills/cc-publish/SKILL.md`
- `.context-circuit/wrapper/contracts/invariants.yaml`: INV-EXTERNAL-01, INV-EXTERNAL-02,
  INV-EXTERNAL-03 (owners map: `external_surface`, `publication_config`,
  `publication_record`, `publication_thread_record`)
- `.context-circuit/wrapper/contracts/schemas/publication-config.yaml`,
  `publication-record.yaml`, `publication-thread-record.yaml`
- No `engine.sh` provider/network code (INV-RUNTIME-01); the record write reuses the
  existing atomic-write path and `cc_digest`.

## Verification

`sh test/acceptance.sh` (external-surface suite, registered at `test/acceptance.sh`);
`sh test/external-surface/test-external-surface.sh` → `PASS: external-surface`. The
scope is additive: the rest of the core acceptance suite is unchanged by the
surface's presence.

## Provenance

Authored from the current wrapper at HEAD `6614841`. Design source
`sources/system-design/context-circuit/v0.6/external-surface/` was named by the
accepting request.

## Acceptance notes

Accepted 2026-08-28 from proposal `0023-domain-external-surface`.
