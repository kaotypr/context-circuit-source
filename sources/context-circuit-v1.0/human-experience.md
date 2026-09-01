# Human experience — the conversational surface

v1.0 changes *what the human decides and when*, so it changes the conversation. This
file describes the human-facing surface: the two gates in plain language, what the
coordinator says when it re-gates, the transparency questions the human can ask, and
the reporting rules. The mechanics stay hidden; the decisions do not — the same
disclosure principle Context Circuit already holds.

## What stays hidden vs. shown

Hidden (as today): branches, worktrees, leases, base selection, candidate digests,
runtime records. Shown (as today, plus the new decisions): the intent and its
criteria, the adversary's findings, the worker's handoff, the verifier's result, and
the two gates. The reader never sees a workspace file name, internal path, or internal
id in normal conversation — internal terms are translated (see `glossary.md`), the
same rule publications already follow (INV-EXTERNAL-03 spirit).

## Gate 1 — approving an intent (new)

The coordinator presents the intent as one readable thing, folds in the adversary's
findings, and asks for a single decision:

> Here's what I understand you want to build: **retry failed checkout charges up to 3
> times before failing the order.** Out of scope: the payment provider, refunds.
> It'll touch the checkout service only.
>
> I had this checked adversarially first. One thing came back: the criteria said
> "retry up to 3 times" but didn't bound the added latency — an implementation could
> pass and still blow your latency budget. I've added a criterion for that.
>
> This is **Standard** risk (single service, reversible, well-covered), so it gets an
> independent check.
>
> Approve this and I'll build it — or tell me what to change.

Approval here is the real decision. No confirmation card, no token — a conversational
yes (INV-APPROVE-01 reworked). After approval the human is not asked to approve the
plan; the plan derives from what they just approved.

## Between the gates — mostly silent, by design

Execution, candidate, and verification run without asking the human anything (that is
the point — pains 3, 4, 5, 7). The human hears about them only as *reporting*, not as
*gates*:

> Built and checked. The independent check passed all four criteria against the exact
> change. Ready for you to try it and decide whether it ships.

The human tries it (they are the acceptance oracle for the human-observable part) and
accepts the candidate — a lightweight "looks right, ship it," not a status flip.

## Gate 2 — authorizing delivery (kept)

Unchanged from today (INV-DELIVER-01), and still the second real decision:

> Ready to open the pull request onto `main` for the checkout service. This merges the
> retry change. Go ahead?

## When it re-gates (the envelope check speaking)

The one place the human is pulled back in mid-flow is scope drift — and the coordinator
says *why*, in plain language:

> Heads up: building this, it turns out the change also needs to touch **payments-lib**,
> which wasn't part of what you approved. I've paused rather than widen the scope on my
> own. Want to include payments-lib (I'll re-check the criteria), or should I keep the
> change to the checkout service?

This is the crown-jewel envelope check surfacing as a question, not an error. It should
feel like the system protecting the human's approval, not nagging.

## Transparency — the questions a human can ask

Because v1.0 removes gates and leans on mechanical checks, the human must be able to
interrogate the machinery in words. The coordinator answers these from runtime facts
(no new gate, just reporting):

- **"Why is this Critical?"** → lists the risk signals that set the tier
  (`crown-jewels.md`), so tiering is never a black box.
- **"What's blocking the next plan?"** → names the unreconciled delivered work
  (knowledge debt, M4) in plain terms.
- **"What did the check actually run?"** → the acceptance criteria and their results
  bound to the current change, never the worker's prose.
- **"Is this verified?"** → honest by tier: "independently checked" at Standard/
  Critical; "human-supervised, not independently checked" at Explore. The word
  "verified" is never used for Explore output.

## Reporting rules (kept and extended)

- **Plain language, no internals.** Translate every internal term (`glossary.md`);
  never surface a file name, path, digest, or internal id in normal conversation.
- **Never overstate assurance.** Explore is "human-supervised"; a passed check is
  "the criteria passed," never "bug-free." Missing evidence is reported as missing,
  never as success (kept from today).
- **Say what happened faithfully.** A blocked or re-gated state is reported as such,
  with the reason and the human's options — not smoothed over.

## The net effect on felt ceremony

Today the human is asked to approve the plan, mark it done, and remember to reconcile —
three touch points per change that feel like bookkeeping. v1.0 replaces them with one
upstream decision that matters (the intent), an accept-and-ship moment, and a system
that refuses to *let* reconciliation be forgotten. Fewer prompts, and the prompts that
remain are real decisions — "that is how you get few questions rather than no
questions."
