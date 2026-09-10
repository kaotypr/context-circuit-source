# Scenario contracts

## Required declaration

Every scenario declares:

- stable identity, title, family, and schema version;
- complete workspace fixture composition;
- immutable Context Circuit artifact selection;
- human prompt content;
- required real-host capabilities and limits;
- expected roles and native child relationships;
- normal or explicitly expanded role-context envelopes and budgets;
- observable successful lifecycle and artifact outcomes;
- deterministic and semantic assertions;
- retention behavior.

```yaml
schema_version: 1
id: established-project-create-intent
family: intent-planning-execution
test_class: real-host
workspace:
  fixture: established-reporting-project
prompt:
  fixture: create-csv-export-intent
host:
  required_capabilities: [filesystem, shell]
context_policy:
  mode: normal
expect:
  roles: [coordinator]
  journey: intent-drafted
  repository_unchanged: true
retention:
  on_pass: result-only
  on_questioned: complete-run
```

Fixture names resolve before launch and their digests become part of the run
manifest. The prompt contains no hidden grading instructions. A changed fixture,
prompt, artifact, host version, or model setting remains visible in comparison.

## Independence

Later-lifecycle scenarios receive valid intent, plan, or execution prerequisites
as fixtures. They do not consume artifacts created by an earlier scenario run.
This keeps scheduling parallelizable and makes a single scenario reproducible.

## Positive terminal condition

Each scenario describes an achievable request against a valid starting state.
Its expected terminal condition is successful completion of that user journey.
Safety and context boundaries are asserted along the successful path, not by
constructing a request whose expected result is refusal or failure.

