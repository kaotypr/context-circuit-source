# Sol 5.6 Plan — Context Circuit routing and context-efficiency redesign

Status: **draft — requires human review and approval before execution**

Scope: the Context Circuit wrapper source repository and the workspace template
that it publishes. This plan changes the product's instruction, routing,
documentation, test, and release surfaces. It does not authorize implementation,
merge, publication, deployment, or a canonical plan status change.

This standalone filename was explicitly requested for design comparison. The
current checkout is also an uninitialized Context Circuit template, so this file
does not pretend to be a canonical `plan.yaml`. Before execution, the human must
choose whether to convert the accepted work into a canonical plan bundle under
`plans/context-circuit-plans/` or treat it as ordinary wrapper-repository work.

## 1. Objective

Redesign Context Circuit so that agents:

1. spend substantially fewer input tokens on wrapper instructions;
2. reach one deterministic, explainable route from user intent and workspace
   state;
3. load only the context required for that route;
4. encounter one canonical owner for each normative rule;
5. distinguish wrapper source, released template, instantiated workspace state,
   host tasks, Context Circuit sessions, plans, and plan tasks;
6. preserve the current safety model unless a later human-approved plan changes
   it explicitly.

## 2. Source and provenance

This draft is grounded in:

- the user's request and the discussion in Codex task
  `01a022ee-0fd8-7482-8866-4a9906fedd68`;
- measured sizes of the current always-read, root-entry, run-plan, writer, and
  verifier document sets;
- `AGENTS.md`, `WORKFLOW.md`, `workspace.yaml`, and `context/INDEX.md`;
- `context/WORKSPACE.md`, `context/PROJECT.md`, `context/ARCHITECTURE.md`,
  `context/CONVENTIONS.md`, and `context/DECISIONS.md`;
- `docs/agent-workspace-workflow.md`, `docs/runtime-contract.md`,
  `docs/planning.md`, `docs/getting-started.md`, and
  `docs/host-capabilities.md`;
- `agents/coordinator.md`, `agents/repository-worker.md`, and
  `agents/reviewer.md`;
- `cc-session-entry`, `cc-whats-next`, `cc-run-plan`, and `cc-run-stack` skill
  contracts;
- the current acceptance and release implementation;
- `Opus-4.8-plan.md`, read only to understand the existing parallel proposal
  and avoid overwriting or accidentally duplicating its artifact.

No file under `sources/` was read. `Opus-4.8-plan.md` is comparative evidence,
not accepted Product Knowledge or an instruction source.

## 3. Observed baseline

The following values are measured from the current checkout. Token ranges use
a conservative estimate of 3.5–5 bytes per model token; acceptance will record
byte and word counts even when an exact model tokenizer is unavailable.

| Context profile | Words | Bytes | Estimated tokens |
| --- | ---: | ---: | ---: |
| Always-read bootstrap | 3,079 | 22,349 | 4.5–6.4k |
| Full root entry | 10,547 | 76,456 | 15.3–21.8k |
| Root entering `cc-run-plan` | 12,950 | 93,873 | 18.8–26.8k |
| Writer wrapper context | 4,075 | 30,387 | 6.1–8.7k |
| Verifier wrapper context | 4,033 | 30,108 | 6.0–8.6k |
| One-plan static total | — | 154,368 | 30.9–44.1k |
| All 15 skill bodies | 7,252 | 51,146 | 10.2–14.6k |

Normative duplication is also measurable:

- the writer-child rule appears in 12 Markdown documents;
- the standard single-plan-entry rule appears in 7 documents;
- the no-global-current-session rule appears in 4 documents;
- human-gate language appears in at least 18 documents, with both necessary
  route-specific detail and repeated general policy mixed together.

The current source checkout has another structural ambiguity: it is the Context
Circuit product repository while `workspace.yaml` and `context/PROJECT.md`
describe an uninitialized template. The release process packages much of the
same root tree as that template. Agents therefore have to infer whether they are
developing the wrapper or operating a workspace created from it.

## 4. Quantitative targets

### 4.1 Context and token budgets

| Context profile | Hard byte budget | Approximate token target | Reduction target |
| --- | ---: | ---: | ---: |
| Bootstrap | 8 KB | 0.9–2.3k | 65–85% |
| Full root entry before task evidence | 16 KB | 3.2–4.6k | 70–85% |
| Root `cc-run-plan` wrapper context | 20 KB | 4–5.7k | 70–85% |
| Writer wrapper context | 10 KB | 2–2.9k | 55–75% |
| Verifier wrapper context | 9 KB | 1.8–2.6k | 55–75% |
| One-plan static total | 39 KB | 7.8–11.2k | 65–80% |

Expected practical outcome:

- save approximately 20–36k static wrapper tokens per one-plan run;
- reduce total input usage by 50–70% after task-specific context selection;
- keep every substantial skill body below 700 words, with a preferred budget
  below 450 words;
- require a recorded reason when a route exceeds its context budget.

### 4.2 Routing and reliability targets

- one canonical route-decision table;
- one canonical definition for every registered invariant;
- at least 50 filesystem-state route fixtures with 100% deterministic expected
  route results;
- at least 50 natural-language routing cases with at least 95% correct route
  selection and no more than 5% ambiguous outcomes;
- exact session lookup in one record read when a host binding is present;
- no route may infer authorization merely from eligibility;
- zero regressions in plan gates, leases, worktree exclusivity, archive
  eligibility, verifier independence, recovery, and cleanup protections.

## 5. Target architecture

### 5.1 Authority layers

| Layer | Responsibility | Authority rule |
| --- | --- | --- |
| `AGENTS.md` | Small host bootstrap, safety, precedence, source boundary | Normative and always loaded |
| `WORKFLOW.md` | Compact lifecycle and link to canonical router | Normative summary, not a duplicate contract |
| `wrapper/manifest.yaml` | Wrapper version, catalog, budgets, schema locations | Canonical machine-readable wrapper catalog |
| `wrapper/contracts/routing.yaml` | Route conditions, precedence, capabilities | Sole routing authority |
| `wrapper/contracts/invariants.yaml` | Stable safety and lifecycle invariant IDs | Sole shared-rule authority |
| `wrapper/contracts/context-sets.yaml` | Bootstrap and route-specific read sets | Sole context-selection authority |
| `wrapper/contracts/schemas/` | Runtime and durable artifact shapes | Canonical record schemas |
| `.agents/skills/` | Host discovery and operation-specific procedure | Thin adapters referencing contracts |
| `agents/` | Role-specific deltas | Thin role adapters referencing contracts |
| `docs/` | Explanation, examples, maintenance guidance | Informative unless explicitly designated |
| `workspace.yaml`, `context/`, `plans/`, `.runtime/` | Mutable instantiated workspace state | Existing state-specific authorities |

Static wrapper metadata belongs in one catalog. Do not add repetitive
frontmatter to every static document. Dynamic Product Knowledge may keep small
metadata only when provenance, freshness, or acceptance state requires it.

### 5.2 Routing model

Routing must separate concepts that are currently blended together:

```yaml
schema_version: 1
intent: execute-plan
session_kind: root
session_phase: orienting
eligibility: ready
next_capability: cc-run-plan
authorization: explicitly-requested
reason_codes:
  - PLAN_APPROVED
  - DEPENDENCIES_READY
  - NO_LIVE_OWNER
context_set: run-plan
human_gate: none
```

- `intent` records what the human requested.
- `session_phase` records current runtime lifecycle.
- `eligibility` records whether state permits the candidate action.
- `next_capability` identifies the single selected operation.
- `authorization` distinguishes recommendation from permission to mutate.
- `reason_codes` make the result testable and explainable.
- `context_set` selects the bounded read packet.

`cc-session-entry` evaluates this contract. `cc-whats-next` is a read-only view
of the same decision and must not maintain a second routing model.

### 5.3 Progressive context loading

```text
User intent + host identity
          ↓
Bounded state probe
          ↓
Structured route decision
          ↓
Named context set
          ↓
Route-specific skill and evidence
          ↓
Compact handoff with exact continuation references
```

Context sets define:

- required files or record fields;
- conditional deep reads and their reason codes;
- prohibited broad scans such as unselected `sources/` content;
- maximum words and bytes;
- evidence that must be carried into a child packet;
- wrapper contract version used to make the decision.

Runtime sections should be split by concern so an ordinary orientation does
not load lease acquisition, stack joins, cleanup, and takeover rules together.

### 5.4 Host and session identity

Add an optional portable binding to runtime session records:

```yaml
host_binding:
  provider: codex
  task_id: 01a022ee-0fd8-7482-8866-4a9906fedd68
```

The binding is not a global current-session pointer and does not override the
Context Circuit session ID. It enables exact lookup and makes documentation
distinguish:

- host task or thread;
- Context Circuit root or child session;
- plan execution;
- task within a plan;
- worktree ownership.

### 5.5 Source repository and released template

Separate maintainer state from the blank workspace distributed to users:

```text
Context Circuit source repository
├── wrapper/                 canonical shipped wrapper contracts
├── template/                uninitialized mutable workspace seed
├── docs/                    maintainer and explanatory documentation
├── test/
├── scripts/
└── source workspace state   describes Context Circuit itself

Published artifact
├── AGENTS.md / WORKFLOW.md  thin entry adapters
├── wrapper/                 immutable versioned contract
├── .agents/skills/          host discovery adapters
├── workspace.yaml           copied from template seed
├── context/                 copied from template seed
├── sources/
└── plans/
```

The release build must verify both the source contract and the staged artifact.
Stable user-facing root paths should remain compatible during migration.

## 6. Implementation tasks

### SOL-001 — Freeze the baseline and build the evaluation corpus

Dependencies: none.

Paths:

- `test/baselines/`
- `test/fixtures/routing/`
- `test/context-budget.sh`
- `docs/benchmarks/`

Work:

- record current byte and word counts for each context profile in Section 3;
- record every normative duplication family and current owner candidates;
- create at least 50 filesystem route fixtures covering fresh root, child,
  resume, draft, approved, done, archived, malformed, blocked dependency, live
  lease, stale lease, completion, cleanup, single-plan, and stack cases;
- create a natural-language route corpus with expected intent and capability;
- make measurement independent of an exact model tokenizer.

Acceptance:

- baseline values reproduce within documented newline/platform tolerance;
- every fixture has a unique expected route or an explicit ambiguous result;
- the task changes no normative behavior.

Verification:

- `sh test/context-budget.sh --baseline`
- `sh test/acceptance.sh`
- `git diff --check`

Expected evidence: baseline YAML/Markdown, fixture inventory, command output.

Stop conditions: a profile cannot be reproduced from explicit file paths, or
an existing route case has no agreed expected result.

### SOL-002 — Define the canonical wrapper contract

Dependencies: SOL-001.

Paths:

- `wrapper/manifest.yaml`
- `wrapper/contracts/routing.yaml`
- `wrapper/contracts/invariants.yaml`
- `wrapper/contracts/context-sets.yaml`
- `wrapper/contracts/schemas/`

Work:

- assign stable IDs to current safety, ownership, lifecycle, and human-gate
  invariants without changing their meaning;
- define the route-decision schema and explicit precedence;
- define bootstrap, orientation, planning, run-plan, run-stack, verification,
  finish, archive, cleanup, and recovery context sets;
- classify every existing normative statement as canonical, route-specific,
  explanatory, obsolete duplication, or unresolved conflict;
- record compatibility and contract versioning rules.

Acceptance:

- every existing safety rule maps to exactly one invariant ID;
- every shipped capability maps to exactly one route;
- intent, eligibility, authorization, and session phase have separate fields;
- no existing human gate or execution constraint changes semantics.

Verification:

- schema and reference checks over `wrapper/`;
- all SOL-001 route fixtures resolve to their expected decisions;
- `sh test/acceptance.sh` and `git diff --check`.

Expected evidence: invariant coverage table, routing matrix, schema-validation
results.

Stop conditions: two existing documents imply materially different safety
semantics, or route precedence requires a new human policy decision.

### SOL-003 — Implement one router and thin route adapters

Dependencies: SOL-002.

Paths:

- `AGENTS.md`
- `WORKFLOW.md`
- `.agents/skills/cc-session-entry/SKILL.md`
- `.agents/skills/cc-whats-next/SKILL.md`
- `agents/coordinator.md`
- relevant routing tests

Work:

- make `cc-session-entry` the sole evaluator of the routing contract;
- make `cc-whats-next` present the same result without duplicating conditions;
- replace mixed route/status vocabulary with the schema enums;
- emit compact reason codes, selected context set, authorization state, and
  human gate;
- keep the root response human-readable without requiring raw YAML output.

Acceptance:

- exactly one canonical routing table exists;
- all filesystem route fixtures pass deterministically;
- recommendation never implies mutation authorization;
- malformed or contradictory state returns a visible recovery decision.

Verification:

- route fixture suite;
- duplicate-route-table lint;
- `sh test/acceptance.sh` and `git diff --check`.

Expected evidence: before/after routing map, fixture results, ambiguity report.

Stop conditions: the host cannot expose enough identity to distinguish root
from delegated child, or a route would silently satisfy a human gate.

### SOL-004 — Add exact host-session binding and bounded discovery

Dependencies: SOL-002, SOL-003.

Paths:

- `wrapper/contracts/schemas/session.yaml`
- `docs/runtime-contract.md`
- `docs/host-capabilities.md`
- session and recovery fixtures

Work:

- add optional `host_binding` without changing portable session identity;
- specify exact lookup, missing binding, duplicate binding, and host migration
  behavior;
- preserve the prohibition on a global current-session pointer;
- define bounded fallback discovery when no binding exists;
- use consistent terms for host task, Context Circuit session, plan, task, and
  worktree.

Acceptance:

- a valid binding resolves with one session-record read;
- duplicate or conflicting bindings block instead of guessing;
- sessions without bindings remain backward compatible;
- host metadata cannot authorize takeover or mutation.

Verification:

- session lookup and recovery fixtures;
- schema compatibility checks;
- `sh test/acceptance.sh` and `git diff --check`.

Expected evidence: lookup-count results, compatibility fixtures, terminology
lint.

Stop conditions: a supported host has no stable task identifier or storing the
identifier violates an established privacy requirement.

### SOL-005 — Compile route-specific context and shrink the wrapper surface

Dependencies: SOL-002, SOL-003.

Paths:

- `AGENTS.md`
- `WORKFLOW.md`
- `context/INDEX.md`
- `docs/agent-workspace-workflow.md`
- `docs/runtime-contract.md`
- `docs/planning.md`
- `docs/getting-started.md`
- `.agents/skills/*/SKILL.md`
- `agents/*.md`
- `wrapper/contracts/context-sets.yaml`

Work:

- reduce the always-read spine to the bootstrap budget;
- split deep runtime guidance by route-selectable concern;
- make `context/INDEX.md` navigation-only rather than another fact store;
- replace repeated normative prose with invariant references and concise local
  consequences;
- keep skills focused on triggers and unique procedure;
- ensure child packets carry exact context references, invariant IDs, wrapper
  version, acceptance criteria, and stop conditions;
- preserve direct reads of original evidence rather than copying large source
  summaries into packets.

Acceptance:

- every context profile meets the Section 4 hard byte budget;
- writer-child normative ownership reduces from 12 copies to 1;
- single-plan-entry normative ownership reduces from 7 copies to 1;
- each substantial skill is at most 700 words and preferably below 450;
- no route loses a safety invariant required by its previous contract.

Verification:

- `sh test/context-budget.sh`;
- invariant ownership and broken-reference lint;
- route-by-route context-set audit;
- `sh test/acceptance.sh` and `git diff --check`.

Expected evidence: before/after profile table, invariant coverage report, skill
size report.

Stop conditions: meeting a budget would require deleting a safety rule rather
than routing it conditionally, or a skill cannot remain understandable from its
declared context set.

### SOL-006 — Separate wrapper source from template state

Dependencies: SOL-002, SOL-005.

Paths:

- `template/`
- `wrapper/`
- source `workspace.yaml` and `context/`
- `scripts/release-manifest.txt`
- `scripts/release-artifact.sh`
- `docs/release.md`
- release and smoke-test fixtures

Work:

- move uninitialized mutable seed state under `template/`;
- make source-repository Product Knowledge describe Context Circuit itself;
- assemble the published root workspace from static wrapper files plus the
  template seed;
- retain stable artifact entry paths required by supported hosts;
- verify that maintainer-only state and `.runtime/` never enter the artifact;
- document upgrade boundaries between wrapper version and workspace state.

Acceptance:

- source entry identifies the Context Circuit project rather than claiming it
  is uninitialized;
- a staged release starts as an uninitialized workspace;
- published root paths remain compatible;
- source-only, dirty, ignored, runtime, and credential material remain absent
  from the artifact;
- the staged template independently passes its smoke checks and context
  budgets.

Verification:

- release-artifact fixture in a temporary directory;
- source and staged-artifact acceptance suites;
- inventory comparison and credential scan;
- `git diff --check`.

Expected evidence: source/artifact identity comparison, artifact inventory,
smoke-test results.

Stop conditions: a host requires a path that the migration would remove, or
source/template separation would silently break upgrades for existing
workspaces.

### SOL-007 — Replace prose-presence tests with semantic contract tests

Dependencies: SOL-002, SOL-003, SOL-005.

Paths:

- `test/acceptance.sh`
- `test/contracts/`
- `test/routing/`
- `test/lifecycle/`
- `test/concurrency/`
- `test/release/`
- `test/context-budget/`

Work:

- retain behavioral fixtures for gates, leases, worktrees, recovery, archives,
  stacks, completion, and cleanup;
- replace assertions that require repeated phrases in multiple documents with
  invariant, schema, or route-result assertions;
- split the monolithic acceptance suite into focused scripts while retaining
  one top-level command;
- validate all catalog paths, invariant references, context sets, schemas, and
  budgets;
- fail when explanatory documentation becomes an accidental second authority.

Acceptance:

- one top-level acceptance command remains available;
- all existing behavior scenarios continue to pass;
- canonical rules may be rewritten once without requiring matching prose in
  several documents;
- a deliberately duplicated invariant, wrong route, oversized context set, or
  stale catalog path makes the suite fail.

Verification:

- negative mutation fixtures for each lint class;
- `sh test/acceptance.sh`;
- individual suite execution;
- `git diff --check`.

Expected evidence: test inventory mapping old assertions to new semantic
checks, positive and negative results.

Stop conditions: an existing prose assertion encodes behavior that has not yet
been represented in the canonical contract.

### SOL-008 — Run the A/B evaluation and prepare the migration handoff

Dependencies: SOL-004, SOL-005, SOL-006, SOL-007.

Paths:

- `docs/benchmarks/`
- release notes and migration documentation
- no product-code paths outside the approved redesign scope

Work:

- compare the frozen baseline with the redesigned source and staged template;
- execute all filesystem and natural-language route cases;
- measure root, writer, verifier, and full-plan context profiles;
- record safety regressions, ambiguous routes, tool-read counts, and context
  overruns;
- prepare an upgrade and rollback guide without publishing or merging.

Acceptance:

- one-plan static wrapper context is at most 39 KB;
- total wrapper-token reduction is at least 65%, with 75% as the target;
- filesystem routes pass 100%;
- natural-language route accuracy is at least 95%;
- no safety or lifecycle regression remains;
- staged artifact identity and source identity are both correct;
- remaining limitations and any missed stretch targets are explicit.

Verification:

- all top-level acceptance, routing, budget, source, and staged-artifact checks;
- `git diff --check`;
- independent read-only verifier reproduction of reported measurements.

Expected evidence: final benchmark table, verification handoff, migration and
rollback guide, unresolved limitation list.

Stop conditions: any safety regression, route accuracy below the minimum,
static reduction below 65%, or mismatch between the staged artifact and its
manifest.

## 7. Dependencies and execution sequence

```text
SOL-001 Baseline
   ↓
SOL-002 Canonical contract
   ↓
SOL-003 Router ───────────────┐
   ↓                         │
SOL-004 Host binding         │
                             ↓
SOL-005 Context compaction
   ├──────────────→ SOL-006 Source/template separation
   └──────────────→ SOL-007 Semantic tests
                         \       /
                          SOL-008 A/B verification
```

Recommended delivery batches after plan approval:

1. Baseline and additive canonical contract: SOL-001–SOL-002.
2. Router, host binding, and context compaction: SOL-003–SOL-005.
3. Source/template separation, test restructuring, and final benchmark:
   SOL-006–SOL-008.

Each batch must remain reviewable and preserve a passing compatibility path.
Batching is not authorization to commit, push, open a PR, merge, publish, or
deploy.

## 8. Repository scope

Expected repository key: `context-circuit`, subject to workspace initialization
or explicit confirmation before canonical execution.

Implementation scope:

- root agent instruction and workflow files;
- wrapper contracts and schemas;
- Context Circuit skills and role instructions;
- relevant explanatory and maintainer documentation;
- context-navigation and source/template identity files;
- acceptance, routing, budget, lifecycle, and release tests;
- release staging and manifest code.

Non-goals:

- weakening human approval, completion, cleanup, takeover, publication, merge,
  or deployment gates;
- removing exclusive plan leases or writable worktrees;
- removing independent verification;
- adding a scheduler, daemon, database, Node runtime, or package manager;
- adding a low-risk or small-change execution bypass in this plan;
- changing canonical plan/task lifecycle values;
- scanning or rewriting user files under `sources/`;
- modifying or replacing `Opus-4.8-plan.md`;
- publishing a release or delivering external changes.

## 9. Plan-level acceptance criteria

The redesign is acceptable only when:

1. the wrapper/source and instantiated-workspace boundary is explicit;
2. one canonical router produces structured, explainable decisions;
3. intent, phase, eligibility, capability, and authorization remain distinct;
4. all Section 4 context budgets pass in both source and staged template;
5. one-plan static wrapper tokens fall by at least 65%, with 75% targeted;
6. duplicated normative copies are replaced by one owner plus references;
7. at least 50 filesystem route fixtures pass 100%;
8. natural-language route accuracy reaches at least 95%;
9. session binding enables exact lookup without a global current pointer;
10. all existing safety and lifecycle behavior remains verified;
11. source and release artifact identities are no longer contradictory;
12. migration and rollback steps are documented and independently verified.

## 10. Risks and mitigations

| Risk | Mitigation |
| --- | --- |
| Compression drops a safety condition | Map every old rule to an invariant ID before deleting prose; test route coverage. |
| Machine-readable contracts become another duplicate authority | Declare them canonical and make prose reference IDs rather than restate rules. |
| Byte budgets encourage unreadable documents | Apply budgets to context sets, not arbitrary files; route deep detail conditionally. |
| Router chooses an eligible but unauthorized mutation | Keep `authorization` separate and test recommendation-only cases. |
| Source/template split breaks host entry paths | Preserve root adapters in the artifact and test every supported host path. |
| Host task identifiers are missing or unstable | Make binding optional and retain bounded fallback discovery. |
| Semantic test migration loses existing protections | Maintain an assertion ledger mapping every removed prose check to a contract or behavior check. |
| Natural-language benchmark overfits wording | Include paraphrases, incomplete requests, contradictions, and adversarially similar intents. |
| Refactor becomes an execution-topology redesign | Keep writer/verifier topology and gates explicitly out of scope. |
| Concurrent untracked planning files are overwritten | Treat them as user-owned and modify only this named plan. |

## 11. Assumptions

- The wrapper may add portable YAML contracts without adding a user-facing
  command runtime.
- Shell-based validation remains acceptable for the release artifact.
- Stable root adapter filenames are required for host discovery.
- Existing safety semantics are intentional unless contradicted by an explicit
  human decision.
- Byte and word budgets are stable acceptance measures; exact token counts are
  supplementary because model tokenizers differ.
- The source/template identity split is desirable, but its final filesystem
  layout remains subject to human review.

## 12. Open human decisions

1. Should this standalone draft be converted into a canonical
   `plans/context-circuit-plans/<number>-<slug>/` bundle, or remain a normal
   wrapper-repository implementation plan?
2. Confirm the hard byte budgets in Section 4, especially 8 KB bootstrap and
   39 KB one-plan static total.
3. Confirm `wrapper/` and `template/` as the physical source layout, or request
   a different naming convention.
4. Should migration be delivered in the recommended three review batches or
   one larger pull request?
5. Should the natural-language routing benchmark be required for local
   acceptance, release acceptance, or both?

## 13. Human gates and next action

This document remains a draft. Creating it does not approve or begin any task.

Required gates:

- plan-shape decision: canonical Context Circuit plan bundle versus ordinary
  wrapper-repository plan;
- plan approval before implementation;
- material scope review if safety topology or lifecycle semantics would change;
- independent verification before completion;
- explicit merge and publication authorization for any delivered release;
- explicit status-change confirmation before any canonical plan is marked
  done.

Recommended next action: perform a read-only plan review, compare this draft
with `Opus-4.8-plan.md`, resolve the five open decisions, and then either revise
this file or convert the accepted version into the canonical plan structure.
