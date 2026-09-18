# How workspace behavior is taught

A generated workspace teaches its agent through three layers, and which layer a
rule belongs in is decided by one question: **what will bring this rule back at
the moment it applies?**

## The three layers

**The entry instruction** is always loaded. It carries the gates — what needs a
person — and the standing prohibitions — what is never done without one. Nothing
here can rely on being retrieved, because the moments these rules fire are
moments nobody announces.

**Skills**, one per lifecycle stage, carry the procedure: how to write an
intent, how to prepare a worktree, how to run several plans at once. Each is
loaded when its stage comes due.

**Reference** carries what you look up once you already know you need it: the
command surface, the file contract, the record structures, worktree recovery
mechanics. It is the oldest of the three layers and the easiest to strand —
material shipped there reaches an agent only if something it loads names it.

One shipped document has a different reader. `how-it-works.md` explains the
division between person, agent, workspace, worktree, and executable to somebody
deciding whether to adopt the product. The entry instruction deliberately does
not name it. The agent already holds that division normatively, so a pointer
would spend always-loaded context on prose it never acts on, and leave two
accounts of the same rules free to drift — with the descriptive one reading as
though it explained the authoritative one. Its reader arrives through the
workspace README instead.

## Why a prohibition never becomes a skill

A skill is retrieved by an agent that has recognized its moment. That works for
procedure, because an agent about to write a plan knows it is about to write a
plan. It fails for prohibitions, because nothing will make an agent load the
do-not-push skill before pushing. So "commit, push, merge require explicit
authorization" and "never add agent attribution" stay resident, and a skill
never becomes the only place a gate is written.

The layering is enforced rather than trusted: every skill and every agent-facing
document is checked for reachability from the entry instruction, and the entry
instruction is checked for the gates that must not have migrated out of it. A
document written for a person is checked in both directions — reachable from the
README, and absent from the entry instruction — so the audience split is a
decision the suite holds rather than a convention someone remembers.

## What makes a procedure safe to defer

Two things, and either is enough.

A **user word** that reliably precedes it — "run these plans", "review this",
"initialize", "add a member". The person's own request is the trigger.

A **command that names the obligation in its own output** — `setup_required`,
`planning_required`, `reconcile_required`. The pointer then arrives at the
moment the work comes due rather than depending on a document read once at
session start. This is why those fields exist and why the procedures behind
those gates can live outside the entry instruction; see
[authorization boundaries](../domains/authorization.md).

A procedure with neither is left resident, or given a diagnostic that catches
the violation after the fact — which is what makes the knowledge-writing
conventions safe to defer, since `check` reports a note that crosses them.

## The host constraint

Only some hosts register skills by description. Claude Code and Cursor get a
thin routing file under their own skill directory pointing at the canonical body;
everywhere else, including Codex, a skill is reached because the entry
instruction names its path. So every skill costs a resident line, and the number
of skills is bounded by what that line budget can carry — the split is by
lifecycle stage rather than by topic for exactly this reason.

Owner:

- `context-circuit-source@product/AGENTS.md.in` — the always-loaded entry
  instruction: the gates, the prohibitions, and the pointer to each stage's skill.
- `context-circuit-source@product/skills/` — the canonical skill bodies and the
  per-host routing files that reach them.
- `context-circuit-source@product/docs/` — the reference layer, read on demand
  once an agent knows which question it has.
- `context-circuit-source@scripts/release-manifest.txt` — where each skill and
  doc is installed in a generated workspace.
- `context-circuit-source@internal/cli/documentation_test.go` — the reachability
  and placement checks that keep a layer from drifting.
