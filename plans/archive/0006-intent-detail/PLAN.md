# 0006 — Intent detail via the system-design skill

Plan ID: `0006-intent-detail`  
Intent: `i004-intent-detail`  
Status: `draft`

## Original request and coverage

When someone is about to agree to a change, they can see it written out in full, split by topic, and kept with that change. The short approval page stays short. The existing system-design helper writes that fuller picture beside the change; later planning follows those topics. This plan covers the whole approved change: both authoring homes, the optional offer and plan use, and proof of the two paths.

## Objective and desired behavior

The system-design helper still writes a product-level design under `sources/system-design/`, and can also write an intent's fuller by-concern write-up under `intent/<id>/detail/` with the same three-tier layout and split-when-earned rule. A human can request that write-up, or be recommended it when one change has several concerns; skipping it does not block approval, and a small single-outcome change stays short. Approval remains the one conversational decision. When the fuller write-up exists, planning uses it as the confirmed topic shape without replacing the post-approval read of the real code. The human-facing intent page stays the five-section approval surface.

## Constraints and non-goals

- Folder name is `detail/` under `intent/<id>/`, not `design/`.
- Split by concern, never by repository; reuse the existing three-tier rubric.
- `cc-system-design` owns the authoring rubric for both roots; `cc-intent` recommends or honors a request and does not duplicate that rubric.
- `INTENT.md` does not gain a sixth section or internal filenames.
- Detail is drafted without reading the codebase; `contract.yaml` remains the frozen authority and detail is not part of `contract_digest`.
- Prefer skill and guidance changes; do not add an engine verb, schema field, or invariant.
- Discard the unsolicited `sources/system-design/context-circuit/v1.1.0/` tree rather than treating it as specification.
- Do not add a new skill, a second approval gate, or a required detail for every intent.
- Do not change Explore's planless path, tracing-before-planning, delivery, or scope-safety ownership.
- Do not turn intent detail into Product Knowledge or a lifecycle object with its own status.
- Do not modify `wrapper/`, `template/`, or `agent-harness/` in this plan.
- Preserve unrelated and dirty work; do not commit, publish, or deliver as part of this plan.

## Product Knowledge grounding

- `system-design-authoring` (`context/domains/system-design-authoring/README.md`) — owns the three-tier rubric and must gain the second home without a runtime surface.
- `intent` (`context/domains/intent/README.md`) — owns Gate 1, the five-section approval page, and the optional recommend/request path.
- `tracing` (`context/domains/tracing/README.md`) — keeps the post-approval code read and done-checks out of the frozen contract; detail is not a tracing output.
- `plan-authorization` (`context/domains/plan-authorization/README.md`) — preserves a single intent gate and automatic plan derivation.
- `planning` (`docs/planning.md`) — consumes detail when present as confirmed topic shape alongside the trace.

Grounding summary: the approved change is feasible on skill, guidance, and tests. There is no required modify outside the bound scope, no unresolved intent-level question, and no tier raise. The unsolicited `v1.1.0` draft files are already gone; remaining empty directories must still be removed and not revived as spec.

## Repositories and source evidence

`context-circuit-source` is the sole repository. The code look-over maps the change to `.agents/skills/cc-system-design/`, `.agents/skills/cc-intent/`, `.agents/skills/cc-plan/`, `.agents/skills/cc-trace/`, `agents/coordinator.md`, `docs/`, `context/domains/`, `test/contracts/`, `test/intent/`, and the leftover empty `sources/system-design/context-circuit/v1.1.0/` tree. Engine validate already allows extra files under an intent directory; the contract digest hashes only `contract.yaml`.

## Tasks

1. **IDT-001** (context-circuit-source; paths: `.agents/skills/cc-system-design/`, `context/domains/system-design-authoring/`, `docs/terminology.md`, `sources/system-design/context-circuit/v1.1.0/`; depends on: none) — teach the design helper both homes and discard the unsolicited `v1.1.0` tree. Acceptance: `IDT-AC-001`, `IDT-AC-004`. Verification: `IDT-VT-001`, `IDT-VT-004`.

2. **IDT-002** (context-circuit-source; paths: `.agents/skills/cc-intent/`, `.agents/skills/cc-plan/`, `.agents/skills/cc-trace/`, `agents/coordinator.md`, `context/domains/intent/`, `context/domains/tracing/`, `context/domains/plan-authorization/`, `docs/templates/intent.md`, `docs/planning.md`; depends on: IDT-001) — offer or honor a request for the fuller write-up while preparing a change, keep the short approval page, and use the write-up when planning. Acceptance: `IDT-AC-002`, `IDT-AC-003`, `IDT-AC-005`, `IDT-AC-006`. Verification: `IDT-VT-002`, `IDT-VT-003`, `IDT-VT-005`, `IDT-VT-006`.

3. **IDT-003** (context-circuit-source; paths: `test/contracts/`, `test/intent/`, `test/acceptance.sh`, `test/acceptance/criteria-map.yaml`; depends on: IDT-002) — prove a short change can stay short and a richer change can get the fuller write-up, with no second approval. Acceptance: `IDT-AC-001`–`IDT-AC-006`. Verification: `IDT-VT-007`, `IDT-VT-008`.

## Acceptance criteria

- `IDT-AC-001` — The system-design skill still authors a product-level design under `sources/system-design/`, and can also author an intent's fuller write-up under `intent/<id>/detail/` with the same three-tier layout and split-when-earned rule.
- `IDT-AC-002` — A human can request the fuller write-up, or be recommended it when one intent has several concerns; skipping it does not block approving the intent, and a small single-outcome intent can stay short.
- `IDT-AC-003` — Approving the intent remains the one conversational decision; there is no separate approval of the detail, and plan derivation after approval is unchanged.
- `IDT-AC-004` — This capability is specified and shipped through the system-design skill and the related intent/plan/trace guidance, not through a new `sources/system-design` scope written for this feature.
- `IDT-AC-005` — When an intent has a detail write-up, later planning uses it as the confirmed shape of what to build so plans and tasks follow those topics, without replacing the post-approval read of the real code.
- `IDT-AC-006` — The human-facing intent page stays the short five-section approval surface; the fuller write-up is extra files with the intent, not a longer approval page.

## Verification

- `IDT-VT-001` — `rg -n 'sources/system-design/|intent/<id>/detail/|three-tier|split when|by concern' .agents/skills/cc-system-design/SKILL.md context/domains/system-design-authoring/README.md && rg -n 'sources/system-design/' .agents/skills/cc-system-design/SKILL.md && rg -n 'intent/.*/detail|intent/<id>/detail' .agents/skills/cc-system-design/SKILL.md`; skill and domain page document both roots and the shared rubric; the product-level path remains.
- `IDT-VT-002` — `rg -n -i 'recommend|request|optional|several concerns|does not block|not required|skip' .agents/skills/cc-intent/SKILL.md context/domains/intent/README.md agents/coordinator.md`; guidance states request or recommend, skip does not block Gate 1, and a small change may stay short.
- `IDT-VT-003` — `rg -n -i 'no separate approval|no second.*(gate|approval)|Gate 1|single.*(upstream|conversational)|detail.*(not|never).*approv' .agents/skills/cc-intent/SKILL.md .agents/skills/cc-plan/SKILL.md context/domains/intent/README.md context/domains/plan-authorization/README.md docs/planning.md && sh test/intent/test-feasibility.sh`; single intent approval; no detail-approval gate; feasibility suite passes.
- `IDT-VT-004` — `test ! -e sources/system-design/context-circuit/v1.1.0 && test ! -e .agents/skills/cc-intent-detail && rg -n 'intent/<id>/detail|intent/.*/detail' .agents/skills/cc-system-design/SKILL.md .agents/skills/cc-intent/SKILL.md && ! rg -n 'v1.1.0/intent-detail' .agents/skills/cc-system-design/SKILL.md .agents/skills/cc-intent/SKILL.md .agents/skills/cc-plan/SKILL.md context/domains/system-design-authoring/README.md context/domains/intent/README.md`; `v1.1.0` tree gone; no new skill; capability specified in skill/guidance.
- `IDT-VT-005` — `rg -n -i 'detail/|confirmed shape|topic|when.*detail.*exist|without replacing|trace manifest' .agents/skills/cc-plan/SKILL.md docs/planning.md context/domains/tracing/README.md`; planning consumes detail when present and still requires the post-approval code read.
- `IDT-VT-006` — `rg -n 'five sections|exactly five|Intention|Expectations|The plans|How carefully this is checked|Open questions' docs/templates/intent.md .agents/skills/cc-intent/SKILL.md context/domains/intent/README.md && ! rg -n '## Detail|sixth section|detail/design.md' docs/templates/intent.md docs/templates/intent.example.md`; five-section approval page; detail not folded into `INTENT.md`.
- `IDT-VT-007` — `sh test/contracts/test-contracts.sh`; contracts suite passes with extended dual-root rubric assertions and continued no-runtime-surface checks.
- `IDT-VT-008` — `sh test/acceptance.sh`; full semantic acceptance suite passes, including any new intent-detail coverage wired into `acceptance.sh` / `criteria-map.yaml`.

## Assumptions, open questions, and risks

- Assumption: recommend the fuller write-up when one draft has several concerns — several numbered plans, several independent outcomes, or the human already described several parts of one change. Always honor an explicit request. Never require it.
- Assumption: blank-seed layout docs under `template/` and conversation plots under `agent-harness/` can wait; skill, domain, and test checks are enough to prove the approved outcomes.
- Assumption: no engine verb, schema field, or new invariant is required; extra files under an intent directory are already allowed and are not digested.
- Open question: none remaining at the intent level. Folder name (`detail/`) and helper choice (extend the existing design helper) are already decided.
- Risk: a dual-root edit could imply intent detail is Product Knowledge, a lifecycle object, or a second approval.
- Risk: without a clear contrast, agents may still write a product-level design, or duplicate the rubric inside the intent helper.
- Risk: planning could ignore an existing detail write-up, or treat it as a replacement for reading the real code.
- Risk: extending the contracts suite must keep the product-level `sources/system-design/` assertion and the no-runtime-surface checks.
- Risk: unrelated dirty work on other intents must not be normalized.

## Expected commits and delivery notes

One source-only Conventional Commit is expected for this repository. Delivery, merge, push, and publication remain separate explicit actions.

## Expected Product Knowledge impact

In-place updates are expected on `system-design-authoring`, `intent`, `tracing`, and `plan-authorization`. Review those units at completion. No new context unit is required.
