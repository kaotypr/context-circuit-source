# skill-and-contract

Continues [design.md](./design.md). This file specifies the **`cc-pair` skill**
surface, how it is triggered, how it relates to the existing skills, and the
**contract delta** v0.7.0 introduces. Mechanism lives in
[session-and-isolation.md](./session-and-isolation.md); this file does not restate
it.

## The `cc-pair` skill

`cc-pair` is one new product skill shipped at `.agents/skills/cc-pair/SKILL.md`,
resolved by path by the coordinator (INV-SKILL-01). It owns the direct interactive
collaboration **mode**: the user ↔ coordinator ↔ worker loop, base selection and
session-worktree bracketing, `host-blocked`, convergence, and the plain-language
narration of what changed. It grants no route, role, or authority — the coordinator
still writes nothing, and there is no verifier because the mode is not an execution.

What it does **not** own:

- **Repo-specific knowledge** — how to build/run the project, conventions — is
  discovered through the existing repository-grounding path and honored by the
  worker; it is not baked into the skill.
- **Any authority** — completion and delivery stay exactly where the core puts them;
  a pairing branch is delivered by the normal, separate `cc-deliver` action.

## Triggering

- **Anytime a repo/project is bound** — `cc-pair` is available to work directly on
  the project with no plan.
- **By intent or `/cc-pair`** — the user asks to just work on something together.
- **As a suggested next step after a plan or stack execution completes** — the
  coordinator may *offer* "refine this together before delivery?"; it is an option
  presented, never entered automatically.

The mode is scoped to interactive, human-supervised work; it is never used to
substitute for plan execution or independent verification.

## Relation to the existing skills

| Skill | Relationship to `cc-pair` |
| --- | --- |
| `cc-plan` / `cc-execute` / `cc-run-stack` / `cc-verify` | The plan lifecycle. `cc-pair` is **orthogonal** — it neither invokes nor is invoked by them. After they finish, `cc-pair` may be *offered*. |
| `cc-deliver` | Delivers a `cc-pair/<session>` branch like any other branch (source = the pairing branch, target = `anchor_branch`), as a separate explicit action. It reports the work as human-supervised, not verified. |
| `cc-complete` | Unaffected. `cc-pair` never marks a plan done. |

## Contract delta

v0.7.0 introduces a bounded, **plan-contract-free** change; owners settle in
`wrapper/contracts/`:

- **New skill `cc-pair`** at `.agents/skills/cc-pair/SKILL.md`.
- **New invariant `INV-PAIR-01`** (owner: `wrapper/contracts/invariants.yaml`) —
  the mode's boundary: outside the plan lifecycle, three actors, no verifier
  (human is the live oracle), no lease, own `cc-pair/<session>` branch + worktree
  from a base commit (never the active branch, never a plan branch in place),
  light resumable pointer only, output human-supervised and never "verified,"
  never auto-completes or delivers, worker commits follow INV-COMMIT-01,
  `host-blocked` when the worker child or worktree cannot be created.
- **New session-pointer schema** (owner: `wrapper/contracts/schemas/`, e.g.
  `pairing-session.yaml`) — the light pointer described in
  [session-and-isolation.md](./session-and-isolation.md) (repo, worktree, branch,
  base). A pointer, not an execution/verifier record.
- **New runtime folder `.runtime/pairing/<session>/`** — worktree/branch/pointer
  operations owned by `wrapper/runtime/engine.sh` (the coordinator invokes them;
  it never reads the engine implementation, INV-RUNTIME-01 corollary).
- **No change to `plan.yaml`, `verifier-result.yaml`, or any plan-lifecycle
  schema; no `schema_version` bump.** This is the payoff of decoupling.
- **`runtime_version` `0.7.0`.**
- **`host_evidence`.** The ability to create a worker child and a worktree is
  bounded host evidence that drives the `host-blocked` outcome; it authorizes no
  route (INV-HOST-01).

Owners-map additions (in `wrapper/contracts/invariants.yaml`): `pairing_mode`
→ `.agents/skills/cc-pair/SKILL.md`; `pairing_session` →
`wrapper/contracts/schemas/pairing-session.yaml`; plus the `INV-PAIR-01` entry.

## Compatibility

Opt-in and orthogonal. Absent any `cc-pair` use, everything is v0.6, and no plan
contract changes — a v0.6 reader is unaffected because the plan schema is untouched.

## Boundary

This scope drafts source only. It authorizes no implementation, delivery, or
publication, and reaches Product Knowledge and a plan through the normal path — the
coordinator gathers context from this named design source, proposes the invariant
and schema as contract changes through the existing proposal path, a human accepts,
and implementation follows. No separate design-acceptance gate.
