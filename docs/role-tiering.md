# Per-role model & effort

This adapter surface reduces per-agent inference cost and latency by running the
right `(model, effort)` for each execution role. It is **invariant-neutral by
construction**: model and effort are bounded, provider-neutral host evidence
(INV-HOST-01) and must never live in the runtime (INV-RUNTIME-01). The runtime
stays model-blind; selection is a coordinator/host decision that changes cost and
speed, never meaning. No new invariant governs it.

The unit is a concrete `(model, effort)` pair — a real model id and a real effort
level (`low` … `max`) — not an abstract "tier". Both are latency levers: a
smaller model *and* a lower effort each finish sooner. The user names the actual
values; there is no `cheap`/`top` vocabulary and no hidden mapping.

## Where the config lives — host-local, per-user, grouped by host

Model availability differs per host and provider, and each person's preference
differs, so this config is **host-local and per-user**, in the same category as
`repositories.local.yaml`: it lives at the workspace root, is **gitignored**, is
never part of shipped template or workspace state, and an upgrade neither creates
nor preserves it (`wrapper/manifest.yaml` never-ship boundary). Read it, when
present, from **`role-tiering.local.yaml`** at the workspace root. The
`*.local.yaml` suffix is this workspace's convention for a gitignored, per-user,
host-local file (the same convention as `repositories.local.yaml`).

Because a model id only means something on a host that offers it, the config is
**grouped by host**. Each host group names the models available on that host. One
user who runs different plans on different hosts (or works across several hosts)
keeps a single file with one group per host; the coordinator reads the group
matching the host it is running on — its own `host_evidence` id — and falls back
to the adapter defaults when its host has no group. Shape:

```yaml
# host-local, per-user, gitignored — one group per host you use
hosts:
  <host-id>:                 # the host adapter's own id (the CLI/agent you run)
    worker:
      model: <a-capable-model-on-this-host>
      effort: high
      escalate_on_repair: true
    verifier:
      model: <a-cheaper-model-on-this-host>
      effort: medium
      escalate_on_repair: false   # hard pin — same (model, effort) every attempt
  <another-host-id>:
    worker:   { model: <a-model-on-that-host>, effort: high }
    verifier: { model: <a-model-on-that-host>, effort: medium }
```

- **Grouped by host.** `hosts.<id>` names the models for that host; the same
  workspace run from a different host reads that host's own group. A host with no
  group uses the adapter defaults.
- **Only `worker` and `verifier`** appear in a group. The coordinator is the root
  session the user already controls at the host level (its own model/effort
  controls), so it needs no entry and is never auto-retiered.
- **`escalate_on_repair` is per role.** `false` is a hard pin; `true` makes the
  configured `(model, effort)` a start-and-floor a repair may raise, never lower.
- An unset role — or an unlisted host — falls back to the **adapter default**.

## Adapter defaults (unset host or role)

When the current host has no group, or a role is absent from its group, the
coordinator uses the adapter's default for that host. These are sensible starting
points, not authority; a host that exposes only one model runs everything at that
model.

| Role | Default model | Default effort | escalate_on_repair |
| --- | --- | --- | --- |
| worker | the session's current model | high | true |
| verifier | the session's current model | medium | false |

"The session's current model" means: absent an explicit choice, spawn the child at
whatever model the coordinator session is already running, so the default is
always a model the host actually offers. There is no abstract default to resolve
and no per-host mapping table.

## Escalate on repair

When `escalate_on_repair: true`, the configured `(model, effort)` is the **start
and the floor**; a repair may raise it, never lower it. The signal is free — the
work failed once — so the repair attempt spends more capability:

- Attempt 1 (initial implementation) runs the role at its configured `(model,
  effort)`.
- On each verifier rejection, the repair attempt raises effort and/or model above
  the configured start.

**The adapter owns the escalation ladder above the user's start**, because only
the host knows the ordering of its own model catalog. Effort has a host-neutral
order (`low` < `medium` < `high` < `max`) the coordinator walks directly; model
ordering is adapter-specific — the adapter raises effort first, then steps the
model up its own catalog. `escalate_on_repair: false` disables the ladder
entirely; the pin holds for every attempt.

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

## Applying the tier is host-specific

Reading and recording a `(model, effort)` does not, by itself, change what a
worker or verifier runs — the value must be **applied to the spawn**, and how is a
host-adapter concern:

- The coordinator records the configured `(model, effort)` as host evidence, and
  **the host adapter sets it on the child** when launching the worker/verifier
  (see the host adapter for the concrete mechanism). Recording without setting it
  on the spawn is intent only; the child then silently inherits the coordinator's
  session model.
- A host may support a per-child **model** but not a per-child **effort**. Where
  per-child effort is unavailable, effort is recorded as evidence but not applied,
  and tiering degrades to **model-only** — no block, no error. A host that exposes
  only one model runs everything on it.

## Recording and honesty

- The coordinator records the `(model, effort)` used per attempt as host evidence
  with the runtime `attempt-evidence-record` action; read against the failure
  counter this gives a first-try-pass-rate. The runtime stores it and never
  interprets it. It records no wall-clock — an attempt already carries `started_at`
  and `checked_at` timestamps whose span is its duration.
- The values are never surfaced to a lay user as mechanism (`docs/terminology.md`
  reporting rules); they appear only under explicit diagnostics.
- Selection lives in the coordinator/host, never in `wrapper/runtime/engine.sh`
  (INV-RUNTIME-01). `(model, effort)` authorizes nothing (INV-HOST-01) and never
  changes the failure counter (INV-REPAIR-01) or the independence requirement.
