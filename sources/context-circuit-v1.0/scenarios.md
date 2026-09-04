# Worked scenarios

Four end-to-end walkthroughs make the design concrete and double as acceptance
fixtures. Internal artifacts are named for the reader; the human sees plain language
(`human-experience.md`).

## Scenario A — a Standard feature, one repository

Request: "Retry failed checkout charges up to 3 times before failing the order."

1. **Intent (Gate 1).** `cc-intent` drafts `INTENT.md` + `contract.yaml`, no code
   read: goal, non-goals (no provider change, no refund change), outcome-level
   criteria `ac-1` (transient → retried), `ac-2` (permanent decline → not
   retried), envelope scope `checkout-service:[src/checkout/]` (coarse), tier
   `standard` (provisional). The human approves; this also confirms the
   coordinator understood the ask. Approval **freezes** `contract_digest`.
2. **Discovery.** Approval spawns one read-only discovery child over
   `checkout-service`. It reads the real retry path and reports a manifest back to
   the coordinator: the exact call-sites, a proposed task split, and — the risk
   only a real-code read catches — "the retry loop as written would retry
   *synchronously*, which would blow the stated checkout latency budget; the fix
   needs an added-latency check." It turns `ac-1` and `ac-2` into executable
   `grep`/test checks, proposes `ac-3` (added latency within budget) to close the
   gap it found, and confirms the touched paths are `src/checkout/` and
   `test/checkout/` — within the approved envelope.
3. **Plan.** `cc-plan` creates the plan from the manifest, including the new
   `ac-3` check; `intent-envelope-check` confirms the plan touches only
   `src/checkout/` and `test/checkout/` → **WITHIN**. The human sees the plan at
   informal review — no second approval.
4. **Execute (human-triggered).** The human asks for execution; one worker child
   implements in an isolated worktree, commits.
5. **Candidate + verify.** `candidate-digest` = `cand-…` over the commit map + base +
   contract digest. The independent verifier (Standard) runs `ac-1..3` against the
   candidate → `passed`, bound to `cand-…`.
6. **Accept.** The human accepts the candidate (`human-acceptance.yaml`).
7. **Deliver (Gate 2).** The human authorizes the pull request onto `base_branch`.
8. **Reconcile.** Delivery emits a knowledge-debt marker; before the next plan
   grounds, `cc-plan` blocks until the human accepts/defers the reconciliation
   proposals. Done is inferred from accepted + delivered.

Human touch points: **two gates** (approve intent, authorize delivery) + trigger
execution + accept the candidate + resolve reconciliation. Compare to today: approve
plan, complete, deliver, each per plan.

## Scenario B — Explore → promote (the cc-pair path)

Request: "Quick — the checkout error banner shows the wrong copy."

1. **Explore.** No intent, no candidate. Coordinator + **one worker child** in a fresh
   `cc-pair/<session>` worktree; the human watches live. No verifier, and no
   discovery child — at Explore the human reads the real code alongside the agent,
   which is discovery enough for a change this small.
2. The worker fixes the copy; the human eyeballs it. Output is *human-supervised*,
   never "verified."
3. **"Actually, this touches the retry logic too."** The human **promotes** in place:
   `cc-intent` attaches a small contract (criteria for the copy + the retry note)
   and the human approves it; the tier rises to `standard`, which spawns a
   discovery child over `checkout-service` — it reads the pairing-branch code,
   reports back its manifest (including the executable checks), and `cc-plan`
   creates a plan of record from it. A candidate now exists; the independent
   verifier spawns and checks it. From here it is Scenario A from step 5.

No stop-and-restart: the same work becomes a candidate-bearing change. The cliff is a
ramp.

## Scenario C — two stacked plans, one pull request

Request produced two dependent plans (`0012` API change, `0013` consumer update),
both in `checkout-service`, to land together.

1. Each plan executes independently: own lease, own base (`0013` stacks on `0012`'s
   branch per INV-CONCURRENCY-02), own worktree, own worker. **Unchanged.**
2. At delivery the human groups `{0012, 0013}` into **one change set** (one PR).
3. `candidate-digest` is computed over the **integration tip** of the change set. The
   independent verifier runs **once** against that integration candidate — not twice.
   (Pain 6.)
4. One human acceptance, one delivery. If the integration tip will not build,
   `BASE_UNBUILDABLE` blocks the change set (not a worker failure); the human splits
   or reorders.

## Scenario D — envelope drift re-gates

During Scenario A, discovery turns out to also find a change needed in
`payments-lib` (a repository **not** in the approved intent scope).

1. The envelope check runs against discovery's findings and returns **EXCEEDS**
   (new repository) — before any plan is written.
2. Work is **held and re-gated**: the coordinator tells the human, in plain
   language, "this work now needs to touch payments-lib, which wasn't in what you
   approved — widen the intent or narrow the change?"
3. The human widens the intent scope (a new decision), which **re-freezes**
   `contract_digest`. Only after re-approval does `cc-discover` fold in
   `payments-lib` and `cc-plan` proceed. Scope creep never ships under the
   original approval.

## What these scenarios also test

A → the happy path and inferred completion. B → pairing-as-Explore and promotion. C →
change-set candidate and single verification. D → the crown-jewel envelope check
firing. Together they exercise every new mechanism and both crown jewels — the
starting set for the semantic acceptance suite (`migration-and-build-order.md`).
