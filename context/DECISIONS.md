# Decisions

## 2026-08-21 — wrapper/template split

Decision: keep shipped wrapper ownership under `.context-circuit/wrapper/` and the blank mutable
seed under `template/`.

Rationale: source identity must not be confused with an instantiated workspace.
Consequence: release assembly overlays root adapters and template state while
preserving no user data.

## 2026-08-21 — separate lifecycle gates

Decision: review, intent approval, execution, completion, delivery, archive,
takeover, and cleanup remain separate human actions.

Rationale: eligibility is not authorization and runtime evidence is not Done.
Consequence: the router emits an exact gate or read-only recommendation.

## 2026-08-24 — built-template harness location

Decision: the source-only built-template behavior laboratory lives at top-level
`template-harness/`, relocated and renamed from `test/template-runtime/`.

Rationale: it assembles and drives the released product, distinct in kind from
the `test/` engine and contract suites; the name reflects its assemble → drive →
grade role.
Consequence: the release allowlist still excludes it (it cannot ship);
`test/lib/assert.sh` resolves the repo root git-based so suites may live at any
nesting. Accepted from proposal `0002-template-harness-and-terminology-decision`.

## 2026-08-24 — terminology as a context category

Decision: Terminology is a standard context category. A per-project glossary
seed lives at `context/TERMINOLOGY.md`; the canonical Context Circuit glossary
and its internal → user-facing translation live in `.context-circuit/docs/terminology.md`,
projected from design chapter 08, and the coordinator references that projection
instead of hardcoding the term list.

Rationale: one owner per rule — chapter 08 owns term meaning and translation,
`.context-circuit/docs/terminology.md` is the shipped projection, the coordinator role references
it rather than duplicating policy.
Consequence: the shipped template gains `context/TERMINOLOGY.md`,
`.context-circuit/docs/terminology.md`, and `.context-circuit/docs/templates/terminology-context.md`. Accepted from
proposal `0002-template-harness-and-terminology-decision`.

## 2026-08-24 — current-state context refresh

Decision: after `plans/context-circuit-plans/` was deleted in commit `4b8ac0b`
as obsolete previous-version plans, re-ground the entire `context/` Product
Knowledge on the current wrapper (skills, `engine.sh`, schemas, `invariants.yaml`,
.context-circuit/docs/adapters) rather than on the deleted plans, and accept the result.

Rationale: the domain/role pages were reverse-engineered from previous-version
plans and had drifted from the shipped contract (consolidated invariant ids,
renamed skills, deleted `routes.yaml`/`context-sets.yaml`/`.context-circuit/docs/gates.md`, and an
approval card/token mechanism that contradicted `INV-APPROVE-01`). The current
wrapper is the authoritative present-day evidence.

Consequence: the domain set is now the nine domains that mirror the shipped
lifecycle — repository-binding (broadened with orientation), plan-review
(broadened to planning), plan-authorization, plan-execution, verification,
completion, plan-organization, delivery, and host-adapters — plus the refreshed
maintainer role, INDEX owner pointers, and terminology authority pointer.
`context/sources.yaml` provenance was retired (the wrapper is not a `sources/`
read). The proposals for this refresh were consumed from `context/proposals/`
on acceptance, as usual.

## 2026-08-24 — design↔context reconciliation

Decision: reconcile the accepted `context/` Product Knowledge against the v0.5
design under `sources/system-design/context-circuit/v0.5/`, under the policy
"shipped is truth; log deltas". Keep "wrapper" as an accepted synonym for the
universal project workspace product and update the design source to accept it
(rather than aligning context to the design's deprecation).

Rationale: the current shipped state is built toward that design, so context
should describe shipped behavior while design↔implementation divergences stay
visible instead of silently overwriting either side.

Consequence: accepted the reconciliation proposals. Fixed contradictions in
ARCHITECTURE.md (removed the two-stage-router-as-canonical framing and the
"lease" wording; added the execution model), PROJECT.md (stated the universal
multi-repository workspace identity; removed the "operating system" metaphor and
router framing), and INDEX.md (retrieval-catalog role). Enriched the lifecycle domain pages from the shipped wrapper. Added the
`source-release-and-upgrade` domain, a `context/DESIGN-DELTAS.md` log (evidence
layers, repair-limit, plan-id reuse, terminology authority), and CONVENTIONS
policy-change escalation. The design source (`08-terminology.md`,
`09-source-and-template.md`) was updated to accept "wrapper", and the glossaries
gained the term.

## 2026-08-27 — v0.6 coordinated contract bump

Decision: `runtime_version` 0.5.0 → 0.6.0 and plan `accepted_schema_versions`
becomes `[1, 2]`; `execution.yaml` stays schema 1.

Rationale: `plan_dependencies` is load-bearing, so a v0.5 engine must refuse a
`schema_version: 2` plan rather than schedule it dependency-blind (INV-PLAN-05);
`execution.yaml` is private single-version evidence with no cross-version reader.
Consequence: only plans using `plan_dependencies` stamp 2; existing plans stay 1
and read identically on both engines. Accepted from proposal
`0021-change-decisions`.

## 2026-08-27 — concurrency is orchestration, not authority (run-stack)

Decision: executing a set of intent-authorized plans in one run changes only order
and overlap; conflicts are prevented (dependencies order waves, path leases
serialize file overlaps, a dependent's base already contains its prerequisites),
not resolved afterward.

Rationale: intent approval, verification, completion, and delivery gates must be
untouched; the runtime detects readiness/leases/bases deterministically and the
coordinator decides how many ready plans to launch — no scheduler heuristic in the
runtime (INV-RUNTIME-01, INV-CONCURRENCY-01/02).
Consequence: a failed or blocked plan holds only its descendants; unrelated
verified plans are unaffected. Accepted from proposal `0021-change-decisions`.

## 2026-08-27 — delivery drift guard

Decision: a plan whose recorded base has diverged from the current base tip is
rebased onto the tip and re-verified before its pull request (INV-DELIVER-01,
extended).

Rationale: a plan must never merge from a base that no longer reflects the branch
it will land on; the only merge the runtime authors is the integration base on a
plan's own branch, never a delivery merge.
Consequence: `cc_delivery_drift` / `cc_delivery_rebase`, with
`DELIVERY_REBASE_CONFLICT` reported as blocked. Accepted from proposal
`0021-change-decisions`.

## 2026-08-27 — repository grounding: reference, not capture

Decision: the worker reads and honors the target repository's own agent guidance,
discovered live from the worktree as data (INV-GROUND-01); precedence is CC
scope/safety on what/where and repo guidance on how within that scope
(INV-GROUND-02); the worker brief is a fixed template filled from the manifest and
delivered, never authored (INV-GROUND-03).

Rationale: repository knowledge should be discovered and referenced, not
hand-injected per prompt or captured into a per-repo profile.
Consequence: no per-repo profile and no `plan.yaml` field; worker friction returns
as `repository_friction` and becomes a proposal on the repo's own agent docs.
Accepted from proposal `0021-change-decisions`.

## 2026-08-27 — a system design is a source, not a lifecycle stage

Decision: the v0.6 system-design scope ships only the `cc-system-design` authoring
skill; a system design lives under `sources/system-design/` with no status,
acceptance gate, or runtime record, and feeds Product Knowledge and plans through
the existing flow.

Rationale: a system design is deliberation, authored as structured source; the
durable accepted residue still lives in `context/` via the existing proposal path.
Consequence: no engine change, no new schema, no new invariant, no WORKFLOW action,
and no first-class `design/` area. Accepted from proposal `0021-change-decisions`.

## 2026-08-27 — template-harness efficiency ledger made real

Decision: template-harness dimension D now measures per-action usage from the
runner's own result and compares it to case budgets (units fixed: `max_tokens` =
generated output tokens, `max_turns` = conversational turns; optional
`max_agent_turns`, `max_cost_usd`), staying soft (never gates a run).

Rationale: dimension D was declared in v0.5 but never measured; a budget with no
defined unit is unfalsifiable.
Consequence: maintainer tooling only — no product surface change. Accepted from
proposal `0021-change-decisions`.

## 2026-08-28 — the external surface is orthogonal to the core workflow

Decision: publishing Context Circuit data outward is a peer command (`cc-publish`),
never a phase, trigger, gate, dependency, or side effect of one; it runs only on
explicit manual invocation, every time, against the workspace as found.

Rationale: coupling an external, credentialed, provider-specific side effect to the
deterministic core loop would add weight to the runtime and the agent and hand an
outside system a foothold on plan state.
Consequence: INV-EXTERNAL-01; no runtime or network code, no hook, no lifecycle
listener — nothing in the workflow triggers a publication, so there is nothing to
hook. Accepted from proposal `0026-change-decisions`.

## 2026-08-28 — external surface: export-first, non-authoritative, self-contained

Decision: data flows Context Circuit → outward only; the external copy is one-way,
idempotent on re-run, and self-contained (no workspace file, path, id, or internal
mechanism leaks; a plan id in a title is the one allowed cross-reference);
`plan.yaml` stays canonical and nothing is written under `plans/`.

Rationale: import would bypass the human intent-authoring and approval gate; a leaked
internal makes the external copy unreadable to a lay reader.
Consequence: INV-EXTERNAL-02, INV-EXTERNAL-03; any future import must pass through the
normal authoring gate (INV-APPROVE-01). Records live under a user-owned `publication/`
folder created on first use; config is credential-free (INV-SEC-01) and the host / MCP
layer carries all provider weight (INV-RUNTIME-01 unchanged). Accepted from proposal
`0026-change-decisions`.

## 2026-08-28 — "publish" is the external surface's word; git delivery is push / pull request

Decision: the product reserves "publish"/"publication" for the external surface and
vacates it from git delivery, which speaks only in push / pull request / merge /
deliver.

Rationale: one word, one meaning everywhere, so no per-mention qualifier is ever
needed.
Consequence: a wording-only edit to INV-DELIVER-01, INV-RUNTIME-01,
`.context-circuit/wrapper/manifest.yaml`, and AC-16 (no behavior, version, or authority change); the
delivery page and glossaries follow (proposals `0024-change-delivery`,
`0025-change-terminology`). Accepted from proposal `0026-change-decisions`.

## 2026-08-28 — the external surface adds no core contract bump

Decision: the scope ships a skill (`cc-publish`) + two record schemas + a config
convention (`publication-config.yaml`) + three invariants, and changes no plan-schema,
execution, or runtime-version.

Rationale: it binds to data, never to control flow, so it needs no coordinated bump —
comparable to how system-design-authoring ships only a skill.
Consequence: a workspace that configures no publication is a v0.5-shaped workspace plus
the availability of `cc-publish`; kinds (`plan`, `thread`, future `docs`) are new
`cc-publish` behavior plus a `config.yaml`, never new authority. Accepted from proposal
`0026-change-decisions`.

## 2026-09-03 — direct collaboration is the Explore tier, not a separate mode

Decision: adopt `cc-pair` as the **Explore tier** of the single assurance ladder
(Explore / Standard / Critical, INV-ASSURE-01), not an orthogonal working mode.
One connected repository per session; user plus coordinator plus one worker; no
verifier, lease, execution record, failure counter, plan status, or completion
gate; a fresh isolated `cc-pair/<session>` branch and worktree from a chosen base;
a light resumable pointer; output labeled human-supervised, never verified;
explicit commits only; closure preserves the branch and worktree; separate
delivery that blocks on base drift rather than silently rebasing. An explicit
**promote** step (attach an intent via `cc-intent`, raise the tier so the verifier
appears, author a lightweight plan of record via `cc-plan` bound to the existing
pairing commits) turns a session into a candidate-bearing change in place.

Rationale: the v0.7.0 work shipped `cc-pair` and its context proposals
(`0029`–`0033`) framed it as a standalone mode "outside the plan lifecycle." v1.0
reworked INV-PAIR-01 and added INV-ASSURE-01 so the same mechanics became the
bottom rung of one ladder with a ramp upward; Product Knowledge must describe the
shipped v1.0 system, so the standalone-mode framing is superseded rather than
recorded.

Consequence: accepted proposals `0029`–`0033`, **rewritten to the v1.0
assurance-ladder framing** — a new [direct-collaboration](domains/direct-collaboration/README.md)
domain, and extensions to [host-adapters](domains/host-adapters/README.md)
(worker child serves pairing; no verifier; host-blocked fails closed),
[delivery](domains/delivery/README.md) (pairing delivery blocks on base drift,
never auto-rebases), and ARCHITECTURE.md (the ladder replaces the
"two orthogonal modes" framing). INV-PAIR-01, INV-ASSURE-01, and INV-HOST-01 are
unchanged. The `0034-add-conversation-spec-library` proposal was left pending by
explicit request.

## 2026-09-03 — agent-harness rename

Decision: the source-only built-template behavior laboratory, previously at
`template-harness/`, is renamed to `agent-harness/`; its two-word prose form
becomes "agent harness". This supersedes the location/name half of the
2026-08-24 "built-template harness location" decision above; that entry stays as
the historical record.

Rationale: the harness assembles and drives the released product through the
host agents (coordinator, worker, verifier) and the human simulator, so "agent
harness" names what it exercises rather than that it is seeded from the template.

Scope of the rename: the directory and all live references (its own files,
`README.md`, `.gitignore`, `scripts/release-manifest.txt`, `test/`, `.claude/`,
the live intents and context pages). The immutable proposal/decision id
`0002-template-harness-and-terminology-decision` is preserved verbatim (ids are
never rewritten), and historical `sources/` design docs and archived `plans/`
are left as they were.

## 2026-09-04 — context/ re-grounded to Context Circuit v1.0 (tracer/feasibility)

Decision: reconcile the Product Knowledge base with the shipped v1.0 trust-core.
The v1.0 layer (intent gate, post-approval tracing, feasibility check, consequence
tiering, scope-free authorization) had landed in `.context-circuit/wrapper/` without ever passing
through `context/` — no proposal, no domain pages. This re-ground (a) retires the
earlier v1.0 design's spec-adversary and scope-envelope vocabulary from the domain
pages and `roles/maintainer.md`, replacing it with the scope-free
`intent-authorized` model (approved intent + unchanged criteria); (b) adds three
`status: proposed` domain pages — [intent](domains/intent/README.md) (Gate 1),
[tracing](domains/tracing/README.md) (tracer + feasibility check), and
[assurance](domains/assurance/README.md) (the consequence-tier ladder, the one
safety-critical automated check); and (c) re-grounds delivery (named as Gate 2,
where scope-safety is now settled), verification (tier-conditional, candidate-bound),
and repository-grounding (disambiguated from the new tracer).

Rationale: one rule, one owner, and "shipped is truth" — the knowledge base must
describe the product that ships. v1.0 removed the automated scope gate, so
scope-safety lives entirely at delivery; that shift must be visible where an agent
routes work, not only in `.context-circuit/wrapper/`.

Consequence: the three new pages were accepted by the maintainer on 2026-09-04.
Edits to existing accepted pages are corrections toward shipped truth and preserve
their accepted decisions, each noted in its acceptance notes. The retired
spec-adversary artifact `intent/i001-conversation-spec-library/adversary.md` was
removed and its human-accepted criteria residuals preserved in that intent's
`contract.yaml` machine record. The stale `release/requests/0.0.1-beta.1.md` — which
still described the earlier spec-adversary + scope-envelope design — was dropped
because the design has since shifted; the built `dist/…beta.1` artifact remains a
stale build to be re-cut as a separate release concern.

## 2026-09-06 — pending proposal 0034 dropped

Decision: drop pending proposal `0034-add-conversation-spec-library`. Do not
apply it as `context/domains/conversation-spec-library/`. Product Knowledge is
now updated by writing live `context/` files in place; the staging path is
removed rather than left empty. The 2026-09-03 note that 0034 was left pending
stands as the historical record of that earlier request.

Rationale: gathering and mark-done reconcile write live files; there is no
sidecar and no extra knowledge-acceptance gate.

Consequence: the pending conversation-spec-library knowledge was dropped, not
applied as a domain page.
