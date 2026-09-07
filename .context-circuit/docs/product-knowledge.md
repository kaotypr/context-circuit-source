# Product Knowledge

`context/` holds accepted, durable, agent-oriented knowledge about the project:
purpose, architecture, conventions, decisions, domain rules, roles, terminology,
and provenance. It is human-readable, but it exists first so an agent can locate
and apply the right knowledge when creating a plan.

`context/TERMINOLOGY.md` is this project's own domain glossary — the words the
team and codebase use, with an agreed meaning for each. It is distinct from the
Context Circuit product glossary in `.context-circuit/docs/terminology.md`; use the template at
`.context-circuit/docs/templates/terminology-context.md` when proposing glossary units.

`context/INDEX.md` is the retrieval catalog: it maps concepts, aliases,
repositories, domains, decisions, and constraints to knowledge units so an agent
selects the relevant units without scanning the directory. Retrieval metadata
and the entry shape are owned by `.context-circuit/wrapper/contracts/schemas/context-index.yaml`.

Knowledge changes by writing live `context/` files in place. Gathering context
edits those files and keeps the catalog consistent. After a plan is marked done,
the coordinator performs the same in-place reconcile when the plan affected
Product Knowledge. There is no sidecar staging path and no extra
knowledge-acceptance gate. Completion, verification, worker claims, and
delivery never write Product Knowledge. A later plan may start even if a
knowledge update has not landed. `sources/` is passive raw evidence and is read
only when named.
