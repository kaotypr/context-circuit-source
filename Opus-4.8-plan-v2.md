# Opus 4.8 Plan v2 — Context Circuit redesign (measured, staged, cost-gated)

Status: **draft — human review and approval required before execution.**

Scope: the Context Circuit wrapper source repository and the workspace template
it publishes — instruction, routing, docs, tests, and release surfaces. This
plan authorizes no merge, publication, deployment, or lifecycle status change.

Goals, unchanged: agents load **less context**, follow **one explainable
route**, and read a document set where **every rule has exactly one owner** —
with **zero change to safety semantics.**

---

## 0. What makes this plan different

Three failure modes kill redesigns like this. Each goal below is chosen to
prevent one, and these principles bind every task in Section 4.

1. **Abstraction must pay for itself (anti-over-engineering).** The core risk of
   a redesign is curing doc-bloat by adding a heavier substrate that agents must
   *also* read. So no new machine-readable layer ships until a **spike proves it
   nets fewer tokens on a real route** than well-owned prose. We reach for the
   lightest mechanism that passes the test, and promote to heavier machinery only
   on measured evidence — never on aesthetics.

2. **Behavior is frozen by a differential, not by trust.** Before any edit we
   capture the *current* routing decision and safety outcome for every fixture as
   a **golden record**, then assert byte-identical decisions after each change.
   "Semantics unchanged" becomes a test that fails loudly, not a promise.

3. **Value lands first, risk lands last.** Phases are ordered so the **largest
   token saving ships in Phase 1** with no new substrate, and every phase is
   independently reviewable and revertable. Nothing valuable is gated behind the
   most speculative work.

The measure of success is not "we built the contract system." It is "an agent
does the same job correctly for far fewer tokens, and we can prove both."

---

## 1. Measured baseline (this checkout)

| Context profile | Words | Bytes |
| --- | ---: | ---: |
| Always-read spine | 3,068 | 22.3 KB |
| Deep contracts it points to | 5,806 | 42.2 KB |
| Thorough root entry (spine + deep) | 8,874 | 64.5 KB |
| All 15 skill bodies | 7,252 | 51.1 KB |
| Agent role files | 1,339 | 9.8 KB |

Duplication (owner-count): the **writer-child rule spans 12 files**;
"not a scheduler" spans 5; the reading-order, gate, and lifecycle rules each
repeat across 4–18 files. `context/INDEX.md:1-4` already forbids this
duplication in prose; nothing enforces it.

Structural defect: this repo is simultaneously the **product** and, per
`workspace.yaml` + `context/PROJECT.md`, an **"uninitialized workspace."** An
agent cannot tell whether it is developing the wrapper or operating an instance
of it. Neither goal (tokens, routing) is fully solvable while that ambiguity
stands.

Measurement method: byte and word counts are the stable acceptance unit; an
approximate model-token count is recorded alongside when a tokenizer is
available. Budgets are asserted in bytes so results are reproducible across
tokenizers.

---

## 2. Target model

### 2.1 Two-tier reading with named route context sets

- **Tier 0 spine — always read, hard budget 8 KB.** `AGENTS.md` (safety,
  precedence, source boundary), `WORKFLOW.md` (lifecycle summary + link to the
  one router), `workspace.yaml`, `context/INDEX.md` (navigation only, zero
  facts), `context/WORKSPACE.md` + `context/PROJECT.md` (identity).
- **Tier 1 — loaded only when a route selects it:** deep contracts, skills, role
  deltas. Each route names a **context set**: the exact files/fields to read, a
  byte budget, and the reason codes that unlock any conditional deep read.

Runtime guidance is split by concern so ordinary orientation never loads lease
acquisition, stack joins, takeover, and cleanup together.

### 2.2 One routing decision, with authority separated

A single canonical route decision, evaluated by `cc-session-entry` and merely
*viewed* by `cc-whats-next` (no second model). It keeps four things distinct
that the current docs blend:

- **intent** — what the human asked for;
- **eligibility** — whether state permits the candidate action;
- **authorization** — whether the agent may *mutate* (recommendation ≠
  permission);
- **phase** — current session lifecycle.

The decision emits **reason codes** so every route is explainable and testable.
Separating eligibility from authorization is a safety improvement, not just
tidiness: it is exactly the confusion that lets agents cross gates.

**Representation is decided empirically (Principle 1).** It starts as one
canonical table plus a machine-checkable fixture set. It is promoted to a
machine-readable `routing.yaml` **only if** the Phase 6 spike shows the YAML form
nets fewer tokens across routes than the owned-prose form. We do not assume the
answer.

### 2.3 One owner per rule, enforced semantically

Each shared rule gets a stable **invariant ID** in a single flat registry
(`docs/invariants.md` — one line per invariant, its owner, and canonical text).
Every other document *references the ID*; none restates the rule. This is a flat
registry, not a schema tree — the lightest thing that lets tests assert the
*rule*, not its *wording*. Tests then stop pinning prose (the reason prose was
copied everywhere), which removes the incentive to duplicate.

### 2.4 Product vs. template identity

Fix the ambiguity in the smallest sufficient step, spike-gated:

- **First, surgically:** make source-repo Product Knowledge describe Context
  Circuit itself, and add a manifest declaring which paths ship in the released
  template vs. which are maintainer-only. This alone removes the "am I the
  product or an instance?" ambiguity for agents.
- **Structurally (only if the spike shows the manifest is insufficient):**
  physically separate a `template/` seed from shipped wrapper files. The
  directory move is the expensive option and must justify itself.

---

## 3. Quantitative targets

| Profile | Current bytes | Hard budget | Reduction |
| --- | ---: | ---: | ---: |
| Tier-0 spine (always read) | 22.3 KB | **≤ 8 KB** | ~64% |
| Thorough root entry | 64.5 KB | **≤ 20 KB** | ~69% |
| Any single route's context set | — | **≤ its declared budget** | — |
| Each substantial skill body | up to ~6 KB | **≤ 700 words, prefer ≤ 450** | — |

Routing & safety targets:

- exactly **one** canonical route decision; `cc-whats-next` duplicates none of it;
- **≥ 50 filesystem route fixtures**, resolving **100% identical to the golden
  baseline** (behavior freeze), covering fresh/child/resume, draft/approved/done/
  archived, malformed, blocked-dependency, live/stale lease, completion,
  cleanup, single-plan, and stack;
- the writer-child rule drops from **12 owners to 1**; every registered invariant
  has exactly one owner;
- **no route infers authorization from eligibility** (explicit
  recommendation-only fixtures);
- a natural-language route sample is a **reported metric, not a launch gate** —
  its cost/benefit is a Phase 6 decision, not an assumed requirement.

Every budget applies to **both** the source repo and the staged template.

---

## 4. Phases (value-first, each independently reviewable & revertable)

Each phase states its gate. A phase that fails its gate stops for a human
decision rather than proceeding.

### P0 — Freeze the baseline (no normative edits)
Record byte/word/token profiles (Section 1). Build the golden record: current
route decision + safety outcome for ≥ 50 filesystem fixtures. Add
`test/context-budget.sh` (byte budgets) and a differential runner that will fail
if any fixture's decision changes.
**Gate:** every profile reproduces from explicit paths; every fixture has a
recorded golden decision; `sh test/acceptance.sh` green.

### P1 — Spine compaction + single-owner registry (ships the biggest win now)
Create `docs/invariants.md` (IDs + owners). Rewrite the Tier-0 spine to ≤ 8 KB:
fold the canonical routing table into `WORKFLOW.md`, make `INDEX.md`
navigation-only, replace restated rules with invariant references. No new
substrate — pure prose consolidation.
**Gate:** spine ≤ 8 KB; writer-child + reading-order rules each have one owner;
golden fixtures still 100% identical; acceptance green.

### P2 — Canonical route decision + eligibility/authorization split
Make `cc-session-entry` the sole evaluator; `cc-whats-next` a read-only view.
Introduce the intent/eligibility/authorization/phase/reason-code vocabulary as a
table + fixtures (no YAML yet).
**Gate:** one routing table; all fixtures resolve **identically to golden**;
recommendation-only cases never authorize mutation; malformed state returns a
visible recovery decision.

### P3 — Route context sets + CI budgets
Define named context sets per route with byte budgets and conditional-deep-read
reason codes. Split runtime guidance by concern.
**Gate:** every route meets its budget in source and staged template; a route
exceeding budget must record a reason or fail.

### P4 — Semantic tests replace prose-presence tests
Migrate acceptance assertions from "phrase appears in N docs" to invariant/route/
schema assertions. Keep all behavioral fixtures. Maintain an **assertion ledger**
mapping every removed prose check to its replacement, so no protection is lost.
**Gate:** one top-level `sh test/acceptance.sh`; a deliberately duplicated
invariant, wrong route, or oversized context set makes the suite fail; all prior
behavior scenarios pass.

### P5 — Product vs. template identity (surgical)
Fix source-repo identity; add the ship-vs-maintainer manifest; assert the staged
artifact starts uninitialized and excludes maintainer-only, `.runtime/`, dirty,
and credential material.
**Gate:** source entry identifies the project; staged template passes its own
budgets and smoke checks; published root paths stay host-compatible.

### P6 — Contract-layer spike + go/no-go (the cost gate)
Spike: express **one** route (`cc-run-plan`) as a machine-readable contract and
measure end-to-end tokens vs. the owned-prose form from P1–P3. Same spike tests
whether the identity manifest (P5) suffices or a `template/` move is warranted.
**Gate (explicit human decision):** promote the contract layer / directory split
**only if** it measurably nets fewer tokens or removes a real safety ambiguity;
otherwise keep the lighter form and record why. This is where this plan refuses
to over-build.

### P7 — Final A/B, migration & rollback
Independent read-only verifier reproduces all measurements; produce upgrade +
rollback guide. No merge or publish.
**Gate:** targets met in source and template; zero safety/lifecycle regression;
artifact matches its manifest.

Delivery batches: (P0–P1) win + freeze, (P2–P4) router + tests, (P5–P7)
identity + cost gate + verification. Batching authorizes no commit, push, PR,
merge, or publish.

---

## 5. Acceptance criteria

1. Tier-0 spine ≤ 8 KB; thorough root entry ≤ 20 KB; both in source and template.
2. Exactly one canonical route decision; `cc-whats-next` duplicates none of it.
3. intent / eligibility / authorization / phase are distinct; routes carry reason
   codes; no route infers authorization from eligibility.
4. ≥ 50 filesystem fixtures resolve **100% identical to the P0 golden baseline**
   (proven behavior freeze).
5. Every registered invariant has exactly one owner; writer-child rule = 1 owner;
   tests assert invariants, not prose.
6. Product-vs-template identity is unambiguous; staged artifact is clean and
   host-compatible.
7. Any heavier machinery (routing YAML, `template/` split) ships **only** with
   recorded evidence it nets fewer tokens or removes a safety ambiguity.
8. One top-level acceptance command; assertion ledger maps every removed prose
   check to a semantic replacement.

---

## 6. Risks & mitigations

| Risk | Mitigation |
| --- | --- |
| Cure adds a substrate that eats the savings | P6 spike gate: no heavy contract layer without measured net reduction |
| Compression silently drops a safety rule | P0 golden differential + invariant-ID mapping before any prose deletion |
| Redesign drifts into behavior change | Any non-path acceptance failure halts as a regression; fixtures assert against golden |
| Router recommends an eligible-but-unauthorized mutation | authorization is a separate field with recommendation-only fixtures |
| Invariant IDs become a second authority to read | Flat one-line registry; prose references IDs; budgets counted with IDs included |
| Source/template split breaks host paths | Manifest-first; directory move only if spike-justified; test every host entry path |
| Value stranded behind speculative work | Biggest token win is P1, before any new substrate |
| NL-routing benchmark overfits / over-costs | Reported metric, not a launch gate; deterministic fixtures are the bar |

---

## 7. Non-goals

Weakening any human gate; removing exclusive leases, writable worktrees, or
independent verification; adding a scheduler, daemon, DB, Node runtime, or
package manager; changing canonical plan/task lifecycle values; scanning or
rewriting `sources/`; renaming `AGENTS.md` / `WORKFLOW.md` / `workspace.yaml`
(host-wired); publishing or merging.

---

## 8. Open decisions

1. Confirm the 8 KB spine / 20 KB root-entry budgets.
2. Confirm invariant IDs live in one flat `docs/invariants.md` (vs. a richer
   schema) as the default, with promotion gated by P6.
3. Confirm identity is fixed manifest-first, `template/` move only if P6 warrants.
4. Confirm the NL-routing sample is reported-only, not a launch gate.
5. Deliver in the three batches above, or one PR?

---

## 9. Why this beats a contract-first redesign

A contract-first plan front-loads its most speculative, highest-churn work (a new
YAML substrate + a directory migration + an NL benchmark) and only measures the
payoff at the end. This plan inverts that: it **measures first, ships the certain
win first, and makes every heavy abstraction pass a token-cost gate before it
exists.** It reaches the same target architecture where that architecture is
justified — and stops short of it where prose with one owner already wins. Less
machinery for the same correctness is the whole point of the project; the plan
should embody it.
