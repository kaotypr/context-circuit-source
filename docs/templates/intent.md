# Intention — <id>

_Status: draft, waiting for your approval._

<!--
  This file is what the HUMAN reads. Keep it plain and short. Use exactly the five
  sections below — do not add out-of-scope, history, or assurance-rationale
  sections here. Keep technical terms as they are; do not rename them. Do not
  dump runtime commands, digest hashes, or engine invocations. The detailed,
  machine-checkable record lives in contract.yaml, which the human is not
  expected to open.

  `_Status:` is the human mirror. Keep the line; do not remove it.
  Draft: `draft, waiting for your approval.`
  After Gate 1: `approved.`
  After a feasible tracer with no intent-level questions:
  `approved, look complete, feasible.`
  contract.yaml remains the authority for approval identity (draft|approved).

  For a filled-in example of the voice and level of detail, see
  docs/templates/intent.example.md.
  A fuller by-concern write-up, when present, lives as separate files beside this
  file — not as another heading here.
-->

## Intention

<One short paragraph reflecting back what they want, then how it works in a few
plain lines. Prefer a small mermaid graph over prose for a flow. See
docs/templates/intent.example.md for the shape.>

## Expectations

- <What the human will have once every plan is done — the end result, one plain
  line each.>

## The plans

<Numbered plans. Each is a short bold title plus a one-line "_After this:_" saying
what changes once that plan lands. Mark an already-finished plan "— done.">

1. **<Plan title>**
   _After this:_ <what changes once this plan lands.>

## How carefully this is checked

**`<Explore | Standard | Critical>`**

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

<!--
  This section is phase-aware. At draft time, record every known unresolved human
  decision from the request and Product Knowledge. The line above does not claim
  that tracing cannot reveal a question. After tracing, revisit this section:
  - an intent-level question changes the goal, non-goals, constraints, acceptance,
    scope, tier, authority, or lifecycle; update the intent and re-enter Gate 1;
  - a plan-level implementation question belongs in the trace and plan instead;
  - a question already answered by the request is applied without asking again.
  When there is a genuine undecided question the human must settle, write it in bold.
  Once it is answered, keep the question and add an italic answer line beneath it;
  never delete an answered question:

  **<the question>**
  _Answer: <the decision>._
-->
