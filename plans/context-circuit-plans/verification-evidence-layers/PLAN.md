# Acceptance evidence-layer enforcement

Status: draft
Repository: context-circuit-source
Source: sources/wrapper-action-cost-and-integrity-prd.md

## Review summary

Current plans map acceptance criteria to commands but do not state whether a
command proves schema, store, API, process, browser, or human behavior. This
plan adds that missing integrity check after engine-owned runtime evidence
exists.

## What approval authorizes

Approval authorizes evidence-layer fields and compatibility rules, verifier
comparison of required and observed layers, explicit limitation/waiver
evidence, thin guidance updates, and bounded runtime, host, security, upgrade,
release, and semantic fixtures.

Approval does not run product verification, invoke a browser or provider,
change implementation, approve a waiver, deliver, publish, deploy, merge, or
clean runtime.

## Scope and non-goals

`plan.yaml` owns exact scope. The verifier stays read-only and cannot turn a
waiver into a pass. Existing completed plans and historical evidence remain
unchanged.

## Proposed solution

Add a canonical evidence-layer vocabulary to acceptance and verification
mappings. Carry required and observed layers through task, delegation, handoff,
and completion records. The verifier compares them and reports pass, fail,
blocked, or waived per criterion. Host capability is bounded evidence only.
Legacy unfinished work must add a mapping before another verification run.

## Tasks and dependencies

| Task | Outcome | Depends on |
| --- | --- | --- |
| VEL-001 | Evidence layers and legacy compatibility are contract-owned | — |
| VEL-002 | Verifier rejects weaker evidence and reports outcomes distinctly | VEL-001 |
| VEL-003 | Host, security, upgrade, release, and semantic behavior is proven | VEL-002 |

The completed multi-host plan and draft engine-owned runtime-record plan are
prerequisites.

## Acceptance criteria

- VEL-AC-01: New criteria and verification mappings declare evidence layers.
- VEL-AC-02: Weaker observed evidence is rejected.
- VEL-AC-03: Build or text checks cannot prove browser or human behavior.
- VEL-AC-04: Missing capability and waivers never become false passes.
- VEL-AC-05: Independent verification remains read-only and non-repairing.
- VEL-AC-06: Legacy history remains readable without invented evidence.
- VEL-AC-07: Host, security, upgrade, and release boundaries remain valid.

## Verification

VEL-VT-01 through VEL-VT-07 cover contracts, runtime comparison, host
limitations, security, upgrades, release assembly, and complete acceptance.

## Risks, assumptions, and open decisions

- Evidence layers need an explicit ordering or compatibility matrix; a simple
  numeric hierarchy may be wrong when human and browser evidence are
  incomparable.
- Waiver approval ownership must be defined without turning host capability or
  verifier judgment into authorization.
- Browser tooling is outside this plan; the plan validates evidence claims,
  not the provider used to obtain them.

## Delivery boundary

Commit, push, merge, publication, deployment, and cleanup remain unapproved.

## Provenance

Primary source: `sources/wrapper-action-cost-and-integrity-prd.md`, Could
requirement 14 and related acceptance criteria. Static host evidence from
CC-002 is treated as implemented groundwork. Confirmation-token requirement 15
is intentionally excluded because CC-004 established the current exact
confirmation flow. The source PRD remains passive and is not Product Knowledge.
