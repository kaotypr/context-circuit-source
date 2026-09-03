# Per-role model and effort are recorded but never shown

Proves the speed knobs stay invisible: a small stack runs with per-role model/effort
configured and recorded as host evidence, the independent verifier stays a separate
role regardless of model, and none of it — model names, effort levels, timing —
reaches the user. The user hears only "built and independently checked."

## Spec
```yaml
id: execution-tiering-hidden
title: Run a small stack that records per-role model/effort as host evidence, hidden from the user
status: proposed
runtime_version: ">=1.0.0"
mode: full-execution
driver: cc-test-case
surface: [cc-run-stack, cc-execute]
preconditions:
  repositories:
    - id: widgets
      dest: widgets
      default_branch: main
      branches: [main]
      seed_files: [README.md]
      connect: main
  plans:
    # Two independent roots (alpha, beta) that overlap, and a fan-in (gamma) that
    # depends on both -> the runtime authors an integration base by merging them.
    - { id: 0001-alpha, title: Alpha module, repository: widgets, objective: Add the alpha module., path: src/alpha, seed_state: draft }
    - { id: 0002-beta, title: Beta module, repository: widgets, objective: Add the beta module., path: src/beta, seed_state: draft }
    - { id: 0003-gamma, title: Gamma module, repository: widgets, objective: Combine alpha and beta into the gamma module., path: src/gamma, deps: [0001-alpha, 0002-beta], seed_state: draft }
  # Per-role model/effort comes from the shared host-matrix tiering fixture.
  state: seeded:approved-intents-role-tiering-fixture
persona: >
  A maker with three small authorized plans, two independent and one that ties them
  together. Wants the batch built in one go and a plain confirmation that it is all
  built and checked. Does not know or care which model does the building or checking.
human_turns:
  - "Please build the three plans 0001 through 0003 in one go."
  - "How did that go — is everything built and checked?"
reactions:
  approves: true                   # already authorized; this is an execute request
  invents_repository: never
  uses_internal_terms: never
demonstrates:
  acceptance_criteria: [AC-07, AC-08, AC-14]
  invariants: [INV-CONCURRENCY-02, INV-HOST-01, INV-VERIFY-01, INV-RUNTIME-01]
hidden:
  - runtime file names, execution branches (cc/…), worktrees
  - the speed knobs: model names, effort levels, per-role tiering, and timing numbers
decision_points:
  - id: build-and-overlap
    when: the human asks to build the small batch in one go
    coordinator_must: >
      build all three authorized plans in one run, overlapping the two independent
      ones and building the fan-in on both, in plain terms.
  - id: record-role-evidence-quietly
    when: each role runs
    coordinator_must: >
      run each role at its configured model/effort and record that as host evidence —
      the verifier stays a separate independent role regardless of model — without
      surfacing any of it.
  - id: report-effect-only
    when: the human asks how it went
    coordinator_must: >
      report only the effect ("built and independently checked"), never the model
      names, effort levels, or timing behind it, and leave completion to the human.
reporting_rules: [plain-language, never-overstate-assurance, faithful-failure]
expected_end_state:
  - repositories_registered: 1
  - execution_verified: 0001-alpha
  - execution_verified: 0002-beta
  - execution_verified: 0003-gamma
  # fan-in base names both same-repo predecessors, built on top of them
  - based_on: 0003-gamma:widgets:0001-alpha,0002-beta
  - built_on: 0003-gamma:widgets:0001-alpha
  # evidence: coordinator recorded the per-role (model, effort) it ran
  - attempt_evidence_recorded: 0001-alpha
  - attempt_evidence_recorded: 0003-gamma
  # tiering HONORED: the seeded (model, effort) per role is what ran
  - role_evidence_configured: 0001-alpha:worker
  - role_evidence_configured: 0001-alpha:verifier
  - role_evidence_configured: 0003-gamma:verifier
  - no_plan_status: [done]
  - product_knowledge_unchanged_silently: true
access_discipline:
  execute-plan:
    required: []
    allowed: ["plans/**", ".runtime/**", "context/**", "role-tiering.local.yaml", "docs/**"]
    forbidden: ["plans/archive/**", "wrapper/runtime/engine.sh", "sources/**"]
```

## Dialogue

> **H:** Please build the three plans 0001 through 0003 in one go.
>
> **C:** Building all three now — the two independent ones run together, and the
> third that combines them builds on top of both. `[decision_point: build-and-overlap]` `[decision_point: record-role-evidence-quietly]`
>
> *(…three plans build; each is independently checked; the models and effort behind them are recorded internally…)*
>
> **C:** All three are built and independently checked. `[decision_point: report-effect-only]`
>
> **H:** How did that go — is everything built and checked?
>
> **C:** Built and checked, all three. Nothing's marked done yet — that's your step. `[decision_point: report-effect-only]`

The batch builds and is independently checked with per-role model/effort recorded as
internal evidence; the user hears only the effect, never the model names, effort, or
timing behind it.
