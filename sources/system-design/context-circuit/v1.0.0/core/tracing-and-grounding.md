# Tracing and grounding — how a plan earns its detail

Part of Mechanism 1 (`intent-and-criteria.md`). Tracing is the phase that reads the
real codebase and turns an approved plain intent into a detailed, grounded plan. It is
the answer to "where does the agent learn, in detail, what it is going to build?"

## The problem, from real sessions

Two observed sessions frame this:

- **A design-first session** produced a detailed multi-plan change — but only because
  the human had already driven a trace/design phase (an `Explore` child read the
  code four times; a system design was written) *before* asking for plans. The detail
  came from that phase, funnelled through `context/` and `sources/`. The plan step
  itself read no code.
- **A direct session** ("hash everything in localStorage; route every access through a
  hook") produced a detailed nine-task plan with *no* prior design phase — because the
  coordinator improvised tracing inline: it grepped every `localStorage` call-site,
  grouped them by feature area, checked import aliases. It worked, but it was ad-hoc,
  it bloated the coordinator's own context, it captured nothing durable, and its file
  count *drifted mid-plan* — evidence the tracing was not authoritative.

The lesson: detailed plans already depend on a tracing read of the real code. Today
that read is *luck* — it happens when the human drives it, or when the model chooses
to improvise it. v1.0 makes it a **structured, recorded, reusable phase.**

## What the tracer is

On intent approval, the coordinator spawns the tracers: **one read-only child per
repository in the intent's scope, in parallel.** Each reads its repository
first-hand for this change and **reports back a manifest** to the coordinator. It does
not write plans and does not talk to the human — it finds, records, and reports.

Why these properties:

- **Post-approval.** Tracing is expensive and must aim at a *confirmed* target. The
  human's Gate-1 approval is what confirms the coordinator understood the plain ask,
  so the tracer never fires on an unconfirmed guess. (See the flow in
  `intent-and-criteria.md`.)
- **One child per repository.** Each repository is its own codebase with its own
  grounding rules (INV-GROUND-01) and its own eventual worktree. Per-repo children
  keep each child's context bounded and let the reads run concurrently.
- **Read-only, so no lease or worktree.** The lease/worktree machinery exists to
  protect *writers* (`concurrency-and-candidate.md`). The tracer writes nothing, so it
  needs none of it and can fan out widely and cheaply. This mirrors execution's
  per-repo worker fan-out, moved earlier in the lifecycle.
- **Reports back; the coordinator stays the sole planner.** The coordinator authors
  the plan and runs the gates. Letting the tracer write plans directly would scatter
  authorship and bypass the coordinator's feasibility check; the fidelity that a
  direct-writer would gain is preserved instead by making the manifest **rich enough to
  plan from.**

## The manifest — rich enough to plan from

The tracer's report is not a thin summary; the coordinator plans from it without having
read the code itself, so it must carry:

- the **file and call-site map** — exact paths, the specific sites that change,
  relevant signatures and integration points;
- a **proposed task partition** — how the work naturally splits (by feature area, by
  layer), for the coordinator to ratify or adjust;
- **concrete risks** named against the real code — data, security, irreversibility,
  coupling — including risks only a code read reveals (see the worked example);
- the **executable "done" checks** — the runnable commands/tests that *prove the
  intent's outcome criteria hold* (below);
- a **tier signal** — evidence that may raise the provisional tier;
- **open questions** — anything the human must settle, for the coordinator to relay;
- a **completeness proof** where the change has a "change every X" obligation — the
  command and count showing the found set is the *whole* set.

## Executable "done" is a tracing output

The intent's `contract.yaml` states acceptance at the **outcome** level ("no direct
`localStorage` access remains"). The tracer turns each outcome criterion into a
**runnable check** grounded in the real code — e.g. `grep` proves zero direct
`localStorage` call-sites remain outside the hook. That check is carried into the plan
and re-run by the verifier.

This is what fixes the drift seen in the direct session: completeness stops being a
file count the coordinator hand-maintains and becomes a command that either passes or
fails. It is also *why the spec adversary was removed* — you cannot wordsmith an
executable-precise criterion into existence before the code is read; the tracer *earns*
it against the real code. The intent holds the stable, human-level target; the tracer
produces the machine-precise measure.

## Capture and freshness — pay for tracing once

The manifest is **durable**, not ephemeral grep output lost in a context window. It is
recorded as grounding evidence, so:

- a later session does not re-trace from zero — it loads the recorded manifest and
  runs a **bounded freshness check** against current code (have the files it describes
  changed since it was written?), re-reading only what drifted;
- captured design knowledge (`sources/`, `context/`) is reused across sessions exactly
  as it is today — the freshness check is what keeps a reused snapshot from silently
  going stale, the one failure mode of grounding on `sources/` alone.

This is the closed knowledge loop (`closed-knowledge-loop.md`) applied to *code-level*
tracing, not just design docs: a workspace accumulates a living map of itself, and
tracing gets cheaper the more it is used.

## Feedback edges

- **Kick back to the intent.** The tracer usually feeds the plan. When it finds the
  *intent itself* is wrong — infeasible, or larger than approved — it re-opens the
  intent rather than planning the wrong goal in detail. Nothing has executed, so this
  is cheap.
- **Out-of-scope reach.** If the human bound a repository and the tracer finds the work
  must *modify* a repository or area beyond it, the feasibility check surfaces that to
  the human before any plan is written ("doable, but it also needs to change X — include
  it?"). This is not a scope gate — scope-safety is settled at delivery (Gate 2); it is
  the feasibility check flagging a reach the human should confirm
  (`intent-feasibility.md`).
- **Tier raise.** A tier signal (auth, secrets, data, irreversibility) raises the
  provisional tier; the raise surfaces at plan review.

## Depth scales with tier

- **Explore** (`cc-pair`, M3): no separate tracer child — the human reads alongside
  the agent live; tracing is the pairing itself.
- **Standard:** the tracer spawns; a proportionate read grounded by a triage pass that
  sizes the change first, so effort matches breadth.
- **Critical:** the deepest read — exhaustive call-site and dependency mapping, with
  completeness proofs required for every "change every X" obligation.

A triage pass inside tracing sizes the change before the deep read, so a one-file
change gets a glance and a cross-cutting change gets a full survey.

## Relationship to the worker's grounding

The worker still reads the repository first-hand at execution (INV-GROUND-01) — but
that read now **confirms** the tracer's findings rather than being the code's first
contact. With the executable completeness check in hand, the worker's job is to
implement and let the check prove nothing was missed: confirmation, not re-tracing.

## Worked example — the risk only the tracer catches

Intent: "anything written to localStorage should be hashed or encoded." The tracer
reads the code and finds the values are **auth tokens**. It reports a risk a
criteria-wording pass could never find: *base64 encoding is trivially reversible and
gives no protection at rest; real security needs encryption and key handling* — which
may be a bigger change than the goal assumed. That reaches the human at plan review
(or kicks back to the intent), before a single line is written. Reading reality, not
polishing wording, is what surfaces it.

## Skill

- **New `cc-trace`** — spawned on intent approval: fan out one read-only tracer
  child per repository, collect the manifests, run the feasibility check on the
  findings, and hand the coordinator a grounded basis for the plan(s). Read-only over
  the code; never edits the contract; never writes plans.
