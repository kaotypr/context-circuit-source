# 002 — CC-010 execute agent trace

This is a passive source artifact. It records what agents actually did while
executing approved plan `verification-evidence-layers` (CC-010) in one
product-source Cursor session on 2026-08-23. It is not accepted Product
Knowledge, not a canonical workspace plan, and not authorization to change
`wrapper/`, `template/`, skills, or runtime.

Drafted 2026-08-23 from named session evidence only. Do not scan the rest of
`sources/`.

The point of this note is the **execute ceremony**: which steps were about
lifecycle rules (route, packet, lease, child, verifier, isolation) versus
steps that implemented evidence-layer contracts. Times below are wall clock
in UTC+7 unless noted.

## Provenance

| ID | Path or locator | Why it was read |
| --- | --- | --- |
| CC-010 | `plans/context-circuit-plans/verification-evidence-layers/plan.yaml`, `PLAN.md`, `tasks/001-contracts.md`, `002-verifier-enforcement.md`, `003-verification.md` | The approved plan that was executed |
| HOST | Cursor Agent session `b651ce73-e628-446d-bafd-e0ef535a2874`, 2026-08-23 | Root coordinator turns and tool calls |
| WRITER | Cursor Task child `6511de0b-f5ef-4b3c-9651-16eb0786189f` | One writer reused for VEL-001, VEL-002, VEL-003 |
| VERIFIER | Cursor Task child `78dc5ee4-a70a-4b06-be79-ddcf42ded586` | Independent verifier; `write_worktree: false` |
| RT-ROOT | `.runtime/sessions/cc010-root-20260823/` | Route dump, root `run-plan` receipt |
| RT-WRITER | `.runtime/sessions/cc010-writer-vel001-20260823/` | Writer receipts, delegations, handoffs |
| RT-VERIFIER | `.runtime/sessions/cc010-verifier-20260823/` | Verifier receipt, delegation, handoff |
| TX | `.runtime/.transactions/cc010-tx-vel00{1,2,3,verify}-20260823/` | Engine-owned runtime graphs and commit markers |
| LEASE | `.runtime/plans/verification-evidence-layers/lease.yaml` | Exclusive write lease; still `active` at end |
| WT | `.runtime/worktrees/context-circuit-source/verification-evidence-layers` | Detached worktree at `655fcf5` |
| SKILLS | `.agents/skills/cc-entry/SKILL.md`, `cc-execute/SKILL.md`, `cc-verify/SKILL.md` | Rules the root loaded before touching the plan |
| ENGINE | `wrapper/runtime/engine.sh` | `cc_route`, packet load, lease, worktree, graph |

Related but not copied here: `001-session-agent-behavior-review.md` (CC-003
lifecycle walk) and the frozen byte baseline in
`test/baselines/profile-baseline.yaml`.

## Clock

| Mark | UTC+7 | Evidence |
| --- | --- | --- |
| Human: execute CC-010, and watch CC-006/008/009 | 18:16 | User query timestamp |
| First runtime files (route, lease, worktree) | 18:27:57 | `route.yaml`, `lease.yaml`, worktree mtime |
| VEL-001 graph committed | 18:28:10 | `cc010-tx-vel001-20260823/commit.marker` |
| VEL-001 writer handoff | 18:36:14 | `handoff.md` |
| VEL-002 graph committed | 18:39:07 | `cc010-tx-vel002-20260823/commit.marker` |
| VEL-002 writer handoff | 18:46:48 | `handoff-vel002.md` |
| VEL-003 graph committed | 18:49:40 | `cc010-tx-vel003-20260823/commit.marker` |
| VEL-003 writer handoff | 18:53:37 | `handoff-vel003.md` |
| Verifier graph committed | 18:58:23 | `cc010-tx-verify-20260823/commit.marker` |
| Verifier handoff | 19:09:01 | verifier `handoff.md` |
| Root execute report in chat | ~19:18 | Next human query at 19:18 |

**53 minutes** is 18:16 → 19:09 (execute request to independent verifier
handoff). About **9 more minutes** of root wrap followed. Full request-to-report
wall clock is about **62 minutes**.

This was not a one-task docs plan. CC-010 has three sequential tasks plus an
independent verifier. The coordinator also spent time checking whether
CC-006/008/009 were live, which is extra relative to a bare execute.

## Actors and tool counts

| Actor | Tool calls (approx.) | Role |
| --- | ---: | --- |
| Root coordinator | 116 in the whole chat; most during execute | Route, packets, lease, graphs, launch children, observe prior plans |
| Writer child | 185 | Implement VEL-001–003 in the worktree only |
| Verifier child | 89 | Re-run canonical commands; write only its handoff |

Writer writes/patches: 66. Verifier writes: 1 (its handoff). Main checkout
porcelain stayed empty.

## Process vs product

Two kinds of work are tagged below:

- **Process / rules** — reading skills, routing, budgets, leases, graphs,
  isolation checks, “do not self-verify”, “do not finish”. This is thinking
  about how Context Circuit says an execute must run, not about evidence
  layers.
- **Product** — changing schemas, engine comparison, fixtures, and tests that
  CC-010 asked for.

Rough wall split of the 53 minutes to verifier handoff:

| Kind | Approx. wall | Where |
| --- | ---: | --- |
| Process / rules | ~32 min | 18:16–18:28 preflight; 18:36–18:39 and 18:46–18:50 between tasks; 18:53–18:58 verifier setup |
| Product (writer) | ~19 min | 18:28–18:36 VEL-001; 18:39–18:46 VEL-002; 18:49–18:53 VEL-003 |
| Product (verifier tests) | ~11 min | 18:58–19:09, mixed with read-only rules |

The verifier’s 11 minutes are the plan’s verification work *and* process
(graph validation, refuse repair). They are not implementation.

## Step-by-step

### 0. Before execute (18:08–18:16) — not CC-010

The human asked this session to read conversation `25fa9f9a-…`. The root
summarized CC-006/008/009 and the leftover CC-010 draft/approved state. That
orient pass is **outside** the 53-minute execute clock.

### 1. Human execute request (18:16)

Exact intent: execute CC-010, and observe whether already-done CC-006, CC-008,
and CC-009 take effect. That second clause caused extra packet measurements
and status checks that a minimal execute would skip.

### 2. Root loads execute rules (18:16–~18:22) — process

The root did not start a worktree. It read, in order:

1. `.agents/skills/cc-entry/SKILL.md`
2. `.agents/skills/cc-execute/SKILL.md`
3. `.agents/skills/cc-verify/SKILL.md`
4. Adapter `AGENTS.md` / `WORKFLOW.md`, `wrapper/contracts/routes.yaml`,
   `wrapper/contracts/context-sets.yaml`, `wrapper/contracts/tier0.yaml`
5. Large slices of `wrapper/runtime/engine.sh` (`cc_route`, packet loader,
   lease, worktree, `cc_construct_runtime_graph`)

This is **rule reconstruction**. The coordinator re-learned how to execute
instead of trusting a short launch packet. Same pattern as CC-003’s “~150
root tool calls to re-learn routing,” now on execute rather than approve.

### 3. Route and packet math (18:22–18:27) — process, and CC-008 live

The root ran `cc_route` on the human text. Result: `execute-plan` /
`run-plan` / `ROUTE_SELECTED` / `explicitly-requested`.

Then it measured packets instead of loading everything:

| Attempt | Result |
| --- | --- |
| Root `run-plan` + `plan.yaml` + VEL-001 only | **20,657 / 22,528** — fit |
| Writer + static files + `plan.yaml` | **15,859 > 12,288** — `CONTEXT_BUDGET_EXCEEDED` |
| Writer + `plan.yaml` + VEL-002 + VEL-003 | **13,597 > 12,288** — rejected |
| Writer + delegation schema + `plan.yaml` + one task | **11,605–11,748 / 12,288** — fit |

CC-008 took effect here: the engine refused undeclared/over-budget sets. The
coordinator spent several minutes finding a legal subset. That is process
time caused by a real product constraint.

Other preflight:

- Identity `kind: product-source` → CC-006 projection **does not apply**
- `git status` clean; HEAD `655fcf5` (`Approve CC-010 evidence-layer plan`)
- CC-009 `plan.yaml` still `approved` while implementation is on the branch
  (lifecycle lag). Root proceeded on human attestation, without marking
  CC-009 `done`.
- No `repositories.local.yaml`; product-source used `cc_prepare_worktree`
  on this repo (detached HEAD). Named-branch worktrees still not used.

### 4. Lease, worktree, VEL-001 graph (18:27:57–18:28:44) — process / CC-009

In ~50 seconds the engine wrote:

- exclusive lease (`session_id: cc010-writer-vel001-20260823`,
  `acquired_at` still hardcoded `2026-08-21T00:00:00Z`)
- detached worktree at `655fcf5`
- root receipt `run-plan` 20,657 bytes
- writer receipt 11,748 bytes
- runtime graph `cc010-tx-vel001-20260823` → `runtime-graph-authoritative`
- launch projection: `role: writer`, `assigned_root` = worktree,
  locators `delegation.yaml` / `handoff.md`

Root copied graph records into the writer session dir so the child could
read them. It did not edit product files.

### 5. Writer VEL-001 (18:28:44–18:36:14) — product, ~7.5 min

Child `6511de0b-…` resumed later for 002/003 (same worktree, same lease).

Recorded steps:

1. **Process:** “Reading writer session records” (delegation, receipt, launch)
2. **Product:** “Inspecting current contract schemas”
3. **Product:** Write/patch `plan.yaml`, `task.yaml`, `delegation.yaml`,
   `handoff.yaml`, `completion.yaml`, `invariants.yaml`, `manifest.yaml`,
   fixtures under `test/contracts/fixtures/evidence-layers/`
4. **Product:** `sh test/contracts/test-contracts.sh` (writer claim, PASS)
5. **Process:** `git status` in the worktree; write `handoff.md`; do not
   change `plan.yaml` status

Isolation check after this child (root, ~18:36–18:38): main checkout empty
porcelain; only the worktree dirty. That check is process, requested by
“see if done plans take effect.”

### 6. VEL-002 graph (18:38:59–18:39:10) — process

New transaction `cc010-tx-vel002-20260823` on the **same lease and session
id** (lease session must match graph session). New writer packet 11,605
bytes. ~10 seconds of engine, plus ~3 minutes of root wrapping after
VEL-001.

### 7. Writer VEL-002 (18:39:10–18:46:48) — product, ~7.5 min

Same child, resumed:

1. Process: read VEL-002 records
2. Product: exact-match layer comparison in `wrapper/runtime/engine.sh`
3. Product: thin pointers in `cc-verify` skill and `agents/verifier.md`
   (explicitly not a second policy owner — that sentence is a **rule**
   inside product work)
4. Product: runtime fixtures + `sh test/runtime/test-runtime.sh`
5. Process: `handoff-vel002.md`

### 8. VEL-003 graph (18:49:35–18:49:43) — process

Same pattern. Writer packet 11,711 bytes. Gap 18:46:48–18:49:35 is root
measuring the next packet and constructing the graph (~3 min process).

### 9. Writer VEL-003 (18:49:43–18:53:37) — product, ~4 min

1. Process: read VEL-003 records
2. Product: host limitation fixtures, security negatives, upgrade copy
   behavior, release exclusions, `test/acceptance/evidence-layer-mapping.yaml`
3. Product: host/security/upgrade/release scripts PASS; `sh test/acceptance.sh`
   **FAIL** `8466 > 8192` on the tier-0 ledger (`wrapper/manifest.yaml`
   grew). Writer did not repair that; it is outside VEL-003’s delegated
   “do not repair unrelated failures” stop — mixed product + rule.
4. Process: `handoff-vel003.md`

### 10. Verifier graph (18:53:37–18:58:26) — process

Root refused to repair in-session. It:

- measured a verifier packet; a writer handoff path with a non-canonical
  filename was `CONTEXT_PATH_UNDECLARED` (CC-008 again)
- loaded 10,581 / 11,264 bytes (`verifier.md` + schemas + `plan.yaml`)
- could not `cc_acquire_lease` a second lock; built a **read-only lease
  input file** for graph construction without taking `lease.lock`
- committed `cc010-tx-verify-20260823` with `write_worktree: false`

~5 minutes, no product edits.

### 11. Verifier child (18:58:26–19:09:01) — product verification + rules, ~11 min

Child `78dc5ee4-…`:

1. **Process:** “Reading verifier skill”; validate graph/receipts with
   worktree `engine.sh`
2. **Product:** rerun `test-contracts`, `test-runtime`, host, security,
   upgrade, release, `test/acceptance.sh`
3. **Process:** treat writer handoffs as claims; do not repair; confirm
   main checkout still clean
4. **Product result:** VEL-AC-01–06 pass with exact layer match; VEL-AC-07
   fail, same `8466 > 8192`
5. **Process:** write only
   `.runtime/sessions/cc010-verifier-20260823/handoff.md`

### 12. Root report (~19:09–19:18) — process

Root did not finish the plan, did not merge, did not repair. It told the
human: next safe action is bounded repair of the tier-0 overrun, then
another independent verifier.

## What the agents were thinking about (rules, not the problem)

These thoughts consumed clock and tokens without moving evidence-layer
semantics:

1. Which skill applies (`cc-entry` then `cc-execute`); never infer execute
   from a generic resume.
2. Load only the route-selected packet; stop on undeclared path or overrun.
3. Product-source dirty-base vs `MAINTAINER_APPROVAL_COMMIT_REQUIRED`
   (not applicable; approval already committed).
4. CC-009 `approved` vs human “it’s done” — do not silently `done` it.
5. `cc_safe_id` rejects slashes in `codex/development/v0.5`; use detached
   `HEAD` instead of a named branch.
6. One writer lease; sequential tasks must reuse `session_id`.
7. Children get launch projection only; root must not implement in the
   main checkout.
8. Verifier must not self-verify in root; missing child would be
   `host-blocked`.
9. After fail, do not finish; do not auto-repair.

The evidence-layer *problem* (vocabulary, exact match, waived ≠ pass,
legacy unread-mutation) only appears in writer steps 5/7/9 and verifier
step 11.

## Direct-edit counterfactual

A single agent that created a worktree and edited files would skip steps
2–4, 6, 8, 10, and most of 12: on this trace that is on the order of
**half an hour of process**. Implementation + tests (steps 5, 7, 9, 11)
were about **20–30 minutes**. Cowboy would still need those, minus
independent replay, so perhaps ~15–25 minutes of product work.

That faster path would not have proven CC-008 budgets, would likely have
dirtied the main checkout, and would have let one agent both implement
and declare VEL-AC-07 fixed.

## Outcome at end of execute

- Plan status: still `approved` (correct; finish is a later gate).
- Tasks: still `ready`.
- Lease: still `active`.
- Main checkout: clean.
- Worktree: dirty with VEL-001–003; detached `655fcf5`.
- Independent verification: **fail** VEL-AC-07 / VEL-VT-07.
- CC-006: projection skipped for `product-source` (as designed).
- CC-008: live (registered set, overruns and undeclared paths stopped).
- CC-009: live graphs/launch; lease timestamps still fake; worktree still
  detached.

## Why this execute was slow

1. Root re-read execute machinery instead of a short already-routed packet.
2. Observing CC-006/008/009 added measurements and isolation tours.
3. Three tasks got **three graphs and three launches** instead of one
   writer delegation covering VEL-001–003.
4. Packet budgets forced combinatorial measurement.
5. Independent verifier is a full second agent (~11 min) after the writer
   already ran the same tests.

None of those minutes are “thinking about whether schema and browser are
the same layer.” That comparison is cheap once the ceremony is done.
