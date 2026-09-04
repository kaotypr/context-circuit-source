# Expected conversations

Reference dialogues for the human↔coordinator surface, one per scenario in
`scenarios.md`. They pin down the reporting language and double as fixtures for the
human-simulated test harness (the same way template test cases drive the coordinator
through natural turns). Rules from `human-experience.md` hold throughout: plain
language, no internal file names or paths, never overstate assurance, report failures
faithfully. "C" is the coordinator; "H" is the human.

## A — a Standard feature (single repository)

> **H:** Failed checkout charges should retry a few times before the order fails.
>
> **C:** Here's what I understand you want: **retry a failed card charge up to 3 times
> with backoff before marking the order failed.** Out of scope: the payment provider
> and refunds. It touches the checkout service only, and I'd call this **Standard**
> risk — one service, reversible, well-covered — so it'll get an independent check.
> Approve this and I'll get it grounded in the real code, or tell me what to change.
>
> **H:** Approve.
>
> **C:** Approved. Let me read the real checkout service against this before I plan
> anything.
>
> *(…discovery reads the checkout-service code…)*
>
> **C:** One thing came back: "retry up to 3 times" didn't bound the added latency —
> an implementation could pass every check and still blow your checkout latency
> budget. I've added a criterion that the added latency stays within budget, and
> broken the work into tasks from what it found. Ready for me to build it?
>
> **H:** Go ahead.
>
> **C:** Building it now.
>
> *(…work happens; the human is not asked to approve a plan…)*
>
> **C:** Built and independently checked — all four criteria pass against the exact
> change, including the latency one. Want to try it and decide whether it ships?
>
> **H:** Looks right. Ship it.
>
> **C:** Ready to open the pull request onto `main` for the checkout service. Go ahead?
>
> **H:** Yes.
>
> **C:** Opened. One thing before your next change here: this merged work hasn't been
> folded into what the project knows yet — want me to reconcile that now, or note it
> as no update needed?
>
> **H:** Reconcile it.

Human decisions: approve the intent, trigger the build, accept the result, authorize
delivery, resolve reconciliation — each a real decision, none a rubber stamp.

## B — Explore, then promote

> **H:** Quick one — the checkout error banner shows the wrong copy.
>
> **C:** I'll just work on this with you directly — a separate working copy, no formal
> plan or independent check, so I'll call it human-supervised, not verified. Here's the
> fix to the banner copy. Look right?
>
> **H:** Yeah. Actually — while you're there, the retry count in that path is wrong too.
> That's not just copy.
>
> **C:** That's real logic, not a quick tweak. I'd rather promote this to a proper
> change so it gets checked and remembered. I'll write up what we're doing — fix the
> banner copy and correct the retry count — and treat it as Standard so it gets an
> independent check. Sound good?
>
> **H:** Go.
>
> **C:** Promoted and approved. Let me read the real code in that path before anything
> else.
>
> *(…discovery reads the pairing-branch code…)*
>
> **C:** The work so far is captured, and I've turned what discovery found into a plan.
> Spinning up the independent check now. …Checked and passing. Try it and tell me if
> it ships.

The promotion is the moment an intent is approved; discovery then reads the code and
the first plan comes out of what it finds — before it, nothing was recorded beyond the
working copy.

## C — two stacked plans, one pull request

> **H:** The rate-limit change needs the API update and the consumer update together.
>
> **C:** Right — that's two pieces that land as one change: the API change first, then
> the consumer that depends on it. I'll plan both under the one intent.
>
> *(…intent approved; discovery reads both repositories; both are built…)*
>
> **C:** Both pieces are built. Because they ship as one pull request, I ran **one**
> independent check over the combined result, not one per piece — it passed. Accept and
> I'll open a single pull request for both.
>
> **H:** Accept.
>
> **C:** Before I do — heads up, the two pieces don't combine cleanly yet: there's a
> conflict where they meet. I've stopped rather than guess. Want to reorder them, or
> should I bring them back to you to resolve?

Shows the change-set candidate (one check for one pull request) and the honest
`BASE_UNBUILDABLE` outcome when the combined result will not build.

## D — envelope drift re-gates

> *(mid-build on Scenario A's approved intent)*
>
> **C:** Heads up: building this, it turns out the change also needs to touch the
> **payments library**, which wasn't part of what you approved — you'd scoped this to
> the checkout service only. I've paused rather than widen the scope on my own.
>
> Do you want to include the payments library — I'll re-check the criteria against the
> wider scope — or should I keep the change to the checkout service and find another way?
>
> **H:** Include it.
>
> **C:** Widening the intent to cover the payments library. That's a change to what you
> approved — want me to re-check the widened criteria against the real payments-lib
> code first, or just re-approve and continue?
>
> **H:** Just re-approve.

The re-gate reads as the system protecting the human's approval, not as an error or a
nag.

## E — knowledge debt blocks the next change

> **H:** New change: add a coupon field to checkout.
>
> **C:** Before I plan this — there's merged work in the checkout area I haven't folded
> into what the project knows yet (the retry change from earlier). If I plan on top of
> the old knowledge, this could be grounded on something already out of date. Want me to
> reconcile that first — it's quick — or explicitly mark it as no update needed?
>
> **H:** Reconcile first.
>
> **C:** Done — the retry behavior is now part of what the project knows. Now, the coupon
> field…

The debt gate refuses to let reconciliation be silently skipped; it never accepts
knowledge on its own.

## What these fixtures assert

Together they exercise the two gates (A), promotion from Explore (B), the change-set
candidate and honest non-atomic multi-repo outcome (C), the crown-jewel envelope
re-gate (D), and the closed knowledge loop (E) — the same coverage as `scenarios.md`,
expressed as the conversation a lay human would actually have. A human-simulated case
passes when the coordinator's turns match this language and its decisions land at the
same points.
