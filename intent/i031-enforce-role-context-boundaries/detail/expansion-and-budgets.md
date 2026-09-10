# Explicit expansion and budgets

## Expansion authority

An expanded scenario records the exact human prompt clause, affected role,
additional named surface, purpose, stopping condition, and revised budget. The
harness derives no expansion from what the agent later chose to read.

```yaml
mode: expanded
authority:
  prompt_clause: "Compare the PRD with the active API repository."
role: coordinator
additional_surfaces:
  - resource: source:billing-prd
    purpose: extract product requirements
  - resource: repository:api
    purpose: gather current implementation behavior
stopping_condition: compatible durable knowledge recorded
```

Authority is exact. One PRD does not mean all sources; one repository does not
mean every connected repository; diagnostics do not authorize changes; an
expansion for the coordinator does not automatically reach a planner or worker.

## Budgets

Every created role has hard ceilings for provided bytes, active-read bytes, and
file or resource count. Input-token ceilings are added when the host provides
trustworthy role-level usage.

Budgets belong to the scenario and normal request shape, not to a universal
constant. Results show absolute usage, the largest contributors, and change from
the scenario baseline. Exceeding a hard ceiling fails the context criterion;
material movement below it is diagnostic variation.

Token usage supplements resource evidence because tokenization and reporting
differ by host. A low token count cannot excuse a forbidden read, and unavailable
token reporting cannot erase observed delivered or active context.

