# Intention — i031-enforce-role-context-boundaries

_Status: approved, look complete, feasible.

_Status: approved.

_Status: draft, waiting for your approval._

## Intention

Make **role-specific context usage a strict Agent Harness test outcome**. During
a normal request, the coordinator, planner, worker, and verifier should receive
and read only the context required by their distinct responsibility.

Different or broader reading is valid only when the human prompt explicitly asks
for it. The harness will distinguish what a role could access, what it was given,
what it actively read, and what the host charged as input usage.

## Expectations

- Every scenario declares a normal context envelope for each role it creates.
- Context delivered to a role and context actively read by that role are tracked
  separately.
- Unexpected or forbidden reads prevent a strict scenario from passing.
- Broader context is allowed only when the exact human request and added scope are
  recorded.
- File, byte, and available token usage remain within role-specific scenario
  budgets.
- A host that cannot prove a required context measurement reports it as
  inconclusive rather than pretending no context was used.

## The plans

1. **Enforce minimal context for every agent role.**
   _After this:_ successful harness runs prove that each role stayed inside its
   normal or explicitly expanded context boundary and budget.

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

_At draft time: no known unresolved human decisions._

