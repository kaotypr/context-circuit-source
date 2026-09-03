---
id: 0034-add-conversation-spec-library
target_context_unit: context/domains/conversation-spec-library/README.md
operation: add
statement: >-
  The agent test harness is specified by example: agent-harness/conversations/
  is a versioned, whole-surface library of expected human↔coordinator dialogues
  ("plots", each a `## Spec` YAML plus a `## Dialogue`) that is the intended
  generative source for the human-simulated harness — each scenarios/*/case.yaml is a
  generated build artifact derived from a plot (generate, not reference), so a stale
  hand-edited case cannot linger. The library is descriptive only (it references
  INV-*/AC-*, defines no rule) and covers the whole lay-user conversational surface
  across every lifecycle phase, complementary to the deterministic end-to-end suite
  (test/scenarios/test-scenarios.sh, AC-36). As of this proposal it is FULLY AUTHORED
  (38 plots / 38 generated cases / 7 coverage docs) but PROPOSAL-STAGE: no generator or
  suite is wired yet and the 22 live scenarios/*/case.yaml remain untouched.
evidence_refs:
  - agent-harness/conversations/README.md
  - agent-harness/conversations/plots/
  - agent-harness/conversations/generated/
  - agent-harness/conversations/coverage/
  - agent-harness/scenarios/README.md
  - test/scenarios/test-scenarios.sh
  - test/acceptance/criteria-map.yaml
  - wrapper/contracts/invariants.yaml
affected_repositories:
  - context-circuit-source
confidence: high
status: review-needed
---

## Proposed knowledge

Add a bounded `conversation-spec-library` domain describing the specification-by-example
("conversation-spec-first") approach for the agent test harness, without duplicating
any rule it merely demonstrates.

The accepted page should record:

- **What it is.** `agent-harness/conversations/` holds one plot per expected
  conversation under `plots/<id>.md` (a machine-readable `## Spec` the grader can
  consume plus a human-readable `## Dialogue` annotated with `[decision_point: …]`),
  and a generated `generated/<id>.case.yaml` per plot in the existing conversation-case
  shape. The plot is the source; the case is a build artifact carrying a
  `# generated from plots/<id>.md — do not edit` header.

- **Why.** A delta is the wrong spec unit; the expected conversation is authored whole
  for the current version, so editing the plot first makes the harness fail against a
  live coordinator until the skills catch up — an old flow cannot silently persist
  across a version bump.

- **Two load-bearing disciplines.** Descriptive, never a second policy owner (a plot
  demonstrates behavior that `wrapper/contracts/invariants.yaml` and
  `test/acceptance/criteria-map.yaml` already own; "one rule, one owner" stays intact);
  and versioned, a real maintenance commitment (each plot names the `runtime_version`
  range it describes).

- **Scope and coverage.** The whole lay-user conversational surface across the six
  lifecycle phases — orientation & context, the intent gate, the Explore tier and
  promotion, execute & verify, candidate/completion/delivery, the closed knowledge loop
  and organization — plus the peripheral surfaces (external publication, system-design
  authoring) and the end-to-end whole-flow. It is the LIVE, whole-surface, lay-user
  complement to the deterministic five-family suite that owns AC-36.

- **Status and open work (do not overstate).** Fully authored and committed
  (source-only), but not yet wired to a generator or any suite; the 22 live
  `scenarios/*/case.yaml` are untouched and are the plots' intended future source.
  Authoring surfaced five acceptance-criterion mapping corrections in the existing
  cases (01, 16, 17, 19, 20) and a consolidated set of harness-capability gaps — chiefly
  id-agnostic post-conditions for fresh whole-arc cases, a few new seed states/faults,
  and new state predicates — recorded per phase under
  `agent-harness/conversations/coverage/`.

The page should reference the invariants and acceptance criteria the plots demonstrate
rather than restating them, and should not be treated as a lifecycle stage: the library
is source-only test material, not product behavior.
