# Roles and spawning

v1.0 has **four structural agent roles** (three unchanged from Context Circuit
today, plus discovery, new): discovery is spawned automatically, per repository,
in parallel, on intent approval — it is a real read-only child, not a
human-invoked advisory step. This file states the cast, when each is spawned, and
the independence rules — once, explicitly, so no other file has to reconstruct it.

## The four structural roles

| Role | New? | What it does | Spawned as a child? | Writes product files? |
| --- | --- | --- | --- | --- |
| **Coordinator** | kept | Root conversational session; interprets, routes, runs the gates, delivers briefs | No — it *is* the session | Never |
| **Discovery** | **new** | Reads one repository first-hand for the approved intent and reports a manifest back to the coordinator | Yes (host child), one per repository, in parallel, spawned automatically on intent approval | No (read-only over the repository; no lease or worktree needed) |
| **Worker** | kept | The one bounded implementer, in an isolated worktree | Yes (host child) | Yes, within scope |
| **Independent verifier** | kept | Read-only check bound to the candidate | Yes (host child) | No (read-only, structurally enforced) |

The coordinator is the root you talk to; discovery, the worker, and the
independent verifier are all spawned as children when the lifecycle calls for
them (Task/subagent on Claude Code, `spawn_agent` on Codex, the equivalent on
Cursor), mapped through the existing host-adapter rule that a native child maps to
a single bounded role and carries only bounded `host_evidence` (INV-HOST-01).
Discovery is spawned **automatically** on intent approval, at Standard and
Critical — it is never optional and never human-invoked; it reports only to the
coordinator, never directly to the human.

## Discovery in depth

Discovery is not advisory and not optional at Standard/Critical — it is the
mandatory reality check that runs *after* Gate 1, once the target is confirmed.
Full mechanism in `discovery-and-grounding.md`; here is its shape as a role:

- **Input:** the approved `INTENT.md` + `contract.yaml` (frozen at Gate 1), scoped
  to one repository named in the intent's envelope. Discovery is the first thing
  in the intent phase that reads the real **codebase** — drafting the intent upstream
  reads only existing `context/` Product Knowledge, never the code. Intent grounds in
  existing knowledge; discovery grounds in the code.
- **Task:** read the repository first-hand and build the file/call-site map,
  name concrete risks against the real code (data, security, irreversibility,
  coupling), propose a task partition, turn each outcome-level acceptance
  criterion into an **executable "done" check**, signal whether the provisional
  tier should rise, and produce a completeness proof for any "change every X"
  obligation.
- **Output shape:** a rich, durable **manifest** reported back to the
  coordinator — not a verdict, and not a file the human reads directly. It is
  recorded as grounding evidence and reused (with a bounded freshness check)
  rather than redone from zero next time.
- **What it does not do:** it does not write plans, does not edit the frozen
  contract, and does not talk to the human. It can, however, **kick back to the
  intent** (if the intent itself is wrong) or trigger the **envelope re-gate**
  (if the change reaches outside the approved scope) — both routed through the
  coordinator, never direct to the human.
- **Independence (structural, per repository):** one child per repository, in
  parallel, each read-only. Because it writes nothing, it needs no lease and no
  worktree — the machinery that protects writers is simply not needed here.

## Who spawns, by tier

The number of children is not fixed — it scales with consequence (M3).

| Tier | Discovery | Worker | Independent verifier | Children spawned |
| --- | :--: | :--: | :--: | --- |
| **Explore** (pairing) | **none** — the human reads alongside the agent live (`cc-pair`) | 1 | **none** | **1** — worker only |
| **Standard** | one child per repository, proportionate depth | 1 | 1 (candidate-bound) | up to N, across the lifecycle (discovery fan-out + worker + verifier) |
| **Critical** | one child per repository, exhaustive depth, completeness proofs required | 1 | 1 (candidate-bound) | up to N, across the lifecycle (discovery fan-out + worker + verifier) |

They do not all run at once even at Standard/Critical: discovery runs at
**intent-approval time** (after Gate 1, before any plan or worktree), the worker
at **execution time**, the verifier **after each candidate**.

## Does Explore spawn an agent? Yes — one worker

Explore is `cc-pair` folded in (M3). The coordinator interprets and delegates but
**never writes** (INV-PAIR-01 preserved); it spawns **one worker child** that writes
in the isolated `cc-pair/<session>` branch/worktree. It spawns **no verifier** and
**no discovery child** — at Explore the human reads the real code alongside the
agent live, so a separate discovery read would be redundant
(`discovery-and-grounding.md`, "Depth scales with tier"). So Explore = coordinator
(root) + 1 worker child, and none of them is a verifier — the direct answer to "do
I need a verifier for this?" at the bottom tier is "no," structurally.

If the host cannot create even that one worker child (or the isolated worktree), the
outcome is `host-blocked` and read-only — unchanged from `cc-pair` today.

## Promotion adds agents as stakes rise

Promoting an Explore session to Standard/Critical is the moment new roles can
appear: attaching an intent and raising the tier spawns discovery (per repository,
in parallel) on approval, and the resulting plan(s) spawn the **independent
verifier** on the candidate. The ramp adds agents with consequence, replacing
today's cliff between "pairing (worker only)" and "plan (always a verifier)."

## Independence rules (what must be distinct)

The structurally load-bearing separations are preserved and extended:

- **Worker ≠ verifier** — the verifier is a distinct child, read-only, never
  self-verifies (INV-VERIFY-02, kept). This is enforced by the runtime.
- **Discovery is read-only and reports only to the coordinator** — it never
  writes plans, never edits the frozen contract, and never talks to the human
  directly; the coordinator is always the one relaying its findings, kicking back
  to the intent, or triggering an envelope re-gate.
- Everything else (which model, which host) is bounded `host_evidence` and
  authorizes nothing — a bigger model for the worker never grants broader scope or
  makes a check less independent (INV-HOST-01, kept).

Per-role `(model, effort)` tiering from `role-tiering.local.yaml` still applies and
is orthogonal to the consequence tier: consequence tier decides *which roles run*;
model/effort tiering decides *how each spawned role is powered*.
