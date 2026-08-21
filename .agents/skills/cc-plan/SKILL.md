---
name: cc-plan
description: Draft, review, or prepare a human-facing plan bundle without executing it.
---

Use for idea briefs, PRDs, context-grounded plans, and read-only plan review.
The human entry is `PLAN.md`; `plan.yaml` owns canonical fields and lifecycle
status. Tasks reference acceptance and verification IDs instead of copying
commands. Review reports readiness and unresolved decisions but never approves,
claims a lease, starts children, or changes status.

Task Markdown files must begin with a line containing exactly `---`, contain
their YAML task frontmatter, close it with a second line containing exactly
`---`, and only then start the Markdown body. `plan.yaml` is pure YAML and does
not use Markdown frontmatter.

For approval, show the exact card from `docs/gates.md`; only a current explicit
confirmation may change `draft` to `approved` and project tasks to `ready`.
