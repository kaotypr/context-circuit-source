# Roles and spawning

vNext has **four agent roles** (Context Circuit today has three). This file states
the cast, when each is spawned, and the independence rules — once, explicitly, so no
other file has to reconstruct it.

## The four roles

| Role | New? | What it does | Spawned as a child? | Writes product files? |
| --- | --- | --- | --- | --- |
| **Coordinator** | kept | Root conversational session; interprets, routes, runs the gates, delivers briefs | No — it *is* the session | Never |
| **Worker** | kept | The one bounded implementer, in an isolated worktree | Yes (host child) | Yes, within scope |
| **Independent verifier** | kept | Read-only check bound to the candidate | Yes (host child) | No (read-only, structurally enforced) |
| **Spec adversary** | **new** | Attacks the acceptance criteria before approval | Yes (host child) | No (read-only over the intent contract) |

The coordinator is the root you talk to; the other three are spawned children
(Task/subagent on Claude Code, `spawn_agent` on Codex, the equivalent on Cursor),
mapped through the existing host-adapter rule that a native child maps to a single
bounded role and carries only bounded `host_evidence` (INV-HOST-01).

## The spec adversary in depth

The adversary is the highest-leverage new role, so its contract is worth stating:

- **Input:** the intent's `contract.yaml` (goal, non-goals, constraints,
  acceptance criteria, scope). It does **not** see any implementation — it runs
  before code exists.
- **Task:** produce, in `adversary.md`, (a) any way to *satisfy every criterion and
  still be wrong*, and (b) the missing edge / error / security / concurrency /
  data-loss paths the criteria do not cover.
- **Output shape:** a list of findings, each `{severity, statement, suggested
  criterion}`, plus a verdict `criteria_sound: yes | needs-work`. Findings become
  new/revised criteria or explicit open questions before approval.
- **Independence:** a distinct child from whatever authored the criteria; it never
  writes the contract itself (it proposes; the human decides at the gate).
- **Calibration (anti-theatre):** the adversary is measured by how often its
  findings actually change a criterion. A low change-rate means it is fabricating
  implausible cases and is miscalibrated — see `risks-and-open-questions.md`. Its
  depth is tiered (below).

## Who spawns, by tier

The number of children is not fixed at four — it scales with consequence (M3). This
is the whole point of tiering.

| Tier | Spec adversary | Worker | Independent verifier | Children spawned |
| --- | :--: | :--: | :--: | --- |
| **Explore** (pairing) | skipped (or one light inline pass) | 1 | **none** | **1** — the worker only |
| **Standard** | 1 (on the criteria) | 1 | 1 (candidate-bound) | up to 3, across the lifecycle |
| **Critical** | 1 (full battery) | 1 | 1 (candidate-bound) | up to 3, across the lifecycle |

They do not all run at once even at Standard/Critical: the adversary fires at
**intent time** (before any worktree), the worker at **execution time**, the
verifier **after each candidate**.

## Does Explore spawn an agent? Yes — one worker

Explore is `cc-pair` folded in (M3). The coordinator interprets and delegates but
**never writes** (INV-PAIR-01 preserved); it spawns **one worker child** that writes
in the isolated `cc-pair/<session>` branch/worktree. It spawns **no verifier** and
**no adversary**. So Explore = coordinator (root) + 1 worker child, and none of them
is a verifier — the direct answer to "do I need a verifier for this?" at the bottom
tier is "no," structurally.

If the host cannot create even that one worker child (or the isolated worktree), the
outcome is `host-blocked` and read-only — unchanged from `cc-pair` today.

## Promotion adds agents as stakes rise

Promoting an Explore session to Standard/Critical is the moment new roles appear:
attaching an intent runs the **adversary**; raising the tier spawns the
**independent verifier** on the resulting candidate. The ramp adds agents with
consequence, replacing today's cliff between "pairing (worker only)" and "plan
(always a verifier)."

## Independence rules (what must be distinct)

The one structurally load-bearing separation is preserved and extended:

- **Worker ≠ verifier** — the verifier is a distinct child, read-only, never
  self-verifies (INV-VERIFY-02, kept).
- **Adversary ≠ the criteria author** — the adversary must not be the same actor
  that wrote the criteria it attacks; otherwise it is self-review of the spec.
- Everything else (which model, which host) is bounded `host_evidence` and
  authorizes nothing — a bigger model for the worker never grants broader scope or
  makes a check less independent (INV-HOST-01, kept).

Per-role `(model, effort)` tiering from `role-tiering.local.yaml` still applies and
is orthogonal to the consequence tier: consequence tier decides *which roles run*;
model/effort tiering decides *how each spawned role is powered*.
