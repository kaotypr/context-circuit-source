# Intention — i029-establish-agent-harness-foundation

_Status: approved.

_Status: draft, waiting for your approval._

## Intention

Create **Agent Harness as a standalone product**, separate from Context Circuit,
with the core runner and disposable-workspace factory needed to exercise an
immutable Context Circuit release as an external user would.

Each run starts from a fully declared project condition, records exactly what it
was given, and can preserve or reproduce the resulting workspace without using
the Context Circuit maintainer checkout as its test environment.

## Expectations

- Agent Harness has its own product and repository boundary.
- A scenario can declare a complete starting condition from small reusable
  fixtures.
- Every run receives a unique disposable project, repositories, inputs, and
  observation area.
- The tested Context Circuit artifact and all resolved fixtures are identified
  immutably in the run record.
- Passing workspaces can be cleaned up, while questioned runs can be preserved
  and reproduced from recorded inputs.

## The plans

1. **Establish the standalone harness and workspace factory.**
   _After this:_ Agent Harness can materialize a declared test world around an
   immutable Context Circuit artifact and manage its full run lifecycle.

## How carefully this is checked

**`Standard`**

Explanations:
- **Explore:** you check it yourself as you work alongside the agent — no separate
  verification. Meant for trying things out.
- **Standard:** a second, independent agent verifies the finished work against your
  definition of done before it's called done.
- **Critical:** the strictest — independent verification plus a repair-and-recheck
  loop. For anything risky or impossible to undo (money, security, production, data
  you can't get back).

## Open questions

1. **Where should the standalone Agent Harness repository be created or connected,
   and which branch should its work start from?**
   _Answer: Use the standalone `kaotypr/agent-harness` repository at
   `/Users/kao/Workspace/kaotypr/repositories/agent-harness`, starting from
   `main`._
