# Lifecycle and gates

The resulting flow, with exactly two human gates. Compare to today's four
independent gates (approve, execute, complete, deliver).

## The flow

```
request
  │
  ▼
INTENT                     cc-intent
  ├─ draft INTENT.md + contract.yaml, no code read (goal, non-goals, constraints,
  │                                    acceptance_criteria, scope, tier)
  └─ ★ GATE 1: human approves the intent           → freezes contract_digest
  │
  ▼
DISCOVERY                  cc-discover                (spawned automatically on approval)
  ├─ one read-only child per repository in the envelope, in parallel — no lease/worktree
  ├─ reports a manifest: file/call-site map, risks, task partition, executable
  │                       done-checks, tier signal, open questions, completeness proof
  ├─ kick-back: intent itself is wrong → re-opens the intent, nothing planned
  └─ envelope check against findings: within scope → continues
                                       exceeds scope → re-gate to human
  │
  ▼
PLAN(S)                    cc-plan
  ├─ create tasks from the discovery manifest
  ├─ plan.intent = <id>
  └─ envelope check against the plan: within scope → authorized (no gate)
                     exceeds scope → re-gate to human (widen intent or narrow plan)
  │
  ▼
REVIEW                     informal, no gate
  ├─ human may read the plan(s) and course-correct — a plan has no approval status
  └─ execution is a separate, human-triggered action
  │
  ▼
EXECUTE                    cc-execute / cc-run-stack   (kept mechanics)
  ├─ leases, base selection, worktrees, worker commits, repairs
  └─ tier decides verifier: Explore = none (human-supervised)
                            Standard/Critical = independent, candidate-bound
  │
  ▼
CANDIDATE = digest(commit map + bases + contract_digest)
  ├─ independent check result binds to candidate   (Standard/Critical)
  └─ human-acceptance-record binds to candidate
  │      any new commit or criteria change → new candidate → prior evidence VOID
  ▼
★ GATE 2: human authorizes DELIVERY (PR / merge / push / deploy)   (kept, INV-DELIVER-*)
  │      drift guard: rebase + re-verify → new candidate → re-accept
  ▼
DONE   (inferred from accepted + delivered at Explore/Standard;
        explicit human completion at Critical)
  │
  ▼
RECONCILE (M4): completion/delivery emits reconciliation debt;
                next plan's grounding blocks (or warns at Explore) until cleared
```

## The two gates, precisely

- **Gate 1 — Intent.** The human decides *what "correct" means* and *what scope is
  in bounds*, from the plain intent alone — no code has been read yet. Approving
  also confirms the coordinator understood the ask correctly, which is what lets
  discovery spawn immediately after, aimed at a confirmed target. This is the
  decision they are already making today when they agree to a design — now it is
  first-class and it counts. Approval freezes `contract_digest`.
- **Gate 2 — Delivery.** The human authorizes *the irreversible act*. Unchanged
  from today (INV-DELIVER-01): separate from verification and completion, target is
  the recorded `base_branch`, blocks on missing remote/branch.

Everything between the gates is mechanical: discovery, envelope check, execution,
candidate, tiered verification, acceptance binding, drift rebase, reconciliation
debt.

## What happened to today's four gates

| Today's gate | v1.0 |
| --- | --- |
| **Approve** (draft→approved, INV-APPROVE-01) | moved up to **Gate 1 (intent)**; plan approval is automatic within the envelope |
| **Execute** (separate authorization, INV-EXEC-01) | authorized by the approved intent's envelope; no separate human step within scope |
| **Complete** (explicit human, INV-COMPLETE-01) | **inferred** from candidate acceptance + delivery at Explore/Standard; explicit only at Critical |
| **Deliver** (INV-DELIVER-01) | **kept — Gate 2** |

Two human gates instead of four, and — crucially — the surviving upstream gate is
on the decision the human actually makes, not on a derivative they rubber-stamp.

## The Explore fast path (pairing)

```
request → Explore session (cc-pair as tier)      no intent, no candidate, no verifier
  │        user + coordinator + one worker, live human oracle, git-isolated
  │
  ├─ done here → human-supervised output, never "verified"   (INV-PAIR-01 preserved)
  │
  └─ "this is real" → PROMOTE in place:
         attach intent + criteria → raise tier → discovery spawns and reads the code
         → plan created from the manifest → spawn verifier
         → now a candidate exists → normal flow from EXECUTE onward
```

Promotion is the ramp that replaces today's cliff between pairing and plans.

## When each artifact is created

A plan is **not** the thing the human approves — it is the derivation of an
already-approved intent. So artifacts appear in this order:

| Stage | Artifact created |
| --- | --- |
| Request → intent (`cc-intent`) | `intent/<id>/` — `INTENT.md`, `contract.yaml`. **No plan yet.** |
| Gate 1 — approve intent | `contract_digest` frozen. Discovery spawns. Still no plan. |
| Discovery (`cc-discover`) reports | a durable manifest recorded as grounding evidence — file/call-site map, risks, executable done-checks, tier signal. **Still no plan** until the envelope check against these findings passes. |
| `cc-plan` creates the plan | `plans/<id>/` — `PLAN.md`, `plan.yaml` (with `intent: <id>`), `tasks/`. **This is when a plan file is created**, only after the envelope check against the plan passes. |
| Execute → candidate → accept → deliver | runtime records under `.runtime/…` (candidate, acceptance, debt marker). |

**A plan file is created only at Standard/Critical, after intent approval and
discovery.** It is the mechanical elaboration of the approved intent, grounded in the
discovery manifest + Product Knowledge, and it must stay within the intent's scope
envelope. Creation is **automatic** — the coordinator's next action, not a gate — and
one intent may yield **one or more** stacked plans (`intent-and-criteria.md`).

## Explore creates no plan — and the promotion decision

**Explore (cc-pair) writes no intent, no plan, and no candidate.** Its only state is
the light pairing-session pointer plus the `cc-pair/<session>` branch and worktree. It
is recordless by design; that is what keeps the fast path fast.

**Promotion is the first moment an intent and a plan exist for that work.** When the
human promotes an Explore session, `cc-intent` authors the intent and the tier rises,
and to enter the trust pipeline (candidate → verifier → delivery → reconciliation) the
work needs the normal plan/execution scaffolding.

There is a real design decision here, stated so it is not silently assumed:

- **Recommended — promotion authors a lightweight plan of record.** `cc-plan` creates a
  `plans/<id>/` whose tasks describe the change already made in the session; its
  execution binds to the existing pairing-branch commits, which produces the candidate.
  Everything downstream then runs on the standard scaffolding, unchanged. The cost is a
  plan authored partly retroactively.
- **Alternative — a plan-less candidate.** Allow an execution/candidate to be formed
  directly from the pairing branch + the intent, with no `plans/<id>`. This is lighter
  but makes the candidate the one execution kind with no plan behind it, which
  complicates the otherwise-uniform "every candidate has a plan" rule.

This study takes the first option (a plan is always the home of an execution), and flags
the choice as an open question in `risks-and-open-questions.md`. Either way, the
invariant that matters is unchanged: **no plan file is ever created for un-promoted
Explore work.**

## Why this is safe

The safety property is no longer "a human approved every transition." It is:

1. a human approved the **definition of correct** and the **scope envelope**
   (Gate 1) — and that approval doubles as confirmation the coordinator understood
   the ask;
2. discovery then reads the real code and no plan or candidate may exceed that
   envelope without re-gating (the envelope check, run against both discovery's
   findings and each plan — crown jewel 1);
3. the appropriate independent check ran and is bound to the exact candidate (tier
   + M2), with tiering that fails upward (crown jewel 2);
4. a human authorized the irreversible act (Gate 2).

If both crown-jewel checks are correct, this is *stronger* on correctness than
today (because the criteria are now first-class and reviewable, discovery grounds
the plan in the real code before anything executes, and evidence cannot float free
of the code) and *lighter* on ceremony (because the derivative gates are gone).
