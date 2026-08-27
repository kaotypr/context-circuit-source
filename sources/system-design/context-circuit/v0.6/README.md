# Context Circuit v0.6

Source design for Context Circuit v0.6 — a **delta on v0.5**. Read the v0.5 design
first ([../v0.5/](../v0.5/)); v0.6 changes only what it names. This README is the
version index; each scope owns its own design.

## Scopes

- [run-stack/](./run-stack/) — execute a set of approved plans in one run (a *plan
  stack*): inter-plan dependencies, path leases, execution bases, the scheduling
  loop, failure containment, and the delivery drift guard. Start at
  [run-stack/design.md](./run-stack/design.md).
- [system-design-authoring/](./system-design-authoring/) — a **skill** for
  authoring and structuring a **system design** as source material: the three-tier
  layout, how much detail per file, and scope separation by concern. It adds no
  lifecycle, status, gate, or runtime — a system design is one kind of source.
  Start at [system-design-authoring/design.md](./system-design-authoring/design.md).
- [repository-grounding/](./repository-grounding/) — ground the **writer** in the
  target repository's own agent guidance (discovered live) and hand it a prepared
  worktree, delivered through a generated brief instead of hand-authored prose.
  Start at [repository-grounding/design.md](./repository-grounding/design.md).
- [external-surface/](./external-surface/) — a **manually-triggered, config-driven**
  way to publish Context Circuit data to external systems (task trackers, chat,
  docs), **fully orthogonal to the core workflow** — a publication is a peer
  command, never a phase. Publications live under `publication/<name>/` (e.g.
  `plans-clickup`); the publish command is the `cc-publish` skill and the first
  *kind* is `plan` (a plan and its tasks → ClickUp / Jira / GitHub / Notion / …).
  Start at [external-surface/design.md](./external-surface/design.md).

These four scopes are independent product capabilities released together in v0.6.
run-stack and repository-grounding share a coordinated contract bump (plan schema
`[1, 2]`, `runtime_version 0.6.0`); system-design-authoring ships only a skill and
external-surface ships a skill plus config/record schemas, both with no core
contract change.

## Maintainer tooling (not a product scope)

- [template-harness/](./template-harness/) — a delta on the v0.5 human-simulated
  test harness: make **dimension D (the efficiency ledger)** real by emitting
  per-action usage telemetry from the runner's own result, defining the budget
  units, and having the grader compare observed usage to each case's `budgets`
  (still soft, warning-only). Source-only maintainer tooling; changes no product
  surface and ships nothing. Start at
  [template-harness/design.md](./template-harness/design.md).

## Layout convention

`sources/system-design/<product>/<version>/<scope>/`. Every folder has a
`README.md` landing/index; each scope's normative design is `design.md`, with
detail split into files or sub-folders as it grows. See the v0.5 version README's
"Layout convention" for the full rule.

## Authority

v0.6 adds no owner and duplicates no rule. It specifies intent and points to the
canonical owners under `wrapper/`. Contract deltas are summarized in each scope's
`schema-and-*` file.
