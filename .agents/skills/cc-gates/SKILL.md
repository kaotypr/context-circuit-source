---
name: cc-gates
description: Handle explicit human confirmation cards for lifecycle, recovery, delivery, archive, and cleanup actions.
---

Use when a consequential mutation is requested. Show action, target, observed
state, exact effects, non-effects, risks, and the current-session confirmation.
Bind confirmation to the displayed target and state; stale or vague “yes” is
not authorization. Apply archive dependency checks, takeover consequences,
offline fallback, delivery boundaries, and dirty cleanup inspection from the
canonical contracts. No child role may satisfy a human gate.
