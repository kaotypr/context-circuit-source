# Intention — <id>

_Status: draft, waiting for your approval._

<!--
  This file is what the HUMAN reads. Keep it plain, short, and free of jargon:
  no internal ids, file names, digests, branch/model names, or runtime commands,
  and no fancy terms. The detailed, machine-checkable record lives in
  contract.yaml, which the human is not expected to open. Use exactly the five
  sections below — do not add out-of-scope, history, or assurance-rationale
  sections here.

  For a filled-in example of the voice and level of detail, see
  docs/templates/intent.example.md.
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

None.

<!--
  When there is a genuine undecided question the human must settle, replace "None."
  with the question in bold. Once it is answered, keep the question and add an
  italic answer line beneath it; never delete an answered question:

  **<the question>**
  _Answer: <the decision>._
-->
