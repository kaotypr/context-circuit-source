# Connect an existing repository — in place, on the human's own branch

Proves binding an existing checkout: the coordinator connects the folder the user
already has (never clones or recreates it), works from the branch the human names
(not the repository's default), and confirms the connection in plain language
without ever exposing the binding machinery. The base branch is user-selected, not
inferred.

## Spec
```yaml
id: connect-existing-repo
title: Connect an existing local repository, on the branch the human works from
status: proposed
runtime_version: ">=1.0.0"
mode: conversation-only
driver: claude-p
surface: [cc-workspace]
preconditions:
  repositories:
    # Seeded INSIDE the workspace before the coordinator runs, so it is
    # pre-existing code the coordinator only has to CONNECT. default_branch is
    # master, but the human works on develop — proving the base is user-selected,
    # never inferred from default_branch (INV-REPO-02).
    - id: my-notes
      dest: my-notes
      default_branch: master
      branches: [master, develop]
      seed_files: [notes.md]
  state: seeded:unconnected-checkout
persona: >
  A solo maker who already has a small notes project in a folder. Has never heard
  of Context Circuit and does not know what a workspace, binding, or base branch
  is. Knows their code lives in a folder called "my-notes" and that they work on
  the "develop" branch. Speaks in plain terms.
human_turns:
  - "I've already got a little notes project — the code's in a folder called my-notes right here."
  - "Can you hook it up so we can work on it together?"
  - "I do all my work on the develop branch, so use that one."
  - "Great, thanks."
reactions:
  approves: false                  # only connecting; no plan, no approval, no execution
  invents_repository: never        # never invents a repo it does not have
  uses_internal_terms: never
demonstrates:
  acceptance_criteria: [AC-01, AC-26, AC-27]
  invariants: [INV-REPO-01, INV-REPO-02, INV-REPO-04, INV-SEC-01]
hidden:
  - internal binding file names/paths (workspace.yaml, repositories.local.yaml)
  - the words "workspace", "binding", "base branch" as jargon
decision_points:
  - id: connect-not-clone
    when: the human points at an existing folder of code
    coordinator_must: >
      offer to connect the folder in place rather than recreating, cloning, or
      copying it, in plain language.
  - id: settle-the-branch
    when: it needs to know which branch to work from
    coordinator_must: >
      ask which branch to work from, or clearly confirm the branch the human named
      — never silently adopt the repository's default branch.
  - id: use-named-branch
    when: the human names their working branch
    coordinator_must: >
      use exactly the branch the human named (develop), not the default it could
      have picked (master).
  - id: confirm-plainly
    when: the connection is made
    coordinator_must: >
      confirm the connection by its effect, in plain language, without exposing
      internal file names, paths, or the phrase "base branch"; and never claim to
      have built, run, planned, or approved anything.
reporting_rules: [plain-language, never-overstate-assurance, faithful-failure]
expected_end_state:
  - repositories_registered: 1             # the named repo is connected...
  - base_branch: my-notes:develop          # ...on the human's branch, not default master
  - no_repository_clone: true              # connected in place; nothing under repositories/
  - plans_created: 0                       # connecting is not planning
  - no_execution_records: true             # nothing ran
  - no_credentials_in_workspace: true      # INV-SEC-01
  - no_machine_path_in_identity: true      # portable identity stays path-free (INV-REPO-01)
access_discipline:
  orient:
    required: [workspace.yaml]
    allowed: [AGENTS.md, WORKFLOW.md, "context/*.md", context/INDEX.md]
    forbidden: ["plans/archive/**", ".context-circuit/wrapper/runtime/engine.sh", "sources/**", ".context-circuit/wrapper/contracts/**"]
  connect-repo:
    required: [workspace.yaml]
    allowed: ["context/**"]
    forbidden: ["plans/archive/**", ".context-circuit/wrapper/runtime/engine.sh", "sources/**"]
```

## Dialogue

> **H:** I've already got a little notes project — the code's in a folder called
> my-notes right here.
>
> **C:** Got it — since the code already exists, I'll connect that folder as-is so
> we can work on it together. I won't copy it or start a new one. `[decision_point: connect-not-clone]`
>
> **H:** Can you hook it up so we can work on it together?
>
> **C:** Will do. One thing I need from you first: which branch do you want me to
> work from? I won't just guess one. `[decision_point: settle-the-branch]`
>
> **H:** I do all my work on the develop branch, so use that one.
>
> **C:** Connected your my-notes project — I'll work from the develop branch you
> named. Nothing's been changed, built, or planned; it's just hooked up and ready
> when you are. `[decision_point: use-named-branch]` `[decision_point: confirm-plainly]`
>
> **H:** Great, thanks.

The coordinator binds the existing checkout in place, settles the working branch
by asking rather than defaulting, uses the branch the human named, and confirms
the connection by its effect — no clone, no exposed binding files, no claim of
work done.
