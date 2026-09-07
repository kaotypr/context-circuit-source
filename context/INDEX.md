# Context Circuit Product Knowledge index

This is the agent retrieval catalog and a human navigation aid
(INV-KNOWLEDGE-01). Per `.context-circuit/wrapper/contracts/schemas/context-index.yaml` the
catalog carries per-unit retrieval metadata (context id, summary, topics,
aliases, domains, repositories, decisions, constraints, status, freshness,
provenance) for live context files. The product source identity is in
`workspace.yaml`;
the one-rule-one-owner map is `.context-circuit/wrapper/contracts/invariants.yaml`; conversational
routing is owned by `.context-circuit/agents/coordinator.md` and the `cc-*` skills; the retrieval
catalog shape is owned by `.context-circuit/wrapper/contracts/schemas/context-index.yaml`; runtime
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
- external references: `context/references/`

The source inbox and instantiated workspace plans are separate layers. Raw
sources are passive and request-scoped; ordinary entry never scans them.

Domain pages added and accepted in the 2026-09-04 v1.0 re-ground, closing the gap
where v1.0's core additions had no owning domain page:
[intent](domains/intent/README.md) (Gate 1 front door),
[tracing](domains/tracing/README.md) (post-approval tracer + feasibility check), and
[assurance](domains/assurance/README.md) (consequence-tier ladder).

The direct-collaboration proposals `0029`–`0033` were accepted on 2026-09-03,
rewritten to the v1.0 assurance-ladder framing (`cc-pair` is the Explore tier, not
a separate mode): a new [direct-collaboration](domains/direct-collaboration/README.md)
domain plus extensions to host-adapters, delivery, ARCHITECTURE.md, and
DECISIONS.md. See the 2026-09-03 decision in `context/DECISIONS.md`.
