---
kind: domain
status: accepted
title: Direct collaboration
slug: direct-collaboration
owners: []
sources: []
source_revisions:
  - wrapper: HEAD
    commit: 13a8e92
    basis: current-wrapper
generated_at: 2026-09-03T00:00:00Z
review_date: 2026-12-03
freshness: accepted-from-current-wrapper
assumptions:
  - Direct collaboration adds no authority; promotion is the only path that
    attaches an intent, a plan, a candidate, and a verifier.
unknowns: []
contradictions: []
acceptance:
  state: accepted
  accepted_at: 2026-09-03
  accepted_by: maintainer
workflows: []
---

# Direct collaboration

## Summary

`cc-pair` is the **Explore tier** of the one assurance ladder
(Explore / Standard / Critical, INV-ASSURE-01) — not a separate mode. It is the
bottom rung: a small, live, reversible change worked directly with the user in
one connected repository, with a coordinator and one worker and **no independent
verifier**. Route "let's work on this together", `/cc-pair`, and "just try it"
here, and offer it as an optional next step after a plan or stack execution.
Never enter it automatically. Owned by the `cc-pair` skill; its mechanics are
owned by `.context-circuit/wrapper/contracts/invariants.yaml` (INV-PAIR-01), and its place in the
ladder by INV-ASSURE-01.

The defining move is **promote in place**: when the work turns out to be real,
an intent, a plan, a candidate, and a verifier are attached without stopping and
restarting — the cliff is a ramp ([plan-authorization](../plan-authorization/README.md),
[verification](../verification/README.md)).

## Scope

Inside: entry and the one-repository rule, the coordinator↔worker live loop,
single-repository branch/worktree isolation, the light resumable pointer,
human-supervised (never "verified") labeling, convergence on explicit human
say-so, host-blocked behavior, the promote ramp into the checked pipeline, and
separate delivery.

Outside: the plan-execution loop and its verifier, lease, execution record, and
completion gate ([plan-execution](../plan-execution/README.md),
[verification](../verification/README.md), [completion](../completion/README.md)) —
none of which exist before promotion; and the tier ladder itself, owned by
INV-ASSURE-01.

## Behavior

- **The Explore tier, not a separate world (INV-PAIR-01, INV-ASSURE-01).** Before
  promotion, direct collaboration is an interactive working mode *outside* the
  plan lifecycle: it neither invokes nor is invoked by intent approval, execution,
  verification, completion, or delivery. It is nonetheless the Explore rung of the
  single ladder — the same ladder whose Standard and Critical tiers require an
  independent verifier. Explore is planless and has no grounding preflight of its
  own.
- **One repository, three actors.** A session resolves exactly one connected
  repository and cannot span repositories; ambiguity is one focused question
  before any state is created. The user is the live acceptance oracle, the
  coordinator interprets and delegates but never writes, and one worker (the
  existing `.context-circuit/agents/worker.md` role) makes the concrete change for each turn.
- **Isolation.** The worker writes only in a fresh `cc-pair/<session>` branch and
  isolated worktree under `.runtime/explore/<human-name>/`, created from a chosen
  base commit (the recorded base tip for
  fresh work, or a completed execution tip when offered after a plan/stack, or a
  commit the user names). The human chooses the short name; the agent does not
  invent it. Explore worktrees are not mixed with plan-execution worktrees under
  `.runtime/worktrees/`. It never reuses the active checkout or edits a plan
  execution branch in place.
- **No trust scaffolding.** There is no verifier, lease, execution record, failure
  counter, plan status, or completion gate; the only runtime state is a light
  resumable pairing-session pointer. Build/test checks are feedback for the live
  loop and are never reported as verified status.
- **Commits and convergence.** Worker commits occur only on explicit human request
  (INV-COMMIT-01). The session converges only when the user says the result is
  finished; a dirty worktree is never committed automatically. Closing a clean
  session preserves its branch and worktree — cleanup is a separate explicit
  action (`runtime-cleanup`) that also removes Explore worktrees under
  `.runtime/explore/`. Closing does not delete the worktree by itself, and a
  still-live session is not removed unless cleanup was explicitly asked. The
  result is described as **human-supervised, not independently
  verified**.
- **Promote — the ramp into the trust system.** When the work is real, promote in
  place rather than restart: attach an intent (`cc-intent`, so the human approves
  it at Gate 1 and tracers ground it against the real code on approval), raise the tier
  to Standard/Critical (the independent verifier appears — a risk surface refuses
  to stay Explore), and author a lightweight plan of record (`cc-plan`) whose
  execution binds to the existing pairing-branch commits and so produces the
  candidate. Everything downstream then runs on the standard scaffolding,
  unchanged. Un-promoted Explore work never produces a plan file or a candidate.

## Interfaces

- Human request: "let's work on this together" / `/cc-pair` / "try it and I'll tell
  you" — or an offer after a completed plan or stack execution.
- Runtime (invoke, never read the engine): `pair-begin . <repo> <session> [base]`,
  `pair-inspect . <session>`, `pair-close . <session>`, `runtime-cleanup .`.
- Promotion path: `cc-intent` → tier raise → `cc-plan` → the ordinary
  candidate/verify/deliver/reconcile pipeline.

## Data

A single light resumable pointer per session (repository, session slug, branch,
worktree, base commit, resumable flag). No leases, execution records, evidence,
handoffs, failure counters, or plan state are created.

## Constraints and edge cases

- **Delivery blocks, never silently rebases.** A closed pairing branch is
  delivered only by a separate explicit action ([delivery](../delivery/README.md)),
  labeled human-supervised. If the current base tip is not contained in the pairing
  branch at delivery time, delivery **blocks** and the work must be brought forward
  in a new human-supervised session — the plan delivery drift guard's automatic
  rebase-and-re-verify does *not* apply here (there is no verifier to re-run).
- **Host-blocked is read-only.** If the host cannot create the worker child or the
  isolated worktree, the outcome is `host-blocked` and the coordinator remains
  read-only; the coordinator never performs the worker's edits itself and there is
  no coordinator write fallback ([host-adapters](../host-adapters/README.md)).

## Implementation references

- `.agents/skills/cc-pair/SKILL.md`, `.context-circuit/agents/worker.md`
- `.context-circuit/wrapper/runtime/engine.sh`: `pair-begin`, `pair-inspect`, `pair-close`
- `.context-circuit/wrapper/contracts/invariants.yaml`: INV-PAIR-01 (Explore-tier mechanics),
  INV-ASSURE-01 (the tier ladder and promote step), INV-COMMIT-01
- `.context-circuit/wrapper/contracts/schemas/` (pairing-session shape)

## Verification

`sh test/acceptance.sh` — `test/pairing/test-pairing.sh` (pairing isolation,
resumability, close-only-clean, block-on-base-drift) and the tier suite that
enforces the Explore floor (no verifier; Explore never treated as verified;
approval refuses Explore when a risk signal is present).

## Provenance

Authored from the current wrapper at HEAD `13a8e92`, where INV-PAIR-01 and
INV-ASSURE-01 place direct collaboration as the Explore tier. The v0.7.0 design
source `sources/system-design/context-circuit/v0.7.0/direct-collaboration/` was
named by the original proposal and describes the standalone-mode ancestor of this
behavior; the shipped v1.0 invariants are the authoritative present-day evidence.

## Acceptance notes

Accepted 2026-09-03 from proposal `0029-add-direct-collaboration`, **rewritten to
the v1.0 assurance-ladder framing**. The original v0.7.0 proposal described
`cc-pair` as a mode "outside the plan lifecycle" and orthogonal to plan execution;
the accepted page instead records it as the Explore tier of the single ladder
(INV-ASSURE-01), with the promote step as the ramp into the checked pipeline,
because the shipped invariants supersede the standalone-mode framing. The
Explore-tier mechanics (one repository, user-as-oracle, no verifier/lease/record,
light resumable pointer, human-supervised labeling, block-on-base-drift,
host-blocked read-only) stand in full.
