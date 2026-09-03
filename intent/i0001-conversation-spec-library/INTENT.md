# Intent i0001 — author the conversation-spec library

_Status: draft (retroactive record, awaiting Gate 1 ratification). Tier: standard.
Scope: `agent-harness/conversations/`._

## The bigger picture

Context Circuit's agent test harness needed a single, versioned source of truth for
how the product is *expected to converse* with a lay user — authored whole for the
current version, not reconstructed from the v0.x→v1.0 deltas. This intent records the
decision to author that source: `agent-harness/conversations/` as a
**specification-by-example** library where each plot is a concrete expected dialogue
that *is* the specification, intended to generate the human-simulated harness cases.

## What a correct change achieved

- **Whole-surface coverage:** one plot per lay-user conversation across the entire
  lifecycle — orientation & context, the intent gate, the Explore tier and promotion,
  execute & verify, candidate/completion/delivery, the closed knowledge loop,
  organization, the external surfaces, and the whole-flow — each with a `## Spec` the
  grader can consume and a `## Dialogue` a person can read.
- **Generated worked cases:** a `generated/<id>.case.yaml` per plot in the existing
  conversation-case shape, marked as a build artifact derived from its plot.
- **Coverage cross-checks:** a per-phase doc mapping the existing scenario cases onto
  plots and recording the harness-capability gaps the authoring surfaced.
- **A flow guide:** a lifecycle-order reading lens over the catalog.
- **Descriptive discipline held:** every plot references INV-*/AC-* and defines no rule.

## Deliberately out of scope

Wiring the library to a generator or suite, regenerating the live scenario cases, and
applying the acceptance-criteria corrections it surfaced — all of that is the separate,
dependent decision **i0002**. No product behavior, runtime, or contract change; the
deterministic five-family suite (AC-36) stays the owner of the end-to-end families.

## Retroactive record

This intent is recorded **after** the work was authored and committed — the source
repository had no intents until now, so the library was authored intent-less. It is
captured here as the first-class decision that authorized it, so i0002 has a parent to
depend on and the history reads correctly. Because the artifacts already exist and are
inspectable, the spec adversary's *before-code* role does not apply; the criteria above
are instead directly and mechanically checkable against the committed library (see
`adversary.md`). Approving this ratifies what was built.

## Assurance

Standard tier: a single repository, reversible, source-only documentation change. It
carries no security, money, migration, production, or irreversibility signal.
