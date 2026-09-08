# Context index

This is the agent retrieval catalog for Product Knowledge. It maps concepts,
aliases, repositories, domains, decisions, and constraints to knowledge units so
an agent can locate the right knowledge without scanning the whole directory.

Retrieval metadata and entry shape are owned by
`.context-circuit/wrapper/contracts/schemas/context-index.yaml`.

Knowledge homes: owned concepts live under `context/domains/` (scoped by
`repositories`) and `context/roles/`. Knowledge about an external service the
workspace consumes but does not own lives under `context/references/<service>/`
— no `repositories` tag; that directory is created when the first such service
is documented.

## Knowledge units

No accepted context units yet. Add units during context gathering, each with a
stable context ID, summary, topics, domains, repositories, decisions,
constraints, status, freshness, and provenance.
