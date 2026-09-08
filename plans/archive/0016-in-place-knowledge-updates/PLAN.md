# 0016 — Update Product Knowledge in place when a plan is marked done

Plan ID: `0016-in-place-knowledge-updates`  
Intent: `i014-direct-context-updates`  
Status: `done`

## Original request and coverage

Product Knowledge is updated by changing the live `context/` files, not by staging a proposal. Gathering context already does that in place. A plan becomes done only when you ask to mark it done; that ask then updates live context files if the plan affected Product Knowledge. Delivery stays delivery only. The next plan can start without waiting. Pending proposal `0034` is dropped.

## Objective and desired behavior

- Gathering context and reconciling Product Knowledge write live `context/` files and keep the retrieval catalog consistent. They do not create `context/proposals/` documents.
- There is no separate “accept the context update” step.
- A plan becomes done only on an explicit mark-done ask. Delivery, verification, and candidate acceptance do not mark it done. Explore stays planless.
- Mark-done then updates live context files when the plan affected Product Knowledge; if it did not, context files stay as they are.
- Delivery does not start a knowledge update and does not mark the plan done.
- A later plan can start even if Product Knowledge from a prior plan has not been updated yet.
- Pending proposal `0034` is dropped, not applied as a domain page.

## Constraints and non-goals

- Gate 1 (intent approval) and Gate 2 (delivery) stay the two human gates.
- Delivery remains an explicit ask.
- Explore stays planless.
- The runtime still must not write or interpret Product Knowledge. In-place context edits are a coordinator act after mark-done.
- Do not apply `0034` as `context/domains/conversation-spec-library/`.
- Preserve unrelated and dirty work, including untracked `intent/i015-lean-workspace-context/`. Do not commit, publish, or deliver as part of this plan.
- Run the semantic acceptance suite after the proving task.

## Product Knowledge grounding

- `completion` (`context/domains/completion/README.md`) — owns done-trigger and reconcile; currently stages proposals and infers Standard done from delivery.
- `delivery` (`context/domains/delivery/README.md`) — Gate 2; must not start reconcile or mark done.
- `plan-authorization` (`context/domains/plan-authorization/README.md`) — next-plan start must not wait on knowledge debt.
- `docs/product-knowledge.md` — currently teaches proposal staging and a separate accept.

Grounding summary: derived from the approved i014 contract and the post-approval look at this repository. One lifecycle identity change (when a plan becomes done, and how knowledge is written).

Decomposition decision: keep **one plan with five embedded tasks**. One repository, one worker, one independent check. Splitting contracts, runtime, and tests into stacked plans would leave contradictory owners mid-stack.

## Repositories and source evidence

`context-circuit-source`. The look mapped contracts, runtime, skills, shipped docs, live `context/`, `template/`, tests, and `agent-harness/`. `template/` and `agent-harness/` live in this same repository and must change with the product path.

## Tasks

1. **DCU-001** — rework knowledge and completion contracts. Depends on: none.
2. **DCU-002** — make mark-done the only done trigger in the runtime. Depends on: DCU-001.
3. **DCU-003** — teach coordinator, skills, and shipped docs the in-place path. Depends on: DCU-001.
4. **DCU-004** — drop pending proposal `0034` and rewrite live Product Knowledge. Depends on: DCU-001.
5. **DCU-005** — prove the old path is gone. Depends on: DCU-002, DCU-003, DCU-004.

## Acceptance criteria

- `DCU-AC-001` — Gathering and reconcile write live `context/` files; they do not stage `context/proposals/`.
- `DCU-AC-002` — There is no separate “accept the context update” decision.
- `DCU-AC-003` — A plan becomes done only when a human asks to mark it done. Delivery, verification, and candidate acceptance do not.
- `DCU-AC-004` — Asking for a pull request, merge, or sync does not start reconciliation and does not mark the plan done.
- `DCU-AC-005` — Mark-done marks the plan done, then updates live context files only if the plan affected Product Knowledge.
- `DCU-AC-006` — A later plan can start even if Product Knowledge from a prior plan has not been updated.
- `DCU-AC-007` — Pending proposal `0034` is gone and was not applied as a domain page.
- `DCU-AC-008` — Gathering still writes in place; mark-done starts that same in-place reconcile, not a proposal sidecar.

## Verification

- `DCU-VT-001` — `rg` over skills, adapters, and docs finds no remaining `context/proposals/` staging procedure for gather/reconcile.
- `DCU-VT-002` — `rg` finds no remaining “accept the context update” route.
- `DCU-VT-003` — `sh test/completion/test-inferred.sh` — delivery plus inferred completion must not mark a Standard plan done; change-set complete must not mark members done.
- `DCU-VT-004` — delivery skill and delivery-record do not start reconcile or inferred done.
- `DCU-VT-005` — `sh test/knowledge/test-debt.sh` — overlapping later plan is not blocked solely because Product Knowledge is not yet updated.
- `DCU-VT-006` — proposal `0034` file gone, no `conversation-spec-library` domain page, `context/INDEX.md` has no pending `0034` row.
- `DCU-VT-007` — `sh test/acceptance.sh` exits 0 after the rewritten completion, knowledge, scenario, contract, and release checks.

## Assumptions, open questions, and risks

Assumptions (plan-level, already decided for implementation):

- After mark-done, the coordinator compares the plan’s knowledge references, changed paths, and evidence to live context units, then edits `context/` and `INDEX.md` only when there is a durable knowledge change. The runtime does not write `context/`.
- Remove the `context/proposals/` staging path (including its README) from the product and the shipped seed; stop `workspace-init` from recreating it.
- Stop shipping `context-proposal.yaml` as a live schema; update contract and release checks to match.
- Knowledge-debt markers must not block the next plan and must not be started by delivery. Prefer removing the block from the plan-start check rather than keeping a silent pending marker that still gates work.
- `completion-infer` must not mark a plan done. Keep it only as a refusal, or drop the verb once tests no longer call it.

Open questions: none remaining at the intent level.

Risks:

- Lifecycle identity is implemented in several owners at once (invariants, runtime, skills, tests, harness). Changing only one leaves contradictions.
- Change-set complete currently fuses delivery recording with marking every member done. Those must split.
- Harness generated cases are build artifacts; plots and generated files must change together.
- Plan `0001-conversation-spec-library` currently cites the `0034` proposal path; drop the citation, do not apply the proposal to keep the citation alive.

## Expected commits and delivery notes

Source-only Conventional Commits on this branch. Delivery remains a separate explicit action.

## Expected Product Knowledge impact

In-place update expected on `completion`, `delivery`, `context/INDEX.md`, and `docs/product-knowledge.md`. Proposal `0034` is dropped, not accepted. Review when this plan is marked done.
