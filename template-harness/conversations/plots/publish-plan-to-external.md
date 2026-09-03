# Publish a plan to an external tracker — one-way, self-contained, preview-first

Proves the external surface's `plan` kind: reflecting a plan and its tasks outward to a
team tracker. The coordinator treats it as a manual, export-only action orthogonal to
the workflow — it previews what would appear in plain language a teammate could read
(no internal file names, paths, or ids), pushes nothing until an explicit go, and never
changes the plan or its status.

## Spec
```yaml
id: publish-plan-to-external
title: Publish a plan and its tasks to an external tracker, export-only and preview-first
status: proposed
runtime_version: ">=1.0.0"
mode: conversation-only
driver: claude-p
surface: [cc-publish]
preconditions:
  repositories:
    - id: reports
      dest: reports
      default_branch: main
      branches: [main]
      seed_files: [README.md]
      connect: main
  plans:
    - id: 0021-csv-export
      title: Add CSV export to the reports page
      repository: reports
      objective: Add a CSV export button to the reports page with a couple of tasks.
      path: src/reports
      seed_state: draft
  # No publication configured yet; the coordinator would set one up on first use.
  state: seeded:one-plan-no-publication
persona: >
  A maker who wants teammates to follow a plan in the team's tracker. Has never heard
  of Context Circuit; does not know what a publication, a field intent, or a record is.
  Speaks in plain terms.
human_turns:
  - "Can you publish the CSV export plan to our team's tracker so everyone can follow it?"
  - "Show me what it'd look like first, before you send anything."
  - "Looks good — but hold off actually sending it for now."
  - "Thanks."
reactions:
  approves: false                  # holds off the actual push; a preview only
  invents_repository: never
  uses_internal_terms: never
demonstrates:
  # No dedicated acceptance criterion: the external surface is owned by
  # INV-EXTERNAL-01/02/03 (a v0.7 addition); the v0.5-rooted criteria map has no
  # publish AC. See coverage/phase-6b.md.
  acceptance_criteria: []
  invariants: [INV-EXTERNAL-01, INV-EXTERNAL-02, INV-EXTERNAL-03, INV-PLAN-01]
hidden:
  - internal file names/paths (plan.yaml, publication/, config.yaml), internal ids (AC-…/VR-…)
  - the words "publication", "field intent", "mirror", "source of truth" as jargon
decision_points:
  - id: export-only-orthogonal
    when: the human asks to publish a plan to a tracker
    coordinator_must: >
      recognize this as sending data outward — one-way, and separate from the actual
      work — and say plainly it won't change the plan or how the work proceeds.
  - id: self-contained-preview
    when: it shows what would appear
    coordinator_must: >
      preview what a teammate would see in the tracker in plain language — the plan as
      one item, its tasks as items under it, acceptance as plain checklist lines — with
      NO internal file names, paths, or ids.
  - id: consult-before-push
    when: it has a preview ready
    coordinator_must: >
      push nothing until an explicit go; a preview sends nothing to the tracker. (If
      there's no way to reach the tracker, it says so and changes nothing.)
  - id: never-touch-plan-state
    when: throughout
    coordinator_must: >
      never change the plan, its status, or any workspace state as a result of
      publishing or previewing.
reporting_rules: [plain-language, never-overstate-assurance, faithful-failure]
expected_end_state:
  # Mostly transcript-borne (self-contained text, orthogonality). The state assertions
  # confirm publishing changed nothing in the workspace. `no_provider_push` /
  # `plans_untouched` are NEW predicates (harness gaps) — see coverage/phase-6b.md.
  - repositories_registered: 1
  - plan_status: 0021-csv-export:draft       # unchanged by publish/preview
  - plans_untouched: true                     # nothing written under plans/
  - no_provider_push: true                    # preview only; nothing sent outward
  - product_knowledge_unchanged_silently: true
access_discipline:
  publish:
    # A publication reads plans/** and writes only under publication/**; never plans/.
    required: []
    allowed: ["plans/**", "publication/**", "intent/**"]
    forbidden: ["plans/archive/**", "wrapper/runtime/engine.sh", "sources/**", "wrapper/contracts/**"]
```

## Dialogue

> **H:** Can you publish the CSV export plan to our team's tracker so everyone can
> follow it?
>
> **C:** Sure — to be clear on what that does: it sends a copy outward to your tracker
> so people can follow along. It's one-way and separate from the actual work; it won't
> change the plan or how anything proceeds here. Want me to show you what it'd look
> like before anything goes out? `[decision_point: export-only-orthogonal]`
>
> **H:** Show me what it'd look like first, before you send anything.
>
> **C:** Here's what a teammate would see:
>
> - **[0021] Add CSV export to the reports page** — a card describing the goal in plain
>   terms, with one sub-item per task and each task's acceptance as a plain checklist.
>
> Nothing's been sent — this is just a preview. `[decision_point: self-contained-preview]` `[decision_point: consult-before-push]`
>
> **H:** Looks good — but hold off actually sending it for now.
>
> **C:** Held off — nothing went to the tracker, and the plan here is exactly as it
> was. Say the word whenever you want it sent. `[decision_point: never-touch-plan-state]`
>
> **H:** Thanks.

The coordinator treats publishing as a manual, one-way, orthogonal action: it previews
self-contained content, sends nothing without an explicit go, and leaves the plan and
its status untouched.
