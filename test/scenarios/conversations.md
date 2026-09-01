# v1.0 human-simulated conversations (fixture)

Expected human↔coordinator dialogue for each v1.0 scenario. "H" is the human, "C"
the coordinator. These pin the reporting language: plain words, no internal file
names, paths, digests, or runtime terms; the two gates land as real decisions;
assurance is never overstated (Explore is human-supervised, never "verified"). The
scenario driver (`test-scenarios.sh`) plays the matching engine flow and asserts
these lines carry no internal vocabulary.

## A — a Standard feature (single repository)

> **H:** Failed checkout charges should retry a few times before the order fails.
>
> **C:** Here's what I understand you want: retry a failed card charge up to 3
> times with backoff before marking the order failed. Out of scope: the payment
> provider and refunds. It touches the checkout service only.
>
> I had this checked adversarially first. One thing came back: "retry up to 3
> times" didn't bound the added latency — an implementation could pass every check
> and still blow your checkout latency budget. I've added a criterion that the
> added latency stays within budget.
>
> This is Standard risk — one service, reversible, well-covered — so it gets an
> independent check. Approve this and I'll build it, or tell me what to change.
>
> **H:** Approve.
>
> **C:** Building it now.
>
> **C:** Built and independently checked — all four criteria pass against the exact
> change, including the latency one. Want to try it and decide whether it ships?
>
> **H:** Looks right. Ship it.
>
> **C:** Ready to open the pull request onto main for the checkout service. Go ahead?
>
> **H:** Yes.
>
> **C:** Opened. One thing before your next change here: this merged work hasn't
> been folded into what the project knows yet — want me to reconcile that now, or
> note it as no update needed?
>
> **H:** Reconcile it.

## B — Explore, then promote

> **H:** Quick one — the checkout error banner shows the wrong copy.
>
> **C:** I'll just work on this with you directly — a separate working copy, no
> formal plan or independent check, so I'll call it human-supervised, not verified.
> Here's the fix to the banner copy. Look right?
>
> **H:** Yeah. Actually — while you're there, the retry count in that path is wrong
> too. That's not just copy.
>
> **C:** That's real logic, not a quick tweak. I'd rather promote this to a proper
> change so it gets checked and remembered. I'll write up what we're doing, check
> the criteria adversarially, and treat it as Standard so it gets an independent
> check. Sound good?
>
> **H:** Go.
>
> **C:** Promoted. The work so far is captured, and I've spun up the independent
> check on it. Checked and passing. Try it and tell me if it ships.

## C — two stacked plans, one pull request

> **H:** The rate-limit change needs the API update and the consumer update together.
>
> **C:** Right — that's two pieces that land as one change: the API change first,
> then the consumer that depends on it. I'll plan both under the one intent.
>
> **C:** Both pieces are built. Because they ship as one pull request, I ran one
> independent check over the combined result, not one per piece — it passed. Accept
> and I'll open a single pull request for both.
>
> **H:** Accept.

## D — envelope drift re-gates

> **C:** Heads up: building this, it turns out the change also needs to touch the
> payments library, which wasn't part of what you approved — you'd scoped this to
> the checkout service only. I've paused rather than widen the scope on my own.
>
> Do you want to include the payments library — I'll re-check the criteria against
> the wider scope — or should I keep the change to the checkout service?
>
> **H:** Include it.
>
> **C:** Widening what you approved to cover the payments library. That's a change
> to the criteria, so I re-ran the adversarial check — nothing new came back.
> Re-approve the wider scope and I'll continue.
>
> **H:** Approved.

## E — knowledge debt blocks the next change

> **H:** New change: add a coupon field to checkout.
>
> **C:** Before I plan this — there's merged work in the checkout area I haven't
> folded into what the project knows yet (the retry change from earlier). If I plan
> on top of the old knowledge, this could be grounded on something already out of
> date. Want me to reconcile that first — it's quick — or explicitly mark it as no
> update needed?
>
> **H:** Reconcile first.
>
> **C:** Done — the retry behavior is now part of what the project knows. Now, the
> coupon field.
