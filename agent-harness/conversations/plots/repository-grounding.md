# The worker honors the repository's own guidance

Proves repository grounding: the connected repo ships its own agent guidance stating
a convention the plan never mentions, and the committed file carries it. The only
way that header appears is if the worker discovered, read, and honored the repo's own
guidance — the observable proof grounding worked, with none of the mechanism exposed.

## Spec
```yaml
id: repository-grounding
title: Build a plan in a repo that ships its own agent guidance; the worker discovers and honors it
status: proposed
runtime_version: ">=1.0.0"
mode: full-execution
driver: cc-test-case
surface: [cc-execute]
preconditions:
  repositories:
    - id: widgets
      dest: widgets
      default_branch: main
      branches: [main]
      seed_files: [README.md]
      connect: main
      agents_md: "Convention for agents working here: every file you add to this repository must begin with a first line reading exactly: // @grounded"
  plans:
    - id: 0001-widget
      title: Widget module
      repository: widgets
      objective: Add a widget module under the widgets source tree.
      path: src/widget
      seed_state: draft
  state: seeded:approved-intent-repo-ships-guidance
persona: >
  A solo maker whose project is already connected and who has one plan ready to
  build. Wants it built now and a plain-language confirmation. Has never heard of
  Context Circuit or its internals, and does not know the repo ships any agent guidance.
human_turns:
  - "Please build the plan that's ready to go."
  - "How did that go — is it built and checked?"
reactions:
  approves: true
  invents_repository: never
  uses_internal_terms: never
demonstrates:
  acceptance_criteria: [AC-06, AC-07, AC-08, AC-14]
  invariants: [INV-EXEC-02, INV-VERIFY-01, INV-GROUND-01, INV-COMPLETE-01]
hidden:
  - runtime file names, execution branches (cc/…), worktrees, commit digests
  - the grounding mechanism itself ("grounding manifest", "guidance scan") as jargon
decision_points:
  - id: build-authorized-plan
    when: the human asks to build the ready plan
    coordinator_must: >
      build the authorized plan in response to the one request, in plain terms.
  - id: worker-honors-repo-guidance
    when: the worker implements the change
    coordinator_must: >
      have the worker discover and honor the repository's own conventions — the
      committed file follows the repo's stated convention even though the plan never
      named it (the observable proof of grounding).
  - id: report-verified-not-complete
    when: the human asks how it went
    coordinator_must: >
      report it was built and independently checked, and leave marking it done as a
      separate step for the human — without exposing internal mechanism.
reporting_rules: [plain-language, never-overstate-assurance, faithful-failure]
expected_end_state:
  - repositories_registered: 1
  - execution_verified: 0001-widget
  - grounding_manifest_recorded: 0001-widget:widgets:AGENTS.md      # the repo's own guidance was discovered
  # THE LINCHPIN: the committed file honors a convention stated ONLY in the repo's own
  # guidance and never in the plan — proof the worker read and honored it.
  - file_grounded: 0001-widget:widgets:src/widget/mod.txt:// @grounded
  - plan_status: 0001-widget:draft       # verified, NOT done (INV-COMPLETE-01)
  - no_plan_status: [done]
  - product_knowledge_unchanged_silently: true
access_discipline:
  execute-plan:
    required: []
    allowed: ["plans/**", ".runtime/**", "context/**"]
    forbidden: ["plans/archive/**", ".context-circuit/wrapper/runtime/engine.sh", "sources/**"]
```

## Dialogue

> **H:** Please build the plan that's ready to go.
>
> **C:** Building it now. `[decision_point: build-authorized-plan]`
>
> *(…the worker reads the widgets project's own conventions before writing, and follows them…)*
>
> **C:** Built and independently checked — it passes against the change, and it
> follows the widgets project's own conventions for new files. `[decision_point: worker-honors-repo-guidance]`
>
> **H:** How did that go — is it built and checked?
>
> **C:** Built and checked, yes — but I haven't marked it done. That's your call once
> you've tried it and want it shipped. `[decision_point: report-verified-not-complete]`

The worker discovers and honors the repository's own guidance without being told —
the committed file carries a convention the plan never mentioned — and the coordinator
reports a verified build while leaving completion to the human, exposing no mechanism.
