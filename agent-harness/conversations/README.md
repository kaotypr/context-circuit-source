# Expected-conversation library (PROPOSAL — not yet adopted)

Status: **proposal**. This folder is a draft for review. It changes no product
behavior, owns no rule, and is not wired into any suite yet. It is source-only
(the release manifest excludes `test/` and the harness; nothing here ships in
`context-circuit-template`).

> **Two ways in.** This README is the **coverage** view (plots grouped by phase,
> with gap analysis). To read the library **in lifecycle order** — the golden path
> top-to-bottom, with each branch and guardrail shown where it hangs off — see
> **[FLOW.md](FLOW.md)**.

## Why this exists

Today the human-facing behavior of the template is specified in three
delta-shaped places:

- `sources/context-circuit-v1.0/conversations.md` — five worked dialogues (A–E),
  written to make the **v1.0 changes** concrete.
- `sources/context-circuit-v1.0/human-experience.md` — the **changed** surface
  (the two gates, re-gating, transparency questions).
- `agent-harness/scenarios/*/case.yaml` — 22 cases whose expectations are
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
harness scenario  (agent-harness/scenarios/*/case.yaml — the executable check)
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
  a case is **generated from** a plot instead of re-encoding what the coordinator
  should say. The plot is the source; the case is a build artifact.

## Plot file format

One file per conversation *plot* under `plots/<id>.md`. Each has a
machine-readable `## Spec` (YAML) the grader can consume, and a human-readable
`## Dialogue`. The `## Spec` shape:

```yaml
id: <kebab-id>                 # stable; the generated case derives from this
title: <plain title>
status: proposed | adopted
runtime_version: ">=1.0.0"     # the version range this dialogue describes
mode: conversation-only | full-execution
driver: claude-p | cc-test-case           # which harness driver the generated case wires
surface: [<skill or operation>, ...]      # what template area this exercises
preconditions:
  repositories: [...]          # workspace seed -> case.yaml setup.repositories
  state: fresh | seeded:<what>
persona: >
  A lay user with a goal, who has never heard of Context Circuit.
human_turns:                   # the turns the driver plays -> case.yaml human.turns
  - "<what the human says, turn 1>"
  - "<turn 2>"
reactions:                     # -> case.yaml human.reactions
  approves: true | false
  invents_repository: never
  uses_internal_terms: never
demonstrates:                  # -> grader.acceptance_criteria / invariants
  acceptance_criteria: [AC-..]
  invariants: [INV-..]
hidden:                        # -> grader.transcript_checks.forbids_regex (never surface these)
  - internal file names, paths, digests, branch names, model names, tier labels
decision_points:               # ordered beats the coordinator MUST land
  - id: <beat-id>              #   -> human.visible_expectations + transcript_checks.requires_any
    when: <what the human just did / asked>
    coordinator_must: <the plain-language behavior that must occur>
reporting_rules: [plain-language, never-overstate-assurance, faithful-failure]
expected_end_state:            # -> grader.post_conditions (only what the harness can assert)
  - <post-condition>: <value>
access_discipline:             # -> grader.access_policy (right files, only necessary files)
  <phase>: { required: [...], allowed: [...], forbidden: [...] }
```

The `## Dialogue` is the turn-by-turn expected exchange, each coordinator beat
annotated `[decision_point: <beat-id>]` so a reader — and a lint that keeps the
`**H:**` lines in sync with `human_turns` — can line the prose up with the spec.

## How a plot generates a case (proposed shape)

The plot is the source; the `case.yaml` is a **generated artifact**, not a
hand-authored file that references the plot. A small deterministic generator maps
plot fields to case fields:

| plot field | generated case field |
| --- | --- |
| `preconditions.repositories` | `setup.repositories` |
| `human_turns` | `human.turns` (each `say:`) |
| `persona`, `reactions` | `human.persona`, `human.reactions` |
| `decision_points` | `human.visible_expectations` + `transcript_checks.requires_any` |
| `hidden` | `transcript_checks.forbids_regex` |
| `demonstrates` | `grader.acceptance_criteria`, `grader.invariants` |
| `expected_end_state` | `grader.post_conditions` |
| `access_discipline` | `grader.access_policy` |

The generated case carries a `# generated from plots/<id>.md — do not edit`
header. The **only** thing not sourced from the plot is environmental tuning that
is not part of the conversation spec: per-host `budgets` and the host lane, which
live in a small overlay beside the generated file. Everything semantic has one
source of truth, so a case cannot drift from its spec — change the plot,
regenerate, and the harness reflects the new expected conversation. A worked
example lives in `generated/` (one plot, its generated case).

## The whole-surface plot catalog (proposed)

Grouped by phase. "case" = the current `scenarios/` case that maps to the plot;
**NEW** = a real coverage gap with no live case today; *(opt)* = optional
enhancement, not a hole (the mechanic is covered deterministically elsewhere).

### Orientation & context foundation
| plot | proves | case |
| --- | --- | --- |
| onboarding-what-is-this | first contact on a blank workspace: explain what it is, what it can do, and how to start, in plain terms | **NEW** |
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
| feasibility-question-then-approve | approve → reading the real code surfaces a decision the goal didn't settle → resolve → derive plan | 15 |
| refuse-before-gate1 | fail closed: no build before intent approval | 04 |
| tier-fails-upward-refuse-explore | decline Explore on a security surface (the one safety-critical automated check) | 20 |
| scope-reach-feasibility-question | delivering would reach outside the approved scope → surfaced as a feasibility question, human decides | 16 |
| approve-and-build-one-turn | approve intent AND build in one turn (INV-EXEC-01 compound) | **NEW** |

### Explore tier & promotion
| plot | proves | case |
| --- | --- | --- |
| direct-collaboration-explore | human-supervised, uncommitted, no verifier | 14 (host-neutralize) |
| explore-promote-to-standard | promote in place → intent+trace+verifier+candidate | **NEW** |

### Execute & verify
| plot | proves | case |
| --- | --- | --- |
| execute-standard-verify-not-complete | verified ≠ complete | 05 |
| repository-grounding | worker honors the repo's own agent guidance | 11 |
| verifier-host-blocked | host can't spawn verifier → blocked, never self-verify | 09 |
| run-stack-single-repo | dependency-ordered single-repo batch | 10 |
| run-stack-multi-repo | cross-repo ordering gate, same-repo stacking | 12 |
| execution-tiering-hidden | per-role model/effort recorded, hidden from user | 13 |
| three-failure-stop | third failed attempt → stop and report honestly, everything preserved (AC-13, INV-REPAIR-01) | **NEW (conv. probe)** |
| interrupted-recovery | an interrupted run resumes or stops read-only, nothing lost (AC-19, INV-PRESERVE-01) | **NEW (conv. probe)** |
| lease-ownership-conflict | work another execution holds → read-only/blocked, never a silent steal (INV-CONCURRENCY-01, INV-OWN-01) | **NEW (conv. probe)** |

The three **conv. probe** entries above exercise how a deterministically-owned
mechanic (owned by `test/runtime/test-runtime.sh` and the engine tests) is
*explained to a lay user in plain language* — the live conversational surface that
has no case today. They add no mechanic coverage; they cross-reference the
deterministic owner.

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
| reconcile-and-proceed | the happy path: human reconciles the debt, it clears, the next plan proceeds (AC-34, INV-KNOWLEDGE-02) | **NEW** (positive complement of 19) |

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

1. Author the plot library **fresh from the whole surface** — one lifecycle phase
   at a time — using the existing 22 cases and the v1.0 delta docs as *inputs and
   a coverage cross-check* (every existing case must map to a plot), not as a
   1:1 transcription source. The spec is authored as intent, not reverse-engineered
   from whatever cases happen to exist.
2. Stand up the generator and regenerate `scenarios/*/case.yaml` from the plots;
   retire the inline `visible_expectations`/`transcript_checks` duplication. The
   calibrated executable layer (seed states, budgets) is preserved as the per-case
   environmental overlay — evolve the cases, do not recreate them.
3. Author the **NEW** plots and generate their cases (the real coverage gaps).
4. Add the *(opt)* whole-flow and intent-archive plots if the maintainer wants
   the library to carry the lay-user's whole-conversation experience, not only
   discrete probes.
