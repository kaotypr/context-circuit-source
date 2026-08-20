# Document-system evaluation and conformance

Status: implementation evidence; no rollout threshold or publication decision
has been accepted.

This document describes reproducible measurements for the integrated
document-system fixtures. It reports information architecture and reviewer
understanding, not just file or token counts. The fixture evidence is
structural: deterministic YAML and schema execution remains the approved
portable host capability.

## Evidence inputs

- `test/fixtures/document-system/conformance-matrix.yaml` is the end-to-end
  result matrix. Every case names its owning parent plan, fixture, and source
  contract.
- `test/fixtures/document-system/integration/fixture-attribution.yaml` prevents
  a joined report from relabeling or merging independently verified evidence.
- `test/fixtures/document-system/measurements/route-reads.yaml` records
  required files and lines for the entry, planning, gathering, execution,
  verification, resume, and review routes.
- `test/fixtures/document-system/measurements/repeated-material.yaml` records
  repeated authority claims and whether rationale, safety, and evidence remain.
- `test/fixtures/document-system/measurements/schema-deviations.yaml` records
  hard deviations separately from advisory diagnostics.
- `test/fixtures/document-system/measurements/resume-correctness.yaml` records
  structured, legacy, malformed, and contradictory resume outcomes.
- `test/fixtures/document-system/measurements/human-review.yaml` records
  whether a reviewer can identify purpose, status, authority, uncertainty, and
  the next decision from each prototype.
- `test/fixtures/document-system/package-boundary.yaml` records the approved
  distinction between the complete maintainer source package and the separate
  publishable artifact boundary.

## Measurement rules

For each route, count the mandatory files and source lines read before an agent
can safely continue. Count repeated material only when it repeats an authority
claim, lifecycle field, acceptance state, verification result, or rationale
that has a canonical owner. A shared link or a short navigation pointer is not
automatically duplication.

Schema deviations are counted by category:

- base OKF hard failures;
- Context Circuit profile hard failures;
- the applicable non-OKF artifact-schema failures;
- advisory diagnostics such as broken links, optional-index absence,
  freshness, metadata, copy, and retrieval concerns.

The categories are independent. An advisory cannot become a base OKF failure,
and a non-OKF runtime or plan failure cannot be relabeled as a knowledge
failure. A workflow stop caused by source contradiction or a human gate is
recorded as a separate authority observation.

Resume correctness requires all of the following:

1. use complete `handoff.yaml` when present;
2. use Markdown-only `handoff.md` only when structured state is absent;
3. reject malformed structured state without falling back to a conflicting
   Markdown summary;
4. preserve legacy history without rewriting it;
5. stop on contradictory evidence rather than reconstructing state.

Human-review comprehension is a hard safety constraint. A prototype fails its
measurement if a reviewer cannot identify all five fields: purpose, current
status, authority, uncertainty, and next decision. Review notes are
observations, not evidence that a human gate was satisfied.

## Before and after discipline

The before values are bounded fixture baselines, not claims about every
workspace. The after values describe the current prototypes. A route passes only
when it reads no more required context than its baseline and when a targeted
prototype reduces repeated authority material without dropping safety,
evidence, uncertainty, or rationale.

Token counts may be recorded directionally for comparison, but they are never a
release gate and never substitute for file/line counts, schema correctness,
resume correctness, or comprehension. A smaller number is invalid evidence when
it was achieved by omitting a required read, safety instruction, evidence
reference, uncertainty, rationale, or human decision.
A smaller artifact that loses one of those is a failed measurement.

No quantitative improvement threshold has been accepted yet. The human
`metric-threshold-decision` gate must choose a target after baseline review.
Until then, the measurements are evidence only and do not authorize migration,
status change, publication, or release.

## Package-layer evidence

The approved decision
`.runtime/plans/document-system-conformance-rollout/package-boundary-decision.yaml`
defines the maintainer source package as including schemas, templates, skills,
agents, fixtures, and attribution. The publishable template remains governed by
the existing release manifest and excludes `test/`, `docs/release.md`,
`.github/`, and `scripts/`. Acceptance checks measure source-package
completeness separately from publishable-boundary exclusions; neither check
publishes or changes the release boundary.

## Reproduction

From a clean source checkout:

```sh
git diff --check
sh -n test/acceptance.sh
sh test/acceptance.sh
```

The acceptance suite checks fixture attribution, independent result categories,
measurement invariants, release-package contents, and attribution text. It
does not install a parser, package manager, or command runtime.
