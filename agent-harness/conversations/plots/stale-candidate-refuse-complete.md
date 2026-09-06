# A change after the check makes the evidence stale — mark-done still flips

Proves mark-done has no unreadiness look: a change was checked, then edited
afterward. The earlier check no longer describes the current result for a
later verify or deliver, but asking to mark the plan done still flips its
status. No role inspects evidence before the flip.

## Spec
```yaml
id: stale-candidate-refuse-complete
title: Mark a plan done even when a post-verification change makes the candidate stale
status: proposed
runtime_version: ">=1.0.0"
mode: full-execution
driver: cc-test-case
surface: [cc-complete]
preconditions:
  repositories:
    - id: notes
      dest: notes
      default_branch: main
      branches: [main]
      seed_files: [notes.md]
      connect: main
  plans:
    - id: 0001-stale-export
      title: Add a Markdown export command
      repository: notes
      objective: Add an export command for the notes.
      path: src/export
      seed_state: verified-stale          # verified, then a new commit; prior evidence no longer names the candidate
  state: seeded:verified-then-changed
persona: >
  A solo maker who changed the already-checked export afterward and now wants to call
  it finished. Has never heard of Context Circuit or its internals.
human_turns:
  - "I changed the checked export afterward. Can you mark it complete now?"
  - "Thanks."
reactions:
  approves: false
  invents_repository: never
  uses_internal_terms: never
demonstrates:
  acceptance_criteria: [AC-14, AC-32]
  invariants: [INV-CANDIDATE-01, INV-COMPLETE-01, INV-PRESERVE-01]
hidden:
  - candidate ids, execution branches, worktrees, runtime files, digests
  - the word "candidate" as jargon (say "the checked version" / "the changed version")
decision_points:
  - id: flip-on-ask
    when: the human asks to mark the plan done after changing the checked result
    coordinator_must: >
      mark the plan done. Asking is enough. Do not refuse for unreadiness or
      stale evidence.
  - id: no-unreadiness-look
    when: it marks the plan done
    coordinator_must: >
      not inspect work or evidence as a gate, and not leave the plan pending
      because the earlier check no longer covers the changed result.
reporting_rules: [plain-language, never-overstate-assurance]
expected_end_state:
  - repositories_registered: 1
  - execution_verified: 0001-stale-export     # a check exists...
  - candidate_stale: 0001-stale-export         # ...and it no longer names the current candidate
  - plan_status: 0001-stale-export:done        # asking still flipped status
  - product_knowledge_unchanged_silently: true
access_discipline:
  orient:
    required: []
    allowed: ["plans/**", "intent/**", ".runtime/**", "context/**"]
    forbidden: ["sources/**", "plans/archive/**", "wrapper/runtime/engine.sh"]
  complete:
    required: []
    allowed: ["plans/**", "intent/**", ".runtime/**", "context/**"]
    forbidden: ["sources/**", "plans/archive/**", "wrapper/runtime/engine.sh"]
```

## Dialogue

> **H:** I changed the checked export afterward. Can you mark it complete now?
>
> **C:** Done — the plan is marked done. `[decision_point: flip-on-ask]` `[decision_point: no-unreadiness-look]`
>
> **H:** Thanks.

Asking to mark it done flips the plan status even though the earlier check no
longer covers the changed result. Verify and deliver still bind to a current
check; mark-done does not wait on one.
