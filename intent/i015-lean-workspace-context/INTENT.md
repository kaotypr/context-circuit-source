# Intention — i015

_Status: draft, waiting for your approval._

## Intention

What you want: **keep the current Context Circuit agent loop — it is already about right — and remove the waste around it**, so a workspace agent spends context on the work, not on hunting and rereading.

The gates, the skill-as-procedure path, and “read the catalog then only the matching knowledge” stay. What changes is the waste that does not change correctness:

- **The catalog actually selects.** An agent matches the request to named knowledge units from the catalog. It does not grep or list the knowledge tree to discover them.
- **Each procedure fact has one home.** Shared instructions, workflow, coordinator, and skills stop retelling the same lifecycle. They point at the owner.
- **A skill names files and invokes.** It does not restate the contract, and the invoke line matches the runtime, so the agent does not search tests or open the engine.

A fuller write-up by topic sits beside this page. You still approve once.

```mermaid
flowchart TD
  A["Agent starts a workspace action"] --> B["Same loop as today:<br/>catalog, named skill, your gate"]
  B --> C["Catalog names the knowledge to open"]
  B --> D["One home per procedure fact"]
  B --> E["Skill names files and invokes"]
  C --> F["The work, with less hunting"]
  D --> F
  E --> F
```

## Expectations

- Ordinary workspace work still follows today’s loop: catalog, skill, tracer after approval, independent check on Standard and Critical.
- An agent can pick the right knowledge from the catalog without searching the knowledge tree.
- The same lifecycle fact is not retold in shared instructions, workflow, coordinator, and every skill.
- Each skill names the next files and the runtime action; a wrong or restated invoke line is not something the agent has to discover.
- Safety spine, human gates, and one-owner-per-rule are not dropped or copied into skills to save a read.
- Knowledge pages a typical route does not open are not rewritten.
- This is what an instantiated Context Circuit workspace gets.

## The plans

1. **Make the catalog actually select.**
   _After this:_ the catalog carries per-unit summaries, topics, and aliases; a typical route opens the named units and does not grep or list the knowledge tree to find them.
2. **Give each procedure fact one home.**
   _After this:_ shared instructions, workflow, coordinator, and skills point at the owner of a fact instead of restating it.
3. **Make skills name files and invokes.**
   _After this:_ a skill is a short procedure with the files to open and runtime actions that match the runtime; the agent does not open a schema or search tests to learn the next step.
4. **Prove it is an improvement, not a thinner wrong loop.**
   _After this:_ checks fail if the catalog still requires a knowledge-tree search, if a skill restates a contract or documents a wrong invoke, or if a gate or independent check was dropped.

## How carefully this is checked

**`Standard`**

This changes the instruction surface every workspace agent reads. The point is to keep today’s correct loop, so a second agent should confirm the waste is gone and the gates are not.

Explanations:
- **Explore:** you check it yourself as you work alongside the agent — no separate
  verification. Meant for trying things out.
- **Standard:** a second, independent agent verifies the finished work against your
  definition of done before it's called done.
- **Critical:** the strictest — independent verification plus a repair-and-recheck
  loop. For anything risky or impossible to undo (money, security, production, data
  you can't get back).

## Open questions

**Should Product Knowledge pages themselves be shortened, or is the win mainly cutting always-on instructions, naming files instead of grepping, and dropping extra hops?**
_Answer: Do not rewrite knowledge pages. The win is a catalog that selects, one home per procedure fact, and skills that name files and invokes. A page a typical route still has to open may be left as it is._
