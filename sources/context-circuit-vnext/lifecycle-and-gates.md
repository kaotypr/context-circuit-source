# Lifecycle and gates

The resulting flow, with exactly two human gates. Compare to today's four
independent gates (approve, execute, complete, deliver).

## The flow

```
request
  │
  ▼
INTENT                     cc-intent
  ├─ author INTENT.md + contract.yaml (goal, non-goals, constraints,
  │                                    acceptance_criteria, scope, tier)
  ├─ spec adversary attacks the criteria           → adversary.md
  └─ ★ GATE 1: human approves the intent           → freezes contract_digest
  │
  ▼
PLAN(S)                    cc-plan
  ├─ derive tasks from approved contract + knowledge + grounding
  ├─ plan.intent = <id>
  └─ envelope check: within scope → authorized (no gate)
                     exceeds scope → re-gate to human (widen intent or narrow plan)
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
  in bounds*, after an independent adversary has tried to break the criteria. This
  is the decision they are already making today when they agree to a design — now
  it is first-class and it counts. Approval freezes `contract_digest`.
- **Gate 2 — Delivery.** The human authorizes *the irreversible act*. Unchanged
  from today (INV-DELIVER-01): separate from verification and completion, target is
  the recorded `anchor_branch`, blocks on missing remote/branch.

Everything between the gates is mechanical: envelope check, execution, candidate,
tiered verification, acceptance binding, drift rebase, reconciliation debt.

## What happened to today's four gates

| Today's gate | vNext |
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
         attach intent + criteria (adversary runs) → raise tier → spawn verifier
         → now a candidate exists → normal flow from EXECUTE onward
```

Promotion is the ramp that replaces today's cliff between pairing and plans.

## Why this is safe

The safety property is no longer "a human approved every transition." It is:

1. a human approved the **definition of correct** and the **scope envelope**
   (Gate 1), after adversarial challenge;
2. no plan or candidate may exceed that envelope without re-gating (the envelope
   check — crown jewel 1);
3. the appropriate independent check ran and is bound to the exact candidate (tier
   + M2), with tiering that fails upward (crown jewel 2);
4. a human authorized the irreversible act (Gate 2).

If both crown-jewel checks are correct, this is *stronger* on correctness than
today (because the criteria are now challenged and evidence cannot float free of
the code) and *lighter* on ceremony (because the derivative gates are gone).
