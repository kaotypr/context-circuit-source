# Context Circuit v0.6 — External-surface (overview)

Status: authoritative source design for the v0.6 external-surface scope (delta on v0.5)
Revision: 2 — 2026-08-27

This is the **overview** of the external-surface scope: the capability, the
principles, and the shape of the solution. Each mechanism has its own detail file
(see [Detailed design](#detailed-design)). Read [README.md](./README.md) first for
the index, and [../../v0.5/core/design.md](../../v0.5/core/design.md) for
everything this delta builds on.

## The one capability

v0.5 keeps everything a workspace produces — plans, tasks, completion records,
context — **inside the workspace**. v0.6 adds a way to take that data and publish it
to an **external system** the team already lives in (a task tracker, a chat channel,
a docs space), on demand.

A configured pipeline is a **publication**: a user-declared, manually-triggered
command that *reads Context Circuit artifacts as it finds them and reflects them
outward*. A publication is **not** a lifecycle stage, not a step of any phase, and
not something the core workflow ever calls.

```mermaid
flowchart LR
  subgraph core["Core workflow — unchanged, self-contained"]
    direction LR
    P["plan"] --> A["approve"] --> E["execute"] --> V["verify"] --> D["deliver"]
  end
  subgraph ext["Publications — separate, user-driven"]
    direction LR
    U(["user runs<br/>/cc-publish 0023"]) --> K["cc-publish reads<br/>the workspace as-found"] --> X["ClickUp / Jira /<br/>GitHub / Slack / …"]
  end
  core -.->|"NO edge — the two never connect"| ext
```

The dashed line is crossed out on purpose: **there is no edge** between the core
workflow and the publications, in either direction. They are peers that happen to
share a workspace. A publication runs only when, and exactly when, a human invokes
it.

## Problem

A team running Context Circuit still lives in ClickUp / Jira / GitHub / Slack.
Today the only way to reflect a plan there is to retype it by hand, which drifts
immediately and scales badly. The naïve fix — have `deliver` (or any phase) also
post to the tracker — is exactly wrong: it couples an external, credentialed,
provider-specific side effect to the deterministic core loop, adds weight to the
runtime and the agent, and quietly hands an outside system a foothold on plan
state. v0.6 gives the reflection a home that touches none of that.

## Goals

1. Let a user publish Context Circuit data (starting with a plan and its tasks) to
   an external system with one explicit command.
2. Keep it **opt-in and manual** — nothing happens unless the user asks, every time.
3. Make it **general**: one backbone that hosts many *kinds* of publication
   (plans → a tracker, docs → a docs space, open questions → a chat thread), with
   plan publishing as the first.
4. Add **zero weight** to the core workflow — no runtime code, no agent load, no
   reference anywhere in the plan/approve/execute/verify/deliver path.
5. Keep the external copy **self-contained and non-authoritative**: it reads as
   ordinary project work to someone who has never heard of Context Circuit, and
   `plan.yaml` remains the single source of truth for plan identity and status.

## Non-goals

- **No automatic triggering.** Never on approve, never on deliver, never on a
  timer. Manual invocation is the only entry.
- **No coupling to the workflow.** No hook, no callback, no lifecycle listener,
  no "emit on phase X." Nothing to wire, because nothing triggers it.
- **No import (this scope is export-first).** Data flows Context Circuit →
  outward. Pulling external items *into* plans would bypass the human
  plan-authoring and approval gates and is out of scope; any future import must go
  through the normal authoring gate, never around it.
- **No bidirectional sync / no status flow-back.** The external item may drift
  (a human edits it); that is cosmetic and never returns to `plan.yaml`.
- **No provider code in the runtime** (v0.5 INV-RUNTIME-01). Network and
  credentialed calls live in the host / MCP layer.
- **No new authority.** A publication never changes plan status, approves,
  executes, verifies, completes, or delivers.

## Principles

- **Bind to data, never to control flow.** This is what makes "zero core-flow
  weight" a structural fact rather than a coding discipline. A publication is a
  *downstream reader* of workspace artifacts; the workflow is a *producer* that is
  unaware anything reads it.
- **A publication is a command, not a link.** It has no standing relationship to a
  plan — it runs once, when invoked, against whatever is on disk at that instant.
- **The external copy is a self-contained reflection.** One-way, non-authoritative,
  idempotent on re-run, and written in plain language a reader without workspace
  access understands; the workspace stays canonical.
- **The host carries the provider weight.** The runtime stays the small
  deterministic library it already is; adapters reach providers through the host
  (MCP), and configuration files never hold credentials.

## The backbone (what this scope actually adds)

Three small pieces, and nothing more:

1. **A user-owned `publication/` folder.** One sub-folder per configured
   publication, named for what it publishes and where — `plans-clickup`,
   `plans-github`, `docs-clickup`, `thread-slack`. Each holds a `config.yaml` and a
   `published/` folder of records. It is **created on first use**, not shipped empty
   into every workspace. **Everything a publication reads or writes lives here;
   nothing is ever written under `plans/`.** Detail in
   [configuration-and-records.md](./configuration-and-records.md).
2. **A manual trigger** — the **`cc-publish`** skill (`.agents/skills/cc-publish/SKILL.md`),
   invoked by name or slash command `/cc-publish` (v0.5 INV-SKILL-01). It publishes
   a publication according to the `kind` in its `config.yaml` — the `plan` kind
   here, other kinds as the surface grows. Skills are resolved only when explicitly
   invoked, so the trigger is isolated by construction. The coordinator runs
   `cc-publish` as a read-as-procedure skill (INV-SKILL-01), the same way it runs
   `cc-deliver` and `cc-archive`.
3. **An isolation contract** — invariants stating the orthogonality, the
   export-only data boundary, and that every published artifact is self-contained
   (no workspace file, path, id, or internal mechanism leaks outward). Detail in
   [contracts.md](./contracts.md).

No event system is needed: because nothing in the workflow triggers a publication,
there is nothing to hook.

## Publications are named for what they publish and where

Each folder under `publication/` is one configured publication whose name is,
conventionally, `<subject>-<provider>`:

| Publication (folder) | kind | provider | publishes |
| --- | --- | --- | --- |
| `plans-clickup` | plan | clickup | a plan and its tasks |
| `plans-github` | plan | github | a plan and its tasks |
| `thread-slack` | thread | slack | a plan's open questions, as a discussion |
| `docs-clickup` | docs *(future)* | clickup | context docs / decisions |

The `plan` and `thread` kinds are designed in this scope ([publish-plan.md](./publish-plan.md),
[thread.md](./thread.md)); `docs` is illustrative. A new kind is new behavior in
`cc-publish` plus a `config.yaml`, nothing more — the same backbone hosts them all
without touching the core workflow.

## What changes relative to v0.5

| Area | v0.5 | v0.6 (this scope) |
| --- | --- | --- |
| Reach of workspace data | stays inside the workspace | can be published outward, on demand |
| Trigger | — | the manual `cc-publish` skill / slash command, never the workflow |
| Config + records | workspace / plan / context files | adds a user-owned `publication/` (create-on-first-use) |
| External records | — | per-plan records under `publication/<name>/published/`; nothing under `plans/` |
| Core workflow | plan → … → deliver | **unchanged**; gains no reference to publications |
| Runtime | deterministic library | **unchanged**; no provider/network code (INV-RUNTIME-01) |

Everything else in v0.5 is unchanged. In this product, "publish" and "publication"
name sending data to an external system; git delivery is "push the branch" and
"open a pull request" (INV-DELIVER-01), never "publish." The two never share a word,
so no qualifier is needed (see [contracts.md](./contracts.md)).

## Detailed design

- [configuration-and-records.md](./configuration-and-records.md) — the
  `publication/<name>/` folders, `config.yaml`, create-on-first-use, the per-plan
  records under `published/`, on-demand cross-plan lookup, credentials boundary.
- [publish-plan.md](./publish-plan.md) — the `plan` kind: plan → work-item,
  task → child-item across ClickUp / Jira / GitHub / Notion (provider is a config
  field, list open-ended); containment vs dependency; self-contained external text;
  the `[NNNN]` title convention; one-way idempotent reflection; provider wrinkles;
  a worked trace.
- [thread.md](./thread.md) — the `thread` kind: a plan's open questions → a chat
  thread (Slack first); a `[thread]`-prefixed parent plus one fully-described reply
  per question; discussion-safe idempotency (edit only its own messages, never human
  replies); its own record shape.
- [contracts.md](./contracts.md) — proposed invariants, owner-map additions, the
  publish/git-delivery vocabulary split, and why no core contract bumps.

## Compatibility (summary)

Fully additive and opt-in. No existing plan, schema, runtime version, or workflow
phase changes. A workspace that configures no publication is byte-for-byte a v0.5
workspace plus the availability of the `cc-publish` skill. The plan schema does **not**
bump and `execution.yaml` does not change; unlike run-stack and repository-grounding
this scope needs no coordinated contract bump.

## Implementation order

Bounded, independently reviewable phases:

1. **Isolation contract** — add the external-surface invariants and owner-map
   entries; assert the core workflow references nothing here.
2. **Config + record schemas** — `config.yaml` and the per-plan record; the
   `publication/<name>/` layout; create-on-first-use; on-demand cross-plan lookup.
3. **`cc-publish` skill** — the `plan` kind: the plan→work-item mapping,
   self-contained external text, idempotent re-run, and the provider realizations,
   driven through host/MCP.
4. **Semantic verification** — a worked trace (configure a publication, publish a
   plan, re-publish and update in place) plus an assertion that the core acceptance
   suite is unchanged by the surface's presence.

This design does not authorize implementation, delivery, or publication by itself.

## Final design decisions

- The external surface is **orthogonal to the core workflow** — a publication is a
  peer command, never a phase, trigger, gate, dependency, or side effect of one.
- A publication runs **only on explicit manual invocation**, every time; no
  automatic, scheduled, or workflow-driven trigger exists.
- A publication is a **downstream reader**: it reads declared workspace artifacts
  and writes only its own records under `publication/`; it never mutates core
  Context Circuit state or plan status, and **nothing is written under `plans/`**.
- The scope is **export-first**; import (external → plan) is out of scope and, if
  ever added, must flow through the normal plan-authoring gate.
- The external copy is **one-way, non-authoritative, and self-contained**;
  `plan.yaml` stays canonical (INV-PLAN-01) and no workspace file, path, id, or
  internal mechanism ever appears in an external artifact (a plan id in a title is
  the one allowed cross-reference).
- Configuration lives in a **user-owned `publication/<name>/config.yaml`, created on
  first use**, and is **credential-free**; credentials stay at the host / MCP layer
  (INV-SEC-01).
- Publications are triggered by the **`cc-publish`** skill (INV-SKILL-01), which
  publishes per the publication's `kind`; the host carries all provider/network
  weight (INV-RUNTIME-01 unchanged).
- The concept and this source scope are named **external-surface**; the product
  root folder is **`publication/`**; a publication folder is named `<subject>-<provider>`
  (e.g. `plans-clickup`); its config file is **`config.yaml`**; the publish command
  is the single **`cc-publish`** skill; the designed kinds are **`plan`** and
  **`thread`**.
- Every kind honors optional **`instructions`** config — free-text guidance to
  `cc-publish` covering *the language to author in* (default: the plan's language),
  *how* external text reads (tone, phrasing, term handling), and *which optional
  provider fields to enrich* by estimation (e.g. a time estimate per task, a target
  date range), written one-way as best-effort estimates. Instructions guide wording
  and optional field values only; they never override the self-contained rule
  (INV-EXTERNAL-03), expand `reads`, write back to the workspace, change the core
  mapping, or relax any boundary.
- The **`thread`** kind publishes a plan's open questions as a chat discussion — a
  **`[thread]`**-prefixed parent plus **one fully-described reply per question** —
  and is discussion-safe: on re-run it edits only its own messages and never a
  human's reply, and it never deletes messages. Its record shape is its own; record
  shapes are per-kind.
- The scope adds **no core contract bump** — no plan-schema, execution, or
  runtime-version change.
