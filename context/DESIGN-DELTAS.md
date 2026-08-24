# Design ↔ implementation deltas

Where the shipped wrapper diverges from the v0.5 design under
`sources/context-circuit-v0.5-design/`. Policy: the context pages describe what
ships today ("shipped is truth"); the divergences are logged here as a design/
implementation backlog and kept visible. Each entry: what the design says, what
ships, where seen, and disposition.

## Open

1. **Evidence layers.** Design (05, verifier prompt) requires an evidence layer
   per acceptance criterion, enforced by the verifier. The shipped runtime does
   not implement it (the `verification-evidence-layers` plan was deleted). Seen:
   design 05 §10–11 vs `wrapper/contracts/schemas/verifier-result.yaml`,
   `agents/verifier.md`. Disposition: the verification domain page describes the
   shipped verifier; revisit if/when evidence layers are implemented.

2. **Repair limit.** Design treats the limit as an authored `plan.yaml` field
   (`repair limit`). The shipped `INV-REPAIR-01` hard-codes a maximum of three.
   Seen: design 05 §3/§13 vs `wrapper/contracts/invariants.yaml`. Disposition:
   pages state the shipped hard limit of three.

3. **Plan-id reuse.** `INV-PLAN-03` (never reused / next after highest ever) is a
   shipped invariant not stated in the design, which says only "stable id". Seen:
   `wrapper/contracts/invariants.yaml` vs design 05 §3. Disposition: pages keep
   the shipped invariant.

4. **Terminology authority.** The design makes the design chapter the owner of
   term meaning and translation, projected to `docs/terminology.md`. The shipped
   `docs/terminology.md` and `context/TERMINOLOGY.md` say authority is "settled
   by the runtime contracts under `wrapper/contracts/`". Seen: design 08 vs
   `docs/terminology.md`. Disposition: context follows the shipped projection;
   the ownership claim is the delta.

## Resolved

5. **"Two-stage router" language.** The design says the experience must not
   expose a two-stage router; the shipped wrapper has no router artifact. The old
   `ARCHITECTURE.md` router-as-canonical framing matched neither and was removed.
   Resolved 2026-08-24 in context.

6. **"wrapper" naming.** Resolved 2026-08-24 by decision: "wrapper" is an
   accepted synonym for the universal project workspace product. The design
   source (08, 09) was updated to accept it, and the glossaries add the term.
