# Product Knowledge

`context/` holds accepted, durable, agent-oriented knowledge about the project:
purpose, architecture, conventions, decisions, domain rules, roles, terminology,
and provenance. It is human-readable, but it exists first so an agent can locate
and apply the right knowledge when creating a plan.

`context/TERMINOLOGY.md` is this project's own domain glossary — the words the
team and codebase use, with an agreed meaning for each. It is distinct from the
Context Circuit product glossary in `docs/terminology.md`; use the template at
`docs/templates/terminology-context.md` when proposing glossary units.

`context/INDEX.md` is the retrieval catalog: it maps concepts, aliases,
repositories, domains, decisions, and constraints to knowledge units so an agent
selects the relevant units without scanning the directory. Retrieval metadata
and the entry shape are owned by `wrapper/contracts/schemas/context-index.yaml`.

Knowledge changes only through explicit human acceptance. Gathering context,
completion reconciliation, repository evidence, and worker claims can produce
proposals under `context/proposals/`, but a proposal is never applied
automatically. A plan may be `done` while a proposal is pending; relevant
pending proposals surface during future plan creation. `sources/` is passive
raw evidence and is read only when named.
