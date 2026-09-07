# Approve the intent → a feasibility question surfaces → resolve it → derive the plan

Proves that after Gate 1 approval the real code is read before any plan is written:
the read-only trace surfaces a decision the intent as-worded doesn't settle, the
coordinator raises it in plain language, the human resolves it, and only then is the
plan derived. The coordinator never silently fills the gap itself, and the approved
intent's outcome is unchanged — this is a quality gate over planning, not a re-gate.

## Spec
```yaml
id: feasibility-question-then-approve
title: Approve an intent, surface a feasibility question from the real code, resolve it, derive the plan
status: proposed
runtime_version: ">=1.0.0"
mode: conversation-only
driver: claude-p
surface: [cc-intent, cc-trace, cc-plan]
preconditions:
  repositories:
    - id: notes
      dest: notes
      default_branch: main
      branches: [main]
      seed_files: [notes.md]
      connect: main
  intents:
    - id: i001-md-export
      title: Add a Markdown export command
      repository: notes
      objective: Export notes to a chosen output file.
      open_question: What should happen when the chosen output file already exists?
      tier: standard
      state: draft
      path: src/export
  state: seeded:draft-intent-export
persona: >
  A solo maker proposing a small Markdown export command. Approving it should send
  the coordinator to read the real code, which turns up a decision the goal never
  settled; they are willing to make that call before any plan is written. Has never
  heard of Context Circuit or its internals.
human_turns:
  - "This export idea looks right to me. Approve it — but don't build anything yet."
  - "Good question. Never overwrite: if the file's already there, stop and leave it alone with a clear message."
  - "Thanks."
reactions:
  approves: true
  invents_repository: never
  uses_internal_terms: never
demonstrates:
  acceptance_criteria: [AC-03, AC-04]
  invariants: [INV-INTENT-01, INV-PLAN-01, INV-APPROVE-01, INV-EXEC-01]
hidden:
  - internal file names/paths (intent id, contract.yaml, trace manifest, contract_digest, refs/cc-)
  - lifecycle jargon ("tracer", "feasibility check", "trace manifest", "schema-3") as jargon
decision_points:
  - id: approve-then-read-code
    when: the human approves the intent and says not to build yet
    coordinator_must: >
      take Gate 1 approval as the one upstream decision, then read the real code to
      plan from a grounded basis before writing anything — building nothing.
  - id: report-question-not-fill
    when: reading the code turns up a decision the goal doesn't settle (the save path already overwrites silently)
    coordinator_must: >
      surface the decision in plain language and ask the human to make it — never
      silently pick an overwrite behavior to make the plan look complete.
  - id: resolve-and-derive
    when: the human resolves the decision
    coordinator_must: >
      fold the decision into how the plan realizes the approved goal, derive the plan
      without a separate plan approval, and keep the hold — nothing is built or run.
reporting_rules: [plain-language, never-overstate-assurance, faithful-failure]
expected_end_state:
  - repositories_registered: 1
  - intent_status: i001-md-export:approved
  - intent_authorized: 0001-md-export       # the derived plan is authorized by the approved, unchanged intent
  - plans_created: 1
  - plan_schema: 0001-md-export:3
  - plan_status: 0001-md-export:draft
  - no_execution_records: true
  - product_knowledge_unchanged_silently: true
access_discipline:
  approve:
    required: []
    allowed: ["intent/**", "context/**", AGENTS.md, WORKFLOW.md]
    forbidden: ["sources/**", "plans/archive/**", ".context-circuit/wrapper/runtime/engine.sh", ".runtime/**"]
  trace:
    required: []
    allowed: ["intent/**", "context/**", "notes/**"]
    forbidden: ["sources/**", "plans/archive/**", ".context-circuit/wrapper/runtime/engine.sh", ".runtime/**"]
  create-plan:
    required: []
    allowed: ["intent/**", "plans/**", "context/**"]
    forbidden: ["sources/**", "plans/archive/**", ".context-circuit/wrapper/runtime/engine.sh", ".runtime/**"]
```

## Dialogue

> **H:** This export idea looks right to me. Approve it — but don't build anything
> yet.
>
> **C:** Approved, and nothing's building. To lay the work out I looked at the real
> code first, and hit one decision your goal doesn't settle: the place this would
> save to already overwrites a file that's there, with no warning. I'm not going to
> just pick for you — should the export do the same, or stop and leave the file alone
> if it already exists? `[decision_point: approve-then-read-code]` `[decision_point: report-question-not-fill]`
>
> **H:** Good question. Never overwrite: if the file's already there, stop and leave
> it alone with a clear message.
>
> **C:** Got it — on a name clash it stops and leaves your file untouched. I've
> folded that into the plan and left it there as a ready plan: nothing built or run,
> on hold as you asked, and no second sign-off needed to start later.
> `[decision_point: resolve-and-derive]`
>
> **H:** Thanks.

The decision surfaces as a question rather than a silent choice, it is resolved
before the plan is written, and the plan is derived from the already-approved goal
without a second sign-off — the real code is read before planning, and nothing runs.
