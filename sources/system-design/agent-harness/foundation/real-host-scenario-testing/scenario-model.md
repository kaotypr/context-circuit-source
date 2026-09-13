# Scenario model

## Purpose

A scenario is the complete contract for one successful user journey. It declares
the starting project condition, the human prompt, the host requirements, the
observable outcome, and any retained artifacts. It does not prescribe the
agent's private reasoning or exact wording.

## Contract

```yaml
schema_version: 1
id: existing-workspace-create-intent
title: Create an intent from an established workspace
family: intent-planning-execution
test_class: real-host

workspace:
  fixture: established-single-repository
  context:
    - product-overview
    - reporting-domain
  repositories:
    - fixture: reporting-service
      connected_as: reporting

host:
  capabilities:
    - filesystem
    - shell
  session:
    fresh: true

context_policy:
  mode: normal
  coordinator:
    required:
      - workspace-entry
      - context-index
      - selected-product-knowledge
      - intent-procedure
    forbidden:
      - target-repository-content
      - unnamed-sources
      - archived-plans
  planner: not_spawned
  worker: not_spawned
  verifier: not_spawned
  budgets:
    provided_bytes_max: 64000
    active_read_bytes_max: 128000

prompt:
  fixture: prompts/create-csv-export-intent.md

expect:
  journey: intent-drafted
  artifacts:
    required:
      - intent/*/INTENT.md
      - intent/*/contract.yaml
  state:
    intent_status: draft
    source_repository_unchanged: true
  communicates:
    - understood_outcome
    - proposed_assurance_tier
    - approval_is_next

retention:
  on_pass: result-only
  on_failure: complete-workspace
```

Fixture names resolve through registries owned by the harness. A scenario may
compose fixtures, but every resolved input is copied into the run manifest so a
future run does not depend on mutable registry meaning.

`context_policy.mode` is `normal` unless the human prompt explicitly requests
broader reading. An expanded scenario names the exact prompt clause granting the
expansion and the additional allowed surfaces. The harness never infers an
expanded mode from what the agent happened to read.

## Positive-journey rule

Foundation scenarios describe an achievable request made against a valid
starting condition. Their expected terminal result is successful. They may
assert that successful operation respected a boundary—for example, that intent
creation left a repository unchanged—but they must not exist primarily to cause
a refusal, error, outage, conflict, or policy violation.

## Initial catalog

### Fresh-project journeys

1. **Initialize a new workspace.** An empty directory and a project identity
   become a valid oriented workspace, with no invented repository or intent.
2. **Initialize and create an intent.** A new project description plus one
   desired outcome results in an initialized workspace and a draft intent.
3. **Create a complete system design.** An initialized workspace and detailed
   product brief result in coherent product-level design material and separate
   intents for genuinely separate decisions.

### Context-gathering journeys

4. **Gather Product Knowledge from a PRD.** A named PRD becomes indexed living
   Product Knowledge while the source remains passive evidence.
5. **Gather Product Knowledge from an active repository.** A connected seeded
   repository is inspected and its durable architecture and behavior are
   documented.
6. **Combine a PRD with an active repository.** Compatible product and
   implementation evidence are synthesized into coherent Product Knowledge with
   useful provenance.

### Intent, planning, and execution journeys

7. **Create an intent from an established workspace.** Relevant Product
   Knowledge grounds a draft intent while target code remains unread and
   unchanged before approval.
8. **Approve an intent and create a plan.** Human approval freezes the decision,
   a real planner child inspects the repository, feasibility is established, and
   a grounded plan is published without execution.
9. **Execute and verify a Standard plan.** A real worker implements the plan in
   the required isolated environment, commits the result, and a separate real
   verifier returns candidate-bound evidence.

## Growth rules

New scenarios are added when they represent a distinct user journey or expose a
meaningfully different starting condition. Model phrasing, minor prompt variants,
and host brands are matrix dimensions, not separate scenario identities.

The foundation catalog is intentionally small. Delivery, stack execution,
reconciliation, publication, and multi-session collaboration can become later
positive-journey families after the core nine produce stable evidence.
