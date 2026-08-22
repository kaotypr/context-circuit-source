---
schema_version: 2
id: WPE-002
plan: workspace-projection-and-gate-effects
status: ready
repository: context-circuit-source
paths:
  - wrapper/runtime/engine.sh
  - docs/gates.md
  - template/context/WORKSPACE.md
  - template/context/PROJECT.md
  - template/context/INDEX.md
  - test/gates/test-gates.sh
  - test/recovery/test-recovery.sh
depends_on: [WPE-001]
acceptance: [WPE-AC-01, WPE-AC-02, WPE-AC-03, WPE-AC-04, WPE-AC-05]
verification: [WPE-VT-02, WPE-VT-03, WPE-VT-04]
expected_evidence:
  - Atomic render and validation results for the identity region in all three summaries.
  - Identity-acceptance and repository-registration card templates with displayed defaults.
  - Existing cards updated with immediate and later effect sections.
  - Failure fixtures proving no partial identity-region publish and no post-confirmation field invention.
stop_conditions:
  - Rendering overwrites authored Product Knowledge outside the identity region.
  - A partial write can leave identity regions disagreeing with workspace.yaml.
  - Card wording changes authorization or combines separate gates.
---

# Render identity regions and present complete gate cards

## Objective

Implement the contract through host-neutral engine operations and canonical
gate cards while preserving authored Product Knowledge and separate human
gates.

## Work

Add bounded render and validation operations that stage and publish only the
identity region in the three summaries. Integrate validation into entry and
write preflight. Keep repository-local paths only in the ignored binding
file. Registering a logical repository must update `workspace.yaml` shared
identity and all three identity regions in one atomic operation.

Add full card templates in `docs/gates.md` for identity-acceptance and
repository-registration. Identity-acceptance must display proposed defaults
for mode, repositories or project items, roles, and default branches before
confirmation. Repository-registration is the shared-identity card consumed by
later clone or create-empty gates; this task does not add a new router action
or clone/create behavior. Update existing cards so each current confirmation
lists immediate effects and later authorized effects using the locked
identifiers. Confirmation records displayed defaults; it does not invent
fields afterward.

## Non-goals

Do not create a repository, activate create-empty, change the plan lifecycle,
combine approval with commit or execution, rewrite Product Knowledge, or
introduce provider-specific UI state.

## Verification

Use WPE-VT-02, WPE-VT-03, and WPE-VT-04.

## Expected evidence

Atomic-update fixtures, mismatch failures, identity-acceptance and
repository-registration cards, interruption tests, and a path audit showing
only declared mutable seed and wrapper files changed.

## Stop conditions

Stop if the engine rewrites content outside the identity region, publishes a
partial region, or any gate presentation implies authorization beyond the
current card.
