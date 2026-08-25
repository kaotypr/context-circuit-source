# Context Circuit Product Knowledge index

This is the agent retrieval catalog and a human navigation aid
(INV-KNOWLEDGE-01). Per `wrapper/contracts/schemas/context-index.yaml` the
catalog also carries per-unit retrieval metadata (context id, summary, topics,
aliases, domains, repositories, decisions, constraints, status, freshness,
provenance) and indexes pending context proposals and stale-context warnings by
target context id, related plan, topic, repository, and impact status. The
product source identity is in `workspace.yaml`;
the one-rule-one-owner map is `wrapper/contracts/invariants.yaml`; conversational
routing is owned by `agents/coordinator.md` and the `cc-*` skills; the retrieval
catalog shape is owned by `wrapper/contracts/schemas/context-index.yaml`; runtime
shapes are owned by the wrapper schemas.

Tier 0 summaries: `context/WORKSPACE.md`, `context/PROJECT.md`.

Route-selected knowledge:

- architecture: `context/ARCHITECTURE.md`
- conventions: `context/CONVENTIONS.md`
- accepted decisions: `context/DECISIONS.md`
- terminology: `context/TERMINOLOGY.md`
- provenance rules: `context/SOURCES.md` and `context/sources.yaml`
- design deltas: `context/DESIGN-DELTAS.md`
- bounded domains: `context/domains/`
- cross-domain roles: `context/roles/`

The source inbox and instantiated workspace plans are separate layers. Raw
sources are passive and request-scoped; ordinary entry never scans them.
