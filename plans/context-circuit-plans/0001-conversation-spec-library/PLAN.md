# 0001 — Establish the conversation-spec library as the living harness specification

- **Plan ID:** `0001-conversation-spec-library`
- **Intent:** `i0001-conversation-spec-library`
- **Status:** draft
- **Repository:** `context-circuit-source`
- **Assurance:** Standard — one worker and one independent verifier

## Original request

Make `agent-harness/conversations/` the living, descriptive specification and
generative source of the agent harness. The plots must describe the complete
lay-user surface, generate the live cases, and keep the plot, generated case,
and executed scenario from silently drifting apart.

## Objective and desired behavior

- Each conversation plot has a machine-readable `## Spec` and a readable
  `## Dialogue`.
- Each plot states the coordinator's expected behavior and output in its
  `decision_points.coordinator_must` descriptions. For intent creation, this
  includes producing the human-readable intent record and machine-readable
  intent record; the user-facing Dialogue still hides internal filenames.
- The generated case is derived from the plot, not hand-authored or copied as
  an independent source.
- The harness executes the generated cases against the live coordinator and
  checks both conversation behavior and grounded end state.
- The existing 22 scenario cases retain their calibrated seeds, host budgets,
  and assertion strength while becoming generated projections.
- Drift is mechanically rejected across the external surface inventory, plots,
  generated cases, and live scenarios.

## Constraints and non-goals

- All implementation changes stay under `agent-harness/`.
- Do not change `wrapper/`, `template/`, `sources/`, the invariant owner map,
  or the acceptance-criteria map.
- The library remains descriptive: it cites `INV-*` and `AC-*` but defines no
  product rule of its own.
- Preserve the 22 existing scenario identities and calibrated seeds/budgets;
  do not recreate them from scratch.
- The environmental overlay is limited to the enumerated host lane and
  per-host budgets.
- No delivery, merge, push, publication, or deployment is included.

## Grounded findings

The approved intent is Standard and feasible. The tracer found 39 plots, 39
generated case files, and 22 currently executed scenario cases. The current
runner executes `agent-harness/scenarios/*/case.yaml`; it does not consume the
generated tree, and no tracked generator or drift guard currently exists.

The current plots already express coordinator expectations through
`coordinator_must`, Dialogue annotations, and post-conditions, but artifact
outputs such as the intent files, trace manifest, and derived plan are not yet
generically machine-checked. The feasibility plots also need to be realigned
to the post-approval, pre-plan tracing order.

## Product Knowledge and source grounding

- `wrapper/contracts/invariants.yaml` — owner map for all product rules the
  descriptive plots cite.
- `test/acceptance/criteria-map.yaml` — externally owned AC-* acceptance
  surface.
- `context/ARCHITECTURE.md` — current coordinator, worker, verifier, and
  assurance model.
- `context/domains/intent/README.md` — Gate 1 and intent-to-plan relationship.
- `context/domains/tracing/README.md` — post-approval tracing and feasibility.
- `context/domains/assurance/README.md` — Explore/Standard/Critical ladder.
- `context/proposals/0034-add-conversation-spec-library.md` — pending
  description of the conversation-spec library and its intended guardrails.
- `sources/context-circuit-v1.0/design.md` — named v1.0 design source for the
  two-gate lifecycle, candidate evidence, and reconciliation behavior.

## Tasks

### 1. Define and refit the conversation source

Create the external-surface inventory and parser-facing conventions for plots.
Refit stale plots, coverage documents, and harness orientation to current v1.0
semantics. Make `coordinator_must` descriptions explicit about expected
coordinator outputs, including intent artifacts, while keeping internal names
out of the human-facing Dialogue. Correct the feasibility fixtures, stale
operation names, legacy labels, delivery wording, and over-broad AC mappings.

### 2. Build deterministic generation and drift protection

Add the generator and structural checks under `agent-harness/`. It must read
plots only, propagate every semantic Spec field into the generated case, keep
the overlay bounded, be byte-deterministic, and reject inventory/plot/case/live
scenario disagreement. The canonical relationship between
`conversations/generated/*.case.yaml` and `scenarios/*/case.yaml` must be made
explicit rather than leaving two competing outputs.

### 3. Extend grading and fixtures for grounded outputs

Add the harness capabilities required by the plot set: artifact-aware intent
outcomes, dynamic entity IDs, trace and plan state checks, falsifiable
post-conditions, negative fixtures, new seed states, and host-aware access
auditing. Preserve the existing host distinctions and disposable-workspace
boundary.

### 4. Rebuild and wire the corpus

Regenerate the 22 existing cases from their plots, preserve their calibrated
inputs, and wire the live runner to generated cases. Add the frozen net-new
conversation cases and gate them through the existing in-scope acceptance
bridge. Keep AC-36's deterministic five-family suite unchanged while adding
the library checks and live-execution floor.

## Acceptance and verification

- Every inventory entry has exactly one plot and one generated case.
- Every plot has faithful Dialogue/spec alignment, valid citations, and
  explicit coordinator expectations.
- Intent-flow cases describe and check the expected intent artifacts and later
  trace/plan outputs without leaking internal details to the simulated human.
- Regeneration is deterministic and reads no committed generated case as input.
- The overlay contains only host lane and per-host budgets.
- Live scenarios are a bijective generated projection; hand-edited or orphaned
  cases fail the suite.
- Every generated post-condition is individually falsifiable.
- Existing 22 cases preserve seeds/budgets and AC-36 remains unchanged.
- Net-new successful conversation cases pass against the live coordinator when
  the required host lanes are available.
- The diff remains fenced to `agent-harness/`.

Verification commands are carried forward from the trace manifest:

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

## Assumptions, open questions, and risks

- The implementation must choose and document one canonical relationship
  between the generated directory and the live scenario directory while
  preserving the 22 scenario directory identities.
- Artifact expectations are expressed in plot `coordinator_must` descriptions
  and enforced by generated state/post-condition checks; no new product
  contract field is introduced.
- The feasibility plots must be resolved as either genuinely code-discovered
  post-approval questions or pre-Gate-1 questions; the current mixed fixtures
  cannot remain ambiguous.
- Claude Code exposes dependable file-access traces; Codex and Cursor do not,
  so access assertions must remain host-aware.
- Live host execution depends on available CLIs and nested-agent support;
  deterministic checks must remain runnable without credentials or provider
  payloads.

## Expected Product Knowledge impact

No durable Product Knowledge is written by this plan. The work may later
support reconciliation of the pending conversation-spec-library proposal, but
that is separate from implementation and acceptance.

## Delivery notes

This plan does not authorize delivery. Completion and any later commit, push,
pull request, merge, or publication remain separate human-requested actions.
