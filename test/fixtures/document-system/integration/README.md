# Integrated document-system fixtures

This fixture index joins the independently verified Plan 0016 foundation,
Plan 0017 knowledge, Plan 0018 plan/task, and Plan 0019 runtime evidence. Each
scenario keeps its owning plan and source contract visible in
`conformance-matrix.yaml`; the join does not redefine any parent authority.

The matrix reports these categories independently:

- `base_okf` — hard OKF v0.2 conformance for explicitly declared bundles;
- `context_circuit_profile` — the stricter knowledge-profile result;
- `non_okf_artifact_schema` — the owning schema for plans, tasks, runtime, or
  normative instruction artifacts;
- `advisory_diagnostics` — links, optional indexes, freshness, copy, and
  retrieval findings that never fail base OKF;
- `workflow_gate` — a separate observation of human approval, status, or
  contradiction handling.

The repository supplies structural evidence only. Deterministic YAML and schema
execution remains the approved portable host capability. A source contradiction
blocks continuation and is surfaced to a human; it is never repaired by
reconstructing missing facts.
