# Retired machinery

What v1 had, why each piece went, and what protects that concern now. This file
is what separates a rewrite from a regression: every removal below is a claim
that the concern still holds, by a different mechanism.

The governing principle is P6 — *ceremony must earn its cost.* A mechanism that
adds a step to every increment must prevent a failure that would otherwise
actually occur.

## The accounting

| v1 mechanism | What it did | Why retired | What covers it now |
| --- | --- | --- | --- |
| **Consequence tiers** (Explore / Standard / Critical) | Classified each change and scaled required assurance | Classification is a judgment dressed as a computation; it decided how much review a human would get without the human choosing, and a wrong classification was invisible | The human requests review when they want it; the agent recommends isolation for risky work and says what it left uncertain |
| **`contract_digest`** | Froze the identity of an intent's criteria at approval; downstream evidence self-invalidated on change | Solved drift nobody had: criteria change rarely, and when they do the agent re-approves. The digest's real effect was invalidating work after an innocuous wording edit | Re-approval only when the outcome or its success criteria **materially** change; the agent judges materiality and says so |
| **Candidate identity** | Bound evidence to an exact commit/base set so acceptance referred to a precise change | The Git revision already is that identity. A parallel identity had to be kept in sync with Git, and Git won every disagreement | `worktree inspect`, real branches, real diffs. P1: Git is authoritative |
| **Execution records** | Recorded that a plan ran, with what result | Recorded that the machinery ran, not that the work was right; on resume they were reliably staler than the diff | The diff, which is what happened; the worker's report reaches the person, and `record complete` carries the result into the plan when they ask for it |
| **Verification records** | Recorded a verifier's result against a candidate | Made an optional read-only opinion look like a certificate | The reviewer's findings, reported to the human, who decides |
| **Host-evidence records** | Captured host identity, permissions, model, provider status | Bounded evidence that never affected authority — so it was overhead that looked like provenance | Settings are a request; substitutions are reported. Nothing claims a host loaded or ran anything (P3) |
| **Path leases** | Let several writers share a repository by reserving paths | A lease is a lock over a namespace nobody could enumerate in advance; workers still collided at a granularity leases could not see | One worker per plan owns its whole worktree; `--shared` is an explicit, documented exception with a shared-ownership brief |
| **Automatic repair loops** (three-failure limit) | Re-dispatched work after a failed check | Turned a failure a human should see into three more attempts and a larger diff; the loop's own limit was evidence it was the wrong shape | Stop and report, preserving all work. The human decides whether to retry |
| **Risk-triggered review** | Spawned a verifier based on tier | Review that starts by itself becomes review nobody reads | Explicitly requested, read-only, fresh context, never blocking |
| **External publication (`cc-publish`)** | One-way publication of plans and discussions to external systems | A provider surface inside the core that the core had to keep neutral toward; used rarely, maintained always | An explicitly requested task through whatever host tools exist. No automatic provider synchronization |
| **Compulsory delegation** | Required spawning a child agent at certain tiers | Delegation is an efficiency choice, not a safety property; a required spawn for a two-line change bought nothing | Subagents stay central — `cc-dispatch`, four roles, a worker per plan in a stacked run. What went is the *requirement*: a small task may stay in the main session |
| **Plan approval gate** | A second approval between plan and execution | Asked a human to ratify a derivation they were not positioned to check — which is how ratification becomes a rubber stamp | Plans are earned by intent approval, presented, and proceed |
| **Inferred completion** | Derived done-ness from evidence | A plan is done when a human says the work landed; inferring it meant inferring it wrong sometimes, silently | Explicit `record complete`, which also stamps the date ordering reads |
| **The one-plan-per-repository *constraint*** | Forced every plan to name exactly one repository | Fragmented one coherent change across records that had to be reassembled by hand | Per-repository plans remain first-class and are often the right shape — `--repo` is repeatable, and `record order` reports which plans in a wave share a repository. What went is the *constraint*: a plan **may** span repositories when that reads better |
| **Invariant catalog (`INV-*`)** + contract schemas | A registry mapping each rule to one canonical owner | Real value for a large shell runtime with many parallel readers; disproportionate for one instruction file plus a typed Go package | The shared instruction states each rule once; Go's types and tests own the mechanical half |
| **Shell runtime (`engine.sh`)** | ~2,000 lines driving every state transition | Windows needed WSL; YAML edits were text manipulation; the engine was the ceiling on the product | A Go executable, ~2,500 lines across typed packages, tested on three operating systems |
| **Acceptance harness** | Scenario fixtures asserting product behavior end to end | Asserted that scripted scenarios passed, not that a host follows instructions — precisely the thing v2 declines to claim | Deterministic Go tests for file and Git behavior; `check-release.sh` for assembly; and an explicit statement that host behavior is unverified |

## One removal that was reversed: member bands

Bands were retired in the first v2 draft, on the argument that they put a
permanent accident of authorship into every public ID to solve a conflict that
is rare and resolvable by hand. The first half of that argument was right; the
second half was wrong. Two people allocating from separate clones is not rare in
the workspace this product is *for*, and "resolve it afterwards" means renaming
records and fixing every reference to them — the worst possible time to do it.

Bands are back, with the shape the original objection actually asked for:
**optional**. A member holding band N allocates intents from `N*100` and plans
from `N*1000`; a member without one allocates from the numbers no band claims. A
solo workspace never sees a band and still counts from `i001` and `p0001`, so the
accident of authorship is only in IDs belonging to workspaces that chose it in
exchange for offline safety.

What stayed retired is the *mandatory* band — v1 required one, which is how every
ID acquired a member's fingerprint whether or not anyone needed it.

This row is left in this file rather than deleted because the reversal is part of
the design record: P6 says ceremony must earn its cost, and the correct test is
whether the mechanism prevents a failure that would *actually occur*. Bands pass
that test; requiring them does not.

## Retiring a requirement is not retiring a capability

Two rows above name a *rule*, not a feature, and the distinction is the whole
point of those rows:

- **Subagents are not retired.** They are more central in v2 than in v1: four
  roles, per-host model and effort settings, a dedicated `cc-dispatch` skill, and
  a worker dispatched per plan in every stacked run. Only the *obligation* to
  spawn one is gone.
- **Per-repository plans are not retired.** They remain first-class and are
  frequently the right shape — separate plans deliver separately, depend on each
  other, and get their own worktree and branch. Only the *obligation* that every
  plan be exactly one repository is gone.

A design that removed either capability would be a regression. What v2 removed is
the compulsion, so the shape follows the work instead of the rule.

## What was kept, and why

The removals only make sense beside what survived unchanged in substance:

- **The knowledge circuit.** Not merely kept — tightened, with the durable-only
  boundary, catalog consistency, and reconciliation on explicit completion. It is
  the product; see [knowledge-circuit.md](./knowledge-circuit.md).
- **Intent before code, approved by a human.** The one gate that guards a
  decision only a person can make.
- **Delivery as a separate explicit authorization.** Unchanged.
- **Isolation by worktree, with unconditional preservation of work.** Extended
  with CoW reuse to make it cheap enough to actually use.
- **Read passive sources only when named.** Unchanged, and now enforced inside
  `context/` as well.
- **Four role shapes and real delegation.** Kept as capability descriptions,
  stripped of authority — and used throughout stacked execution.
- **Per-repository plans.** Kept as a first-class shape, no longer compulsory.
- **Member allocation bands.** Restored after being dropped, now optional rather
  than mandatory; see above.
- **No AI attribution in commits, PRs, reviews, or comments.** Unchanged, with
  the added obligation to inspect and remove injected attribution.

## The honest risk

Every removal trades a mechanical guarantee for a behavioral one. v1 could point
at a verification record; v2 points at an instruction that says *report what you
left uncertain*, and at a human who can read a diff.

That trade is defensible only if two things hold, and the design states both
rather than assuming them:

1. **The real artifacts are genuinely available.** Git state, diffs, check
   output, and the reviewer's findings are what the human inspects. They are, and
   the executable's job is to keep them reachable.
2. **The instruction is actually followed.** This is the unverified half (P3).
   No test in this repository establishes it, and the design says so in the
   release notes, in `WORKFLOW.md`, and here.

If the second assumption fails, v2 degrades to an agent with good bookkeeping and
no ceremony. v1's failure mode under the same assumption was an agent with good
bookkeeping, no ceremony, **and** a pile of records asserting otherwise. The
argument for v2 is not that the assumption is safe — it is that v2 does not
manufacture evidence against the possibility that it is false.
