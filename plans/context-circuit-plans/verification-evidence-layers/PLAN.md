# Acceptance evidence-layer enforcement

Status: draft
Repository: context-circuit-source
Source: sources/wrapper-action-cost-and-integrity-prd.md

## Review summary

Current plans map acceptance criteria to commands but do not state whether a
command proves schema, store, API, process, browser, or human behavior. This
plan adds that missing integrity check after engine-owned runtime evidence
exists. The comparison rule is intentionally exact: a criterion passes only
when its required layer equals the observed layer.

## What approval authorizes

Approval authorizes evidence-layer fields and compatibility rules, verifier
comparison of required and observed layers, compact evidence references, thin
guidance updates, and bounded runtime, host, security, upgrade, release, and
semantic fixtures.

Approval does not run product verification, invoke a browser or provider,
change implementation, authorize a waiver, deliver, publish, deploy, merge,
or clean runtime.

## Scope and non-goals

`plan.yaml` owns exact scope. The verifier stays read-only and cannot turn a
waiver into a pass. Existing completed plans and historical evidence remain
unchanged.

## Proposed solution

Add a canonical evidence-layer vocabulary to acceptance and verification
mappings. `plan.yaml` owns `required_layer` and `produced_layer`; tasks retain
only acceptance and verification IDs. Runtime records carry compact
`observed_layer`, outcome, and evidence references. The verifier uses exact
matching and reports passed, failed, blocked, or waived per criterion. Only
passed evidence can satisfy completion; waived is a human-visible,
non-passing outcome, not a new gate. Legacy unfinished work must add a mapping
before another verification run.

## Tasks and dependencies

| Task | Outcome | Depends on |
| --- | --- | --- |
| VEL-001 | Evidence layers and legacy compatibility are contract-owned | — |
| VEL-002 | Verifier rejects non-matching evidence and reports outcomes distinctly | VEL-001 |
| VEL-003 | Host, security, upgrade, release, and semantic behavior is proven | VEL-002 |

CC-002 is done. CC-009 is approved and is the implementation prerequisite;
CC-010 must not execute until CC-009 reaches done with completion evidence.

## Acceptance criteria

- VEL-AC-01 (`schema`): New criteria and verification mappings declare evidence layers.
- VEL-AC-02 (`process`): Non-matching observed evidence is rejected.
- VEL-AC-03 (`process`): Build or text checks cannot prove browser or human behavior.
- VEL-AC-04 (`process`): Missing capability never becomes a false pass.
- VEL-AC-05 (`process`): Independent verification remains read-only and non-repairing.
- VEL-AC-06 (`process`): Legacy history remains readable without invented evidence.
- VEL-AC-07 (`process`): Host, security, upgrade, and release boundaries remain valid.

## Verification

| Criterion | Required layer | Verification | Produced layer |
| --- | --- | --- | --- |
| VEL-AC-01 | schema | VEL-VT-01 | schema |
| VEL-AC-02, VEL-AC-03, VEL-AC-05 | process | VEL-VT-02 | process |
| VEL-AC-04 | process | VEL-VT-03 | process |
| VEL-AC-06 | process | VEL-VT-05 | process |
| VEL-AC-07 | process | VEL-VT-04, VEL-VT-06, VEL-VT-07 | process |

Each verification mapping produces one declared layer. A command that checks
multiple layers uses separate mappings; no layer is inferred from another.

## Risks and assumptions

- Exact matching is conservative and may report blocked when a human would
  consider one layer informative about another; that tradeoff keeps the rule
  deterministic and cheap to load.
- Browser tooling is outside this plan; the plan validates evidence claims,
  not the provider used to obtain them.
- A waived outcome remains non-passing and is recorded in the existing
  human-decision evidence, without adding a waiver gate.

## Delivery boundary

Commit, push, merge, publication, deployment, and cleanup remain unapproved.

## Provenance

Primary source: `sources/wrapper-action-cost-and-integrity-prd.md`, Could
requirement 14 and related acceptance criteria. Static host evidence from
CC-002 is treated as implemented groundwork. CC-009 is approved but remains an
execution prerequisite until done. Confirmation-token requirement 15
is intentionally excluded because CC-004 established the current exact
confirmation flow. The source PRD remains passive and is not Product Knowledge.
