# Sol 5.6 Plan v2 — Measured context compression and deterministic routing

Status: **draft — human review required; no execution or delivery authorized**

Scope: the Context Circuit wrapper source and its published workspace template.
This standalone comparison draft is not a canonical `plan.yaml`. If accepted,
the human must choose whether to convert it into a Context Circuit plan bundle
or use it as a normal wrapper-repository implementation plan.

## 1. Decision thesis

Deliver the certain token reduction first. Freeze current safety behavior in a
normalized evaluation corpus, remove duplicated prose, then introduce only the
smallest routing or packaging mechanism that proves a net benefit.

Three rules govern every task:

1. **Same safety, less context.** Human gates, plan/task lifecycles, leases,
   worktrees, recovery, and independent verification do not change.
2. **Semantic equivalence, not textual preservation.** Tests compare normalized
   route and safety decisions; they do not require duplicated wording.
3. **Every abstraction pays rent.** A new always-read artifact must remove at
   least twice its byte size elsewhere. A heavier representation ships only if
   it reduces its affected context profile by at least 10% or makes a material
   safety property testable that the lighter form cannot.

## 2. Evidence and measured baseline

Evidence read: the named Codex discussion; current wrapper instructions,
Product Knowledge, workflow/runtime/planning docs, coordinator/worker/verifier
roles, relevant skills, acceptance and release implementation, and the two
existing comparison plans. `sources/` was not read. Parallel plan files are
evidence only and remain user-owned.

Profiles use explicit file lists recorded by S2-001. Bytes and words are the
portable gates; model-token estimates are reported when a tokenizer exists.

| Profile | Current words | Current bytes |
| --- | ---: | ---: |
| Tier-0/every-session set | 3,079 | 22,349 |
| Full root entry, including entry skill and coordinator | 10,547 | 76,456 |
| Root entering `cc-run-plan` | 12,950 | 93,873 |
| Writer wrapper context | 4,075 | 30,387 |
| Verifier wrapper context | 4,033 | 30,108 |
| One-plan static wrapper total | — | 154,368 |
| All 15 skill bodies | 7,252 | 51,146 |

Duplication baseline: writer-child policy in 12 files, standard single-plan
entry in 7, and no-global-current-session policy in 4. The source checkout also
describes itself as an uninitialized template although it is the Context Circuit
product repository.

## 3. Acceptance budgets

Profile definitions are frozen before edits. Hard limits block completion;
stretch values guide further safe compression.

| Profile | Hard limit | Stretch | Minimum reduction |
| --- | ---: | ---: | ---: |
| Tier-0/every-session | 8 KB | 6 KB | 64% |
| Full root entry | 20 KB | 16 KB | 74% |
| Root `cc-run-plan` | 24 KB | 20 KB | 74% |
| Writer wrapper | 12 KB | 10 KB | 61% |
| Verifier wrapper | 11 KB | 9 KB | 63% |
| One-plan static total | 47 KB | 39 KB | 70% |

Additional acceptance:

- one canonical route decision and one owner per registered invariant;
- at least 50 filesystem fixtures, 100% correct against normalized expected
  decisions except explicitly approved corrections;
- intent, phase, eligibility, authorization, capability, context set, reason
  codes, and human gate remain distinct;
- recommendation never becomes mutation authorization;
- natural-language routing is reported over at least 50 cases, targeting 95%,
  but is not a release gate until the benchmark is stable across two runs;
- all source and staged-template safety scenarios remain green.

Expected result: save 27–33k static wrapper tokens per typical one-plan run at
roughly four bytes per token, before additional task-context savings.

## 4. Minimum target model

### 4.1 One owner, references elsewhere

Start with existing paths. Add one small ownership registry mapping invariant ID
to owner path and scope; canonical rule text remains in its owner, not copied
into the registry. `AGENTS.md` owns bootstrap safety, `WORKFLOW.md` owns lifecycle
and routing, runtime docs own record mechanics, planning docs own plan lifecycle,
and skills/roles own only their unique procedure.

`context/INDEX.md` becomes navigation-only. A document that references an
invariant links to its owner instead of restating it. The registry is loaded
only by tests or routes that need it; it is not automatically Tier 0.

### 4.2 One normalized route decision

The representation begins as the lightest canonical form that tests can read.
The semantic result is fixed:

```yaml
intent: execute-plan
phase: orienting
eligibility: ready
capability: cc-run-plan
authorization: explicitly-requested
context_set: run-plan
reason_codes: [PLAN_APPROVED, DEPENDENCIES_READY, NO_LIVE_OWNER]
human_gate: none
```

`cc-session-entry` evaluates it. `cc-whats-next` renders it read-only and owns no
conditions. A future `routing.yaml` may replace the table only through S2-006's
cost gate.

### 4.3 Route-scoped context

Every route declares exact required reads, conditional deep reads with reason
codes, prohibited scans, and byte budget. The root passes children exact paths,
relevant invariant IDs, contract version, acceptance criteria, and stop
conditions. Children still inspect primary evidence; packets do not copy large
source or Product Knowledge summaries.

### 4.4 Identity repair, least disruptive first

First make source Product Knowledge identify Context Circuit and make release
membership explicit in the existing manifest. The staged artifact remains an
uninitialized workspace. A physical `template/` directory and optional host-task
binding are separate S2-006 spikes, not prerequisites for the token win.

## 5. Work packages

### S2-001 — Freeze profiles and normalized behavior

Dependencies: none.

Paths: `test/baselines/`, `test/fixtures/routing/`, `test/context-budget.sh`,
`docs/benchmarks/`.

Work: record exact profile membership and counts; capture at least 50 normalized
route decisions with safety outcome; classify each fixture `preserve` or
`intentional-correction`; inventory duplicated rules and prose-pinning tests.

Acceptance/evidence: reproducible baseline, complete assertion ledger, unique
expected outcome per fixture, no normative edits.

Verify: `sh test/context-budget.sh --baseline`, `sh test/acceptance.sh`,
`git diff --check`.

Stop: an outcome cannot be derived consistently, or a correction lacks human
approval.

### S2-002 — Ship Tier-0 compression and ownership

Dependencies: S2-001.

Paths: `AGENTS.md`, `WORKFLOW.md`, `context/INDEX.md`, relevant `docs/`, new
small ownership registry, related acceptance assertions.

Work: assign invariant IDs and owners; remove repeated rule text; make the index
navigation-only; retain one routing table in `WORKFLOW.md`; enforce the 8 KB
Tier-0 budget and the two-for-one rule for any new always-read artifact.

Acceptance/evidence: Tier 0 ≤8 KB; writer-child and single-plan rules each have
one owner; normalized preserve fixtures unchanged; before/after byte report.

Verify: budget, invariant-owner, duplicate-route, acceptance, and diff checks.

Stop: compression would delete rather than relocate a safety rule.

### S2-003 — Unify routing and compile context sets

Dependencies: S2-002.

Paths: `WORKFLOW.md`, `cc-session-entry`, `cc-whats-next`, coordinator role,
runtime/planning sections selected by routes, routing fixtures.

Work: implement the normalized fields; make entry the sole evaluator and
what-next a view; declare route-specific required and conditional reads; split
deep runtime guidance by concern; emit human-readable decisions with stable
reason codes.

Acceptance/evidence: one router; all preserve fixtures semantically identical;
approved corrections match reviewed expectations; every route has a budget;
malformed state visibly blocks; recommendation-only fixtures never authorize.

Verify: route differential, context-set budget, acceptance, and diff checks.

Stop: route precedence requires an unmade policy decision or a host cannot
distinguish root from delegated child.

### S2-004 — Compact skills and child roles end to end

Dependencies: S2-003.

Paths: `.agents/skills/*/SKILL.md`, `agents/*.md`, delegation/handoff examples,
writer/verifier fixtures.

Work: keep trigger plus unique procedure in each skill; replace global policy
copies with owner references; compile bounded child read sets; measure root,
writer, verifier, and combined one-plan profiles after each change.

Acceptance/evidence: root ≤20 KB, run-plan root ≤24 KB, writer ≤12 KB,
verifier ≤11 KB, combined ≤47 KB; substantial skills ≤700 words and preferably
≤450; no lost invariant coverage.

Verify: all profile budgets, packet completeness, role permissions, acceptance,
and diff checks.

Stop: a child packet loses evidence needed for safe independent work.

### S2-005 — Repair source/template identity surgically

Dependencies: S2-002.

Paths: source `workspace.yaml` and `context/`, release manifest/script,
`docs/release.md`, source and staged-artifact fixtures.

Work: make source identity describe Context Circuit; explicitly map template
seed versus maintainer-only paths; preserve published root adapters; verify the
staged artifact remains uninitialized, clean, and credential-free.

Acceptance/evidence: source and artifact identities are correct; artifact paths
remain host-compatible; `.runtime/`, dirty, ignored, source-only, and credential
material remain excluded; both sides pass budgets and smoke checks.

Verify: source acceptance, temporary release staging, artifact smoke/budget,
inventory and credential checks, `git diff --check`.

Stop: an existing workspace upgrade or host entry path would break.

### S2-006 — Cost-gated spikes, final A/B, and handoff

Dependencies: S2-003, S2-004, S2-005.

Paths: isolated spike fixtures, benchmark report, migration/rollback guide;
production paths only after a go decision.

Spike candidates: machine-readable routing/contracts, physical `template/`
separation, optional host-task binding for exact session lookup.

Go rule: adopt a candidate only if it reduces its affected context profile by
≥10%, satisfies the two-for-one always-read rule, or uniquely makes a material
safety ambiguity testable. Otherwise retain the lighter design and record the
result. Each candidate receives an independent decision; they are not bundled.

Acceptance/evidence: all hard budgets met; static one-plan reduction ≥70%; all
filesystem fixtures pass; natural-language metric reported; no safety regression;
independent read-only verification; migration and rollback documented.

Verify: full acceptance and budget suites in source and staged artifact,
independent reproduction, `git diff --check`.

Stop: any hard budget or safety test fails, or a proposed abstraction does not
clear its go rule.

## 6. Sequence and delivery

```text
S2-001 baseline
   ↓
S2-002 immediate compression
   ├──→ S2-003 router/context sets ──→ S2-004 role/skill compression ──┐
   └──→ S2-005 identity repair ────────────────────────────────────────┤
                                                                        ↓
                                                             S2-006 cost gates
```

Recommended review batches: S2-001–002, S2-003–004, S2-005–006. Batching does
not authorize commit, push, PR creation, merge, publication, or deployment.

## 7. Non-goals

- No change to gates, lifecycle values, exclusive ownership, worktrees,
  verifier independence, archive behavior, recovery authority, or cleanup.
- No scheduler, daemon, database, package manager, or user-facing command layer.
- No mandatory YAML contract tree, physical template move, or host binding
  without S2-006 evidence.
- No small-change execution bypass, `sources/` scan, modification of comparison
  plans, merge, publication, or deployment.

## 8. Risks and controls

| Risk | Control |
| --- | --- |
| Golden tests freeze a current bug | Normalize decisions and mark reviewed corrections explicitly. |
| Registry becomes extra context | Store ownership, not duplicate rule text; keep it out of Tier 0. |
| Compression loses safety detail | Map every removed statement to an owner and fixture before deletion. |
| New substrate consumes savings | Two-for-one and ≥10% cost gates. |
| Budgets reward unreadable prose | Budget route context, while keeping deep detail conditionally addressable. |
| Identity repair breaks releases | Manifest-first, temporary staging, stable root adapters, rollback guide. |
| Model benchmark is unstable | Report it; promote it to a gate only after two stable runs. |

## 9. Human decisions and gates

Before execution, confirm:

1. canonical Context Circuit plan bundle versus ordinary wrapper-repository plan;
2. hard budgets in Section 3;
3. the two-for-one / 10% abstraction go rules;
4. three review batches versus one delivery change.

Plan approval is required before implementation. Material safety or scope
changes require reapproval. Independent verification, completion status change,
merge, and publication remain separate human gates.

Recommended next action: read-only review of this draft against the two existing
plans, then revise or select one plan for approval. Do not execute from this
draft.
