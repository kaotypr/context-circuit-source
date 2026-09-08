---
kind: domain
status: accepted
title: Assurance and consequence tiering
slug: assurance
owners: []
sources: []
source_revisions:
  - wrapper: HEAD
    commit: current
    basis: current-wrapper
generated_at: 2026-09-04T00:00:00Z
review_date: 2026-12-04
freshness: accepted-from-current-wrapper
assumptions:
  - Consequence tiering is the single safety-critical automated check in v1.0.
unknowns: []
contradictions: []
acceptance:
  state: accepted
  accepted_at: 2026-09-04
  accepted_by: maintainer
workflows:
  - .context-circuit/docs/planning.md
---

# Assurance and consequence tiering

## Summary

How much checking a change gets is summarized by a single **consequence tier** —
**Explore / Standard / Critical** — declared on the intent and enforced
mechanically (INV-ASSURE-01). Tiering is **the one safety-critical automated
check** in v1.0: with the spec adversary and the automated scope-envelope check
removed, the tier classifier and its floor are the only automated gate the runtime
enforces on the path between the two human gates. Route "how carefully is this
checked", tier selection/raising, and the verifier floor here. Owned by the
tiering rule in `invariants.yaml`.

## Scope

Inside: the three tiers and what each buys, the model-blind risk-signal
classifier, fail-upward behavior, the lowering guard, and the floor enforced at
approval and verifier spawn.

Outside: the intent that declares the tier ([intent and Gate 1](../intent/README.md)),
the tier *signal* the planner contributes ([tracing and feasibility](../tracing/README.md)),
the verifier mechanism itself ([verification](../verification/README.md)),
completion gating ([completion](../completion/README.md)), and the Explore-tier
pairing mechanics ([direct collaboration](../direct-collaboration/README.md)).

## Behavior

Assurance is composed from parts and **summarized by a consequence tier**
(INV-ASSURE-01), declared on the intent from transparent risk signals — repository
count, security/secrets, money, data migration, production/deploy, irreversibility,
novelty — and **raisable by the human**.

- **Explore** is human-supervised: **no independent verifier**, and its output is
  never labeled "verified". It is the bottom rung of the ladder; direct
  collaboration (`cc-pair`) *is* the Explore tier, not a separate mode, plus an
  explicit promote step that attaches an intent and raises the tier to turn a
  session into a candidate-bearing change without restart (INV-PAIR-01).
- **Standard** requires one **independent, read-only verifier** bound to the
  current candidate (INV-VERIFY-01).
- **Critical** requires the independent verifier **and** an explicit human
  completion (INV-COMPLETE-01) — verification alone never completes it.

**The runtime is model-blind** (INV-RUNTIME-01): it **never selects a tier**. It
provides a deterministic, **fail-upward** signal classifier over declared facts
and **enforces the floor**:

- `tier-classify` reads the intent's scope signals and returns a classified tier
  plus whether Explore is still permissible;
- **approval refuses the Explore tier** (which would drop the verifier) whenever a
  risk signal is present;
- **Standard and Critical execution** requires a candidate-bound independent pass,
  and the runtime refuses to treat an Explore result as verified.
  `completion-ready` reports that eligibility as a query; it does not refuse
  mark-done.

**Tiering fails upward: when unsure, tier higher.** A benign path that shares scope
with a risk surface still tiers to the risk surface; the highest signal across all
scope paths wins. The human may *raise* a tier freely, but **lowering onto a risk
surface is refused** (Critical → Explore for a security path fails; Critical →
Standard keeps the verifier and is allowed).

## Interfaces

- Human-facing: the intent's "How carefully this is checked" line (`Explore` / `Standard` / `Critical`)
- Engine (report/enforce only, never selects): `tier-classify`, `tier-lower-check`; floor enforced inside `intent-approve` and verifier spawn. `completion-ready` is an eligibility query, not a mark-done gate.
- Signals: repository count, security/secrets, money, data migration, production/deploy, irreversibility, novelty

## Data

The intent's `tier` field (provisional at draft, set at feasibility, frozen with
the contract at approval). The planner contributes a `tier_signal` in the
finding; the human's declared tier and any raise are recorded on the intent.

## Constraints and edge cases

A no-signal single-repo change classifies **Standard by default**, not Explore —
the engine cannot see reversibility, coverage, or novelty, so no-signal is not
auto-Explore; Explore stays available only as an explicit human lowering when no
risk signal is present. An unrecognized or repo-wide scope classifies Standard,
not Explore. More than one repository classifies Standard, not Explore. If the host
cannot create an independent verifier at Standard/Critical, the execution is
**blocked** — the worker or coordinator must not self-verify (INV-VERIFY-02).

## Implementation references

- `.context-circuit/wrapper/runtime/engine.sh`: `cc_tier_classify`, `cc_tier_signals`, `cc_tier_lower_check`, `cc_intent_scope_repos`, `cc_intent_scope_paths`; floor in `cc_intent_approve` and verifier spawn. `cc_completion_ready` remains an eligibility query.
- `.context-circuit/wrapper/contracts/invariants.yaml`: INV-ASSURE-01, INV-VERIFY-01, INV-VERIFY-02, INV-COMPLETE-01, INV-PAIR-01 (owner map: `assurance_tiering`)

## Verification

`sh test/acceptance.sh` (tier suite — the deepest fixtures, pinning fail-upward
corners and the per-tier verifier floor).

## Acceptance notes

Accepted 2026-09-04 (maintainer) to close a knowledge gap: the consequence-tier
ladder — the one safety-critical automated check — had no owning domain page
(only glancing mentions in `ARCHITECTURE.md` and plan-authorization).
In the earlier v1.0 design this was framed as "crown jewel 2" alongside the
scope-envelope check; the envelope check was removed, leaving consequence
tiering as the single safety-critical automated check.
