# 0001 — Establish the conversation-spec library as the living harness specification

- **Plan ID:** `0001-conversation-spec-library`
- **Intent:** `i001-conversation-spec-library`
- **Status:** draft
- **Repository:** `context-circuit-source`
- **Assurance:** Standard — one worker and one independent verifier

## Original request and coverage

Make `agent-harness/conversations/` the living, descriptive specification and
generative source of the agent harness. The plots must describe the complete
lay-user surface, generate the cases that the live harness executes, and keep
the plot, generated case, and executed scenario from silently drifting apart.

This covers the complete library lifecycle: the external-surface inventory and
plot format; deterministic plot-only generation; artifact-aware grading and
fixtures; conversion of the existing 22 scenario cases without changing their
calibrated behavior; the frozen net-new plot set; live execution; and the
acceptance/drift bridge. It does not cover product-rule changes or delivery.

## Objective and desired behavior

- Each conversation plot has a machine-readable `## Spec` and a readable
  `## Dialogue`, with faithful actors, turns, decision points, and terminal
  outcome.
- Every plot is descriptive: it cites current `INV-*` and/or `AC-*` ids and
  does not create a parallel product rule. Coordinator expectations, including
  expected intent, trace, and plan artifacts, live in `coordinator_must`; the
  simulated human Dialogue remains free of internal filenames and mechanics.
- The generator reads plots only and materializes deterministic case artifacts.
  `conversations/generated/<plot-id>.case.yaml` is the canonical generated
  artifact, and `scenarios/<plot-id>/case.yaml` is the runner's identity-
  preserving projection produced by the same generator. Neither is
  hand-authored; the suite checks their bijection and content relationship.
- The fixed environmental overlay contains only the host lane and per-host
  budgets. All other case semantics derive traceably from the plot, and
  regeneration is byte-deterministic.
- The harness executes the generated scenario projection against the live
  coordinator and checks grounded artifacts, dynamic entity binding,
  individually falsifiable post-conditions, host-aware access evidence, and
  the required negative fixtures.
- The existing 22 scenario identities retain their calibrated seeds, host
  budgets, and assertion strength. AC-36's deterministic five-family suite is
  unchanged and remains green.
- Drift is mechanically rejected across the external inventory, plots,
  generated artifacts, and live scenarios. The frozen net-new successful cases
  reach their expected live grade when their host lanes are available.

## Constraints and non-goals

- All implementation changes stay under `agent-harness/`.
- Do not change `wrapper/`, `template/`, `sources/`, the invariant owner map,
  or the externally owned acceptance-criteria map.
- The library remains descriptive and references existing `INV-*` and `AC-*`
  rules; it does not define product behavior or introduce a new product
  contract field.
- Preserve the 22 existing scenario directory identities, calibrated seeds,
  per-host budgets, and AC-36 behavior. The deferred rename, diagram-format,
  and flow-metadata organization enhancements remain out of scope.
- The environmental overlay is limited to the enumerated host lane and
  per-host budgets; no nondeterministic or semantic field may be placed there.
- The net-new plot set is frozen when this plan is implemented. Successful
  lifecycle completions cannot be reclassified as expected failures to evade
  the live-execution floor.
- No delivery, merge, push, publication, deployment, or typed-language port is
  included. No AI attribution is added to any source-only commit.

## Grounded findings

The approved Standard intent is feasible and scoped to `agent-harness/`. The
existing trace identifies 39 plots, 39 generated-looking case files, and 22
currently executed scenario cases. The repository has no tracked generator or
drift guard, and the live runner currently resolves cases from
`agent-harness/scenarios/` rather than from the generated directory.

The trace maps the main change surfaces to the conversation documentation and
plots, generator and case trees, human runner/grader/matrix, host drivers and
fixtures, the in-scope acceptance bridge, and the read-only product/acceptance
owners cited by the plots. Existing plots already carry coordinator
expectations, Dialogue annotations, and post-conditions, but generic checks
for intent/trace/plan artifacts, dynamic ids, falsifiability, and plot/case
drift are absent. Several plots and coverage documents contain stale v1.0
operation names or lifecycle wording and need refitting before generation is
made authoritative.

The latest decomposition guidance confirms that task count and assurance tier
do not determine plan count. This work remains one plan: every partition is in
the same repository and bounded `agent-harness/` surface, the generator and
runner changes must agree as one change-set, and the work is best executed and
independently verified through one Standard worker/verifier lifecycle. The
partition dependencies are therefore intra-plan task dependencies; no
`plan_dependencies` edge is needed.

## Product Knowledge and source grounding

- `wrapper/contracts/invariants.yaml` — owner map for all product rules that
  descriptive plots may cite; read-only.
- `test/acceptance/criteria-map.yaml` — externally owned `AC-*` acceptance
  surface and suite ownership; read-only.
- `context/ARCHITECTURE.md` — current coordinator, worker, verifier, and
  assurance ownership.
- `context/domains/intent/README.md` — Gate 1, intent artifacts, and the
  intent-to-plan relationship.
- `context/domains/tracing/README.md` — post-approval tracing, feasibility,
  manifest evidence, and question classification.
- `context/domains/assurance/README.md` — Explore/Standard/Critical ladder and
  the independent-verifier expectation.
- `context/domains/plan-authorization/README.md` — intent-derived plan
  authorization and the absence of a second plan-approval gate.
  - `context/domains/plan-review/README.md` — plan derivation, grounding, task
  dependencies, and plan-record responsibilities.

## Trace questions and dispositions

All trace questions are implementation-level and are resolved here without
changing the approved intent:

1. **Generated/live relationship — plan resolution.** The plot is the only
   source. One generator run emits the canonical file under
   `conversations/generated/` and the same identity-preserving projection under
   `scenarios/`; the runner consumes the latter, while hash/bijection checks
   prove it is generated output. Existing 22 scenario directory names remain
   unchanged.
2. **Artifact expectations — already answered by the approved intent and
   existing plot conventions.** Keep expected coordinator outputs in
   `coordinator_must`, express the machine-checkable evidence through generated
   state/post-condition checks, do not add a product contract field, and keep
   internal filenames out of Dialogue.
3. **Feasibility fixtures — plan resolution.** Refit them to the current v1.0
   sequence: approval comes before tracing/feasibility, and tracing occurs
   before plan derivation. A newly discovered question is represented in the
   post-approval trace outcome; no fixture may imply the retired pre-plan
   adversary or an automatic system-design intent spawn.

No intent-level question or out-of-scope required write remains.

## Tasks

### 1. Define the plot source and refit stale v1.0 conversations

- **Task ID:** `CSL-001`
- **Repository:** `context-circuit-source`
- **Paths:**
  `agent-harness/conversations/README.md`,
  `agent-harness/conversations/FLOW.md`,
  `agent-harness/conversations/coverage/`,
  `agent-harness/conversations/plots/`,
  `agent-harness/README.md`,
  `agent-harness/scenarios/README.md`
- **Depends on:** none

Establish the externally anchored inventory and parser-facing conventions for
`## Spec`/`## Dialogue`. Refit stale feasibility, Product Knowledge,
system-design, delivery, tier, legacy-label, count, and acceptance-mapping
content to v1.0. Make `coordinator_must` explicit about expected coordinator
outputs, including intent artifacts and downstream trace/plan outcomes, while
keeping human-facing Dialogue plain and internal-name-free. Freeze and record
the net-new plot set.

Acceptance: every inventory entry has one complete plot with Spec and Dialogue;
every plot names coordinator expectations, reporting rules, citations, and a
terminal outcome; intent-flow plots cover expected artifacts without leaking
internals; and current invariant/acceptance citations and lifecycle wording
resolve cleanly.

Verification:

- `sh agent-harness/conversations/test-library.sh --check whole-surface`
- `sh agent-harness/conversations/test-library.sh --check surface-coverage`
- `sh agent-harness/conversations/test-library.sh --check descriptive-only`
- `sh agent-harness/conversations/test-library.sh --check reporting-rules`
- `sh agent-harness/conversations/test-library.sh --check dialogue-faithful`

### 2. Generate cases deterministically and guard inventory-to-live drift

- **Task ID:** `CSL-002`
- **Repository:** `context-circuit-source`
- **Paths:**
  `agent-harness/conversations/generate.sh`,
  `agent-harness/conversations/test-library.sh`,
  `agent-harness/conversations/generated/`,
  `agent-harness/scenarios/`
- **Depends on:** `CSL-001`

Implement plot-only generation with full semantic-field propagation, explicit
decorative-field handling, the bounded host/budget overlay, and stable ordering.
Emit the canonical generated artifacts and runner projections while preserving
the 22 existing scenario identities. Enforce inventory↔plot↔generated↔live
bijective wiring, orphan detection, and hand-edit/plot-edit drift failures.
Include full-set forward mutation and decorative reverse-mutation checks.

Acceptance: two generations are byte-identical outside the overlay; no
committed case is read as generator input; every non-decorative Spec field is
traceably load-bearing; decorative fields are proven inert; the overlay is
exactly the fixed set; and all inventory/plot/generated/live mismatches fail.

Verification:

- `sh agent-harness/conversations/generate.sh --check-all`
- `sh agent-harness/conversations/test-library.sh --check overlay`
- `sh agent-harness/conversations/test-library.sh --check wiring`
- `sh agent-harness/conversations/test-library.sh --check drift`

### 3. Add artifact-aware grading and required harness capabilities

- **Task ID:** `CSL-003`
- **Repository:** `context-circuit-source`
- **Paths:**
  `agent-harness/human/run-scenario.sh`,
  `agent-harness/human/grade.sh`,
  `agent-harness/human/run-matrix.sh`,
  `agent-harness/human/drivers/`,
  `agent-harness/fixtures/`,
  `agent-harness/human/test-*.sh`
- **Depends on:** `CSL-002`

Add the setup and grading capabilities required by the plot set: intent,
trace, and plan artifact/state predicates; id-agnostic entity binding; wrong
entity/kind/count negatives; individually falsifiable post-condition fixtures;
new seeds, faults, and state predicates; and host-aware access auditing.
Preserve disposable workspace setup, host distinctions, and the Codex/Cursor
access-audit degradation. Keep credentials and provider payloads out of all
harness artifacts.

Acceptance: intent-flow cases verify expected artifacts without fixed generated
ids; every generated post-condition has a targeted failing negative; every
enumerated new fixture is exercised and has a failing negative; and host
evidence remains honest.

Verification:

- `sh agent-harness/conversations/test-library.sh --check id-agnostic`
- `sh agent-harness/conversations/test-library.sh --check falsifiable`
- `sh agent-harness/conversations/test-library.sh --check new-fixtures`

### 4. Rebuild the corpus and wire the live execution floor

- **Task ID:** `CSL-004`
- **Repository:** `context-circuit-source`
- **Paths:**
  `agent-harness/conversations/generated/`,
  `agent-harness/scenarios/`,
  `agent-harness/conversations/coverage/`,
  `agent-harness/test-template-runtime.sh`,
  `agent-harness/human/`
- **Depends on:** `CSL-003`

Regenerate all 22 existing cases from their plots and prove seeds, budgets,
identities, and assertion strength are preserved. Add and execute the frozen
net-new cases through the existing acceptance bridge and live coordinator when
host lanes are available. Keep AC-36's deterministic five-family suite
unchanged, and do not edit `test/acceptance.sh` or the external criteria map.

Acceptance: the existing-case equivalence and live-execution floors hold; every
frozen net-new successful completion reaches its expected live grade; expected
failure plots assert their specific triaged mode; AC-36 remains unchanged and
green; and the implementation diff is fenced to `agent-harness/`.

Verification:

- `sh agent-harness/conversations/test-library.sh --check regeneration-equivalence`
- `sh agent-harness/conversations/test-library.sh --check corrections`
- `sh test/scenarios/test-scenarios.sh`
- `sh test/acceptance.sh`
- `sh agent-harness/human/run-matrix.sh`

## Full acceptance and verification matrix

The plan carries every executable done-check from the refreshed trace:

```text
sh agent-harness/conversations/test-library.sh --check whole-surface
sh agent-harness/conversations/test-library.sh --check surface-coverage
sh agent-harness/conversations/test-library.sh --check descriptive-only
sh agent-harness/conversations/test-library.sh --check reporting-rules
sh agent-harness/conversations/test-library.sh --check dialogue-faithful
sh agent-harness/conversations/generate.sh --check-all
sh agent-harness/conversations/test-library.sh --check overlay
sh agent-harness/conversations/test-library.sh --check wiring
sh agent-harness/conversations/test-library.sh --check id-agnostic
sh agent-harness/conversations/test-library.sh --check falsifiable
sh agent-harness/conversations/test-library.sh --check new-fixtures
sh agent-harness/conversations/test-library.sh --check regeneration-equivalence
sh agent-harness/conversations/test-library.sh --check corrections
sh agent-harness/conversations/test-library.sh --check drift
sh test/scenarios/test-scenarios.sh
sh test/acceptance.sh
sh agent-harness/human/run-matrix.sh
```

The scope proof is read-only and must be evaluated against the implementation
base/tip by the verifier:

```text
git diff --name-only "$WORKER_BASE" "$WORKER_TIP" | awk 'index($0, "agent-harness/") != 1 {bad=1} END {exit bad}'
```

## Assumptions and risks

Assumptions:

- The plot-only generator can emit both the canonical generated artifact and
  the runner projection without changing scenario directory identities.
- Artifact expectations remain expressed in `coordinator_must` and executable
  generated state/post-condition checks; no product contract field is needed.
- The existing acceptance bridge in `agent-harness/test-template-runtime.sh`
  is sufficient to gate library checks without touching `test/acceptance.sh`.
- Claude Code access traces remain dependable, while Codex and Cursor retain
  their explicitly degraded access-audit behavior.
- Live host execution may be unavailable in a credential-free environment;
  deterministic checks remain runnable, and the live floor reports required
  host availability honestly.

Risks:

- The present 39 generated-looking cases are not evidence of generation; the
  implementation must prove that no committed case is used as input.
- The runner's current `scenarios/` topology and the generated directory can
  otherwise become competing sources; wiring and content-hash checks must make
  the projection relationship unambiguous.
- Full-field mutation and per-condition negative fixtures may expose plot
  fields or grader predicates that are currently implicit; those gaps must be
  resolved within `agent-harness/` without weakening the acceptance suite.
- Refit work could accidentally restore retired adversary, scope-envelope, or
  automatic-intent-spawn semantics; the v1.0 sequence and descriptive-only
  checks must guard against that regression.
- The five accepted residual loopholes in the intent contract remain known
  risks and are not silently turned into new product criteria by this plan.

## Expected commits and Product Knowledge impact

One or more source-only Conventional Commits are expected while preserving the
task order, using `type(scope): imperative subject` and no AI attribution. The
implementation is not authorized to commit merely by deriving this plan; any
explicit source-only commit is a separate human-requested delivery action.

No durable Product Knowledge is written by this plan. The existing pending
conversation-spec-library proposal may be reconciled later, but that is outside
implementation and acceptance.

## Delivery notes

This plan does not authorize delivery. Completion, commit, push, pull request,
merge, publication, and deployment remain separate explicit actions.
