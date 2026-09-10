# Real-host scenario testing

This scope describes the first complete Agent Harness capability: create a
temporary project condition, start a real agent-host session, observe what the
agent does, and produce durable evidence about a successful user journey.

## Reading order

1. [System design](design.md) — normative architecture and fixed decisions.
2. [Scenario model](scenario-model.md) — scenario contracts and the initial
   positive journey catalog.
3. [Workspace factory](workspace-factory.md) — disposable project and repository
   construction.
4. [Role context boundaries](role-context-boundaries.md) — normal-request read
   envelopes, explicit expansion, and context-usage accounting.
5. [Host execution](host-execution.md) — real Codex, Claude Code, and Cursor
   session drivers.
6. [Observation and grading](observation-and-grading.md) — evidence capture,
   assertions, results, and reproduction.
