---
id: 0001-architecture-router-lease-execution
target_context_unit: context/ARCHITECTURE.md
operation: change
statement: >
  ARCHITECTURE.md frames "one two-stage router" as a canonical layer and says
  writers are isolated by "atomic plan leases". The shipped wrapper exposes no
  router artifact (routing is the coordinator + cc-* skills) and uses ownership
  locking (cc_lock_*, locks/<plan-id>/owner.yaml), not lease ceremony; the design
  (spine §8/§16) also says the experience must not expose a two-stage router.
  The four-layer framing also omits the core execution model.
evidence_refs:
  - agents/coordinator.md                       # routing owner; no router artifact
  - wrapper/runtime/engine.sh                    # cc_lock_acquire/owner/release
  - wrapper/contracts/invariants.yaml            # INV-OWN-01 (lock), INV-EXEC-02/03/04, INV-REPAIR-01
  - sources/context-circuit-v0.5-design/context-circuit-v0.5-design.md  # §8/§16 no exposed router
affected_repositories:
  - context-circuit-source
affected_commits:
  - 4b8ac0b
confidence: high
status: review-needed
---

# Fix ARCHITECTURE.md: router framing, lease term, missing execution model

## Corrections (design + shipped wrapper agree)

- Drop "one two-stage router ... The router and invariant catalogs are canonical."
  The shipped wrapper has no router file; conversational routing is owned by
  `agents/coordinator.md` and the `cc-*` skills. Describe routing as the
  coordinator selecting one bounded action, not a canonical router layer.
- Replace "atomic plan leases" with ownership locking: one active writer per
  plan enforced by an atomic exclusive-create lock (`cc_lock_*`), a live lock
  never silently stolen (INV-OWN-01). The word "lease" is legacy.
- Add the core execution model the current framing omits: one worker executes
  all tasks of one approved plan in dependency order across mapped repositories,
  commits each repository before verification, repairs with new commits, and
  stops at three worker failures (INV-EXEC-02/03/04, INV-REPAIR-01).

## Acceptance action

Rewrite the four-layer paragraph to: compact entry, coordinator-selected bounded
action, human-facing plan artifacts, isolated worker execution under an ownership
lock, independent read-only verification, and preserved runtime evidence.
