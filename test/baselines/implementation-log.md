# Context Circuit 1.0 implementation log

## FINAL-001 — baseline frozen

Status: complete.

Evidence: the pre-edit `sh test/acceptance.sh` run passed all existing
scenarios. The supplied static measurements are recorded in
`test/baselines/profile-baseline.yaml`; 60 normalized routing cases are in
`test/routing/fixtures.yaml`; the old assertions are mapped in
`test/baselines/assertion-ledger.yaml` and rule owners in
`test/baselines/ownership-map.yaml`.

Remaining work packages append evidence below. This log records observed test
results, implementation decisions, deviations, and unresolved risks; it is
not a lifecycle authority or plan approval record.

## FINAL-002 — wrapper/template and release foundation

Status: complete. `wrapper/manifest.yaml`, `wrapper/adapters/`, `template/`,
and the staged release assembler define source, shipped, mutable, and
preserved boundaries. The release suite proves an uninitialized artifact,
excludes maintainer/runtime/source state and old adapters, and reports dirty
source state without publishing.

## FINAL-003 — canonical router and bounded context

Status: complete. The owner catalogs, compact Tier-0 probe index, route/context
contracts, schemas, and `wrapper/runtime/engine.sh` implement one Stage A/Stage
B evaluator. 60 fixtures pass with one normalized action each. Tier 0 measures
5,932 bytes; the hard ceiling is 8,192 bytes.

## FINAL-004 — human plan experience

Status: complete. Plan/task v2 schemas, additive v1 migration, human `PLAN.md`
and task templates, Review Card, and separate approval/run/finish language are
shipped. Lifecycle tests prove approval requires confirmation, changes only
status/projections, and cannot execute.

## FINAL-005 — roles, receipts, and runtime

Status: complete. Thin coordinator/writer/verifier adapters, delegation and
receipt schemas, atomic leases, exclusive scope checks, handoffs, completion
and stack schemas are shipped. Runtime tests prove contention, verifier
read-only permissions, receipt validation, and resumable evidence.

## FINAL-006 — gates and operational fallback

Status: complete. Confirmation cards, append-only archive events, source
selection, takeover confirmation, cleanup inspection, provider-neutral fallback,
delivery documentation, and security boundaries are implemented and tested.

## FINAL-007 — compatibility, upgrade, rollback

Status: complete. Compatible, migration-needed, blocked, and legacy-unknown
classification plus additive v1→v2 migration, preserve boundaries, and
wrapper-only rollback helpers are shipped. No migration rewrites accepted
context, canonical plan intent/status, runtime evidence, dirty worktrees, or
repositories.

## FINAL-008 — semantic acceptance suites

Status: complete. `test/acceptance.sh` is a small suite runner replacing the
former monolithic prose harness. Contracts, routing, budgets, lifecycle,
ownership, runtime, gates, recovery, stacks, upgrades, security, behavior
matrix, and release suites all pass. The behavior matrix records 40 positive /
negative obligations and routing coverage is 60 cases.

## FINAL-009 — final A/B and release evidence

Status: complete with live-token publication waiver. Measured selected context:

| Profile | Measured bytes | Hard ceiling |
| --- | ---: | ---: |
| Tier 0 | 5,932 | 8,192 |
| Orientation | 6,521 | 20,480 |
| Root run-plan | 5,359 | 22,528 |
| Writer | 3,374 | 12,288 |
| Verifier | 1,503 | 11,264 |
| L1 root + writer + verifier | 10,236 | 46,080 |
| L2 three-task resume estimate | 17,311 | 58,368 |
| Resume delta | 5,116 | 12,288 |

L1 deterministic reduction is 93% against the frozen 154,368-byte baseline.
Live host-token usage was unavailable; the explicit publication waiver and
retained deterministic gates are in `token-evidence.yaml`. Release staging
passed and no external publication or deployment was performed.

## Release artifact layout correction

Status: complete. The source-only `template/` directory is no longer copied
into the generated workspace. Its mutable seed is overlaid at the artifact
root, matching DESIGN-SPEC §4.2. The release manifest and acceptance suite now assert
that the nested directory is absent while root `workspace.yaml`, `context/`,
`sources/`, and `plans/` remain present.

## Source design and adapter naming cleanup

Status: complete. Comparative design drafts and the active specification live
under `sources/context-circuit-design/`, with the active file named
`DESIGN-SPEC.md`. The redesigned host adapters use the `cc-*` discovery prefix;
their new contract-driven behavior is unchanged and the obsolete skill set is
not restored.

## Maintainer self-hosting approval guard

Status: implemented and tested. Self-hosting initially exposed that approval
made the maintainer checkout dirty, after which the clean-base worktree guard
correctly blocked execution. The runtime now classifies only the exact
`product-source` approval projection as `MAINTAINER_APPROVAL_COMMIT_REQUIRED`,
routes it through a focused human commit card, and continues to reject staged,
untracked, unrelated, or normal product-repository changes. The source safety
instructions now allow a commit only when the maintainer explicitly requests
it; commits remain separate from approval and are never automatic.

Evidence: ownership fixtures cover the exact projection and unrelated dirty
file failure; routing coverage includes the maintainer commit card; the full
semantic acceptance result is recorded after this change.

## Uninitialized workspace write guard

Status: implemented and tested. A fresh released workspace now carries an
explicit `instantiated-workspace` identity with `status: uninitialized`.
Router calls inspect that state and force every non-read-only request—including
generic requests such as “help me build this”—to initialization with the
`identity-acceptance` gate. Read-only orientation remains available, and an
accepted workspace retains normal plan routing. Entry-skill guidance and the
route precedence contract now forbid implementation, plan, task, runtime, or
repository writes before identity acceptance.

Evidence: routing tests cover uninitialized write, uninitialized read-only,
and accepted-workspace paths; the complete semantic acceptance suite passes.

The follow-up screenshot review exposed a second bypass after identity
acceptance: a generic build request could still be interpreted as direct
implementation. The accepted-workspace router now sends unnamed write-like
requests to `draft-plan`; `cc-entry` and `cc-plan` explicitly prohibit
implementation until plan approval and the separate execution trigger.

## Versioned dist artifact naming

Status: implemented and tested. The convenience build now defaults to the
current product version, `v0.5.0`, and produces
`dist/context-circuit-v0.5.0` with a matching archive. Alternate versions can
be supplied explicitly to the builder; the unversioned `preview` artifact name
is no longer generated by the default command.

Evidence: the release acceptance test passes with the versioned directory and
archive assertions, and a repository-level build was verified after assembly.

## Known deviations and risks

- The repository implements host-neutral filesystem primitives and contracts;
  it does not launch external agents or provide a user-facing CLI.
- Live-token reduction is not measured in this offline checkout; publication
  requires the recorded waiver or later host evidence.
- The shell engine intentionally validates required record fields rather than
  being a general YAML parser. Hosts must preserve the schema contracts and
  semantic suites.
