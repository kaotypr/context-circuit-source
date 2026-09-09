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

Later work added more facts that now get the same treatment: durable-only
writing, member identity and band-scoped ids, and host-native trees as routes.
A host-native stub that copies policy is a fifth telling, not a host
convention.

An agent that follows “read shared instructions, then coordinator, then
workflow, then the skill” pays for the novel several times before doing the
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
`.context-circuit/docs`), at `.agents/skills/` for product skills, and as
thin routes under `.claude/`, `.codex/`, and `.cursor/`. This intent does
not move them.

| Fact | Home (kind) | Who points |
| --- | --- | --- |
| Authority order and “you say → action” | Workflow | Shared instructions, coordinator |
| Safety spine (what never to do) | Shared instructions | Coordinator, skills |
| Conversation routing and reporting language | Coordinator | Shared instructions (already send the agent there) |
| Per-route procedure | The named skill | Coordinator (already names skills by path) |
| One owner per rule | Invariant catalog | Shared instructions (already point there); skills do not restate |
| Runtime invoke-not-read | Shared instructions runtime section | Each skill that invokes: one pointer, not a restated corollary |
| Intent / trace / plan / execute / deliver mechanics | The skill that owns that route, plus the invariant owner | Other skills do not retell a neighbor’s route |
| Durable-only writing | Knowledge invariant / product-knowledge | Gather and complete skills point; they do not paste the citation rule |
| Member identity and band allocation | Member invariant; workspace skill for the one-time roster choice | Intent and plan skills invoke allocate and point; they do not retell band mechanics |
| Host-native discovery | Host adapters / host instruction files | Native stubs name the owner; they do not copy role bodies or policy |

Domain pages are **not** in this table. They own domain knowledge for humans
and for retrieval after the catalog selects them. This intent does not rewrite
them to stop retelling the lifecycle.

## What a pointer looks like

Enough to act:

- “Gates and conversational actions: workflow.”
- “Do not read the runtime; invoke it. See shared instructions, Runtime.”
- “Commit messages: read INV-COMMIT-01. Do not restate it here.”

Not enough (restatement):

- a new paragraph that explains Gate 1, the tracer, feasibility, and “no
  second plan approval” inside a skill whose job is archive or verify.
- a host-native rule file that pastes the lifecycle instead of naming the
  owner.

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

Host-native stubs stay the host’s discovery format. A Cursor rule that already
points at an invariant is the pattern; a stub that retells Gate 1 is not.

## What this is not

- Not a line-count quota on always-on files.
- Not deleting the safety spine.
- Not moving invariant ownership into a skill.
- Not a rewrite of Product Knowledge domain pages.
- Not a redesign of host-native folders, nested home, durable-only writing,
  or band allocation.

## Edge cases

- **Host adapter extras** (Claude slash-skill symlinks, Codex spawn flags,
  Cursor spawn-model rules) stay on the host adapter. They are host facts,
  not lifecycle restatements.
- **A skill that must mention a neighbor route** (intent mentions tracer)
  points at the tracer skill; it does not paste tracer procedure.
- **If two files both look like the owner**, that is a bug this intent
  fixes: pick one home, convert the other to a pointer. Do not leave two
  essays “in case one is missed.”
- **New facts from later work** (durable-only, member identity, native
  trees as routes) are in this table so they get the same pointer treatment.
  Their owners stay the owners they already are.
