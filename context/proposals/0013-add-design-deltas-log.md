---
id: 0013-add-design-deltas-log
target_context_unit: context/DESIGN-DELTAS.md
operation: add
statement: >
  Per the accepted "shipped is truth; log deltas" policy, add a context reference
  file that records where the shipped wrapper diverges from the v0.5 design, so
  the domain pages can describe shipped behavior while the divergences stay
  visible as an implementation/design backlog.
evidence_refs:
  - sources/context-circuit-v0.5-design/05-planning-and-execution.md
  - sources/context-circuit-v0.5-design/08-terminology.md
  - wrapper/contracts/invariants.yaml
  - docs/terminology.md
affected_repositories:
  - context-circuit-source
affected_commits:
  - 4b8ac0b
confidence: high
status: review-needed
---

# Add context/DESIGN-DELTAS.md (design <-> shipped divergences)

Create a reference log with these initial entries (each: design says / shipped
does / where seen / disposition):

1. **Evidence layers.** Design 05 §10-11 require an evidence layer per acceptance
   criterion, enforced by the verifier. Shipped runtime does not implement it
   (the `verification-evidence-layers` plan was deleted). Disposition: verification
   page describes the shipped verifier; revisit if/when implemented.
2. **Repair limit.** Design 05 §3/§13 treats the limit as an authored `plan.yaml`
   field. Shipped `INV-REPAIR-01` hard-codes a maximum of three. Disposition:
   pages state the shipped hard limit.
3. **Plan-id reuse.** `INV-PLAN-03` (never reused / next after highest ever) is a
   shipped invariant not stated in the design (design only says "stable id").
   Disposition: pages keep the shipped invariant.
4. **Terminology authority.** Design 08 §7 makes the design chapter the owner of
   term meaning/translation (projected to `docs/terminology.md`). Shipped
   `docs/terminology.md` says authority is "settled by the runtime contracts
   under wrapper/contracts/". Disposition: context follows the shipped projection;
   the ownership claim is a delta.
5. **"Two-stage router" language.** Design 08/§16 say the experience must not
   expose a two-stage router; the shipped wrapper has no router artifact. The old
   ARCHITECTURE.md router-canonical framing matched neither and is being removed
   (0001). Disposition: resolved in context.
6. **"wrapper" naming.** RESOLVED 2026-08-24 by decision: "wrapper" is an
   accepted synonym for the universal project workspace product. The design
   source (08 §1/§5, 09) was updated to accept it rather than deprecate it, and
   the glossaries add the term (0012). No longer an open delta; recorded here for
   history.

## Acceptance action

Create `context/DESIGN-DELTAS.md` with the entries above and add a row to
`context/INDEX.md` ("design deltas: context/DESIGN-DELTAS.md"). Keep it updated
as deltas are opened or closed.
