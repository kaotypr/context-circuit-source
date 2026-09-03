# Intent i0001 — specification-by-example as the living spec of the agent harness

_Status: draft (re-gated — awaiting Gate 1). Tier: standard. Scope: `agent-harness/`._

## The bigger picture — one decision

The real decision is a single thing: **make specification-by-example the living
specification of the agent harness.** The plots under `agent-harness/conversations/`
are the expected human↔coordinator conversations, authored whole for the current
version. This intent takes them the whole distance — author, make them *generate* the
scenario cases, make the harness able to *run* them, run them against the live
coordinator, and let them *grow* without drifting. Authoring, wiring, and the grow-loop
are not separate decisions you would review and ship independently; the authoring only
has value because the plots become the live, self-revalidating source. So it is one
intent with several plans, where the completed authoring is simply plan #1.

## What a correct change achieves

- **The library exists, whole-surface** (done): one plot per lay-user conversation
  across the lifecycle, each with a `## Spec` the grader consumes and a readable
  `## Dialogue`, a generated worked case per plot, a coverage cross-check per phase, and
  a flow-order reading guide. Descriptive only — references INV-*/AC-*, defines no rule.
- **The plots generate the cases:** a deterministic generator that provably *derives*
  each case from its plot (fidelity for all plots, mutation-sensitive, byte-deterministic),
  never echoes a committed case.
- **The harness can run them:** the capability gaps the library surfaced are closed —
  first id-agnostic post-conditions (without which the flagship whole-arc cases cannot
  run at all), plus the new seed states, faults, and predicates — each *exercised* and
  *falsifiable*, never a no-op that always passes.
- **The existing cases are evolved:** the 22 are regenerated from their plots with seeds
  and budgets provably unchanged, assertion power non-regressed, the five caught
  acceptance-criteria corrections mechanically confirmed, and the suite green.
- **It runs live:** every frozen net-new case reaches a recorded grade against the live
  coordinator.
- **It grows safely:** a bidirectional drift guard gates the suite — a hand-edited case
  fails, and a plot edited without regeneration fails. Editing the plot first is the only
  way to change a case.

## The plans (one intent, several plans)

Authored under `cc-plan` after approval, each naming this intent, all inside `agent-harness/`:

1. **Author the whole-surface library** — *already done and committed* (plan #1).
2. Build the generator + fidelity/mutation/determinism checks.
3. Close the harness-capability gaps (id-agnostic post-conditions first) with negative fixtures.
4. Regenerate the 22 cases + mechanically confirm the five corrections.
5. Run the frozen net-new cases + add the bidirectional drift guard (the grow-loop).

## Deliberately out of scope

No product behavior, runtime, contract, or invariant change (the library stays
descriptive). The deterministic five-family suite (AC-36) stays the owner of the
end-to-end families. The 22 cases are evolved, not recreated. The file-rename / diagram
/ flow-metadata polish is deferred. Delivery (merge/push/publish) is a separate Gate 2
action. No typed-language runtime port.

## Assurance

Standard tier: a single repository, reversible, test-infrastructure change with real
novelty (a new generator and grader capabilities) — not Explore-eligible, so it gets an
independent verifier, but it carries no security, money, migration, production, or
irreversibility signal, so it is not Critical.

## History and dogfooding

This intent previously existed split as i0001 (authoring, retroactive) and i0002
(go-live). They were merged into this single decision because the seam was chronology
(part of it was already done), not decision structure — and the whole point is one
living specification, not two shipments. i0002 is removed as subsumed. This is itself
dogfooding: we use Context Circuit's own intent mechanism to find the best decision
shape for what will ship to the template.

Because merging changed the criteria of a previously-approved intent, i0001 is re-gated
to `draft`, its `contract_digest` cleared; a fresh spec-adversary pass runs on the
merged criteria (see `adversary.md`), and Gate 1 re-approval is yours to give. The
runtime was not run, so no id was machine-allocated and no digest is frozen.
