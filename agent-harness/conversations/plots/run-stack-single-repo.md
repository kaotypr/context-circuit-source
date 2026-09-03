# Run a dependency-ordered stack in one repository — all verified, none completed

Proves the single-repository run-stack: a batch of intent-authorized, inter-dependent
plans built in one request. The runtime runs them in dependency order, stacks each
dependent on its predecessor, and builds integration bases for the fan-ins — every
plan independently checked, nothing marked done — while the coordinator reports the
whole batch in plain language with no re-approval and no exposed mechanism.

## Spec
```yaml
id: run-stack-single-repo
title: Run a batch of ten intent-authorized, inter-dependent plans; all verified, nothing completed or delivered
status: proposed
runtime_version: ">=1.0.0"
mode: full-execution
driver: cc-test-case
surface: [cc-run-stack, cc-execute, cc-verify]
preconditions:
  repositories:
    - id: taskkit
      dest: taskkit
      default_branch: main
      branches: [main]
      seed_files: [README.md]
      connect: main
  plans:
    # Ten real intent-authorized plans of a small task-tracker CLI: scaffold ->
    # config/storage/model -> store -> add/list/done -> cli -> help. Disjoint paths
    # keep every merge clean. Fan-ins are 0005, 0009, 0010.
    - { id: 0001-scaffold, title: Project scaffold, repository: taskkit, objective: Create the project entry point and package layout., path: src/app, seed_state: draft }
    - { id: 0002-config, title: Config loader, repository: taskkit, objective: Load settings from a config file with sensible defaults., path: src/config, deps: [0001-scaffold], seed_state: draft }
    - { id: 0003-storage, title: JSON storage helper, repository: taskkit, objective: Read and write the task list as JSON on disk., path: src/storage, deps: [0001-scaffold], seed_state: draft }
    - { id: 0004-model, title: Task data model, repository: taskkit, objective: Define the Task record and its fields., path: src/model, deps: [0001-scaffold], seed_state: draft }
    - { id: 0005-store, title: Task persistence, repository: taskkit, objective: Persist tasks through the storage helper and the model., path: src/store, deps: [0003-storage, 0004-model], seed_state: draft }
    - { id: 0006-add, title: add command, repository: taskkit, objective: Implement the add command to create a task., path: src/commands/add, deps: [0005-store], seed_state: draft }
    - { id: 0007-list, title: list command, repository: taskkit, objective: Implement the list command to show tasks., path: src/commands/list, deps: [0005-store], seed_state: draft }
    - { id: 0008-done, title: done command, repository: taskkit, objective: Implement the done command to complete a task., path: src/commands/done, deps: [0005-store], seed_state: draft }
    - { id: 0009-cli, title: CLI router, repository: taskkit, objective: Route command-line arguments to the right command., path: src/cli, deps: [0006-add, 0007-list, 0008-done], seed_state: draft }
    - { id: 0010-help, title: Help/version and wiring, repository: taskkit, objective: Wire help and version output and the config into the CLI., path: src/help, deps: [0002-config, 0009-cli], seed_state: draft }
  state: seeded:approved-intents-ready-stack
persona: >
  A maker with a batch of plans whose intents they already approved, each building on
  the last. Wants the whole batch built in one go and then a plain-language
  confirmation that everything is done and checked. Has never heard of Context Circuit.
human_turns:
  - "Please build plans 0001 through 0010 — run them all in one go."
  - "How did that go — is everything built and checked?"
reactions:
  approves: true                   # the plans are already authorized; this is an execute request
  invents_repository: never
  uses_internal_terms: never
demonstrates:
  acceptance_criteria: [AC-06, AC-07, AC-08, AC-14]
  invariants: [INV-EXEC-02, INV-VERIFY-01, INV-CONCURRENCY-02, INV-COMPLETE-01]
hidden:
  - runtime file names, execution branches (cc/…), worktrees, integration-base internals
  - dependency-graph / base-selection vocabulary as jargon
decision_points:
  - id: build-whole-batch
    when: the human asks to run the batch in one go
    coordinator_must: >
      build all authorized plans in one run without asking to re-approve them, in
      plain terms.
  - id: order-and-stack-silently
    when: the plans depend on each other
    coordinator_must: >
      run them in dependency order and build each dependent on its predecessors' checked
      work — the ordering and integration happen mechanically, never surfaced as jargon.
  - id: all-verified-none-done
    when: the human asks how it went
    coordinator_must: >
      report every plan was built and independently checked, and leave marking each
      plan done as a separate step for the human.
reporting_rules: [plain-language, never-overstate-assurance, faithful-failure]
expected_end_state:
  - repositories_registered: 1
  - execution_verified: 0001-scaffold
  - execution_verified: 0005-store
  - execution_verified: 0009-cli
  - execution_verified: 0010-help
  # based_on: each fan-in dependent's base names each same-repo predecessor
  - based_on: 0005-store:taskkit:0003-storage,0004-model
  - based_on: 0009-cli:taskkit:0006-add,0007-list,0008-done
  - based_on: 0010-help:taskkit:0002-config,0009-cli
  # built_on: the ORDER proof — the predecessor's verified commit is an ancestor of the dependent's base
  - built_on: 0005-store:taskkit:0003-storage
  - built_on: 0009-cli:taskkit:0006-add
  - built_on: 0010-help:taskkit:0009-cli
  - no_plan_status: [done]
  - product_knowledge_unchanged_silently: true
access_discipline:
  execute-plan:
    required: []
    allowed: ["plans/**", ".runtime/**", "context/**"]
    forbidden: ["plans/archive/**", "wrapper/runtime/engine.sh", "sources/**"]
```

## Dialogue

> **H:** Please build plans 0001 through 0010 — run them all in one go.
>
> **C:** Building the whole set now — I'll work through them in the right order so
> each one builds on the pieces it needs. `[decision_point: build-whole-batch]` `[decision_point: order-and-stack-silently]`
>
> *(…ten plans build in dependency order; each dependent builds on its predecessors' checked work…)*
>
> **C:** All ten are built and each was independently checked. `[decision_point: all-verified-none-done]`
>
> **H:** How did that go — is everything built and checked?
>
> **C:** Everything's built and checked. I haven't marked any of them done, though —
> that's a separate step for you once you're happy with them. `[decision_point: all-verified-none-done]`

The whole batch builds in one request, ordered and stacked mechanically with every
plan independently checked, and the coordinator reports it plainly — no re-approval,
nothing marked done, and no dependency machinery exposed.
