# Approval-to-plan latency detail

This detail explains the intended shape of `i018-approval-plan-latency`. It is
additional design context for the intent, not a separate approval or lifecycle
stage.

Reading order:

1. [design.md](design.md) — normative overview and fixed decisions.
2. [incremental-grounding.md](incremental-grounding.md) — revision-bound reuse
   and safe fallback to full tracing.
3. [plan-shaped-trace.md](plan-shaped-trace.md) — trace output that is ready for
   plan derivation.
4. [atomic-materialization.md](atomic-materialization.md) — one-pass plan-stack
   creation, validation, and timing evidence.
