---
name: cc-plan
description: Draft, review, or prepare a human-facing plan bundle without executing it, including named-plan read-only review.
---

Use for idea briefs, PRDs, context-grounded plans, and named-plan read-only
review — for example `Review plan <id>` or `Walk me through plan <id>`. The
human entry is `PLAN.md`; `plan.yaml` owns canonical fields and lifecycle
status. Tasks reference acceptance and verification IDs instead of copying
commands. There is no separate review skill; this is the only discovery
adapter for both drafting and reviewing a plan.

When entry routes a generic build or implementation request to `draft-plan`,
create only the bounded draft plan bundle. Do not create implementation files,
claim a lease, create a worktree, or treat the original build request as plan
approval or execution authorization.

## Read-only review mode

When the engine selects `review-plan`, load only the `plan-review` context
set and emit the Review Card from `docs/plan-review.md`. In this mode never
draft, approve, claim a lease, start a child, or change plan/task status —
review only reports outcome, risks, and open decisions. A review request with
no usable plan id routes to `clarify-target` instead; ask which plan, do not
guess or inspect an unrelated bundle.

When the Review Card carries one to three focused human decisions, also offer
those decisions through the current host's optional native question prompt
(see `docs/plan-review.md`). A missing, denied, or failed prompt is not
`host-blocked`; fall back to the card text and never retry the prompt as a
required child. A chosen next-action label may continue only into that
route's existing first card in the same session — selecting it never skips
confirmation or changes status.

Task Markdown files must begin with a line containing exactly `---`, contain
their YAML task frontmatter, close it with a second line containing exactly
`---`, and only then start the Markdown body. `plan.yaml` is pure YAML and does
not use Markdown frontmatter.

For approval, hand off to `cc-gates` and show the exact card from
`docs/gates.md`; only a current explicit confirmation there may change `draft`
to `approved` and project tasks to `ready`. `cc-plan` does not keep its own
copy of that confirmation procedure.
