# Observation and grading

## Responsibility

Observation converts a real host run and its temporary workspace into normalized
evidence. Grading evaluates that evidence against the scenario contract and
produces a verdict with reasons. Neither component changes the workspace after
the session finishes.

## Evidence bundle

```text
observations/
  manifest.yaml              resolved immutable run inputs
  before/
    filesystem.json
    git.json
    lifecycle.json
  after/
    filesystem.json
    git.json
    lifecycle.json
  events/
    host.jsonl
    children.jsonl
    commands.jsonl
    context.jsonl
  transcript.md
  changes.patch
  assertions.json
result.yaml
```

Collectors take the before snapshot after fixture construction and the after
snapshot after the host reaches a terminal state. Snapshots include content
digests and normalized relative paths. Sensitive host configuration and private
runtime/provider payloads are excluded.

## Observation domains

The normalized model covers:

- filesystem additions, modifications, deletions, and unexpected writes;
- Git worktree state, branches, commits, authorship, and remotes;
- Context Circuit artifact state, including Product Knowledge, intents, plans,
  execution evidence, and their transitions;
- host session identity and terminal state;
- child-agent topology and roles;
- per-role context delivered, active reads, accessible surfaces, and explicit
  expansion authority;
- host-visible commands, tools, and accessed targets;
- final human-facing response;
- duration and usage when the host exposes them.

An unavailable observation is distinct from an observed negative fact.

## Assertion types

### Deterministic assertions

Deterministic assertions establish mechanical facts directly from evidence:

- required or forbidden path changes;
- schema and index validity;
- Git cleanliness, branches, and commits;
- lifecycle state transitions;
- source evidence remaining unchanged;
- target repository remaining unchanged before authorization;
- planner, worker, and verifier child topology;
- normal-request role context envelopes and zero unexpected role reads;
- declared byte or token ceilings for provided and actively read context;
- absence of writes outside the run root;
- recorded Context Circuit and host versions.

These assertions are authoritative. A semantic grader cannot override them.

### Semantic assertions

Semantic assertions are limited to human-facing meaning that cannot be reduced
to a stable mechanical fact, such as whether the final answer clearly explains
what was created and what the user should do next.

The semantic grader receives only the criterion, relevant response fragments,
and selected normalized evidence. It returns structured output:

```yaml
criterion: approval-is-next
verdict: pass
confidence: high
evidence:
  - "Approve this intent and I will inspect the repository and create the plan."
reason: The response identifies approval as the next action and does not claim execution began.
```

Semantic criteria avoid style preferences and exact wording. Each criterion
must describe one meaning a human can recognize.

## Verdicts

A run has one of four outcomes:

- **pass:** the real session completed successfully and every required assertion
  passed;
- **fail:** the session completed but one or more required assertions were
  contradicted by evidence;
- **inconclusive:** required evidence was unavailable or the host could not
  establish that this was the requested real-session class;
- **harness error:** workspace construction, collection, or grading failed before
  a trustworthy product verdict could be made.

The product under test is never blamed for a harness error. Failed and
inconclusive results preserve the workspace and evidence.

## Positive-journey grading

All foundation scenarios expect `pass`. Assertions describe the successful end
state and the boundaries respected on the way there. The suite does not count a
correct refusal, injected outage, expected conflict, or expected policy failure
as a positive pass case.

For example, the intent-creation journey can assert that the source repository
remained unchanged. This observes correct behavior within a successful intent
draft; it is not a negative scenario constructed to provoke premature execution.

## Reproducibility and comparison

Results are comparable by scenario revision, fixture digests, Context Circuit
artifact digest, host version, and declared model settings. Host and model
variance is retained as evidence rather than normalized away.

Repeated runs may produce different prose or valid implementation details. The
harness compares declared outcomes and normalized state, reporting non-contractual
differences separately as diagnostic variation.

The command-line report leads with the verdict, failed criteria, and preserved
workspace path. Machine-readable results support CI aggregation by scenario,
host, Context Circuit revision, duration, and assertion.
