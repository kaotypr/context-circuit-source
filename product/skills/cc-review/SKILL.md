---
name: cc-review
description: Run an independent read-only code review of Context Circuit work when a person asks for one, and say so honestly when independence is unavailable.
---

# Independent review, on request

Independent verification is a manually requested read-only code review, usually
after pull-request creation and optionally after delivery as an audit. It is
never triggered by risk classification, and it is not a condition for opening a
pull request, delivering, or completing.

Use `.agents/skills/cc-dispatch/SKILL.md` to invoke an independent
review-capable agent or session when it is requested and available:

```sh
context-circuit-cli --workspace <root> --json agent dispatch --host codex --role reviewer --path . --review-requested --task 'Review the billing change against its success criteria'
```

`--review-requested` represents the actual user request. It supplies no human
approval by itself.

## What to supply

The requested diff, the current revision, relevant surrounding code, and the
intent's success criteria. A read-only reviewer cannot run Git on most hosts, so
the diff is text you hand it rather than something it fetches.

## What a reviewer returns, and does not do

It reports actionable findings with locations and its own limitations. It does
not modify code, automatically dispatch fixes, or post external comments unless
that was requested. Tests that change files belong to implementation, not to a
read-only review.

## When independence is unavailable

Explain that limitation and offer an ordinary review instead. Never call the
implementing session's own inspection an independent review — a fresh, separate
context is the whole point, and reusing the implementer's session produces
agreement rather than scrutiny.

The user decides whether to request fixes or proceed. This optional capability
does not make child agents mandatory anywhere else.
