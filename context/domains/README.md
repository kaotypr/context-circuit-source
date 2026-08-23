# Domain Knowledge

`context/domains/<domain>/` is the canonical home for knowledge about one
bounded project area. The domain page owns exact behavior, workflows, rules,
interfaces, data, constraints, edge cases, implementation references, and
verification when those facts are supported by evidence.

Generated domain pages use [the domain template](../../docs/templates/domain-context.md)
and normally live at `context/domains/<domain>/README.md`. Detailed workflow
pages may live below that domain at `workflows/<workflow>.md`.

## Selecting domain context

For a domain request, first read this index and then only the named domain page
and the workflow pages it links to. Do not recursively read every domain or the
raw source inbox. If the request spans domains, select the smallest set that
covers the stated outcome and record the selection in the plan or handoff.

Domain pages can be `proposed`, `accepted`, or `needs-review`. Proposed pages
are useful evidence but are not accepted Product Knowledge until a human gate
is satisfied. A refresh of an accepted page must preserve accepted decisions;
conflicting evidence is surfaced for review rather than silently replacing it.

## Domain entries

- [Repository binding and bootstrap](repository-binding/README.md) — shared identity, host-local bindings, explicit bootstrap, isolated worktrees. Status: proposed.
- [Host adapters](host-adapters/README.md) — Codex, Claude Code, and Cursor Agent CLI as evidence-only hosts. Status: proposed.
- [Named-plan review](plan-review/README.md) — read-only Review Card and optional host question prompts. Status: proposed.
- [Plan approval and product-source commit](plan-approval/README.md) — two-turn status-only approval and the maintainer commit card. Status: proposed.
