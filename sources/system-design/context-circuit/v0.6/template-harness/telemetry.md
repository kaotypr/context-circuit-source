# Template harness v0.6 — telemetry

## Where the numbers come from

No new instrumentation is needed. The Claude Code host runner, driven with
`--output-format stream-json`, ends each `claude -p` turn with a `result` event
that already carries the usage for that turn:

- `num_turns` — agent-loop turns in that invocation (tool-use rounds, including
  the coordinator's spawned worker/verifier sub-agents);
- `usage.output_tokens` — tokens the model **generated**;
- `usage.cache_read_input_tokens` — tokens **re-read** from the prompt cache (the
  size of the growing conversation/context each turn);
- `total_cost_usd` — the fullest aggregate cost, including sub-agents.

The driver already accumulates the full stream to `coordinator-stream.jsonl` and
already computes a per-turn **action tag** (the same tag it stamps on the
file-access trace). Telemetry reuses both: after each coordinator turn, parse that
turn's `result` event and fold its numbers into the action's running totals.

The human-simulator's verdict call is a **separate** `claude -p` written to its
own file, never appended to `coordinator-stream.jsonl`, so it never contaminates
the coordinator's ledger.

## The `$CC_TELEMETRY` file

A TSV the driver writes and `grade.sh` reads. v0.5 defined a three-column row
`action⇥turns⇥tokens`; v0.6 extends it **additively** so the legacy shape stays
readable:

```
# action  conv_turns  agent_turns  output_tokens  context_peak  cost_usd
execute-plan   1   90   37008   5188343   8.5009
```

- `conv_turns` — how many human↔coordinator exchanges mapped to this action.
- `agent_turns` — summed `num_turns` (the agent-loop count).
- `output_tokens` — summed `usage.output_tokens` (the token **budget** unit).
- `context_peak` — the **max** `cache_read_input_tokens` seen for the action (how
  large the context window grew); summing would double-count, so peak is used.
- `cost_usd` — summed `total_cost_usd`.

One row per distinct action tag. `grade.sh` maps the first two/three columns for
backward compatibility and uses the rest when present.

## Per-action attribution

Attribution uses the driver's existing per-turn action tag — the same best-effort
mechanism dimension C already relies on. It is deliberately coarse: a turn's whole
usage is charged to that turn's action. Because attribution is heuristic (a turn
can do more than its label suggests), the ledger is reported, never gated (v0.5
§4.2; this scope keeps it soft). If finer attribution is ever needed, the v0.5
"optional efficiency-analyst" note already covers an out-of-loop, advisory pass —
it remains out of the pass/fail path.

## Boundaries

- Deterministic parse of the runner's own output; **no model call**, no new probe,
  no change to how the coordinator runs.
- Host-specific: only the `claude-code` driver emits this today. Other hosts (or a
  run with no usable result usage) simply omit `$CC_TELEMETRY`, and dimension D
  degrades to "unavailable" exactly as before.
- Written under the disposable per-run directory, never the source `.runtime/`.
