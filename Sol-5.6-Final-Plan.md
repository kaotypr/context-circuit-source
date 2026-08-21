# Sol 5.6 Final Plan — Context Circuit 1.0

Status: **draft — human review required**

Revision: **2 — 2026-08-21**, incorporating the final comparative review.

Purpose: define the complete next Context Circuit product: its released
filesystem, document authority, context-loading model, routing, human workflow,
agent behavior, plan format, execution/recovery rules, compatibility contract,
tests, and staged implementation.

This file is a requested standalone design and migration plan. It does not
approve itself, execute work, change canonical plan status, merge, publish, or
deploy. If accepted, convert it into one or more canonical plan bundles under
`plans/context-circuit-plans/` before execution.

This full design is maintainer review material, not workspace bootstrap
context. Routine agents must not load it; implementation sessions receive only
the accepted canonical bundle and route-selected sections needed for their work.

## Review this plan in five minutes

The proposed product has four layers:

1. **A tiny entry spine** identifies the workspace, applies safety rules, and
   selects a bounded evidence probe.
2. **One two-stage router** first determines what must be read and then chooses
   exactly one eligible and authorized action.
3. **Human-facing plans** separate review, approval, execution, completion, and
   delivery. Approval never starts execution.
4. **Private runtime evidence** gives each writer an exclusive worktree, gives
   verification to a separate read-only child, and makes interruption safe.

The implementation is split into nine dependency-ordered work packages. The
release gate requires semantic parity, a Tier-0 ceiling of 8 KiB, a complete
single-plan lifecycle ceiling of 45 KiB, at least 70% static-context reduction,
and measured live-token evidence or an explicit publication waiver.

An ordinary human uses four surfaces: the workspace `README.md`, a plan's
`PLAN.md`, the conversational “what's next” card, and the latest surfaced
handoff. Everything else is machine state or optional reference material.

Before approval, decide only the three items in Section 23: source/release
layout, budgets and evidence threshold, and lifecycle/schema compatibility.

## 1. Product outcome

Context Circuit is a filesystem-backed operating system for AI-assisted project
work. A human speaks naturally; the wrapper determines the smallest safe route,
loads bounded evidence, performs read-only work automatically, requests explicit
confirmation at meaningful gates, delegates implementation into exclusive
worktrees, verifies independently, and leaves durable continuation evidence.

The redesign must deliver:

- one obvious human entry point: “Start or resume work in this workspace”;
- one two-stage router: first select evidence, then select the final action;
- one owner per normative rule;
- route-specific context instead of whole-workspace reading;
- plan documents that a human can understand from one `PLAN.md`;
- explicit, separate human actions for review, approval, execution, completion,
  delivery, archive, takeover, and cleanup;
- safe concurrency, interruption, and old-workspace upgrade;
- at least 70% less static wrapper material in a standard plan lifecycle;
- a complete offline filesystem workflow with optional host integrations.

## 2. Product principles and fixed decisions

1. **Conversation is the interface.** Capability names help discovery; humans
   are never required to type internal commands.
2. **Eligibility is not authorization.** An approved plan may be eligible to
   run, but execution starts only from an explicit human run request.
3. **One fact, one owner.** Other files reference the owner and explain only
   their local consequence.
4. **Load progressively.** Tier 0 chooses the next evidence packet; final routing
   happens only after that packet is read.
5. **Primary evidence wins.** Context receipts and handoffs locate evidence;
   they never replace accepted context, plans, repository state, or Git.
6. **Plans are intent, sessions are execution, runtime is evidence.** These
   remain separate authorities.
7. **Humans control gates.** Acceptance, plan approval, execution trigger,
   material scope change, completion, takeover, delivery, archive, publication,
   deployment, and destructive cleanup are explicit.
8. **Writers are isolated; verification is independent.** One writing owner per
   plan, one writer per worktree, separate read-only verifier.
9. **The core works offline.** Integrations may fail without blocking the
   filesystem lifecycle.
10. **Every abstraction pays rent.** A new always-read artifact removes at least
    twice its bytes elsewhere; heavier machinery must improve an affected
    lifecycle context by at least 10% or uniquely make a safety property
    testable.

## 3. Vocabulary humans and agents use consistently

| Term | Meaning |
| --- | --- |
| Wrapper | Versioned instructions, route contract, context sets, schemas, skills, and role contracts shipped by Context Circuit. |
| Workspace | Mutable project identity, Product Knowledge, sources, plans, runtime evidence, and registered repositories. |
| Host task | The Codex/Claude/Cursor conversation or execution context. It may optionally bind to a Context Circuit session. |
| Session | One root or child execution context recorded under `.runtime/sessions/`. |
| Plan | Human-reviewed intended work. `plan.yaml` owns lifecycle status. |
| Task | Bounded unit inside a plan. Its status projects the plan lifecycle; it is not a second approval gate. |
| Worktree | Exclusive writable Git checkout assigned to one plan execution. |
| Context set | Exact route-selected files/fields plus conditional reads and a byte budget. |
| Context receipt | Versioned list of evidence references and digests already resolved for a session. It is a cache key, not truth. |
| Handoff | Durable report of work, evidence, changes, tests, blockers, and next action. |
| Implemented | Runtime evidence: writer finished, verifier passed, clean committed worktree. Plan remains `approved`. |
| Done | Human-confirmed canonical plan status after completion evidence. |

Do not call host tasks, Context Circuit sessions, plan tasks, and plans all
“tasks” in durable documents. User-facing conversation may use the host's normal
word, but every persisted record uses the precise term.

## 4. Source repository versus released workspace

The current repository incorrectly acts as both the Context Circuit product and
an uninitialized template. The new source repository separates authored wrapper
files, mutable template seed, maintainer material, and its own Product Knowledge.

### 4.1 Maintainer source layout

```text
context-circuit/
├── AGENTS.md                         # maintainer/source instructions
├── WORKFLOW.md                       # maintainer workflow summary
├── workspace.yaml                    # identifies the Context Circuit product
├── context/                          # Product Knowledge for Context Circuit
├── wrapper/                          # canonical shipped wrapper
│   ├── manifest.yaml
│   ├── adapters/
│   │   ├── AGENTS.md                 # released root adapter
│   │   ├── WORKFLOW.md               # released root adapter
│   │   └── README.md                 # released human quick start
│   ├── contracts/
│   │   ├── invariants.yaml
│   │   ├── routes.yaml
│   │   ├── context-sets.yaml
│   │   └── schemas/
│   │       ├── workspace.yaml
│   │       ├── session.yaml
│   │       ├── delegation.yaml
│   │       ├── context-receipt.yaml
│   │       ├── lease.yaml
│   │       ├── handoff.yaml
│   │       ├── plan.yaml
│   │       ├── task.yaml
│   │       ├── archive.yaml
│   │       ├── completion.yaml
│   │       └── stack.yaml
│   └── migrations/
├── .agents/skills/                   # thin host-discoverable adapters
├── agents/                           # thin coordinator/writer/verifier deltas
├── template/                         # blank mutable workspace seed
│   ├── workspace.yaml
│   ├── context/
│   ├── sources/
│   └── plans/
├── docs/                             # explanatory and maintainer docs
├── test/                             # semantic, budget, upgrade, release tests
└── scripts/                          # maintainer-only release assembly
```

The physical `wrapper/` and `template/` split is the preferred target. It ships
only if migration tests prove host compatibility and the contract files satisfy
the context-rent rule. Until then, the release manifest provides the same
ownership boundary without moving paths.

### 4.2 Released and instantiated workspace

```text
workspace-root/
├── README.md                         # human entry point
├── AGENTS.md                         # tiny safety/bootstrap adapter
├── WORKFLOW.md                       # tiny lifecycle/router summary
├── wrapper/
│   ├── manifest.yaml                 # wrapper version, catalog, budgets
│   └── contracts/                    # route-selected, normally not all read
├── .agents/skills/                   # natural-language capability discovery
├── agents/                           # execution-role deltas
├── workspace.yaml                    # mutable workspace identity/configuration
├── context/
│   ├── INDEX.md                      # navigation only
│   ├── WORKSPACE.md                  # accepted workspace identity
│   ├── PROJECT.md                    # accepted project purpose/boundaries
│   ├── ARCHITECTURE.md               # project architecture
│   ├── CONVENTIONS.md                # project conventions
│   ├── DECISIONS.md                  # accepted decisions
│   ├── SOURCES.md                    # provenance rules
│   ├── sources.yaml                  # provenance registry
│   ├── domains/                      # canonical bounded product behavior
│   └── roles/                        # cross-domain business-role perspectives
├── sources/                          # passive raw evidence and authored briefs/PRDs
├── plans/
│   ├── README.md                     # human guide to plan lifecycle
│   └── <repository-key>-plans/
├── .runtime/                         # private resumable execution evidence
└── <registered repositories>         # actual product code
```

Wrapper upgrades may replace wrapper-owned files only. They never overwrite
workspace identity, accepted Product Knowledge, sources, plans, `.runtime/`, or
product repositories.

## 5. Which files matter to a human

Humans normally need only four surfaces:

| Priority | Surface | Human use |
| ---: | --- | --- |
| 1 | `README.md` | Natural-language quick start, common requests, gates, and where to review results. |
| 2 | `plans/<repo>-plans/<plan>/PLAN.md` | The complete human plan: outcome, scope, tasks, risks, acceptance, status, and exact next gate. |
| 3 | Conversational “what's next” card | Current state, eligible action, required confirmation, blockers, and one exact next request. |
| 4 | Latest handoff surfaced in conversation | Work performed, verification, limitations, unresolved decisions, and recovery action. |

`workspace.yaml`, `context/PROJECT.md`, `context/DECISIONS.md`, and `plan.yaml`
remain inspectable authorities, but the agent summarizes their relevant fields
and offers direct links; ordinary operation must not require humans to edit or
cross-read them. Configuration and accepted-context changes happen through
focused conversational confirmation cards.

For recovery, a human may inspect the latest
`.runtime/sessions/<session-id>/handoff.md` and plan `completion.yaml`. Humans
should not manually edit leases, worktree ownership, stack progress, or another
session's records. Wrapper maintainers inspect `AGENTS.md`, `WORKFLOW.md`,
`wrapper/manifest.yaml`, and the contracts; ordinary workspace users need not.

## 6. Authority and context-loading model

### 6.1 Authority order

1. Host/system safety instructions.
2. `AGENTS.md` bootstrap safety and source boundary.
3. `WORKFLOW.md` lifecycle and routing contract.
4. Repository-local instructions for code behavior.
5. Human-approved plans and accepted decisions.
6. Accepted, source-cited Product Knowledge.
7. Selected raw sources and observed repository evidence.
8. Runtime state and handoffs.
9. Agent assumptions and proposals.

Runtime may describe progress; it cannot redefine intent. A handoff may locate
evidence; it cannot overrule a plan. An external provider status is never
canonical lifecycle state.

### 6.2 Tier 0

Tier 0 exists only to identify role, workspace identity, the correct evidence
probe, safety conflicts, and whether a human gate is implicated. It contains:

- released `AGENTS.md`;
- released `WORKFLOW.md`;
- `wrapper/manifest.yaml`;
- compact route and context-set indexes;
- `workspace.yaml` identity fields;
- `context/INDEX.md`, `WORKSPACE.md`, and `PROJECT.md` summaries;
- exact current session/delegation record when the host binding provides it.

Hard budget: 8 KiB; stretch: 6 KiB. All context budgets use
1 KiB = 1,024 bytes.

Tier 0 does not read all Product Knowledge, plans, runtime sessions, sources,
repository instructions, or deep contracts.

### 6.3 Route context sets

Each route declares:

```yaml
id: run-plan
required:
  - selected plan.yaml fields
  - selected task headers
  - declared Product Knowledge references
  - repository instructions
conditional:
  lease-conflict:
    - runtime lease contract
  resume:
    - current context receipt
    - latest handoff
prohibited:
  - broad sources scan
  - unrelated plans
budget_bytes: 24000
```

The agent reports any budget overrun with the reason code and newly selected
evidence. It never silently broadens context.

### 6.4 Context receipts and continuation

Every consequential session may record:

```yaml
schema_version: 1
wrapper_version: 1.0.0
session_id: sess-0042
route_decision_digest: sha256:...
context_set: run-plan
references:
  - path: context/domains/checkout/README.md
    revision: <git-sha-or-content-digest>
invariants:
  - PLAN_APPROVED
  - EXCLUSIVE_WORKTREE
created_at: 2026-08-21T12:00:00Z
```

On resume, the agent verifies the wrapper version and digests, then reloads only
changed, missing, or newly required primary evidence. A receipt never copies raw
sources or accepted documents. An invalid receipt cannot authorize a write.

## 7. The two-stage router

### 7.1 Stage A: probe decision

From the human request, host/session identity, and Tier 0, choose exactly one:

- a bounded context set;
- a specific missing human choice;
- a safety/recovery blocker;
- a read-only answer requiring no deeper workspace evidence.

### 7.2 Stage B: action decision

After the selected evidence is loaded, emit:

```yaml
intent: execute-plan
session_kind: root
phase: orienting
probe: run-plan-preflight
eligibility: ready
capability: cc-run-plan
authorization: explicitly-requested
reason_codes:
  - PLAN_APPROVED
  - DEPENDENCIES_READY
  - NO_LIVE_OWNER
context_set: run-plan
human_gate: none
```

The route fields mean:

- `intent`: the human's requested outcome;
- `phase`: runtime lifecycle, never a capability name;
- `probe`: evidence packet selected at Stage A;
- `eligibility`: whether observed state permits the candidate;
- `capability`: the one operation selected;
- `authorization`: read-only, explicitly requested, confirmed gate, or absent;
- `reason_codes`: stable explanations and test assertions;
- `human_gate`: exact confirmation still required, or `none`.

`cc-session-entry` is the sole evaluator. `cc-whats-next` presents its read-only
recommendation. Skills do not maintain parallel route rules.

### 7.3 Precedence

1. Conflicting instructions, unsafe paths, corrupt records, or dirty ownership
   state route to `blocked/recovery` before mutation.
2. A valid child delegation packet routes to that bounded child role.
3. An explicitly bound non-terminal session routes to resume after compatibility
   and ownership checks.
4. An explicit human capability request is preflighted without substitution.
5. Without an explicit action, state produces one read-only next-action
   recommendation.
6. If evidence remains insufficient, ask one focused question or orient.

An archived plan is filtered before approval/execution/finish selection. A live
foreign owner blocks writing. A generic “start or resume” may recommend an
approved plan but does not itself authorize execution.

## 8. Complete behavior matrix

| Situation or human request | Route | Required agent behavior | Human involvement |
| --- | --- | --- | --- |
| “What is this workspace?” | orient | Read Tier 0, summarize identity, known project, repositories, state, and one next action. Do not create runtime for a purely read-only answer. | None. |
| Uninitialized workspace | initialize | Ask mode, repositories/project items, roles, default branches in one compact set. Zero repositories is valid. | Confirm identity before accepted write. |
| No repository, uncertain idea | idea-brief | Offer brief, selected source intake, both, or defer. Do not invent a repository. | Accept/revise/leave draft. |
| Clear idea, no durable requirements needed | idea-brief or direct plan | Choose the smallest artifact; explain why PRD is unnecessary. | Accept artifact or request direct plan. |
| Human requests PRD | create-prd | Read only selected evidence, draft testable requirements/non-goals/scenarios/provenance. | Accept or revise PRD; no plan approval implied. |
| Human names raw sources | selected-source intake | State exact files and reasons before reading; never scan `sources/`. Record provenance. | Confirm source selection if scope is ambiguous. |
| Product Knowledge missing/stale | gather-context | Select one domain or role, read linked evidence only, draft proposed page with freshness and contradictions. | Accept proposed context before canonical use. |
| Sources contradict accepted context | context-review | Preserve accepted facts, show contradiction and proposed refresh. | Decide whether to accept revised context. |
| Clear request with sufficient evidence | create-plan | Draft one coherent plan; do not force Idea Brief or PRD. | Review/approve later. |
| “Review plan X” | review-plan | Read-only readiness review; inspect plan/tasks/evidence/state; produce Review Card. | Resolve listed decisions; no status change. |
| “Approve plan X” | approve-plan | Preflight completeness/archive/dependencies, display Approval Card, wait for exact confirmation. | Confirm named plan; status only changes to approved. |
| Human says “yes” without a current card | clarification | Do not infer target or gate from old conversation. | Name or reconfirm action and target. |
| Draft plan requested to run | blocked | Explain it needs approval; do not auto-approve. | Review/revise/approve. |
| Approved plan merely discovered | recommend-run | Report eligibility and exact run prompt. Do not claim lease. | Explicitly request run. |
| “Run approved plan X” | run-plan | Preflight, create root record, claim lease atomically, create exclusive worktree, direct writer then verifier. | Run request authorizes execution within approved scope. |
| One plan has sequential tasks | run-plan | One writer child and worktree execute tasks in dependency order; one later verifier. | No per-task approval ceremony. |
| Independent approved plans | separate run-plan executions | Separate leases/worktrees/children; may run concurrently; report overlap before delivery. | Explicitly select plans or confirm proposed set. |
| Connected approved plans requested | run-stack | Freeze DAG/progress, use parent SHAs and deterministic joins, writer+verifier per member. | Explicit run request; no separate stack approval artifact. |
| Live same-session owner | resume | Validate receipt, handoff, lease, worktree, wrapper compatibility; continue bounded work. | None unless new gate/blocker. |
| Live foreign owner | coordinate/blocked | Remain read-only; do not steal lease or write worktree. | Choose wait, coordinate, or explicit takeover process. |
| Stale/ambiguous lease | recovery | Preserve evidence; inspect heartbeat/worktree; present takeover consequences. | Explicit takeover confirmation. |
| Dirty base repository | blocked | Do not create normal execution worktree from uncertain base or discard changes. Report exact dirt. | Resolve/preserve base state. |
| Child lacks complete delegation | child-blocked | Refuse work; report missing identity, scope, paths, permissions, acceptance, or stop conditions. | Parent repairs packet. |
| Host cannot create child | host-blocked | Do not silently let root implement or self-verify. Preserve claim/state and explain limitation. | Choose supported host or an explicitly redesigned plan. |
| Writer reaches unapproved scope | scope-blocked | Stop, preserve dirty work and evidence, return handoff. | Revise/reapprove plan if material. |
| Verification passes | ready-for-finish | Record task evidence and completion request; do not mark done or deliver. | Review and explicitly finish. |
| Verification fails within scope | repair-loop | Return to same writer for bounded repair, then independent re-verification. | None unless hard stop/material scope. |
| Verification needs scope change | blocked | Preserve failure; do not broaden plan. | Revise and reapprove. |
| “Finish plan X” | finish-plan | Check completion/verifier/evidence/Git, show Finish Card, wait for confirmation, then mark plan/tasks done and release owned lease. | Explicit status-change confirmation. |
| “Archive plan X” | archive-plan | Check ownership, worktree, dependents, compatibility; show effects; append event only. | Explicit archive confirmation and reason. |
| “Restore plan X” | restore-plan | Recheck current compatibility; append restored event; do not approve/run. | Separate restore confirmation. |
| “Configure delivery/integration” | configure | Ask only capability-specific questions, show fallback, store intent not credentials. | Confirm configuration; later actions remain gated. |
| Provider unavailable | offline fallback | Report unavailable/denied; continue filesystem workflow. | None unless user chooses another path. |
| Verified work needs delivery | delivery-choice | Apply configured manual/remote-review/local-target proposal; never infer Git/external authorization. | Confirm commit/push/merge/publication/deployment as applicable. |
| “Clean runtime” | cleanup | Inspect every runtime worktree and live session; show dirty/unpushed consequences. | Confirm workspace-wide deletion; extra explicit discard confirmation if dirty. |
| Wrapper upgrade with live work | compatibility | Read wrapper version/context receipt, classify compatible/migration-needed/blocked; preserve runtime and worktrees. | Confirm any breaking migration. |

## 9. Agent roles and required outputs

### 9.1 Root coordinator

Owns the human request, router, gates, plan lease, worktree creation, child
delegation, evidence consolidation, and human-facing cards. It may write its own
session state and the lease it owns. It never replaces a writer or verifier for
approved execution.

Every consequential root update reports:

```text
Observed state
Selected route and reason
Evidence read
Action performed or proposed
Files/state changed
Verification
Blockers and human decisions
One next safe action
```

### 9.2 Writer child

Reads its delegation, exact context set, approved plan/tasks, repository
instructions, and assigned worktree. Writes only delegated paths in that
worktree. It never changes wrapper context, plan/task status, another session,
external systems, or the base checkout.

Stops on scope expansion, contradictory evidence, ownership mismatch, dirty
uncertainty, missing dependency, or material acceptance change.

### 9.3 Verifier child

Receives a separate session and read-only worktree permission. Reproduces the
approved verification commands and checks acceptance criteria, scope, ownership,
regressions, and limitations. It writes only its own handoff. It never repairs,
changes statuses, or authorizes delivery.

### 9.4 Research/planning child

Optional and read-only unless the explicit artifact scope grants a workspace
write. It receives selected sources, questions, expected evidence, and stop
conditions. Findings return to the root; a child cannot accept Product
Knowledge or satisfy a human gate.

### 9.5 Delegation contract

Every child packet includes session/parent/root IDs, role, objective, plan/task,
paths, non-goals, context set/receipt, repository/worktree, permissions,
acceptance criteria, verification, expected evidence, stop conditions, and
handoff schema. Missing required fields block child entry.

## 10. How humans use Context Circuit

Humans can use plain language. The agent maps it to capabilities and always
names the resulting effect.

### 10.1 Common prompts

```text
Start or resume work in this workspace.
What is this project and what should happen next?
Set up this workspace for these repositories.
Capture this idea as a short brief.
Create a PRD from these selected files: ...
Refresh the checkout domain context using only ...
Draft a plan for ...
Review plan checkout-validation.
Approve plan checkout-validation.
Run approved plan checkout-validation.
Run the connected approved plans A, B, and C.
Show the latest handoff and blockers for plan X.
Finish plan X.
Archive plan X because it is replaced by Y.
Configure manual delivery for repository app.
Prepare the verified work for remote review.
Clean runtime after showing me dirty and unpushed work.
```

The human does not need to say `cc-...`; names appear in responses so advanced
users can be precise.

### 10.2 Standard confirmation card

Before any human-gated mutation, show:

```text
Action: approve-plan
Target: plans/app-plans/0042-checkout-validation
Observed state: draft, active, no owner
Will change: plan status draft→approved; included tasks draft→ready
Will not change: no lease, worktree, execution, Git, delivery, or publication
Risks/open decisions: none
Confirmation requested: Approve this named plan?
```

A confirmation is valid only for the shown action, target, effects, and current
session. If those change, show a new card.

### 10.3 Human journey patterns

**New idea:** initialize → optional Idea Brief → optional PRD → accepted context
→ draft plan → review → approve → explicit run → verify → finish → delivery →
optional cleanup.

**Clear maintenance change:** orient → selected repository evidence → draft
plan → review/approve → explicit run. Idea Brief and PRD are skipped.

**Existing approved plan:** start/resume → eligibility recommendation → human
“Run plan X” → execution.

**Interrupted work:** start/resume with session or plan → compatibility/lease/
worktree inspection → resume or human-gated takeover.

**Research only:** ask question → selected evidence → answer/proposed context;
no plan/session mutation unless needed.

## 11. Product artifacts before a plan

### 11.1 Idea Brief

Lives at a human-selected path under `sources/`. Minimal frontmatter:

```yaml
kind: idea-brief
status: draft
title: Checkout reliability
```

Sections: intent, desired outcome, users/context, constraints, assumptions, open
questions, Product Knowledge references, provenance. Human acceptance makes it
accepted but does not require a PRD or authorize implementation.

### 11.2 PRD

Also lives at a human-selected path under `sources/`. Sections: problem/outcome,
users/scenarios, numbered requirements, non-goals, acceptance criteria,
constraints/dependencies, assumptions, questions, Product Knowledge references,
provenance. Acceptance defines requirements; it does not approve a plan.

### 11.3 Product Knowledge

- `PROJECT.md`: project purpose and boundaries.
- Domain page: canonical bounded behavior; workflow pages own exact flows.
- Role page: cross-domain business perspective, linking rather than copying.
- `DECISIONS.md`: accepted decisions with date, rationale, consequences, and
  supersession.
- `sources.yaml`: selected-source provenance and freshness.

Generated context is `proposed` until human acceptance. Contradictions stay
visible; unsupported sections are omitted rather than invented.

## 12. Canonical plan bundle

### 12.1 Directory shape

```text
plans/app-plans/0042-checkout-validation/
├── plan.yaml                         # canonical lifecycle and machine index
├── PLAN.md                           # sole human review entry point
├── tasks/
│   ├── APP-0042-01.md
│   └── APP-0042-02.md
└── archive.yaml                      # optional; absent means active-compatible
```

Review reports, execution evidence, prompts, handoffs, completion, and leases
belong under `.runtime/`, not in the plan bundle.

### 12.2 `plan.yaml`

```yaml
schema_version: 2
id: checkout-validation
number: 42
title: Make checkout validation deterministic
status: draft
repository: app
source:
  kind: accepted-prd
  reference: sources/checkout-prd.md
  revision: sha256:...
product_knowledge:
  references:
    - context/PROJECT.md
    - context/domains/checkout/README.md
objective: Detect invalid checkout state before payment submission.
scope:
  paths:
    - src/checkout/
    - test/checkout/
non_goals:
  - Redesign payment-provider integration.
dependencies:
  - plan: checkout-foundation
    required: done
tasks:
  - APP-0042-01
  - APP-0042-02
acceptance:
  - id: AC-01
    outcome: Invalid checkout state is rejected before payment submission.
  - id: AC-02
    outcome: Valid checkout behavior remains unchanged.
verification:
  - id: VT-01
    command: sh test/checkout.sh
    proves: [AC-01, AC-02]
human_gates:
  - plan-approval
  - status-change
created_at: 2026-08-21T12:00:00Z
updated_at: 2026-08-21T12:00:00Z
```

`plan.yaml` owns IDs, status, repository, source/context references, bounded
paths, dependencies, task membership, acceptance outcomes, verification IDs
and commands, and gates. Task documents reference acceptance and verification
IDs instead of copying commands. Schema v1 remains readable; v2 is additive and
never mass-rewrites old plans without a human-approved migration.

### 12.3 Human-facing `PLAN.md`

```markdown
# Make checkout validation deterministic

Status: draft
Repository: app
Source: accepted checkout PRD

## Review summary
Outcome, affected users, and why this plan is needed.

## What approval authorizes
The exact implementation scope. Approval does not start execution.

## Scope and non-goals
Human-readable boundaries linked to `plan.yaml` paths.

## Proposed solution
Behavior and architecture at reviewable depth, without task-level repetition.

## Tasks and dependencies
| Task | Outcome | Depends on |

## Acceptance criteria
- AC-01: observable outcome and evidence.

## Verification
Verification IDs, what they prove, test scope, and independent-review
expectations. Exact executable commands remain canonical in `plan.yaml`.

## Risks, assumptions, and open decisions
Only items that may affect approval or execution.

## Delivery boundary
What remains unapproved: commit/push/merge/publication/deployment.

## Provenance
Exact accepted context, source, and repository evidence used.
```

A human should understand the proposal, boundaries, risks, tasks, and proof by
reading this file alone. Its displayed status is a projection of `plan.yaml`;
review fails if the two differ. It explains other canonical fields without
becoming a second owner.

### 12.4 Task document

```markdown
---
schema_version: 2
id: APP-0042-01
plan: checkout-validation
status: draft
repository: app
paths:
  - src/checkout/validation.ts
  - test/checkout/validation.test.ts
depends_on: []
acceptance: [AC-01]
verification: [VT-01]
---

# Validate checkout state before submission

## Objective
One bounded implementation outcome.

## Work
Expected behavioral changes, not speculative code instructions.

## Non-goals
Paths and behavior the writer must not change.

## Verification
Verification IDs and task-specific observations. Do not copy the canonical
commands from `plan.yaml`.

## Expected evidence
Changed files, test results, relevant limitations, handoff path.

## Stop conditions
Scope expansion, contradiction, ownership conflict, or missing evidence.
```

Tasks own implementation detail. Task status is `draft`, `ready`, or `done`,
projected from the plan. It never represents implementation progress.

## 13. Plan review, approval, execution, and completion

### 13.1 Review

Human prompt: `Review plan checkout-validation.`

The agent reads the plan bundle, referenced accepted context, source provenance,
repository instructions/evidence, dependencies, archive state, active ownership,
and relevant delivery boundaries. It returns:

```text
Plan: checkout-validation
Outcome: ready-for-approval | needs-revision | blocked | approved-for-execution
Summary: what changes and why
Approval would authorize: exact bounded implementation intent
Approval would not authorize: execution, Git delivery, publication, deployment
Evidence inspected: paths/revisions
Tasks: ID, outcome, dependency
Acceptance mapping: criterion → task → verification evidence
Risks/assumptions: material items only
Human decisions: preferably 0–3 focused choices
Contradictions/blockers: explicit or none
Next action: revise, approve, run, or resolve blocker
```

Review is read-only. A review result may be stored under
`.runtime/plans/<plan-id>/reviews/<session-id>.md`; it never changes status.

### 13.2 Approval

Human prompt: `Approve plan checkout-validation.`

After preflight, show an Approval Card. On explicit confirmation only:

- `plan.yaml`: `draft` → `approved`;
- all included task projections: `draft` → `ready` atomically/idempotently;
- preserve every other field;
- do not claim a lease, create a worktree, start children, or perform Git work.

If already approved, report it and offer the exact run action. If incomplete,
archived, contradictory, or unknown, refuse and name the smallest repair.

### 13.3 Execution trigger

Human prompt: `Run approved plan checkout-validation.`

This separate explicit request is required even immediately after approval. It
authorizes only execution inside the approved scope. Preflight may still block
on dependencies, archive state, dirty base, ownership, incompatible wrapper,
missing child primitive, or unsafe assumptions.

On success: create/resume root session, atomically claim lease, create exclusive
worktree, write context receipt and writer delegation, run writer, then create a
separate verifier session. Approval is never inferred from the run request.

### 13.4 Completion

When all task evidence exists and the verifier passes, record
`completion.yaml: ready-for-human-status-change`. Human prompt:
`Finish plan checkout-validation.`

Show a Finish Card. After confirmation: plan `approved` → `done`, tasks →
`done`, completion → `completed`, release owned lease. Preserve worktree,
handoffs, sessions, branches, and runtime until separate cleanup/delivery.

### 13.5 Archive and restore

`archive.yaml` is an append-only eligibility sidecar. Archive/restore never
changes lifecycle status, deletes evidence, moves the bundle, or cascades.
Active owners, dirty/unpushed work, stack membership, and unresolved dependents
block the action. Each direction requires a separate confirmation.

## 14. Runtime, concurrency, recovery, and stacks

### 14.1 Runtime layout

```text
.runtime/
├── sessions/<session-id>/
│   ├── session.yaml
│   ├── delegation.yaml              # child only
│   ├── context-receipt.yaml
│   └── handoff.md
├── plans/<plan-id>/
│   ├── lease.lock/owner.yaml
│   ├── lease.yaml
│   ├── prompt.md
│   ├── completion.yaml
│   ├── reviews/
│   ├── evidence/
│   └── handoffs/
├── stacks/<stack-id>/
│   ├── graph.yaml
│   ├── progress.yaml
│   └── lease.lock/owner.yaml
└── worktrees/<repository-key>/<plan-id>/
```

There is no global current-session, current-plan, or current-stack pointer.
Optional host binding enables exact lookup but cannot authorize mutation.

### 14.2 Concurrency

- Same plan: one active writing lease; contenders become read-only/blocked.
- Different plans: separate worktrees and children; overlap is an integration
  risk reported before delivery.
- Sequential tasks: same writer/worktree.
- Verification: concurrent read-only only after writer handoff.
- Multi-repository plan: explicit repository/path/worktree boundary per target.

### 14.3 Recovery

On interruption preserve records, dirty state, lease, branch, worktree, receipt,
and handoff. Same session may resume after validation. A new session needs
explicit takeover evidence if ownership is live or ambiguous. Heartbeat expiry
alone never grants ownership.

### 14.4 Stack execution

Connected approved plans run through one runtime stack:

- freeze `graph.yaml` once;
- `progress.yaml` is the resume cursor;
- no-parent member starts from configured branch;
- one-parent member starts from parent frozen SHA;
- multi-parent member joins sorted parent SHAs in-run;
- each member still has its own lease, worktree, writer, verifier, completion
  evidence, and later per-plan finish gate;
- implemented parents unblock dependents; `done` is not required inside the run;
- no scheduler, hidden stack plan, auto-finish, merge, push, or deploy.

## 15. Configuration, delivery, and integrations

Initialization records identity only. Later, a human may configure:

- `manual`: preserve verified worktree and ask when delivery matters;
- `remote-review`: prepare a reviewable remote path after explicit commit/push
  authorization;
- `local-target`: prepare local integration and stop at the merge gate;
- optional external activity integration with explicit read/write boundary;
- host capability mapping with filesystem fallback.

`workspace.yaml` stores provider-neutral intent and authorization state, never
credentials or provider payloads. A changed repository, branch, risk, or stale
authorization requires focused reconfirmation. Provider failure returns
`disabled`, `denied`, or `unavailable` and never blocks core execution evidence.

Delivery happens after verification and remains separate from finishing. No
configuration silently grants commit, push, PR creation, merge, publication, or
deployment.

## 16. Wrapper versioning and live-workspace upgrades

`wrapper/manifest.yaml` owns `wrapper_version`, supported artifact schema
versions, contract paths, and context budgets. New sessions record the wrapper
version in `session.yaml` and receipts.

Classify upgrades:

- **compatible:** prose compression, additive invariant IDs, optional fields,
  new explanatory docs; resume after validation;
- **migration-needed:** renamed route enums, record field meaning changes,
  required schema changes; preserve state and ask for migration confirmation;
- **blocked:** unknown/malformed version, contradictory migration, live state
  that cannot be mapped safely; remain read-only.

Legacy records without `wrapper_version` remain readable as `legacy-unknown`.
Do not rewrite them merely to make them current. Upgrade fixtures include an
old live lease, dirty worktree, pending handoff, approved plan, and interrupted
stack. Wrapper publication includes upgrade and rollback notes.

## 17. Measurement and acceptance system

### 17.1 Deterministic context ledger

CI records exact wrapper-selected paths, revisions, bytes, words, packet/handoff
bytes, and optional tokenizer totals for:

- L0 orientation;
- L1 root + writer + verifier plan run;
- L2 three tasks + one resume;
- U1 legacy in-flight workspace upgrade/resume.

The pre-redesign inventory is frozen before normative edits. Current measured
values and target budgets are:

| Profile | Current bytes | Hard | Stretch | Hard-budget reduction |
| --- | ---: | ---: | ---: | ---: |
| Tier 0 | 22,349 | 8 KiB | 6 KiB | 63.3% |
| Full root entry | 76,456 | 20 KiB | 16 KiB | 73.2% |
| Root run-plan | 93,873 | 22 KiB | 20 KiB | 76.0% |
| Writer | 30,387 | 12 KiB | 10 KiB | 59.6% |
| Verifier | 30,108 | 11 KiB | 9 KiB | 62.6% |
| L1 total | 154,368 | 45 KiB | 39 KiB | 70.1% |
| L2 total | baseline frozen by FINAL-001 | 57 KiB | 47 KiB | reported after baseline |
| Resume delta | baseline frozen by FINAL-001 | 12 KiB | 8 KiB | reported after baseline |

L1 is exactly root run-plan + writer + verifier. Its 45 KiB ceiling is 46,080
bytes, a 70.1% reduction from the 154,368-byte baseline. Budgets are ceilings,
not instructions to pad packets.

### 17.2 Host token evidence

When a host exposes usage, record host/model/reasoning, three comparable runs,
median input tokens, cached tokens when known, and scenario. Core CI never
launches external agents or stores credentials. Publication target: ≥50% median
input-token reduction for L1/L2 with identical outcome. If usage is unavailable,
publication requires an explicit waiver; deterministic budgets still apply.

### 17.3 Test layout

```text
test/
├── acceptance.sh                    # one top-level command
├── contracts/                       # schema/invariant/catalog validation
├── routing/                         # probe/action fixtures and corrections
├── context-budget/                  # L0/L1/L2 ledgers and negative budgets
├── lifecycle/                       # draft/approved/done projection and gates
├── ownership/                       # leases, worktrees, child permissions
├── recovery/                        # receipts, stale/live ownership, takeover
├── stacks/                          # DAG, joins, resume, no auto-finish
├── upgrades/                        # wrapper compatibility and legacy state
├── release/                         # staged artifact identity/inventory
└── security/                        # traversal, symlink, credential/source boundaries
```

At least 50 routing fixtures pass 100% for preserved behavior. Intentional
corrections require a human-reviewed expected result. Negative tests prove that
duplicated invariants, unauthorized gates, bad leases, stale receipts, context
overruns, forbidden source scans, dirty cleanup, or artifact leaks fail.

## 18. Implementation work packages

### FINAL-001 — Freeze behavior, profiles, and ownership inventory

Dependencies: none.

Paths: `test/baselines/`, routing/upgrade fixtures, context ledger, assertion
ledger, benchmark docs.

Work: freeze L0/L1/L2/U1 profile membership; normalize probe/action decisions;
classify preserve/correction cases; map every normative statement and prose test
to an owner candidate.

Acceptance/evidence: reproducible baselines; ≥50 unique cases; no normative
edits; current acceptance green. Stop: contradictory expected behavior or
unapproved correction.

### FINAL-002 — Establish wrapper/template ownership and release assembly

Dependencies: FINAL-001.

Paths: `wrapper/`, `template/`, source identity, release manifest/script/docs.

Work: add wrapper manifest/version/catalog; separate source Product Knowledge
from blank seed; stage released root adapters and mutable seed; preserve host
paths and exclude maintainer/runtime/dirty/credential state.

Acceptance/evidence: source identifies Context Circuit; artifact is
uninitialized; exact inventory and rollback; both pass smoke/budget checks.
Stop: host compatibility or existing workspace upgrade cannot be preserved.

### FINAL-003 — Canonicalize invariants, router, and context sets

Dependencies: FINAL-001, FINAL-002.

Paths: `AGENTS.md`, `WORKFLOW.md`, wrapper contracts, `context/INDEX.md`, entry
and what-next skills, coordinator, routing tests.

Work: assign one owner per rule; compact Tier 0; implement probe/action contract,
precedence, reason codes, route budgets, and single evaluator/read-only view.

Acceptance/evidence: Tier 0 ≤8 KiB; one router; duplicated policies one owner;
all preserved cases pass; unauthorized and malformed cases block. Stop: budget
requires removing safety or route precedence needs a new policy decision.

### FINAL-004 — Deliver the human artifact and plan experience

Dependencies: FINAL-003.

Paths: released `README.md`, Idea Brief/PRD/context templates, `plans/README.md`,
plan `PLAN.md` and task schemas/templates, review/approval/finish/archive skills
and fixtures.

Work: implement natural-language quickstart, confirmation cards, canonical plan
bundle, Review Card, separate approval/run/finish triggers, typed v2 schema with
v1 compatibility, and human file guide.

Acceptance/evidence: human understands a plan from `PLAN.md`; schema owns all
canonical fields; review maps criteria→task→verification; approval never runs;
run requires explicit target; ≤3 unresolved human decisions recommended per
review. Stop: migration would silently rewrite old plan intent/status.

### FINAL-005 — Implement bounded roles, receipts, and execution runtime

Dependencies: FINAL-003, FINAL-004.

Paths: roles, run-plan/run-stack skills, runtime schemas/docs, delegation,
receipt/handoff/evidence fixtures.

Work: thin role instructions; bounded child packets; context receipts and delta
resume; atomic leases/worktrees; independent verifier; completion and stack
records.

Acceptance/evidence: root ≤20 KiB, run-plan ≤22 KiB, writer ≤12 KiB, verifier
≤11 KiB; L1 ≤45 KiB, L2 ≤57 KiB; contention loser cannot write; verifier cannot
repair; resume reloads changed evidence only. Stop: packet compaction weakens
scope, evidence, ownership, or verifier independence.

### FINAL-006 — Implement every human gate and operational fallback

Dependencies: FINAL-004, FINAL-005.

Paths: initialize/context/approve/finish/archive/configure/cleanup skills,
configuration/delivery/integration docs, gate and negative fixtures.

Work: standard confirmation cards; current-session target/effect binding;
archive dependency checks; dirty/unpushed cleanup inspection; delivery policy and
offline integration fallback; explicit takeover and scope-change recovery.

Acceptance/evidence: every behavior-matrix row has a fixture; no child satisfies
a gate; no configuration grants external mutation; cleanup cannot discard
unconfirmed work. Stop: any consequential path lacks an exact human gate.

### FINAL-007 — Add compatibility, upgrade, and rollback

Dependencies: FINAL-002, FINAL-005.

Paths: wrapper manifest/migrations, session/receipt schemas, upgrade fixtures,
release notes.

Work: compatible/migration-needed/blocked classifier; optional version records;
legacy-unknown handling; in-flight plan/lease/stack upgrade; wrapper-only update
boundary; rollback instructions.

Acceptance/evidence: U1 resumes without canonical/runtime loss; breaking change
requires confirmation; rollback restores wrapper files without reverting
workspace data. Stop: migration needs destructive or ambiguous state rewrite.

### FINAL-008 — Replace monolithic prose tests with semantic suites

Dependencies: FINAL-003 through FINAL-007.

Paths: `test/` suites in Section 17 and top-level acceptance command.

Work: migrate each old assertion through the assertion ledger; validate schemas,
invariants, routes, budgets, gates, ownership, sources, stacks, upgrades, and
release; add deliberate negative mutations.

Acceptance/evidence: one top-level command; every old protection mapped; prose
may be rewritten once without copies; each negative fixture fails for the
expected reason. Stop: an old assertion has no semantic replacement.

### FINAL-009 — Cost-gated optional machinery and final A/B

Dependencies: FINAL-008.

Candidates: alternative route representation, physical template layout if not
already proven, host-task binding, optional tokenizer/usage adapters.

Work: decide candidates independently using 2×/10% rules; run deterministic and
host evidence; independent verifier reproduces results; prepare migration,
limitations, and publication handoff.

Acceptance/evidence: static lifecycle reduction ≥70%; L1/L2 live-token reduction
≥50% or explicit publication waiver; routing 100% preserved/corrected;
natural-language metric reported; zero unresolved safety regression; staged
artifact matches manifest. Stop: any hard gate fails or machinery does not pay
rent.

## 19. Delivery sequence

```text
FINAL-001 Baseline
   ↓
FINAL-002 Wrapper/template foundation
   ↓
FINAL-003 Router + context economy
   ↓
FINAL-004 Human artifacts + plan UX
   ↓
FINAL-005 Runtime + roles + execution
   ├──→ FINAL-006 Gates + operations
   └──→ FINAL-007 Compatibility + upgrade
                  ↓
          FINAL-008 Semantic suites
                  ↓
          FINAL-009 Final A/B
```

Recommended review batches:

1. Foundation: FINAL-001–003.
2. Human workflow and execution: FINAL-004–005.
3. Gates and compatibility: FINAL-006–007.
4. Verification and release evidence: FINAL-008–009.

Each batch remains independently reviewable and revertible. It authorizes no
commit, push, PR, merge, publication, or deployment by itself.

## 20. Plan-level acceptance criteria

The new Context Circuit is acceptable only when:

1. source wrapper, template seed, instantiated workspace, and runtime are
   unambiguous and versioned;
2. ordinary humans can operate every lifecycle path through natural language;
3. ordinary humans operate through the root `README.md`, a plan's `PLAN.md`,
   the “what's next” card, and surfaced handoffs, without manually editing
   runtime or cross-reading machine authorities;
4. Tier 0 chooses the correct bounded probe and never pretends to know unread
   state;
5. one router produces one normalized action decision after evidence;
6. each normative rule, status, field, and artifact has one owner;
7. plan review is read-only, approval changes status only, and execution needs a
   separate explicit run request;
8. plan/task documents have bounded paths, dependencies, acceptance,
   verification, expected evidence, and stop conditions;
9. every behavior-matrix case has positive and relevant negative verification;
10. writer/verifier isolation, leases, worktrees, recovery, stacks, archive,
    finish, delivery, and cleanup retain their gates;
11. static wrapper lifecycle material falls by at least 70%;
12. old in-flight workspaces upgrade or visibly block without state loss;
13. offline operation remains complete and credentials never enter workspace
    state;
14. release artifact and source identities are independently verified;
15. independent verification reproduces the final evidence and limitations.

## 21. Non-goals

- No automatic scheduler or plan queue.
- No global current-session/current-plan pointer.
- No database, package runtime, mandatory cloud service, or user-facing CLI.
- No small-change bypass around plan approval, child writer, or verifier.
- No implicit source ingestion, plan approval, execution, completion, takeover,
  Git delivery, publication, deployment, archive, or cleanup.
- No external provider status as canonical lifecycle state.
- No automatic migration that overwrites accepted context, plan intent, runtime
  evidence, dirty work, or product repositories.

## 22. Risks and controls

| Risk | Control |
| --- | --- |
| New contract tree recreates context bloat | Tier/route budgets, 2×/10% rent tests, non-Tier-0 deep contracts. |
| Compression removes a safety nuance | Assertion ledger, invariant owner map, normalized positive/negative fixtures. |
| Tier 0 is judged on unavailable facts | Two-stage probe/action routing. |
| Context receipt becomes stale truth | Digests, primary-evidence validation, blocked mutation on mismatch. |
| Plan `PLAN.md` and YAML drift | Schema-owned fields, references/IDs, projected status, review lint, no duplicated status or command authority. |
| Humans accidentally run immediately after approval | Separate run request and explicit effect cards. |
| Source/template split breaks hosts | Stable released root adapters and staged-host smoke tests. |
| Upgrade corrupts live work | Version compatibility classifier, U1 fixtures, human-gated migration, rollback. |
| Real token benchmark is nonportable | Deterministic CI ledger plus host-neutral, credential-free usage evidence. |
| Integration failure blocks work | Offline filesystem fallback and provider-neutral outcomes. |

## 23. Human decisions before approval

1. Confirm the preferred source `wrapper/` + `template/` structure, subject to
   migration/rent verification, and `PLAN.md` as the sole human plan entry.
2. Confirm the Section 17 ceilings, the ≥70% deterministic L1 reduction, and
   the ≥50% live-token publication target with an explicit evidence-unavailable
   waiver.
3. Confirm separate approve/run/finish/delivery/cleanup actions and additive
   plan/task schema v2 while preserving v1 reads.

The four implementation review batches in Section 19 are the proposed default
and may be regrouped without changing product intent.

Recommended next action: resolve these three decisions, run a final read-only
plan review, then convert the accepted design into canonical draft plan bundles.
Do not execute this standalone artifact.
