# 0002 — Change intent identifiers to three digits

- **Plan ID:** `0002-intent-id-width-migration`
- **Intent:** `i003-intent-id-width`
- **Status:** draft
- **Repository:** `context-circuit-source`
- **Assurance:** Standard — one worker and one independent verifier

## Original request

Change intent identifiers from four digits after `i` to three digits, for example
`i001`, and update the existing intent records and every current reference to
them. Old or unrelated intents may be archived or deleted only through an
explicit human request or manual human action.

## Objective and desired behavior

- The canonical form is `i<NNN>-slug`, with exactly three zero-padded digits.
- The three existing intent records move together to their canonical three-digit
  names, including directories, canonical fields, indexes, traces,
  plan parent links, and all current references.
- Runtime validation, allocation, approval, authorization, archive, restore,
  tracing, planning, and related lifecycle consumers agree on the new form.
- Approved intent digests are recomputed with the existing digest algorithm after
  the identity-only rename; the approved criteria and lifecycle semantics stay
  unchanged.
- Allocation never silently archives, deletes, or recycles an intent. If the
  exact three-digit range is exhausted, it fails clearly and requires an explicit
  human decision.

## Constraints and non-goals

- Keep plan identifiers in their existing `NNNN-<slug>` form.
- Preserve existing intent statuses, criteria, relationships, trace evidence,
  plan meaning, approval state, and history except for the identifier spelling and
  the required digest recomputation.
- Include the shipped `agent-harness/` and `template/` references; the migration
  is workspace-wide for current references.
- Do not change approval gates, plan authorization, execution, verification,
  completion, delivery, or tier semantics.
- Do not rewrite historical Git objects or commit messages.
- Do not automatically archive, delete, or reuse identifiers. Manual archive or
  deletion remains a direct human action.
- Preserve unrelated and dirty work. Do not scan or modify `sources/`; use the
  approved trace and current Product Knowledge as grounding.

## Grounded findings

The refreshed trace found one repository in scope, three active intent directories,
zero archived intent directories, and 148 old-form token matches across 130 lines
and 46 files. The change is feasible at Standard. The existing conversation-library
intent has an
active plan and trace manifest; the plan and trace must move with the renamed intent.
The approved i003 contract now includes `agent-harness/` and `template/`.

The existing digest helper hashes the canonical `intent` field, so each approved
intent whose identity changes receives a new digest using the unchanged algorithm.
The current trace found no candidate-bearing intent, but authorization must be
checked after the migration before any plan proceeds.

## Product Knowledge and source grounding

- `wrapper/contracts/invariants.yaml` — owns intent anchoring, approval, feasibility,
  and the no-silent-lifecycle-change rules.
- `wrapper/contracts/schemas/intent-contract.yaml` — owns the canonical intent id
  grammar and contract shape.
- `wrapper/contracts/schemas/plan.yaml` — owns parent-intent references in plans.
- `wrapper/contracts/schemas/trace-manifest.yaml` — owns trace record paths and
  post-approval grounding evidence.
- `context/domains/intent/README.md` — describes intent identity and Gate 1.
- `context/domains/tracing/README.md` — describes trace freshness, question
  disposition, and feasibility before planning.
- `docs/planning.md` and `docs/templates/intent.md` — coordinator planning and
  intent-record guidance.
- `agents/coordinator.md`, `agents/tracer.md`, and `.agents/skills/` — host-facing
  lifecycle guidance and current intent references.
- `intent/i003-intent-id-width/trace/context-circuit-source.yaml` — exact traced
  call sites, reference inventory, risks, and executable done checks.

## Tasks

### 1. Update the owned identity contract and runtime mechanics

Change the canonical grammar to exactly three digits in the intent schema,
invariant, runtime validator, allocator, and all shared lifecycle consumers.
Keep plan ids four-digit. Add explicit overflow behavior at `i999`; do not add
automatic archive, delete, or id-recycling behavior. Update the shipped guidance
and contract assertions that describe the intent form.

**Acceptance:** new ids validate and old four-digit ids fail; allocation emits the
next three-digit id and fails clearly when no exact three-digit id is representable;
all lifecycle paths continue to use the shared intent helpers.

**Verification:** runtime id-shape and overflow tests; contract/schema checks;
`sh test/acceptance.sh`.

### 2. Migrate current intent records and relationships

Rename the three current intent directories and update their canonical `intent`
fields and index rows for all three existing intent records.
Move trace directories and preserve their evidence. Update every current plan
parent link, including this plan's parent link after the migration. Recompute the
frozen digest for the renamed approved records with the existing algorithm while
preserving status and criteria. Do not alter historical Git objects.

**Acceptance:** every active or archived intent directory matches its canonical
three-digit contract field; current traces, plans, indexes, relationships, status,
and approval state resolve through the new names; approved digest recomputation is
provable and authorization remains valid.

**Verification:** directory/contract/index consistency checks; digest recomputation
checks for renamed approved records; `intent-validate` for each migrated intent;
`intent-authorized` for each affected plan.

### 3. Migrate shipped fixtures, templates, and current references

Update deterministic test fixtures, dynamic plan-to-intent derivation, human
scenario setup, conversation plots/generated cases/scenarios, documentation, and
the shipped template. Replace only intent identifiers and intent derivation logic;
leave four-digit plan ids unchanged. Include a completeness check that finds no
old-form active reference outside explicitly historical grounding evidence.

**Acceptance:** every current workspace reference to the three existing intents
resolves to their canonical three-digit identifiers; dynamic fixture derivation emits the
three-digit form; plan ids remain four-digit; `agent-harness/` and `template/` are
covered.

**Verification:** the traced completeness commands; fixture and scenario suites;
template/runtime laboratory; stale-reference and plan-id regression checks.

### 4. Prove lifecycle preservation and manual cleanup boundaries

Add or update tests for allocation across active and archived intents, old-form
rejection, archive/restore with three-digit ids, approved digest behavior,
authorization, and explicit overflow at `i999`. Prove that no lifecycle operation
automatically archives, deletes, or recycles an intent, and run the full semantic
acceptance suite. Preserve the existing migration trace as evidence while excluding
historical evidence from active-reference checks where necessary.

**Acceptance:** the complete suite remains green; the migration adds no gate and
does not change authorization, tier, execution, verification, completion, or
delivery semantics; manual archive/delete remains human-triggered and identifiers
are not silently reused.

**Verification:**

```text
sh test/acceptance.sh
sh test/scenarios/test-scenarios.sh
sh test/intent/test-intent.sh
sh test/intent/test-feasibility.sh
sh agent-harness/test-template-runtime.sh
```

## Acceptance and verification summary

- Exactly three digits are accepted after `i`; four-digit intent ids are rejected.
- Current intent directories, contract fields, indexes, traces, plan parent links,
  and current references migrate consistently.
- Approved frozen digests are recomputed with the existing algorithm.
- The next available id is allocated across active and archived records, with a
  clear exact-three-digit exhaustion failure at `i999`.
- Archive/delete is never automatic and never silently recycles an identifier.
- Plan identifiers remain four-digit and all lifecycle behavior remains unchanged.
- The full semantic acceptance suite remains green.

## Assumptions, open questions, and risks

- The human has answered the scope question by requiring every current workspace
  reference, including the harness and template.
- The human has chosen existing-algorithm digest recomputation for approved intents.
- The human has chosen explicit human/manual archive or deletion for old or
  unrelated intents; no automatic cleanup or recycling is authorized.
- The exact-three-digit namespace has a hard ceiling at `i999`; allocation must
  fail clearly at that ceiling rather than invent a fourth digit.
- The trace manifest records old-form matches as historical evidence. Active
  reference checks must exclude or classify that evidence without hiding any live
  reference.
- No candidate-bearing intent was found during tracing; if a candidate appears
  before execution, authorization and evidence invalidation must be handled before
  proceeding.

## Expected Product Knowledge impact

The owned intent, tracing, and planning guidance will be reconciled to the new
identifier grammar and question-disposition rule. No new Product Knowledge domain
is required.

## Delivery notes

This plan does not authorize delivery. Completion, commit, push, pull request,
merge, or publication remain separate human-requested actions.
