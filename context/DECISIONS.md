# Decisions

## 2026-09-10 — Intent open questions are numbered like The plans

Decision: on a new human-facing intent, Open questions use the same numbered 1, 2, 3 list as The plans — number, bold question, italic answer beneath. Answering keeps the number. A question added later takes the next unused number. When there are none, the empty-state line stays unnumbered; no dummy numbered item. Already-written intents stay as authored. Numbering is human-facing only and is not a machine identifier on the contract.

Rationale: a person can answer by number instead of repeating the question, and the number remains a stable handle for the life of that intent.

Consequence: [intent](domains/intent/README.md) owns the convention with the intent template and the `cc-intent` skill. The approval ask invites an answer by number. Plan open questions and thread questions are unchanged.

## 2026-09-10 — GitHub names match source vs published template

Decision: this maintainer source lives on GitHub as `kaotypr/context-circuit-source`. The published template lives on GitHub as `kaotypr/context-circuit`. The conceptual product identity remains `context-circuit-template` (the distributable universal project workspace). GitLab remains an optional mirror of the already-published tree; its project path is not renamed with GitHub.

Rationale: the short GitHub name is the product people use. The factory must not occupy that name.

Consequence: [source-release-and-upgrade](domains/source-release-and-upgrade/README.md) names both GitHub destinations. Publication binding, the publish Action, and this workspace's canonical URL follow those slugs. GitHub stays the canonical published template.

## 2026-09-09 — Published template is GitHub-canonical with an optional GitLab mirror

Decision: official template publication still lands on GitHub (`main`, the
version tag, and the GitHub Release). The same publish Action then pushes that
same assembled tree and tag to GitLab and creates a GitLab Release with the
same notes and archive. GitHub stays canonical; GitLab is a mirror of the
already-published tree, not a second assembly. The GitLab project path and
token live only in GitHub Actions settings (`GITLAB_TEMPLATE_PROJECT`,
`GITLAB_TEMPLATE_TOKEN`). The GitLab project must already exist. A GitLab push
or Release failure fails the job.

Rationale: the published template has one canonical GitHub home. A second host
can receive the same commit without a second assembler and without credentials
in workspace files.

Consequence: [source-release-and-upgrade](domains/source-release-and-upgrade/README.md)
documents the operator variables and token scopes (`write_repository` and
`api`). Runtime and `cc-publish` are unchanged.

## 2026-09-08 — Intent and plan ids allocate from per-member bands

Decision: new intent ids (`i<NNN>-slug`) and plan ids (`NNNN-slug`) are allocated
from non-overlapping numeric bands on a committed member roster. Each machine
names one roster member once in gitignored `member.local.yaml`. Next id is the
highest number ever used in that member's band among active and archived work,
plus one. A stack from one intent still receives consecutive plan ids inside
that band. Exhaustion fails clearly; archived numbers are never reused. Member
names never appear in the id or in execution branches
(`cc/<plan-id>/<repository-id>`).

Rationale: teammates sharing one workspace can create intents and plans in
parallel without waiting for each other to push, while keeping short speakable
sequential ids and a one-time local identity instead of a block number in chat.

Consequence: INV-MEMBER-01 owns the roster and local identity. INV-INTENT-01 and
INV-PLAN-03 are band-scoped never-reuse. Missing identity fails closed. Domain
pages: [repository-binding](domains/repository-binding/README.md),
[plan-review](domains/plan-review/README.md), [intent](domains/intent/README.md),
[run-stack](domains/run-stack/README.md).

## 2026-09-08 — Product Knowledge is durable-only

Decision: live context files hold durable product knowledge only. They never
name a particular plan, intent file, or sources file, and never name a path
into the sources tree. `DECISIONS.md` records what is now true about the
product — decision, rationale, and consequence — not edited paths or the
ephemeral artifact behind the change. Explaining the concepts plan, intent,
and sources remains allowed; path patterns that describe product structure
(for example `intent/<id>/`) are not particular-file cites. Invariant IDs and
shipped contract or adapter paths remain nameable.

Rationale: retrieval and in-place updates already owned how knowledge is found
and written; they did not ban ephemeral citations. Provenance footnotes and
decision history that named particular artifacts made the catalog brittle as
plans and source layouts moved.

Consequence: INV-KNOWLEDGE-03 owns the writing rule. There is no provenance
ledger file beside the concept page. Domain pages keep Implementation
references to shipped contract paths. New workspaces receive the same
durable-only seed.

## 2026-09-08 — host-native integration surface

Decision: Claude Code, Codex CLI, and Cursor Agent CLI each integrate through
committed host-native project folders — `.claude/`, `.codex/`, `.cursor/` — as
thin routes into Context Circuit owners (`.context-circuit/agents`, owning
invariants, `.agents/skills/cc-*`). Root instruction adapters remain; the native
trees are how each host discovers worker, verifier, and planner roles and
standing rules in that host's format. The product host set ships through release
assembly and the blank template seed; maintainer-only Claude extras stay out of
the artifact.

Rationale: each host already searches its own project tree; Context Circuit
should meet it there instead of treating those folders as optional host-local
convenience. The routes must not copy role bodies or invent a second
authorization policy (INV-HOST-01).

Consequence: extensions to [host-adapters](domains/host-adapters/README.md)
(native trees as the integration surface, child mapping, maintainer-only
filtering) and [source-release-and-upgrade](domains/source-release-and-upgrade/README.md)
(host folders join the shipped set). `ARCHITECTURE.md` names the dual entry
(root adapters plus native routes). Missing required native child support
remains read-only and reports `host-blocked`.

## 2026-09-07 — nested product home under `.context-circuit/`

Decision: nest the shipped wrapper, role files, and product docs under
`.context-circuit/` (`.context-circuit/wrapper`, `.context-circuit/agents`,
`.context-circuit/docs`). Keep `.agents/` at the workspace root. Keep
`template/` as the source-root seed. Source-checkout `AGENTS.md` /
`WORKFLOW.md` stay maintainer-specific; source `CLAUDE.md` / `CURSOR.md`
import that root `AGENTS.md` and must not import nested adapters (that would
replace maintainer instructions with the shipped product copy). Instantiated
workspaces get adapter copies at the root. An upgrade moves template-owned
product folders under `.context-circuit/` and preserves workspace-owned files.

Rationale: the workspace root belongs to the project, not the product.
Consequence: engine lookups, assembler staging, skills, tests, and Product
Knowledge name `.context-circuit/...`; leftover root `wrapper/`, `agents/`,
or `docs/` are gone.

## 2026-09-06 — pending conversation-spec-library knowledge dropped

Decision: do not apply the pending conversation-spec-library knowledge as a
domain page. Product Knowledge is now updated by writing live `context/` files
in place; the staging path is removed rather than left empty. The earlier note
that this knowledge was left pending stands as the historical record of that
request.

Rationale: gathering and mark-done reconcile write live files; there is no
sidecar and no extra knowledge-acceptance gate.

Consequence: the pending conversation-spec-library knowledge was dropped, not
applied as a domain page.

## 2026-09-04 — context re-grounded to Context Circuit v1.0 (planner/feasibility)

Decision: reconcile the Product Knowledge base with the shipped v1.0 trust-core.
The v1.0 layer (intent gate, post-approval tracing, feasibility check, consequence
tiering, scope-free authorization) had landed in `.context-circuit/wrapper/` without
ever passing through `context/` — no domain pages. This re-ground (a) retires the
earlier v1.0 design's spec-adversary and scope-envelope vocabulary from the domain
pages and `roles/maintainer.md`, replacing it with the scope-free
`intent-authorized` model (approved intent + unchanged criteria); (b) adds three
domain pages — [intent](domains/intent/README.md) (Gate 1),
[tracing](domains/tracing/README.md) (planner + feasibility check), and
[assurance](domains/assurance/README.md) (the consequence-tier ladder, the one
safety-critical automated check); and (c) re-grounds delivery (named as Gate 2,
where scope-safety is now settled), verification (tier-conditional, candidate-bound),
and repository-grounding (disambiguated from the post-approval planner).

Rationale: one rule, one owner, and "shipped is truth" — the knowledge base must
describe the product that ships. v1.0 removed the automated scope gate, so
scope-safety lives entirely at delivery; that shift must be visible where an agent
routes work, not only in `.context-circuit/wrapper/`.

Consequence: the three new pages were accepted by the maintainer. Edits to
existing accepted pages are corrections toward shipped truth and preserve their
accepted decisions. The retired spec-adversary artifact was removed and its
human-accepted criteria residuals preserved in that intent's machine record. A
stale release request that still described the earlier spec-adversary +
scope-envelope design was dropped because the design has since shifted; the
built beta artifact remains a stale build to be re-cut as a separate release
concern.

## 2026-09-03 — agent-harness rename

Decision: the source-only built-template behavior laboratory is named the
agent harness; its two-word prose form is "agent harness". This supersedes the
location/name half of the 2026-08-24 "built-template harness location" decision
above; that entry stays as the historical record.

Rationale: the harness assembles and drives the released product through the
host agents (coordinator, worker, verifier) and the human simulator, so "agent
harness" names what it exercises rather than that it is seeded from the template.

Consequence: live references use the new name. Immutable proposal and decision
ids are preserved verbatim (ids are never rewritten). Historical design docs and
archived plans are left as they were.

## 2026-09-03 — direct collaboration is the Explore tier, not a separate mode

Decision: adopt `cc-pair` as the **Explore tier** of the single assurance ladder
(Explore / Standard / Critical, INV-ASSURE-01), not an orthogonal working mode.
One connected repository per session; user plus coordinator plus one worker; no
verifier, lease, execution record, failure counter, plan status, or completion
gate; a fresh isolated pairing branch and worktree from a chosen base;
a light resumable pointer; output labeled human-supervised, never verified;
explicit commits only; closure preserves the branch and worktree; separate
delivery that blocks on base drift rather than silently rebasing. An explicit
**promote** step (attach an intent via `cc-intent`, raise the tier so the verifier
appears, author a lightweight plan of record via `cc-plan` bound to the existing
pairing commits) turns a session into a candidate-bearing change in place.

Rationale: an earlier framing treated pairing as a standalone mode outside the
plan lifecycle. v1.0 reworked INV-PAIR-01 and added INV-ASSURE-01 so the same
mechanics became the bottom rung of one ladder with a ramp upward; Product
Knowledge must describe the shipped v1.0 system, so the standalone-mode framing
is superseded rather than recorded.

Consequence: a new [direct-collaboration](domains/direct-collaboration/README.md)
domain, and extensions to [host-adapters](domains/host-adapters/README.md)
(worker child serves pairing; no verifier; host-blocked fails closed),
[delivery](domains/delivery/README.md) (pairing delivery blocks on base drift,
never auto-rebases), and ARCHITECTURE.md (the ladder replaces the
"two orthogonal modes" framing). INV-PAIR-01, INV-ASSURE-01, and INV-HOST-01 are
unchanged.

## 2026-08-28 — the external surface adds no core contract bump

Decision: the scope ships a skill (`cc-publish`) + two record schemas + a config
convention + three invariants, and changes no plan-schema, execution, or
runtime-version.

Rationale: it binds to data, never to control flow, so it needs no coordinated bump —
comparable to how system-design-authoring ships only a skill.
Consequence: a workspace that configures no publication is a v0.5-shaped workspace plus
the availability of `cc-publish`; kinds (`plan`, `thread`, future `docs`) are new
`cc-publish` behavior plus a config file, never new authority.

## 2026-08-28 — "publish" is the external surface's word; git delivery is push / pull request

Decision: the product reserves "publish"/"publication" for the external surface and
vacates it from git delivery, which speaks only in push / pull request / merge /
deliver.

Rationale: one word, one meaning everywhere, so no per-mention qualifier is ever
needed.
Consequence: a wording-only edit to INV-DELIVER-01, INV-RUNTIME-01,
`.context-circuit/wrapper/manifest.yaml`, and AC-16 (no behavior, version, or authority change); the
delivery page and glossaries follow.

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
layer carries all provider weight (INV-RUNTIME-01 unchanged).

## 2026-08-28 — the external surface is orthogonal to the core workflow

Decision: publishing Context Circuit data outward is a peer command (`cc-publish`),
never a phase, trigger, gate, dependency, or side effect of one; it runs only on
explicit manual invocation, every time, against the workspace as found.

Rationale: coupling an external, credentialed, provider-specific side effect to the
deterministic core loop would add weight to the runtime and the agent and hand an
outside system a foothold on plan state.
Consequence: INV-EXTERNAL-01; no runtime or network code, no hook, no lifecycle
listener — nothing in the workflow triggers a publication, so there is nothing to
hook.

## 2026-08-27 — template-harness efficiency ledger made real

Decision: the agent-harness efficiency dimension measures per-action usage from the
runner's own result and compares it to case budgets (units fixed: `max_tokens` =
generated output tokens, `max_turns` = conversational turns; optional
`max_agent_turns`, `max_cost_usd`), staying soft (never gates a run).

Rationale: the dimension was declared but never measured; a budget with no
defined unit is unfalsifiable.
Consequence: maintainer tooling only — no product surface change.

## 2026-08-27 — a system design is a source, not a lifecycle stage

Decision: system-design authoring ships only the `cc-system-design` skill; a
system design is passive source material with no status, acceptance gate, or
runtime record, and feeds Product Knowledge and plans through the existing flow.

Rationale: a system design is deliberation, authored as structured source; the
durable accepted residue still lives in `context/` via in-place updates.
Consequence: no engine change, no new schema, no new invariant, no WORKFLOW action,
and no first-class `design/` area.

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

## 2026-08-27 — delivery drift guard

Decision: a plan whose recorded base has diverged from the current base tip is
rebased onto the tip and re-verified before its pull request (INV-DELIVER-01,
extended).

Rationale: a plan must never merge from a base that no longer reflects the branch
it will land on; the only merge the runtime authors is the integration base on a
plan's own branch, never a delivery merge.
Consequence: delivery drift and rebase outcomes, with a rebase conflict reported
as blocked.

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
verified plans are unaffected.

## 2026-08-27 — v0.6 coordinated contract bump

Decision: `runtime_version` 0.5.0 → 0.6.0 and plan `accepted_schema_versions`
becomes `[1, 2]`; `execution.yaml` stays schema 1.

Rationale: `plan_dependencies` is load-bearing, so a v0.5 engine must refuse a
`schema_version: 2` plan rather than schedule it dependency-blind (INV-PLAN-05);
`execution.yaml` is private single-version evidence with no cross-version reader.
Consequence: only plans using `plan_dependencies` stamp 2; existing plans stay 1
and read identically on both engines.

## 2026-08-24 — design↔context reconciliation

Decision: reconcile the accepted `context/` Product Knowledge against the v0.5
system design, under the policy "shipped is truth; log deltas". Keep "wrapper"
as an accepted synonym for the universal project workspace product and update
the design source to accept it (rather than aligning context to the design's
deprecation).

Rationale: the current shipped state is built toward that design, so context
should describe shipped behavior while design↔implementation divergences stay
visible instead of silently overwriting either side.

Consequence: accepted the reconciliation. Fixed contradictions in
ARCHITECTURE.md (removed the two-stage-router-as-canonical framing and the
"lease" wording; added the execution model), PROJECT.md (stated the universal
multi-repository workspace identity; removed the "operating system" metaphor and
router framing), and INDEX.md (retrieval-catalog role). Enriched the lifecycle
domain pages from the shipped wrapper. Added the `source-release-and-upgrade`
domain, a design-delta log (evidence layers, repair-limit, plan-id reuse,
terminology authority), and CONVENTIONS policy-change escalation. The design
source was updated to accept "wrapper", and the glossaries gained the term.

## 2026-08-24 — current-state context refresh

Decision: after obsolete previous-version plans were removed, re-ground the
entire `context/` Product Knowledge on the current wrapper (skills, `engine.sh`,
schemas, `invariants.yaml`, adapters) rather than on the deleted plans, and
accept the result.

Rationale: the domain/role pages were reverse-engineered from previous-version
plans and had drifted from the shipped contract (consolidated invariant ids,
renamed skills, deleted routing and context-set artifacts, and an
approval card/token mechanism that contradicted `INV-APPROVE-01`). The current
wrapper is the authoritative present-day evidence.

Consequence: the domain set is the lifecycle domains that mirror the shipped
product — repository-binding (broadened with orientation), plan-review
(broadened to planning), plan-authorization, plan-execution, verification,
completion, plan-organization, delivery, and host-adapters — plus the refreshed
maintainer role, INDEX owner pointers, and terminology authority pointer.
The earlier provenance ledger for those deleted plans was retired (the wrapper
is not a source read).

## 2026-08-24 — terminology as a context category

Decision: Terminology is a standard context category. A per-project glossary
seed lives at `context/TERMINOLOGY.md`; the canonical Context Circuit glossary
and its internal → user-facing translation live in `.context-circuit/docs/terminology.md`,
projected from the terminology design chapter, and the coordinator references that
projection instead of hardcoding the term list.

Rationale: one owner per rule — the terminology design chapter owns term meaning
and translation, `.context-circuit/docs/terminology.md` is the shipped projection, the
coordinator role references it rather than duplicating policy.
Consequence: the shipped template gains `context/TERMINOLOGY.md`,
`.context-circuit/docs/terminology.md`, and `.context-circuit/docs/templates/terminology-context.md`.

## 2026-08-24 — built-template harness location

Decision: the source-only built-template behavior laboratory lives at the
repository root as a top-level harness, relocated and renamed from a nested
test path.

Rationale: it assembles and drives the released product, distinct in kind from
the `test/` engine and contract suites; the name reflects its assemble → drive →
grade role.
Consequence: the release allowlist still excludes it (it cannot ship);
`test/lib/assert.sh` resolves the repo root git-based so suites may live at any
nesting.

## 2026-08-21 — separate lifecycle gates

Decision: review, intent approval, execution, completion, delivery, archive,
takeover, and cleanup remain separate human actions.

Rationale: eligibility is not authorization and runtime evidence is not Done.
Consequence: the router emits an exact gate or read-only recommendation.

## 2026-08-21 — wrapper/template split

Decision: keep shipped wrapper ownership under `.context-circuit/wrapper/` and the blank mutable
seed under `template/`.

Rationale: source identity must not be confused with an instantiated workspace.
Consequence: release assembly overlays root adapters and template state while
preserving no user data.
