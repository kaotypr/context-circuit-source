# Migration and build order

A strangler-fig sequence, ordered by leverage-per-risk and grounded in the coupling
map. Each step is shippable on its own and additive; nothing requires the mechanics
port.

## Build order

1. **Intent + discovery + envelope check** (M1). Purely additive — a new
   `cc-intent` skill (drafts `INTENT.md` + `contract.yaml`, takes Gate-1
   approval), a new `cc-discover` skill (spawned on approval, fans out one
   read-only discovery child per repository, collects manifests, checks findings
   against the envelope), the `intent-*` verbs, a `plan.intent` field, and
   `intent-envelope-check` (`cc-plan` creates the plan(s) from the discovery
   manifest). Touches **no** execution mechanics. Delivers pains 1, 3, 4 and adds
   the single highest-leverage safety mechanism (the envelope check, now run
   against both discovery findings and plans). Lowest risk; do first.
   *Interim state:* plan approval can remain as-is until the envelope check is
   trusted, then flip INV-APPROVE-01 to auto-within-envelope.
2. **Candidate identity** (M2). Add `candidate-digest`, `candidate-current`,
   `human-acceptance-record`; bind `verifier-result-record` to the candidate.
   Additive to the two coupled execution drivers. Delivers mechanical staleness and
   sets up pains 6, 7. *Interim state:* candidate recorded and checked before it
   gates anything, so evidence-void behavior can be observed before it blocks.
3. **Tiered assurance** (M3). Add `tier` to the intent contract; make
   `completion-ready` tier-aware; spawn the verifier only at Standard/Critical;
   reframe `cc-pair` as the Explore tier with a promote step. Delivers pain 5 and
   resolves the pairing island. *Interim state:* default every intent to Standard
   (today's always-verify behavior) and enable Explore only once tiering is trusted.
4. **Closed knowledge loop** (M4). Emit the reconciliation-debt marker on
   completion/delivery; add `knowledge-debt` / `knowledge-reconciled`; make
   `cc-plan` grounding block (Standard/Critical) or warn (Explore) on debt. Delivers
   pain 8. *Interim state:* warn-only everywhere first, then escalate to block.
5. **Inferred completion** (M2/M3 finish). Flip completion to a projection of
   accepted + delivered at Explore/Standard; keep explicit at Critical. Delivers
   pain 7. Do last, because it most changes human-facing semantics and should ride
   on top of trusted candidate + tier machinery.
6. **(Optional, later, separate track)** behavior-preserving port of the mechanics
   to a typed language behind the same seam. Not required for M1–M5.

## Data migration

- **Existing plans** gain no parent intent retroactively. Two options, chosen per
  workspace: (a) leave completed/archived plans as-is (historical, no intent
  required); (b) for in-flight plans, author a lightweight intent that captures the
  already-agreed criteria and bind it. A one-time helper can scaffold an intent
  stub from an existing `plan.yaml` for the human to confirm — it never fabricates
  criteria silently.
- **Existing executions** keep their attempt-centered records; the candidate is
  computed from the same recorded commit maps + bases, so no record rewrite is
  needed — the candidate is an identity *over* existing data plus the (new) contract
  digest. Pre-intent executions use a null/`legacy` contract digest.
- **`context-impact` files** already exist; M4 adds the debt marker alongside them,
  so reconciliation history is preserved.
- Follows the release/upgrade preservation posture the product already has: additive
  schema versions, old readers degrade to read-only rather than guessing, and no
  upgrade makes stale evidence current under a new candidate algorithm.

## Rollback

Each step is behind its own verbs and invariant flags. Because the seam and the
mechanics are unchanged, any step can be disabled (revert the skill preflight, keep
the old gate) without touching leases, base selection, or delivery. The riskiest
semantic flips (auto-approval within envelope; inferred completion; debt blocking)
are each staged interim → enforced, so they can sit in warn/observe mode until
trusted.

## Testing hooks

The existing semantic acceptance suite and human-simulated scenarios extend
naturally: add cases for envelope drift (a plan that exceeds scope must re-gate),
candidate void (a new commit voids a prior pass and acceptance), tier floor (a
Critical candidate cannot complete without a candidate-bound independent pass), and
knowledge-debt blocking (grounding refuses while a delivered candidate is
unreconciled). The two crown-jewel checks (envelope, tiering) get the deepest
fixtures — see `risks-and-open-questions.md`.
