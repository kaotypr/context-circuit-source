# Orient a brand-new project — no repository, a simple idea

Proves the safest possible first contact: a lay user arrives with a goal and no
project set up at all. The coordinator orients — offers to get started, asks the
one thing it cannot infer, and fabricates nothing: no repository the user never
named, no plan, no claim that anything was built. This is the entry every other
conversation grows out of.

## Spec
```yaml
id: orient-new-project
title: Brand-new project, no repository, a simple idea
status: proposed
runtime_version: ">=1.0.0"
mode: conversation-only
driver: claude-p
surface: [cc-workspace]
preconditions:
  repositories: []                 # nothing seeded
  state: fresh                     # no repository, no intent, no plan
persona: >
  A solo maker who wants a small command-line tool. Has never heard of Context
  Circuit. Does not know what a workspace, plan, or repository binding is.
  Speaks in plain goals.
human_turns:
  - "Hey, I want to build a little tool that turns a Markdown file into a PDF."
  - "It's just for me, nothing fancy. Can you help me get started?"
  - "I don't have an existing code folder yet — what would you need from me before we decide what to build?"
  - "Thanks, I'll think about where I want the project to live."
reactions:
  approves: false                  # this persona is only scoping; it approves nothing here
  invents_repository: never        # never supplies a repo name or path
  uses_internal_terms: never
demonstrates:
  # Narrowed from case 01's historical [AC-01, AC-04, AC-05] to what THIS dialogue
  # actually demonstrates: registering zero repositories, and the fail-closed
  # invariants. AC-04 (a plan preserves request detail) and AC-05 (envelope
  # re-gate) are not exercised here — nothing is planned or executed — and are
  # owned by the intent/plan plots (03, 04, 16, 20). See coverage/phase-1.md.
  acceptance_criteria: [AC-01]
  invariants: [INV-REPO-04, INV-PLAN-04, INV-EXEC-01, INV-SEC-02]
hidden:
  - internal file names/paths (workspace.yaml, context/INDEX.md, engine.sh)
  - the words "workspace", "repository binding", "base branch", "intent" as jargon
decision_points:
  - id: recognize-no-project
    when: the human states a goal with nothing set up yet
    coordinator_must: >
      recognize that no project exists yet and offer to help get started, in plain
      language — never pretend a project, plan, or code already exists, and never
      claim to have built or run anything.
  - id: no-fabricated-repo
    when: the work would need a home for code
    coordinator_must: >
      ask whether to start a new project or connect an existing one, rather than
      inventing a name or creating a repository the human never named.
  - id: ask-dont-invent
    when: it needs a fact it cannot safely infer
    coordinator_must: >
      ask one plain question instead of guessing — it may describe a possible next
      step, but it does not settle an unstated fact on the human's behalf.
  - id: nothing-persisted
    when: the human leaves the decision open
    coordinator_must: >
      persist nothing consequential — no repository, no plan — and never say the
      work is approved, executed, done, or delivered.
reporting_rules: [plain-language, never-overstate-assurance, faithful-failure]
expected_end_state:
  - repositories_registered: 0       # the human named no repository -> none created
  - plans_created: 0                 # no plan before a project home + approved intent
  - no_plan_status: [done]           # nothing advanced past draft
  - no_execution_records: true       # nothing ran
  - product_knowledge_unchanged_silently: true
access_discipline:
  orient:
    required: [workspace.yaml, context/INDEX.md]
    allowed: [AGENTS.md, WORKFLOW.md, "context/*.md"]
    forbidden: ["plans/archive/**", "wrapper/runtime/engine.sh", "sources/**", "wrapper/contracts/**"]
```

## Dialogue

> **H:** Hey, I want to build a little tool that turns a Markdown file into a PDF.
>
> **C:** Nice — that's a clear little goal. Right now there's nothing set up on my
> side for this yet, so I'm starting from scratch with you, not from an existing
> project. Want me to help you get it off the ground?
> `[decision_point: recognize-no-project]`
>
> **H:** It's just for me, nothing fancy. Can you help me get started?
>
> **C:** Happy to. The first thing to settle is where the code should live. Do you
> already have a folder with code in it you'd like to use, or should we start a
> fresh one? I won't create anything until you tell me which.
> `[decision_point: no-fabricated-repo]` `[decision_point: ask-dont-invent]`
>
> **H:** I don't have an existing code folder yet — what would you need from me
> before we decide what to build?
>
> **C:** Just two things, in plain terms: a home for the code — a fresh project we
> create, or an existing folder you point me at — and then a clear description of
> what "done" looks like for the tool. Once there's a home and we've agreed on
> what it should do, I can turn that into something to build. Nothing's been
> created or started yet; the decision's entirely yours.
> `[decision_point: ask-dont-invent]` `[decision_point: nothing-persisted]`
>
> **H:** Thanks, I'll think about where I want the project to live.
>
> **C:** Sounds good — whenever you're ready, tell me whether to start fresh or
> connect a folder you already have, and we'll go from there. `[decision_point: nothing-persisted]`

The coordinator orients rather than fabricates: it names that nothing exists yet,
offers to get started, asks the single fact it cannot infer (where code lives),
and leaves the workspace untouched — no repository, no plan, no execution, and no
claim of assurance.
