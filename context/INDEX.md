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
- external references: `context/references/`

The source inbox and instantiated workspace plans are separate layers. Raw
sources are passive and request-scoped; ordinary entry never scans them.

Pending context proposals:

- `0029-add-direct-collaboration` → `context/domains/direct-collaboration/README.md`
  (add; direct collaboration, pairing, human supervision; repository:
  `context-circuit-source`; review needed).
- `0030-change-host-adapters-for-pairing` →
  `context/domains/host-adapters/README.md` (change; worker-child mapping and
  host-blocked pairing; repository: `context-circuit-source`; review needed).
- `0031-change-delivery-for-pairing` → `context/domains/delivery/README.md`
  (change; pairing delivery and anchor-drift block; repository:
  `context-circuit-source`; review needed).
- `0032-change-architecture-for-pairing` → `context/ARCHITECTURE.md` (change;
  orthogonal repository-changing modes; repository: `context-circuit-source`;
  review needed).
- `0033-change-decisions-for-pairing` → `context/DECISIONS.md` (change; fixed
  direct-collaboration decisions; repository: `context-circuit-source`; review
  needed).
