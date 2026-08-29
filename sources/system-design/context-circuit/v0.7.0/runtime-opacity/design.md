# Runtime opacity

A **delta on v0.6 / v0.6.1** ([../README.md](../README.md)). This scope changes
only the invoke-not-read boundary it names; everything else about the runtime and
the coordinator is unchanged.

## Capability

The coordinator treats `wrapper/runtime/engine.sh` as an **opaque tool**: it
invokes actions (`sh wrapper/runtime/engine.sh <action> <args>`) and consumes
their printed results, and it never opens the engine's source to read how an
action is implemented. Everything a coordinator needs to drive an action is
already carried by the invoking skill and the execution brief.

This is the coordinator-side corollary of **INV-RUNTIME-01**: because the runtime
is a small host-neutral deterministic library with no prompts, no Product
Knowledge interpretation, and no routing policy, there is nothing in it for the
coordinator to *interpret* — only actions to *call*. Reading it can therefore only
waste tokens and latency, tempt the coordinator to couple to internal mechanism,
and leak internals into a surface that is supposed to stay conversational.

## The problem

The discipline exists but is not enforced where it matters, and it is not obeyed
reliably.

- **It is worded as unnecessary, not forbidden.** `wrapper/adapters/AGENTS.md`
  says *"You never need to read its implementation to understand or execute a
  plan; the execution brief is enough."* That tells a model reading is
  *unnecessary* — it does not tell it reading is *not allowed*. A model that
  believes reading might help still may, without contradicting the text.
- **It is scattered, not owned.** The stronger phrasing *"invoke, never read the
  engine"* lives inside individual skills (`cc-execute`, `cc-run-stack`,
  `cc-publish`). There is no single normative statement of the boundary, so its
  strength varies by which skill is in play and it can drift between them.
- **It is enforced only by a maintainer test.** The only hard check is dimension C
  of the human-simulated harness (`template-harness/human/grade.sh`), which lists
  `wrapper/runtime/engine.sh` as a forbidden read. That gate runs in the source
  repo, not in a released workspace — so nothing in the shipped product prevents,
  detects, or discourages the read at run time.

### Observed failure (the trigger)

A live v0.6 run of the ten-plan run-stack case (`10-run-approved-stack`, driven
through the real headless coordinator) built and independently verified all ten
plans in correct dependency order, with every fan-in integration base correct —
**and still failed the grade**, solely because the coordinator issued one `Read`
of `wrapper/runtime/engine.sh` early in the run:

```
[FAIL] READ forbidden path (any point): wrapper/runtime/engine.sh
```

The behavior is **non-deterministic**: the same case has passed dimension C on
prior runs. That is the crux — an access boundary obeyed *most* of the time is not
a boundary. A released template must not depend on the model choosing not to read
the engine on a given run.

## Principles

- **Opacity is a boundary, not a suggestion.** The engine is a tool with a
  documented action surface; its implementation is out of scope for every product
  role that invokes it.
- **Remove the reason before adding a rule.** A coordinator reaches for the engine
  source when the skill did not tell it how to invoke an action it needs. The first
  fix is to close that information gap so there is never a motive to open the file;
  the prohibition is the backstop, not the primary mechanism.
- **Say "must not," not "need not."** Wherever the boundary is stated, it forbids
  the read; it does not merely call it unnecessary.
- **One owner, referenced everywhere.** The boundary is stated once, normatively,
  and every skill references it rather than restating a variant.
- **Enforcement should not depend on a maintainer-only test.** The released product
  should discourage or prevent the read on its own; the harness gate stays as
  regression proof, not as the only line of defense.

## Fixed decisions

1. **No new invariant.** This is a strengthening of the INV-RUNTIME-01 corollary
   and a wording/ownership cleanup, not a new contract rule. The canonical owner of
   runtime policy stays `wrapper/contracts/invariants.yaml`; the canonical prose
   owner of the boundary stays `wrapper/adapters/AGENTS.md`.
2. **The boundary is a prohibition.** The released wording states that a role which
   invokes the engine **must not read** `wrapper/runtime/engine.sh` (nor any
   runtime implementation file), and that the skill plus brief are sufficient by
   construction.
3. **Skills stay self-sufficient.** Every engine action a coordinator may need is
   documented at its point of use (as `cc-run-stack` already does — "everything the
   loop needs is below"), so no skill leaves a coordinator with a reason to open
   the engine. Any gap found is closed in the skill, not by permitting a read.
4. **The harness gate stays.** Dimension C keeps `wrapper/runtime/engine.sh` (and
   runtime implementation files generally) as a forbidden read, and this scope is
   the acceptance target: the fix is done when the live run-stack and execute cases
   pass dimension C reliably, not occasionally.

## Design options (the fix is deferred to a plan)

The mechanism is left to the implementing plan ("investigate + fix" — later). The
options, in preference order set by the "remove the reason, then forbid" principle:

- **A. Close the wording gap (required).** Promote the single normative statement
  in `wrapper/adapters/AGENTS.md` from "you never need to read" to an explicit
  "must not read the runtime implementation," and have each invoking skill
  (`cc-execute`, `cc-run-stack`, `cc-publish`) reference that one statement instead
  of carrying its own softer variant. This is the minimum and is always part of the
  fix.
- **B. Close the information gap (required where found).** Audit each invoking
  skill for any engine action whose usage is not fully specified at the call site,
  and document it inline, so a coordinator never lacks the invocation detail that
  tempts a read. `cc-run-stack` is the reference for the target completeness.
- **C. Host-level guardrail (optional, host-shaped).** Where a host supports it,
  deny-list `wrapper/runtime/engine.sh` for read tools (e.g. a Claude Code
  permission entry). This is a per-host convenience recorded as `host_evidence`
  only; it never becomes a route or a contract requirement (INV-HOST-01) and must
  not be the sole mechanism, since not every host can enforce it.

Investigation should first pull the coordinator transcript around the observed
read to confirm *why* it reached for the source (most likely an action whose
invocation it was unsure of — which would point at option B), then apply A + any B
gap it reveals, and leave C to the host adapter.

## Boundaries

- This scope governs only **reading the runtime implementation**. It changes
  nothing about how actions are invoked, what they do, or the run-stack /
  execute control flow.
- It adds no authority: it neither approves, executes, verifies, completes, nor
  delivers anything, and it introduces no new role or skill.
- It is orthogonal to the [execution-latency](../execution-latency/) scope. That
  scope makes run-stack overlap independent plans; this one keeps the engine opaque
  while it does so. Neither depends on the other.
