# Per-role model & effort (v0.7.0)

This adapter surface reduces per-agent inference cost and latency by running the
right `(model, effort)` for each execution role. It is **invariant-neutral by
construction**: model and effort are bounded, provider-neutral host evidence
(INV-HOST-01) and must never live in the runtime (INV-RUNTIME-01). The runtime
stays model-blind; selection is a coordinator/host decision that changes cost and
speed, never meaning. No new invariant governs it.

The unit is a concrete `(model, effort)` pair — a real model id and a real effort
level (`low` … `max`) — not an abstract "tier". Both are latency levers: a
smaller model *and* a lower effort each finish sooner. The user names the actual
values; there is no `cheap`/`top` vocabulary and no hidden per-host mapping.

## Where the config lives — host-local, never shipped

Model ids are host/provider-specific, and each person's preference differs, so
this config is **host-local and per-user**, in the same category as
`repositories.local.yaml`: it lives at the workspace root, is **gitignored**, is
never part of shipped template or workspace state, never travels with the
workspace, and an upgrade neither creates nor preserves it
(`wrapper/manifest.yaml` never-ship boundary). Each person configures it once for
the host they work on; a different host or teammate has its own file, or just the
adapter defaults below. It is not portable precisely because it holds
host-specific values — unlike the workspace's logical repository identity
(INV-REPO-01), which is portable because it holds none.

Read it, when present, from **`role-tiering.local.yaml` at the workspace root**.
The `*.local.yaml` suffix is this workspace's convention for a gitignored,
per-user, host-local file (the same convention as `repositories.local.yaml`).
Shape:

```yaml
# host-local — NOT portable workspace state, NOT shipped
worker:
  model: opus
  effort: high
  escalate_on_repair: true
verifier:
  model: sonnet
  effort: medium
  escalate_on_repair: false   # hard pin — same (model, effort) every attempt
```

- **Only `worker` and `verifier` appear.** The coordinator is the root session
  the user already controls at the host level (Claude Code `/model`, `/fast`,
  effort), so it needs no entry and is never auto-retiered.
- **`escalate_on_repair` is per role.** `false` is a hard pin; `true` makes the
  configured `(model, effort)` a start-and-floor a repair may raise, never lower.
- A role that is unset falls back to the **adapter default** below.

## Adapter defaults (unset roles)

When a role is absent from the host-local config, the coordinator uses the
adapter's concrete default for this host. These are sensible starting points, not
authority; a host that exposes only one model runs everything at that model.

| Role | Default model | Default effort | escalate_on_repair |
| --- | --- | --- | --- |
| worker | the session's current model | high | true |
| verifier | the session's current model | medium | false |

"The session's current model" means: absent an explicit host-local choice, spawn
the child at whatever model the coordinator session is already running, so the
default is always a model the host actually offers. There is no abstract default
to resolve and no per-host mapping table.

## Escalate on repair

When `escalate_on_repair: true`, the configured `(model, effort)` is the **start
and the floor**; a repair may raise it, never lower it. The signal is free — the
work failed once — so the repair attempt spends more capability:

- Attempt 1 (initial implementation) runs the role at its configured `(model,
  effort)`.
- On each verifier rejection, the repair attempt raises effort and/or model
  above the configured start.

**The adapter owns the escalation ladder above the user's start**, because only
the host knows the ordering of its own model catalog. Effort has a host-neutral
order (`low` < `medium` < `high` < `max`) the coordinator walks directly; model
ordering is adapter-specific — for the Claude Code adapter, raise effort first,
then step the model up the catalog (…→ sonnet → opus). `escalate_on_repair:
false` disables the ladder entirely; the pin holds for every attempt.

### It must not touch the accounting

INV-REPAIR-01 is unchanged: the worker-failure counter increments on **every**
verifier rejection including the initial implementation, and three still stops
execution. Escalation changes *which model runs attempt N*, never *what a
rejection costs*. A rejection at the configured start is a real failure and
counts; escalation never buys extra attempts. Each repair is still a new commit
(INV-EXEC-04).

### The honest consequence of a hard pin

If a user pins the worker to a small model + low effort with escalation off, the
three-failure limit may be reached more often, because the system can no longer
add capability on repair. That is the user's tradeoff to own, not a bug — but the
report must make it visible in plain language: "stopped after three attempts; the
plan was pinned to a fixed setting, so no extra capability was added on repair."
The system respects the pin and tells the truth about its cost; it never silently
escalates past a pin to rescue the run.

## Optional per-plan complexity hint

Escalate-on-repair reacts *after* one failed attempt. When a plan is known hard
up front, `cc-plan` may record an optional `complexity: high` on the plan
(`plan.yaml`, additive, default absent) that nudges the worker's start **one step
above** the configured `(model, effort)` for attempt 1, then escalates as usual.
It is a hint, not a gate: absent, behavior is exactly the config above; present,
it only shifts the starting point up. A hard pin (`escalate_on_repair: false`)
ignores the hint.

## Verifier independence is preserved

Independence is defined by role and access, not model (INV-VERIFY-01/02): one
independent verifier, read-only on committed state, never repairing or
self-verifying. A verifier on a smaller model is still a **separate agent
invocation** inspecting the worker's committed result, so it remains independent.
A user may even set the **same** `(model, effort)` for worker and verifier — that
does not break independence, because independence is separate-agent-reading-
committed-state, never different-model. What cost pressure must never do is
collapse the worker and verifier into one invocation, or let a verifier reuse the
worker's context. The rule is: separate agent always; `(model, effort)` is free
to vary.

## Recording and honesty

- The coordinator records the `(model, effort)` used per attempt as host evidence
  with the runtime `attempt-evidence-record` action (it pairs naturally with
  `worker_wall_s` / `verifier_wall_s`), enabling a first-try-pass-rate read.
  The runtime stores it and never interprets it.
- The values are never surfaced to a lay user as mechanism (`docs/terminology.md`
  reporting rules); they appear only under explicit diagnostics.
- Selection lives in the coordinator/host, never in `wrapper/runtime/engine.sh`
  (INV-RUNTIME-01). `(model, effort)` authorizes nothing (INV-HOST-01) and never
  changes the failure counter (INV-REPAIR-01) or the independence requirement.
