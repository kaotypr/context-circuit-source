# Start a fresh repository — initialize a new project home

Proves the other half of "where does the code live": the user has no code yet and
wants to start fresh. The coordinator initializes a new repository under the
workspace's own project area, sets the working branch the user names (never the
suggested default), makes the first commit that a fresh repository needs, and
stops there — creating a remote or pushing is a separate action it does not take.

## Spec
```yaml
id: clone-or-init-new-repo
title: Initialize a new project repository from scratch
status: proposed
runtime_version: ">=1.0.0"
mode: conversation-only
driver: claude-p
surface: [cc-workspace]
preconditions:
  repositories: []                 # nothing seeded; the repository is created live
  state: fresh
persona: >
  A solo maker with a project idea and no code yet, who wants to start a brand-new
  repository for it. Has never heard of Context Circuit; does not know what a base
  branch, a default branch, or a remote is. Names the project themselves.
human_turns:
  - "I don't have any code yet — can you start a fresh repo for this project? Call it inkwell."
  - "Yeah, create it new, don't clone anything from anywhere."
  - "I'll work on main, use that."
  - "Great, that's all for now."
reactions:
  approves: false                  # setup only; no plan, no approval, no execution
  invents_repository: never        # the human names "inkwell"; the coordinator invents nothing
  uses_internal_terms: never
demonstrates:
  acceptance_criteria: [AC-01, AC-27]
  invariants: [INV-REPO-01, INV-REPO-02, INV-REPO-03, INV-REPO-04, INV-DELIVER-01]
hidden:
  - internal file names/paths (workspace.yaml, repositories.local.yaml, repositories/, engine.sh)
  - the words "base branch", "default branch", "remote", "binding" as jargon
decision_points:
  - id: init-not-connect
    when: the human asks to start a fresh repository they named
    coordinator_must: >
      offer to create a new repository (not connect an existing one and not clone),
      and say plainly where it will live and what branch it will use before creating it.
  - id: new-repo-default-home
    when: it creates the new repository
    coordinator_must: >
      place the new repository in the workspace's own project area by default,
      keeping the workspace's own files separate from it — the human need not know
      the path, but a new repository never lands loose in the workspace root.
  - id: user-branch-not-default-guidance
    when: the human names the branch to work on
    coordinator_must: >
      use exactly the branch the human named (main) as the branch work happens on;
      any portable "suggested first branch" is guidance only and never silently
      becomes the working branch.
  - id: local-only-no-remote
    when: the repository is initialized
    coordinator_must: >
      create it locally with the first commit it needs, and stop there — it does
      not create a remote, push, or publish, and says so if asked; those are
      separate actions taken only on request.
reporting_rules: [plain-language, never-overstate-assurance, faithful-failure]
expected_end_state:
  # A live-created repo needs id-agnostic / path-shape post-conditions, not the
  # id-keyed ones seeded cases use. `repository_under_project_area` and
  # `no_remote_created` are the positive inverses of case 02's `no_repository_clone`
  # and are HARNESS GAPS today (see coverage/phase-1.md, gaps).
  - repositories_registered: 1
  - repository_under_project_area: true    # created in repositories/, not the workspace root (INV-REPO-03)
  - base_branch: inkwell:main              # user-selected working branch (INV-REPO-02)
  - initial_commit_present: true           # a fresh repo has the base commit execution later needs
  - no_remote_created: true                # creating a remote/pushing is a separate delivery action (INV-DELIVER-01)
  - plans_created: 0
  - no_execution_records: true
  - no_credentials_in_workspace: true
  - no_machine_path_in_identity: true
access_discipline:
  orient:
    required: [workspace.yaml, context/INDEX.md]
    allowed: [AGENTS.md, WORKFLOW.md, "context/*.md"]
    forbidden: ["plans/archive/**", ".context-circuit/wrapper/runtime/engine.sh", "sources/**", ".context-circuit/wrapper/contracts/**"]
  create-repo:
    required: [workspace.yaml]
    allowed: ["context/**"]
    forbidden: ["plans/archive/**", ".context-circuit/wrapper/runtime/engine.sh", "sources/**"]
```

## Dialogue

> **H:** I don't have any code yet — can you start a fresh repo for this project?
> Call it inkwell.
>
> **C:** Sure. To be clear on what I'll do: create a brand-new, empty project
> called inkwell — not connect an existing folder and not copy one from anywhere —
> living in this workspace's project area, kept separate from my own working files.
> Want me to go ahead? `[decision_point: init-not-connect]` `[decision_point: new-repo-default-home]`
>
> **H:** Yeah, create it new, don't clone anything from anywhere.
>
> **C:** One thing before I make it: which branch do you want to work on? I'll use
> whatever you name rather than picking one for you. `[decision_point: user-branch-not-default-guidance]`
>
> **H:** I'll work on main, use that.
>
> **C:** Done — inkwell is created and ready, and I'll work on the main branch you
> chose. It's local only: I haven't set up anywhere remote to push it or published
> it anywhere — those are separate steps I'll only take if you ask. Nothing's been
> planned or built yet. `[decision_point: local-only-no-remote]`
>
> **H:** Great, that's all for now.

The coordinator creates the named repository fresh (not a connect, not a clone),
puts it in the project area rather than loose in the workspace, uses the branch the
human chose over any suggested default, gives it the first commit a new repository
needs, and stops at local — no remote, no push, no plan, no execution.
