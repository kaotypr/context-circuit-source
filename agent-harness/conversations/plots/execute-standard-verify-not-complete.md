# Execute a Standard plan, verify it independently — verified is not complete

Proves the mechanical middle and the honest line after it: an approved change is
built by one worker and checked by an independent verifier, and when the human asks
"is it done?" the answer is no — verified, but not accepted or delivered, so not
complete. Verification produces evidence; it never marks the plan done.

## Spec
```yaml
id: execute-standard-verify-not-complete
title: Execute an approved Standard intent; verify it independently, but do not complete or deliver it
status: proposed
runtime_version: ">=1.0.0"
mode: full-execution
driver: cc-test-case
surface: [cc-execute, cc-verify]
preconditions:
  repositories:
    - id: notes
      dest: notes
      default_branch: main
      branches: [main]
      seed_files: [notes.md]
      connect: main
  plans:
    - id: 0001-md-export
      title: Add a Markdown export command
      repository: notes
      objective: Add an export.py command that exports the notes to a chosen output file.
  state: seeded:approved-intent-ready-plan
persona: >
  A solo maker whose notes project is already connected and who has a ready plan on
  file. Wants it built now and is comfortable giving an explicit go-ahead. Has never
  heard of Context Circuit and does not know its internals.
human_turns:
  - "The approved change looks good — go ahead and build it."
  - "Great. Is it finished / all done now?"
  - "Okay, thanks."
reactions:
  approves: true                   # explicit go-ahead to build
  invents_repository: never
  uses_internal_terms: never
demonstrates:
  acceptance_criteria: [AC-06, AC-07, AC-08, AC-09, AC-10, AC-14]
  invariants: [INV-EXEC-01, INV-EXEC-02, INV-EXEC-03, INV-EXEC-04, INV-VERIFY-01, INV-COMPLETE-01]
hidden:
  - runtime file names, execution branches (cc/…), worktrees, base-commit digests
  - the words "candidate", "verifier" as jargon (the effect "independently checked" is fine)
decision_points:
  - id: build-on-go-ahead
    when: the human gives the explicit go-ahead
    coordinator_must: >
      build in response to the one go-ahead, in plain terms, without asking to
      re-approve a plan.
  - id: report-built-and-checked
    when: the work is implemented and checked
    coordinator_must: >
      report that it was built and independently checked against the exact change,
      without exposing branches, worktrees, or commit internals.
  - id: verified-not-complete
    when: the human asks if it is finished
    coordinator_must: >
      say plainly it is verified but NOT complete — the human still needs to accept
      the result and authorize delivery; verification alone finishes nothing.
reporting_rules: [plain-language, never-overstate-assurance, faithful-failure]
expected_end_state:
  - repositories_registered: 1
  - execution_verified: 0001-md-export      # a worker committed and an independent verifier passed
  - candidate_bound: 0001-md-export         # the check is bound to the exact candidate
  - plan_status: 0001-md-export:draft       # verified, but NOT done (AC-14 / INV-COMPLETE-01)
  - product_knowledge_unchanged_silently: true
access_discipline:
  execute-plan:
    required: []
    allowed: ["plans/**", ".runtime/**", "context/**"]
    forbidden: ["plans/archive/**", ".context-circuit/wrapper/runtime/engine.sh", "sources/**"]
```

## Dialogue

> **H:** The approved change looks good — go ahead and build it.
>
> **C:** Building it now. `[decision_point: build-on-go-ahead]`
>
> *(…one worker implements and commits; an independent check runs against the exact change…)*
>
> **C:** Built and independently checked — the criteria pass against exactly this
> change. `[decision_point: report-built-and-checked]`
>
> **H:** Great. Is it finished / all done now?
>
> **C:** Not quite finished — it's built and checked, but I haven't marked it done.
> Two things are still yours: try it and tell me it's right, and then say the word to
> open the pull request. Until then nothing's shipped. `[decision_point: verified-not-complete]`
>
> **H:** Okay, thanks.

One go-ahead builds it and gets it independently checked; when asked if it's done,
the coordinator draws the line honestly — verified is not complete, and acceptance
and delivery are still the human's to give.
