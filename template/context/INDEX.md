# Context index

This is the agent retrieval catalog for Product Knowledge. It maps concepts,
aliases, repositories, domains, decisions, and constraints to knowledge units so
an agent can locate the right knowledge without scanning the whole directory.

Retrieval metadata and entry shape are owned by
`wrapper/contracts/schemas/context-index.yaml`.

## Knowledge units

No accepted context units yet. Add units during context gathering, each with a
stable context ID, summary, topics, domains, repositories, decisions,
constraints, status, freshness, and provenance.

## Pending context impacts

None. Completion reconciliation and stale-context warnings are listed here by
target context ID, related plan, topic, repository, and impact status.
