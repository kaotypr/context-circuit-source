# Agent-context runtime fixtures

These fixtures exercise the Plan 0019 route and structured handoff contract.
They are runtime evidence, not OKF concepts, Product Knowledge, plan status, or
a replacement for the portable host validation capability.

The fixture set covers:

- explicit `entry`, `planning`, `gathering`, `execution`, `verification`,
  `resume`, and `review` read manifests;
- a missing-contract stop when a required read is absent;
- structured-only and structured-plus-summary handoffs;
- legacy Markdown-only handoff fallback without rewriting history;
- ownership, recovery, verifier isolation, stack, task evidence, and completion
  boundaries.

`handoff.yaml` is preferred whenever it exists. `handoff.md` is explanatory
only beside structured state and is a read-only legacy fallback when no YAML
handoff exists. Runtime YAML remains outside OKF and cannot override canonical
plan/task status or human gates.
