# Sol 5.6 Plan v3 — Minimize controllable context, prove routing, preserve live work

Status: **draft — review required; no implementation or delivery authorized**

Scope: Context Circuit wrapper instructions, routing, roles, skills, tests,
release template, and upgrade compatibility. This comparison artifact is not a
canonical `plan.yaml`; acceptance must precede conversion or execution.

## 1. Thesis

Optimize what the wrapper actually controls: selected files, tool-read output,
delegation packets, handoffs, and repeated evidence discovery. Measure those
deterministically in CI; validate real model tokens through a host-neutral trace
protocol when a host exposes usage. Never pretend the wrapper controls system
prompt injection, conversation replay, tokenizer choice, or prompt caching.

Ship certain compression first. Add machinery only when it proves a net win.

Rules:

1. Safety semantics remain unchanged and semantically tested.
2. A new always-read artifact removes at least 2× its bytes elsewhere.
3. A heavier mechanism ships only if it reduces an affected lifecycle scenario
   by ≥10% or uniquely makes a material safety property testable.
4. Existing workspaces and in-flight plans upgrade without silent state loss.

## 2. Evidence and baseline

Evidence: current instructions, Product Knowledge, workflow/runtime/planning
contracts, roles, relevant skills, acceptance/release code, the named Codex
discussion, `Sol-5.6-plan-v2.md`, and `Opus-4.8-plan-v3.md`. No `sources/` file
was read; comparison plans remain untouched and are evidence, not authority.

| Static profile | Current bytes | Hard | Stretch |
| --- | ---: | ---: | ---: |
| Tier 0 / every-session | 22,349 | 8 KB | 6 KB |
| Full root entry, including entry skill and coordinator | 76,456 | 20 KB | 16 KB |
| Root entering `cc-run-plan` | 93,873 | 24 KB | 20 KB |
| Writer wrapper | 30,387 | 12 KB | 10 KB |
| Verifier wrapper | 30,108 | 11 KB | 9 KB |
| One-plan static total | 154,368 | 47 KB | 39 KB |

Duplication baseline: writer-child policy in 12 files, single-plan entry in 7,
and no-global-current-session policy in 4. The source checkout also identifies
itself as an uninitialized template instead of the Context Circuit product.

## 3. What gets measured

### 3.1 Deterministic context ledger — mandatory CI gate

For each scenario, record exact selected paths, revisions/digests, bytes, words,
packet bytes, handoff bytes, and optional tokenizer count. Count only material
the wrapper selects or emits. This is portable, credential-free, and repeatable.

Scenarios:

- `L0`: fresh read-only orientation;
- `L1`: one approved plan, one writer, one verifier, no resume;
- `L2`: three sequential tasks sharing the writer, verifier, one root resume;
- `U1`: workspace created on the old wrapper, live plan/lease, upgrade, resume.

Hard lifecycle budgets: `L1` wrapper material ≤47 KB; `L2` ≤59 KB, including a
resume delta ≤12 KB. No task iteration within one session reloads Tier 0 through
wrapper-directed reads.

### 3.2 Host usage trace — publication evidence, not fake portability

Define a provider-neutral evidence shape:

```yaml
host: codex
model: <exact model>
reasoning: <setting>
scenario: L1
runs: 3
median_input_tokens: <reported>
median_cached_input_tokens: <reported-or-unknown>
```

A host adapter or human-run benchmark supplies usage; core shell tests never
spawn an agent, store credentials, or invent unavailable token counts. Compare
the same model/settings over three runs. Publication target: median input-token
reduction ≥50% for `L1` and `L2`, with no routing regression. If the host cannot
report usage, publication requires an explicit human waiver; deterministic
ledger budgets still gate implementation.

### 3.3 Two-stage routing clarity

Tier 0 cannot decide facts it has not read. Therefore measure:

1. **Probe decision:** from request + Tier 0, choose the correct context set or
   blocking question.
2. **Action decision:** after that bounded evidence, choose the final capability,
   eligibility, authorization, reason codes, and human gate.

At least 50 filesystem fixtures must achieve 100% probe and action accuracy for
`preserve` cases. Reviewed `intentional-correction` cases use their new expected
decision. A separate reviewer sees the same minimal contract and state but not
the expected answer; it is not deprived of the routing rules it must apply.
Natural-language routing targets 95% over 50 cases and remains reported until
stable across two runs.

## 4. Minimum design

### 4.1 One owner per rule

Keep existing paths initially. A small registry maps invariant ID → owner path →
scope, without copying canonical text and without joining Tier 0. `AGENTS.md`
owns bootstrap safety; `WORKFLOW.md` owns lifecycle and the sole routing table;
runtime/planning docs own deep mechanics; skills and roles own only unique
procedure. `context/INDEX.md` becomes navigation-only.

### 4.2 One route contract, representation cost-gated

`cc-session-entry` alone evaluates; `cc-whats-next` only renders. The normalized
decision keeps these independent:

```text
intent | phase | probe/context_set | eligibility | capability
authorization | reason_codes | human_gate
```

Start with the smallest testable table. YAML/schema promotion is optional and
must clear the rent rule.

### 4.3 Delta-safe continuation

Handoffs do not copy resolved documents. They carry a compact context receipt:

```yaml
wrapper_version: <version>
route_decision_digest: <digest>
context_set: run-plan
references:
  - path: context/PROJECT.md
    revision: <digest-or-commit>
invariants: [PLAN_APPROVED, EXCLUSIVE_WORKTREE]
```

Resume validates the receipt and reloads only changed, missing, or newly needed
evidence. Primary evidence remains authoritative; a receipt is a cache key, not
truth. Missing/contradictory receipts block unsafe mutation rather than guessing.

### 4.4 Compatible identity repair

Fix source identity and release membership manifest-first. New sessions record
optional `wrapper_version`; absent versions remain valid legacy records and
enter compatibility preflight. Classify changes as compatible, migration-needed,
or blocked. A physical `template/` move and host-task binding remain independent
cost-gated spikes.

## 5. Work packages

### F-001 — Freeze ledgers, routes, and upgrade state

Dependencies: none. Paths: `test/baselines/`, `test/fixtures/routing/`,
`test/fixtures/upgrades/`, `test/context-ledger.sh`, `docs/benchmarks/`.

Work: freeze static profiles and L0/L1/L2/U1 ledgers; capture ≥50 normalized
probe/action decisions, each `preserve` or approved correction; inventory rule
duplication and prose assertions.

Accept/evidence: reproducible ledgers, complete assertion map, no normative edit.
Verify: ledger baseline, current acceptance, `git diff --check`.
Stop: an expected decision is inconsistent or a correction lacks human review.

### F-002 — Deliver Tier-0 compression and ownership

Dependencies: F-001. Paths: `AGENTS.md`, `WORKFLOW.md`, `context/INDEX.md`,
relevant docs, small ownership registry, tests.

Work: assign owners/IDs, delete duplicate statements, make the index navigational,
retain one routing table, enforce 8 KB and 2× rules.

Accept/evidence: Tier 0 ≤8 KB; duplicated policies each one owner; L0 probe
accuracy unchanged; before/after ledger. Verify: budget, ownership, routing,
acceptance, diff checks. Stop: compression would remove a safety rule.

### F-003 — Implement two-stage routing and context receipts

Dependencies: F-002. Paths: `WORKFLOW.md`, entry/what-next skills, coordinator,
delegation/handoff contracts, selected runtime guidance, routing/resume fixtures.

Work: implement probe then action decisions; route-scoped reads and budgets;
context receipts with digest validation; recommendation/authorization separation.

Accept/evidence: one evaluator; 100% preserve fixtures; approved corrections
match; malformed state blocks; resume delta ≤12 KB; no stale receipt overrides
primary evidence. Verify: semantic differential, packet/receipt, budget,
acceptance, diff checks. Stop: precedence needs a new policy decision.

### F-004 — Compact skills/roles and prove lifecycle budgets

Dependencies: F-003. Paths: `.agents/skills/*/SKILL.md`, `agents/*.md`, deep
contract sections, lifecycle fixtures and ledger.

Work: keep triggers and unique procedures; compile bounded child context;
replace prose-presence tests with invariant/decision tests; measure after each
change.

Accept/evidence: all static limits met; L1 ≤47 KB; L2 ≤59 KB; substantial skill
≤700 words, preferably ≤450; no lost safety coverage. Verify: full ledger,
semantic negative fixtures, acceptance, diff checks. Stop: child independence or
required evidence is weakened.

### F-005 — Repair identity and verify live-workspace compatibility

Dependencies: F-002, F-003. Paths: source `workspace.yaml`/`context/`, runtime
session schema, release manifest/script/docs, U1 fixture.

Work: source identifies Context Circuit; artifact remains uninitialized;
introduce optional wrapper version and compatibility matrix; exercise old live
plan/lease → new wrapper → safe resume without rewriting old evidence.

Accept/evidence: identities correct; artifact clean and host-compatible; U1
preserves lease, worktree, handoff, and canonical statuses; ambiguous migration
requires a human gate. Verify: source/artifact acceptance, temporary staging,
U1 recovery, credential/inventory, diff checks. Stop: an existing workspace or
host entry path would silently break.

### F-006 — Cost-gated candidates, host trace, and final handoff

Dependencies: F-004, F-005. Paths: isolated spike fixtures, benchmark evidence,
migration/rollback guide; production paths only after separate go decisions.

Candidates: machine-readable router/schema, physical `template/`, host-task
binding. Apply 2×/10% rules independently. Run three comparable host traces when
usage is available; independent verifier reproduces deterministic evidence.

Accept/evidence: hard ledgers pass; static reduction ≥70%; probe/action fixtures
pass; host median reduction ≥50% or explicit publication waiver; no safety
regression; rollback documented. Stop: any hard gate fails or a candidate does
not pay rent.

## 6. Sequence, non-goals, and gates

```text
F-001 → F-002 → F-003 → F-004 ─┐
                 └────→ F-005 ─┴→ F-006
```

Batches: F-001–002, F-003–004, F-005–006. No batch authorizes commit, push, PR,
merge, publication, deployment, or canonical status change.

Non-goals: changing gates/lifecycles/leases/worktrees/verifier independence;
adding a scheduler, database, package runtime, user-facing CLI, small-change
bypass, `sources/` scan, or mandatory heavy contract/template/host machinery.

Human decisions before execution: canonical plan bundle versus ordinary wrapper
plan; static/lifecycle budgets; abstraction rent rules; host-trace publication
gate and waiver policy; three batches versus one delivery change.

Recommended next action: read-only comparison review, resolve those decisions,
then select or revise one draft. This file remains unapproved.
