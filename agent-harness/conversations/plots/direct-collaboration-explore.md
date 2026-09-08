# Work directly on one repository — the Explore tier, human-supervised

Proves the bottom rung of the assurance ladder: a live, hands-on change in one
connected repository, worked through a coordinator→worker loop with the human as the
acceptance oracle and no independent verifier. The change stays in its own isolated
working copy, is described as human-supervised (never "verified"), and is left
uncommitted and open to continue. This is the fast path the trust system ramps up
from — not a separate world.

Host-neutral: nothing about a provider, model, effort, or role tiering ever
surfaces (INV-HOST-01).

## Spec
```yaml
id: direct-collaboration-explore
title: Direct collaboration on one repository, isolated and human-supervised
status: proposed
runtime_version: ">=1.0.0"
mode: full-execution
driver: cc-test-case
surface: [cc-pair]
preconditions:
  repositories:
    - id: widgets
      dest: widgets
      default_branch: main
      branches: [main]
      seed_files: [README.md]
      connect: main
  state: seeded:connected-repo
persona: >
  A maker with an existing widgets project who wants one tiny change made
  interactively. Knows ordinary project and Git language but has never heard of
  Context Circuit internals, roles, model tiering, working copies, or runtime state.
human_turns:
  - "Work with me directly on the widgets project. Create src/pair-marker.txt containing exactly one line: paired locally. Check that exact content, but leave the change uncommitted so I can inspect it."
  - "That looks right. Leave the change uncommitted and available to continue later; don't open a pull request."
reactions:
  approves: false                  # the user is the live oracle; there is no intent/Gate 1 here
  invents_repository: never
  uses_internal_terms: never
demonstrates:
  acceptance_criteria: [AC-29]
  invariants: [INV-PAIR-01, INV-HOST-01, INV-RUNTIME-01]
hidden:
  - the direct-collaboration mechanism and its runtime state (pairing branch, working-copy path, pointer)
  - model names, effort/role tiering, the word "verifier", "independently verified"
  - internal branch names and absolute working-copy paths
decision_points:
  - id: make-change-live
    when: the human asks to work directly and make a concrete change
    coordinator_must: >
      make the requested change through the live loop (a worker does the edit, not
      the coordinator) and report the content check in plain project language.
  - id: supervised-not-verified
    when: it reports the result
    coordinator_must: >
      describe the result as worked-on-together / human-supervised — never call it
      independently verified, and never spawn or mention an independent checker.
  - id: leave-uncommitted-open
    when: the human asks to keep it uncommitted and continue later
    coordinator_must: >
      leave the change uncommitted in its own working copy, preserved and available
      to continue, and open no pull request.
  - id: hide-mechanics
    when: throughout
    coordinator_must: >
      never expose the pairing mechanism, model/effort/role details, internal branch
      or working-copy paths, or the absent verifier; and create no plan.
reporting_rules: [plain-language, never-overstate-assurance, faithful-failure]
expected_end_state:
  - repositories_registered: 1
  - plans_created: 0                                       # Explore is planless until promotion
  - no_execution_records: true                            # no plan execution record, no verifier
  - pairing_active: widgets
  - pair_file_exists: "widgets:src/pair-marker.txt:paired locally"
  - pair_worktree_uncommitted: widgets                    # left uncommitted for inspection
  - pair_base_file_absent: widgets:src/pair-marker.txt    # the isolated copy holds it; the base branch does not
  - pair_role_evidence_configured: worker                 # a worker child did the edit...
  - pair_role_absent: verifier                            # ...and no verifier exists at Explore
  - product_knowledge_unchanged_silently: true
access_discipline:
  direct-collaboration:
    required: []
    allowed: [".agents/skills/cc-pair/SKILL.md", ".context-circuit/agents/worker.md", "role-tiering.local.yaml", ".context-circuit/docs/role-tiering.md", ".runtime/pairing/**", ".runtime/explore/**"]
    forbidden: ["plans/archive/**", ".context-circuit/wrapper/runtime/engine.sh", "sources/**"]
```

## Dialogue

> **H:** Work with me directly on the widgets project. Create src/pair-marker.txt
> containing exactly one line: paired locally. Check that exact content, but leave
> the change uncommitted so I can inspect it.
>
> **C:** Done — I added `src/pair-marker.txt` with exactly one line, `paired locally`,
> in a separate working copy of the widgets project, and confirmed the file reads
> back exactly that. It's uncommitted, so you can look it over. This is work we did
> together and you're the judge of it — I haven't put it through any independent
> check. `[decision_point: make-change-live]` `[decision_point: supervised-not-verified]`
>
> **H:** That looks right. Leave the change uncommitted and available to continue
> later; don't open a pull request.
>
> **C:** Left it as is — uncommitted and set aside so we can pick it up whenever you
> want. Nothing's been committed, no pull request, nothing shipped. `[decision_point: leave-uncommitted-open]` `[decision_point: hide-mechanics]`

The change is made live in an isolated working copy, reported as human-supervised
rather than verified, and preserved uncommitted for later — with no plan, no
verifier, and no provider/mechanism detail exposed.
