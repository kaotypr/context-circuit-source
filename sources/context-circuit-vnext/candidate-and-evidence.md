# Mechanism 2 — Candidate identity + mechanical evidence staleness

Fixes pain 6 (multiple verifiers for one pull request) and pain 7 (marking done
feels like bookkeeping). Makes "it passed earlier" impossible by construction
instead of by discipline.

## The problem it solves

Today evidence attaches to an **execution/attempt**
(`.runtime/executions/<plan>/<execution>/attempts/NNN/verifier.yaml`). Each plan
runs its own execution and spawns its own verifier. So N stacked plans that
converge to one pull request produce N verifications — and freshness ("does this
green still describe the current code?") is enforced conversationally and
per-attempt rather than as one mechanical fact. The unit of trust is the attempt;
it should be the change that actually lands.

## The object: a candidate

A **candidate** is a deterministic digest identifying the exact proposed result:

```
candidate = digest over {
  repository -> latest_commit    (the per-repo commit map, already recorded)
  base_commit per repository     (already recorded)
  contract_digest                (the frozen acceptance criteria from M1)
}
```

New verb: `candidate-digest <plan|change> <execution>` → prints a stable
`cand-<hash>`. Everything the engine already writes per repository
(`latest_commit`, `base_commit`) is the input; the only new ingredient is the
intent's `contract_digest`, which folds the *definition of correct* into the
identity of the *thing being proven*.

## The rule: evidence binds to a candidate, and any change voids it

New invariant **INV-CANDIDATE-01**: the independent-check result and the human
acceptance record both name the candidate they observed. A new commit **or** a
criteria change produces a new candidate, and all prior evidence and human
acceptance for the old candidate become **void** — not "stale, use judgment,"
void.

The engine already detects a commit change (`cc_verifier_result_record` compares
tips; `VERIFIER_MODIFIED_PRODUCT`). vNext generalizes that from "tips unchanged
since worker commit" to "candidate unchanged since evidence was recorded," and
extends it to human acceptance. This is the same idea three of the surveyed
designs reached independently — Bench's "greens go grey," Drydock's "not stale,
void," Lattice's freshness algebra — expressed in Context Circuit's own terms.

Because `contract_digest` is *inside* the candidate, hardening or changing the
criteria (M1) correctly re-invalidates proofs too — you cannot silently move the
goalposts and keep the old green.

## One pull request → one candidate → one verification (pain 6)

When several stacked plans converge to a single pull request against the base
branch, their combined tip set forms **one candidate**. The independent check runs
once against that candidate; the human accepts once. You stop paying for a verifier
per plan when the delivered thing is one change. Concretely, the run-stack loop
(which already composes plans onto integration bases via
`cc_base_prepare`/INV-CONCURRENCY-02) computes the candidate over the *combined*
delivered tip set and verifies that, rather than each plan's attempt in isolation.

The stacked-plan mechanics (leases, integration merge, base rebuild) are unchanged;
only the *unit that evidence and acceptance attach to* moves from the attempt to
the candidate.

## First-class human acceptance (pain 7)

New verb `human-acceptance-record <candidate>` writes a first-class acceptance bound
to the exact candidate and the checklist it was accepted against — adopting
Codex-PC's move. "Done" then becomes **inferred** from `candidate accepted +
delivered` rather than a separate status flip the human resents or forgets:

- At Explore/Standard tiers, completion is a projection of "this candidate was
  accepted and delivered," not a manual `plan-complete` call.
- At Critical tier, an explicit human completion is still required (M3), because
  the stakes justify a deliberate act.

For teams, the acceptance record is the durable, candidate-bound artifact that
tells everyone what was accepted and by whom — so the team-visibility that
`plan-complete` provided today is preserved, without the solo-user bookkeeping.

## Interaction with delivery drift (kept mechanic)

The existing drift guard (`cc_delivery_rebase`, INV-DELIVER-01) already rebases a
drifted branch and forces re-verification before a pull request opens. Under vNext
that is exactly right and gets sharper: a rebase changes the commit map, which
changes the candidate, which **voids** the prior verification and acceptance by
INV-CANDIDATE-01 — so the re-verification the drift guard demands is no longer a
special-case rule, it is the general candidate rule doing its job.

## Skills and engine

- **New verbs:** `candidate-digest`, `human-acceptance-record`,
  `candidate-current` (report the current candidate for a plan/change).
- **Changed:** `verifier-result-record` and `completion-ready` read and check the
  candidate id.
- **Skills:** `cc-execute` / `cc-run-stack` record evidence against the candidate
  (additive to their existing `.runtime/executions/` writes — the coupling they
  already have).

## What stays the same

The `.runtime/executions/` layout, the attempt/repair loop, the three-failure
counter, and the verifier's structural read-only enforcement are all unchanged.
The candidate is an *identity computed over records the engine already keeps* plus
the contract digest — it adds a binding, not a new state machine.
