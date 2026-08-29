# skill-and-schema

Continues [design.md](./design.md). This file specifies the **`cc-refine` skill**
surface, how it relates to the existing skills, how it is triggered, and the
**contract delta** v0.7.0 introduces. Mechanism lives in
[refinement-loop.md](./refinement-loop.md) and
[acceptance-and-verification.md](./acceptance-and-verification.md); this file does
not restate it.

## The `cc-refine` skill

`cc-refine` is one new product skill shipped at `.agents/skills/cc-refine/SKILL.md`,
resolved by path by the coordinator (INV-SKILL-01). It owns the interactive
human-gated refinement **mode** and carries the **universal** UI-engineering
procedure (reading taste language into edits, responsive/a11y defaults, the
render-channel discipline, convergence on approval). It grants no route, role, or
authority.

What it does **not** own:

- **Repo-specific frontend knowledge** — component library, tokens, dev-server
  command, breakpoints, VRT tooling — is grounded per project through the existing
  repository-grounding path and reaches the worker via its brief, not baked into the
  skill.
- **Any authority** — completion, delivery, and the acceptance oracle stay exactly
  where the core puts them.

## Relation to the existing skills

`cc-refine` is a **sibling execution mode**, not a replacement for `cc-execute`:

| Skill | Role in a refinement arc |
| --- | --- |
| `cc-plan` | Produces the plan whose acceptance is the **human-gated visual** kind. |
| `cc-refine` | Runs the interactive loop (session-scoped worker + live render), holds one lease, converges on the human's approval, commits the blessed state. |
| `cc-verify` | The **single end verify** — recognizes the human-gated kind, verifies the residue, records the approval fact, freezes the VRT baseline. |
| `cc-deliver` | Separate delivery, normal gates, unchanged. |
| `cc-execute` | Unchanged for runnable-acceptance work; `cc-refine` is chosen instead when acceptance is taste-driven. |

The arc is one continuous conversation: plan lightly → refine interactively → the
system quietly verifies and freezes the baseline → deliver.

## Triggering

- **By intent** — the coordinator recognizes taste-driven refinement ("make it look
  good," "tighten the layout," "refine the UI") and enters/offers the mode.
- **By name** — `/cc-refine` explicitly, for users who want to invoke it directly.

The mode is scoped to interactive, human-gated work; it is not entered for runnable
acceptance.

## Contract delta

v0.7.0 introduces a bounded contract change; final owners and integers settle in
`wrapper/contracts/` on acceptance:

- **Acceptance kind (plan contract).** Acceptance criteria gain an optional
  **human-gated visual** kind (targets, breakpoints, objective residue, approval-fact
  reference). A plan without it is unchanged and behaves exactly as under v0.6; a
  plan using it is the next plan schema revision. Owner: the plan schema under
  `wrapper/contracts/schemas/`.
- **Provisional `INV-REFINE-01`.** The interactive refinement mode adds no
  authority; the terminal authority for a human-gated visual criterion is the human
  gate plus one independent end verify. No per-tweak independent verification; no
  headless "looks fine" (`host-blocked` instead).
- **Provisional `INV-REFINE-02`.** A human-gated visual acceptance freezes a
  visual-regression baseline on the human-blessed commit at the end verify; the
  verifier verifies the objective residue and the fact of approval, never the
  aesthetics.
- **`runtime_version`.** Provisional `0.7.0`; the template release takes a minor
  (pre-1.0) bump at publish.
- **`host_evidence`.** Gains no new authority; render capability and any
  render/verify block are reported as host evidence and drive the `host-blocked`
  outcome, they never authorize a route.

## Compatibility

Opt-in and backward compatible. Absent the acceptance kind, everything is v0.6. A
v0.6 engine that meets a plan carrying the new kind refuses it rather than running it
taste-blind (same posture as the run-stack `schema_version` guard). Template-upgrade
preservation follows the v0.5 convention.

## Boundary

This scope drafts source only. It authorizes no implementation, delivery, or
publication, and reaches Product Knowledge and a plan through the normal path — the
coordinator gathers context from this named design source, proposes the acceptance
kind and invariants as context/contract changes through the existing proposal path,
a human accepts, and a plan grounds in the result. No separate design-acceptance
gate.
