# Product Knowledge and context lifecycle

Product Knowledge is an agent-oriented understanding system for the project.
Its files are readable storage and human review surfaces, but its primary job
is to help an agent find, retrieve, and apply the right product and project
knowledge when creating plans. This document defines how context is gathered,
updated, and used without allowing every plan or repository observation to
silently rewrite accepted knowledge.

The workspace is not designed to store Markdown for its own sake. A context
file is useful when its meaning is easy for an agent to locate, understand,
scope, and connect to evidence.

## 1. Context categories

| Category | Examples | Authority |
| --- | --- | --- |
| Identity | workspace name, purpose, repository keys | workspace.yaml and identity pages |
| Architecture | services, modules, data flows, integration boundaries | accepted context pages |
| Conventions | language, testing, naming, review, delivery conventions | accepted context pages |
| Decisions | accepted tradeoffs and constraints | DECISIONS.md and linked pages |
| Domain knowledge | billing, checkout, permissions, workflows | context/domains |
| Role perspectives | support, operations, finance, platform | context/roles |
| Terminology | project-specific vocabulary, entity and state names | context/TERMINOLOGY.md and glossary pages |
| Provenance | source references, revisions, freshness | SOURCES.md and sources.yaml |
| Raw evidence | source documents, repository files, traces | sources/ and repositories |

A task description, runtime handoff, or single code observation is not accepted
Product Knowledge without evaluating scope and provenance.

## 2. Context page shape

A context page should be easy for an agent to retrieve and understand, while
remaining readable and reviewable by a human.

~~~markdown
# Billing domain

Context ID: domain.billing
Summary: Subscription state, lifecycle, and billing ownership.
Topics: billing, subscription, cancellation, invoice
Repositories: api, web, contracts
Applies when: a plan changes subscription behavior or billing UI

Status: accepted
Last reviewed: 2026-08-20
Owner: product-and-engineering
Freshness: review when billing API changes

## Stable facts

- A customer has at most one active subscription per account.
- Cancellation takes effect at the end of the paid period.

## Constraints

- Billing events are append-only.
- The dashboard must not calculate invoice totals independently of the API.

## Decisions

- Subscription state is owned by the API service.
- The contracts repository publishes event schemas.

## Open questions

- Whether trial conversion creates a separate event.

## Provenance

- sources/billing-requirements.md#subscription-lifecycle
- repositories/api/src/subscriptions/model.ts at abc123
- plans/0001-billing-v2/PLAN.md
~~~

Each page distinguishes stable facts, constraints, decisions, open questions, and
provenance. It should not contain a full implementation plan or copied source
article.

## 3. Evidence sources

A context proposal may use:

- a human-named source file;
- an accepted PRD or brief;
- a completed or approved plan;
- repository code and tests at a recorded revision;
- a verified worker or verifier handoff;
- an explicit human statement.

Evidence is not equal in scope. Repository implementation can show what code
does, but not necessarily what the product should do. A plan expresses intent,
but not necessarily what was implemented. A verifier can confirm behavior, but
not replace a domain decision.

## 4. Gathering context

The root agent's context-gathering capability follows this sequence. It is a
root/coordinator responsibility, not a separate child agent:

1. Define the topic and evidence boundary.
2. Read the relevant context page and provenance.
3. Read the named sources, plan sections, or repository paths.
4. Separate observed facts, proposed interpretation, and unresolved conflict.
5. Draft a context update with citations and freshness.
6. Show the proposal when it changes accepted knowledge.
7. Apply the accepted update atomically.
8. Record evidence revision and update the context index when needed.

A request that names exact source files may read those files. A general request
to “understand everything” should be narrowed to a domain or question; it must
not trigger a broad sources or repository scan.

## 5. Context updates from plans

Plans can produce context candidates after execution or review:

- a new architecture fact discovered while planning;
- a decision the human explicitly accepted;
- a convention demonstrated consistently by repositories;
- a domain rule clarified during review;
- a reusable verification or acceptance pattern.

The plan remains the authority for its execution. It does not automatically
become Product Knowledge because it is approved or verified.

A plan-derived proposal cites plan ID and revision, task or acceptance ID,
whether the fact was proposed, implemented, or independently verified, affected
repository and commit, and the target context page.

When a human marks a verified plan done, the workspace creates an
implementation completion record and performs a post-implementation knowledge
reconciliation. This compares the actual repository commits and verifier
evidence with the Product Knowledge that grounded the plan. It may conclude
that no durable update is needed, or create proposals for new, changed,
retired, stale, or conflicting knowledge.

## 6. Context updates from sources

Source material remains passive until the human accepts a structured proposal.
The proposal preserves:

- source path;
- source revision or digest;
- relevant section or locator;
- extracted statement;
- confidence;
- conflict with accepted context;
- proposed context page;
- whether the source is authoritative, informative, or unresolved.

Do not paste an entire source into context. Summarize the stable project-relevant
fact and preserve a provenance link.

## 7. Context updates from repositories

Repository evidence is read at a recorded revision. The agent may inspect:

- architecture and module boundaries;
- public interfaces;
- test conventions;
- configuration conventions;
- dependency and integration boundaries;
- decisions repeated intentionally across the codebase.

Do not automatically promote generated output, temporary branches, failing
experiments, secrets, a single accidental pattern, an unreviewed choice, or a
worker claim not confirmed by independent evidence.

A repository-derived proposal identifies whether it is observed behavior,
recommended convention, or accepted decision.

## 8. Conflicts and staleness

When new evidence disagrees with accepted context, do not silently overwrite
either side. Record:

- accepted statement;
- observed statement;
- evidence references;
- affected plan or repository;
- practical impact;
- proposed resolution.

The human chooses whether to retain accepted context, update it, or record a
decision explaining the exception.

A page is stale when its freshness rule expired, its cited repository revision
materially changed, or a relevant decision changed. Stale context may guide
read-only orientation, but a consequential plan should surface the stale
reference and refresh it before execution.

## 9. Context acceptance

A context update is accepted when the human explicitly accepts the proposed
facts or decision. Acceptance updates the page and provenance together.

Plan completion does not accept a context update. A verified implementation can
be evidence for a proposal, but code, worker claims, verifier success, and the
human's decision that the plan is done are not by themselves acceptance of
Product Knowledge. The plan may remain done while its context proposals are
pending or deferred.

If the update changes repository identity, branch policy, security boundary, or
execution behavior, it also updates the owning workspace or contract file
through that owner's normal action. A policy change must not live only in a
context page.

## 10. Context index

context/INDEX.md is an agent retrieval catalog and a human navigation aid. It
lists:

- page path;
- stable context ID;
- topic;
- concise summary;
- keywords and aliases;
- related domains and repositories;
- decisions and constraints covered;
- status;
- owner;
- freshness;
- primary related repositories;
- primary provenance.

The index is updated when a page is added, renamed, archived, or materially
changed. It does not duplicate full page content; it makes the content findable
without requiring an agent to scan the directory.

Pending context proposals and stale-context warnings are also indexed by target
context ID, related plan, topic, repository, and impact status so future plan
creation can find them without scanning all proposals.

## 11. Context brief for plan creation

Before creating a plan, the agent builds a context brief:

~~~yaml
topic: recurring-billing
workspace: acme-commerce
repositories: [api, web, contracts]
accepted_context:
  - id: domain.billing
    path: context/domains/billing/README.md
    reason: Governs billing behavior.
  - id: architecture.system
    path: context/ARCHITECTURE.md
    reason: Maps service ownership.
  - id: conventions.project
    path: context/CONVENTIONS.md
    reason: Defines project conventions.
decisions:
  - context/DECISIONS.md#shared-contracts
grounding_summary: API owns billing state; web renders it; contracts own schemas.
sources:
  - sources/billing-requirements.md
repository_evidence:
  - api: src/subscriptions at abc123
  - web: app/billing at xyz789
open_conflicts:
  - trial conversion behavior is unresolved
pending_context_impact:
  - plan: 0001-billing-v2
    status: review-needed
    reason: The latest implementation changed subscription ownership references.
~~~

The plan is drafted from this brief, not from a full workspace dump. Each
reference is traceable to a context page, source, repository path, or decision.

The root/coordinator uses context to answer:

- what the product is supposed to do;
- which repositories own which behavior;
- which conventions and constraints apply;
- which decisions must not be reopened casually;
- what evidence is needed for acceptance;
- what risks or conflicts require human attention.

The root/coordinator also performs a drift check. It compares the proposed objective,
repository map, task scope, constraints, decisions, and acceptance against the
selected knowledge units. A contradiction becomes an explicit open question or
plan risk; it is never hidden by choosing a convenient implementation.

The root/coordinator also performs a request-fidelity check. It compares the original
request with the plan's requirements, desired behavior, constraints, non-goals,
and expected outcomes. Missing detail becomes an open question or risk so the
human can resolve it during plan review.

## 12. Context during execution

The worker receives the plan's context references and reads relevant pages. It
does not receive all historical context or sources. If it discovers a context
gap, it stops at the smallest safe boundary and reports a proposed question.

The verifier uses context to understand promised behavior and required evidence.
It does not turn an implementation detail into a new product requirement during
verification.

## 13. Context after execution

After verification, the agent may prepare a preview of likely context impact,
but the final implementation reconciliation happens when the human marks the
plan done. This proposal is separate from the plan-status transition. A failed
execution can reveal a useful conflict or missing convention, but the failure
is not accepted as fact automatically.

### 13.1 Completion reconciliation

The reconciliation reads:

- the final plan and task files plus their revisions;
- the Product Knowledge references and grounding summary;
- changed paths and commits for every affected repository;
- worker handoffs and verifier evidence;
- current revisions of relevant context units.

It compares planned intent with implemented behavior and looks for:

- new product, domain, architecture, or convention facts;
- changed decisions or constraints;
- accepted context whose cited implementation is now stale or invalid;
- durable implementation patterns worth preserving;
- no durable knowledge change.

Each result has an impact status:

~~~text
not-assessed → review-needed → accepted
                         ↘ deferred
                         ↘ conflict
no-update-needed
~~~

For `review-needed`, `deferred`, or `conflict`, the workspace stores a
reviewable proposal under `context/proposals/` and records its relationship to
the implementation completion record. A proposal contains the target context
unit, proposed add/change/retire/stale operation, concise statement, evidence
references, affected repositories and commits, confidence, and conflict or
open-question details.

The agent reports the result after completion and responds to requests such as
“review context updates for 0001-billing-v2.” A later explicit acceptance updates the
context page, provenance, and retrieval index atomically. Deferring a proposal
does not undo plan completion; it remains visible and is included in future
context briefs when relevant.

Context updates should be small, reviewable, reversible, and fully sourced.

## 14. Plan-grounding invariant

Every executable plan must identify the Product Knowledge units that grounded
its objective, repository ownership, constraints, decisions, and acceptance.
Plan creation fails or asks a focused question when the agent cannot find
relevant context, finds contradictory accepted context, or cannot explain why
the selected context applies.

Before drafting a new plan, the agent also checks relevant pending context
proposals and stale-impact warnings from completed plans. It surfaces a
material unresolved impact as a plan risk or open question rather than treating
the old context as unquestionably current.

The plan stores references and a retrieval summary, not a full copy of every
context file. This keeps the plan readable while allowing later agents to
reproduce the grounding decision.
