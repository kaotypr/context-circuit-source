# Plan artifacts and their reading map

A developer may ask for a grounded plan of an already specified outcome without
creating an intent. An approved intent also leads to planning. Either route
creates a folder with `plan.md` as its entry; only the linked route writes an
`intent` field and backlink. A planning request authorizes investigation and
drafting. Implementation starts after the complete plan is presented and the
person separately requests execution.

`plan.md` carries the outcome, repository scope, approach, ordered tasks, risks,
checks, and a Details section. Add supporting files in that folder only when
they make the design easier to inspect. Each file is linked from Details and
assigned in `required_files.shared` or `required_files.by_repository.<repo>`.
Paths are relative to the plan folder and may not escape it or pass through a
symlink. The same file may serve several repositories. A worker receives the
entry's body plus only its assigned files as paths relative to a stated absolute
workspace root; the worker reads them before editing. A multi-repository plan
needs an explicit repository assignment at dispatch.

The CLI resolves each plan ID to one entry, whether it is a folder plan or a
legacy single-file plan. It does not migrate old plans. The folder is
one record and one completion boundary; supporting files have no IDs or status.
`check` reports structural drift between the entry, reading map, and files.

An execution or delivery run names its exact plans with repeatable `record order
--plan ID`, or selects an intent's linked plans with `--intent ID`. Unfinished
dependencies outside the selection block work; they are never added silently.

Owner:

- `context-circuit-source@internal/workspace/records.go` — folder creation and record lookup.
- `context-circuit-source@internal/workspace/planfiles.go` — reading-map validation.
- `context-circuit-source@internal/workspace/agents.go` — repository-specific worker briefs.
- `context-circuit-source@product/AGENTS.md.in` — the planning and execution decisions.
