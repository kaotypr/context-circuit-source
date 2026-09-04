# The two crown-jewel checks, specified

v1.0 moves safety from "gate every transition" to "gate the two that matter," which
concentrates the entire safety property into two deterministic checks. This file
specifies them precisely, because everything rests on them. Both **fail upward**:
when the check cannot conclude "safe," it escalates rather than proceeds.

## Crown jewel 1 — the intent-envelope check

**Purpose.** Let per-plan approval collapse safely by guaranteeing that a plan
derived from an approved intent cannot exceed what the human approved.

**Checked at two points.** First, against **discovery's findings** — before any plan
is written, when discovery reports a manifest that reaches outside the approved
scope, the change is **held and re-gated** to the human, right there
(`discovery-and-grounding.md`). Second, against **each plan**, specified below.

**Verb.** `intent-envelope-check <plan>`, run as a `cc-plan` preflight and again as
an execution preflight (`execution-begin`).

**Inputs.**
- The plan's declared repositories and, per repository, its path regions
  (`cc_plan_repo_paths`, already computed today).
- The parent intent's `scope.repositories[*].{id, paths}` and its
  `contract_digest`.

**Algorithm.**
```
for each repository R the plan touches:
    if R is not named in intent.scope           -> EXCEEDS (new repository)
    for each path region P the plan declares in R:
        if no scope path S in R with region_overlap(S, P) covering P
                                                  -> EXCEEDS (path outside scope)
if plan.intent's current contract_digest != the digest approved
                                                  -> EXCEEDS (criteria changed)
if any comparison is indeterminate (unresolvable path, missing scope entry)
                                                  -> RE-GATE (fail upward)
otherwise                                         -> WITHIN
```
`region_overlap` reuses the **existing** lease-overlap semantics
(`cc_region_overlap`: equal, prefix-ancestor either way, or repository-wide `.`), so
envelope and lease share one definition of "inside." Coverage means the scope path is
an ancestor-or-equal of the plan's region (the scope must *contain* the plan region,
not merely touch it) — a plan region broader than any scope path EXCEEDS.

**Outcomes.**
- `WITHIN` → the plan is authorized by the intent's approval; no human step.
- `EXCEEDS` → **held and re-gated**: the human either widens the intent (a new
  decision, re-frozen digest — discovery re-runs for the widened scope) or narrows
  the plan.
- `RE-GATE` (indeterminate) → treated as EXCEEDS. Ambiguity never passes.

**What it does NOT catch (stated honestly).** The envelope guards *where* work
happens (repos/paths), not *what* it does. A plan that stays inside the paths yet
implements something the criteria did not anticipate is not caught here — that is the
job of the acceptance criteria themselves, sharpened by what discovery finds when it
reads the real code (`discovery-and-grounding.md`). Envelope + criteria together
cover where and what; neither alone is sufficient.

**Test fixtures (must-pass).**
- A plan touching a repo not in scope → EXCEEDS.
- A plan whose path region is broader than the scope path → EXCEEDS.
- A plan whose region is a strict child of a scope path → WITHIN.
- A criteria edit after approval → EXCEEDS until re-approved.
- An unresolvable/relative path in the plan → RE-GATE.

## Crown jewel 2 — consequence tier classification

**Purpose.** Decide how much assurance a change gets — including whether an
independent verifier spawns at all — deterministically, not by model mood.

**Where it is set.** On the intent's `contract.yaml` (`tier`), proposed at intent
time from transparent risk signals and confirmable/raisable by the human.

**Risk signals (the transparent inputs).** Each is a fact the coordinator already
has or can cheaply derive:

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

**Rules.**
- The classifier proposes a tier; the **default for anything uncertain or
  unrecognized is Standard** (independent verifier on).
- The human may **raise** the tier freely; **lowering** below the classified tier is
  an explicit, logged human decision (never silent, never the default).
- The runtime is model-blind: it never selects a tier. It only **enforces the floor**
  — `completion-ready` refuses a Critical candidate without a candidate-bound
  independent `passed`, and refuses to label an Explore result "verified"
  (INV-ASSURE-01, INV-RUNTIME-01).
- **Fail upward:** a signal set that does not clearly qualify for Explore is not
  Explore.

**Test fixtures (must-pass).**
- A scope path under an auth/secrets directory → Critical, verifier required.
- A single-repo, well-covered, reversible change → Explore eligible, no verifier.
- An unrecognized risk shape → Standard (not Explore).
- A human lowering Critical→Standard → allowed only as an explicit recorded
  decision; Critical→Explore for a security surface → refused.

## Why both must be the most-tested components

Everything else in v1.0 degrades gracefully; these two do not. A too-loose envelope
ships scope creep under an approval never given; a mis-low tier ships an unverified
Critical change. They are small and deterministic *on purpose* — so they can carry
the deepest fixtures and be audited in isolation, which is the payoff of keeping the
mechanics kernel small (`preserved-core.md`). If either proves hard to make reliable
in practice, the honest fallback is to keep more gates — the current
gate-every-transition model is the safe default this design bets against
(`risks-and-open-questions.md`).
