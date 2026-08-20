# Measurement fixtures

These files are reproducible, bounded measurements for Plan 0020. They compare
route reads, repeated authority material, schema deviations, resume outcomes,
and human-review comprehension while preserving safety and evidence.

The values are directional fixture baselines, not a universal benchmark. A
measurement is invalid when a smaller count comes from omitting required
instructions, evidence, uncertainty, rationale, or a human-gate decision.
Token counts are optional context and never the sole threshold.

Files:

- `route-reads.yaml` — required files and lines per route;
- `repeated-material.yaml` — canonical claim repetition and retained rationale;
- `schema-deviations.yaml` — independent hard and advisory result categories;
- `resume-correctness.yaml` — structured, legacy, malformed, and contradiction
  resume outcomes;
- `human-review.yaml` — purpose, status, authority, uncertainty, and next
  decision comprehension observations.

The human `metric-threshold-decision` gate remains unsatisfied until baseline
data is reviewed. These fixtures do not approve migration, change plan status,
publish, or release an artifact.
