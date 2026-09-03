# Expected-conversation library (PROPOSAL — not yet adopted)

Status: **proposal**. This folder is a draft for review. It changes no product
behavior, owns no rule, and is not wired into any suite yet. It is source-only
(the release manifest excludes `test/` and the harness; nothing here ships in
`context-circuit-template`).

## Why this exists

Today the human-facing behavior of the template is specified in three
delta-shaped places:

- `sources/context-circuit-v1.0/conversations.md` — five worked dialogues (A–E),
  written to make the **v1.0 changes** concrete.
- `sources/context-circuit-v1.0/human-experience.md` — the **changed** surface
  (the two gates, re-gating, transparency questions).
- `template-harness/scenarios/*/case.yaml` — 22 cases whose expectations are
  **hand-encoded inline** (`visible_expectations`, `transcript_checks`).

Two problems follow from that shape:

1. **A delta is the wrong spec unit.** v1.0 moved the human gate upstream to the
   intent, made completion inferred, and folded pairing into a tier. So the
   end-to-end *flow of pre-existing operations* (approve, complete, deliver, even
   orient) reads differently now — not just the new mechanisms. "What changed"
   cannot express "how the whole conversation goes now." The expected conversation
   must be authored **whole for the current version**, or new machinery gets
   validated against half-remembered old dialogue.

2. **Expectations live inside the checks.** Each `case.yaml` re-states its own
   expected behavior, so there is no single canonical surface to audit coverage
   against, and the cases drift case-by-case as the product evolves.

This library is the fix: one **versioned, whole-surface corpus of expected
human↔coordinator conversations**, from which the harness scenarios are derived
and against which they are revalidated.

## The model: specification by example

Each plot is a concrete example dialogue that *is* the specification — the
established discipline of **specification by example** (a.k.a. executable
specification / living documentation). The example is authored as intent, and the
same example is what the harness executes against a live coordinator. The method,
in one word, is **conversation-spec-first**.

```
expected conversation (the specification, by example)
        │  derives + validates
        ▼
harness scenario  (template-harness/scenarios/*/case.yaml — the executable check)
        │  runs against
        ▼
coordinator + skills + engine  (the implementation)
```

- The **conversation plot** is the specification of expected behavior, expressed
  as a worked example.
- The **harness scenario** is the executable encoding that drives a live
  coordinator and grades the transcript against that example.
- When Context Circuit evolves, you edit the **plot first**. The harness then
  fails against the live coordinator until the skills catch up — so an old flow
  can never silently persist across a version bump. That property is exactly what
  the delta-shaped v0.x→v1.0 material lacked, and it is why the current cases
  ended up scattered and inconsistently mapped. The plots are also **living
  documentation**: kept current, they are a readable, authoritative account of how
  the product actually converses.

## Two disciplines this library must keep

1. **Descriptive, never a second policy owner.** A plot *demonstrates* behavior
   that `wrapper/contracts/invariants.yaml` and `test/acceptance/criteria-map.yaml`
   already own. It references `INV-*`/`AC-*`; it never redefines a rule. Context
   Circuit's "one rule, one owner" discipline stays intact — the plot owns the
   *expected dialogue*, invariants own the *rules*, the criteria map owns
   *AC → suite*.
2. **Versioned, and a real maintenance commitment.** A whole-surface corpus rots
   if it drifts from the skills. That is the cost; the conversation-spec-first
   loop above is what pays it back. Each plot names the `runtime_version` range it
   describes.

## Relationship to the existing material

- `sources/context-circuit-v1.0/conversations.md` stays a **design study** (the
  rationale for the delta). It has no authority and is not the test spec.
- `test/scenarios/test-scenarios.sh` stays the **deterministic** owner of the
  five end-to-end families (AC-36), engine-driven. This library is the **live,
  lay-user, whole-surface** complement — it is not a duplicate of that suite.
- This library **supersedes the inline expectations** in `scenarios/*/case.yaml`:
  a case references a plot id instead of re-encoding what the coordinator should
  say.

## Plot file format

One file per conversation *plot* under `plots/<id>.md`. Each has a
machine-readable `## Spec` (YAML) the grader can consume, and a human-readable
`## Dialogue`. The `## Spec` shape:

```yaml
id: <kebab-id>                 # stable; the case.yaml references this
title: <plain title>
status: proposed | adopted
runtime_version: ">=1.0.0"     # the version range this dialogue describes
mode: conversation-only | full-execution
surface: [<skill or operation>, ...]      # what template area this exercises
preconditions:
  repositories: [...]          # workspace seed (same vocabulary as case.yaml setup)
  state: fresh | seeded:<what>
persona: >
  A lay user with a goal, who has never heard of Context Circuit.
demonstrates:
  acceptance_criteria: [AC-..] # the precise ACs this plot proves
  invariants: [INV-..]
hidden:                        # must NEVER surface to the lay user
  - internal file names, paths, digests, branch names, model names, tier labels
decision_points:               # the ordered beats the coordinator MUST land
  - id: <beat-id>
    when: <what the human just did / asked>
    coordinator_must: <the plain-language behavior that must occur>
reporting_rules: [plain-language, never-overstate-assurance, faithful-failure]
```

The `## Dialogue` is the turn-by-turn expected exchange, each coordinator beat
annotated `[decision_point: <beat-id>]` so a grader (or a reader) can line the
prose up with the checkable spec.

## How a `case.yaml` references a plot (proposed shape)

Instead of re-encoding expectations, a case points at a plot and supplies only
what the driver needs (seed + budgets); the grader validates the live transcript
against the plot's `decision_points`, `hidden`, and `reporting_rules`:

```yaml
id: 05-approve-and-execute
conversation: execute-standard-verify-not-complete   # -> plots/<id>.md
mode: full-execution
setup: { ... }          # workspace seed (unchanged)
grader:
  from_plot: true       # acceptance_criteria / invariants / transcript rules come from the plot
  post_conditions: [ ... ]   # state assertions stay in the case (engine-level, per-seed)
  budgets: { ... }
```

State post-conditions (engine-level, seed-specific) remain in the case; the
*conversational* expectations move to the plot. One canonical surface to audit,
one place to edit when the flow changes.

## The whole-surface plot catalog (proposed)

Grouped by phase. "case" = the current `scenarios/` case that maps to the plot;
**NEW** = a real coverage gap with no live case today; *(opt)* = optional
enhancement, not a hole (the mechanic is covered deterministically elsewhere).

### Orientation & context foundation
| plot | proves | case |
| --- | --- | --- |
| orient-new-project | no project/plan exists yet; orient, don't fabricate | 01 |
| connect-existing-repo | bind in place on the user's branch | 02 |
| clone-or-init-new-repo | clone/init a new repo under `repositories/` (AC-27, INV-REPO-03) | **NEW** |
| build-product-knowledge | first-time knowledge build (`init-project-context`) | **NEW** |
| query-product-knowledge | "what does the project know about X" via retrieval | **NEW** |
| accept-or-defer-context-proposal | knowledge acceptance is an explicit human decision (INV-KNOWLEDGE-02) | **NEW** |

### Intent gate (Gate 1)
| plot | proves | case |
| --- | --- | --- |
| review-intent-approve-derive | review intent, settle open question, approve, auto-derive plan | 03 |
| adversary-revise-reapprove | adversary finds a gap → revise → re-challenge → approve | 15 |
| refuse-before-gate1 | fail closed: no build before intent approval | 04 |
| tier-fails-upward-refuse-explore | decline Explore on a security surface (crown jewel 2) | 20 |
| scope-envelope-regate | plan exceeds approved scope → held, re-gated (crown jewel 1) | 16 |
| approve-and-build-one-turn | approve intent AND build in one turn (INV-EXEC-01 compound) | **NEW** |

### Explore tier & promotion
| plot | proves | case |
| --- | --- | --- |
| direct-collaboration-explore | human-supervised, uncommitted, no verifier | 14 (host-neutralize) |
| explore-promote-to-standard | promote in place → intent+adversary+verifier+candidate | **NEW** |

### Execute & verify
| plot | proves | case |
| --- | --- | --- |
| execute-standard-verify-not-complete | verified ≠ complete | 05 |
| repository-grounding | worker honors the repo's own agent guidance | 11 |
| verifier-host-blocked | host can't spawn verifier → blocked, never self-verify | 09 |
| run-stack-single-repo | dependency-ordered single-repo batch | 10 |
| run-stack-multi-repo | cross-repo ordering gate, same-repo stacking | 12 |
| execution-tiering-hidden | per-role model/effort recorded, hidden from user | 13 |

### Candidate, completion, delivery
| plot | proves | case |
| --- | --- | --- |
| stale-candidate-refuse-complete | new commit voids evidence; refuse completion | 17 |
| critical-repair-then-complete | Critical explicit completion after a repair | 06 |
| delivery-boundary-block-no-remote | Gate 2 separate; no remote → block, no silent push | 08 |
| standard-inferred-completion | inferred completion after acceptance + Gate 2 | 18 |
| change-set-one-verification | one integrated candidate, one PR, complete both | 21 |
| change-set-base-unbuildable | honest "these don't combine cleanly, I stopped" (conversation C tail) | **NEW** |

### Closed knowledge loop
| plot | proves | case |
| --- | --- | --- |
| knowledge-debt-blocks-next-plan | reconciliation debt gates the next plan | 19 |

### Organization
| plot | proves | case |
| --- | --- | --- |
| archive-plan | status-blind archive | 07 |
| restore-archived-plan | status-blind restore | 07b |
| archive-restore-intent | intent-side archive/restore (part of AC-30) | *(opt)* |

### Peripheral surfaces (template-invocable, currently uncovered)
| plot | proves | case |
| --- | --- | --- |
| publish-plan-to-external | cc-publish plan kind: self-contained, one-way, idempotent | **NEW** |
| publish-open-questions-thread | cc-publish thread kind | **NEW** |
| author-system-design-spawns-intents | cc-system-design authors `sources/system-design/` that spawns intents | **NEW** |

### Whole-flow (the lay-user experience across phases)
| plot | proves | case |
| --- | --- | --- |
| standard-feature-whole-flow | Scenario A end-to-end as one live conversation (two gates + inferred completion + reconcile) | *(opt)* |

Two plots are fully worked as examples in `plots/`:
`standard-feature-whole-flow.md` (the whole-flow the delta files never play live)
and `build-product-knowledge.md` (the "gathering context" family that is missing
everywhere). The rest are catalog entries to be authored if this proposal is
adopted.

## If adopted — suggested build order

1. Author the plots that map to existing cases (transcribe the inline
   expectations into the plot format; no behavior change) — establishes the
   canonical surface.
2. Switch `case.yaml` to `conversation:` references + `from_plot: true`; delete
   the inline `visible_expectations`/`transcript_checks` duplication.
3. Author the **NEW** plots and add their cases (real coverage gaps).
4. Add the *(opt)* whole-flow and intent-archive plots if the maintainer wants
   the library to carry the lay-user's whole-conversation experience, not only
   discrete probes.
