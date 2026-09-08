# Promote an Explore session in place — the ramp into the trust pipeline

Proves the promote step: an Explore session that turns out to be real work becomes a
properly-checked change *without stopping and starting over*. The coordinator
attaches an intent (so the change goes through the intent gate, and on approval the
real code is read before any plan), raises the tier so an independent check appears,
and authors a plan of record whose execution binds to the commits already made — the
cliff is a ramp. Nothing about the earlier work is thrown away or redone.

## Spec
```yaml
id: explore-promote-to-standard
title: Promote an Explore session in place to a checked Standard change
status: proposed
runtime_version: ">=1.0.0"
mode: full-execution
driver: cc-test-case
surface: [cc-pair, cc-intent, cc-plan]
preconditions:
  repositories:
    - id: accounts
      dest: accounts
      default_branch: main
      branches: [main]
      seed_files: [README.md, src/util.js]
      connect: main
  # A single continuous arc: a short pairing beat establishes real work, then the
  # human asks to make it proper. Authored as one in-session flow rather than a
  # seeded mid-pairing session, because there is no seed shape for pairing-branch
  # commits today (see coverage/phase-3.md, gaps). Kept to a non-risk surface so the
  # raised tier is Standard (one independent check appears), not Critical.
  state: fresh
persona: >
  A maker pairing on a small helper that grows into real logic they care about. They
  want it made into a proper, checked change without starting over. Has never heard
  of Context Circuit; does not know what an intent, a tier, or a plan of record is.
human_turns:
  - "Let's work directly on the accounts project — add a small helper in src/util.js that formats a display name."
  - "Actually this is turning into real logic I care about. Can we make this a proper, checked change without starting over?"
  - "Yes, that goal's right — approve it with the independent check on. Don't build more or ship it yet."
  - "Thanks."
reactions:
  approves: true                   # approves the attached intent (Gate 1) at the explicit turn
  invents_repository: never
  uses_internal_terms: never
demonstrates:
  acceptance_criteria: [AC-29, AC-30, AC-33]
  invariants: [INV-PAIR-01, INV-ASSURE-01, INV-INTENT-01, INV-APPROVE-01]
hidden:
  - the promote mechanism and its runtime state (pairing branch, plan-of-record binding, candidate)
  - tier labels ("Explore", "Standard"), model/effort/role details, the intent id, "verifier" as jargon
decision_points:
  - id: pair-then-notice-real
    when: the human pairs on a small change that grows
    coordinator_must: >
      work the change live (human-supervised, no independent check yet) and, when
      it clearly matters, offer to make it a proper checked change rather than
      restarting from scratch.
  - id: attach-intent-and-approve
    when: the human agrees to make it proper
    coordinator_must: >
      attach a goal with what "correct" means for what the session actually did,
      surface the criteria and scope, and present it for a single approval — never
      fabricating them silently.
  - id: raise-tier-check-appears
    when: it presents the change for approval
    coordinator_must: >
      say plainly that this now gets an independent check (the earlier live work did
      not), because the change is real — without exposing tier labels as jargon.
  - id: promote-in-place-not-restart
    when: the human approves
    coordinator_must: >
      build on the work already done — the approved change is recorded as a proper
      plan bound to what was already made, not a fresh redo — and honor "don't build
      more or ship it yet".
reporting_rules: [plain-language, never-overstate-assurance, faithful-failure]
expected_end_state:
  # Fresh whole-arc case: the intent/plan ids are authored live, so id-keyed
  # assertions assume the first-allocated ids (i001-.../0001-...). Same id-agnostic
  # concern the whole-flow plot raises; `promoted_in_place` is a NEW assertion that
  # the plan of record binds to the pairing commits rather than a fresh redo — a
  # harness gap (see coverage/phase-3.md).
  - repositories_registered: 1
  - intent_status: i001-util-helper:approved
  - classified_tier: i001-util-helper:standard      # raised from Explore; the independent check appears
  - plans_created: 1
  - plan_schema: 0001-util-helper:3
  - promoted_in_place: true                           # plan of record binds to the existing pairing commits
  - no_delivery_records: true                         # "don't ship it yet" honored
  - product_knowledge_unchanged_silently: true
access_discipline:
  direct-collaboration:
    required: []
    allowed: [".agents/skills/cc-pair/SKILL.md", ".context-circuit/agents/worker.md", "role-tiering.local.yaml", ".runtime/pairing/**", ".runtime/explore/**"]
    forbidden: ["plans/archive/**", ".context-circuit/wrapper/runtime/engine.sh", "sources/**"]
  promote:
    required: []
    allowed: ["intent/**", "plans/**", "context/**", ".runtime/pairing/**"]
    forbidden: ["plans/archive/**", ".context-circuit/wrapper/runtime/engine.sh", "sources/**", ".context-circuit/wrapper/contracts/**"]
```

## Dialogue

> **H:** Let's work directly on the accounts project — add a small helper in
> src/util.js that formats a display name.
>
> **C:** Added a small display-name helper to `src/util.js` and tried it out — it's
> working. This is live work between us, not independently checked yet, and it's
> uncommitted so you can see it. `[decision_point: pair-then-notice-real]`
>
> **H:** Actually this is turning into real logic I care about. Can we make this a
> proper, checked change without starting over?
>
> **C:** Yes — and we don't have to redo anything. I'll take what we've already
> built, write down what it's meant to do and where it's allowed to touch, and put it
> to you for approval. Because it's now real work, it'll also get an independent check
> that the quick live version didn't. Here's what I'd write as the goal: a
> display-name formatter in the accounts project's utilities, handling empty and
> unusually long names cleanly. Sound right?
> `[decision_point: attach-intent-and-approve]` `[decision_point: raise-tier-check-appears]`
>
> **H:** Yes, that goal's right — approve it with the independent check on. Don't
> build more or ship it yet.
>
> **C:** Approved. I looked over the code we already wrote, and the work we already
> did is now recorded as the plan for this change — built on, not redone — and it's
> set up to get its independent check. I've not gone further or shipped anything,
> exactly as you said. `[decision_point: promote-in-place-not-restart]`
>
> **H:** Thanks.

The live Explore work is promoted in place: an intent is attached, the change is put
up for one approval, the independent check appears because the work is now real, and
the approved plan binds to the commits already made — no restart, and nothing shipped.
