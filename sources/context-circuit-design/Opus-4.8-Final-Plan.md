# Context Circuit — Final Design (Opus 4.8)

The complete target design for the redesigned Context Circuit: how it is
structured, which files and folders matter, how an agent behaves in every case,
how a human uses it in every case, and how plans are reviewed, approved, and
triggered. This is the reference specification — not itself always-read context.

> Design principles carried from the plan rounds: **less context** (two tiers +
> route-scoped reads), **one clear route** (a single decision model), **one owner
> per rule** (invariant IDs), **same safety** (nothing weakens a human gate), and
> **proof over proxy** (real-token and routing-clarity measurement). Every heavy
> mechanism must pay rent before it ships.

---

## Part 1 — The mental model

Context Circuit turns a directory into an **AI-agent workspace**. A human drives
complex, multi-session software work in plain language; the agent orients,
plans, executes in isolated worktrees, verifies, and hands off — and **the human
holds every consequential gate.**

The whole system is **instruction- and filesystem-driven**. There is no CLI, no
daemon, no runtime pointer. Files are the coordination surface; the human speaks
in sentences.

Five things are kept strictly distinct (the vocabulary the whole design uses):

| Term | Meaning | Lives in |
| --- | --- | --- |
| **Host task** | The chat/thread the human opened in their tool | the host (Codex, Claude Code, …) |
| **Session** | One AI execution context (root or child) | `.runtime/sessions/<id>/` |
| **Plan** | Human-reviewed *intended work* | `plans/<repo>-plans/<n>-<slug>/` |
| **Task** | A bounded unit inside a plan | inside the plan bundle |
| **Worktree** | The writable Git isolation for one plan execution | `.runtime/worktrees/<repo>/<plan>/` |

A **plan** is the unit a human thinks in. A **session** is how an agent executes
it. They are never conflated. Two more terms every persisted record uses
precisely: **Implemented** = runtime evidence (writer done, verifier passed,
clean committed worktree) while the plan is still `approved`; **Done** = the
human-confirmed canonical status. Implemented never means Done.

---

## Part 2 — Repository and workspace structure

Context Circuit is shipped as a **template**. The **source repo** (this repo) is
the product; a **published artifact** is the blank workspace users instantiate.
The redesign makes that split explicit so an agent always knows which it is in.

### 2.1 The source repository (what maintainers work in)

```text
context-circuit/                 ← the product source
├── README.md                    ★ HUMAN front door — "how to use this in 1 page"
├── AGENTS.md                    ◆ Tier-0: safety, precedence, source boundary
├── WORKFLOW.md                  ◆ Tier-0: lifecycle + the one routing table
├── workspace.yaml               ◆ Tier-0: config (repos, mode, branches, gates)
│
├── wrapper/                     the shipped, versioned wrapper (immutable in an instance)
│   ├── manifest.yaml            wrapper_version, catalog, budgets, ship-vs-maintainer map
│   └── context-sets.yaml        per-route read sets + byte budgets (§4)
│
├── docs/                        Tier-1 reference (loaded only when a route needs it)
│   ├── invariants.md            ★ one owner per rule (ID → owner → canonical text)
│   ├── routing.md               the route decision model in detail
│   ├── runtime-contract.md      record shapes, leases, handoffs, wrapper_version
│   ├── planning.md              plan/task lifecycle, archive, stacks
│   ├── product-knowledge.md     how context/ is grounded
│   └── … (configuration, delivery, integrations, release, getting-started)
│
├── .agents/skills/             host-discoverable capability contracts (thin)
│   └── cc-*/SKILL.md            trigger + unique procedure only; reference invariants
├── agents/                     role deltas: coordinator, repository-worker, reviewer
│
├── context/                    ★ Product Knowledge — describes Context Circuit ITSELF
│   ├── INDEX.md                ◆ Tier-0: navigation only (concept → owner file)
│   ├── WORKSPACE.md            workspace identity
│   ├── PROJECT.md              the project (here: Context Circuit the product)
│   ├── ARCHITECTURE.md · CONVENTIONS.md · DECISIONS.md
│   ├── domains/ · roles/       deep, scope-loaded knowledge
│   └── sources.yaml            provenance
│
├── template/                   the blank seed copied into a new instance
│   ├── workspace.yaml          uninitialized-workspace defaults
│   └── context/                starter (uninitialized) knowledge pages
│
├── sources/                    passive human inbox (never auto-scanned)
├── plans/                      human-reviewed intended work (see Part 6)
├── test/                       acceptance + budget + token-truth + routing harnesses
└── scripts/                    release staging (maintainer-only, never shipped)
```

Legend: **★ human reads first**, **◆ Tier-0 always-read spine**.

### 2.2 The published artifact (what a user gets)

```text
my-workspace/
├── README.md · AGENTS.md · WORKFLOW.md   thin entry adapters (host-wired)
├── wrapper/                              immutable versioned contract
├── .agents/skills/ · agents/            capability + role adapters
├── workspace.yaml · context/            copied from template/ (uninitialized)
├── sources/ · plans/                    empty, ready for the user
└── .runtime/                            created on first session (private)
```

The release build asserts the artifact starts **uninitialized**, is **clean**,
and never contains `.runtime/`, maintainer-only paths, or credentials.

### 2.3 The files that matter — by audience

**A human only ever needs these:**

| File | Why a human opens it |
| --- | --- |
| `README.md` | The one-page "how do I use this" |
| `plans/<repo>-plans/<n>-<slug>/PLAN.md` | The human-readable plan: what, why, scope, status, how to approve |
| the **"what's next" view** (from `cc-whats-next`) | The live dashboard of what they can do right now |
| a plan's `handoff` (surfaced in chat) | What the agent did, found, and needs decided |

Everything else (`plan.yaml`, `.runtime/*`, `wrapper/*`, `docs/*`) is machine or
reference material the human may inspect but is never required to.

**An agent always reads the Tier-0 spine (◆), then only what its route selects.**

---

## Part 3 — The layered authority model

When two documents seem to conflict, this order decides — and the agent surfaces
the conflict rather than guessing.

| # | Authority | Owner |
| --- | --- | --- |
| 1 | Host & system instructions | host |
| 2 | Wrapper safety | `AGENTS.md` |
| 3 | Workflow rules & routing | `WORKFLOW.md` |
| 4 | Repository-local instructions (code) | each repo |
| 5 | Human-approved plans & decisions | `plans/`, `context/DECISIONS.md` |
| 6 | Source-cited Product Knowledge | `context/` |
| 7 | Runtime state & handoffs | `.runtime/` |
| 8 | Agent assumptions | (lowest) |

**Runtime never overrides an approved plan, a human gate, or Product Knowledge.**
Every shared rule has exactly one home in `docs/invariants.md`; all other files
reference its ID instead of restating it. Example registry line:

```text
INV-OWN-01  owner: WORKFLOW.md   One active writing owner per plan; one writable session per worktree.
INV-GATE-03 owner: AGENTS.md     A subagent never satisfies a human gate.
INV-RUN-07  owner: docs/runtime-contract.md  No global current-session pointer; sessions are discovered by record.
```

---

## Part 4 — How an agent loads context (two tiers + route sets)

```text
Human intent + host identity
        ↓
Read Tier-0 spine  (≤ 8 KB, always)
        ↓
Bounded state probe  (active sessions, plans, leases — metadata, not full reads)
        ↓
One route decision  (Part 5)  →  names a context_set
        ↓
Load that context_set only  (exact files/fields + budget)
        ↓
Act, then write a compact handoff carrying resolved context forward
```

- **Tier 0 (always):** `AGENTS.md`, `WORKFLOW.md`, `workspace.yaml`,
  `context/INDEX.md`, `context/WORKSPACE.md`, `context/PROJECT.md`. Hard budget
  **8 KB**.
- **Tier 1 (on demand):** deep docs, skills, role deltas — each route declares
  the exact set it needs in `wrapper/context-sets.yaml`, with a byte budget and
  the reason codes that unlock any conditional deep read.
- **Amortization via context receipts.** Tier-0 is read **once per session,
  never per task**. Each consequential session writes a **context receipt** —
  a versioned list of the evidence it already resolved (path + git-sha/digest)
  plus the invariant IDs and route decision it acted on. It is a *cache key, not
  truth*: on resume or hand-down, the agent verifies the wrapper version and
  digests and **reloads only changed, missing, or newly required primary
  evidence** — it never re-runs root orientation and never copies raw sources or
  accepted docs into the receipt. An invalid or stale receipt cannot authorize a
  write; it forces a fresh read. This is the concrete mechanism that cuts the
  *re-read* across a plan's children and resumes.

---

## Part 5 — The one router (two stages)

`cc-session-entry` is the sole evaluator; `cc-whats-next` only **renders** its
recommendation (it owns no logic). Routing is two stages so Tier-0 is never
judged on state it hasn't read:

**Stage A — probe.** From the human request + host/session identity + Tier-0
alone, choose exactly one *evidence packet*: a named context set, a specific
missing human choice, a safety/recovery blocker, or a read-only answer needing
no deeper reads. Tier-0 decides **what to load**, never the final action.

**Stage B — action.** After that packet is loaded, emit the decision:

```yaml
intent:        execute-plan          # what the human asked for
phase:         orienting             # session lifecycle state, never a capability
probe:         run-plan-preflight    # the evidence packet chosen at Stage A
eligibility:   ready                 # does observed state permit the candidate?
authorization: explicitly-requested  # may the agent MUTATE? (recommend ≠ permit)
capability:    cc-run-plan           # the single selected operation
context_set:   run-plan              # the bounded read packet
reason_codes:  [PLAN_APPROVED, DEPENDENCIES_READY, NO_LIVE_OWNER]
human_gate:    none                  # or: plan-approval | status-change | publication | archive
```

`eligibility` ≠ `authorization` is the core safety property: an action can be
*possible* yet still require the human to say the word. **Precedence:** a
corrupt/dirty/conflicting state routes to `blocked/recovery` before any
mutation; a valid child packet routes to that bounded role; a bound live session
routes to resume after compatibility checks; an explicit human request is
preflighted without substitution; otherwise the router yields one read-only
next-action recommendation, or asks one focused question.

### 5.1 The canonical routing table (single source; nothing else duplicates it)

| Situation (evidence) | Route (capability) | Needs human? |
| --- | --- | --- |
| No clear objective | Orient; ask focused questions | — |
| New idea, uncertain | `cc-idea-brief` | — |
| New/updated PRD or source | `cc-gather-context` → draft knowledge | accept knowledge |
| Knowledge but no plan | `cc-create-plan` → draft | **approve plan** |
| Plan drafted, review asked | `cc-review-plan` (read-only) | — |
| Plan `draft` | (blocked from execution) | **approve plan** |
| Plan `approved`, no owner | `cc-run-plan` | — (approval already given) |
| Several connected approved plans | `cc-run-stack` | — |
| Interrupted stack run | resume `cc-run-stack` from `progress.yaml` | — |
| Existing owner/child session | resume or coordinate | — |
| Verification failed | repair within scope, or report blocker | decide if scope change |
| Scope/acceptance changed | pause | **re-approval** |
| Implementation complete + verified | prepare completion evidence | **finish plan** |
| Human asks to delete runtime | `cc-cleanup-runtime` | **confirm cleanup** |
| Malformed/contradictory state | visible recovery decision | decide |

The agent **states the chosen route and its reason codes before any
consequential action.**

---

## Part 6 — Plans: the human-centered lifecycle

Plans are how a human reviews and authorizes work. The design optimizes the
human's four moments: **read, approve, trigger, finish** — each a plain sentence.

### 6.1 The plan bundle

```text
plans/context-circuit-plans/0007-checkout-validation/
├── PLAN.md          ★ HUMAN-READABLE: objective, why, scope, tasks, status, how to approve
├── plan.yaml           machine-canonical lifecycle record (status lives here)
├── tasks.md            task contracts (IDs, scope, acceptance, verification)
└── archive.yaml        optional; append-only archive/restore eligibility (Part 9)
```

`PLAN.md` is the file a human opens. `plan.yaml` is the source of truth for
**status only** and is kept small.

### 6.2 What `PLAN.md` looks like (the human's view)

```markdown
# Plan 0007 — Checkout validation                     Status: DRAFT → needs your approval

**Why.** Cart lets invalid quantities through; PRD sources/checkout-prd.md §3.
**Outcome.** Reject non-positive and over-stock quantities before payment.

## Scope
- May change: services/checkout/validate.*, its tests
- Will NOT change: payment, inventory write path

## Tasks (at a glance)
1. CO-01  Add quantity validation          → unit tests green
2. CO-02  Surface validation errors in API → contract test green

## Acceptance
- Invalid quantity returns 422 with a typed error; valid checkout unaffected.

## Dependencies
- none

---
To approve:  say  “approve plan 0007”
To review first: say  “review plan 0007”
```

### 6.3 `plan.yaml` (machine — small, canonical status)

```yaml
id: checkout-validation
number: 7
title: Checkout validation
status: draft            # draft → approved → done  (ONLY a human gate changes this)
repository: checkout
source: { kind: accepted-prd, reference: sources/checkout-prd.md }
tasks: [CO-01, CO-02]    # task status is a projection of plan status
dependencies: []
wrapper_version: 1
```

### 6.4 The human lifecycle — plain sentences, filesystem effects

| The human says | Capability | What changes | Gate |
| --- | --- | --- | --- |
| "start / resume work" | `cc-session-entry` | routes; creates/reads a session record | — |
| "what can I do?" | `cc-whats-next` | renders the route view (read-only) | — |
| "capture this idea …" | `cc-idea-brief` | draft brief under `sources/` | — |
| "make a plan for X" | `cc-create-plan` | new `draft` plan bundle | — |
| "review plan 0007" | `cc-review-plan` | read-only readiness report | — |
| **"approve plan 0007"** | `cc-approve-plan` | `plan.yaml` draft→**approved**; tasks→ready | **✔ human** |
| **"run plan 0007"** | `cc-run-plan` | claims lease + worktree; writer→verifier children | — |
| "run the stack" | `cc-run-stack` | executes connected approved plans | — |
| **"finish plan 0007"** | `cc-finish-plan` | approved→**done** (needs evidence + passing verify) | **✔ human** |
| "clean up runtime" | `cc-cleanup-runtime` | deletes `.runtime/` after showing dirty work | **✔ human** |

**Approve → Run → Finish are three separate human moments.** Approval authorizes
*intent*; running is execution; finishing is the human confirming *done*. No
test result, Git state, or child report ever crosses those gates automatically.

### 6.5 The confirmation card (every human-gated mutation)

Before any gated change, the agent shows one uniform card and waits:

```text
Action: approve-plan
Target: plans/app-plans/0042-checkout-validation
Observed: draft, active, no owner
Will change: plan draft→approved; tasks draft→ready
Will NOT change: no lease, worktree, execution, Git, delivery, or publication
Open decisions: none
Confirm: Approve this named plan?
```

A confirmation is valid **only** for the shown action, target, effects, and
current session. If any changes, a new card is shown. This closes the
"human said *yes* to something stale" failure: a bare "yes" with no active card
is a clarification, never an inferred gate.

### 6.6 Review, approval, and trigger in detail

- **Review** (`cc-review-plan`, read-only) produces a readiness report a human
  skims: scope bounded? acceptance observable? dependencies ready? conflicts?
  It never changes status. A plan may be approved without a prior review.
- **Approve** (`cc-approve-plan`) is the *only* way `plan.yaml` becomes
  `approved`. It requires an explicit in-session confirmation, shows the task
  list for visibility, and does nothing else (no lease, no worktree, no run).
- **Trigger** (`cc-run-plan`) claims the plan lease and an exclusive worktree
  from the repository's default/active branch, then directs a **writer child**
  and a later **independent read-only verifier child**. Sequential tasks share
  one writer. `cc-run-stack` does the same across connected approved plans,
  freezing parent SHAs — it is not a scheduler.

---

## Part 7 — Expected agent behavior in every case

Every case follows: **read Tier-0 → probe state → state route + reasons → act →
handoff.** The agent always separates *observed / decided / assumed / blocked /
next*.

### 7.1 Root session

- **No objective:** orient; ask focused questions; write a root session record;
  take no consequential action.
- **New PRD/source given:** read only the named source files (never scan
  `sources/`), draft Product Knowledge with provenance, and stop for the human to
  accept it. Raw sources stay in `sources/`.
- **Knowledge but no plan:** draft a `PLAN.md` + `plan.yaml` in `draft`; present
  it; request approval. Never self-approve.
- **Approved plan, no owner:** claim the lease + worktree, run writer then
  verifier children, report. (Approval already covers this — no second gate.)
- **Connected approved plans:** run `cc-run-stack`; base dependents on parent
  frozen SHAs; hand the human the leaf worktrees; do not mark anything done.
- **Resume:** read the explicit session record + latest handoff; reconcile stale
  task projections; **never reconstruct ownership from chat history or a global
  pointer.**
- **Verification failed:** repair within approved scope, or, if the fix needs
  scope/acceptance change, **pause and request re-approval.**
- **Scope change mid-flight:** pause; surface it; request a human decision.
- **Complete + verified:** record completion evidence, mark
  `ready-for-human-status-change`, and ask for `cc-finish-plan`. It does not set
  `done` itself.
- **Blocked/ambiguous/malformed state:** stop, present a visible recovery
  decision, preserve all runtime files.

### 7.2 Child session

- Reads its own record + parent packet + the exact plan/task + only the required
  context set. Returns findings/changes/tests/questions/blockers to the parent.
- **Never broadens scope, changes canonical status, satisfies a human gate, or
  writes another session's state.** If the assignment is contradictory or unsafe,
  it returns a blocker instead of inventing authority.
- **Verifier child** is read-only for implementation/lease/worktree; it may write
  **only its own handoff** so its pass/fail is durable.

### 7.3 Concurrency & recovery

- Two sessions want the same plan → only the valid lease-holder writes; the other
  becomes read-only or blocked.
- Different plans touching overlapping paths → both proceed in separate
  worktrees; the root **reports the integration risk before merge/publication.**
- A crashed session → preserve its worktree + records; offer resume or an
  explicitly recorded takeover; never reset or clean silently.

### 7.4 The rest of the case matrix (compact)

| Situation | Required behavior |
| --- | --- |
| Uninitialized workspace | Ask mode, repos/project, roles, branches in one set; zero repos is valid; confirm before any accepted write |
| Dirty base repository | Do **not** create a worktree from an uncertain base or discard changes; report the exact dirt and block |
| Human says "yes" with no active card | Do not infer target/gate from old chat; ask them to name the action + target |
| Draft plan asked to run | Explain it needs approval; never auto-approve |
| Approved plan merely discovered | Report eligibility + the exact run prompt; do **not** claim the lease |
| Live foreign owner | Stay read-only; never steal a lease or write its worktree; offer wait / coordinate / explicit takeover |
| Stale/ambiguous lease | Preserve evidence; show takeover consequences; heartbeat expiry alone never grants ownership |
| Host cannot spawn a child | Do not let root self-implement or self-verify; preserve the claim and report the missing host primitive |
| Multi-repository plan | Explicit repo/path/worktree boundary per target; report overlap before delivery |
| Provider/integration unavailable | Return `unavailable`/`denied`; continue the filesystem workflow uninterrupted |
| "Restore plan X" | Recheck current compatibility; append a `restored` event; does **not** approve or run |
| Verified work needs delivery | Apply the configured policy (`manual`/`remote-review`/`local-target`); never infer Git/external authorization |

---

## Part 8 — Runtime state (private, resumable)

```text
.runtime/
├── sessions/<id>/    session.yaml, delegation.yaml, context-receipt.yaml, handoff.md
├── plans/<plan-id>/  lease.lock/owner.yaml, lease.yaml, prompt.md,
│                     completion.yaml, reviews/, handoffs/<session>-<seq>.md
├── stacks/<id>/      graph.yaml, progress.yaml, lease.*
└── worktrees/<repo>/<plan-id>/
```

A session record (minimal):

```yaml
session_id: sess-002
parent_session_id: sess-001
root_session_id: sess-001
role: verifier
plan: plans/checkout-plans/0007-checkout-validation
scope: read-only
worktree: .runtime/worktrees/checkout/0007-checkout-validation
status: verifying
wrapper_version: 1
host_binding: { provider: codex, task_id: 01a0…  }   # optional; exact lookup, not a global pointer
next_action: Report verification to sess-001
```

Leases are atomic (a lock directory) so two writers can never claim one plan.
`.runtime/` is preserved until a human runs `cc-cleanup-runtime`.

---

## Part 9 — Archive, delivery, and upgrade

- **Archive** is a separate `archive.yaml` sidecar (append-only
  archived/restored events). It hides a plan from ordinary discovery/execution
  without touching `plan.yaml` status. Only `cc-archive-plan`, human-gated, writes
  it. An archived *done* dependency still resolves; an archived unfinished one
  stays a visible blocker.
- **Delivery policy** (`remote-review`, `local-target`, `manual`) is optional,
  set only via `cc-configure-workspace`, and never stores credentials.
- **Live-workspace upgrade** (the shipped-template concern): the wrapper carries
  `wrapper_version`. Compaction and added invariant IDs are **backward-compatible**
  for a workspace mid-plan; changes to record fields or route enums are
  **breaking** and require a documented, human-gated migration. A workspace with a
  live lease created under an older version either keeps working or is migrated
  explicitly — **never silently corrupted.**

---

## Part 10 — Measurement and guardrails (proof over proxy)

Shipped as `test/`, one top-level `sh test/acceptance.sh`:

| Check | Proves |
| --- | --- |
| `context-budget` | Tier-0 ≤ 8 KB; every route's context set within budget |
| `token-truth` (live agent) | **real** input tokens drop across the plan *lifecycle*, not just bytes |
| Tier-0-only **routing-clarity** | the compressed spine still routes correctly (couples the two goals) |
| independent-reader route check | routes are unambiguous to a fresh agent, not just the author |
| invariant-ownership lint | each rule has exactly one owner; tests assert IDs, not prose |
| behavior fixtures (`preserve` / `intentional-correction`) | gates, leases, worktrees, recovery, archive, completion unchanged |
| `upgrade` fixtures | an old in-flight workspace (live lease, dirty worktree, interrupted stack) resumes or visibly blocks — never corrupts |
| `security` | path traversal, symlink escape, credential/source-boundary leaks fail |
| release/staging | artifact is uninitialized, clean, host-compatible, credential-free |

Rent rule for any new mechanism: a new always-read artifact must remove **≥ 2×**
its bytes from Tier 0; a heavier representation ships only if it cuts its affected
lifecycle profile by **≥ 10% on real tokens** or makes a safety property testable
that the lighter form cannot.

---

## Part 11 — End-to-end walkthrough (a human's full journey)

```text
Human: "Start work here. Here's a PRD: sources/checkout-prd.md — validate cart quantities."
Agent: [reads Tier-0] Route: gather-context (reason: NEW_SOURCE). Reads only that
       PRD, drafts context/domains/checkout.md with provenance. "Drafted knowledge
       from §3. Accept it as canonical?"                              ← human gate 1

Human: "Yes, accept it."
Agent: Route: create-plan. Writes plans/checkout-plans/0007/… (PLAN.md + plan.yaml, DRAFT).
       "Plan 0007 drafted — 2 tasks, scope bounded to validate.*. Review or approve?"

Human: "Review it."
Agent: cc-review-plan → "Scope bounded ✓, acceptance observable ✓, no deps, no
       conflicts. Ready for approval."                                (read-only)

Human: "Approve plan 0007."                                          ← human gate 2
Agent: plan.yaml draft→approved; tasks→ready. "Approved. Run it now?"

Human: "Run plan 0007."
Agent: Claims lease + worktree. Writer child implements CO-01, CO-02. Verifier
       child (read-only) runs tests → pass. completion.yaml =
       ready-for-human-status-change. "Implemented + independently verified. Diff
       in worktree. Finish plan 0007?"                                ← human gate 3

Human: "Finish plan 0007."
Agent: plan.yaml approved→done; tasks→done. "Done. Runtime preserved. Clean up when ready."
```

At every step the human read at most **one** file (`PLAN.md`) and spoke plain
sentences; the agent held the coordination in files and never crossed a gate.

---

## Part 12 — Why this design meets the three goals

- **Less context:** two tiers + route-scoped sets cut the always-read spine ~64%
  and the thorough entry ~74%; context receipts cut the *re-read* across a plan's
  children and resumes — verified on **real tokens**, not bytes.
- **Clearer routing:** one two-stage router, one table, `intent/eligibility/
  authorization/phase` separated, reason codes on every route — and a Tier-0-only
  clarity metric that fails if compression ever hurts routing.
- **Consistent docs:** one owner per rule via invariant IDs, enforced by semantic
  tests instead of duplicated prose; a human-facing `PLAN.md` and `README.md` so
  people read prose and machines read `plan.yaml`.

The human experience reduces to four sentences — **make a plan, approve it, run
it, finish it** — over a filesystem that records everything and a set of agents
that never step past a gate.
