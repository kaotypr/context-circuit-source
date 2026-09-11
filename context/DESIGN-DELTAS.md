# Design ↔ implementation deltas

Where the shipped wrapper diverges from the v0.5 system design. Policy: the
context pages describe what ships today ("shipped is truth"); the divergences
are logged here as a design/implementation backlog and kept visible. Each
entry: what the design says, what ships, where seen, and disposition.

## Open

1. **Evidence layers.** Design (05, verifier prompt) requires an evidence layer
   per acceptance criterion, enforced by the verifier. The shipped runtime does
   not implement it (the `verification-evidence-layers` plan was deleted). Seen:
   design 05 §10–11 vs `.context-circuit/wrapper/contracts/schemas/verifier-result.yaml`,
   `.context-circuit/agents/verifier.md`. Disposition: the verification domain page describes the
   shipped verifier; revisit if/when evidence layers are implemented.

2. **Repair limit.** Design treats the limit as an authored `plan.yaml` field
   (`repair limit`). The shipped `INV-REPAIR-01` forces a stop at three and permits
   one additional attempt per explicit human continuation without resetting the count.
   Seen: design 05 §3/§13 vs `.context-circuit/wrapper/contracts/invariants.yaml`. Disposition:
   pages state the shipped hard limit of three.

3. **Plan-id allocation.** `INV-PLAN-03` (next after the highest active plan in
   the current member's band, with archived prefixes reusable) is a shipped
   invariant not stated in the original design, which says only "stable id". Seen:
   `.context-circuit/wrapper/contracts/invariants.yaml` vs design 05 §3. Disposition: pages keep
   the shipped band-scoped invariant (`INV-MEMBER-01`).

4. **Terminology authority.** The design makes the design chapter the owner of
   term meaning and translation, projected to `.context-circuit/docs/terminology.md`. The shipped
   `.context-circuit/docs/terminology.md` and `context/TERMINOLOGY.md` say authority is "settled
   by the runtime contracts under `.context-circuit/wrapper/contracts/`". Seen: design 08 vs
   `.context-circuit/docs/terminology.md`. Disposition: context follows the shipped projection;
   the ownership claim is the delta.

## Resolved

5. **"Two-stage router" language.** The design says the experience must not
   expose a two-stage router; the shipped wrapper has no router artifact. The old
   `ARCHITECTURE.md` router-as-canonical framing matched neither and was removed.
   Resolved 2026-08-24 in context.

6. **"wrapper" naming.** Resolved 2026-08-24 by decision: "wrapper" is an
   accepted synonym for the universal project workspace product. The design
   source (08, 09) was updated to accept it, and the glossaries add the term.
