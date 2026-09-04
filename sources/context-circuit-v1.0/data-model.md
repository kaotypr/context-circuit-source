# Data model — new schemas and on-disk layout

`invariant-deltas.md` references new schemas; this file specifies them and shows how
they relate to the existing ones. All are additive; existing schemas gain at most one
field. Schemas follow the current style (a controlled subset the engine reads/writes;
see `wrapper/contracts/schemas/`).

## New: `intent-contract.yaml` (owner of INV-INTENT-01)

Lives at `intent/<id>/contract.yaml`.

| Field | Type | Notes |
| --- | --- | --- |
| `schema_version` | int | 2 |
| `intent` | id | `i<NNNN>-<slug>` (e.g. `i0007-checkout-retries`) — distinct from plan ids; its own never-reused sequence |
| `title` | string | human title |
| `goal` | text | one paragraph |
| `non_goals` | list<text> | explicit exclusions |
| `constraints` | list<text> | e.g. no new deps, latency budget |
| `acceptance_criteria` | list<text> | **OUTCOME level** — what must be true, in terms a human can approve; not executable and not frozen upfront. The runnable check that proves each one (a grep, a test id) is a **discovery** output, carried in the plan, never authored here (`discovery-and-grounding.md`) |
| `scope` | object | the **envelope** — COARSE on purpose: candidate repositories + a rough boundary, `repositories: [{id, paths[]}]`. Checked against discovery's findings and against each plan |
| `tier` | enum | `explore \| standard \| critical` — **provisional**; discovery may raise it, and the raise surfaces at plan review |
| `status` | enum | `draft \| approved` |
| `intent_dependencies` | list | optional; `{id, reason}` entries naming other intents this one depends on (mirrors `plan_dependencies`, INV-PLAN-05) |
| `contract_digest` | digest | set on approval; the frozen identity of the decision |

`INTENT.md` (human-facing) sits beside it — the plain-language view of the same
decision, reviewable as one thing.

## New: `discovery-manifest.yaml` (grounding evidence)

Recorded per intent, one per repository in the envelope — discovery's report,
captured as **durable grounding evidence** rather than lost in a context window
(`discovery-and-grounding.md`). Lives at `intent/<id>/discovery/<repo>.yaml`.

| Field | Type | Notes |
| --- | --- | --- |
| `schema_version` | int | 1 |
| `intent` | id | the parent intent |
| `repository` | id | which repository this manifest covers |
| `file_map` | list | exact paths, the specific call sites that change, relevant signatures and integration points |
| `task_partition` | list | proposed split of the work, for the coordinator to ratify or adjust |
| `risks` | list<text> | concrete risks named against the real code (data, security, irreversibility, coupling) |
| `done_checks` | list<check> | executable commands/tests that prove the intent's outcome criteria hold; carried into the plan and re-run by the verifier |
| `tier_signal` | text | evidence that may raise the intent's provisional `tier` |
| `open_questions` | list<text> | anything the human must settle, for the coordinator to relay |
| `completeness_proof` | object | present only for a "change every X" obligation — the command and count showing the found set is the whole set |
| `recorded_at` | timestamp | |
| `freshness_checked_at` | timestamp | last bounded freshness check against current code; re-reads only what drifted |

Reused across sessions exactly as other grounding evidence is (INV-GROUND-*): a later
session loads the recorded manifest instead of re-discovering from zero.

## New: `candidate.yaml` (owner of INV-CANDIDATE-01)

A candidate is an identity computed over records the engine already keeps, plus the
contract digest. Recorded under the execution dir.

| Field | Type | Notes |
| --- | --- | --- |
| `schema_version` | int | 1 |
| `candidate_id` | digest | `cand-<hash>` over the inputs below |
| `repositories` | map | `repo -> latest_commit` (from the repository records) |
| `bases` | map | `repo -> base_commit` (from the repository records) |
| `contract_digest` | digest | the intent's frozen criteria (or `legacy` for pre-intent executions) |
| `computed_at` | timestamp | |

The digest is deterministic: identical inputs → identical `candidate_id`. Any change
to a commit map, a base, or the contract digest yields a new `candidate_id`, which is
what voids prior evidence (INV-CANDIDATE-01).

## New: `human-acceptance.yaml` (co-owned with candidate)

First-class acceptance, bound to a candidate.

| Field | Type | Notes |
| --- | --- | --- |
| `candidate_id` | ref | the accepted candidate |
| `accepted_by` | string | the human (name/handle), for team visibility |
| `checklist` | list | which criteria the human confirmed, incl. any `manual` ones |
| `accepted_at` | timestamp | |

A new candidate voids a prior acceptance — the acceptance names a `candidate_id`, so
it simply no longer matches once the candidate changes.

## Changed existing schemas (one field each)

- `plan.yaml` (+`intent: <id>` required; INV-PLAN-01 reworked). Status becomes a
  projection (see `lifecycle-and-gates.md`).
- `verifier-result.yaml` (+`candidate_id` the result is bound to).
- `context-impact.yaml` (+ the reconciliation-debt marker: `candidate_id`,
  `resolved: pending | reconciled | deferred`).

Everything else in the existing schemas is unchanged.

## On-disk layout additions

```
intent/                          # NEW — parallel to plans/
  i0007-checkout-retries/
    INTENT.md
    contract.yaml
    discovery/                   # NEW — durable manifests, one per repository
      checkout-service.yaml
  INDEX.md                       # active intents (mirrors plans/INDEX.md)
  archive/                       # archived intents (status-blind move, like plans)

plans/0012-.../plan.yaml         # gains `intent: i0007-checkout-retries`

.runtime/executions/<plan>/<exec>/
    candidate.yaml               # NEW — the current candidate identity
    human-acceptance.yaml        # NEW — candidate-bound acceptance
    context-impact.yaml          # gains the reconciliation-debt marker
    ...                          # attempts/, repositories/, snapshot/ unchanged
.runtime/knowledge-debt/         # NEW — markers for delivered-but-unreconciled candidates
```

The `intent/` tree mirrors `plans/` conventions exactly (stable ids never reused,
archive as a status-blind move, an INDEX catalog) so it reuses the existing id
allocation, index, and archive machinery rather than inventing new ones.

## Relationship to Product Knowledge (the boundary)

An **intent** is *the decision for one change* — request-scoped, frozen at approval,
and archived when the change is delivered. **Product Knowledge** (`context/`) is
*durable, cross-change truth*. They are different stores and stay so:

- An intent's criteria are **not** copied into Product Knowledge; they describe one
  change, not a durable fact (mirrors INV-GROUND-01's "referenced, never captured").
- After delivery, reconciliation (M4) may propose durable knowledge units *derived
  from* what the change taught — through the existing proposal/acceptance path
  (INV-KNOWLEDGE-02) — but the intent itself is not a knowledge unit.
- So the flow is `intent (decision) → change → reconciliation → Product Knowledge
  (durable)`, with the human acceptance gate on knowledge unchanged.
