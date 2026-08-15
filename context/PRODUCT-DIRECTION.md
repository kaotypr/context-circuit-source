# Product Direction

Status: accepted Product Knowledge summary.

Context Circuit is a reusable AI-agent workspace template for solo developers
and teams working across one or more repositories. Its primary interface is a
conversational agent session that carries routine organization while the user
decides intent, scope, publication, merge, deployment, and ownership boundaries.

The core journey is:

`start or resume → understand intent → capture the right artifact → ground it
in Project Knowledge → plan → approve (`cc-approve-plan`) → execute and verify
(`cc-run-plan`) → finish (`cc-finish-plan`) → optional cleanup
(`cc-cleanup-runtime`) → handoff`

The core discoverable set includes `cc-approve-plan`, `cc-finish-plan`, and
`cc-cleanup-runtime` alongside the existing planning and execution skills.
These names are discoverability aids, not mandatory ceremonies.

A new project with no repository is valid. Raw files in `sources/` are passive
and may be read only for a requested source-based activity. Product Knowledge,
Idea Briefs, PRDs, plans, and private runtime state remain separate layers.

`cc-run-plan` directs child writers and verifiers in every workspace mode.
The same behavior applies to solo and team workspaces. `workspace.yaml` mode
is initialization identity; it does not select an execution topology.
Delivery policies `solo-local` and `team-review` remain post-verification
choices.

Source: `sources/context-circuit-product-direction.md` in the owning workspace;
this file is a concise summary, not a copy of that source.
