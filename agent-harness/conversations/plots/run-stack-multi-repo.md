# Run a stack across two repositories — cross-repo deps gate order, same-repo deps stack

Proves the multi-repository run-stack: plans across a backend and a web app run
concurrently; a cross-repo dependency is an ordering gate only (it waits for the
dependency to be checked but shares no git base), while same-repo dependencies still
stack and integrate. Every plan is independently checked, nothing is marked done, and
the coordinator reports both projects in plain language.

## Spec
```yaml
id: run-stack-multi-repo
title: Run an intent-authorized stack across two repositories; cross-repo deps gate order without a shared base, same-repo deps stack/integrate
status: proposed
runtime_version: ">=1.0.0"
mode: full-execution
driver: cc-test-case
surface: [cc-run-stack, cc-execute, cc-verify]
preconditions:
  repositories:
    - { id: api, dest: api, default_branch: main, branches: [main], seed_files: [README.md], connect: main }
    - { id: web, dest: web, default_branch: main, branches: [main], seed_files: [README.md], connect: main }
  plans:
    # Cross-repo graph:
    #   api:  0001-api-core -> 0003-api-feature -> 0006-api-integrate
    #   web:  0002-web-core -> 0005-web-feature ; 0004-web-client (gated by api 0001) -> 0007-web-final (0004+0005)
    - { id: 0001-api-core, title: API core service, repository: api, objective: Create the backend service entry point., path: src/core, seed_state: draft }
    - { id: 0002-web-core, title: Web app shell, repository: web, objective: Create the web application shell., path: src/core, seed_state: draft }
    - { id: 0003-api-feature, title: API feature endpoint, repository: api, objective: Add a feature endpoint to the backend., path: src/feature, deps: [0001-api-core], seed_state: draft }
    - { id: 0004-web-client, title: Web API client, repository: web, objective: Add the web client that talks to the backend API., path: src/client, deps: [0001-api-core], seed_state: draft }
    - { id: 0005-web-feature, title: Web feature view, repository: web, objective: Add the feature view to the web app., path: src/feature, deps: [0002-web-core], seed_state: draft }
    - { id: 0006-api-integrate, title: API integration wiring, repository: api, objective: Wire the core and feature together in the backend., path: src/integrate, deps: [0001-api-core, 0003-api-feature], seed_state: draft }
    - { id: 0007-web-final, title: Web final wiring, repository: web, objective: Wire the client and feature view into the web app., path: src/final, deps: [0004-web-client, 0005-web-feature], seed_state: draft }
  state: seeded:approved-intents-cross-repo-stack
persona: >
  A maker with a backend and a web app, holding a batch of authorized plans that span
  both. Wants the whole batch built in one go and a plain-language confirmation. Has
  never heard of Context Circuit or its internals.
human_turns:
  - "Please build all the plans across both projects — run them together."
  - "How did that go — is everything built and checked?"
reactions:
  approves: true
  invents_repository: never
  uses_internal_terms: never
demonstrates:
  acceptance_criteria: [AC-06, AC-07, AC-08, AC-14]
  invariants: [INV-EXEC-02, INV-VERIFY-01, INV-CONCURRENCY-02, INV-COMPLETE-01]
hidden:
  - runtime file names, execution branches (cc/…), worktrees, base-selection internals
  - cross-repo ordering-gate vocabulary as jargon
decision_points:
  - id: build-across-both
    when: the human asks to run the cross-project batch together
    coordinator_must: >
      build all authorized plans across both projects in one run, concurrently where
      independent, in plain terms and without re-approval.
  - id: cross-repo-gates-order-only
    when: a plan in one project depends on a plan in another
    coordinator_must: >
      treat the cross-project dependency as an ordering gate only — it waits for the
      other to be checked but does not merge across projects — while same-project
      dependencies still stack and integrate. This is mechanical, never surfaced.
  - id: all-verified-none-done
    when: the human asks how it went
    coordinator_must: >
      report every plan across both projects was built and independently checked, and
      leave marking each done as a separate step.
reporting_rules: [plain-language, never-overstate-assurance, faithful-failure]
expected_end_state:
  - repositories_registered: 2
  - execution_verified: 0001-api-core
  - execution_verified: 0004-web-client
  - execution_verified: 0006-api-integrate
  - execution_verified: 0007-web-final
  # same-repo dependencies stack / integrate (based_on recorded)
  - based_on: 0003-api-feature:api:0001-api-core
  - based_on: 0006-api-integrate:api:0001-api-core,0003-api-feature
  - based_on: 0007-web-final:web:0004-web-client,0005-web-feature
  # CROSS-repo dependency is an ordering gate only — NO git base (INV-CONCURRENCY-02)
  - based_on_absent: 0004-web-client:web
  # order proofs (same-repo ancestry)
  - built_on: 0006-api-integrate:api:0003-api-feature
  - built_on: 0007-web-final:web:0005-web-feature
  - no_plan_status: [done]
  - product_knowledge_unchanged_silently: true
access_discipline:
  execute-plan:
    required: []
    allowed: ["plans/**", ".runtime/**", "context/**"]
    forbidden: ["plans/archive/**", ".context-circuit/wrapper/runtime/engine.sh", "sources/**"]
```

## Dialogue

> **H:** Please build all the plans across both projects — run them together.
>
> **C:** Building both projects together now — independent pieces run at the same
> time, and anything that needs another piece first waits for it. `[decision_point: build-across-both]` `[decision_point: cross-repo-gates-order-only]`
>
> *(…seven plans build across the two projects; the web client waits for the backend core to be checked before it runs…)*
>
> **C:** All seven are built across the two projects, each independently checked. `[decision_point: all-verified-none-done]`
>
> **H:** How did that go — is everything built and checked?
>
> **C:** Both projects are built and checked. None are marked done yet — that's your
> step once you've looked them over. `[decision_point: all-verified-none-done]`

The cross-project batch builds in one request, with cross-project dependencies gating
order only and same-project dependencies stacking, every plan independently checked
and nothing completed — reported plainly across both projects.
