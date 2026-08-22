# Context Circuit source

This repository builds Context Circuit 1.0, a filesystem-backed wrapper for
AI-assisted project work. The released workspace is a clean, uninitialized
template; this checkout is the product source and maintainer workspace.

## Human entry

In a released workspace, say:

> Start or resume work in this workspace.

The agent reads a compact safety spine, selects one bounded evidence probe, then
selects exactly one eligible action. Humans operate through the root README, a
plan's `PLAN.md`, the conversational “what's next” card, and surfaced handoffs.

Approval, execution, completion, delivery, archive, takeover, publication,
deployment, and cleanup are separate human actions. The core works offline and
stores no credentials.

## Source layout

- `wrapper/` — versioned shipped adapters, contracts, schemas, runtime guards,
  and migrations.
- `template/` — blank mutable seed for new workspaces.
- `plans/` — source-only maintainer implementation plans; not released.
- `repositories.local.yaml` — ignored host-local repository bindings; never released.
- `repositories/` — optional ignored convenience checkouts; never released.
- `.agents/skills/` and `agents/` — thin host and role adapters.
- `context/` — Product Knowledge for the Context Circuit product itself.
- `test/` — semantic suites, budgets, upgrade, and release verification.
- `scripts/` — maintainer-only release assembly; nothing here is user workflow.

Read `sources/context-circuit-design/DESIGN-SPEC.md` only as maintainer design
material. The released artifact uses the shipped wrapper and template, not the
source repository's maintainer state.
