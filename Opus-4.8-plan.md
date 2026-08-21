# PLAN — Context Circuit wrapper redesign

Scope of this plan: **the Context Circuit repository itself.** We are changing
the wrapper product, not operating inside it. There is no `plan.yaml`, no
`cc-approve-plan` gate, and no dogfooding — this is normal repo development on
`development/v0.5`, reviewed by diffs and shipped by PR to `main`.

Goals: make an AI agent load **less context**, make agent **routing explicit and
single-sourced**, and enforce a **consistent, single-owner document structure**.

Change scope: **whole repo, including tests, skills, agents, and release
scripts.** Nothing is off-limits; `CLAUDE.md`, `AGENTS.md`, `WORKFLOW.md`, and
`test/acceptance.sh` are all product surface we can restructure.

---

## 1. Diagnosis

One structural defect drives all three goals: **no single source of truth per
fact.** The same rules are restated across many files, so agents re-read
identical facts, routing logic drifts across copies, and two documents each
claim to be "normative."

Measured redundancy:

| Duplicated fact | Files repeating it |
| --- | --- |
| Session-entry reading order | 18 |
| "cc-run-stack is not a scheduler / not a second way to run one plan" | 7 |
| "mode is identity, not an execution topology" | 4 |
| "never infer completion from tests / Git state" | 4 |

Reading-path weight:

- "Read for every session" set: **~3,000 words**.
- It points into two more self-declared "normative" docs
  (`docs/agent-workspace-workflow.md` 3,270 w + `docs/runtime-contract.md`
  2,536 w), so a thorough entry pulls **~8,000+ words** before any project work.

`context/INDEX.md:4` already says docs "must not duplicate their facts." The
rule exists; nothing enforces it. This plan makes it enforceable.

---

## 2. Target design

### 2.1 Two-tier reading model

**Tier 0 — the spine (always loaded, budget: ≤ 1,500 words total).** Enough to
orient, identify session role, route, and know where everything else lives. No
restated deep contracts — only links.

- `AGENTS.md` — safety invariants, authority/precedence order, source boundary.
- `WORKFLOW.md` — entry sequence, role identification, the single routing table,
  the development-loop diagram, the human-gate list.
- `workspace.yaml` — machine config (repos, mode, branches, gates).
- `context/INDEX.md` — a pure navigation map: concept → owner file. No facts.
- `context/WORKSPACE.md` + `context/PROJECT.md` — identity, role unchanged.

**Tier 1 — deep contracts (loaded only when the route selects them).**

- `docs/runtime-contract.md` — record shapes, leases, handoff format. Execution/
  resume only (today it sits one hop off entry although most sessions never read
  a lease record).
- `docs/planning.md` — plan/task lifecycle, archive, stack semantics.
- `docs/agent-workspace-workflow.md` — collapsed so it is no longer a second
  "normative" doc competing with `WORKFLOW.md` (see 2.5).
- Each `.agents/skills/*/SKILL.md` and `agents/*.md` — loaded on invocation.

### 2.2 One concept, one owner

A new `docs/doc-ownership.md` maps every shared concept to exactly one owner
file plus the canonical phrase the lint watches. Every other file may reference
a concept by link but must not restate the rule.

| Concept | Single owner |
| --- | --- |
| Safety invariants, authority/precedence order, source boundary | `AGENTS.md` |
| Entry sequence, role identification, routing table, dev loop, gate list | `WORKFLOW.md` |
| Config schema (repos, mode, branches, gates) | `workspace.yaml` + `docs/configuration.md` |
| Concept → owner navigation map | `context/INDEX.md` (links to the manifest) |
| Runtime record shapes, leases, handoff format | `docs/runtime-contract.md` |
| Plan/task lifecycle, archive, stack semantics | `docs/planning.md` |
| Product Knowledge structure | `docs/product-knowledge.md` |
| Delivery / integrations / hosts | `docs/delivery-policies.md`, `docs/integrations.md`, `docs/host-capabilities.md` |
| Each skill's procedure | its `SKILL.md` |
| Agent execution roles | `agents/*.md` |

Exact per-phrase assignments come from the audit (Phase A) before any normative
file is edited.

### 2.3 One canonical routing table

The `Situation → Route` table in `docs/agent-workspace-workflow.md` is the good
part; it is duplicated three more times (a prose version right beneath it, a
version in `WORKFLOW.md`, and conceptually inside `cc-whats-next`). Move the
canonical table into `WORKFLOW.md` (Tier 0), delete the restatements, and have
`cc-whats-next` and the contract hub link to it.

### 2.4 Defensive-prose purge

A large share of words are defensive negations ("is not a scheduler," "not a
second way to run one plan," "do not present X as Y") — docs arguing with a past
misread. State each concept once, positively, in its owner; delete the
negations. Keep any negation that encodes a real safety constraint, relocated to
the safety owner.

### 2.5 Collapse the two "normatives"

`WORKFLOW.md` calls itself "the normative entry summary"; `agent-workspace-
workflow.md` calls itself "the normative behavior contract." Two normatives =
drift. `WORKFLOW.md` becomes the single Tier-0 behavior authority;
`agent-workspace-workflow.md` is reduced to a Tier-1 hub holding only content
owned nowhere else (acceptance scenarios, detailed subagent topology) plus
links — or removed entirely if the audit shows all its content has owners.

### 2.6 Guardrail: `test/doc-lint.sh`, wired into acceptance

So redundancy can't creep back, a new lint (run by `sh test/acceptance.sh`):

1. Reads canonical phrases from `docs/doc-ownership.md`; fails if any appears
   outside its owner file.
2. Fails if a routing-table marker appears in more than one file.
3. Fails if the Tier-0 spine exceeds the word budget.

This is the machine-enforced version of the rule `context/INDEX.md:4` states in
prose today.

---

## 3. Constraints

- **Behavior semantics stay intact.** Authority order, human gates, lease/
  worktree ownership, and plan/task lifecycle values (`draft`/`approved`/`done`,
  `draft`/`ready`/`done`) do not change meaning. This is a structural/compression
  redesign. `test/acceptance.sh` stays green — its nine behavior scenarios still
  pass. (Tests themselves may be refactored, but the behavior they assert holds.)
- **`test/acceptance.sh` is a co-deliverable.** It references docs by path
  heavily (README ×35, `agent-workspace-workflow.md` ×15, `WORKFLOW.md`/
  `INDEX.md` ×15, `AGENTS.md` ×13) and pins normative prose in only ~8 spots
  (mostly fixtures). Every moved/renamed doc has its assertions updated in the
  same change.
- **Host wiring stays valid.** `CLAUDE.md`, `.cursor/`, `.github/`, `agents/`
  reference doc paths; keep spine filenames (`AGENTS.md`, `WORKFLOW.md`,
  `workspace.yaml`) stable and re-verify no adapter points at a removed path.

---

## 4. Phases

Normal branch work; each phase is a reviewable batch of diffs.

### Phase A — Audit (no normative edits)
- Inventory every duplicated fact → write `docs/doc-ownership.md`
  (concept → owner + canonical phrase).
- Inventory every `test/acceptance.sh` doc-path/prose assertion and every
  adapter reference (`CLAUDE.md`, `.cursor/`, `.github/`, `agents/`); mark each
  keep / update / remove.
- **Done when:** ownership manifest covers every concept; assertion ledger is
  complete.

### Phase B — Tier-0 spine rewrite
- Fold the canonical routing table into `WORKFLOW.md`; make `INDEX.md` a pure
  concept → owner map; strip restated facts to links; spine ≤ 1,500 words.
- **Done when:** `wc -w` spine ≤ 1,500; exactly one routing table.

### Phase C — Tier-1 contract consolidation
- Collapse `agent-workspace-workflow.md` per 2.5; trim `runtime-contract.md`,
  `planning.md`, `ARCHITECTURE.md`, `CONVENTIONS.md`, remaining `docs/*.md`;
  remove defensive negations; one fact, one home.
- **Done when:** no canonical phrase appears outside its owner.

### Phase D — Skill + agent dedup
- Each `.agents/skills/*/SKILL.md` and `agents/*.md` references owners instead of
  restating lifecycle/routing/gate rules; `cc-whats-next` points to the routing
  table.
- **Done when:** doc-lint passes across skills and agents; each keeps its unique
  procedure.

### Phase E — Guardrail + acceptance update
- Implement `test/doc-lint.sh` (the three checks); wire into
  `test/acceptance.sh`; update doc-path assertions in lockstep with B/C.
- **Done when:** `sh test/acceptance.sh` passes including the new lint.

### Phase F — Verify + measure
- Full suite green; record before/after Tier-0 word/token counts; grep adapters
  for stale paths.
- **Done when:** every acceptance criterion below is met with recorded numbers.

---

## 5. Acceptance criteria

- Tier-0 spine ≤ 1,500 words (from ~3,000), by `wc -w`.
- Exactly one routing decision table in the repo.
- No canonical phrase (per `docs/doc-ownership.md`) appears outside its owner;
  `test/doc-lint.sh` enforces and passes.
- `docs/doc-ownership.md` exists and is linked from `context/INDEX.md`.
- `sh test/acceptance.sh` passes, including all nine behavior scenarios.
- No adapter references a removed or renamed doc path.

---

## 6. Risks

- **R1 — acceptance.sh pins doc paths.** Moving a doc without updating
  assertions breaks the suite. → Phase A inventories; Phase E updates in lockstep.
- **R2 — redesign drifts into behavior change.** → Semantics are a non-goal; any
  acceptance failure for a non-path reason halts as a regression.
- **R3 — adapters reference moved paths.** → Phase A greps all adapters; Phase F
  re-verifies.
- **R4 — doc-lint false positives.** → Lint matches only phrases registered in
  the manifest, chosen distinctive enough to avoid collisions.
- **R5 — 1,500-word budget too tight.** → Budget enforced only after Tier-1
  owners exist; Phase B surfaces an unowned fact rather than dropping it.

---

## 7. Open decisions

1. Confirm the ≤ 1,500-word Tier-0 budget, or set another target.
2. `docs/agent-workspace-workflow.md`: reduce to a Tier-1 hub, or delete and
   redistribute? (Deletion is cleaner but touches 15 acceptance references.)
3. Ship as one PR to `main`, or stage the phases across several PRs?
