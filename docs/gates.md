# Human gates and cards

`wrapper/runtime/engine.sh` provides card formatting; this document owns the
human-facing effect language. A card is valid only for its current session,
action, target, observed state, and listed effects. A changed target or state
requires a new card.

Every consequential mutation names what will change and what will not. The
separate actions are identity acceptance, context acceptance, plan approval,
execution trigger, material scope change, status change, takeover, delivery,
archive/restore, publication/deployment/merge, and destructive cleanup.
These are separate explicit gates, not one reusable confirmation.

Approval changes plan status and task projections only. Execution starts only
from a separate explicit request. A verifier can produce completion evidence but
cannot finish a plan. Cleanup first reports dirty and unpushed work and needs an
additional discard confirmation before anything destructive.
