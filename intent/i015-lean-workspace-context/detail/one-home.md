# One home per procedure fact

Continues from [design.md](design.md). This is the reread win.

## The waste

Always-on files are not too many. They **repeat**. Shared instructions,
workflow, coordinator, and each skill each retell:

- the two human gates
- tracer after approval, then derived plans
- one owner per rule
- invoke the runtime, do not read it
- host-blocked when a required child is missing
- `sources/` is passive
- archive is outside normal context

An agent that follows “read shared instructions, then coordinator, then
workflow, then the skill” pays for the novel four times before doing the
route. That repetition is not extra safety. The owner already has the fact.

## What “one home” means

For each of those facts, **one file owns the prose**. Every other
agent-facing file that needs it uses a **pointer**: the fact’s name and the
path, in one or two lines, not a restatement.

A pointer is valid only if an agent that must apply the fact is told **when
to open the owner**. If dropping the restatement would leave the agent with
no path to the fact, that is a downgrade — put the pointer in, not the essay.

## Suggested homes (shape, not a file map the tracer will refine)

These are the concerns, not an enforced path list. The tracer names the real
files after approval. Those files already live in the nested home
(`.context-circuit/wrapper` adapters, `.context-circuit/agents`,
`.context-circuit/docs`) and at `.agents/skills/` for product skills. This
intent does not move them.

| Fact | Home (kind) | Who points |
| --- | --- | --- |
| Authority order and “you say → action” | Workflow | Shared instructions, coordinator |
| Safety spine (what never to do) | Shared instructions | Coordinator, skills |
| Conversation routing and reporting language | Coordinator | Shared instructions (already send the agent there) |
| Per-route procedure | The named skill | Coordinator (already names skills by path) |
| One owner per rule | Invariant catalog | Shared instructions (already point there); skills do not restate |
| Runtime invoke-not-read | Shared instructions runtime section | Each skill that invokes: one pointer, not a restated corollary |
| Intent / trace / plan / execute / deliver mechanics | The skill that owns that route, plus the invariant owner | Other skills do not retell a neighbor’s route |

Domain pages are **not** in this table. They own domain knowledge for humans
and for retrieval after the catalog selects them. This intent does not rewrite
them to stop retelling the lifecycle.

## What a pointer looks like

Enough to act:

- “Gates and conversational actions: workflow.”
- “Do not read the runtime; invoke it. See shared instructions, Runtime.”

Not enough (restatement):

- a new paragraph that explains Gate 1, the tracer, feasibility, and “no
  second plan approval” inside a skill whose job is archive or verify.

## Always-on stays a spine

Shared instructions still:

- say this is a Context Circuit workspace
- send the coordinator to the coordinator file before the first reply
- list the safety spine as **pointers plus the few rules that file owns**
- name the first reads: workflow, workspace identity, catalog
- forbid reading the runtime

They do **not** need a second “how a planned request flows” that duplicates
workflow, or a second intent-front-door that duplicates the intent skill.

Coordinator still owns routing and how to talk to the human. It should **map**
requests to skills, not paste the intent skill into the coordinator file.

Workflow still owns authority order and the action table. That table is the
human-facing map. Skills should not reprint it.

## What this is not

- Not a line-count quota on always-on files.
- Not deleting the safety spine.
- Not moving invariant ownership into a skill.
- Not a rewrite of Product Knowledge domain pages.

## Edge cases

- **Host adapter extras** (Claude slash-skill symlinks, Codex spawn flags)
  stay on the host adapter. They are host facts, not lifecycle restatements.
- **A skill that must mention a neighbor route** (intent mentions tracer)
  points at the tracer skill; it does not paste tracer procedure.
- **If two files both look like the owner**, that is a bug this intent
  fixes: pick one home, convert the other to a pointer. Do not leave two
  essays “in case one is missed.”
