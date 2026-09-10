# Grading and results

## Evidence model

The run captures before-and-after filesystem, Git, and lifecycle snapshots;
host, child, command, and context events; the human-visible transcript; changes;
and resolved input identities. Relative paths and stable identifiers support
comparison without exposing unrelated machine state.

## Assertion order

Deterministic assertions first establish required artifacts, allowed changes,
schema validity, lifecycle transitions, Git outcomes, native child topology,
role-context boundaries, budgets, and external-write boundaries.

Semantic assertions then evaluate only declared human-facing meanings, such as
whether the result clearly states what happened and the next human action. Each
semantic verdict cites the response evidence it used. It cannot convert a failed
deterministic assertion into a pass.

## Result classes

- **Pass:** the real positive journey completed and every required assertion
  passed.
- **Fail:** the real session completed but evidence contradicts a required
  outcome.
- **Inconclusive:** required host or role evidence was unavailable or could not
  be attributed honestly.
- **Harness error:** setup, observation, or grading failed before a trustworthy
  product verdict could be reached.

A harness error is never blamed on Context Circuit. Fail, inconclusive, harness
error, interruption, and explicit keep all preserve the complete run.

## Comparison and reporting

Reports lead with the result, failed or inconclusive criteria, evidence links,
and preserved workspace. Aggregation groups by scenario revision, fixture and
Context Circuit artifact identity, host and version, model settings, duration,
usage, role context, and assertion outcome.

Repeated prose and valid implementation differences appear as diagnostic
variation unless they violate the declared outcome. Reproduction always starts a
new session from recorded inputs.

