---
id: 0002-index-catalog-role
target_context_unit: context/INDEX.md
operation: change
statement: >
  INDEX.md opens "This is navigation only", which contradicts INV-KNOWLEDGE-01
  and the context-index schema: INDEX is an agent retrieval catalog AND a human
  navigation aid, and it should also index pending context proposals and
  stale-context warnings. Its current plain link list omits the per-unit
  retrieval fields the schema defines.
evidence_refs:
  - wrapper/contracts/schemas/context-index.yaml  # retrieval-catalog fields + indexed pending impacts
  - wrapper/contracts/invariants.yaml              # INV-KNOWLEDGE-01
  - docs/product-knowledge.md                      # "INDEX.md is the retrieval catalog"
affected_repositories:
  - context-circuit-source
affected_commits:
  - 4b8ac0b
confidence: high
status: review-needed
---

# Fix INDEX.md self-description and catalog role

## Corrections

- Replace "This is navigation only." with: INDEX is the agent retrieval catalog
  and a human navigation aid (INV-KNOWLEDGE-01; owned by
  `wrapper/contracts/schemas/context-index.yaml`).
- State that per the schema the catalog carries per-unit retrieval metadata
  (context id, summary, topics, aliases, domains, repositories, decisions,
  constraints, status, freshness, provenance) and that pending context proposals
  and stale-context warnings are indexed by target context id, related plan,
  topic, repository, and impact status.

## Acceptance action

Rewrite the opening line and add a short "retrieval metadata" note pointing to
the schema. The maintainer INDEX may stay lightweight, but its self-description
must not deny the catalog role.
