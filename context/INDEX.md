# Context Circuit Product Knowledge index

This is the agent retrieval catalog and a human navigation aid
(INV-KNOWLEDGE-01). Per `.context-circuit/wrapper/contracts/schemas/context-index.yaml` the
catalog carries per-unit retrieval metadata (context id, summary, topics,
aliases, domains, repositories, decisions, constraints, status, freshness,
provenance) for live context files. Provenance pointers, when present, are
durable retrieval facts only (INV-KNOWLEDGE-03). The product source identity is in
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
- provenance rules: `context/SOURCES.md`
- design deltas: `context/DESIGN-DELTAS.md`
- bounded domains: `context/domains/`
- cross-domain roles: `context/roles/`
- external references: `context/references/`

The source inbox and instantiated workspace plans are separate layers. Raw
sources are passive and request-scoped; ordinary entry never scans them.

Domain pages added and accepted in the v1.0 re-ground, closing the gap where
v1.0's core additions had no owning domain page:
[intent](domains/intent/README.md) (Gate 1 front door),
[tracing](domains/tracing/README.md) (post-approval planner + feasibility check), and
[assurance](domains/assurance/README.md) (consequence-tier ladder).

Direct collaboration was accepted as the Explore tier of the v1.0 assurance
ladder (`cc-pair` is not a separate mode): a new
[direct-collaboration](domains/direct-collaboration/README.md)
domain plus extensions to host-adapters, delivery, ARCHITECTURE.md, and
DECISIONS.md. See the 2026-09-03 decision in `context/DECISIONS.md`.

Committed `.claude/`, `.codex/`, and `.cursor/` trees are the native
integration surface. See the 2026-09-08 host-native decision in
`context/DECISIONS.md` and the updated [host-adapters](domains/host-adapters/README.md)
and [source-release-and-upgrade](domains/source-release-and-upgrade/README.md)
domain pages.

Intent and plan ids allocate from per-member number bands (committed roster,
one-time local identity). See the 2026-09-08 bands decision in
`context/DECISIONS.md` and [repository-binding](domains/repository-binding/README.md),
[plan-review](domains/plan-review/README.md), and [intent](domains/intent/README.md).

Open questions on a new intent are numbered 1, 2, 3 like The plans so a person
can answer by number; already-written intents stay as authored. See the
2026-09-10 decision in `context/DECISIONS.md` and [intent](domains/intent/README.md).

Official template publication is GitHub-canonical; the same Action optionally
mirrors the published tree and Release to GitLab when operators configure
GitHub Actions variables and secrets. See the 2026-09-09 decision in
`context/DECISIONS.md` and
[source-release-and-upgrade](domains/source-release-and-upgrade/README.md).
The GitHub destinations are `kaotypr/context-circuit-source` (this maintainer
source) and `kaotypr/context-circuit` (the published template). See the
2026-09-10 decision in `context/DECISIONS.md`.
