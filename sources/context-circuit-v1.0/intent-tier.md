# Consequence tier classification

After the tracer reads the real code, the coordinator makes two decisions before it writes
any plan: **the feasibility check** (`intent-feasibility.md`) and **tier classification** (this
file). Tier is the **one safety-critical automated check** in v1.0's lighter model: it
decides how much assurance a change gets — including whether an independent verifier
spawns at all — so getting it wrong in the dangerous direction ships an unverified risky
change. It is the most-tested component in the system, and it **fails upward**: when
unsure, it tiers higher, never lower.

**Purpose.** Decide how much assurance a change gets, deterministically, not by model
mood.

## Where it is set

On the intent's `contract.yaml` (`tier`), proposed at intent time from transparent risk
signals and confirmable/raisable by the human, then re-checked with the tracer's findings
before plans are written.

## Risk signals (the transparent inputs)

Each is a fact the coordinator already has or can cheaply derive:

| Signal | Pushes toward |
| --- | --- |
| repositories touched > 1 | Standard+ |
| security / auth / secrets surface in scope paths | Critical |
| money / billing / payments surface | Critical |
| data migration or schema change | Critical |
| production-availability / deploy surface | Critical |
| irreversibility (no obvious revert) | Critical |
| novelty (no prior similar change in knowledge) | Standard+ |
| test coverage of the scope paths is thin/absent | Standard+ |
| single repo, reversible, well-covered, low novelty | Explore eligible |

## Rules

- The classifier proposes a tier; the **default for anything uncertain or unrecognized is
  Standard** (independent verifier on).
- The human may **raise** the tier freely; **lowering** below the classified tier is an
  explicit, logged human decision (never silent, never the default).
- The runtime is model-blind: it never selects a tier. It only **enforces the floor** —
  `completion-ready` refuses a Critical candidate without a candidate-bound independent
  `passed`, and refuses to label an Explore result "verified" (INV-ASSURE-01,
  INV-RUNTIME-01).
- **Fail upward:** a signal set that does not clearly qualify for Explore is not Explore.

## Plan edits do not recompute the tier

A human may hand-edit a plan through the coordinator after it is created; the tier is
**not** recomputed on such an edit. The coordinator gives a simple warning, and if the
edit introduces something risky, that is a human-error scenario, not a silent tier
downgrade — the runtime floor and Gate 2 still apply.

## Test fixtures (must-pass)

- A scope path under an auth/secrets directory → Critical, verifier required.
- A single-repo, well-covered, reversible change → Explore eligible, no verifier.
- An unrecognized risk shape → Standard (not Explore).
- A human lowering Critical→Standard → allowed only as an explicit recorded decision;
  Critical→Explore for a security surface → refused.

## Why tier must be the most-tested component

Tier classification carries the safety bet, so it does not get to degrade gracefully: a
mis-low tier ships an unverified Critical change — the single worst failure the system can
have. It is small and deterministic *on purpose* so it can carry the deepest fixtures and
be audited in isolation, which is the payoff of keeping the mechanics kernel small
(`preserved-core.md`).

The feasibility check (`intent-feasibility.md`) carries the *quality* bet — don't build the
impossible, don't guess — and scope-safety is carried by Gate 2, a human gate that
already exists. If tier classification proves hard to make reliable in practice, the
honest fallback is to keep more gates: the current gate-every-transition model is the safe
default this design bets against (`risks-and-open-questions.md`).
