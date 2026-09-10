# Workspace and fixture model

## Run ownership

A run root contains the project presented to the agent, seeded repositories,
optional local remotes, resolved inputs, before-and-after observation areas, and
the final result. All paths recorded in portable evidence are relative to this
root.

```text
<run>/
  manifest.yaml
  project/
  repositories/
  remotes/
  inputs/
  observations/
  result.yaml
```

## Fixture kinds

The fixture registry supports independently reusable building blocks for:

- project identity;
- request-bounded source evidence;
- accepted Product Knowledge and its valid index;
- real Git repositories with deterministic histories;
- portable and local repository bindings;
- committed roster members and selected local identity;
- valid intent, plan, and execution prerequisites.

Dependencies between fixtures are explicit. Resolution rejects conflicting
destinations and writes the selected fixture revision and content digest into
the manifest. A scenario never relies on whatever a mutable fixture name means
after the run begins.

## Repository fixtures

Seed repositories remain small enough to understand but realistic enough to
require genuine reading. They contain application code, tests, project guidance,
dependencies, and fixed Git history. The initial set supports context gathering,
intent planning, and one executable Standard change.

Repository author, commit, branch, and timestamp inputs are stable where Git
supports control. Fixtures do not depend on or link back to the Context Circuit
source checkout.

## Construction lifecycle

Construction allocates the run root, materializes Context Circuit, initializes
the declared workspace state, installs evidence and knowledge, creates Git
fixtures, records local bindings and member identity, creates declared lifecycle
prerequisites, validates the result, and captures the baseline snapshot.

The factory returns a ready project only after every declared condition is valid.
It performs no host launch and no grading.

