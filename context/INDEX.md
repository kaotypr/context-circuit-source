# Context Circuit Product Knowledge index

This is navigation only. The product source identity is in `workspace.yaml`;
the wrapper route owner is `wrapper/contracts/routes.yaml`; context budgets
are owned by `wrapper/contracts/context-sets.yaml`; runtime shapes are owned by
the wrapper schemas.

Tier 0 summaries: `context/WORKSPACE.md`, `context/PROJECT.md`.

Route-selected knowledge:

- architecture: `context/ARCHITECTURE.md`
- conventions: `context/CONVENTIONS.md`
- accepted decisions: `context/DECISIONS.md`
- provenance rules: `context/SOURCES.md` and `context/sources.yaml`
- bounded domains: `context/domains/`
- cross-domain roles: `context/roles/`

The source inbox and instantiated workspace plans are separate layers. Raw
sources are passive and request-scoped; ordinary entry never scans them.
