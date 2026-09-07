# Structure a larger design as source — then spawn one intent per concern

Proves the design→intents path for a change too big for a single plan: the coordinator
authors a structured system design as passive source material (it drafts files only —
no approval, no plan, no project-knowledge write), and then turns each distinct concern
into its own approvable piece of work — one intent per concern, not one giant intent.

## Spec
```yaml
id: author-system-design-spawns-intents
title: Author a structured system design as source, then spawn one intent per concern
status: proposed
runtime_version: ">=1.0.0"
mode: conversation-only
driver: claude-p
surface: [cc-system-design, cc-intent]
preconditions:
  repositories:
    - id: app
      dest: app
      default_branch: main
      branches: [main]
      seed_files: [README.md]
      connect: main
  state: seeded:connected-repo
persona: >
  A maker planning a sizeable new capability who wants the shape written down before any
  building, then each part turned into its own thing to approve. Has never heard of
  Context Circuit; does not know what a system design source or an intent is.
human_turns:
  - "I want to design a whole notifications system for the app — can you structure the design before we build anything?"
  - "Good structure. Now turn each of those concerns into its own piece of work I can approve separately."
  - "Thanks."
reactions:
  approves: false                  # authoring + spawning drafts; no approval given here
  invents_repository: never
  uses_internal_terms: never
demonstrates:
  acceptance_criteria: [AC-30]     # intents are first-class; genuinely separate concerns become separate intents
  invariants: [INV-INTENT-01, INV-SEC-02, INV-SKILL-01, INV-PLAN-04]
hidden:
  - internal file names/paths (sources/system-design/…, design.md, intent id, contract.yaml)
  - the words "system design source", "three-tier layout", "feasibility check" as jargon
decision_points:
  - id: structure-not-build
    when: the human asks to structure the design first
    coordinator_must: >
      write the design down as organized source material — the shape, the concerns,
      the fixed decisions — and be clear it's a draft to shape: it approves nothing,
      builds nothing, and records no project knowledge yet.
  - id: gaps-become-questions
    when: something in the design is undecided
    coordinator_must: >
      surface an undecided point as an open question rather than silently choosing it.
  - id: one-intent-per-concern
    when: the human asks to turn the concerns into approvable work
    coordinator_must: >
      create a separate approvable piece of work for each genuinely distinct concern —
      one per concern, each reviewable and shippable on its own — never one giant catch-all.
  - id: draft-not-approved
    when: it spawns the pieces
    coordinator_must: >
      leave each as a draft awaiting the human's separate approval; authoring the design
      and the pieces approves, plans, or builds nothing.
reporting_rules: [plain-language, never-overstate-assurance, faithful-failure]
expected_end_state:
  # Authoring the design is a normal write under sources/system-design/; spawning the
  # concerns creates several DRAFT intents. `system_design_authored` / `intents_created`
  # (count) are NEW predicates (harness gaps) — see coverage/phase-6b.md.
  - repositories_registered: 1
  - system_design_authored: true             # structured source written under sources/system-design/
  - intents_created: ">=2"                    # one draft intent per distinct concern
  - intents_approved: 0                       # nothing approved here
  - plans_created: 0                          # no plan before an approved intent
  - no_execution_records: true
  - product_knowledge_unchanged_silently: true # a design is source, not accepted knowledge
access_discipline:
  design:
    # Authoring under sources/system-design/ IS the write here (a normal authoring write);
    # other sources/ subtrees and the runtime stay off-limits.
    required: []
    allowed: ["sources/system-design/**", "context/**", AGENTS.md, WORKFLOW.md]
    forbidden: ["plans/archive/**", ".context-circuit/wrapper/runtime/engine.sh", ".context-circuit/wrapper/contracts/**"]
  spawn-intents:
    required: []
    allowed: ["intent/**", "sources/system-design/**", "context/**"]
    forbidden: ["plans/archive/**", ".context-circuit/wrapper/runtime/engine.sh"]
```

## Dialogue

> **H:** I want to design a whole notifications system for the app — can you structure
> the design before we build anything?
>
> **C:** Good place to start. I've written the shape down as an organized design you can
> review end to end: the overall goal, and separate sections for the concerns — how
> notifications are delivered, how people set their preferences, and how delivery is
> retried when it fails. One thing I left open rather than deciding: whether preferences
> are per-device or per-account — that's a real choice for you. This is a draft to shape;
> I haven't approved, built, or recorded anything as settled. `[decision_point: structure-not-build]` `[decision_point: gaps-become-questions]`
>
> **H:** Good structure. Now turn each of those concerns into its own piece of work I
> can approve separately.
>
> **C:** Done — each concern is now its own piece you can review and approve on its own:
> one for delivery, one for preferences, one for retry. They're kept separate on purpose,
> so you can approve and ship them independently rather than as one big lump. Each is a
> draft waiting for your go; nothing's approved or building. `[decision_point: one-intent-per-concern]` `[decision_point: draft-not-approved]`
>
> **H:** Thanks.

The larger change is written down as structured source (a draft, approving and building
nothing), and each genuinely separate concern becomes its own approvable piece of work —
one per concern, not one giant intent.
